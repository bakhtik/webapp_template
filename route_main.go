package main

import (
	"net/http"

	"github.com/bakhtik/webapp_template/data"
)

func index(w http.ResponseWriter, req *http.Request) {
	if sess, err := session(w, req); err != nil {
		generateHTML(w, nil, "layout", "public.navbar", "index")
	} else {
		user, err := sess.User()
		if err != nil {
			logger.SetPrefix("ERROR ")
			logger.Println(err, "Cannot fetch user")
		}
		data := struct {
			data.User
		}{user}
		generateHTML(w, data, "layout", "private.navbar", "index")
	}
}
