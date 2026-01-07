package service

import (
	"messanger/dbcall"
)

// FriendRequestRequest 친구 요청 요청
type FriendRequestRequest struct {
	FromUserID int `json:"from_user_id"`
	ToUserID   int `json:"to_user_id"`
}

// FriendRequestResponse 친구 요청 응답
type FriendRequestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SendFriendRequest 친구 요청 전송
func SendFriendRequest(req *FriendRequestRequest) (*FriendRequestResponse, error) {
	// 사용자 존재 확인
	toUser, err := dbcall.GetUserByID(req.ToUserID)
	if err != nil {
		return nil, err
	}
	if toUser == nil {
		return &FriendRequestResponse{
			Success: false,
			Message: "사용자를 찾을 수 없습니다",
		}, nil
	}

	// 친구 요청 생성
	if err := dbcall.CreateFriendRequest(req.FromUserID, req.ToUserID); err != nil {
		return nil, err
	}

	// 알림 생성
	message := "새로운 친구 요청이 있습니다"
	if err := dbcall.CreateNotification(req.ToUserID, "friend_request", message); err != nil {
		// 알림 생성 실패해도 친구 요청은 성공으로 처리
	}

	return &FriendRequestResponse{
		Success: true,
		Message: "친구 요청을 보냈습니다",
	}, nil
}

// NotifyUser 사용자에게 알림 전송 (WebSocket)
func NotifyUser(userID int, message string) {
	// 이 함수는 handler에서 websocket.SendToUser를 호출하도록 함
	// service 패키지에서 main 패키지 함수를 직접 호출할 수 없으므로
	// handler에서 처리
	notification := map[string]interface{}{
		"type":    "notification",
		"message": message,
	}
	_ = userID
	_ = notification
	// 실제 전송은 handler에서 처리
}

