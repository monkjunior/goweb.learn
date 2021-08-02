package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	switch r.URL.Path {
	default:
		fmt.Fprint(w, "<h1>Building an awesome web app!</h1>")
	case "/dragon":
		fmt.Fprint(w, "<h1>This is a dragon!</h1>")
	case "/hydra":
		fmt.Fprint(w, "<h1>This is a hydra!</h1>")
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":8080", nil) //Use the built-in serve mux
}
