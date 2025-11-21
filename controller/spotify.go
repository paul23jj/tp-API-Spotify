package controller

import (
	"fmt"
	"net/http"
	"time"
)

func Spotify() {
	token := GetToken()
	if token == "" {
		fmt.Println("Token Spotify indisponible")
		return
	}

	httpClient := http.Client{Timeout: time.Second * 3}

	// Endpoint de vérification léger
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/tracks/67Pf31pl0PfjBfUmvYNDCL", nil)
	if err != nil {
		fmt.Println("Erreur création requête de vérification :", err.Error())
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de la vérification du token :", err.Error())
		return
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode == http.StatusUnauthorized {
		fmt.Println("Token invalide : 401 Unauthorized")
		return
	}
	if res.StatusCode >= 400 {
		fmt.Printf("Échec vérification token, status: %d\n", res.StatusCode)
		return
	}

	// Si le token est valide on lance les fonctions Damso() et Laylow()
	Damso()
	Laylow()
}
