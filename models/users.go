package models

import (
	"errors"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErrNotFound = errors.New("models: resource not found")
)

func NewUserService(connInfo string) (*UserService, error) {
	db, err := gorm.Open(postgres.Open(connInfo), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Disable color
			},
		),
	})
	if err != nil {
		return nil, err
	}

	return &UserService{
		db: db,
	}, nil
}

type UserService struct {
	db *gorm.DB
}

// ByID will look up by the provided ID.
// If the user is found, we will return a nil error.
// If the user is not found, we will return ErrNotFound.
// If there is another error, we might return an error
// that is not generated by the model package.
//
// As a general rule, any error but ErrNotFound should
// probably result in a 500 error.
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Create will create the provided user and backfill data
// like the ID, CreatedAt, and UpdatedAt fields.
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// Closes the UserService database connection.
func (us *UserService) Close() error {
	gDB, err := us.db.DB()
	if err != nil {
		return err
	}
	return gDB.Close()
}

// DestructiveReset drops the user table and rebuilds it
func (us *UserService) DestructiveReset() {
	us.db.Migrator().DropTable(&User{})
	us.db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;uniqueIndex"`
}