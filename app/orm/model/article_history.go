package model
import  "gorm.io/datatypes"
type ArticleHistory struct {
	Id int64 `gorm:"primaryKey" json:"id"`
	Aid int64 `json:"aid"`
	CateId int64 `gorm:"default:0" json:"cate_id"`
	Title string `json:"title"`
	Content string `json:"content"`
	PublishTime string `json:"publish_time"`
	UpdateTime string `json:"update_time"`
	PublishStatus int `json:"publish_status"`
	Views int64 `json:"views"`
	TagIds datatypes.JSON `json:"tag_ids"`
	PAid int64 `json:"p_aid"`
	SortId int64 `json:"sort_id"`
}


func (ArticleHistory) TableName()  string {
	return "bk_article_history"
	
}