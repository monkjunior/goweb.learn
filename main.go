package main

import (
	"fmt"
	"net/http"
)

func myHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	switch r.URL.Path {
	case "/":
		fmt.Fprint(w, "<h1>Building an awesome web app!</h1>")
	case "/dragon":
		fmt.Fprint(w, "<h1>This is a dragon!</h1>")
	case "/hydra":
		fmt.Fprint(w, "<h1>This is a hydra!</h1>")
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>Oop! 404 page!</h1>")
	}
}

func main() {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", myHandlerFunc)
	http.ListenAndServe(":8080", mux)
}
