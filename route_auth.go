package main

import (
	"net/http"
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
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot parse form")
	}
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

// GET /profile
// Show the profile page
func profile(w http.ResponseWriter, req *http.Request) {
	sess, _ := session(w, req)
	user, err := sess.User()
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot fetch user")
	}
	generateHTML(w, user, "layout", "private.navbar", "profile")
}

// POST /change_account
// changes user account (password)
func changeAccount(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot parse form")
	}

	user, err := data.UserByEmail(req.PostFormValue("email"))
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot find user")
		http.Error(w, "Cannot find user", http.StatusForbidden)
		return
	}

	// check if old password matches with existing one
	// does the entered password match the stored password?
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.PostFormValue("old_password"))); err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Old passwords invalid")
		http.Error(w, "Old password invalid", http.StatusForbidden)
		return
	}

	newPassword, confirmPassword := req.PostFormValue("new_password"), req.PostFormValue("confirm_password")
	// check if provided passwords are the same
	if newPassword != confirmPassword {
		logger.SetPrefix("WARNING ")
		logger.Printf("User %s: confirm password mismatch with new password", user.Name)
		http.Error(w, "New passwords must match", http.StatusForbidden)
		return
	}

	// generate hash for the provided password
	bs, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.MinCost)
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot generate hash for new password")
		http.Error(w, "Cannot generate hash for new password", http.StatusForbidden)
		return
	}
	// store new password in the database
	user.Password = string(bs)
	err = user.Update()
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot update user in the database")
		http.Error(w, "Cannot generate hash for new password", http.StatusForbidden)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
