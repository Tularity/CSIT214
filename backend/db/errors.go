package db

import (
	"errors"
	"fmt"

	"github.com/Tularity/CSIT214/backend/models"
)

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

// ShelterFullError is the conflict raised when a shelter cannot take everyone
// from an incident. It carries the shelters that could, so the coordinator can
// pick one straight away instead of trying them one by one (change CR-01).
type ShelterFullError struct {
	Remaining    int
	Alternatives []models.ShelterAlternative
}

func (e *ShelterFullError) Error() string {
	return fmt.Sprintf("Shelter has only %d places remaining", e.Remaining)
}
