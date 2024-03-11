package db

import (
	"database/sql"
	"forum/internal/reaction"
	"log"
	"strconv"
)

type db struct {
	db *sql.DB
}

func NewReactionStorage(dataBase *sql.DB) reaction.ReactionStorage {
	return &db{
		db: dataBase,
	}
}

func (d *db) GetReactionPostByUser(postID, userID int64) (status int, err error) {
	query := `
	SELECT * FROM post_reactions 
	WHERE user_id = ? AND post_id = ?
	`
	r := reaction.Reaction{}
	row := d.db.QueryRow(query, userID, postID)
	err = row.Scan(&r.ID, &r.UserId, &r.PostId, &r.Status)
	switch {
	case err == sql.ErrNoRows:
			return 0, nil
	case err != nil:
		log.Println(err)
			return
	}
	status = r.Status
	return
}

func (d *db) MakeReactionPost(postID, userID int64, status int) (err error) {
	q := `SELECT * FROM post_reactions
	WHERE user_id = ? AND post_id = ?`
	r := reaction.Reaction{}
	row := d.db.QueryRow(q, userID, postID)
	err = row.Scan(&r.ID, &r.UserId, &r.PostId, &r.Status)
	switch{

	case err == sql.ErrNoRows:
		query := `
		INSERT INTO post_reactions (status, user_id, post_id) VALUES
		(?, ?, ?)`
		_, err1 := d.db.Exec(query, status, userID, postID)
		if err1 != nil {
			log.Println(err1)
			return
		}
	case err != nil:
		log.Println(err)
			return
	}
	query := `
	DELETE FROM post_reactions
	WHERE post_id = ? AND user_id = ?
	`
	_, err = d.db.Exec(query, postID, userID)
	if err != nil {
		log.Println(err)
		return
	}
	query = `
		INSERT INTO post_reactions (status, user_id, post_id) VALUES
		(?, ?, ?)`
		_, err1 := d.db.Exec(query, status, userID, postID)
		if err1 != nil {
			log.Println(err1)
			return
		}
	return
}

func (d *db) GetReactionCommentByUser(commentId, userID int64) (status int, err error) {
	query := `
	SELECT id, user_id, comment_id, status FROM comment_reactions 
	WHERE user_id = ? AND comment_id = ?
	`
	r := reaction.ReactionComment{}
	row := d.db.QueryRow(query, userID, commentId)
	err = row.Scan(&r.ID, &r.UserId, &r.CommentId, &r.Status) 
	switch {
	case err == sql.ErrNoRows:
			return 0, nil
	case err != nil:
			return
	}
	status = r.Status
	return
}

func (d *db) MakeReactionComment(postID, userID int64, status int) (err error) {
	q := `SELECT * FROM comment_reactions
	WHERE user_id = ? AND comment_id = ?`
	r := reaction.Reaction{}
	row := d.db.QueryRow(q, userID, postID)
	err = row.Scan(&r.ID, &r.UserId, &r.PostId, &r.Status)
	switch{

	case err == sql.ErrNoRows:
		log.Println("asdsad")
		query := `
		INSERT INTO comment_reactions (status, user_id, comment_id) VALUES
		(?, ?, ?)`
		_, err1 := d.db.Exec(query, status, userID, postID)
		if err1 != nil {
			log.Println(err1)
			return
		}
	case err != nil:
		log.Println(err)
			return
	}
	query := `
	DELETE FROM comment_reactions
	WHERE comment_id = ? AND user_id = ?
	`
	_, err = d.db.Exec(query, postID, userID)
	if err != nil {
		log.Println(err)
		return
	}
	query = `
		INSERT INTO comment_reactions (status, user_id, comment_id) VALUES
		(?, ?, ?)`
		_, err1 := d.db.Exec(query, status, userID, postID)
		if err1 != nil {
			log.Println(err1)
			return
		}
	return
}

func (d *db) GetPostLikeCount(id int64) (count int64, err error) {
	query := `
	SELECT COUNT(status) FROM post_reactions
	WHERE post_id = ? AND status = ?
	`
	var c string
	err = d.db.QueryRow(query, id, 1).Scan(&c)
	switch {
	case err == sql.ErrNoRows:
		log.Println(err)
			return 0, err
	case err != nil:
		log.Println(err)
			return
	default:
		count, _ = strconv.ParseInt(c, 10, strconv.IntSize)
		return
	}
}

func (d *db) GetPostDislikeCount(id int64) (count int64, err error) {
	query := `
	SELECT COUNT(status) FROM post_reactions
	WHERE post_id = ? AND status = ?
	`
	var c string
	err = d.db.QueryRow(query, id, -1).Scan(&c)
	switch {
	case err == sql.ErrNoRows:
		log.Println(err)
			return 0, err
	case err != nil:
		log.Println(err)
			return
	default:

		count, _ = strconv.ParseInt(c, 10, strconv.IntSize)
		return
	}
}

func (d *db)  GetCommentLikeCount(id int64) (count int64, err error) {
	query := `
	SELECT COUNT(status) FROM comment_reactions
	WHERE comment_id = ? AND status = ?
	`
	var c string
	err = d.db.QueryRow(query, id, 1).Scan(&c)
	switch {
	case err == sql.ErrNoRows:
		log.Println(err)
			return 0, err
	case err != nil:
		log.Println(err)
			return
	default:

		count, _ = strconv.ParseInt(c, 10, strconv.IntSize)
		return
	}
}
func (d *db) GetCommentDislikeCount(id int64) (count int64, err error) {
	query := `
	SELECT COUNT(status) FROM comment_reactions
	WHERE comment_id = ? AND status = ?
	`
	var c string
	err = d.db.QueryRow(query, id, -1).Scan(&c)
	switch {
	case err == sql.ErrNoRows:
		log.Println(err)
			return 0, err
	case err != nil:
		log.Println(err)
			return
	default:

		count, _ = strconv.ParseInt(c, 10, strconv.IntSize)
		return
	}
}

func (d *db) GetLikedPosts(userID int64) (p []int64, err error) {
	query := `
	SELECT (post_id) FROM post_reactions
	WHERE user_id = ? AND status = ?
	`
	rows, err := d.db.Query(query, userID, 1)
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
		p = append(p, s.id)
	}
	return
}