package vo

// Blog_vo is the blog item
type Category_vo struct {
	CateId int64 `form:"cateId"`
	Name   string `form:"name" binding:"required"`
}
