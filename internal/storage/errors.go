package storage

import "errors"

var (
	// ErrAlreadyExists error
	ErrAlreadyExists = errors.New("already in storage")

	// ErrNotAccessible error
	ErrNotAccessible = errors.New("not accessible")

	// ErrBadRegisterRequest error
	ErrBadRegisterRequest = errors.New("wrong login or password. please try other")

	// ErrCreateUser error
	ErrCreateUser = errors.New("got error while creating user")

	// ErrCreateToken error
	ErrCreateToken = errors.New("got error while creating token")

	// ErrGetByLogin error
	ErrGetByLogin = errors.New("the user was not found")

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

	//ErrInternalServer error
	ErrInternalServer = errors.New("failed to retrieve todo")

	//ErrTableMigrations error
	ErrTableMigrations = errors.New("migrations failed")
)
