package main

import (
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/monkjunior/goweb.learn/controllers"
	"github.com/monkjunior/goweb.learn/middleware"
	"github.com/monkjunior/goweb.learn/models"
	"github.com/monkjunior/goweb.learn/rand"
	"net/http"
)

func main() {
	cfg := DefaultConfig()
	dbConfig := DefaultPostgresConfig()
	service, err := models.NewServices(dbConfig.ConnectionInfo())
	if err != nil {
		panic(err)
	}
	defer service.Close()
	service.AutoMigrate()

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(service.User)
	galleriesC := controllers.NewGalleries(service.Gallery, service.Image, *r)

	authKey, err := rand.Bytes(32)
	if err != nil {
		panic(err)
	}

	csrfMw := csrf.Protect(authKey, csrf.Secure(cfg.IsProd()))
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

	// Assets
	assetsHandler := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assetsHandler))

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

	fmt.Printf("Starting server on port %v\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMw(userMw.Apply(r)))
}
