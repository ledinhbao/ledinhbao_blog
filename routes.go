package main

import "github.com/gin-gonic/gin"

func initializeRoutes(router *gin.Engine) {

	router.GET("/post/:post-id", func(c *gin.Context) {})

}
