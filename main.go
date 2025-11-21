package main

import (
	"log"
	"net/http"
	"tp-API-Spotify/controller"
	"tp-API-Spotify/routeur"
)

func main() {
	controller.InitTokenRefresh()

	routeur.InitRoutes()
	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
