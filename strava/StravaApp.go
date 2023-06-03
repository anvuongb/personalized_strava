package strava

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type ConfigToken struct {
	CurrentRefreshToken string `json:"currentRefreshToken"`
	CurrentAccessToken  string `json:"currentAccesToken"`
	ClientId            string `json:"clientId"`
	ClientSecret        string `json:"clientSecret"`
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
	req, _ := http.NewRequest("GET", StravaListActivitiesEndpoint, nil)
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
