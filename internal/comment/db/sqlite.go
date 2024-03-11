package db

import (
	"database/sql"
	"forum/internal/comment"
	"log"
)

type db struct {
	db *sql.DB
}

func NewCommentStorage(dataBase *sql.DB) comment.CommentStorage {
	return &db{
		db: dataBase,
	}
}

func (d *db) CreateComment(comment comment.Comment) (id int64, err error) {
	res, err := d.db.Exec("INSERT INTO comments (username, post_id, user_id, message, born) VALUES(?,?,?,?,?)", comment.Username, comment.PostId, comment.UserId, comment.Content, comment.CreateTime)
	if err != nil {
		log.Println(err)
		return
	}
	id, err = res.LastInsertId()
	if err != nil {
		return
	}
	return
}

func (d *db) GetAllCommentsByPostID(id int64) (cs []comment.Comment, err error) {
	rows, err := d.db.Query("SELECT * FROM comments WHERE post_id = ?", id)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var c comment.Comment
		if err = rows.Scan(&c.ID, &c.PostId, &c.UserId, &c.Username, &c.Content, &c.CreateTime); err != nil {
			log.Println(err)
			return
		}
		cs = append(cs, c)
	}
	return
}
