package user

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// type CreateUser struct {
// 	Email    string `json:"email"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }
