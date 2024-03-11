package reaction

type ReactionStorage interface {
	GetReactionPostByUser(int64, int64)	(int, error)
	MakeReactionPost(int64, int64, int) (error)
	GetReactionCommentByUser(int64, int64)	(int, error)
	MakeReactionComment(int64, int64, int) (error)
	GetPostLikeCount(int64) (int64, error)
	GetPostDislikeCount(int64) (int64, error)
	GetCommentLikeCount(int64) (int64, error)
	GetCommentDislikeCount(int64) (int64, error)
	GetLikedPosts(int64) ([]int64, error)
}