package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

var (
	currentToken string
	tokenMutex   sync.RWMutex
)

// InitTokenRefresh obtient un token Spotify via Client Credentials
func InitTokenRefresh() {
	RefreshSpotifyToken()

	// Renouveler automatiquement toutes les heures
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			RefreshSpotifyToken()
		}
	}()
}

// RefreshSpotifyToken obtient un nouveau token via Client Credentials
func RefreshSpotifyToken() {
	clientID := os.Getenv("clientID")
	clientSecret := os.Getenv("clientSecret")

	if clientID == "" || clientSecret == "" {
		fmt.Println("[Erreur] clientID ou clientSecret manquants")
		return
	}

	// Encodage Base64 : clientID:clientSecret
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Requête POST pour obtenir le token
	reqBody := "grant_type=client_credentials"
	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", bytes.NewBufferString(reqBody))
	if err != nil {
		fmt.Println("[Erreur] Création requête token :", err.Error())
		return
	}

	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	httpClient := http.Client{Timeout: time.Second * 5}
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("[Erreur] Appel endpoint token :", err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		fmt.Printf("[Erreur] Status %d: %s\n", res.StatusCode, string(body))
		return
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResp); err != nil {
		fmt.Println("[Erreur] Décodage réponse token :", err.Error())
		return
	}

	// Sauvegarder le token
	tokenMutex.Lock()
	currentToken = tokenResp.AccessToken
	tokenMutex.Unlock()

	// Définir comme variable d'environnement pour les autres contrôleurs
	os.Setenv("SPOTIFY_TOKEN", tokenResp.AccessToken)

	fmt.Printf("[✓] Token Spotify obtenu (expire dans %d secondes)\n", tokenResp.ExpiresIn)
}

// GetToken retourne le token actuel
func GetToken() string {
	tokenMutex.RLock()
	defer tokenMutex.RUnlock()
	return currentToken
}
