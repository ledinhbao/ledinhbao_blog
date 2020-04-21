package main

import (
	"math/rand"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ledinhbao/blog/core"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(bytes), err
}

func randString() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func setSession(c *gin.Context, key string, value string) {
	session := sessions.Default(c)
	session.Set(key, value)
	session.Save()
}

type noUserError struct{}
type invalidUserError struct{}

func (noUserError) Error() string {
	return "No user exists in session"
}
func (invalidUserError) Error() string {
	return "User's data is invalid"
}

func authUserFromSession(c *gin.Context) (core.User, error) {
	s := sessions.Default(c)
	user, ok := s.Get(authUser).(*core.User)
	if !ok {
		return core.User{}, noUserError{}
	}
	if user.ID <= 0 {
		return core.User{}, invalidUserError{}
	}
	return *user, nil
}
