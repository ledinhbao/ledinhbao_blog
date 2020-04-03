package main

import (
	"net/http"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ledinhbao/blog/core"
	"github.com/ledinhbao/blog/packages/models"
	"github.com/ledinhbao/blog/packages/sports/strava"
)

func initNonAuthAdminRoutes(r *gin.RouterGroup) {
	r.GET("/login", showAdminLoginPage)
	r.POST("/postLogin", adminPostLogin)
	r.GET("/admin/register", showAdminRegisterPage)
	r.POST("/admin/register", postAdminRegister)
}

func initAdminRoutes(r *gin.RouterGroup) {
	r.GET("/", displayAdminIndex)
	r.GET("/dashboard", displayAdminDashboard)
	r.GET("/logout", adminLogout)
}

func showAdminRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-register", gin.H{
		"message": "Admin Register",
	})
}

func postAdminRegister(c *gin.Context) {
	var formData = models.User{}
	formData.Username = c.PostForm("username")
	formData.SetPassword(c.PostForm("password"))
	formData.PasswordConfirm = c.PostForm("password2")
	formData.Role = 1

	message := ""

	db := c.MustGet(dbInstance).(*gorm.DB)
	db.Create(&formData)
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func displayAdminIndex(c *gin.Context) {
	ginview.HTML(c, http.StatusOK, "admin-index", gin.H{})
}

func showAdminLoginPage(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes("is-error")
	session.Save()
	ginview.HTML(c, http.StatusOK, "admin-login.html", gin.H{
		"errors": flashes,
	})
}

func adminPostLogin(c *gin.Context) {
	db := c.MustGet(dbInstance).(*gorm.DB)
	user := models.User{}
	passwordFromRequest := c.PostForm("password")
	db.Where("username = ?", c.PostForm("username")).First(&user)

	if user.TryPassword(passwordFromRequest) {
		session := sessions.Default(c)
		session.Set(userkey, user.Username)
		session.Set(authUserID, user.ID)
		session.Save()
		c.Redirect(http.StatusFound, "/admin/dashboard")
		// c.Abort()
	} else {
		session := sessions.Default(c)
		session.AddFlash("Wrong password", "is-error")
		session.Save()
		c.Redirect(http.StatusMovedPermanently, "/admin/login")
		// c.Abort()
	}
}

func adminLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/admin/login")
}

func displayAdminDashboard(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get(authUserID)

	db := c.MustGet(dbInstance).(*gorm.DB)
	userInfo := models.User{}
	stravaInfo := strava.Athlete{}
	stravaLink := strava.Link{}

	db.Where("id = ?", userID).First(&userInfo)
	db.Where(&strava.Link{UserID: userInfo.ID}).First(&stravaLink)
	db.Where(&strava.Athlete{Username: stravaInfo.Username}).First(&stravaInfo)

	var stravaSetting core.Setting
	db.Where(core.Setting{Key: strava.ActiveConfig().SubscriptionDBKey}).First(&stravaSetting)

	var lastRun strava.Activity
	db.Where(strava.Activity{
		AthleteID: stravaInfo.AthleteID,
		Type:      "Run",
	}).Order("start_date desc").First(&lastRun)

	var clubList []strava.StravaClub
	db.Find(&clubList)

	ginview.HTML(c, http.StatusOK, "admin-dashboard", gin.H{
		"user":              userInfo,
		"strava_link":       stravaLink,
		"athelete":          stravaInfo,
		"IsStravaConnected": stravaLink.ID > 0,
		"StravaRevokeURL":   strava.ActiveConfig().GetRevokeURLFor(stravaInfo.Username),

		"IsStravaSubscribed":   stravaSetting.ID > 0,
		"stravaSubscriptionID": stravaSetting.Value,

		"hasLastRun": lastRun.ID > 0,
		"lastRun":    lastRun,

		"hasClubList": len(clubList) > 0,
		"clubList":    clubList,
	})
}
