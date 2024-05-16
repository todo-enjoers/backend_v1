package model

import (
	"errors"
	"strings"
)

type (
	// UserRegisterRequest : :Registration Request from user
	UserRegisterRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	// UserLoginRequest :Authorization Request from user
	UserLoginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	// UserChangePasswordRequest : Changing password Request from user
	UserChangePasswordRequest struct {
		OldPassword      string `json:"old_password"`
		NewPassword      string `json:"new_password"`
		NewPasswordAgain string `json:"new_password_again"`
	}
	// TodoCreateRequest :Creation TodoType Request from user
	TodoCreateRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
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
