package admin

import (
	"github.com/cnmade/bsmi-mail-kernel/app/orm/model"
	"github.com/cnmade/bsmi-mail-kernel/app/utils/admin_utils"
	"github.com/cnmade/bsmi-mail-kernel/app/vo"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	vo2 "github.com/cnmade/bsmi-mail-kernel/pkg/common/vo"
	"github.com/flosch/pongo2/v4"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func ListTagCtr(c *gin.Context) {
	err, username, isAdmin := admin_utils.AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	var tags []model.Tag

	_ = common.NewDb.
		Order("total_nums desc").
		Find(&tags)

	c.HTML(200, "admin/list-tag.html",
		common.Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        common.Config.Site_name,
			"siteDescription": common.Config.Site_description,
			"tags":            tags,
			"username":        username.(string),
			"isAdmin":        isAdmin.(string),
		}))
	return
}

func SaveAddTagCtr(c *gin.Context) {
	err, _, _ := admin_utils.AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	var BI vo.Tag_vo
	err  = c.MustBindWith(&BI, binding.Form)
	if err != nil {
		common.ShowUMessage(c, &vo2.Umsg{Msg: err.Error(), Url: "/"})
		common.LogError(err)
		return
	}
	if BI.Name == "" {
		common.ShowUMessage(c, &vo2.Umsg{Msg: "标签名称不能为空", Url: "/"})
		return
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		common.LogError(err)
		common.ShowUMessage(c, &vo2.Umsg{Msg: "获取时间错误", Url: "/"})
		return
	}
	aid := getLastCateId()
	nextAid := aid + 1

	tag := model.Tag{
		ID:        nextAid,
		CreatedAt: time.Now().In(loc),
		UpdatedAt: time.Now().In(loc),
		DeletedAt: gorm.DeletedAt{},
		Name:      BI.Name,
	}

	result := common.NewDb.Create(&tag)

	if result.Error == nil {
		c.Redirect(http.StatusFound, "/admin/list-tag")
	} else {
		common.LogError(result.Error)
		common.ShowUMessage(c, &vo2.Umsg{Msg: "失败", Url: "/admin/list-tag"})
	}

}
