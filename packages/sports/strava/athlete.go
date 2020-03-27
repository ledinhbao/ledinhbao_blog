package strava

import "github.com/jinzhu/gorm"

// Athlete contains data reflects Strava's Athlete data.
type Athlete struct {
	gorm.Model
	AthleteID     uint64 `json:"id" mapstructure:"id"`
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
