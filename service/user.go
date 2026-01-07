package service

import (
	"messanger/dbcall"
	"messanger/model"
)

// SignupRequest 회원가입 요청
type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// SignupResponse 회원가입 응답
type SignupResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Signup 회원가입 서비스
func Signup(req *SignupRequest) (*SignupResponse, error) {
	// 중복 체크
	existingUser, err := dbcall.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return &SignupResponse{
			Success: false,
			Message: "이미 존재하는 아이디입니다",
		}, nil
	}

	// 사용자 생성
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
		Email:    req.Email,
	}

	if err := dbcall.CreateUser(user); err != nil {
		return nil, err
	}

	return &SignupResponse{
		Success: true,
		Message: "회원가입이 완료되었습니다",
	}, nil
}
