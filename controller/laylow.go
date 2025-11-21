package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LaylowMaladresse struct {
	Results []struct {
		Name        string `json:"name"`
		ReleaseDate string `json:"release_date"`
		Time        int    `json:"duration_minutes"`
		Picture     string `json:"artwork_url"`
	} `json:"items"`
}

func Laylow() {

	urlAPI := "https://api.spotify.com/v1/tracks/67Pf31pl0PfjBfUmvYNDCL"
	//timeout (facultatif)
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}
	//création de la requête
	req, errReq := http.NewRequest(http.MethodGet, urlAPI, nil)
	if errReq != nil {
		fmt.Println("Oupss une erreur est survenue lors de la création de la requête :", errReq.Error())
		return
	}
	//ajout du token d'authorisation
	token := GetToken()
	if token == "" {
		fmt.Println("SPOTIFY_TOKEN manquant")
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)
	//envoi de la requête
	res, errResp := httpClient.Do(req)
	if errResp != nil {
		fmt.Println("Oupss une erreur est survenue lors de l'appel à l'API :", errResp.Error())
		return
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	//lecture du body de la réponse
	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		fmt.Println("Oupss une erreur est survenue lors de la lecture du body de la réponse :", errBody.Error())
		return
	}

	var laylowMaladresse LaylowMaladresse

	json.Unmarshal(body, &laylowMaladresse)

	fmt.Println(laylowMaladresse.Results[0])
}
