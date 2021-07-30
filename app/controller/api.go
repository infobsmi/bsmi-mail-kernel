package controller

import (
	"github.com/cnmade/bsmi-mail-kernel/app/orm/model"
	"github.com/cnmade/bsmi-mail-kernel/app/vo"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	"github.com/gin-gonic/gin"
	"strings"

	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
)

type Api struct {
}

type apiBlogList struct {
	Aid   string `form:"aid" json:"aid"  binding:"required"`
	Title string `form:"title" json:"title"  binding:"required"`
}

func (a *Api) NavAll(c *gin.Context) {

	_, username, _ := UserPermissionCheckDefaultAllow(c)
	if common.Config.PrivateMode == 1 {

		if username == "" {
			c.Redirect(301, "/admin/login")
			return
		}
	}
	var articleList []model.Article
	common.NewDb.Where("p_aid = 0").
		Order("sort_id asc").
		Find(&articleList)

	var na []vo.Nav_item
	for _, s := range articleList {
		na = append(na, vo.Nav_item{
			Name:         s.Title,
			Id:           uint64(s.Aid),
			LoadOnDemand: true,
		})
	}
	c.JSON(http.StatusOK, na)
}

func (a *Api) NavLoad(c *gin.Context) {

	_, username, _ := UserPermissionCheckDefaultAllow(c)
	if common.Config.PrivateMode == 1 {

		if username == "" {
			c.Redirect(301, "/admin/login")
			return
		}
	}
	rawAid := c.Query("node")
	if rawAid == "" {
		c.JSON(http.StatusOK, []string{})
		return
	}
	aid, err := strconv.Atoi(rawAid)
	fmt.Println(aid)
	if err != nil {
		common.Sugar.Fatal(err)
		c.JSON(http.StatusOK, []string{})
		return
	}

	var articleList []model.Article
	common.NewDb.Where("p_aid = ?", aid).
		Order("sort_id asc").
		Find(&articleList)

	var na []vo.Nav_item
	for _, s := range articleList {
		na = append(na, vo.Nav_item{
			Name:         s.Title,
			Id:           uint64(s.Aid),
			LoadOnDemand: true,
		})
	}
	if len(na) <= 0 {
		c.JSON(http.StatusOK, []string{})
	} else {
		c.JSON(http.StatusOK, na)
	}
}

func (a *Api) Resort(c *gin.Context) {

	_, username, _ := UserPermissionCheckDefaultAllow(c)
		if username == "" {
			c.Redirect(301, "/admin/login")
			return
		}
	var req vo.Resort_req

	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.Sugar.Error(err)
		c.JSON(http.StatusOK, gin.H{"msg": "必须是正确的数据结构"})
		return
	}
	common.Sugar.Infof("req: %+v", req)
	common.NewDb.Model(&model.Article{}).
		Where("aid = ?", req.MoveNodeId).
		Update("p_aid", req.NewPaid)

	spl := strings.Split(req.NewSort, ",")

	if len(spl) > 0 {
		for i, s := range spl {

			if s != "" {
				common.NewDb.Model(&model.Article{}).
					Where("aid = ?", s).
					Update("sort_id", i)
			}
		}

	}

	c.JSON(http.StatusOK, []string{})
}

type apiBlogItem struct {
	Aid     string `form:"aid" json:"aid"  binding:"required"`
	Title   string `form:"title" json:"title"  binding:"required"`
	Content string `form:"content" json:"content"  binding:"required"`
}

