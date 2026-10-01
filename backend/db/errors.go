package db

import "errors"

var ErrNotFound = errors.New("record not found")

type NotFoundError struct {
	Entity string
}

func (e *NotFoundError) Error() string {
	return e.Entity + " not found"
}

func (e *NotFoundError) Unwrap() error {
	return ErrNotFound
}

// ConflictError is a request that was well formed but breaks a rule given the
// current state, such as a full shelter or a skipped step. Message is shown to
// the operator as is.
type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}
