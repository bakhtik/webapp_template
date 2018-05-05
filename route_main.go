package main

import "net/http"

func index(w http.ResponseWriter, req *http.Request) {
	if sess, err := session(req); err != nil {
		generateHTML(w, nil, "layout", "public.navbar", "index")
	} else {
		user, err := sess.User()
		if err != nil {
			logger.SetPrefix("ERROR ")
			logger.Println(err, "Cannot fetch user")
		}
		generateHTML(w, user, "layout", "private.navbar", "index")
	}
}
