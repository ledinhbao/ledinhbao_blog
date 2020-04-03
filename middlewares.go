package main

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AuthRequired is a middleware to check if the user is authorized or not.
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// unauthorize will be transfer to /admin/login
		session.AddFlash("Unauthorized!", "is-error")
		session.Save()
		c.Redirect(http.StatusPermanentRedirect, "/admin/login")
	}
}

func dbHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(dbInstance, db)
		c.Next()
	}
}
