package strava

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

const (
	stravaAuthURL   = string("https://www.strava.com/oauth/authorize")
	tokenURL        = string("https://www.strava.com/oauth/token")
	clientID        = string("44814")
	clientSecret    = string("c44a13c4308b3b834320ae5e3648d6c7855980a3")
	revokeURL       = string("https://www.strava.com/oauth/deauthorize")
	subscriptionURL = string("https://www.strava.com/api/v3/push_subscriptions")
	apiURL          = string("https://www.strava.com/api/v3")
)

// InitializeRoutes inits routes with <prefix>/strava/*
func InitializeRoutes(engine *gin.Engine) {
	stravaRoute := engine.Group(config.PathPrefix + "/strava")
	{
		stravaRoute.GET("/", stravaExchangeToken)
		stravaRoute.GET("/revoke/:username", stravaRevokeToken)
		stravaRoute.GET("/list/:username", listActivitiesForYesterday)
		stravaRoute.GET(config.PathSubscription, stravaValidateSubscription)
		stravaRoute.GET(config.PathSubscription+"/delete/:subscription-id", stravaDeleteSubscription)
		stravaRoute.GET(config.PathSubscription+"/create", stravaCreateSubscription)

		stravaRoute.POST(config.PathSubscription, stravaSubscriptionHandle)
	}
}

func init() {
	fmt.Println("Called Strava Package init")
	config = &Config{
		ClientID:          "44814",
		ClientSecret:      "c44a13c4308b3b834320ae5e3648d6c7855980a3",
		PathPrefix:        "/admin",
		PathRedirect:      "/dashboard",
		PathSubscription:  "/subscription",
		GlobalDatabase:    "database",
		SubscriptionDBKey: "strava-subscription",
	}
}

func stravaExchangeToken(c *gin.Context) {
	code := c.Request.URL.Query().Get("code")
	if code == "" {
		// access denied
	} else {
		// Exchange for access token
		data := url.Values{}
		data.Set("client_id", clientID)
		data.Set("client_secret", clientSecret)
		data.Set("code", code)
		data.Set("grant_type", "authorization_code")

		client := &http.Client{}
		request, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
		response, err := client.Do(request)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		respData := make(map[string]interface{})
		err = json.Unmarshal(body, &respData)

		athlete := Athlete{}
		if athleteRaw, err := respData["athlete"]; err {
			mapstructure.Decode(athleteRaw, &athlete)
			// _ = json.Unmarshal([]byte(athleteRaw), &athlete)
		}

		session := sessions.Default(c)

		link := Link{}
		_ = mapstructure.Decode(respData, &link)
		link.AthleteID = athlete.AthleteID
		link.Username = athlete.Username
		link.UserID = session.Get("AuthUserID").(uint)

		db := c.MustGet("database").(*gorm.DB)
		// check if link exist
		db.Where("username = ?", link.Username).Delete(Link{})
		db.Where("username = ?", athlete.Username).Delete(Athlete{})
		db.Create(&athlete)
		db.Create(&link)

		// c.JSON(http.StatusOK, gin.H{
		// 	"link":    link,
		// 	"athlete": athlete,
		// })
		c.Redirect(http.StatusFound, config.getRedirectPath())
	}

}

func getDatabaseInstance(c *gin.Context) *gorm.DB {
	return c.MustGet("database").(*gorm.DB)
}

func stravaRevokeToken(c *gin.Context) {
	// TODO Valid data where username is linked with user
	db := getDatabaseInstance(c)
	username := c.Param("username")
	link := Link{}
	_ = db.Where(Link{Username: username}).First(&link)

	client := &http.Client{}
	urlValues := url.Values{}
	urlValues.Set("access_token", link.AccessToken)

	request, _ := http.NewRequest("POST", revokeURL, nil)
	request.URL.RawQuery = urlValues.Encode()
	response, _ := client.Do(request)
	log.Println("Strava > send revoke token > response code", response.StatusCode)
	// if response.StatusCode >= 200 && response.StatusCode <= 299 {
	// Remove record from database
	go removeStravaRecord(db, username)
	// }

	c.Redirect(http.StatusFound, config.getRedirectPath())
}

func removeStravaRecord(db *gorm.DB, username string) {
	tx := db.Begin()
	tx.Unscoped().Delete(Link{Username: username})
	tx.Unscoped().Delete(Athlete{Username: username})
	tx.Commit()
}

func listActivitiesForYesterday(c *gin.Context) {
	db := getDatabaseInstance(c)
	client := &http.Client{}
	var username = c.Param("username")
	var stravaLink = Link{}
	db.Where("username = ?", username).First(&stravaLink)
	var bearer = "Bearer " + stravaLink.AccessToken

	// TODO check for access_token expiration
	var request, _ = http.NewRequest("GET", "https://www.strava.com/api/v3/athlete/activities", nil)
	request.Header.Add("Authorization", bearer)
	resp, _ := client.Do(request)
	var body, _ = ioutil.ReadAll(resp.Body)
	c.JSON(http.StatusOK, gin.H{
		"data": string(body),
	})
}
