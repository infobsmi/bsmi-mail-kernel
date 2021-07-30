package vo

type AppConfig struct {
	Dbdsn            string
	Admin_user       string
	Admin_password   string
	AdminEmail 		 string
	Site_name        string
	Site_description string
	Site_url         string
	SrvMode          string
	ObjectStorageType uint //1 本地存储 ./vol/oss/  2 aws s3 云存储
	CaptchaEnabled uint
	HCaptchaSiteKey	string
	HCaptchaSecretKey string
	TwoAuthType uint
	PrivateMode uint
	ObjectStorage    struct {
		Aws_access_key_id     string
		Aws_secret_access_key string
		Aws_region            string
		Aws_bucket            string
		Cdn_url               string
	}
	EmailConfig struct {
		FromEmail string
		SmtpHost string
		SmtpPort int
		SmtpUser string
		SmtpPassword string
	}
	MailgunConfig struct {
		FromEmail string
		Domain string
		ApiKey string
	}
	BackupConfig struct {
		BackupEnabled uint
		BackupDir string
	}
	TongjiConfig struct {
		TongjiEnabled uint
		TongjiCode	string
	}
}

