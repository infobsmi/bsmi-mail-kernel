package vo

type Email_account_vo struct {
	Email   string `form:"email" binding:"required"`
	ShortName string `form:"short_name" binding:"required"`
	SmtpHost  string  `form:"smtp_host" `
	SmtpPort    string `form:"smtp_port" `
	SmtpAccount    string  `form:"smtp_account" `
	SmtpPassword    string  `form:"smtp_password" `
	ImapHost  string  `form:"imap_host" `
	ImapPort    string `form:"imap_port" `
	ImapAccount    string  `form:"imap_account" `
	ImapPassword    string  `form:"imap_password" `
}

type Edit_email_account_vo struct {

	Email   string `form:"email" binding:"required"`
	ShortName string `form:"short_name" binding:"required"`
	SmtpHost  string  `form:"smtp_host" `
	SmtpPort    string `form:"smtp_port" `
	SmtpAccount    string  `form:"smtp_account" `
	SmtpPassword    string  `form:"smtp_password" `
	ImapHost  string  `form:"imap_host" `
	ImapPort    string `form:"imap_port" `
	ImapAccount    string  `form:"imap_account" `
	ImapPassword    string  `form:"imap_password" `
}

