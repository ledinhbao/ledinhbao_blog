package core

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (
	// User stores basic user's construction
	User struct {
		ID              uint   `gorm:"primary_key"`
		Username        string `gorm:"unique_index" form:"username"`
		Password        string `form:"password"`
		passwordConfirm string `form:"confirm_password" gorm:"-"`
		Email           string `form:"email"`
		Rank            int    `form:"rank"`

		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time `sql:"index"`
	}
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
