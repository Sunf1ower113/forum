package comment

type Comment struct {
	ID         int64
	PostId     int64
	UserId     int64
	Username   string
	Content    string
	Like       int
	Dislike    int
	CreateTime string
}
