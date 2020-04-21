package main

import (
	"net/http"
	"strconv"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ledinhbao/blog/core"
	"github.com/ledinhbao/blog/packages/sports/strava"
)

func initStravaRoutes(r *gin.RouterGroup) {
	r.GET("/add", stravaAddClub)
	r.POST("/connect", stravaConnectClub)
	r.GET("/remove/:club-id", stravaRemoveClub)
}

func stravaAddClub(c *gin.Context) {
	db := c.MustGet(dbInstance).(*gorm.DB)
	session := sessions.Default(c)
	var userInfo core.User
	db.Where("id = ?", session.Get(authUserID)).First(&userInfo)
	ginview.HTML(c, http.StatusOK, "admin-strava-add-club", gin.H{
		"user": userInfo,
	})
}

func stravaConnectClub(c *gin.Context) {
	db := c.MustGet(dbInstance).(*gorm.DB)
	session := sessions.Default(c)
	userID := session.Get(authUserID).(uint)
	clubID := c.PostForm("club_id")

	var club strava.StravaClub
	db.Where("club_id = ?", clubID).First(&club)
	if club.ID > 0 {
		// There is a group with the same id
		session.AddFlash("There is a club with ID: "+clubID+" had been added.", "strava-clubs-add-error")
		session.Save()
	} else {
		// Save strava club with id, and state of processing is 1 (Processing data from Strava)
		clubIDUInt64, _ := strconv.ParseUint(clubID, 10, 32)
		club.ClubID = uint(clubIDUInt64)
		club.ProcessingState = 1
		db.Create(&club)
		go strava.StravaFetchClub(club, userID, db)
	}
	c.Redirect(http.StatusFound, "/admin/dashboard")

}

func stravaRemoveClub(c *gin.Context) {
	db := c.MustGet(dbInstance).(*gorm.DB)
	clubID, err := strconv.ParseUint(c.Query("club-id"), 10, 64)
	if err != nil {
		session := sessions.Default(c)
		session.AddFlash("Missing club-id for delete.", "is-error")
		session.Save()

	} else {
		db.Delete(strava.StravaClub{ClubID: uint(clubID)}).Delete(&strava.StravaClub{})
	}
	c.Redirect(http.StatusNotFound, "/admin/dashboard")
}
