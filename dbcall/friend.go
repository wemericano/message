package dbcall

import (
	"database/sql"
	"messanger/model"
)

// CreateFriendRequest 친구 요청 생성
func CreateFriendRequest(fromUserID, toUserID int) error {
	query := `INSERT INTO friend_requests (from_user_id, to_user_id, status) VALUES (@p1, @p2, 'pending')`
	_, err := db.Exec(query, fromUserID, toUserID)
	return err
}

// GetUserByID ID로 사용자 조회
func GetUserByID(userID int) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, password, name, email FROM users WHERE id = @p1`
	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// CreateNotification 알림 생성
func CreateNotification(userID int, notificationType, message string) error {
	query := `INSERT INTO notifications (user_id, type, message, is_read) VALUES (@p1, @p2, @p3, 0)`
	_, err := db.Exec(query, userID, notificationType, message)
	return err
}

// GetUnreadNotificationCount 읽지 않은 알림 개수 조회
func GetUnreadNotificationCount(userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = @p1 AND is_read = 0`
	err := db.QueryRow(query, userID).Scan(&count)
	return count, err
}

