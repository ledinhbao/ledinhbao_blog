package strava

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ledinhbao/blog/core"
)

func getSubscriptionToken() string {
	return config.ClientID + "hahaha" + config.ClientSecret
}

func sendDeleteSubscription(subID string) {
	client := &http.Client{}
	urlData := createBaseURLData()
	req, _ := http.NewRequest("DELETE", subscriptionURL+"/"+subID, nil)
	req.URL.RawQuery = urlData.Encode()
	_, _ = client.Do(req)

	log.Println("Strava Delete Subscription URL sent:", req.URL.String())
}

// ViewSubscription send POST request to Strava server to get Subscription ID for your application
func ViewSubscription() string {
	client := &http.Client{}
	urlData := createBaseURLData()
	req, _ := http.NewRequest("GET", subscriptionURL, nil)
	req.URL.RawQuery = urlData.Encode()
	log.Println("View Subscription URL sent:", req.URL.String())
	resp, _ := client.Do(req)
	var jsonBody []map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &jsonBody)

	if len(jsonBody) > 0 && resp.StatusCode == 200 {
		res := fmt.Sprintf("%v", jsonBody[0]["id"])
		log.Println("Subscription exist at Strava endpoint with ID: ", res, " ###")
		return res
	}
	log.Println("Subscription does not exist at Strava endpoint.")
	return ""
}

func sendSubscriptionCreationRequest() (int, map[string]interface{}) {
	client := &http.Client{}
	urlData := createBaseURLData()
	urlData.Set("verify_token", getSubscriptionToken())
	urlData.Set("callback_url", getCallbackURLOrPanic(true))
	req, _ := http.NewRequest("POST", subscriptionURL, nil)
	req.URL.RawQuery = urlData.Encode()
	// resp, _ := client.Post(subscriptionURL, "text/plain; charset=utf-8", strings.NewReader(urlData.Encode()))
	resp, _ := client.Do(req)
	body := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&body)
	log.Println("Strava Subscription Creation Request, URL:", body)
	return resp.StatusCode, body
}

func stravaDeleteSubscription(c *gin.Context) {
	subscriptionID := c.Param("subscription-id")
	sendDeleteSubscription(subscriptionID)
	db := c.MustGet(config.GlobalDatabase).(*gorm.DB)
	db.Where(core.Setting{Value: subscriptionID}).Delete(&core.Setting{})
	c.Redirect(http.StatusFound, config.getRedirectPath())
}

func stravaCreateSubscription(c *gin.Context) {
	db := c.MustGet(config.GlobalDatabase).(*gorm.DB)
	CreateSubscription(db)
	c.Redirect(http.StatusFound, config.getRedirectPath())
}

// CreateSubscription kiểm tra trong bảng table config, kiểm tra SubscriptionDBKey có tồn tại hay không,
// nếu không có thì gởi POST request tới server Strava để đăng ký subscription.
// Cập nhật lại trường SubscriptionDBKey khi nhận dữ liệu về.
//
// Yêu cầu:
//     - CALLED AFTER SERVER RUN
//     - Package ledinhbao/core
//     - Bảng settings (key string, value string) trong database.
func CreateSubscription(db *gorm.DB) {
	log.Println("start... strava.CreateSubscription")
	setting := core.Setting{}
	notFoundSeting := db.Where(core.Setting{Key: config.SubscriptionDBKey}).First(&setting).RecordNotFound()

	if notFoundSeting || setting.Value == "" {
		// Un-sync between app and strava
		subscriptionID := ViewSubscription()
		if subscriptionID != "" {
			sendDeleteSubscription(subscriptionID)
		}

		// Send POST request and create database setting
		respCode, jsonBody := sendSubscriptionCreationRequest()
		if respCode == 201 {
			if subscriptionID, err := jsonBody["id"]; err {
				// Save subscription_id to database
				db.Where(
					core.Setting{Key: config.SubscriptionDBKey},
				).Assign(
					core.Setting{Value: fmt.Sprintf("%v", subscriptionID)},
				).FirstOrCreate(&setting)
			}
		} else {
			log.Panic("Error to create subscription with Strava")
		}
	}
}

func stravaValidateSubscription(c *gin.Context) {
	challenge := c.Query("hub.challenge")
	queryToken := c.Query("hub.verify_token")
	subscriptionToken := getSubscriptionToken()
	if queryToken == subscriptionToken {
		c.JSON(http.StatusOK, gin.H{
			"hub.challenge": challenge,
		})
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"query.token":    queryToken,
			"token.verified": subscriptionToken,
			"challenge":      challenge,
		})
	}
}
