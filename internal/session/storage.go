package session

type SessionStorage interface {
	GetSessionByToken(token string) (Session, error)
	GetSessionByUserID(userId int) (Session, error)
	GetAllSessionsTime() ([]Session, error)
	CreateSession(session Session) error
	DeleteSession(token string) error
}
