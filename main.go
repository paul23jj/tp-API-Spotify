package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var client_id = ""
var client_secret = ""

type ApiSpotify struct {
	Info struct {

}

func main() {
	// URL de L'API
	urlApi := ""

	// Création de la requête HTTP vers L'API avec initialisation de la methode HTTP, la route et le corps de la requête
	req, errReq := http.NewRequest(http.MethodGet, urlApi, nil)
	if errReq != nil {
		fmt.Println("Oupss, une erreur est survenue lors de la création de la requête HTTP : ", errReq.Error())
	}

	// Ajout d'une métadonnée dans le header, User_Agent permet d'identifier l'application, système.
	req.Header.Add("User-Agent", "Ynov Campus Cours")

	// Execution de la requête HTTP vars L'API
	res, errResp := http.DefaultClient.Do(req)
	if errResp != nil {
		fmt.Println("Oupss, une erreur est survenue lors de l'exécution de la requête HTTP : ", errResp.Error())
		return
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	// Lecture et récupération du corps de la requête HTTP
	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		fmt.Println("Oupss, une erreur est survenue lors de la lecture du corps de la réponse : ", errBody.Error())
	}

	var data ApiSpotify
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
		return
	}

	fmt.Println(data.Results[0].Name)
	fmt.Println(data.Results[0].Origin.Name)
	fmt.Println(data.Results[0].Image)
}
}
