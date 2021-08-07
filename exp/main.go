package main

import (
	"fmt"

	"github.com/monkjunior/goweb.learn/models"
	"gorm.io/gorm"
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
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()

	// us.DestructiveReset()

	user, err := us.ByID(1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Predefined %v\n", user)

	newUser := models.User{
		Name:  "duong",
		Email: "duongcho@gmail.com",
	}
	err = us.Create(&newUser)
	if err != nil {
		panic(err)
	}
	fmt.Printf("New user: %v\n", newUser)
}
