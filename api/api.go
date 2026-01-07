package api

import (
	"messanger/handler"
	"net/http"
)

// SetupRoutes API 라우트 설정
func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/signup", handler.SignupHandler)
}
