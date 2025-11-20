package main

import (
	"log"
	"tp-API-Spotify/routeur"
)

func main() {
    router := SetupRouter()

    log.Println("http://localhost:8080")
    http.ListenAndServe(":8080", router)
}