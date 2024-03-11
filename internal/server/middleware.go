package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type KeyUserType string

const KeyUser = "user"

func (s *Server) AuthMiddleware(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/auth/register", http.StatusMovedPermanently)
			return
		}
		session, err := s.sessionService.GetSessionByToken(cookie.Value)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/auth/login", http.StatusMovedPermanently)
			return
		}
		if session.Expiry.Before(time.Now()) {
			log.Println("Session expire")
			http.Redirect(w, r, "/auth/login", http.StatusMovedPermanently)
			return
		}
		user, err := s.userService.GetUserByToken(cookie.Value)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/auth/login", http.StatusMovedPermanently)
			return
		}
		ctx := context.WithValue(r.Context(), KeyUserType(KeyUser), user)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) UnautMiddleware(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			ctx := context.WithValue(r.Context(), KeyUserType(KeyUser), nil)
			handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		session, err := s.sessionService.GetSessionByToken(cookie.Value)
		if err != nil {
			ctx := context.WithValue(r.Context(), KeyUserType(KeyUser), nil)
			handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		if session.Expiry.Before(time.Now()) {
			ctx := context.WithValue(r.Context(), KeyUserType(KeyUser), nil)
			handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		user, err := s.userService.GetUserByToken(cookie.Value)
		if err != nil {
			fmt.Println(err)
			ctx := context.WithValue(r.Context(), KeyUserType(KeyUser), nil)
			handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx := context.WithValue(r.Context(), KeyUserType(KeyUser), user)

		handler.ServeHTTP(w, r.WithContext(ctx))
		// http.Redirect(w, r, "/", http.StatusFound)
	})
}
