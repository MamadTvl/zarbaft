package userRepository

import (
	"bankai/models"
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByUserId(userId uint) (*models.User, error)
	DeleteUser(user *models.User) error
}

type userGormRepository struct {
	db *gorm.DB
}

func NewGormUserRepository() UserRepository {
	return &userGormRepository{
		db: getDbConnection(),
	}
}

func (ur *userGormRepository) CreateUser(user *models.User) error {
	return ur.db.Create(user).Error
}

func (ur *userGormRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := ur.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userGormRepository) DeleteUser(user *models.User) error {
	return ur.db.Delete(user).Error
}

func (ur *userGormRepository) GetUserByUserId(userId uint) (*models.User, error) {
	var user models.User
	result := ur.db.First(&user, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func getDbConnection() *gorm.DB {
	//user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	dbURI := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	fmt.Println(dbURI)
	// Connect to the database

	db, err := gorm.Open("mysql", dbURI)
	if err != nil {
		panic(err)
	}

	// Set up connection pool and other configuration options
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Enable logging in development mode
	db.LogMode(true)

	// Migrate the User model to the database (if necessary)
	db.AutoMigrate(&models.User{})

	// Use the db instance to interact with the database in your application
	return db
}
