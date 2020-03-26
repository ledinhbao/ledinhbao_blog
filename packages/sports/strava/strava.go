package strava

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

const (
	stravaAuthURL = "https://www.strava.com/oauth/authorize"
	tokenURL      = "https://www.strava.com/oauth/token"
	clientID      = "44814"
	clientSecret  = "c44a13c4308b3b834320ae5e3648d6c7855980a3"
	revokeURL     = "https://www.strava.com/oauth/deauthorize"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
}

type Athlete struct {
	gorm.Model
	Profile       string `json:"profile"`
	ProfileMedium string `json:"profile_medium" mapstructure:"profile_medium"`
	Sex           string `json:"sex"`
	State         string `json:"state"`
	Username      string `json:"username"`
	Country       string `json:"country"`
	City          string `json:"city"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
}

type Link struct {
	gorm.Model
	UserID       uint
	Username     string
	AccessToken  string `mapstructure:"access_token"`
	RefreshToken string `mapstructure:"refresh_token"`
	ExpiresAt    int    `mapstructure:"expires_at"`
	ExpiresIn    int    `mapstructure:"expires_in"`
	TokenType    string `mapstructure:"token_type"`
}

// TableName return table's name for strava's link records.
func (Link) TableName() string {
	return "strava_links"
}

// InitializeRoutes inits routes with <prefix>/strava/*
func InitializeRoutes(engine *gin.Engine, prefix string) {
	stravaRoute := engine.Group(prefix + "/strava")
	{
		stravaRoute.GET("/", stravaExchangeToken)
		stravaRoute.GET("/revoke/:username", stravaRevokeToken)
	}
}

// GetAuthURL return authorize link to strava: ?client_id=<>&
func (config Config) GetAuthURL() string {
	res := fmt.Sprintf("https://www.strava.com/oauth/authorize?client_id=%s", config.ClientID)
	res += fmt.Sprintf("&redirect_uri=%s", config.RedirectURI)
	res += fmt.Sprintf("&response_type=code&approval_prompt=auto&scope=%s", strings.Join(config.Scopes, ","))
	return res
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
		c.Redirect(http.StatusFound, "/admin/dashboard")
	}

}

func stravaRevokeToken(c *gin.Context) {
	// TODO Valid data where username is linked with user
	db := c.MustGet("database").(*gorm.DB)
	username := c.Param("username")
	link := Link{}
	_ = db.Where(Link{Username: username}).First(&link)

	client := &http.Client{}
	urlValues := url.Values{}
	urlValues.Set("access_token", link.AccessToken)

	request, _ := http.NewRequest("POST", revokeURL, strings.NewReader(urlValues.Encode()))
	response, _ := client.Do(request)

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		// Remove record from database
		go removeStravaRecord(db, username)
	}

	c.Redirect(http.StatusFound, "/admin/dashboard")
}

func removeStravaRecord(db *gorm.DB, username string) {
	tx := db.Begin()
	tx.Unscoped().Delete(Link{Username: username})
	tx.Unscoped().Delete(Athlete{Username: username})
	tx.Commit()
}
