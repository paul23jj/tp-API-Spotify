package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Structure pour un album
type Album struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
	TotalTracks int    `json:"total_tracks"`
	Images      []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
}

// DamsoAlbums contient la liste des albums
type DamsoAlbums struct {
	Results []Album `json:"items"`
}

func Damso() (*DamsoAlbums, error) {
	urlAPI := "https://api.spotify.com/v1/artists/2UwqpfQtNuhBwviIC0f2ie/albums?limit=50"

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

	var damsoAlbums DamsoAlbums
	err := json.Unmarshal(body, &damsoAlbums)
	if err != nil {
		fmt.Println("Erreur décodage JSON :", err.Error())
		fmt.Println("Réponse reçue:", string(body[:min(500, len(body))]))
		return nil, err
	}

	fmt.Printf("[✓] %d albums de Damso récupérés\n", len(damsoAlbums.Results))

	return &damsoAlbums, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
