package controller

import (
	"errors"
)

var (
	// Unauthenticated error
	Unauthenticated = errors.New("unauthenticated")

	// InvalidPassword error
	InvalidPassword = errors.New("invalid password")

	// ErrPasswordAreNotEqual error
	ErrPasswordAreNotEqual = errors.New("new password are not equal")

	// ErrValidationToken error
	ErrValidationToken = errors.New("could not validate access token from headers")

	// ErrInsertingInDB error
	ErrInsertingInDB = errors.New("could not insert into db")

	// ErrBindingRequest error
	ErrBindingRequest = errors.New("could not bind request")

	// ErrNoContent error
	ErrNoContent = errors.New("no content found in db")
)
