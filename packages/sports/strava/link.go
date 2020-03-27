package strava

import "github.com/jinzhu/gorm"

// Link (aka Strava Link) contains data linked between app's user & strava's athlete data
type Link struct {
	gorm.Model
	UserID       uint
	Username     string
	AthleteID    uint64
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
