package models

import (
	"github.com/jinzhu/gorm"
)

type Image struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Path_or  string `gorm:"not null"`
	Path_neg string `gorm:"not null"`
}
