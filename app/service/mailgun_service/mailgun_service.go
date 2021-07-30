package mailgun_service

import (
	"context"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	"github.com/mailgun/mailgun-go/v4"
	"time"
)

func SendTwoAuth(ToEmail string, KeyB string) {
	mailgunConfig := common.Config.MailgunConfig
	mg := mailgun.NewMailgun(mailgunConfig.Domain, mailgunConfig.ApiKey)

	sender := mailgunConfig.FromEmail
	subject := "两步登录验证"
	body := "您的两步登录验证码为：" + KeyB
	recipient := common.Config.AdminEmail



	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		common.Sugar.Error(err)
	}

	common.Sugar.Infof("ID: %s Resp: %s\n", id, resp)
}
