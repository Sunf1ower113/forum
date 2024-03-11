package user

import (
	"log"
)

type userService struct {
	storage UserStorage
	logger  *log.Logger
}

type UserService interface {
	GetUserByID(id int64) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUserByToken(token string) (User, error)
}

func NewUserService(storage UserStorage, logger *log.Logger) UserService {
	return &userService{
		storage: storage,
		logger:  logger,
	}
}

func (s *userService) GetUserByID(id int64) (user User, err error) {
	user, err = s.storage.GetUserByID(id)
	return
}

func (s *userService) GetUserByEmail(email string) (user User, err error) {
	user, err = s.storage.GetUserByEmail(email)
	return
}

func (s *userService) GetUserByToken(token string) (u User, err error) {
	userID, err := s.storage.GetUserIdByToken(token)
	if err != nil {
		return
	}
	u, err = s.storage.GetUserByID(userID)
	if err != nil {
		return
	}
	return
}
