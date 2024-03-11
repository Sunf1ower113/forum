package post

type PostStorage interface {
	CreatePost(post Post) (int64, error)
	GetAllPost() ([]Post, error)
	GetPostById(postId int64) (Post, error)
	GetPostsByUserID(int64) ([]Post, error)
	GetPostIDByCategory(string) ([]int64, error)
	// CreateCategory(category *Category) error
	// GetCategory() ([]Category, error)
	
	// GetDisLikeStatus(postId, userId int) int
	// DeletePostDisLike(post_id, user_id int) error
	// DisLikePost(post_id, user_id, status int) error
	// GetLikedPostIdByUserId(userId int) ([]int64, error)
	// GetLikeStatus(postId, userId int) int
	// LikePost(post_id, user_id, status int) error
	// UpdatePostLikeDislike(post_id, like, dislike int) error
	// DeletePostLike(post_id, user_id int) error
	
	// GetAllCommentByPostID(postId int64) ([]Comment, error)
	// GetCommentByCommentID(commentId int64) (Comment, error)
	// CommentPost(comment Comment) error

	// GetCommentLikeStatus(comment_id, userId int) int
	// LikeComment(comment_id, user_id, status int) error
	// UpdateCommentLikeDislike(comment_id, like, dislike int) error
	// DeleteCommentLike(comment_id, user_id int) error

	// DisLikeComment(comment_id, user_id, status int) error
	// DeleteCommentDisLike(comment_id, user_id int) error
	// GetDisLikeCommentStatus(comment_id, userId int) int
}