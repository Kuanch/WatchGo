package controller

import (
	"fmt"
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
	if Verify(w, username, password) {
		loginSession.Values["is_authorized"] = true
		// TODO: to video streaming page
	}
}

func Verify(w http.ResponseWriter, user string, pass string) bool {
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
		fmt.Fprintf(w, "Failed to login\n")
		return false
	} else {
		fmt.Fprintf(w, "User: "+user+" successfully login \n")
		return true
	}
}

func VideoStream(w http.ResponseWriter, r *http.Request) {
	loginSession, _ := store.Get(r, "loginSession")
	if authorized := loginSession.Values["is_authorized"].(bool); authorized {
		fmt.Fprintln(w, "Image!")
		// http.ServeFile(w, r, "./template/html/stream.html")
		// image, _ := getImageFromFilePath("data/image.jpg")
		// services.ResponseWithImage(w, &image)
	} else {
		fmt.Fprintln(w, "Not login yet")
	}
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}
