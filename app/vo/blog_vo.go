package vo

type Blog_vo struct {
	Title   string `form:"title" binding:"required"`
	Content string `form:"content" binding:"required"`
	CateId  int64  `form:"cateId" `
	Tags    string `form:"tags" `
	PAid    int64  `form:"paid" `
}

type EditBlog_vo struct {
	Aid     int64  `form:"aid" binding:"required"`
	Title   string `form:"title" binding:"required"`
	Content string `form:"content" binding:"required"`
	CateId  int64  `form:"cateId" `
	Tags    string  `form:"tags" `
	PAid    int64  `form:"paid" `
}

