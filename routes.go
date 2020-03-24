package main

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ledinhbao/blog/packages/models"
)

func initializeRoutes(router *gin.Engine) {
	router.GET("/logout", logout)
	router.GET("/login", showLoginPage)
	router.POST("/login", login)
}

func showLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func login(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	db := c.MustGet(dbInstance).(*gorm.DB)
	user := models.User{}
	db.Where("username = ?", username).First(&user)
	if user.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Username not found",
		})
	} else if user.TryPassword(password) {
		session := sessions.Default(c)
		session.Set(userkey, user.Username)
		session.Save()
		c.Redirect(http.StatusFound, "/admin/dashboard")
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Wrong password",
		})
	}
}

func logout(c *gin.Context) {
	// Remove user from current session
	session := sessions.Default(c)
	session.Delete(userkey)
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
