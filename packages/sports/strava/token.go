package strava

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
