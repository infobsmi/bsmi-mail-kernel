package email_service

import (
	"crypto/tls"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	"github.com/go-mail/mail"
)

func SendTwoAuth(ToEmail string, KeyB string) {
	emailConfig := common.Config.EmailConfig
	d := mail.NewDialer(
		emailConfig.SmtpHost,
		emailConfig.SmtpPort,
		emailConfig.SmtpUser,
		emailConfig.SmtpPassword,
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	m := mail.NewMessage()
	m.SetHeader("From", emailConfig.FromEmail)
	m.SetHeader("To", ToEmail)
	m.SetHeader("Subject", "两步登陆验证")
	m.SetBody("text/html", "您的两步登录验证码为：" + KeyB )
	if err := d.DialAndSend(m); err != nil {
		common.Sugar.Error(err)
	}
}
