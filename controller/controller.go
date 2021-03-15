package controller

import (
	"fmt"
	"os"

	"github.com/gorilla/sessions"

	"net/http"
)

var SESSION_KEY = "Test"
var store = sessions.NewCookieStore([]byte(os.Getenv(SESSION_KEY)))

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "./template/html/login.html")
	} else if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Println(err)
		}
		username := r.Form["username"][0]
		password := r.Form["password"][0]
		loginSession, _ := store.Get(r, "loginSession")
		loginSession.Values["username"] = username
		loginSession.Values["password"] = password
		if err := loginSession.Save(r, w); err != nil {
			fmt.Println(err)
		}

		http.Redirect(w, r, "/api/auth", http.StatusFound)
	}
}

func Auth(w http.ResponseWriter, r *http.Request) {
	loginSession, _ := store.Get(r, "loginSession")
	username := loginSession.Values["username"].(string)
	password := loginSession.Values["password"].(string)
	if Verify(username, password) {
		// http.Redirect(w, r, /api/video)
		fmt.Fprintf(w, "username: "+username+"\n")
		fmt.Fprintf(w, "password: "+password+"\n")
		fmt.Fprint(w, "Login Success")
	}
}

func Verify(user string, pass string) bool {
	return true
}
