package admincontroller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cnmade/bsmi-mail-kernel/app/orm/model"
	"github.com/cnmade/bsmi-mail-kernel/app/service/category_service"
	"github.com/cnmade/bsmi-mail-kernel/app/service/tag_service"
	"github.com/cnmade/bsmi-mail-kernel/app/vo"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	vo2 "github.com/cnmade/bsmi-mail-kernel/pkg/common/vo"
	"github.com/flosch/pongo2/v4"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func DeleteBlogCtr(c *gin.Context) {
	err, _, _ := AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	var BI vo.EditBlog_vo
	err  = c.MustBindWith(&BI, binding.Form)
	if err != nil {
		common.LogError(err)
		common.ShowUMessage(c, &vo2.Umsg{Msg: err.Error(), Url: "/"})
		return
	} else {
		if BI.Aid == 0 {
			common.ShowUMessage(c, &vo2.Umsg{Msg: "文章未找到", Url: "/"})
			return
		}
		common.NewDb.Delete(&model.Article{}, BI.Aid)
		common.ShowUMessage(c, &vo2.Umsg{Msg: "删除成功", Url: "/"})
	}
}

// ListBlogCtr is list blogs for admin
func ListBlogCtr(c *gin.Context) {
	err, username, isAdmin := AdminPermissionCheck(c)
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
	var blogDataList []model.Article

	result := common.NewDb.
		Limit(rpp).
		Offset(offset).
		Order("aid desc").
		Find(&blogDataList)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		common.ShowMessage(c, &vo2.Msg{
			Msg: "文章不存在",
		})
		return
	}

	c.HTML(200, "admin/list-blog.html",
		common.Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        common.Config.Site_name,
			"siteDescription": common.Config.Site_description,
			"blogDataList":    blogDataList,
			"username":        username.(string),
			"isAdmin":        isAdmin.(string),
			"prevPage":        fmt.Sprintf("%d", prev_page),
			"nextPage":        fmt.Sprintf("%d", next_page),
		}))
	return
}

func EditBlogCtr(c *gin.Context) {
	err, username, isAdmin := AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	id := c.Param("id")
	var blogItem model.Article
	result := common.NewDb.First(&blogItem, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		common.ShowMessage(c, &vo2.Msg{
			Msg: "文章不存在",
		})
		return
	}

	categories := category_service.GetCategories()

	var tagStr []string
	var tagIds []int64

	err = json.Unmarshal(blogItem.TagIds, &tagIds)
	if err != nil {
		common.Sugar.Infof(" json decode error %+v", err)
	} else {
		tmpTags := tag_service.BatchGetTagName(tagIds)
		for _, vt := range tmpTags {
			tagStr = append(tagStr, vt.Name)
		}
	}
	c.HTML(200, "admin/edit-blog.html",
		common.Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        common.Config.Site_name,
			"siteDescription": common.Config.Site_description,
			"username":        username.(string),
			"isAdmin":        isAdmin.(string),
			"aid":             fmt.Sprintf("%d", blogItem.Aid),
			"title":           blogItem.Title,
			"content":         blogItem.Content,
			"publishTime":     blogItem.PublishTime,
			"tags":            strings.Join(tagStr, ","),
			"paid": blogItem.PAid,
			"categories":      categories,
			"cateId":      blogItem.CateId,
			"views":           fmt.Sprintf("%d", blogItem.Views),
		}))
	return
}

