package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/monkjunior/goweb.learn/controllers"
	"github.com/monkjunior/goweb.learn/email"
	"github.com/monkjunior/goweb.learn/middleware"
	"github.com/monkjunior/goweb.learn/models"
	"github.com/monkjunior/goweb.learn/rand"
	"log"
	"net/http"
	"os"
)

func main() {
	boolPtr := flag.Bool("prod", false, "Set to true in production. This ensures that a .config file is provided before the application start")
	flag.Parse()
	cfg := LoadConfig(*boolPtr)
	service, err := models.NewServices(
		models.WithGorm(cfg.Database.ConnectionInfo()),
		models.WithUser(cfg.HMACKey, cfg.Pepper),
		models.WithGallery(),
		models.WithImage(),
	)
	if err != nil {
		panic(err)
	}
	defer service.Close()

	err = service.AutoMigrate()
	if err != nil {
		panic(err)
	}

	mailgunCfg := cfg.Mailgun
	emailer := email.NewClient(
		email.WithSender("Goweb.learn support", "support@"+mailgunCfg.Domain),
		email.WithMailgun(mailgunCfg.Domain, mailgunCfg.ApiKey),
	)
	_ = emailer

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(service.User, emailer)
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
	r.HandleFunc("/logout", requireUserMw.ApplyFn(usersC.Logout)).Methods("POST")
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

	log.Printf("Starting server on port %v\n", cfg.Port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMw(userMw.Apply(r))))
}

func LoadConfig(configReq bool) Config {
	if !configReq {
		return DefaultConfig()
	}
	f, err := os.Open(".config")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	var c Config
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&c)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println("Successfully loaded .config file")
	return c
}
