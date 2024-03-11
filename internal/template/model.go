package template

import (
	"forum/internal/user"
	"forum/internal/post"
)

type TemplateData struct {
	Template string
	User     user.User
	Post     post.Post
	Posts    []post.Post
	// Comments []Comment
	// Error    ErrorMsg
}
