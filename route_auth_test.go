package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetLogin(t *testing.T) {
	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()
	login(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
	if strings.Contains(string(body), "Sign in") == false {
		t.Errorf("Body does not contain Sign in")
	}
}

func TestGetSignup(t *testing.T) {
	req := httptest.NewRequest("GET", "/signup", nil)
	w := httptest.NewRecorder()
	signup(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
	want := "Sign up for the account below"
	if strings.Contains(string(body), want) == false {
		t.Errorf("Body does not contain %q", want)
	}
}

func TestSignupAccount(t *testing.T) {
	req := httptest.NewRequest("POST", "/signup_account", nil)
	req.ParseForm()
	req.PostForm.Add("name", "John Doe")
	req.PostForm.Add("email", "john_doe@gmail.com")
	req.PostForm.Add("password", "123")
	req.PostForm.Add("role", "user")
	w := httptest.NewRecorder()
	signupAccount(w, req)

	resp := w.Result()
	// body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}
