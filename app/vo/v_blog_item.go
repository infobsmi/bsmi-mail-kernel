package vo

import "database/sql"

type VBlogItem struct {
	Aid            int
	Title          sql.NullString
	Content        sql.NullString
	Publish_time   sql.NullString
	Publish_status sql.NullInt64
	Views          int
}

