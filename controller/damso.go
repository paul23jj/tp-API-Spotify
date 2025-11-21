package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DamsoAlbums struct {
	Results []struct {
		Name string `json:"album_type"`
		ListeningTime string `json:"release_date"`
		
	} `json:"items"`
}

func Damso() {

	urlAPI := "https://api.spotify.com/v1/artists/2UwqpfQtNuhBwviIC0f2ie/albums"

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, errReq := http.NewRequest(http.MethodGet, urlAPI, nil)
	if errReq != nil {
		fmt.Println("Oupss une erreur est survenue lors de la création de la requête :", errReq.Error())
		return
	}

	req.Header.Add("a modif", "modif")

	res, errResp := httpClient.Do(req)
	if errResp != nil {
		fmt.Println("Oupss une erreur est survenue lors de l'appel à l'API :", errResp.Error())
		return
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		fmt.Println("Oupss une erreur est survenue lors de la lecture du body de la réponse :", errBody.Error())
		return
	}

	var damsoAlbums DamsoAlbums

	json.Unmarshal(body, &damsoAlbums)

	fmt.Println(damsoAlbums.Results[0])
}
