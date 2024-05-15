package storage

import "errors"

var (
	// ErrAlreadyExists error
	ErrAlreadyExists = errors.New("already in storage")

	// ErrNoContent error
	ErrNoContent = errors.New("no content")

	// ErrIsDeleted error
	ErrIsDeleted = errors.New("is deleted")

	// ErrNotAccessible error
	ErrNotAccessible = errors.New("not accessible")

	// ErrAlreadyClosed error
	ErrAlreadyClosed = errors.New("storage is already closed")

	// ErrBadRegisterRequest error
	ErrBadRegisterRequest = errors.New("wrong login or password. please try other")

	// ErrCreateUser error
	ErrCreateUser = errors.New("got error while creating user")

	// ErrCreateToken error
	ErrCreateToken = errors.New("got error while creating token")

	// ErrGetByLogin error
	ErrGetByLogin = errors.New("got error while getting by login")

	// ErrGetByID error
	ErrGetByID = errors.New("got error while getting by id")

	// ErrHashingPass error
	ErrHashingPassword = errors.New("got error while hashing password")

	// ErrComparingPasswords error
	ErrComparingPasswords = errors.New("got error while comparing passwords")
)
