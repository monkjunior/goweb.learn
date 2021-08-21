package models

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Services struct {
	db      *gorm.DB
	Gallery GalleryService
	User    UserService
	Image   ImageService
}

type ServicesConfig func(services *Services) error

func WithGorm(connInfo string) ServicesConfig {
	return func(s *Services) error {
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
			return err
		}
		s.db = db
		return nil
	}
}

func WithUser(hmacKey, pepper string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, hmacKey, pepper)
		return nil
	}
}

func WithGallery() ServicesConfig {
	return func(s *Services) error {
		s.Gallery = NewGalleryService(s.db)
		return nil
	}
}

func WithImage() ServicesConfig {
	return func(s *Services) error {
		s.Image = NewImageService()
		return nil
	}
}

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

// Close the database connection.
func (s *Services) Close() error {
	gDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return gDB.Close()
}

// DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.Migrator().DropTable(&User{}, &Gallery{})
	if err != nil {
		return err
	}
	return s.db.AutoMigrate(&User{}, &Gallery{}, &pwReset{})
}

// AutoMigrate will attempt to automatically migrate all table
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}, &pwReset{})
}
