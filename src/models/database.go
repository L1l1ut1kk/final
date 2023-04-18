package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Database() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("src/models/images.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Create the images table
	db.AutoMigrate(&Image{})

	return db
}
