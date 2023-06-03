package strava

import "time"

const StravaRefreshToAccessTokenEndpoint string = "https://www.strava.com/oauth/token"
const StravaListActivitiesEndpoint string = "https://www.strava.com/api/v3/athlete/activities?page=1&per_page=30"

var StravaRefreshToAccessTokenEndpointBody string = `{
	"client_id": "%s",
	"client_secret": "%s",
	"grant_type": "refresh_token",
	"refresh_token": "%s"
}`

type StravaRefreshToAccessTokenEndpointResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int32  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type StravaListActivitiesResponse struct {
	Distance           float32 `json:"distance"`
	DistanceKm         float32
	DistanceKmStr      string
	Pace               float32
	PaceReMins         int
	PaceReSecs         int
	MovingTime         int `json:"moving_time"`
	MovingTimeReHours  int
	MovingTimeReMins   int
	MovingTimeReSecs   int
	ElapsedTime        int `json:"elapsed_time"`
	ElapsedTimeReHours int
	ElapsedTimeReMins  int
	ElapsedTimeReSecs  int
	StartDate          time.Time `json:"start_date"`
	DateFormatted      string
	HoursFormatted     string
	StartDateLocal     time.Time `json:"start_date_local"`
	StartDateFormatted string
	AverageHeartrate   float32 `json:"average_heartrate"`
	MaxHeartrate       float32 `json:"max_heartrate"`
	SufferScore        float32 `json:"suffer_score"`
}
