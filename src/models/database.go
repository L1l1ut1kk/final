package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Database() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("images.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Create the images table
	db.AutoMigrate(&Image{})

	return db
}
