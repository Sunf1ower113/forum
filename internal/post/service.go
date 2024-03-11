package post

import (
	"errors"
	"forum/internal/category"
	"log"
	"strings"
)

const (
	OK                  = 200
	InternalServerError = 500
	BadRequest          = 400
)

type PostService interface {
	CreatePost(post Post) (int64, error)
	GetAllPost() ([]Post, error)
	GetPostById(postId int64) (Post, error)
	GetPostsByUserID(int64) ([]Post, error)
	GetPostIDByCategory(string) ([]int64, error)
	GetPostsByCategories([]string) ([]Post, error)
}

type postService struct {
	storage  PostStorage
	category category.CategoryService
}

func NewPostService(storage PostStorage, category category.CategoryService) PostService {
	return &postService{
		storage:  storage,
		category: category,
	}
}

func (s *postService) CreatePost(post Post) (id int64, err error) {
	title, ok := validDataString(post.Title)
	if !ok {
		return -1, errors.New("title is invalid")
	}
	post.Title = title
	text, ok := validDataString(post.Content)
	if !ok {
		return -1, errors.New("content is invalid")
	}
	log.Println([]byte(text))
	log.Println(len(text))
	post.Content = text

	for _, n := range post.Category {
		if n == "Liked post" {
			return -1 , errors.New("unavailable category")
		}
		_, err := s.category.GetCategoryByName(n)
		if err != nil {
			return -1, errors.New("unavailable category")
		}
	}

	id, err = s.storage.CreatePost(post)
	if err != nil {
		return -1, errors.New("create post was failed")
	}
	for _, c := range post.Category {
		err = s.category.CreateCategoryPostByPostID(c, id)
		if err != nil {
			return
		}
	}
	return
}

func (s *postService) GetAllPost() (posts []Post, err error) {
	posts, err = s.storage.GetAllPost()
	for i, p := range posts {
		cats, err := s.category.GetCategoryByPostID(p.ID)
		if err != nil {
			return nil, err
		}
		for _, c := range cats {
			posts[i].Category = append(posts[i].Category, c.CategoryName)
		}
	}
	return
}

func (s *postService) GetPostById(postId int64) (p Post, err error) {
	p, err = s.storage.GetPostById(postId)
	cats, err := s.category.GetCategoryByPostID(p.ID)
	if err != nil {
		return
	}
	for _, c := range cats {
		p.Category = append(p.Category, c.CategoryName)
	}
	return
}

func (s *postService) GetPostsByUserID(id int64) (p []Post, err error) {
	p, err = s.storage.GetPostsByUserID(id)
	for _, p := range p {
		cats, err := s.category.GetCategoryByPostID(p.ID)
		if err != nil {
			return nil, err
		}
		for _, c := range cats {
			p.Category = append(p.Category, c.CategoryName)
		}
	}
	return
}

func (s *postService) GetPostIDByCategory(c string) (ids []int64, err error) {
	ids, err = s.storage.GetPostIDByCategory(c)
	return
}

func (s *postService) GetPostsByCategories(cs []string) (ps []Post, err error) {
	all := make(map[int64][]string)
	var ids []int64
	for _, c := range cs {
		ids, err = s.GetPostIDByCategory(c)
		if err != nil {
			return
		}
		for _, id := range ids {
			all[id] = append(all[id], c)
		}
	}
	ids = []int64{}
	for k, v := range all {
		if len(v) == len(cs) {
			ids = append(ids, k)
		}
	}
	var p Post
	for _, id := range ids {
		p, err = s.GetPostById(id)
		if err != nil {
			return
		}
		ps = append(ps, p)
	}
	return
}

func validDataString(s string) (string, bool) {
	str := strings.Trim(s, " \n \r")
	if len(str) == 0 {
		return "", false
	}
	return str, true
}
