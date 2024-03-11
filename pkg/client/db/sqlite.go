package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB(driver, name string) (db *sql.DB, err error) {
	db, err = sql.Open(driver, name)
	if err != nil {
		log.Println("error create db")
		return
	}
	log.Println("create db: OK")
	err = createTable(db)
	if err != nil {
		log.Println("error create tables")
		return
	}
	return
}

func createTable(db *sql.DB) error {
	query := []string{}
	users := `
	CREATE TABLE IF NOT EXISTS users(
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	)
	`
	posts := `
	CREATE TABLE IF NOT EXISTS posts(
		post_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		title TEXT NOT NULL,
		message TEXT NOT NULL,
		born TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(user_id)
	);
	`
	session := `
	CREATE TABLE IF NOT EXISTS sessions(
		user_id INTEGER NOT NULL UNIQUE,
		token TEXT NOT NULL,
		expiry DATE NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
	)
	`
	comments := `
	CREATE TABLE IF NOT EXISTS comments(
		comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_Id INTEGER NOT NULL,
		username TEXT NOT NULL,
		message TEXT NOT NULL,
		born TEXT NOT NULL
	)
	`
	post_reactions := `
	CREATE TABLE IF NOT EXISTS post_reactions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		post_id INTEGER,
		status INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		FOREIGN KEY(post_id) REFERENCES posts(post_id) ON DELETE CASCADE
	);
	`

	comment_reactions := `
	CREATE TABLE IF NOT EXISTS comment_reactions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		comment_id INTEGER,
		status INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		FOREIGN KEY(comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE
	);
	`
	
	categories := `
	CREATE TABLE IF NOT EXISTS categories(
		category TEXT NOT NULL UNIQUE PRIMARY KEY
	);
	INSERT OR IGNORE INTO categories (category)
	VALUES("Backend"), ("Frontend"), ("UX/UI"), ("QA"), ("System dev"), ("Other"), ("My posts"), ("Liked posts");
	`
	categoriesPosts := ` CREATE TABLE IF NOT EXISTS categories_posts(
		category TEXT,
		post_id INTEGER,
		FOREIGN KEY(category) REFERENCES categories(category) ON DELETE CASCADE,
		FOREIGN KEY(post_id) REFERENCES posts(post_id) ON DELETE CASCADE
	)
	`
	query = append(query, users, posts, session, comments, categories, categoriesPosts, post_reactions, comment_reactions)
	for _, v := range query {
		_, err := db.Exec(v)
		if err != nil {
			return err
		}
	}
	return nil
}
