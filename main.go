package main

import (
	"net/http"

	routes "go_rest/router"
)


func main() {

	router := routes.NewRouter()
	http.ListenAndServe(":3000", router)

}
