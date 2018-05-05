package main

import (
	"fmt"
	"net/http"

	"github.com/bakhtik/pob/data"
)

// Check if the user is logged in and has a session, if not err is not nil
func session(r *http.Request) (sess data.Session, err error) {
	cookie, err := r.Cookie("session")
	if err == nil {
		sess = data.Session{Uuid: cookie.Value}
		if err = sess.Check(); err != nil {
			err = fmt.Errorf("Invalid session: %s", err)
		}
	}
	return
}
