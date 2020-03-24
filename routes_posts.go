package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ledinhbao/blog/packages/models"
)

func inititalizePostRoutes(e *gin.Engine) {

	postRoutes := e.Group("/post")
	{
		postRoutes.GET("/list", displayPosts)
	}

}

func displayPosts(c *gin.Context) {
	db := c.MustGet(dbInstance).(*gorm.DB)
	posts := []models.Post{}
	db.Find(&posts)

	c.HTML(http.StatusOK, "post-display.html", gin.H{
		"title": "All Posts",
		"posts": posts,
	})
}
