package main

import (
	"html/template"
	"net/http"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
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

func ginviewBackendTemplateMiddleware() gin.HandlerFunc {
	return ginview.NewMiddleware(goview.Config{
		Root:         "views/backend",
		Extension:    ".html",
		Master:       "layout/master",
		Partials:     []string{},
		DisableCache: true,
		Funcs: template.FuncMap{
			"formatInKilometer": formatInKilometer,
			"formatStravaTime":  formatStravaTime,
		},
	})
}
