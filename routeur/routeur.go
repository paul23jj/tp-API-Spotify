package routeur

import (
	"fmt"
	"net/http"
	"tp-API-Spotify/controller"
)

func InitRoutes() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	tpl := http.FileServer(http.Dir("template"))
	http.Handle("/template/", http.StripPrefix("/template/", tpl))

	http.HandleFunc("/spotify", spotifyHandler)
	http.HandleFunc("/damso", damsoHandler)
	http.HandleFunc("/laylow", laylowHandler)

	http.HandleFunc("/", indexHandler)
}

func spotifyHandler(w http.ResponseWriter, r *http.Request) {
	controller.Spotify()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Spotify check executed")
}

func damsoHandler(w http.ResponseWriter, r *http.Request) {
	controller.Damso()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Damso executed")
}

func laylowHandler(w http.ResponseWriter, r *http.Request) {
	controller.Laylow()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Laylow executed")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/index.html")
}
