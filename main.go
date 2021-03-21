package main

import (
	"net/http"

	routes "watch_go/router"
)

func main() {

	router := routes.NewRouter()
	http.ListenAndServe(":8554", router)

}
