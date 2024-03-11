package session

type SessionService interface {
	CreateSession(session Session) error
	DeleteSession(token string) error
	GetSessionByToken(token string) (Session, error)
	GetSessionByUserID(userId int) (Session, error)
	GetAllSessionsTime() ([]Session, error)
}

type sessionService struct {
	storage SessionStorage
}

func NewSessionService(storage SessionStorage) SessionService {
	return &sessionService{storage: storage}
}

func (s *sessionService) CreateSession(session Session) error {
	err := s.storage.CreateSession(session)
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionService) GetSessionByToken(token string) (Session, error) {
	session, err := s.storage.GetSessionByToken(token)
	if err != nil {
		return Session{}, err
	}
	return session, nil
}

func (s *sessionService) GetSessionByUserID(userId int) (Session, error) {
	session, err := s.storage.GetSessionByUserID(userId)
	if err != nil {
		return Session{}, err
	}
	return session, nil
}

func (s *sessionService) DeleteSession(token string) error {
	
	return s.storage.DeleteSession(token)
}

func (s *sessionService) GetAllSessionsTime() ([]Session, error) {
	return s.storage.GetAllSessionsTime()
}
