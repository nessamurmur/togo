package model

import (
	"errors"
	"fmt"
)

var (
	// ErrTaskNotFound indicates a task ID does not exist.
	ErrTaskNotFound = errors.New("task not found")

	// ErrInvalidStatus indicates an invalid task status value.
	ErrInvalidStatus = errors.New("invalid task status")

	// ErrInvalidStateTransition indicates a state transition is not allowed.
	ErrInvalidStateTransition = errors.New("invalid state transition")

	// ErrEmptyTitle indicates a task title is empty or whitespace-only.
	ErrEmptyTitle = errors.New("task title cannot be empty")

	// ErrDuplicateTaskID indicates a task with the same ID already exists.
	ErrDuplicateTaskID = errors.New("task with this ID already exists")
)

// ValidationError wraps validation failures with field and reason information.
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Reason)
}
