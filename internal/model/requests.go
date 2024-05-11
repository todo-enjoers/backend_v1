package model

type (
	// UserRegisterRequest : : Registration Request from user
	UserRegisterRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	// UserLoginRequest : Authorization Request from user
	UserLoginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	// TodoCreateRequest : Creation TodoType Request from user
	TodoCreateRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)
