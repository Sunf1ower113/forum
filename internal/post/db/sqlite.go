package db

import (
	"database/sql"
	"forum/internal/post"
	"log"
)

type db struct {
	db *sql.DB
}

func NewPostStorage(dataBase *sql.DB) post.PostStorage {
	return &db{
		db: dataBase,
	}
}

func (d *db) CreatePost(post post.Post) (id int64, err error) {
	res, err := d.db.Exec("INSERT INTO posts (title, message, user_id, username, born) VALUES(?,?,?,?,?)", post.Title, post.Content, post.Author.ID, post.Author.Username, post.CreateTime)
	if err != nil {
		log.Println(err)
		return
	}
	id, err = res.LastInsertId()
	if err != nil {
		return
	}
	post.ID =id
	return
}

func (d *db) GetAllPost() (p []post.Post, err error) {
	rows, err := d.db.Query("SELECT * FROM posts ORDER BY born DESC;")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var post post.Post
		if err = rows.Scan(&post.ID, &post.Author.ID, &post.Author.Username, &post.Title, &post.Content, &post.CreateTime); err != nil {
			return
		}
		p = append(p, post)
	}
	return
}

func (d *db) GetPostById(postID int64) (p post.Post, err error) {
	row := d.db.QueryRow("SELECT post_id, title, message, user_id, username, born FROM posts WHERE post_id = ? ", postID)
	if err = row.Scan(&p.ID, &p.Title, &p.Content, &p.Author.ID, &p.Author.Username, &p.CreateTime); err != nil {
		log.Println(err)
		return
	}
	return
}

func (d *db) GetPostsByUserID(UserID int64) (p []post.Post, err error) {
	rows, err := d.db.Query("SELECT * FROM posts WHERE user_id = ?", UserID)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var post post.Post
		if err = rows.Scan(&post.ID, &post.Author.ID, &post.Author.Username, &post.Title, &post.Content, &post.CreateTime); err != nil {
			return
		}
		p = append(p, post)
	}
	return
}

func (d *db) GetPostIDByCategory(c string) (id []int64, err error) {
	query := `SELECT p.post_id FROM posts p
    JOIN categories_posts pc ON p.post_id = pc.post_id
    WHERE pc.category = ?`

	rows, err := d.db.Query(query, c)
	if err != nil {
		log.Println(err)
		return
	}
	s := struct {
		id int64
	}{}
	for rows.Next() {
		err = rows.Scan(&s.id)
		if err != nil {
			log.Println(err)
			return
		}
		id = append(id, s.id)
	}
	return
}
