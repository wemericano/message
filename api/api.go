package api

import (
	"messanger/handler"
	"net/http"
)

// SetupRoutes API 라우트 설정
func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/login", handler.LoginHandler)
	mux.HandleFunc("/api/signup", handler.SignupHandler)
	mux.HandleFunc("/api/search", handler.SearchUsersHandler)
	mux.HandleFunc("/api/friend/request", handler.FriendRequestHandler)
}
