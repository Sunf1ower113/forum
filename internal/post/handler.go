package post

import (
	"forum/internal/category"
	"forum/internal/comment"
	"forum/internal/handlers"
	"forum/internal/reaction"
	"forum/internal/server"

	// t "forum/internal/template"
	u "forum/internal/user"

	// "forum/internal/user"
	e "forum/pkg/error"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	PostFormError            = "Post Form Error"
	AllPostsURL              = "/"
	PostUrl                  = "/post/"
	CreatePostURL            = "/post"
	CreatePostCommentURL     = "/post/comment/"
	CreatePostReactionURL    = "/post/reaction/"
	CreateCommentReactionURL = "/post/comment/reaction/"
)

type PostHandler struct {
	s    PostService
	ss   category.CategoryService
	comm comment.CommentService
	r    reaction.ReactionService
	m    server.Server
}

var _ handlers.Handler = &PostHandler{}

func NewPostHandler(s PostService, ss category.CategoryService, comm comment.CommentService, r reaction.ReactionService, m server.Server) *PostHandler {
	return &PostHandler{
		s:    s,
		ss:   ss,
		m:    m,
		comm: comm,
		r:    r,
	}
}

func (h *PostHandler) Register(mux *http.ServeMux) {
	mux.Handle(AllPostsURL, h.m.UnautMiddleware(h.Home))
	mux.Handle(PostUrl, h.m.UnautMiddleware(h.Post))
	mux.Handle(CreatePostURL, h.m.AuthMiddleware(h.CreatePostHandler))
	mux.Handle(CreatePostCommentURL, h.m.AuthMiddleware(h.CreatePostComment))
	mux.Handle(CreateCommentReactionURL, h.m.AuthMiddleware(h.ReactComment))
	mux.Handle(CreatePostReactionURL, h.m.AuthMiddleware(h.ReactPost))
}

