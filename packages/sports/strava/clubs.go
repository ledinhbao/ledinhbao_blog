package strava

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	clubApiURL = string("https://www.strava.com/api/v3/clubs")
)

// StravaClub present a club info, reflects of Strava Club
type StravaClub struct {
	gorm.Model
	ClubID            uint   `json:"id"`
	Name              string `json:"name"`
	ProfileMedium     string `json:"profile_medium"`
	SportType         string `json:"sport_type"`
	Owner             bool   `json:"owner"`
	Admin             bool   `json:"admin"`
	URL               string `json:"url"`
	Country           string `json:"country"`
	ProcessingState   int
	ProcessingMessage string
}

func stravaAddClub(c *gin.Context) {
	ginview.HTML(c, http.StatusOK, "admin-strava-add-club", gin.H{})
}

// // SetProcessingState indicate current processing state of StravaClubModel
// // -   1: added, currently fetching info from Strava Server
// // - 100: done
// func (club *StravaClub) SetProcessingState(state int) {
// 	club.processingState = state
// }

// // ProcessingState returns current processing state of this instance.
// func (club *StravaClub) ProcessingState() int {
// 	return club.processingState
// }

func StravaFetchClubInfoByID(clubID uint, userID uint, db *gorm.DB) {
	var club StravaClub
	db.Where(StravaClub{ClubID: clubID}).First(&club)
	StravaFetchClub(club, userID, db)
}

func StravaFetchClub(club StravaClub, userID uint, db *gorm.DB) {
	var err error
	token, err := stravaGetAccessTokenForUserID(userID, db)
	if err != nil {
		log.Println(fmt.Sprintf("Fetching club data: error when retrieving athlete's access token (%s)", err.Error()))
		return
	}

	urlData := createBaseURLData()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", clubApiURL, club.ClubID), nil)
	if err != nil {
		log.Println("error when prepare fetching club data, " + err.Error())
		return
	}
	req.URL.RawQuery = urlData.Encode()
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, _ := (&http.Client{}).Do(req)
	if resp.StatusCode == 200 {
		// Successful
		rebody, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(rebody, &club)
		if err != nil {
			log.Println("Failed to unmarshal resp for club data,", err, string(rebody))
		}
		log.Println("Fetching club data: succeed!")
		club.ProcessingState = 100
		db.Save(&club)
	} else {
		log.Println("Failed to get club data, response code: ", resp.StatusCode)
	}
}
