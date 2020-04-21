package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ledinhbao/blog/core"
)

// AuthRequired is a middleware to check if the user is authorized or not.
func AuthRequired(c *gin.Context) {
	fmt.Println("AuthRequired")
	session := sessions.Default(c)
	user, ok := session.Get(authUser).(*core.User)
	if !ok || user.ID == 0 {
		// unauthorize will be transfer to /admin/login
		session.AddFlash("Unauthorized!", "is-error")
		session.Save()
		c.Redirect(http.StatusFound, "/admin/login")
		c.Abort()
	} else if user.Rank < RoleAdmin {
		ginview.HTML(c, http.StatusUnauthorized, "admin-unauthorized.html", gin.H{})
		c.Abort()
	} else {
		c.Next()
	}
}

// UnauthorizedHandler display template for unauthorized request
func UnauthorizedHandler(c *gin.Context) {
	ginview.HTML(c, http.StatusUnauthorized, "admin-unauthorized.html", gin.H{})
	c.Next()
}

// SuperAdminRequired is middleware to check user is SuperAdmin or not
func SuperAdminRequired() gin.HandlerFunc {
	fmt.Println("SuperAdminRequired")
	return func(c *gin.Context) {
		user, err := authUserFromSession(c)
		switch err.(type) {
		case noUserError, invalidUserError:
			ginview.HTML(c, http.StatusUnauthorized, "admin-unauthorized.html", nil)
			c.Abort()
			return
		}
		if user.Rank < RoleSuperAdmin {
			ginview.HTML(c, http.StatusUnauthorized, "admin-unauthorized.html", nil)
			c.Abort()
		} else {
			c.Next()
		}
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
