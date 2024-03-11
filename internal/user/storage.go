package user

type UserStorage interface {
	Create(user User) (User, error)
	GetUserByID(id int64) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUserIdByToken(token string) (int64, error)
	// Update(user User) error
	// Delete(id string) error
}
