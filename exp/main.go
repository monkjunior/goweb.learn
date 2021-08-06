package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "ted"
	password = "your-password"
	dbname   = "goweb_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var id int
	err = db.QueryRow(`
		INSERT INTO users(name, email)
		VALUES($1, $2)
		RETURNING id`,
		"Teddy", "ted@goweb.learn").Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println(id)

	// var name string
	// err = db.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(name)
}
