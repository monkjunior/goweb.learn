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
		User: NewUserService(db),
	}, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
}
