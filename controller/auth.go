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
	tokenExpiry  time.Time
)

// InitTokenRefresh obtient un token et lance une goroutine pour le renouveler toutes les heures
func InitTokenRefresh() {
	// Obtenir le premier token de manière SYNCHRONE (bloquant)
	fmt.Println("[Init] Obtention du token Spotify...")
	RefreshSpotifyToken()

	// Vérifier que le token a bien été obtenu
	if GetToken() == "" {
		fmt.Println("[Erreur] Impossible d'obtenir le token Spotify au démarrage")
		return
	}

	// Lancer la goroutine de renouvellement
	go func() {
		// Renouveler toutes les heures après le démarrage
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			RefreshSpotifyToken()
		}
	}()
}

// RefreshSpotifyToken obtient un nouveau token via Client Credentials
func RefreshSpotifyToken() {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		fmt.Println("Erreur : SPOTIFY_CLIENT_ID ou SPOTIFY_CLIENT_SECRET manquants")
		return
	}

	// Encodage Base64 : clientID:clientSecret
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Préparer la requête POST
	reqBody := "grant_type=client_credentials"
	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", bytes.NewBufferString(reqBody))
	if err != nil {
		fmt.Println("Erreur création requête token :", err.Error())
		return
	}

	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	httpClient := http.Client{Timeout: time.Second * 5}
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Erreur appel endpoint token :", err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		fmt.Printf("Erreur token (status %d): %s\n", res.StatusCode, string(body))
		return
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResp); err != nil {
		fmt.Println("Erreur décodage réponse token :", err.Error())
		return
	}

	// Sauvegarder le token avec expiration
	tokenMutex.Lock()
	currentToken = tokenResp.AccessToken
	tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	tokenMutex.Unlock()

	fmt.Printf("[Token Spotify] Renouvelé à %s (expire à %s)\n", time.Now().Format("15h04"), tokenExpiry.Format("15h04"))
}

// GetToken retourne le token actuel de manière thread-safe
func GetToken() string {
	tokenMutex.RLock()
	defer tokenMutex.RUnlock()
	return currentToken
}
