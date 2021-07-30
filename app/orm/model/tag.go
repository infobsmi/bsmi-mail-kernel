package model

import (
	"gorm.io/gorm"
	"time"
)
type Tag struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name string
	TotalNums int64
}

