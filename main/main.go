package main

import (
	"fmt"
	"html/template"
	"net/http"
	"stravapersonal/strava"
)

type StravaPageData struct {
	PageTitle        string
	AccessToken      string
	ActivitiesList   []strava.StravaListActivitiesResponse
	TotalKmFormatted string
	TotalSession     int
}

func main() {
	accessKeys, err := strava.ParseConfig("configs/accessKeys.json")
	if err != nil {
		panic(err)
	}
	token, err := strava.GetNewTokenFromRefreshToken(accessKeys.CurrentRefreshToken, accessKeys.ClientId, accessKeys.ClientSecret)
	if err != nil {
		fmt.Print(err)
	}
	activitiesList, err := strava.GetActivitiesList(token)
	if err != nil {
		fmt.Print(err)
	}
	// fmt.Print(activitiesList[3])
	totalKm, totalSession := strava.Cal30DaysStats(activitiesList)

	strava.PlotHistogram(activitiesList)

	tmpl := template.Must(template.ParseFiles("web/html/layout.html"))
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := StravaPageData{
			PageTitle:        "My last 30 days Strava",
			AccessToken:      token,
			ActivitiesList:   activitiesList,
			TotalKmFormatted: fmt.Sprintf("%.2f", totalKm),
			TotalSession:     totalSession,
		}
		tmpl.Execute(w, data)
	})
	fmt.Printf("Started HTTP server")
	http.ListenAndServe(":3000", nil)
}
