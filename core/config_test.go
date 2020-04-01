package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringValueForKey(t *testing.T) {
	sampleConfigStr := `{
		"application": {
			"mode": "development",
			"admin-prefix": "/admin"
		},
		"database": {
			"development": {
				"username": "",
				"password": "",
				"host": "/database.db"
			},
			"production": {
				"username": "",
				"password": "",
				"host": "localhost"
			}
		},
		"strava": {
			"webhook-callback": "http://bc7b66a4.ngrok.io"
		}
	}`
	var config Config
	err := json.Unmarshal([]byte(sampleConfigStr), &config)

	actual, err := config.StringValueForKey("application.admin-prefix")
	assert.Nil(t, err)
	assert.Equal(t, "/admin", actual)

	// this is not exist
	actual, err = config.StringValueForKey("application.database.development.host")
	assert.NotNil(t, err)

	actual, err = config.StringValueForKey("database.production.host")
	assert.Nil(t, err)
	assert.Equal(t, "localhost", actual)

	actual, err = config.StringValueForKey("strava-non-exist-key")
	assert.NotNil(t, err)
}

func TestReadFile(t *testing.T) {
	sampleConfig := `{
		"application": {
			"mode": "development",
			"admin-prefix": "/admin"
		},
		"database": {
			"development": {
				"driver": "sqlite3",
				"username": "",
				"password": "",
				"host": "/database.db"
			},
			"production": {
				"driver": "mysql",
				"username": "ledinhbao_axis",
				"password": "L93hxwPc8r",
				"host": "localhost"
			}
		},
		"strava": {
			"webhook-callback": "http://bc7b66a4.ngrok.io"
		}
	}`
	config, err := NewConfigFromJSONFile("../config.json")
	assert.Nil(t, err)

	var expectedConfig Config
	err = json.Unmarshal([]byte(sampleConfig), &expectedConfig)
	assert.Nil(t, err)
	assert.Equal(t, expectedConfig, config)
}
