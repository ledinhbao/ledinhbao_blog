package ngsf

import "time"

type (
	Transaction struct {
		ID        uint `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time `sql:"index"`
	}
)

// TableName add "ngsf_" prefix to table name
func (Transaction) TableName() string {
	return "ngsf_transactions"
}
