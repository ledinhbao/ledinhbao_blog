package strava

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

// Activity contains data for strava's activity. Currently, it presents only distance, moving time, elapsed time and type.
type Activity struct {
	gorm.Model
	ActivityID     uint64    `json:"id"`
	Distance       float64   `json:"distance"`
	MovingTime     uint      `json:"moving_time"`
	ElapsedTime    uint      `json:"elapsed_time"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	StartDate      time.Time `json:"start_date"`
	StartLocalDate time.Time `json:"start_date_local"`
	AthleteID      uint64
}

// TableName return table's name for strava's activity records.
func (Activity) TableName() string {
	return "strava_activities"
}

func GetActivityFromStravaAPIByID(id uint64, accessToken string) (Activity, error) {
	endpoint := fmt.Sprintf("%s/activities/%d", apiURL, id)
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		// Failed here
	}
	log.Println("Get request to", endpoint)
	log.Println("Get activity from strava receive status code", resp.StatusCode)
	rebody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		var activity Activity
		json.Unmarshal(rebody, &activity)
		activity.ActivityID = id
		return activity, nil
	}
	var fault stravaFault
	json.Unmarshal(rebody, &fault)
	log.Println("Receive error", fault)
	return Activity{}, &TokenError{"Expired Token"}
}
