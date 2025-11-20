package controller

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ptet à la place de mettre deux fois les id les déclarer avant (demander a kavtiv)
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, _ := template.ParseFiles("templates/" + tmpl)
	t.Execute(w, data)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index.html", nil)
}

func DamsoHandler(w http.ResponseWriter, r *http.Request) {
	token := GetSpotifyToken("a0d5f1a538194548941f281c8982b9aa", "2f0ba65c5c69454ab44063a026586600")

	artistData := SpotifyGET("https://api.spotify.com/v1/artists/2UwqpfQtNuhBwviIC0f2ie", token)
	RenderTemplate(w, "damso.html", string(artistData))
}

func LaylowHandler(w http.ResponseWriter, r *http.Request) {
	token := GetSpotifyToken("a0d5f1a538194548941f281c8982b9aa", "2f0ba65c5c69454ab44063a026586600")

	artistData := SpotifyGET("https://api.spotify.com/v1/artists/4PN6gPp0vcdQkW7UzDxX2E", token)
	RenderTemplate(w, "laylow.html", string(artistData))
}

var spotifyToken string

func GetSpotifyToken(clientID, clientSecret string) string {
	if spotifyToken != "" {
		return spotifyToken
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	basicAuth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+basicAuth)

	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	spotifyToken = result["access_token"].(string)
	return spotifyToken
}

func SpotifyGET(url string, token string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body
}
