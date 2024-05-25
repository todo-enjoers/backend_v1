package controller

import (
	"errors"
)

var (
	// Unauthenticated error
	Unauthenticated = errors.New("unauthenticated")

	// Internal error
	//Internal = errors.New("internal")

	// NotFound error
	//NotFound = errors.New("not found")

	// BadRequest error
	//BadRequest = errors.New("bad request")

	// Canceled error
	//Canceled = errors.New("canceled")

	// PermissionDenied error
	//PermissionDenied = errors.New("permission denied")

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
)
