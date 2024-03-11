package server

import (
	// "forum/internal/auth"
	"forum/internal/category"
	"forum/internal/session"
	"forum/internal/user"
)

type Server struct {
	userService     user.UserService
	sessionService  session.SessionService
	categoryService category.CategoryService
	// authService    auth.AuthService
}

func NewServer(userService user.UserService, sessionService session.SessionService, categoryService category.CategoryService) *Server {
	return &Server{
		userService:     userService,
		sessionService:  sessionService,
		categoryService: categoryService,
	}
}
