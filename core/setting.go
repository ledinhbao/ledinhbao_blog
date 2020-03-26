package core

import "github.com/jinzhu/gorm"

// Setting contains server configuration, database-based.
type Setting struct {
	gorm.Model
	Key   string
	Value string
}
