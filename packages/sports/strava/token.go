package strava

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// StravaToken represents a Token object with Strava Server
type StravaToken struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func stravaSendRefreshToken(oldRefreshToken string) (StravaToken, error) {
	urlData := createBaseURLData()
	urlData.Set("grant_type", "refresh_token")
	urlData.Set("refresh_token", oldRefreshToken)
	req, _ := http.NewRequest("POST", tokenURL, nil)
	req.URL.RawQuery = urlData.Encode()
	resp, _ := (&http.Client{}).Do(req)
	log.Println("Sent refresh token request, respose status code:", resp.StatusCode)
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		var result StravaToken
		rebody, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(rebody, &result)
		return result, err
	}
	return StravaToken{}, errors.New("Response code is not in range 200...299")
}

func stravaSendRevokeToken(accessToken string) {
	urlData := createBaseURLData()
	urlData.Set("access_token", accessToken)
	req, _ := http.NewRequest("POST", revokeURL, nil)
	req.URL.RawQuery = urlData.Encode()
	resp, _ := (&http.Client{}).Do(req)
	log.Println("Strava > send revoke token > response code", resp.StatusCode)
}

func stravaRevokeToken(c *gin.Context) {
	// TODO Valid data where username is linked with user
	db := getDatabaseInstance(c)
	username := c.Param("username")
	link := Link{}
	_ = db.Where(Link{Username: username}).First(&link)

	// Send http request in another routine.
	go stravaSendRevokeToken(link.AccessToken)
	// Remove database record.
	removeStravaRecord(db, username)

	c.Redirect(http.StatusFound, config.getRedirectPath())
}

func stravaGetAccessTokenForUserID(userID uint, db *gorm.DB) (StravaToken, error) {
	var link Link
	db.Where(Link{UserID: userID}).First(&link)
	if link.ID == 0 {
		return StravaToken{}, fmt.Errorf("Could not find Strava Link, userID %d didn't link Strava to their account", userID)
	}
	return stravaGetAccessTokenForLink(link, db)
}

func stravaGetAccessTokenForAthleteID(athleteID uint64, db *gorm.DB) (StravaToken, error) {
	var link Link
	db.Where(Link{AthleteID: athleteID}).First(&link)
	if link.ID == 0 {
		return StravaToken{}, fmt.Errorf("Cound not find Strava Link for Athlete ID %d", athleteID)
	}
	return stravaGetAccessTokenForLink(link, db)
}

func stravaGetAccessTokenForLink(link Link, db *gorm.DB) (StravaToken, error) {
	if link.ID == 0 {
		return StravaToken{}, errors.New("Cannot get access token for empty Link Object")
	}
	if time.Now().Add(10*time.Minute).Unix() > int64(link.ExpiresAt) {
		token, err := stravaSendRefreshToken(link.RefreshToken)
		if err == nil {
			link.AccessToken = token.AccessToken
			link.RefreshToken = token.RefreshToken
			link.ExpiresAt = token.ExpiresAt
			link.ExpiresIn = token.ExpiresIn
			db.Save(&link)
		}
		return token, nil
	}
	return StravaToken{
		AccessToken:  link.AccessToken,
		RefreshToken: link.RefreshToken,
	}, nil
}

// GetOAuthURL return oath/authorize url with callback_uri=callback
func GetOAuthURL(callback string) string {
	// TODO implement if scope changes
	formatted := "https://www.strava.com/oauth/authorize?client_id=44814&redirect_uri=%s&response_type=code&approval_prompt=auto&scope=activity:read"
	return fmt.Sprintf(formatted, callback)
}
