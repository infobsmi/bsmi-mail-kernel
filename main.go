package main

import (
	. "github.com/cnmade/bsmi-mail-kernel/app/controller"
	"github.com/cnmade/bsmi-mail-kernel/app/controller/admin"
	"github.com/cnmade/bsmi-mail-kernel/app/controller/admin/configuration"
	"github.com/cnmade/bsmi-mail-kernel/app/service/backup_service"
	"github.com/cnmade/bsmi-mail-kernel/app/service/fail_ban_service"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	"github.com/cnmade/pongo2gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"time"
)

//go:generate go run cmd/version_info.go
func main() {
	common.InitApp()

	common.Sugar.Infof("os.Args len: %d", len(os.Args))
	if len(os.Args) > 1 {
		cliArgs := os.Args[1]
		cliArgs1 := ""
		if len(os.Args) > 2 {

			cliArgs1 = os.Args[2]
		}

		switch (cliArgs) {
		case "clearban":
			fail_ban_service.ClearBan(cliArgs1)
			return

		}
	}

	//启动的时候，执行一次备份
	go backup_service.DoBackup()
	s := gocron.NewScheduler(time.UTC)

	s.Cron("5 5 */3 * *").Do(backup_service.DoBackup)

	r := gin.New()
	r.HTMLRender = pongo2gin.New(pongo2gin.RenderOptions{
		TemplateDir: "views",
		ContentType: "text/html; charset=utf-8",
		AlwaysNoCache: true,
	})


	version_bytes, _ := ioutil.ReadFile("./public/version.js")
	common.BsmiKbVersion = string(version_bytes)

	r.Static("/assets", "./public/assets")
	r.Static("/oss", "./vol/oss")

	store := cookie.NewStore([]byte("gssecret"))
	store.Options(sessions.Options{
		Path: "/",
		MaxAge: 999999999,
		HttpOnly: true,

	})
	r.Use(sessions.Sessions("mysession", store))
	fc := new(FrontController)
	r.GET("/", fc.HomeCtr)
	r.HEAD("/", fc.HomeCtr)
	r.GET("/list-tag", fc.ListTagCtr)
	r.GET("/demopongo", fc.DemoPongoCtr)
	r.GET("/about", fc.AboutCtr)
	r.GET("/view/:id", fc.ViewCtr)
	//查看 文章历史记录的详情页面
	r.GET("/view_article_history/:id", fc.ViewArticleHistoryCtr)
	r.GET("/view.php", fc.ViewAltCtr)
	r.GET("/ping", fc.PingCtr)
	r.GET("/search", fc.SearchCtr)
	//查看 文章历史记录的 列表页面
	r.GET("/article_history", fc.ArticleHistoryCtr)
	r.GET("/charge", fc.ChargeCtr)
	r.GET("/user/logout", fc.LogoutCtr)
	r.GET("/countview/:id", fc.CountViewCtr)

	adminRouteGroup := r.Group("/admin")
	{
		admin_co := admin.NewAdminController()

		blog_co := admin.NewBlogController()

		adminRouteGroup.GET("/", blog_co.ListBlogCtr)

		adminRouteGroup.GET("/login", admin.LoginCtr)
		adminRouteGroup.POST("/login-process", admin.LoginProcessCtr)
		adminRouteGroup.POST("/loginStep3", admin.LoginStep3Ctr)
		adminRouteGroup.GET("/logout", admin.LogoutCtr)
		adminRouteGroup.GET("/export", admin_co.ExportCtr)



		adminRouteGroup.GET("/list-cate", admin.ListCateCtr)
		adminRouteGroup.POST("/save-edit-cate", admin.SaveEditCateCtr)
		adminRouteGroup.GET("/edit-cate/:id", admin.EditCateCtr)
		adminRouteGroup.GET("/add-cate", admin.AddCateCtr)
		adminRouteGroup.POST("/save-add-cate", admin.SaveAddCateCtr)

		adminRouteGroup.GET("/list-tag", admin.ListTagCtr)


		adminRouteGroup.GET("/files", admin_co.Files)
		adminRouteGroup.POST("/fileupload", admin_co.FileUpload)
	}




	// 博客相关
	adminBlogRouteGroup := r.Group("/admin/blog")
	{

		//Email account

		blog_co := admin.NewBlogController()
		adminBlogRouteGroup.GET("/list", blog_co.ListBlogCtr)


		adminBlogRouteGroup.GET("/addblog", blog_co.AddBlogCtr)
		adminBlogRouteGroup.POST("/save-blog-add", blog_co.SaveBlogAddCtr)
		adminBlogRouteGroup.GET("/listblog", blog_co.ListBlogCtr)
		adminBlogRouteGroup.GET("/deleteblog/:id", blog_co.DeleteBlogCtr)
		adminBlogRouteGroup.POST("/save-blog-edit", blog_co.SaveBlogEditCtr)
		adminBlogRouteGroup.GET("/editblog/:id", blog_co.EditBlogCtr)

	}

	//配置设置
	adminConfigurationRouteGroup := r.Group("/admin/configuration")
	{

		//Email account

		email_account_co := configuration.NewEmailAccountController()
		adminConfigurationRouteGroup.GET("/email_account/list", email_account_co.ListAction)
		adminConfigurationRouteGroup.GET("/email_account/add", email_account_co.AddAction)
		adminConfigurationRouteGroup.GET("/email_account/saveAdd", email_account_co.SaveAddAction)


	}


	// rss


	rss := new(RSS)
	r.GET("/rss.php", rss.Alter)
	r.GET("/rss", rss.Out)


	a := new(Api)
	api := r.Group("/api")
	{
		api.GET("/nav-all", a.NavAll)
		api.GET("/nav-load", a.NavLoad)
		api.POST("/resort", a.Resort)
	}
	log.Info().Msg("Server listen on 127.0.0.1:3711")
	err := r.Run("127.0.0.1:3711")
	if err != nil {
		common.LogError(err)
	}
}
