package auth

import (
	"forum/internal/handlers"
	"forum/internal/server"
	"forum/internal/session"
	t "forum/internal/template"
	"forum/internal/user"
	e "forum/pkg/error"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

const (
	LoginURL      = "/auth/login"
	SignUpURL     = "/auth/register"
	LogoutURL     = "/auth/logout"
	PostFormError = "Post Form Error"
)

type AuthHandler struct {
	ss session.SessionService
	s  AuthService
	u  user.UserService
	m  server.Server
}

var (
	Sessions []session.Session
	_        handlers.Handler = &AuthHandler{}
)

func NewAuthHandler(s AuthService, ss session.SessionService, u user.UserService, m server.Server) *AuthHandler {
	return &AuthHandler{
		s:  s,
		ss: ss,
		u:  u,
		m:  m,
	}
}

func (h *AuthHandler) Register(mux *http.ServeMux) {
	mux.Handle(LoginURL, h.m.UnautMiddleware(h.Login))
	mux.Handle(SignUpURL, h.m.UnautMiddleware(h.SignUp))
	mux.HandleFunc(LogoutURL, h.Logout)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(server.KeyUserType(server.KeyUser))
	us, ok := ctx.(user.User)
	if !ok {
		us = user.User{}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	switch r.Method {
	case http.MethodGet:

		data := t.TemplateData{
			Template: "sign-in",
			User:     us,
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

	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Printf("%s", err.Error())
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		for k := range r.PostForm {
			if !(k == "email" || k == "password") {
				log.Println(PostFormError)
				e.ErrorHandler(w, http.StatusBadRequest)
				return
			}
		}
		email := strings.ToLower(r.FormValue("email"))

		password := r.FormValue("password")
		u := user.User{
			Username: "",
			Email:    email,
			Password: password,
		}

		session, err := h.s.Login(u)
		if err != nil {
			// Messages.Message = "Wrong password or email"
			http.Redirect(w, r, LoginURL, http.StatusFound)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   session.Token,
			Expires: session.Expiry,
			Path:    "/",
		})
		Sessions = append(Sessions, session)
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		log.Println("Method not allowed")
		e.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(server.KeyUserType(server.KeyUser))
	us, ok := ctx.(user.User)
	if ok {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		us = user.User{}
	}
	switch r.Method {
	case "GET":
		data := t.TemplateData{
			Template: "sign-up",
			User:     us,
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
	case "POST":
		err := r.ParseForm()
		if err != nil {
			log.Printf("%s", err.Error())
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		for k := range r.PostForm {
			if !(k == "username" || k == "email" || k == "password") {
				log.Println(PostFormError)
				e.ErrorHandler(w, http.StatusBadRequest)
				return
			}
		}
		username := r.FormValue("username")
		email := strings.ToLower(r.FormValue("email"))
		password := r.FormValue("password")
		u := user.User{
			Username: username,
			Email:    email,
			Password: password,
		}
		_, err = h.s.SignUp(u)
		if err != nil {
			e.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/auth/login", http.StatusFound)
	default:
		e.ErrorHandler(w, http.StatusMethodNotAllowed)
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not allowed")
		e.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	c, err := r.Cookie("session_token")
	if err == nil {
		h.s.Logout(c.Value)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *AuthHandler) ClearSession() {
	all, err := h.ss.GetAllSessionsTime()
	if err != nil {
		log.Println("error when get all session time", err.Error())
	}
	Sessions = all
	for {
		time.Sleep(time.Second)
		for i, v := range Sessions {
			if v.Expiry.Before(time.Now()) {
				err := h.ss.DeleteSession(v.Token)
				if i == len(Sessions)-1 {
					Sessions = Sessions[:i]
				} else {
					Sessions = append(Sessions[:i], Sessions[i+1:]...)
				}
				if err != nil {
					log.Println("session delete was fail", err.Error())
				} else {
					log.Printf("session for %s was delete\n", v.Username)
				}
			}
		}
	}
}
