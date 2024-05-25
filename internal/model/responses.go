package model

import "github.com/google/uuid"

type (
	// UserRegisterResponse :Registration Response from server
	UserRegisterResponse struct {
		ID           uuid.UUID `json:"id"`
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
	}
	// UserLoginResponse : Authorization Response from server
	UserLoginResponse struct {
		ID           uuid.UUID `json:"id"`
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
	}
	// GroupResponse : Group Response from server
	GroupResponse struct {
		UserID    uuid.UUID `json:"user_id"`
		ProjectID uuid.UUID `json:"project_id"`
	}
	// UserCoupleTokensResponse : Response of generation a couple of tokens
	UserCoupleTokensResponse struct {
		ID           uuid.UUID `json:"id"`
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
	}
	// UserGetMeResponse : Creation ??? Response from server
	UserGetMeResponse struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}
	// ErrorResponse : Creation Error Response from server
	ErrorResponse struct {
		Error string `json:"error"`
	}
	// TodoCreateResponse : Todos Response
	TodoCreateResponse struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		IsCompleted bool      `json:"is_complete"`
		ProjectID   uuid.UUID `json:"project_id"`
		CreatedBy   uuid.UUID `json:"created_by"`
		Column      string    `json:"column"`
	}
	// ColumResponse : Column Response from server
	ColumResponse struct {
		ProjectId uuid.UUID `json:"project_id"`
		Name      string    `json:"name"`
		Order     int       `json:"order"`
	}
	// ProjectResponse : Project Response from server
	ProjectResponse struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		CreatedBy uuid.UUID `json:"created_by"`
	}
)
