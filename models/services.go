package models

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewServices(connInfo string) (*Services, error) {
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

	return &Services{
		db:   db,
		User: NewUserService(db),
	}, nil
}

type Services struct {
	db      *gorm.DB
	Gallery GalleryService
	User    UserService
}

// Closes the database connection.
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
	return s.db.AutoMigrate(&User{}, &Gallery{})
}

// AutoMigrate will attempt to automatically migrate all table
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{})
}
