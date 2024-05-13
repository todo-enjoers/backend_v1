package storage

import "errors"

// vars ...
var (
	// ErrNotFound error
	ErrNotFound = errors.New("not found")
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
)
