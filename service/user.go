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

// LoginRequest 로그인 요청
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 로그인 응답
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

// Login 로그인 서비스
func Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := dbcall.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &LoginResponse{
			Success: false,
			Message: "아이디 또는 비밀번호가 잘못되었습니다",
		}, nil
	}

	if user.Password != req.Password {
		return &LoginResponse{
			Success: false,
			Message: "아이디 또는 비밀번호가 잘못되었습니다",
		}, nil
	}

	return &LoginResponse{
		Success: true,
		Message: "로그인 성공",
		UserID:  user.ID,
	}, nil
}

// SearchRequest 사용자 검색 요청
type SearchRequest struct {
	Name string `json:"name"`
}

// SearchResponse 사용자 검색 응답
type SearchResponse struct {
	Success bool        `json:"success"`
	Users   []*UserInfo `json:"users"`
	Message string      `json:"message"`
}

// UserInfo 사용자 정보 (비밀번호 제외)
type UserInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// SearchUsers 사용자 검색 서비스
func SearchUsers(req *SearchRequest) (*SearchResponse, error) {
	if req.Name == "" {
		return &SearchResponse{
			Success: false,
			Users:   []*UserInfo{},
			Message: "검색어를 입력해주세요",
		}, nil
	}

	users, err := dbcall.SearchUsersByName(req.Name)
	if err != nil {
		return nil, err
	}

	userInfos := make([]*UserInfo, 0, len(users))
	for _, user := range users {
		userInfos = append(userInfos, &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Email:    user.Email,
		})
	}

	return &SearchResponse{
		Success: true,
		Users:   userInfos,
		Message: "검색 완료",
	}, nil
}
