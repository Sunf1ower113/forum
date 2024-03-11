package user

import (
	handlers "forum/internal/handlers"
	e "forum/pkg/error"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	OK                  = "OK"
	MethodNotAllowed    = "Method Not Allowed"
	BadRequest          = "Bad Request"
	InternalServerError = "Internal Server Error"
	NotFound            = "Not Found"

	PostFormError = "Post Form Error"
	usersURL      = "/users"
	userURL       = "/users/"
)

type UserHandler struct {
	s UserService
}

var _ handlers.Handler = &UserHandler{}

func NewUserHandler(s UserService) *UserHandler {
	return &UserHandler{
		s: s,
	}
}

func (h *UserHandler) Register(mux *http.ServeMux) {
	files := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc(usersURL, h.GetUserByID)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/users/"))
		if err != nil {
			e.ErrorHandler(w, http.StatusNotFound)
			return
		}
		_, err = h.s.GetUserByID(int64(id))
		if err != nil {
			e.ErrorHandler(w, http.StatusNotFound)
			return
		}
	default:
		log.Println("method not allowed")
		e.ErrorHandler(w, http.StatusMethodNotAllowed)
	}
}
