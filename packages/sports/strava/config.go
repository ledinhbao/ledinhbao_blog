package strava

import (
	"fmt"
	"strings"

	"github.com/imdario/mergo"
)

// Config struct contains configuration for Strava Package
//   - SubscriptionDBKey	: strava subcsription key, stored in database, indicates exist subscription if present.
type Config struct {
	ClientID          string
	ClientSecret      string
	Scopes            []string
	PathPrefix        string
	PathRedirect      string
	PathSubscription  string
	GlobalDatabase    string
	SubscriptionDBKey string
	URLCallbackHost   string
}

// Active Config object for package
var config *Config

// SetConfig ...
//   - ClientID: 			required
//   - ClientSecret: 		required
//   - PathPrefix:			"/admin" (default, should be "" if you don't want any prefix. "/" will be converted to "")
//   - PathSubscription:	"/subscription" (default, uses for Webhook callback)
//	 - GlobalDatabase:		"database" (default)
//	 - SubscriptionDBKey:	"strava-subscription" (default)
//   - URLCallbackHost:		(Current URL.Host, set different if you want to use another server, etc. ngrok)
//
//   Panic: if URLCallbackHost == "" and user call CreateSubcription()
func SetConfig(c Config) {
	newConfig := Config{}
	mergo.Merge(&newConfig, c)
	mergo.Merge(&newConfig, config)
	config = &newConfig

	if config.PathPrefix == "/" {
		config.PathPrefix = ""
	}
}

func (c *Config) getRedirectPath() string {
	return c.PathPrefix + c.PathRedirect
}

// GetRevokeURLFor return revoke link for username
func (c *Config) GetRevokeURLFor(username string) string {
	return c.PathPrefix + "/strava/revoke/" + username
}

// ActiveConfig return current config object.
func ActiveConfig() *Config { return config }

// GetAuthURL return authorize link to strava: ?client_id=<>&
func (c Config) GetAuthURL() string {
	res := fmt.Sprintf("https://www.strava.com/oauth/authorize?client_id=%s", c.ClientID)
	res += fmt.Sprintf("&redirect_uri=%s", c.getRedirectPath())
	res += fmt.Sprintf("&response_type=code&approval_prompt=auto&scope=%s", strings.Join(c.Scopes, ","))
	return res
}
