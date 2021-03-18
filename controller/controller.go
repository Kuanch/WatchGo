package controller

import (
	"fmt"
	"go_rest/services"
	"image"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"

	"net/http"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

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
		loginSession = Auth(r)
		if err := loginSession.Save(r, w); err != nil {
			fmt.Println(err)
		}

		http.Redirect(w, r, "/api/video_feed", http.StatusFound)
	}
}

func Auth(r *http.Request) *sessions.Session {
	loginSession, _ := store.Get(r, "loginSession")
	username := loginSession.Values["username"].(string)
	password := loginSession.Values["password"].(string)
	if Verify(username, password) {
		loginSession.Values["is_authorized"] = true
	}

	return loginSession
}

func Verify(user string, pass string) bool {
	// TODO: manage user system with database
	savePassword, readPasswordErr := os.ReadFile(user + ".txt")
	if readPasswordErr != nil {
		log.Fatal(readPasswordErr)
	}

	authPasswordStr := []byte(string(savePassword))
	authPasswordByte, _ := bcrypt.GenerateFromPassword(authPasswordStr, bcrypt.DefaultCost)

	hashCompareErr := bcrypt.CompareHashAndPassword(authPasswordByte, []byte(pass))
	if hashCompareErr != nil {
		fmt.Println(hashCompareErr)
		return false
	}
	return true
}

func VideoStream(w http.ResponseWriter, r *http.Request) {
	loginSession, _ := store.Get(r, "loginSession")
	if authorized, _ := loginSession.Values["is_authorized"].(bool); authorized {
		w.WriteHeader(http.StatusOK)
		image, _ := getImageFromFilePath("data/image.jpg")
		services.ResponseWithImageTemp(w, &image)
	} else {
		services.ResponseWithJson(w, http.StatusOK, "Not login yet")
	}
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)

	return img, err
}
