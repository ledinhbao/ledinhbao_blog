package ngsf

import "github.com/jinzhu/gorm"

// MigrateModels call auto migration for all Model in this module
func DatabaseMigration(db *gorm.DB) {
	db.AutoMigrate(&Customer{})
	db.AutoMigrate(&Order{})
	db.AutoMigrate(&Transaction{})
}
