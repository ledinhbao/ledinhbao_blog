package main

import (
	"fmt"
	"net/http"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ledinhbao/blog/core"
	"github.com/ledinhbao/blog/packages/sports/strava"
)

func initNonAuthAdminRoutes(r *gin.RouterGroup) {
	r.GET("/login", showAdminLoginPage)
	r.POST("/login", adminAuthenticate)
	r.GET("/register", showAdminRegisterPage)
	r.POST("/register", postAdminRegister)

	r.GET("/unauthorized", adminUnauthorized)
}

func initAdminRoutes(r *gin.RouterGroup) {
	r.GET("/", displayAdminIndex)
	r.GET("/dashboard", displayAdminDashboard)
	r.GET("/logout", adminLogout)
}

func initSuperAdminRoutes(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		ginview.HTML(c, http.StatusOK, "su-homepage", gin.H{})
	})
}

func showAdminRegisterPage(c *gin.Context) {
	ginview.HTML(c, http.StatusOK, "admin-register.html", gin.H{
		"message": "Admin Register",
	})
}

func postAdminRegister(c *gin.Context) {
	var formData = core.User{}
	formData.Username = c.PostForm("username")
	formData.SetPassword(c.PostForm("password"))
	formData.PasswordConfirm = c.PostForm("password2")
	formData.Rank = 1

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

func adminAuthenticate(c *gin.Context) {
	db := c.MustGet(dbInstance).(*gorm.DB)
	user := core.User{}
	passwordFromRequest := c.PostForm("password")
	db.Where("username = ?", c.PostForm("username")).First(&user)

	if user.ID == 0 {
		// User not found
		ginview.HTML(c, http.StatusUnauthorized, "admin-login.html", gin.H{
			"errors": []string{"Username or password is incorrect."},
		})
	}
	if user.TryPassword(passwordFromRequest) {
		if user.Rank < RoleAdmin {
			ginview.HTML(c, http.StatusUnauthorized, "admin-login.html", gin.H{
				"errors": []string{"You are not authorized to view this page."},
			})
			return
		}
		session := sessions.Default(c)
		session.Set(authUser, user)
		session.Set(userkey, user.Username)
		session.Set(authUserID, user.ID)
		session.Save()
		core.AddUserToSession(user, c)
		c.Redirect(http.StatusFound, "/admin/dashboard")
	} else {
		ginview.HTML(c, http.StatusUnauthorized, "admin-login.html", gin.H{
			"errors": []string{"Username or password is incorrect."},
		})
	}
}

func adminLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/admin/login")
}

func adminUnauthorized(c *gin.Context) {
	ginview.HTML(c, http.StatusUnauthorized, "admin-unauthorized.html", nil)
}

func displayAdminDashboard(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get(authUserID)

	db := c.MustGet(dbInstance).(*gorm.DB)
	userInfo := core.User{}
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

	url := location.Get(c)
	stravaCallbackURL := fmt.Sprintf("%s://%s/admin/strava", url.Scheme, c.Request.Host)

	ginview.HTML(c, http.StatusOK, "admin-dashboard", gin.H{
		"user":              userInfo,
		"strava_link":       stravaLink,
		"athelete":          stravaInfo,
		"IsStravaConnected": stravaLink.ID > 0,
		"StravaRevokeURL":   strava.ActiveConfig().GetRevokeURLFor(stravaInfo.Username),
		"stravaAuthURL":     strava.GetOAuthURL(stravaCallbackURL),

		"IsStravaSubscribed":   stravaSetting.ID > 0,
		"stravaSubscriptionID": stravaSetting.Value,

		"hasLastRun": lastRun.ID > 0,
		"lastRun":    lastRun,

		"hasClubList": len(clubList) > 0,
		"clubList":    clubList,
	})
}
