package models

import (
	"github.com/jinzhu/gorm"
)

type Image struct {
	gorm.Model
	ID       string `gorm:"unique_index"`
	Path_or  string `gorm:"not null"`
	Path_neg string `gorm:"not null"`
}
