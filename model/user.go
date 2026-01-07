package model

// User 사용자 모델
type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
}
