package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Structure pour décoder la réponse de Spotify pour une track
type Track struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Duration int    `json:"duration_ms"`
	Artists  []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Album struct {
		Name   string `json:"name"`
		Images []struct {
			URL string `json:"url"`
		} `json:"images"`
	} `json:"album"`
}

// LaylowMaladresse pour compatibilité
type LaylowMaladresse struct {
	Name     string
	Duration int
	Artists  []struct {
		Name string
	}
	Album struct {
		Name   string
		Images []struct {
			URL string
		}
	}
}

func Laylow() (*Track, error) {
	urlAPI := "https://api.spotify.com/v1/tracks/67Pf31pl0PfjBfUmvYNDCL"

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

	var track Track
	err := json.Unmarshal(body, &track)
	if err != nil {
		fmt.Println("Erreur décodage JSON :", err.Error())
		return nil, err
	}

	fmt.Printf("[DEBUG] Track parsée: %s\n", track.Name)

	return &track, nil
}
