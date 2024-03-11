package db

import (
	"database/sql"
	"forum/internal/user"
	"log"
)

type db struct {
	db     *sql.DB
	logger *log.Logger
}

func NewUserStorage(dataBase *sql.DB, logger *log.Logger) user.UserStorage {
	return &db{
		db:     dataBase,
		logger: logger,
	}
}

func (d *db) Create(u user.User) (user.User, error) {
	_, err := d.db.Exec("INSERT INTO users(email, username, password) values(?,?,?)", u.Email, u.Username, u.Password)
	if err != nil {
		d.logger.Printf("Error:%s\n", err)
		return u, err
	}
	d.logger.Printf("User %s has been created\n", u.Username)
	return u, nil
}

func (d *db) GetUserByEmail(email string) (u user.User, err error) {
	row := d.db.QueryRow("SELECT user_id,email,password,username FROM users WHERE email = ?", email)
	if err = row.Scan(&u.ID, &u.Email, &u.Password, &u.Username); err != nil {
		d.logger.Printf("Error query:%s\n", err)
		return
	}
	d.logger.Printf("Succsed query%s\n", u.Username)
	return
}

func (d *db) GetUserByID(id int64) (u user.User, err error) {
	row := d.db.QueryRow("SELECT user_id, email, password, username FROM users WHERE user_id = ?", id)
	if err = row.Scan(&u.ID, &u.Email, &u.Password, &u.Username); err != nil {
		d.logger.Printf("Error query: %s\n", err)
		return
	}
	d.logger.Printf("Succsed query %s\n", u.Username)
	return
}

func (d *db) GetUserIdByToken(token string) (userID int64, err error) {
	row := d.db.QueryRow("SELECT user_id FROM sessions WHERE token=?", token)
	if err = row.Scan(&userID); err != nil {
		return -1, err
	}
	return
}

// func (d *db) Delete(id string) error {
// 	panic("implement me")
// }

// func (d *db) Update(user user.User) error {
// 	panic("implement me")
// }
