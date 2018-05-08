package main

import (
	"net/http"

	"github.com/bakhtik/webapp_template/data"
	"golang.org/x/crypto/bcrypt"
)

func admin(w http.ResponseWriter, req *http.Request) {
	sess, _ := session(w, req)
	user, err := sess.User()
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot fetch user")
	}
	users, err := data.Users()
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot fetch users")
	}
	data := struct {
		data.User
		Users []data.User
	}{
		user,
		users,
	}
	generateHTML(w, data, "layout", "private.navbar", "admin")

}

// POST /change_account_admin
// changes user account (password)
func changeAccountAdmin(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot parse form")
	}

	user, err := data.UserByEmail(req.PostFormValue("origin_email"))
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot find user")
		http.Error(w, "Cannot find user", http.StatusForbidden)
		return
	}

	user.Name = req.PostFormValue("name")
	user.Email = req.PostFormValue("email")
	user.Role = req.PostFormValue("role")

	newPassword, confirmPassword := req.PostFormValue("new_password"), req.PostFormValue("confirm_password")
	if newPassword != "" {
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
	}
	// update user
	err = user.Update()
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot update user in the database")
		http.Error(w, "Cannot generate hash for new password", http.StatusForbidden)
		return
	}
	http.Redirect(w, req, "/admin", http.StatusSeeOther)
}

// user delete
func deleteUser(w http.ResponseWriter, req *http.Request) {
	user, err := data.UserByEmail(req.FormValue("email"))
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot find user")
		http.Error(w, "Cannot find user", http.StatusForbidden)
		return
	}
	err = user.Delete()
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot delete user %s", user.Name)
		http.Error(w, "Cannot delete user", http.StatusForbidden)
		return
	}
	http.Redirect(w, req, "/admin", http.StatusSeeOther)
}

// for updating users profiles (resetting passwords)
func profileAdmin(w http.ResponseWriter, req *http.Request) {
	sess, _ := session(w, req)
	admin, err := sess.User()
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot fetch user")
	}

	err = req.ParseForm()
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot parse form")
	}

	user, err := data.UserByEmail(req.FormValue("email"))
	if err != nil {
		logger.SetPrefix("ERROR ")
		logger.Println(err, "Cannot find user")
		http.Error(w, "Cannot find user", http.StatusForbidden)
		return
	}

	data := struct {
		data.User
		UserForUpdate data.User
	}{
		admin,
		user,
	}
	generateHTML(w, data, "layout", "private.navbar", "profile_admin")
}
