package handler

import (
	"encoding/json"
	"log"
	"messanger/service"
	"net/http"
)

// SendNotificationToUser WebSocket으로 알림 전송 (main 패키지 함수 호출을 위한 래퍼)
var SendNotificationToUser func(int, []byte)

// FriendRequestHandler 친구 요청 핸들러
func FriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var req service.FriendRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	resp, err := service.SendFriendRequest(&req)
	if err != nil {
		log.Printf("### Friend request error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	// WebSocket으로 알림 전송
	if resp.Success && SendNotificationToUser != nil {
		log.Printf("### Sending notification to user %d", req.ToUserID)
		notification := map[string]interface{}{
			"type":    "notification",
			"message": "새로운 친구 요청이 있습니다",
		}
		msg, err := json.Marshal(notification)
		if err != nil {
			log.Printf("### Failed to marshal notification: %v", err)
		} else {
			SendNotificationToUser(req.ToUserID, msg)
		}
	} else {
		if !resp.Success {
			log.Printf("### Friend request failed, not sending notification")
		}
		if SendNotificationToUser == nil {
			log.Printf("### SendNotificationToUser is nil")
		}
	}

	json.NewEncoder(w).Encode(resp)
}

