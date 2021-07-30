package admincontroller

import (
	"errors"
	"github.com/cnmade/bsmi-mail-kernel/app/orm/model"
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

func SaveEditCateCtr(c *gin.Context) {
	err, _, _ := AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	var BI vo.Category_vo
	err = c.MustBindWith(&BI, binding.Form)
	if err != nil {
		common.LogError(err)
		common.ShowUMessage(c, &vo2.Umsg{Msg: err.Error(), Url: "javascript:history.go(-1)"})
		return
	}
	if BI.CateId == 0 {
		common.ShowUMessage(c, &vo2.Umsg{Msg: "分类未找到", Url: "javascript:history.go(-1)"})
		return
	}
	if BI.Name == "" {
		common.ShowUMessage(c, &vo2.Umsg{Msg: "分类名不能为空", Url: "javascript:history.go(-1)"})
		return
	}

	var cate model.Category

	result := common.NewDb.Find(&cate, BI.CateId)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		common.ShowMessage(c, &vo2.Msg{
			Msg: "文章不存在",
		})
		return
	}

	cate.ID = BI.CateId
	cate.Name = BI.Name
	common.NewDb.
		Where("ID = ? ", cate.ID).
		Save(cate)

	c.Redirect(http.StatusFound, "/admin/list-cate")

}

func ListCateCtr(c *gin.Context) {
	err, username, isAdmin := AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	var categories []model.Category

	_ = common.NewDb.
		Find(&categories)

	c.HTML(200, "admin/list-cate.html",
		common.Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        common.Config.Site_name,
			"siteDescription": common.Config.Site_description,
			"categories":      categories,
			"username":        username.(string),
			"isAdmin":        isAdmin.(string),
		}))
	return
}

func AddCateCtr(c *gin.Context) {
	err, username, isAdmin := AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}

	c.HTML(200, "admin/add-cate.html",
		common.Pongo2ContextWithVersion(pongo2.Context{
			"username":        username.(string),
			"isAdmin":        isAdmin.(string),
			"siteName":        common.Config.Site_name,
			"siteDescription": common.Config.Site_description,
		}))
	return
}

func SaveAddCateCtr(c *gin.Context) {
	err, _, _ := AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	var BI vo.Category_vo
	err  = c.MustBindWith(&BI, binding.Form)
	if err != nil {
		common.ShowUMessage(c, &vo2.Umsg{Msg: err.Error(), Url: "/"})
		common.LogError(err)
		return
	}
	if BI.Name == "" {
		common.ShowUMessage(c, &vo2.Umsg{Msg: "分类名称不能为空", Url: "/"})
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

	cate := model.Category{
		ID:        nextAid,
		CreatedAt: time.Now().In(loc),
		UpdatedAt: time.Now().In(loc),
		DeletedAt: gorm.DeletedAt{},
		Name:      BI.Name,
	}

	result := common.NewDb.Create(&cate)

	if result.Error == nil {
		c.Redirect(http.StatusFound, "/admin/list-cate")
	} else {
		common.LogError(result.Error)
		common.ShowUMessage(c, &vo2.Umsg{Msg: "失败", Url: "/admin/list-cate"})
	}

}

func EditCateCtr(c *gin.Context) {
	err, username, isAdmin := AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	id := c.Param("id")
	var cate model.Category
	result := common.NewDb.First(&cate, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		common.ShowMessage(c, &vo2.Msg{
			Msg: "文章不存在",
		})
		return
	}

	c.HTML(200, "admin/edit-cate.html",
		common.Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        common.Config.Site_name,
			"siteDescription": common.Config.Site_description,
			"username":        username.(string),
			"isAdmin":        isAdmin.(string),
			"cateId":          cate.ID,
			"name":            cate.Name,
		}))
	return
}

func getLastCateId() int64 {
	var cate model.Category
	result := common.NewDb.Last(&cate)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0
	}
	return cate.ID
}

