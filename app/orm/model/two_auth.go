package model

import (
	"gorm.io/gorm"
	"time"
)
type TwoAuth struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	KeyA string
	KeyB string
	KeyBStatus uint
	Status uint //0 未登录 1 已登录
}

func (TwoAuth) TableName()  string {
	return "bk_two_auth"

}