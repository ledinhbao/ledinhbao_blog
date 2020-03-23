package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username        string
	Password        string
	Role            int
	PasswordConfirm string `gorm:"-"`
}

// SetPassword will make a hash version from password
func (user *User) SetPassword(pwd string) {
	hashPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	user.Password = string(hashPwd)
}

func (user User) TryPassword(pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd))
	return err == nil
}
