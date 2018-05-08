package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"time"
)

type Configuration struct {
	Address       string
	Static        string
	SessionLength int
	LogFile       string
}

var config Configuration
var logger *log.Logger

func init() {
	loadConfig()
	file, err := os.OpenFile("webapp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

func version() string {
	return "0.1"
}

func strSliceContains(slice []string, value string) (ok bool) {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return
}

// handler for Apache-style logs
func loggingHandler(writer io.Writer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Create a response wrapper:

		// switch out response writer for a recorder
		// for all subsequent handlers
		c := httptest.NewRecorder()

		next.ServeHTTP(c, req)
		// log
		resp := c.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		// copy everything from response recorder
		// to actual response writer
		for k, v := range c.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(c.Code)
		c.Body.WriteTo(w)

		// write log information
		host, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			host = req.RemoteAddr
		}
		username := "-"
		if url.User != nil {
			if name := req.URL.User.Username(); name != "" {
				username = name
			}
		}

		fmt.Fprintf(writer, "%s - %s [%v] \"%s %s %s\" %d %d\n", host, username, time.Now().Format("02/Jan/2006:15:04:05 -0700"), req.Method, req.RequestURI, req.Proto, resp.StatusCode, len(body))
	})
}

// logged handler changes behaviour depending on configuration file
// if none logfile provided no logging occured
// if "stdout" - logging to console else to provied filename
func logged(h http.Handler) http.Handler {
	switch config.LogFile {
	case "":
		return h
	case "stdout":
		return loggingHandler(os.Stdout, h)
	default:
		logFile, err := os.OpenFile(config.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		return loggingHandler(logFile, h)
	}
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
