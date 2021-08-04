package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func handleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Welcome to my awsome site!</h1>")
	} else if r.URL.Path == "/contact" {
		fmt.Fprint(w, "To get in touch, please send an email to <a href=\"mailto:vungocson998@gmail.com\">vungocson998@gmail.com</a>.")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>404</h1>")
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleFunc)
	r.HandleFunc("/contact", handleFunc)
	http.ListenAndServe(":8080", r)
}
