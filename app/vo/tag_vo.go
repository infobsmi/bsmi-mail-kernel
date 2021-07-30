package vo

// Blog_vo is the blog item
type Tag_vo struct {
	TagId int64 `form:"tagId"`
	Name   string `form:"name" binding:"required"`
}
