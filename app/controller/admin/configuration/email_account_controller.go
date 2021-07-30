package configuration

import (
	"errors"
	"fmt"
	"github.com/cnmade/bsmi-mail-kernel/app/orm/model"
	"github.com/cnmade/bsmi-mail-kernel/app/utils/admin_utils"
	"github.com/cnmade/bsmi-mail-kernel/app/vo"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	vo2 "github.com/cnmade/bsmi-mail-kernel/pkg/common/vo"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/flosch/pongo2/v4"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type email_account_controller struct {

}

func NewEmailAccountController() *email_account_controller {
	return &email_account_controller{}
}



func (co *email_account_controller) ListAction(c *gin.Context) {
	err, username, isAdmin := admin_utils.AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Fatal(err)
	}
	page -= 1
	if page < 0 {
		page = 0
	}

	prev_page := page
	if prev_page < 1 {
		prev_page = 1
	}
	next_page := page + 2
	rpp := 20
	offset := page * rpp
	log.Println(rpp)
	log.Println(offset)
	var emailAccountList []model.EmailAccount

	_ = common.NewDb.
		Limit(rpp).
		Offset(offset).
		Order("id desc").
		Find(&emailAccountList)
	common.Sugar.Infof("emailAccountList: %+v", emailAccountList)


	c.HTML(200, "admin/configuration/email_account/list.html",
		common.Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        common.Config.Site_name,
			"siteDescription": common.Config.Site_description,
			"emailAccountList":    emailAccountList,
			"username":        username.(string),
			"isAdmin":        isAdmin.(string),
			"prevPage":        fmt.Sprintf("%d", prev_page),
			"nextPage":        fmt.Sprintf("%d", next_page),
		}))
	return
}

func (co *email_account_controller) AddAction(c *gin.Context) {
	err, _, _ := admin_utils.AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}

	c.HTML(200, "admin/configuration/email_account/add.html",
		common.Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        common.Config.Site_name,
			"siteDescription": common.Config.Site_description,
		}))
	return
}



func (co *email_account_controller) TestAction(c *gin.Context) {
	err, _, _ := admin_utils.AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}


	id := c.Param("id")
	var emailConfig model.EmailAccount
	result := common.NewDb.First(&emailConfig, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(200, gin.H{"msg": "配置不存在"})
		return
	}

	imapAddr := fmt.Sprintf("%s:%s", emailConfig.ImapHost, emailConfig.ImapPort)
	cmail, err := client.DialTLS(imapAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer cmail.Logout()

	// Login
	if err := cmail.Login(emailConfig.ImapAccount, emailConfig.ImapPassword); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func () {
		done <- cmail.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := cmail.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)

	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 3 {
		// We're using unsigned integers here, only subtract if the result is > 0
		from = mbox.Messages - 3
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- cmail.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	log.Println("Last 4 messages:")
	for msg := range messages {
		log.Println("* " + msg.Envelope.Subject)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")
	c.JSON(200, gin.H{"config": emailConfig})

	return
}



func (co *email_account_controller) SaveAddAction(c *gin.Context) {
	err, _, _ := admin_utils.AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	var BI vo.Email_account_vo
	err = c.MustBindWith(&BI, binding.Form)
	if err != nil {
		common.ShowUMessage(c, &vo2.Umsg{Msg: err.Error(), Url: "/"})
		common.LogError(err)
		return
	}


	blogItem := model.EmailAccount{
		Email: BI.Email,
		ShortName: BI.ShortName,
		SmtpHost: BI.SmtpHost,
		SmtpPort: BI.SmtpPort,
		SmtpAccount: BI.SmtpAccount,
		SmtpPassword: BI.SmtpPassword,
		ImapHost: BI.ImapHost,
		ImapPort: BI.ImapPort,
		ImapAccount: BI.ImapAccount,
		ImapPassword: BI.ImapPassword,
		Status: 0,
	}

	result := common.NewDb.Create(&blogItem)


	if result.Error == nil {
		c.Redirect(http.StatusFound, "/admin/configuration/email_account/list")
	} else {
		common.LogError(result.Error)
		common.ShowUMessage(c, &vo2.Umsg{Msg: "失败", Url: "/"})
	}
}