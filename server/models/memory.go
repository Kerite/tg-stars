package models

import (
	"gorm.io/gorm"
)

type Memory struct {
	gorm.Model
	FileId         string
	Owner          string
	Description    string
	Price          int64
	Anoymous       bool `gorm:"default:false"`
	IsSubscription bool `gorm:"default:false"`
}
