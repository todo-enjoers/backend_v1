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

	// ErrCreateGroup error
	ErrCreateGroup = errors.New("got error while creating group")

	// ErrCreateToken error
	ErrCreateToken = errors.New("got error while creating token")

	// ErrGetByLogin error
	ErrGetByLogin = errors.New("got error while getting by login")

	// ErrGetByID error
	ErrGetByID = errors.New("got error while getting by id")

	// ErrHashingPassword error
	ErrHashingPassword = errors.New("got error while hashing password")

	// ErrComparingPasswords error
	ErrComparingPasswords = errors.New("got error while comparing passwords")

	// ErrNotFound error
	ErrNotFound = errors.New("not found")

	// ErrInserting error
	ErrInserting = errors.New("got error while inserting in database")

	//ErrBadRequestId error
	ErrBadRequestId = errors.New("got error while validating id")

	//ErrForbidden error
	ErrForbidden = errors.New("forbidden")

	//ErrInternalServer error
	ErrInternalServer = errors.New("failed to retrieve todo")
)
