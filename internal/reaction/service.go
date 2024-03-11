package reaction
type ReactionService interface {
	GetReactionPostByUser(int64, int64) (int, error)
	MakeReactionPost(int64, int64, int) error
	GetReactionCommentByUser(int64, int64) (int, error)
	MakeReactionComment(int64, int64, int) error
	GetPostLikeCount(int64) (int64, error)
	GetPostDislikeCount(int64) (int64, error)
	GetCommentLikeCount(int64) (int64, error)
	GetCommentDislikeCount(int64) (int64, error)
	GetLikedPosts(int64) ([]int64, error)
}

type reactionService struct {
	storage ReactionStorage
}

func NewReactionService(storage ReactionStorage) ReactionService {
	return &reactionService{
		storage: storage,
	}
}

func (s *reactionService) GetLikedPosts(userID int64) (p []int64, err error) {
	p, err = s.storage.GetLikedPosts(userID)
	return
}

func (s *reactionService) GetReactionPostByUser(postID, userID int64) (status int, err error) {
	status, err = s.storage.GetReactionPostByUser(postID, userID)
	return
}

func (s *reactionService) MakeReactionPost(postID, userID int64, status int) (err error) {
	err = s.storage.MakeReactionPost(postID, userID, status)
	return
}

func (s *reactionService) GetReactionCommentByUser(postID, userID int64) (status int, err error) {
	status, err = s.storage.GetReactionCommentByUser(postID, userID)
	return
}

func (s *reactionService) MakeReactionComment(postID, userID int64, status int) (err error) {
	err = s.storage.MakeReactionComment(postID, userID, status)
	return
}

func (s *reactionService) GetPostLikeCount(id int64) (count int64, err error) {
	count, err = s.storage.GetPostLikeCount(id)
	return
}
func (s *reactionService) GetPostDislikeCount(id int64) (count int64, err error) {
	count, err = s.storage.GetPostDislikeCount(id)
	return
}

func (s *reactionService) GetCommentLikeCount(id int64) (count int64, err error) {
	count, err = s.storage.GetCommentLikeCount(id)
	return
}
func (s *reactionService) GetCommentDislikeCount(id int64) (count int64, err error) {
	count, err = s.storage.GetCommentDislikeCount(id)
	return
}