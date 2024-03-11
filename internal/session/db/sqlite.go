package db

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/session"
)

type db struct {
	db *sql.DB
}

func NewSessionStorage(dataBase *sql.DB) session.SessionStorage {
	return &db{
		db: dataBase,
	}
}

func (d *db) CreateSession(session session.Session) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO sessions(user_id, token, expiry) VALUES(?,?,?)", session.UserId, session.Token, session.Expiry)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (d *db) GetSessionByToken(token string) (s session.Session, err error) {
	row := d.db.QueryRow("SELECT user_id, token, expiry FROM sessions WHERE token = ?", token)
	if err = row.Scan(&s.UserId, &s.Token, &s.Expiry); err != nil {
		return
	}
	return
}

func (d *db) GetSessionByUserID(userId int) (s session.Session, err error) {
	row := d.db.QueryRow("SELECT user_id, token, expiry FROM sessions WHERE user_id = ?", userId)
	if err = row.Scan(&s.UserId, &s.Token, &s.Expiry); err != nil {
		fmt.Println(err)
		return
	}
	return
}

func (d *db) DeleteSession(token string) (err error) {
	res, err := d.db.Exec("DELETE FROM sessions WHERE token = ?", token)
	if err != nil {
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return
	}
	if rowsAffected == 0 {
		return errors.New("delete session was failed")
	}
	return
}

func (d *db) GetAllSessionsTime() (s []session.Session, err error) {
	rows, err := d.db.Query("SELECT expiry, token FROM sessions")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var session session.Session
		if err = rows.Scan(&session.Expiry, &session.Token); err != nil {
			return
		}
		s = append(s, session)
	}
	return
}
