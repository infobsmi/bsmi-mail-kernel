package model

import (
	"gorm.io/gorm"
	"time"
)
type EmailAccount struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Email string
	ShortName string
	SmtpHost string
	SmtpPort string
	SmtpAccount string
	SmtpPassword string
	ImapHost string
	ImapPort string
	ImapAccount string
	ImapPassword string
	Status uint //0 未登录 1 已登录
}

func (EmailAccount) TableName()  string {
	return "email_account"

}