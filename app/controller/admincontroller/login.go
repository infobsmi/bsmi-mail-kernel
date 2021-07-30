package admincontroller

import (
	"github.com/cnmade/bsmi-mail-kernel/app/orm/model"
	"github.com/cnmade/bsmi-mail-kernel/app/service/email_service"
	"github.com/cnmade/bsmi-mail-kernel/app/service/fail_ban_service"
	"github.com/cnmade/bsmi-mail-kernel/app/service/mailgun_service"
	"github.com/cnmade/bsmi-mail-kernel/app/utils/http_utils"
	"github.com/cnmade/bsmi-mail-kernel/app/vo/login_vo"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common/vo"
	"github.com/flosch/pongo2/v4"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	gonanoid "github.com/matoous/go-nanoid"
	"time"
)

func LoginCtr(c *gin.Context) {

	session := sessions.Default(c)
	session.Delete("username")
	session.Delete("isAdmin")

	c.HTML(200, "admin/login.html",
		common.Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        common.Config.Site_name,
			"siteDescription": common.Config.Site_description,
		}))
	return
}

func LoginProcessCtr(c *gin.Context) {
	clientIp := http_utils.GetClientIp(c)
	if fail_ban_service.CheckBan(clientIp) == true {
			common.ShowUMessage(c, &vo.Umsg{Msg: "登录失败,您已经被记录"})
			//fail ban
			return
	}
	//检查fail ban，如果被封禁了，就果断拦截
	if common.Config.CaptchaEnabled == 1 {
		hcaptchaResp := common.HCaptchClient.SiteVerify(c.Request)
		if !hcaptchaResp.Success {

			common.Sugar.Errorf("error: %+v", hcaptchaResp)
			common.ShowUMessage(c, &vo.Umsg{Msg: "登录失败,验证码不对"})
			fail_ban_service.BanIp(clientIp)
			//fail ban
			return
		}
	}
	var form AdminLoginForm
	err := c.MustBindWith(&form, binding.Form)
	if err != nil {
		common.LogError(err)
		common.ShowUMessage(c, &vo.Umsg{Msg: "登录失败"})
		//fail ban
		fail_ban_service.BanIp(clientIp)
		return
	}
	session := sessions.Default(c)
	if form.Username == common.Config.Admin_user && form.Password == common.Config.Admin_password {
		//这里显示：您开启了两步登录，请完成两步登录流程

		keyA, _ := gonanoid.Nanoid(20)
		keyB, _ := gonanoid.Nanoid(20)

		switch common.Config.TwoAuthType {
		case 1:
			email_service.SendTwoAuth(common.Config.AdminEmail, keyB)
			break;
		case 2:
			mailgun_service.SendTwoAuth(common.Config.AdminEmail, keyB)
			break;
		}

		loc, _ := time.LoadLocation("Asia/Shanghai")
		var item model.TwoAuth
		item.CreatedAt = time.Now().In(loc)
		item.KeyA = keyA;
		item.KeyB = keyB;
		common.NewDb.Create(&item)

		c.HTML(200, "admin/login-process.html",
			common.Pongo2ContextWithVersion(pongo2.Context{
				"siteName":        common.Config.Site_name,
				"siteDescription": common.Config.Site_description,
			}))
		return
	} else {
		session.Delete("username")
		session.Delete("isAdmin")
		err := session.Save()
		if err != nil {
			common.LogError(err)
		}
		//fail ban
		fail_ban_service.BanIp(clientIp)
		common.ShowUMessage(c, &vo.Umsg{Msg: "登录失败", Url: "/"})
	}
}

func LoginStep3Ctr(c *gin.Context) {
	clientIp := http_utils.GetClientIp(c)
	if fail_ban_service.CheckBan(clientIp) == true {
		common.ShowUMessage(c, &vo.Umsg{Msg: "登录失败,您已经被记录"})
		//fail ban
		return
	}
	if common.Config.CaptchaEnabled == 1 {
		hcaptchaResp := common.HCaptchClient.SiteVerify(c.Request)
		if !hcaptchaResp.Success {

			common.Sugar.Errorf("error: %+v", hcaptchaResp)
			common.ShowUMessage(c, &vo.Umsg{Msg: "登录失败,验证码不对"})
			//fail ban
			return
		}
	}
	var form login_vo.TwoAuthForm
	err := c.MustBindWith(&form, binding.Form)
	if err != nil {
		common.LogError(err)
		common.ShowUMessage(c, &vo.Umsg{Msg: "登录失败"})
		//fail ban
		fail_ban_service.BanIp(clientIp)
		return
	}



	session := sessions.Default(c)

	var twoAuthItem model.TwoAuth;
	common.NewDb.Find(&twoAuthItem, "key_b = ?", form.TwoAuthCode)

	//登录失败
	//根据key a 查 key b是否被点过， 如果点过，就把这条记录标记为已登录
	if  twoAuthItem.KeyB == form.TwoAuthCode && twoAuthItem.Status == 0 {
		//登录成功
		session.Set("username", common.Config.Admin_user)
		session.Set("isAdmin", "yes")
		err := session.Save()
		if err != nil {
			common.LogError(err)
		}
		twoAuthItem.Status = 1
		common.NewDb.Save(&twoAuthItem)

		c.Redirect(301, "/")
	} else {
		session.Delete("username")
		session.Delete("isAdmin")
		err := session.Save()
		if err != nil {
			common.LogError(err)
		}
		//fail ban
		fail_ban_service.BanIp(clientIp)
		common.ShowUMessage(c, &vo.Umsg{Msg: "登录失败", Url: "/"})
	}

}
func LogoutCtr(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("username")
	session.Delete("isAdmin")
	err := session.Save()
	if err != nil {
		common.LogError(err)
	}
	c.Redirect(301, "/")
}
