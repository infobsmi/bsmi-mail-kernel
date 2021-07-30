package main

import (
	. "github.com/cnmade/bsmi-mail-kernel/app/controller"
	"github.com/cnmade/bsmi-mail-kernel/app/controller/admincontroller"
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

	admin := r.Group("/admin")
	{
		admin.GET("/", admincontroller.ListBlogCtr)
		admin.GET("/login", admincontroller.LoginCtr)
		admin.POST("/login-process", admincontroller.LoginProcessCtr)
		admin.POST("/loginStep3", admincontroller.LoginStep3Ctr)
		admin.GET("/logout", admincontroller.LogoutCtr)
		admin.GET("/addblog", admincontroller.AddBlogCtr)
		admin.POST("/save-blog-add", admincontroller.SaveBlogAddCtr)
		admin.GET("/listblog", admincontroller.ListBlogCtr)
		admin.GET("/export", admincontroller.ExportCtr)
		admin.GET("/deleteblog/:id", admincontroller.DeleteBlogCtr)
		admin.POST("/save-blog-edit", admincontroller.SaveBlogEditCtr)
		admin.GET("/editblog/:id", admincontroller.EditBlogCtr)


		admin.GET("/list-cate", admincontroller.ListCateCtr)
		admin.POST("/save-edit-cate", admincontroller.SaveEditCateCtr)
		admin.GET("/edit-cate/:id", admincontroller.EditCateCtr)
		admin.GET("/add-cate", admincontroller.AddCateCtr)
		admin.POST("/save-add-cate", admincontroller.SaveAddCateCtr)

		admin.GET("/list-tag", admincontroller.ListTagCtr)


		admin.GET("/files", admincontroller.Files)
		admin.POST("/fileupload", admincontroller.FileUpload)
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
