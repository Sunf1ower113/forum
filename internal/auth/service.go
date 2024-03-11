package auth

import (
	"errors"
	"fmt"
	"forum/internal/session"
	"forum/internal/user"
	"log"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type AuthService interface {
	Login(user user.User) (session.Session, error)
	SignUp(user user.User) (user.User, error)
	Logout(token string) error
}

type authService struct {
	sessionService session.SessionService
	userService    user.UserService
	userStorage    user.UserStorage
}

func NewAuthService(sessionService session.SessionService, userService user.UserService, userStorage user.UserStorage) AuthService {
	return &authService{
		sessionService: sessionService,
		userService:    userService,
		userStorage:    userStorage,
	}
}

func (a *authService) Login(user user.User) (s session.Session, err error) {
	userDB, err := a.userStorage.GetUserByEmail(user.Email)
	if err != nil {
		log.Printf("user %s sign in was failed\n", user.Email)
		return s, errors.New("wrong password or email")
	}
	if user.Password != userDB.Password {
		return session.Session{}, errors.New("Wrong password")
	}
	sessionDb, err := a.sessionService.GetSessionByUserID(int(userDB.ID))
	if err != nil {
		log.Printf("session for user_id %d is not found\n", userDB.ID)
	} else {
		err := a.sessionService.DeleteSession(sessionDb.Token)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("session for user_id %d is deleted\n", user.ID)
		}
	}
	sessionToken := uuid.NewString()
	expiry := time.Now().Add(72 * time.Hour)
	session := session.Session{
		UserId: userDB.ID,
		Token:  sessionToken,
		Expiry: expiry,
	}
	err = a.sessionService.CreateSession(session)
	if err != nil {
		return s, fmt.Errorf("session for user %d was failed\nerror: %w", user.ID, err)
	}
	log.Printf("user %s sign in was successfully\n", user.Email)
	return session, nil
}

func (a *authService) SignUp(u user.User) (user.User, error) {
	nameRegex, err := regexp.Compile("[a-zA-Z0-9_-]{3,16}")
	if err != nil {
		log.Println("Name fail")
		return user.User{}, errors.New("name regex fail")
	}
	emailRegex, err := regexp.Compile(`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,6}`)
	if err != nil {
		log.Println("email fail")

		return user.User{}, errors.New("name regex fail")
	}
	passwordRegex, err := regexp.Compile("[a-zA-Z0-9!@#$%^&*()_+=-]{8,}")
	if err != nil {
		log.Println("pass fail")

		return user.User{}, errors.New("pass regex fail")
	}
	usernameIsValid := nameRegex.MatchString(u.Username)
	emailIsValid := emailRegex.MatchString(u.Email)
	passwordIsValid := passwordRegex.MatchString(u.Password)
	if !(passwordIsValid && usernameIsValid && emailIsValid) {
		log.Println("invalid user data for sign up")
		return u, errors.New("invalid user data for sign up")
	}
	if _, err = a.userService.GetUserByEmail(u.Email); err == nil {
		log.Println("An account with this email already exists")
		return u, errors.New("An account with this email already exists")
	}
	_, err = a.userStorage.Create(u)
	if err != nil {
		log.Println(err)
		return u, err
	}
	return u, nil
}

func (a *authService) Logout(token string) error {
	return a.sessionService.DeleteSession(token)
}
