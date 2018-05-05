package main

import "net/http"

func index(w http.ResponseWriter, req *http.Request) {
	if _, err := session(req); err != nil {
		generateHTML(w, nil, "layout", "public.navbar", "index")
	} else {
		generateHTML(w, nil, "layout", "private.navbar", "index")
	}
}
