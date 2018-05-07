package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/bakhtik/webapp_template/data"
	"golang.org/x/crypto/bcrypt"
)

// GET /login
// Show the login page
func login(w http.ResponseWriter, req *http.Request) {
	if _, err := session(w, req); err == nil { // user already logged in
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	generateHTML(w, nil, "layout", "public.navbar", "login")
}

// GET /signup
// Show the signup page
func signup(w http.ResponseWriter, req *http.Request) {
	if _, err := session(w, req); err == nil { // user already logged in
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	generateHTML(w, nil, "layout", "public.navbar", "signup")
}

// POST /singup_account
// Create the user account
func signupAccount(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot parse form")
	}
	user := data.User{
		Name:     req.PostFormValue("name"),
		Email:    req.PostFormValue("email"),
		Password: req.PostFormValue("password"),
		Role:     req.PostFormValue("role"),
	}
	if err = user.Create(); err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot create user")
	}
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

// POST /authenticate
// Authenticate the user given the email and password
func authenticate(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	user, err := data.UserByEmail(req.PostFormValue("email"))
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot find user")
	}

	// does the entered password match the stored password?
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.PostFormValue("password"))); err == nil {
		session, err := user.CreateSession()
		if err != nil {
			logger.SetPrefix("ERROR ")
			logger.Println(err, "Cannot create session")
		}
		cookie := http.Cookie{
			Name:     "session",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, req, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
	}
}

// GET /logout
// Logs the user out
func logout(w http.ResponseWriter, req *http.Request) {

	sess, err := session(w, req)
	// delete the session
	if err = sess.DeleteByUUID(); err != nil {
		logger.SetPrefix("WARNING ")
		logger.Println(err, "Failed to delete sesssion")
	}
	// remove the cookie
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	// Clean up sessions
	if time.Now().Sub(sessionsCleaned) > (time.Second * time.Duration(config.SessionLength)) {
		go data.CleanSessions(config.SessionLength)
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

// for authorized access only to handlers
func authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// check if authenticated
		_, err := session(w, req)
		if err != nil {
			//http.Error(w, "not logged in", http.StatusUnauthorized)
			logger.SetPrefix("WARNING ")
			logger.Println(err, `Failed to get/verify cookie "session"`)
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return // don't call original handler
		}
		next.ServeHTTP(w, req)
	})
}

// permission check
func authorized(next http.Handler, roles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		sess, _ := session(w, req)
		if roles != nil {
			user, err := sess.User()
			if !strSliceContains(roles, user.Role) {
				logger.SetPrefix("WARNING ")
				logger.Printf("%v: User %s has not permission for requested page", err, user.Name)
				http.Error(w, "You must have admin rights to enter the page", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, req)
	})
}

// log handler
func logged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Create a response wrapper:

		// switch out response writer for a recorder
		// for all subsequent handlers
		c := httptest.NewRecorder()

		next.ServeHTTP(c, req)

		// copy everything from response recorder
		// to actual response writer
		for k, v := range c.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(c.Code)
		c.Body.WriteTo(w)

		// log
		resp := c.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		// fmt.Print(req.RemoteAddr, req.Method, req.RequestURI, "response: ", w.)
		fmt.Printf("%s %s %d %v %q\n", req.Method, req.RequestURI, resp.StatusCode, resp.Header, string(body))
	})
}