func SaveBlogEditCtr(c *gin.Context) {
	err, _, _ := AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	var BI vo.EditBlog_vo
	err  = c.MustBindWith(&BI, binding.Form)
	if err != nil {
		common.LogError(err)
		common.ShowUMessage(c, &vo2.Umsg{Msg: err.Error(), Url: "javascript:history.go(-1)"})
		return
	}
	if BI.Aid == 0 {
		common.ShowUMessage(c, &vo2.Umsg{Msg: "文章未找到", Url: "javascript:history.go(-1)"})
		return
	}
	if BI.Title == "" {
		common.ShowUMessage(c, &vo2.Umsg{Msg: "标题不能为空", Url: "javascript:history.go(-1)"})
		return
	}
	if BI.Content == "" {
		common.ShowUMessage(c, &vo2.Umsg{Msg: "内容不能为空", Url: "javascript:history.go(-1)"})
		return
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		common.LogError(err)
		common.ShowUMessage(c, &vo2.Umsg{Msg: "获取时间错误", Url: "/"})
		return
	}

	var blogItem model.Article

	result := common.NewDb.Find(&blogItem, BI.Aid)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		common.ShowMessage(c, &vo2.Msg{
			Msg: "文章不存在",
		})
		return
	}


	//保存 文章历史


	ahItem := model.ArticleHistory{
		Aid:           blogItem.Aid,
		Title:         blogItem.Title,
		Content:       blogItem.Content,
		PublishTime:   blogItem.PublishTime,
		UpdateTime:   blogItem.UpdateTime,
		PublishStatus: 1,
		CateId:        blogItem.CateId,
		TagIds:        blogItem.TagIds,
		PAid: blogItem.PAid,
		SortId: blogItem.SortId,
	}

	_ = common.NewDb.Create(&ahItem)


	//处理保存编辑帖子

	tagIdStr, tagIds := processTags(BI.Tags)

	blogItem.Aid = BI.Aid
	blogItem.Title = BI.Title
	blogItem.CateId = BI.CateId
	blogItem.Content = BI.Content
	blogItem.TagIds = tagIdStr
	blogItem.PAid = BI.PAid
	blogItem.UpdateTime = time.Now().In(loc).Format("2006-01-02 15:04:05")
	common.NewDb.
		Where("aid = ?", blogItem.Aid).
		Save(blogItem)

	for _, tmpTagId := range tagIds {
		tag_service.RefreshCountOfArticle(tmpTagId)
	}

	CKey := fmt.Sprintf("blogitem-%d", BI.Aid)
	common.LogInfo("Remove cache Key:" + CKey)

	c.Redirect(http.StatusFound, fmt.Sprintf("/#view/%d", BI.Aid))

}

func AddBlogCtr(c *gin.Context) {
	paid, err := strconv.Atoi(c.DefaultQuery("paid", "0"))
	if err != nil {
		common.Sugar.Fatal(err)
	}
	err, _, _ = AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}

	categories := category_service.GetCategories()
	c.HTML(200, "admin/add-blog.html",
		common.Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        common.Config.Site_name,
			"siteDescription": common.Config.Site_description,
			"categories":      categories,
			"paid": paid,
		}))
	return
}

func SaveBlogAddCtr(c *gin.Context) {
	err, _, _ := AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	var BI vo.Blog_vo
	err = c.MustBindWith(&BI, binding.Form)
	if err != nil {
		common.ShowUMessage(c, &vo2.Umsg{Msg: err.Error(), Url: "/"})
		common.LogError(err)
		return
	}
	if BI.Title == "" {
		common.ShowUMessage(c, &vo2.Umsg{Msg: "标题不能为空", Url: "/"})
		return
	}
	if BI.Content == "" {
		common.ShowUMessage(c, &vo2.Umsg{Msg: "内容不能为空", Url: "/"})
		return
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		common.LogError(err)
		common.ShowUMessage(c, &vo2.Umsg{Msg: "获取时间错误", Url: "/"})
		return
	}
	aid := getLastAid()
	nextAid := aid + 1

	tagIdStr, tagIds := processTags(BI.Tags)

	blogItem := model.Article{
		Aid:           nextAid,
		Title:         BI.Title,
		Content:       BI.Content,
		PublishTime:   time.Now().In(loc).Format("2006-01-02 15:04:05"),
		UpdateTime:   time.Now().In(loc).Format("2006-01-02 15:04:05"),
		PublishStatus: 1,
		CateId:        BI.CateId,
		TagIds:        tagIdStr,
		PAid: BI.PAid,
	}

	result := common.NewDb.Create(&blogItem)


	for _, tmpTagId := range tagIds {
		tag_service.RefreshCountOfArticle(tmpTagId)
	}
	if result.Error == nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/#view/%d", nextAid))
	} else {
		common.LogError(result.Error)
		common.ShowUMessage(c, &vo2.Umsg{Msg: "失败", Url: "/"})
	}

}

func processTags(tmpTags string) ([]byte, []int64) {
	var tagIds []int64
	var tagIdStr []byte
	if len(tmpTags) > 0 {
		tmpTags = strings.ReplaceAll(tmpTags, "，", ",")
		tagSplited := strings.Split(tmpTags, ",")
		for _, tagName := range tagSplited {
			tagName = strings.TrimSpace(tagName)
			common.Sugar.Info("tagName: %s", tagName)
			tagIds = append(tagIds, tag_service.GetOrNewTagId(tagName))
		}

		tagIdStr, _ = json.Marshal(tagIds)
		common.Sugar.Infof("tagIdStr: %v", tagIdStr)
	}
	return tagIdStr, tagIds
}

func getLastAid() int64 {
	var blogItem model.Article
	result := common.NewDb.Last(&blogItem)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0
	}
	return blogItem.Aid
}
