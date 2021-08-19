package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/monkjunior/goweb.learn/controllers"
	"github.com/monkjunior/goweb.learn/middleware"
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
	defer service.Close()
	service.AutoMigrate()

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(service.User)
	galleriesC := controllers.NewGalleries(service.Gallery, service.Image, *r)

	userMw := middleware.User{
		UserService: service.User,
	}
	requireUserMw := middleware.RequireUser{User: userMw}

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", requireUserMw.Apply(staticC.Contact)).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/login", usersC.GetLogin).Methods("GET")
	r.HandleFunc("/login", usersC.PostLogin).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	//Image route
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	//Gallery route
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")
	r.HandleFunc("/galleries/new", requireUserMw.ApplyFn(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries/new", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", requireUserMw.ApplyFn(galleriesC.Show)).Methods("GET").Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.GetUpdate)).Methods("GET").Name(controllers.UpdateGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.PostUpdate)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMw.ApplyFn(galleriesC.ImageDelete)).Methods("POST")

	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", userMw.Apply(r))
}
