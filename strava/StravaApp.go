package strava

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
	"time"
)

type ConfigToken struct {
	CurrentRefreshToken string `json:"currentRefreshToken"`
	CurrentAccessToken  string `json:"currentAccesToken"`
	ClientId            string `json:"clientId"`
	ClientSecret        string `json:"clientSecret"`
}

type StravaPageData struct {
	PageTitle        string
	AccessToken      string
	LastUpdated      string
	ActivitiesList   []StravaListActivitiesResponse
	TotalKmFormatted string
	TotalSession     int
}

func ParseConfig(configPath string) (ConfigToken, error) {
	jsonFile, err := os.Open(configPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		return ConfigToken{}, err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return ConfigToken{}, err
	}
	var config ConfigToken
	json.Unmarshal(byteValue, &config)
	return config, nil
}

func GetNewTokenFromRefreshToken(refreshToken string, clientId string, clientSecret string) (StravaRefreshToAccessTokenEndpointResponse, error) {
	// JSON body
	body := []byte(fmt.Sprintf(
		StravaRefreshToAccessTokenEndpointBody,
		clientId,
		clientSecret,
		refreshToken,
	))
	resp, err := http.Post(StravaRefreshToAccessTokenEndpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return StravaRefreshToAccessTokenEndpointResponse{}, err
	}
	defer resp.Body.Close()
	parse := &StravaRefreshToAccessTokenEndpointResponse{}
	err = json.NewDecoder(resp.Body).Decode(parse)
	if err != nil {
		return StravaRefreshToAccessTokenEndpointResponse{}, err
	}

	return *parse, nil
}

func GetActivitiesList(accessToken string) ([]StravaListActivitiesResponse, error) {
	t := time.Now()
	now := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location())
	before := now.AddDate(0, 0, -30)
	endpoint := fmt.Sprintf(StravaListActivitiesEndpoint, before.Unix())
	fmt.Print(endpoint)
	req, _ := http.NewRequest("GET",
		endpoint,
		nil)
	bearer := "Bearer " + accessToken
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var parse []StravaListActivitiesResponse
	err = json.Unmarshal(body, &parse)
	if err != nil {
		return nil, err
	}
	for i, j := 0, len(parse)-1; i < j; i, j = i+1, j-1 {
		parse[i], parse[j] = parse[j], parse[i]
	}
	for k, v := range parse {
		// fmt.Print(v)
		parse[k].StartDateFormatted = v.StartDateLocal.Format("Mon, 02 Jan 2006 15:04:05")
		parse[k].DistanceKm = v.Distance / 1000
		parse[k].DistanceKmStr = fmt.Sprintf("%.2f", parse[k].DistanceKm)

		parse[k].Pace = float32(v.ElapsedTime) / parse[k].DistanceKm
		parse[k].PaceReMins = int(parse[k].Pace) / 60
		parse[k].PaceReSecs = int(parse[k].Pace) - parse[k].PaceReMins*60

		parse[k].ElapsedTimeReHours = v.ElapsedTime / 3600
		parse[k].ElapsedTimeReMins = (v.ElapsedTime - parse[k].ElapsedTimeReHours*3600) / 60
		parse[k].ElapsedTimeReSecs = (v.ElapsedTime - parse[k].ElapsedTimeReHours*3600 - parse[k].ElapsedTimeReMins*60)
		parse[k].MovingTimeReHours = v.MovingTime / 3600
		parse[k].MovingTimeReMins = (v.MovingTime - parse[k].MovingTimeReHours*3600) / 60
		parse[k].MovingTimeReSecs = (v.MovingTime - parse[k].MovingTimeReHours*3600 - parse[k].MovingTimeReMins*60)

		parse[k].DateFormatted = v.StartDateLocal.Format("Mon, 02 Jan 2006")
		parse[k].HoursFormatted = v.StartDateLocal.Format("15:04:05")
	}
	return parse, nil
}

func Cal30DaysStats(data []StravaListActivitiesResponse) (float32, int) {
	var totalKm float32 = 0
	totalSession := len(data)
	for _, v := range data {
		totalKm += v.DistanceKm
	}
	return totalKm, totalSession
}

func GenerateStaticHTML() {
	accessKeys, err := ParseConfig("configs/realAccessKeys.json")

	if err != nil {
		panic(err)
	}
	resp, err := GetNewTokenFromRefreshToken(accessKeys.CurrentRefreshToken, accessKeys.ClientId, accessKeys.ClientSecret)
	if err != nil {
		fmt.Print(err)
	}
	token := resp.AccessToken

	// update token file
	accessKeys.CurrentAccessToken = token
	accessKeys.CurrentRefreshToken = resp.RefreshToken
	file, _ := json.MarshalIndent(accessKeys, "", " ")
	_ = ioutil.WriteFile("configs/realAccessKeys.json", file, 0644)
	activitiesList, err := GetActivitiesList(token)
	if err != nil {
		fmt.Print(err)
	}
	// fmt.Print(activitiesList[3])
	totalKm, totalSession := Cal30DaysStats(activitiesList)

	PlotTrend(activitiesList)

	// gen data
	data := StravaPageData{
		PageTitle:        "My last 30 days Strava",
		AccessToken:      token,
		LastUpdated:      time.Now().Format("January 02, 2006 15:04:05"),
		ActivitiesList:   activitiesList,
		TotalKmFormatted: fmt.Sprintf("%.2f", totalKm),
		TotalSession:     totalSession,
	}

	// generate template
	f, _ := os.Create("web/html/index.html")
	t := template.Must(template.ParseFiles("web/html/layout.html"))
	t.Execute(f, data)
	f.Close()
}

func DoGenerateStaticHTML() {
	// each 10 minutes
	for range time.Tick(time.Second * 10 * 60) {
		// do the interval task
		GenerateStaticHTML()
	}
}
