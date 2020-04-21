package core

import (
	"fmt"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type (
	// User stores basic user's construction
	User struct {
		ID              uint   `gorm:"primary_key"`
		Username        string `gorm:"unique_index" form:"username"`
		Password        string `form:"password"`
		PasswordConfirm string `form:"confirm_password" gorm:"-"`
		Email           string `form:"email"`
		Rank            int    `form:"rank"`

		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time `sql:"index"`
	}

	// UserRank determines a rank of user
	UserRank int
)

const (
	// RankSuperAdmin is the one who own the website
	RankSuperAdmin = UserRank(89)
	// RankAdmin who has power
	RankAdmin = UserRank(55)
	// RankModerator who work as Robin where Admin is batman
	RankModerator = UserRank(34)
	// RankWriter : who use their typing skill as good as Ernest Hemingway
	RankWriter = UserRank(21)
	// RankViewer : who just sit and watch.
	RankViewer = UserRank(1)

	defaultUserKeystring = string("ledinhbao-default-user-key-string")
)

// SetPassword will make a hash version from password
func (user *User) SetPassword(pwd string) {
	hashPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	user.Password = string(hashPwd)
}

// TryPassword compares the given password (pwd) with user's hash password.
func (user User) TryPassword(pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd))
	return err == nil
}

// Validate makes sure the data is not corrupted
func (user User) Validate() (err error) {
	if user.Username == "" {
		err = fmt.Errorf("Username cannot be empty")
	}
	return
}

// AddUserToSession add authorized user to session with defaultUserKeystring
func AddUserToSession(user User, c *gin.Context) {
	s := sessions.Default(c)
	// Clear out password field
	user.Password = ""
	user.PasswordConfirm = ""
	s.Set(defaultUserKeystring, user)
	s.Save()

	logrus.WithFields(logrus.Fields{
		"module": "core",
		"action": "AddUserToSession",
		"data":   user,
	}).Info("User added to Session using default key")
}

// UserRankRequired Use this middleware to require any user's rank
func UserRankRequired(rank UserRank, unauthorizedCallback func(*gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessions.Default(c)
		user, ok := s.Get(defaultUserKeystring).(*User)
		if !ok || user.ID == 0 || UserRank(user.Rank) < RankAdmin {
			if unauthorizedCallback != nil {
				unauthorizedCallback(c)
			}
			c.Abort()
		} else {
			c.Next()
		}
	}
}
