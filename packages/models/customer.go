package models

import (
	"time"

	"github.com/ledinhbao/blog/core"
)

type (
	Customer struct {
		ID        uint       `json:"customer_id" form:"customer_id" gorm:"primary_key"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at" sql:"index"`
		Fullname  string     `json:"fullname" form:"fullname"`
		DOB       time.Time  `json:"dob" form:"dob"`
		UserID    uint       `json:"user_id" form:"user_id"`
		User      core.User  `gorm:"association_autoupdate:false;association_autocreate:false"`
	}
)
