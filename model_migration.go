package main

import (
	"github.com/jinzhu/gorm"
	"github.com/ledinhbao/blog/core"
	"github.com/ledinhbao/blog/packages/models"
	"github.com/ledinhbao/blog/packages/sports/strava"
)

func modelMigration(db *gorm.DB) {
	db.AutoMigrate(core.Setting{})

	// Model migration
	db.AutoMigrate(&core.User{})
	db.AutoMigrate(&models.Post{})

	// Strava Module
	db.AutoMigrate(&strava.Link{})
	db.AutoMigrate(&strava.Athlete{})
	db.AutoMigrate(&strava.Activity{})
	db.AutoMigrate(&strava.Athlete{})
	db.AutoMigrate(&strava.StravaClub{})
}
