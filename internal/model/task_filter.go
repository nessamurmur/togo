package model

import "time"

// TaskFilter encapsulates criteria for filtering tasks in queries.
// It supports filtering by status, tags (AND semantics), and due date ranges.
// A nil or zero value for a field means no filtering on that criterion.
type TaskFilter struct {
	Status    *TaskStatus
	Tags      []string
	DueAfter  *time.Time
	DueBefore *time.Time
	Limit     int
}

// Matches returns true if the task satisfies all filter criteria.
// Uses fail-fast pattern: returns false as soon as any criterion fails.
//
// Filtering semantics:
//   - Status: nil matches any status; non-nil requires exact match
//   - Tags: nil or empty matches any tags; non-empty requires task to have ALL filter tags (AND semantics)
//   - DueAfter: nil matches any date; non-nil requires task.DueDate >= DueAfter (inclusive, rejects nil DueDate)
//   - DueBefore: nil matches any date; non-nil requires task.DueDate <= DueBefore (inclusive, rejects nil DueDate)
//   - Limit: completely ignored by Matches (caller's responsibility to apply limit)
func (f TaskFilter) Matches(t *Task) bool {
	if f.Status != nil && t.Status != *f.Status {
		return false
	}

	if len(f.Tags) > 0 {
		if !containsAllTags(t.Tags, f.Tags) {
			return false
		}
	}

	if f.DueAfter != nil {
		if t.DueDate == nil {
			return false
		}
		if t.DueDate.Before(*f.DueAfter) {
			return false
		}
	}

	if f.DueBefore != nil {
		if t.DueDate == nil {
			return false
		}
		if t.DueDate.After(*f.DueBefore) {
			return false
		}
	}

	return true
}

// containsAllTags returns true if taskTags contains all tags in filterTags.
// Uses map-based lookup for O(n) performance.
// Empty filterTags always returns true.
func containsAllTags(taskTags, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}

	tagSet := make(map[string]bool, len(taskTags))
	for _, tag := range taskTags {
		tagSet[tag] = true
	}

	for _, requiredTag := range filterTags {
		if !tagSet[requiredTag] {
			return false
		}
	}

	return true
}
