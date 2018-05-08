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

	handlers := map[string]http.Handler{
		"/":               http.HandlerFunc(index),
		"/favicon.ico":    http.NotFoundHandler(),
		"/login":          http.HandlerFunc(login),
		"/signup":         http.HandlerFunc(signup),
		"/signup_account": http.HandlerFunc(signupAccount),
		"/authenticate":   http.HandlerFunc(authenticate),
		"/logout":         authenticated(http.HandlerFunc(logout)),
		"/profile":        authenticated(http.HandlerFunc(profile)),
		"/change_account": authenticated(http.HandlerFunc(changeAccount)),
		"/admin":          authenticated(authorized(http.HandlerFunc(admin), "admin")),
	}

	for pattern, handler := range handlers {
		mux.Handle(pattern, logged(handler))
	}

	log.Fatal(http.ListenAndServe(config.Address, mux))
}
