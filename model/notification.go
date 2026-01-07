package model

// Notification 알림 모델
type Notification struct {
	ID        int    `json:"id" db:"id"`
	UserID    int    `json:"user_id" db:"user_id"`
	Type      string `json:"type" db:"type"` // friend_request
	Message   string `json:"message" db:"message"`
	IsRead    bool   `json:"is_read" db:"is_read"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

