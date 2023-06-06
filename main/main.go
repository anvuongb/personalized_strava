package main

import (
	"fmt"
	"net/http"
	"stravapersonal/strava"
)

func main() {

	go strava.DoGenerateStaticHTML()

	// tmpl := template.Must(template.ParseFiles("web/html/layout.html"))
	http.Handle("/strava/web/", http.StripPrefix("/strava/web/", http.FileServer(http.Dir("web"))))
	http.Handle("/strava/", http.StripPrefix("/strava/", http.FileServer(http.Dir("web/html"))))
	// http.HandleFunc("/strava", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl.Execute(w, data)
	// })
	fmt.Printf("Started HTTP server")
	http.ListenAndServe(":3000", nil)
}