func (h *PostHandler) ReactPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(server.KeyUserType(server.KeyUser))
	user, ok := ctx.(u.User)
	if !ok {
		w.Header().Set("session_token", "")
		http.Redirect(w, r, "/auth/login", http.StatusMovedPermanently)
		return
	}
	if r.Method != http.MethodPost {
		log.Println("method not allowed")
		e.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Printf("%s", err.Error())
		e.ErrorHandler(w, http.StatusBadRequest)
		return
	}
	for k := range r.PostForm {
		if !(k == "status") {
			log.Println(PostFormError)
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
	}
	status, err := strconv.Atoi(r.FormValue("status"))
	if err != nil {
		e.ErrorHandler(w, http.StatusBadRequest)
		return
	}
	log.Println(status)
	if status > 1 || status < -1 {
		e.ErrorHandler(w, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/reaction/"))
	if err != nil || id < 0 {
		e.ErrorHandler(w, http.StatusNotFound)
		return
	}
	v, err := h.r.GetReactionPostByUser(int64(id), user.ID)
	if err != nil {
		log.Println(err)
		e.ErrorHandler(w, http.StatusBadRequest)
		return
	}
	if v == status {
		err = h.r.MakeReactionPost(int64(id), user.ID, 0)
		if err != nil {
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
	} else {
		h.r.MakeReactionPost(int64(id), user.ID, status)
		if err != nil {
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
	}
	http.Redirect(w, r, PostUrl+strconv.Itoa(id), http.StatusFound)
}

func (h *PostHandler) ReactComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(server.KeyUserType(server.KeyUser))
	user, ok := ctx.(u.User)
	if !ok {
		w.Header().Set("session_token", "")
		http.Redirect(w, r, "/auth/login", http.StatusMovedPermanently)
		return
	}
	if r.Method != http.MethodPost {
		log.Println("method not allowed")
		e.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Printf("%s", err.Error())
		e.ErrorHandler(w, http.StatusBadRequest)
		return
	}
	for k := range r.PostForm {
		if !(k == "status") {
			log.Println(PostFormError)
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
	}
	status, err := strconv.Atoi(r.FormValue("status"))
	if err != nil {
		e.ErrorHandler(w, http.StatusBadRequest)
		return
	}
	if status > 1 || status < -1 {
		e.ErrorHandler(w, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/comment/reaction/"))
	if err != nil || id < 0 {
		e.ErrorHandler(w, http.StatusNotFound)
		return
	}
	v, err := h.r.GetReactionCommentByUser(int64(id), user.ID)
	if err != nil {
		log.Println(err)
		e.ErrorHandler(w, http.StatusBadRequest)
		return
	}
	if v == status {
		err = h.r.MakeReactionComment(int64(id), user.ID, 0)
		if err != nil {
			log.Println(err)
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
	} else {
		h.r.MakeReactionComment(int64(id), user.ID, status)
		if err != nil {
			log.Println(err)
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
	}
	http.Redirect(w, r, PostUrl+strconv.Itoa(id), http.StatusFound)
}

func (h *PostHandler) CreatePostComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(server.KeyUserType(server.KeyUser))
	user, ok := ctx.(u.User)
	if !ok {
		w.Header().Set("session_token", "")
		http.Redirect(w, r, "/auth/login", http.StatusMovedPermanently)
		return
	}
	if r.Method != http.MethodPost {
		log.Println("method not allowed")
		e.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	sid := strings.TrimPrefix(r.URL.Path, "/post/comment/")
	id, err := strconv.Atoi(sid)
	if err != nil || id < 0 {
		e.ErrorHandler(w, http.StatusNotFound)
		return
	}
	p, err := h.s.GetPostById(int64(id))
	if err != nil {
		log.Println(err)
		e.ErrorHandler(w, http.StatusNotFound)
		return
	}
	err = r.ParseForm()
	if err != nil {
		log.Printf("%s", err.Error())
		e.ErrorHandler(w, http.StatusBadRequest)
		return
	}
	for k := range r.PostForm {
		if !(k == "comment") {
			log.Println(k)
			log.Println(PostFormError)
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
	}
	body := r.FormValue("comment")
	c := comment.Comment{
		PostId:   p.ID,
		UserId:   user.ID,
		Username: user.Username,
		Content:  body,
	}
	_, err = h.comm.CreateComment(c)
	if err != nil {
		log.Printf("%s", err.Error())
		e.ErrorHandler(w, http.StatusBadRequest)
		return
	}
	log.Println(PostUrl + sid)
	r.Method = http.MethodGet
	log.Println(r.Method)
	http.Redirect(w, r, PostUrl+sid, http.StatusFound)
}

func (h *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(server.KeyUserType(server.KeyUser))
	user, ok := ctx.(u.User)
	if !ok {
		w.Header().Set("session_token", "")
		http.Redirect(w, r, "/auth/login", http.StatusMovedPermanently)
		return
	}

	switch r.Method {
	case http.MethodGet:
		data := TemplateData{
			Template: "create_post",
			User:     user,
		}
		c, err := h.ss.GetCategory()
		if err != nil {
			log.Println(err)
		}
		var cat []string
		for _, v := range c {
			if v.CategoryName == "My posts" || v.CategoryName == "Liked posts" {
				continue
			}
			cat = append(cat, v.CategoryName)
		}
		data.Category = cat
		tmpl, err := template.ParseGlob("./templates/*.html")
		if err != nil {
			log.Println(err)
			e.ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "index", data)
		if err != nil {
			log.Println(err)
			e.ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		return
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Printf("%s", err.Error())
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		for k := range r.PostForm {
			if !(k == "title" || k == "body" || k == "category") {
				log.Println(PostFormError)
				e.ErrorHandler(w, http.StatusBadRequest)
				return
			}
		}
		title := r.FormValue("title")
		body := r.FormValue("body")
		cat := r.Form["category"]
		if len(cat) == 0 {
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		p := Post{
			Title:      title,
			Content:    body,
			Category:   cat,
			Author:     user,
			CreateTime: time.Now().Format(time.RFC822),
		}
		_, err = h.s.CreatePost(p)
		if err != nil {
			log.Println(err)
			switch err.Error() {
			case "create post was failed":
				e.ErrorHandler(w, http.StatusInternalServerError)
				return
			default:
				e.ErrorHandler(w, http.StatusBadRequest)
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)

	default:
		log.Println("method not allowed")
		e.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *PostHandler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		e.ErrorHandler(w, http.StatusNotFound)
		return
	}
	if !(r.Method == http.MethodGet || r.Method == http.MethodPost) {
		log.Println("Method not allowed")
		e.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context().Value(server.KeyUserType(server.KeyUser))
	user, ok := ctx.(u.User)
	c, err := h.ss.GetCategory()
	if err != nil {
		log.Println(err)
	}
	var cat []string
	for _, v := range c {
		if (v.CategoryName == "My posts" || v.CategoryName == "Liked posts") && !ok {
			continue
		}
		cat = append(cat, v.CategoryName)
	}
	data := TemplateData{
		Template: "allposts",
		User:     user,
		Category: cat,
	}
	var p []Post
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Printf("%s", err.Error())
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		for k := range r.PostForm {
			if k != "category" {
				log.Println(PostFormError)
				e.ErrorHandler(w, http.StatusBadRequest)
				return
			}
		}
		cat := r.Form["category"]
		if len(cat) == 0 {
			p, err = h.s.GetAllPost()
			if err != nil {
				log.Println(err)
				e.ErrorHandler(w, http.StatusInternalServerError)
				return
			}

		} else {
			for _, n := range cat {
				if n == "My posts" {
					p, err = h.s.GetPostsByUserID(user.ID)
					cat = []string{}
					break
				} else if n == "Liked posts" {
					ids, err := h.r.GetLikedPosts(user.ID)
					if err != nil {
						p = []Post{}
						cat = []string{}
						break
					}
					for _, i := range ids {
						pos, err := h.s.GetPostById(i)
						if err == nil {
							p = append(p, pos)
						}
					}
					cat = []string{}
					break
				} else {
					_, err := h.ss.GetCategoryByName(n)
					if err != nil {
						log.Println(err)
						e.ErrorHandler(w, http.StatusBadRequest)
					}
				}
			}
			if len(cat) != 0 {
				p, err = h.s.GetPostsByCategories(cat)
				log.Println(len(p))
				if err != nil {
					log.Println(err)
					e.ErrorHandler(w, http.StatusInternalServerError)
				}
			}
		}
	} else {
		p, err = h.s.GetAllPost()
		if err != nil {
			log.Println(err)
			e.ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	}
	data.Posts = p
	tmpl, err := template.ParseGlob("./templates/*.html") // template.New("").ParseFiles("./templates/html/*.html") // "templates/header.html"
	if err != nil {
		log.Println(err)
		e.ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		log.Println(err)
		e.ErrorHandler(w, http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Method not allowed")
		e.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context().Value(server.KeyUserType(server.KeyUser))
	user, ok := ctx.(u.User)
	if !ok {
		user = u.User{}
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/"))
	if err != nil || id < 0 {
		e.ErrorHandler(w, http.StatusNotFound)
		return
	}
	p, err := h.s.GetPostById(int64(id))
	if err != nil {
		log.Println(err)
		e.ErrorHandler(w, http.StatusNotFound)
		return
	}
	comms, err := h.comm.GetAllCommentsByPostID(p.ID)
	if err != nil {
		return
	}
	for _, c := range comms {
		likeCount, err := h.r.GetCommentLikeCount(c.ID)
		if err != nil {
			log.Println(err)
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		dislikeCount, err := h.r.GetCommentDislikeCount(c.ID)
		if err != nil {
			log.Println(err)
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		p.Comment = append(p.Comment, Comment{
			ID:         c.ID,
			PostId:     c.PostId,
			UserId:     c.UserId,
			Username:   c.Username,
			Content:    c.Content,
			Like:       int64(likeCount),
			Dislike:    int64(dislikeCount),
			CreateTime: c.CreateTime,
		})
	}
	likeCount, err := h.r.GetPostLikeCount(p.ID)
	if err != nil {
		log.Println(err)
		e.ErrorHandler(w, http.StatusBadRequest)
		return
	}

	dislikeCount, err := h.r.GetPostDislikeCount(p.ID)
	if err != nil {
		log.Println(err)
		e.ErrorHandler(w, http.StatusBadRequest)
		return
	}
	p.Like = int64(likeCount)
	p.Dislike = int64(dislikeCount)
	data := TemplateData{
		Template: "post",
		User:     user,
		Post:     p,
	}
	tmpl, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		log.Println(err)
		e.ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		log.Println(err)
		e.ErrorHandler(w, http.StatusInternalServerError)
		return
	}
}
