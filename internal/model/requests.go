package model

import (
	"errors"
	"github.com/google/uuid"
	"strings"
)

type (
	// UserRegisterRequest : :Registration Request from user
	UserRegisterRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	// UserLoginRequest : Authorization Request from user
	UserLoginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	// UserCoupleTokensRequest : Request of generation a couple of tokens
	UserCoupleTokensRequest struct {
		ID           uuid.UUID `json:"id"`
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
	}
	// UserChangePasswordRequest : Changing password Request from user
	UserChangePasswordRequest struct {
		OldPassword      string `json:"old_password"`
		NewPassword      string `json:"new_password"`
		NewPasswordAgain string `json:"new_password_again"`
	}
	// TodoCreateRequest :Creating TodoType Request from user
	TodoCreateRequest struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		IsCompleted bool      `json:"is_completed"`
		CreatedBy   uuid.UUID `json:"created_by"`
		ProjectID   uuid.UUID `json:"project_id"`
		Column      string    `json:"column"`
	}
	// TodoUpdateRequest :Updating TodoType Request from user
	TodoUpdateRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsCompleted bool   `json:"is_completed"`
	}
	// ColumRequest :Updating ColumnType Request from user
	ColumRequest struct {
		ProjectId uuid.UUID `json:"project_id"`
		Name      string    `json:"name"`
		Order     int       `json:"order"`
	}
	// ProjectRequest :Updating ProjectType Request from user
	ProjectRequest struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		CreatedBy uuid.UUID `json:"created_by"`
	}
)

func (req *UserRegisterRequest) Validate() (ok bool, err error) {
	if !strings.Contains(req.Login, "@") && len(req.Login) > 7 {
		err = errors.New("wrong email address")
		return false, err
	}

	if len(req.Password) < 7 {
		err = errors.New("password is required")
		return false, err
	}

	return true, nil
}
