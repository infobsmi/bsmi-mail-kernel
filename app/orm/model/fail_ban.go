package model

import (
	"gorm.io/gorm"
	"time"
)
type FailBan struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Ip string
	Count uint
	Status uint //0 未登录 1 已登录
}

func (FailBan) TableName()  string {
	return "bk_fail_ban"

}