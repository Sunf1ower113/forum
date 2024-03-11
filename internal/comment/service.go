package comment

import (
	"errors"
	"log"
	"strings"
)

type CommentService interface {
	GetAllCommentsByPostID(int64) ([]Comment, error)
	CreateComment(Comment) (int64, error)
}

type commentService struct {
	storage CommentStorage
}

func NewCategoryService(storage CommentStorage) CommentService {
	return &commentService{
		storage: storage,
	}
}

func (s *commentService) GetAllCommentsByPostID(id int64) (cs []Comment, err error) {
	cs, err = s.storage.GetAllCommentsByPostID(id)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (s *commentService) CreateComment(comment Comment) (id int64, err error) {
	var str string
	str, ok := validDataString(comment.Content)
	if !ok {
		return -1, errors.New("comment message is invalid")
	}
	comment.Content = str
	id, err = s.storage.CreateComment(comment)
	return
}

func validDataString(s string) (string, bool) {
	str := strings.Trim(s, " \n \r")
	if len(str) == 0 {
		return "", false
	}
	return str, true
}
