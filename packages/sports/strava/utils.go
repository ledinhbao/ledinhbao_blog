package strava

import (
	"log"
	"net/url"
)

// createBaseURLData create url.Values content client_id and client_secret
func createBaseURLData() url.Values {
	res := url.Values{}
	res.Set("client_id", config.ClientID)
	res.Set("client_secret", config.ClientSecret)
	return res
}

func getCallbackURLOrPanic(panic bool) string {
	if config.URLCallbackHost == "" {
		if panic {
			log.Panic("Failed to obtain URLCallbackHost for Strava Module.")
		}
	}
	return config.URLCallbackHost + config.PathPrefix + "/strava" + config.PathSubscription
}
