package controller

import (
	"encoding/json"
	"fmt"
	"github.com/cnmade/bsmi-mail-kernel/app/orm/model"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"time"
)

// RSS controll group
type RSS struct {
}

func (rss *RSS) Alter(c *gin.Context) {
	c.Redirect(301, "/rss")
}

// Out Render and output RSS
//TODO 增加缓存
func (rss *RSS) Out(c *gin.Context) {

	_, username, _ := UserPermissionCheckDefaultAllow(c)
	if common.Config.PrivateMode == 1 {

		if username == "" {
			c.Redirect(301, "/admin/login")
			return
		}
	}
	var blogItems []model.Article
	result := common.NewDb.Limit(20).Order("aid desc").Find(&blogItems)
	if result.Error != nil {
		common.LogError(result.Error)
		return
	}
	common.Sugar.Info(json.Marshal(c.Request))
	hostname := common.Config.Site_url
	now := time.Now()
	feed := &feeds.Feed{
		Title:       common.Config.Site_name,
		Link:        &feeds.Link{Href: hostname},
		Description: common.Config.Site_description,
		Created:     now,
	}
	feed.Items = make([]*feeds.Item, 0)
	for _, blog := range blogItems {

		itemTime, _ := time.Parse("2006-01-02 15:04:05", blog.PublishTime)
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       blog.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/#view/%d", hostname, blog.Aid)},
			Description: blog.Content,
			Created:     itemTime,
		})
	}
	c.XML(http.StatusOK, (&feeds.Atom{feed}).FeedXml())
}
