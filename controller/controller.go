package controller

import (
	"fmt"
	"os"
	"watch_go/services"

	"github.com/Kuanch/mjpeg"
	"github.com/gorilla/sessions"

	"watch_go/utils"

	"net/http"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var (
	stream *mjpeg.Stream
)

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
		loginSession.Values["is_authorized"] = false
		loginSession = Auth(w, r)
		if err := loginSession.Save(r, w); err != nil {
			fmt.Println(err)
		}

		http.Redirect(w, r, "/api/video_stream", http.StatusFound)
	}
}

func Auth(w http.ResponseWriter, r *http.Request) *sessions.Session {
	loginSession, _ := store.Get(r, "loginSession")
	username := loginSession.Values["username"].(string)
	password := loginSession.Values["password"].(string)
	if utils.Verify(username, password) {
		loginSession.Values["is_authorized"] = true
	} else {
		services.ResponseWithJson(w, http.StatusOK, "Login failed")
	}

	return loginSession
}

func VideoStream(w http.ResponseWriter, r *http.Request) {
	loginSession, _ := store.Get(r, "loginSession")
	deviceID := 0
	if authorized, _ := loginSession.Values["is_authorized"].(bool); authorized {
		stream = mjpeg.NewStream()
		go utils.VideoFeed(deviceID, stream)
		stream.ServeHTTP(w, r)
	} else {
		services.ResponseWithJson(w, http.StatusOK, "Not login yet")
	}
}
