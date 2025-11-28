package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
	"tp-API-Spotify/controller"
	"tp-API-Spotify/routeur"
)

func main() {
	// Charger les variables du fichier .env
	loadEnv()

	routeur.InitRoutes()
	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		log.Println("Impossible de lire .env, utilisation des variables système")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}

	// Obtenir le token avec les credentials du .env
	controller.InitTokenRefresh()
}
