package main

import (
	auth "forum/internal/auth"
	"forum/internal/category"
	categorydb "forum/internal/category/db"
	"forum/internal/comment"
	commentdb "forum/internal/comment/db"
	config "forum/internal/config"
	post "forum/internal/post"
	postdb "forum/internal/post/db"
	"forum/internal/reaction"
	reactiondb "forum/internal/reaction/db"
	"forum/internal/server"
	"forum/internal/session"
	sessiondb "forum/internal/session/db"
	user "forum/internal/user"
	userdb "forum/internal/user/db"
	sqlitedb "forum/pkg/client/db"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	cfg, err := config.LoadConfiguration("config.json")
	if err != nil {
		log.Fatal(err)
	}
	database, err := sqlitedb.NewDB(cfg.Database.DbDriver, cfg.Database.DbName)
	if err != nil {
		log.Fatal(err)
	}
	err = database.Ping()
	if err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}

	userStorage := userdb.NewUserStorage(database, logger)
	userService := user.NewUserService(userStorage, logger)
	userHandler := user.NewUserHandler(userService)

	sessionStorage := sessiondb.NewSessionStorage(database)
	sessionService := session.NewSessionService(sessionStorage)

	categoryStorage := categorydb.NewCategoryStorage(database)
	categoryService := category.NewCategoryService(categoryStorage)

	commentStorage := commentdb.NewCommentStorage(database)
	commentService := comment.NewCategoryService(commentStorage)

	reactionStorage := reactiondb.NewReactionStorage(database)
	reactionService := reaction.NewReactionService(reactionStorage)

	s := server.NewServer(userService, sessionService, categoryService)

	authService := auth.NewAuthService(sessionService, userService, userStorage)
	authHandler := auth.NewAuthHandler(authService, sessionService, userService, *s)

	postStorage := postdb.NewPostStorage(database)
	postService := post.NewPostService(postStorage, categoryService)
	postHandler := post.NewPostHandler(postService, categoryService, commentService, reactionService, *s)

	router := http.NewServeMux()
	userHandler.Register(router)
	authHandler.Register(router)
	postHandler.Register(router)
	start(router, cfg)
}

func start(router *http.ServeMux, cfg *config.Config) {
	log.Println("Start the application...")

	listner, err := net.Listen(cfg.Listner.Protocol, cfg.Listner.Ip+cfg.Listner.Port)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Handler:      router,
		IdleTimeout:  time.Duration(cfg.Listner.IdleTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Listner.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Listner.ReadTimeout) * time.Second,
	}

	log.Printf("Server is listening port %s:%s\n", cfg.Listner.Ip, cfg.Listner.Port)
	log.Fatal(server.Serve(listner))
}
