package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "ted"
	password = "your-password"
	dbname   = "goweb_dev"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;uniqueIndex"`
	Color string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
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
		panic(err)
	}

	// db.Migrator().DropTable(&User{})
	db.AutoMigrate(&User{})

	name, email, color := getInfo()
	user := &User{
		Name:  name,
		Email: email,
		Color: color,
	}

	db.Create(&user)

	fmt.Printf("%v\n", user)
}

func getInfo() (name, email, color string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What is your name?")
	name, _ = reader.ReadString('\n')
	fmt.Println("What is your email address?")
	email, _ = reader.ReadString('\n')
	fmt.Println("What is your favourite color?")
	color, _ = reader.ReadString('\n')

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	color = strings.TrimSpace(color)

	return name, email, color
}
