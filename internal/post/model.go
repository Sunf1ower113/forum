package post

import (
	user "forum/internal/user"
)

type Post struct {
	ID         int64
	Title      string
	Content    string
	Category   []string
	Comment    []Comment
	Author     user.User
	Like       int64
	Dislike    int64
	CreateTime string
}

type TemplateData struct {
	Template string
	User     user.User
	Post     Post
	Posts    []Post
	Category []string
	// Comments []Comment
	// Error    ErrorMsg
}

type Comment struct {
	ID         int64
	PostId     int64
	UserId     int64
	Username   string
	Content    string
	Like       int64
	Dislike    int64
	CreateTime string
}