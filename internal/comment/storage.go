package comment

type CommentStorage interface {
	GetAllCommentsByPostID(int64) ([]Comment, error)
	CreateComment(Comment) (int64, error)
}
