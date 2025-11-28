package routeur

import (
	"fmt"
	"net/http"
	"text/template"
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
	data, err := controller.Damso()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Erreur: %v\n", err)
		return
	}

	// Debug: afficher les données reçues
	fmt.Printf("[DEBUG] Damso data: %+v\n", data)

	tmpl, err := template.ParseFiles("template/damso.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Erreur template: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

func laylowHandler(w http.ResponseWriter, r *http.Request) {
	data, err := controller.Laylow()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Erreur: %v\n", err)
		return
	}

	tmpl, err := template.ParseFiles("template/laylow.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Erreur template: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/index.html")
}
