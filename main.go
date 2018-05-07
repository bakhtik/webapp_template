package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Webapp template", version(), "started at", config.Address)
	mux := http.NewServeMux()

	// handle static assets
	files := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	mux.Handle("/", logged(http.HandlerFunc(index)))
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/signup_account", signupAccount)
	mux.HandleFunc("/authenticate", authenticate)
	mux.Handle("/logout", authenticated(http.HandlerFunc(logout)))
	mux.Handle("/admin", authenticated(authorized(http.HandlerFunc(admin), "admin")))

	log.Fatal(http.ListenAndServe(config.Address, mux))
}
