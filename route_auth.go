package main

import (
	"net/http"

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
	if err != nil { // user already logged out
		logger.SetPrefix("WARNING ")
		logger.Println(err, `Failed to get cookie "session"`)
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
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

	http.Redirect(w, req, "/", http.StatusSeeOther)
}
