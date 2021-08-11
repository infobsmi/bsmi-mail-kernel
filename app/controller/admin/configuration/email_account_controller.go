package configuration

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/cnmade/bsmi-mail-kernel/app/orm/model"
	"github.com/cnmade/bsmi-mail-kernel/app/vo"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	vo2 "github.com/cnmade/bsmi-mail-kernel/pkg/common/vo"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/flosch/pongo2/v4"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"strconv"
)

type email_account_controller struct {

}

func NewEmailAccountController() *email_account_controller {
	return &email_account_controller{}
}



func (co *email_account_controller) ListAction(c *gin.Context) {

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
			"prevPage":        fmt.Sprintf("%d", prev_page),
			"nextPage":        fmt.Sprintf("%d", next_page),
		}))
	return
}

func (co *email_account_controller) AddAction(c *gin.Context) {

	c.HTML(200, "admin/configuration/email_account/add.html",
		common.Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        common.Config.Site_name,
			"siteDescription": common.Config.Site_description,
		}))
	return
}


func ConvertToStr(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func (co *email_account_controller) TestAction(c *gin.Context) {


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


	dec :=new(mime.WordDecoder)
	dec.CharsetReader= func(charset string, input io.Reader) (io.Reader, error) {
		common.Sugar.Infof("charset: %+v", charset)
		switch charset {
		case "gb2312","gbk","gb18030":
			content, err := ioutil.ReadAll(input)
			if err != nil {
				return nil, err
			}
			//ret:=bytes.NewReader(content)
			//ret:=transform.NewReader(bytes.NewReader(content), simplifiedchinese.HZGB2312.NewEncoder())

			utf8str:=ConvertToStr(string(content),"gbk","utf-8")
			t:=bytes.NewReader([]byte(utf8str))
			//ret:=utf8.DecodeRune(t)
			//log.Println(ret)
			return t, nil
		default:
			content, err := ioutil.ReadAll(input)
			if err != nil {
				return nil, err
			}
			t:=bytes.NewReader(content)
			//ret:=utf8.DecodeRune(t)
			//log.Println(ret)
			return t, nil

		}
	}

	for msg := range messages {
		tmpSubject := msg.Envelope.Subject
		common.Sugar.Infof("tmpSubject: %+v", tmpSubject)
		b, _ := dec.Decode(tmpSubject)
		log.Println("* " + b)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")
	c.JSON(200, gin.H{"config": emailConfig})

	return
}



func (co *email_account_controller) SaveAddAction(c *gin.Context) {
	var BI vo.Email_account_vo
	err := c.MustBindWith(&BI, binding.Form)
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