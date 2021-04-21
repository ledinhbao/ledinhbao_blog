package ngsf

import (
	"time"

	"github.com/ledinhbao/blog/core"
)

type (
	Order struct {
		ID        uint `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time `sql:"index"`
		Buy       bool
		Amount    int
		Price     float64
		Date      time.Time
		UserID    uint
		User      core.User `gorm:"association_autoupdate:false;association_autocreate:false"`
	}
)

// TableName add "ngsf_" prefix to table name
func (Order) TableName() string {
	return "ngsf_orders"
}
