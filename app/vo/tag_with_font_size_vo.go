package vo

// Blog_vo is the blog item
type Tag_with_font_size_vo struct {
	TagId int64 `form:"tagId"`
	Name   string `form:"name" binding:"required"`
	TotalNums int64
	FontSize int
}
