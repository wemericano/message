package model

// FriendRequest 친구 요청 모델
type FriendRequest struct {
	ID         int    `json:"id" db:"id"`
	FromUserID int    `json:"from_user_id" db:"from_user_id"`
	ToUserID   int    `json:"to_user_id" db:"to_user_id"`
	Status     string `json:"status" db:"status"` // pending, accepted, rejected
	CreatedAt  string `json:"created_at" db:"created_at"`
}

