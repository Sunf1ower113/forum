package session

import "time"

type Session struct {
	ID       int64
	UserId   int64
	Email    string
	Username string
	Token    string
	Expiry   time.Time
}
