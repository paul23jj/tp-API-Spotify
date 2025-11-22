package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Structure pour décoder la réponse complète de Spotify
type SpotifyResponse struct {
	Items []Album `json:"items"`
}

// Structure pour un album
type Album struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
	TotalTracks int    `json:"total_tracks"`
	Images      []struct {
		URL string `json:"url"`
	} `json:"images"`
}

// DamsoAlbums pour compatibilité avec le template
type DamsoAlbums struct {
	Results []Album
}

func Damso() (*DamsoAlbums, error) {
	urlAPI := "https://api.spotify.com/v1/artists/2UwqpfQtNuhBwviIC0f2ie/albums"
	
	httpClient := http.Client{
		Timeout: time.Second * 5,
	}
	
	req, errReq := http.NewRequest(http.MethodGet, urlAPI, nil)
	if errReq != nil {
		fmt.Println("Erreur création requête :", errReq.Error())
		return nil, errReq
	}
	
	token := os.Getenv("SPOTIFY_TOKEN")
	if token == "" {
		fmt.Println("SPOTIFY_TOKEN manquant")
		return nil, fmt.Errorf("SPOTIFY_TOKEN manquant")
	}
	
	req.Header.Add("Authorization", "Bearer "+token)
	
	res, errResp := httpClient.Do(req)
	if errResp != nil {
		fmt.Println("Erreur appel API :", errResp.Error())
		return nil, errResp
	}
	defer res.Body.Close()
	
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		fmt.Printf("Erreur API (status %d): %s\n", res.StatusCode, string(body))
		return nil, fmt.Errorf("erreur API status %d", res.StatusCode)
	}
	
	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		fmt.Println("Erreur lecture body :", errBody.Error())
		return nil, errBody
	}

	fmt.Println("[DEBUG] Réponse brute Spotify:", string(body[:min(500, len(body))]))
	
	var spotifyResp SpotifyResponse
	err := json.Unmarshal(body, &spotifyResp)
	if err != nil {
		fmt.Println("Erreur décodage JSON :", err.Error())
		return nil, err
	}

	fmt.Printf("[DEBUG] Nombre d'albums parsés: %d\n", len(spotifyResp.Items))
	
	return &DamsoAlbums{
		Results: spotifyResp.Items,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

