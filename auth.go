package main

import (
	"encoding/gob"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ledinhbao/blog/core"
)

func initSession(engine *gin.Engine) {
	cookieStore := sessions.NewCookieStore([]byte(appName + "-" + appUUID))
	engine.Use(sessions.Sessions(appName+"-sessions", cookieStore))

	// Register User model
	gob.Register(&core.User{})
}
