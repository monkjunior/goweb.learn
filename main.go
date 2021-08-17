package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/monkjunior/goweb.learn/controllers"
	"github.com/monkjunior/goweb.learn/models"
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
	service, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	// TODO
	// defer us.Close()
	// us.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(service.User)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/login", usersC.GetLogin).Methods("GET")
	r.HandleFunc("/login", usersC.PostLogin).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", r)
}
