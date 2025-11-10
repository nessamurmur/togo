package model

import (
	"strings"
	"time"
)

// Task represents a unit of work in the bullet journal system.
// It is the core aggregate root with identity (ID) and lifecycle management.
//
// Invariants enforced by this entity:
//   - ID must be a valid UUID (non-nil)
//   - CreatedAt must be set and non-zero
//   - Status must be a valid TaskStatus value
//   - Title must not be empty (after trimming whitespace)
//   - DeferredCount must be >= 0
//
// The Task entity follows these design principles:
//   - Immutable Identity: ID and CreatedAt never change after construction
//   - Factory Pattern: Use NewTask() to create; do not construct directly
//   - Encapsulation: State transitions happen through methods (added in Task 8)
//   - Value Object Composition: Uses TaskID and TaskStatus value objects
type Task struct {
	ID            TaskID     `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	Title         string     `json:"title"`
	Notes         string     `json:"notes,omitempty"`
	Status        TaskStatus `json:"status"`
	Tags          []string   `json:"tags,omitempty"`
	DueDate       *time.Time `json:"due_date,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	DeferredCount int        `json:"deferred_count"`
}

// NewTask creates a new Task with the given title and tags.
// The task is initialized with a newly generated UUID as ID, current timestamp
// as CreatedAt, StatusPool as initial status, and zero values for optional fields.
//
// The title is trimmed of leading/trailing whitespace before validation.
// If the trimmed title is empty, returns ErrEmptyTitle.
//
// The tags slice is defensively copied to prevent external mutation.
// If tags is nil or empty, the Task.Tags field will be nil (for JSON omitempty).
//
// Example:
//
//	task, err := NewTask("Write documentation", []string{"docs", "high-priority"})
//	if err != nil {
//	    return err
//	}
//
// Returns:
//   - A pointer to the newly created Task
//   - ErrEmptyTitle if the title is empty or whitespace-only
func NewTask(title string, tags []string) (*Task, error) {
	trimmedTitle := strings.TrimSpace(title)
	if trimmedTitle == "" {
		return nil, ErrEmptyTitle
	}

	id := NewTaskID()

	var taskTags []string
	if len(tags) > 0 {
		taskTags = make([]string, len(tags))
		copy(taskTags, tags)
	}

	task := &Task{
		ID:            id,
		CreatedAt:     time.Now(),
		Title:         trimmedTitle,
		Notes:         "",
		Status:        StatusPool,
		Tags:          taskTags,
		DueDate:       nil,
		CompletedAt:   nil,
		DeferredCount: 0,
	}

	return task, nil
}
