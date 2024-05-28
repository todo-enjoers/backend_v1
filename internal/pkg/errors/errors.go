package errors

import (
	"errors"
)

var (
	// Next errors are errors with auth:

	// Unauthenticated error
	Unauthenticated = errors.New("unauthenticated")

	// InvalidPassword error
	InvalidPassword = errors.New("invalid password")

	// ErrPasswordAreNotEqual error
	ErrPasswordAreNotEqual = errors.New("new password are not equal")

	// ErrValidationToken error
	ErrValidationToken = errors.New("could not validate access token from headers")

	// ErrCreateToken error
	ErrCreateToken = errors.New("got error while creating token")

	// ErrGetByLogin error
	ErrGetByLogin = errors.New("the user was not found")

	// ErrBadRegisterRequest error
	ErrBadRegisterRequest = errors.New("wrong login or password. please try other")

	// ErrCreateUser error
	ErrCreateUser = errors.New("got error while creating user")

	// ErrHashingPassword error
	ErrHashingPassword = errors.New("got error while hashing password")

	// ErrComparingPasswords error
	ErrComparingPasswords = errors.New("got error while comparing passwords")

	// Next errors are errors with database:

	// ErrInserting error
	ErrInserting = errors.New("got error while inserting in database")

	// ErrNotAccessible error
	ErrNotAccessible = errors.New("not accessible")

	// ErrNoContent error
	ErrNoContent = errors.New("no content found in db")

	// ErrAlreadyExists error
	ErrAlreadyExists = errors.New("already exist, try other")

	// ErrGetByID error
	ErrGetByID = errors.New("got error while getting by id")

	//ErrTableMigrations error
	ErrTableMigrations = errors.New("migrations failed")

	// Don't sorting errors:

	// ErrBindingRequest error
	ErrBindingRequest = errors.New("could not bind request")
	
	// ErrNotFound error
	ErrNotFound = errors.New("not found")

	//ErrBadRequestId error
	ErrBadRequestId = errors.New("got error while validating id")

	//ErrInternalServer error
	ErrInternalServer = errors.New("failed to retrieve todo")
)
