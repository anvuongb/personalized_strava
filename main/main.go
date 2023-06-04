package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
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
	accessKeys, err := strava.ParseConfig("configs/realAccessKeys.json")

	if err != nil {
		panic(err)
	}
	resp, err := strava.GetNewTokenFromRefreshToken(accessKeys.CurrentRefreshToken, accessKeys.ClientId, accessKeys.ClientSecret)
	if err != nil {
		fmt.Print(err)
	}
	token := resp.AccessToken

	// update token file
	accessKeys.CurrentAccessToken = token
	accessKeys.CurrentRefreshToken = resp.RefreshToken
	file, _ := json.MarshalIndent(accessKeys, "", " ")
	_ = ioutil.WriteFile("configs/realAccessKeys.json", file, 0644)
	activitiesList, err := strava.GetActivitiesList(token)
	if err != nil {
		fmt.Print(err)
	}
	// fmt.Print(activitiesList[3])
	totalKm, totalSession := strava.Cal30DaysStats(activitiesList)

	strava.PlotTrend(activitiesList)

	tmpl := template.Must(template.ParseFiles("web/html/layout.html"))
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	http.HandleFunc("/strava", func(w http.ResponseWriter, r *http.Request) {
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
