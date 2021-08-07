package models

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm/logger"
)

func testingUserService() (*UserService, error) {
	var (
		host     = "localhost"
		port     = 5432
		user     = "ted"
		password = "your-password"
		dbname   = "goweb_test"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	us, err := NewUserService(psqlInfo)
	if err != nil {
		return nil, err
	}
	us.db.Logger.LogMode(logger.Silent)
	// Clear the users table between tests
	us.DestructiveReset()
	return us, nil
}

func TestCreateUser(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()

	user := User{
		Name:  "hien",
		Email: "hien@gmail.com",
	}

	err = us.Create(&user)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Errorf("Expected ID > 0. Received %d", user.ID)
	}

	if time.Since(user.CreatedAt) > time.Duration(5*time.Second) {
		t.Errorf("Expected CreatedAt to be recent. Received %s", user.CreatedAt)
	}

	if time.Since(user.UpdatedAt) > time.Duration(5*time.Second) {
		t.Errorf("Expected UpdatedAt to be recent. Received %s", user.UpdatedAt)
	}
}
