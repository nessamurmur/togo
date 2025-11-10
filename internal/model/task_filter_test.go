package model

import (
	"testing"
	"time"
)

// Temporary Task definition for testing (remove when Task 7 is complete)
// This minimal struct allows parallel development of TaskFilter before Task entity is implemented
type Task struct {
	Status  TaskStatus
	Tags    []string
	DueDate *time.Time
}

// Test fixtures - reusable test data
var (
	statusPool  = StatusPool
	statusToday = StatusToday
	statusDone  = StatusDone

	now       = time.Now()
	yesterday = now.Add(-24 * time.Hour)
	tomorrow  = now.Add(24 * time.Hour)
	nextWeek  = now.Add(7 * 24 * time.Hour)
)

func TestTaskFilter_Matches_StatusFilter(t *testing.T) {
	tests := []struct {
		name   string
		filter TaskFilter
		task   *Task
		want   bool
	}{
		{
			name:   "nil status filter matches any status",
			filter: TaskFilter{Status: nil},
			task:   &Task{Status: StatusPool},
			want:   true,
		},
		{
			name:   "status filter matches exact status - pool",
			filter: TaskFilter{Status: &statusPool},
			task:   &Task{Status: StatusPool},
			want:   true,
		},
		{
			name:   "status filter matches exact status - today",
			filter: TaskFilter{Status: &statusToday},
			task:   &Task{Status: StatusToday},
			want:   true,
		},
		{
			name:   "status filter matches exact status - done",
			filter: TaskFilter{Status: &statusDone},
			task:   &Task{Status: StatusDone},
			want:   true,
		},
		{
			name:   "status filter rejects non-matching status",
			filter: TaskFilter{Status: &statusPool},
			task:   &Task{Status: StatusToday},
			want:   false,
		},
		{
			name:   "status filter rejects different status",
			filter: TaskFilter{Status: &statusDone},
			task:   &Task{Status: StatusPool},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.task)
			if got != tt.want {
				t.Errorf("TaskFilter.Matches() = %v, want %v (filter status: %v, task status: %v)",
					got, tt.want, tt.filter.Status, tt.task.Status)
			}
		})
	}
}

func TestTaskFilter_Matches_TagsFilter(t *testing.T) {
	tests := []struct {
		name   string
		filter TaskFilter
		task   *Task
		want   bool
	}{
		{
			name:   "nil tags filter matches any tags",
			filter: TaskFilter{Tags: nil},
			task:   &Task{Tags: []string{"work", "urgent"}},
			want:   true,
		},
		{
			name:   "empty tags filter matches any tags",
			filter: TaskFilter{Tags: []string{}},
			task:   &Task{Tags: []string{"work", "urgent"}},
			want:   true,
		},
		{
			name:   "single tag filter matches task with that tag",
			filter: TaskFilter{Tags: []string{"work"}},
			task:   &Task{Tags: []string{"work"}},
			want:   true,
		},
		{
			name:   "single tag filter matches task with multiple tags including that tag",
			filter: TaskFilter{Tags: []string{"work"}},
			task:   &Task{Tags: []string{"work", "urgent", "important"}},
			want:   true,
		},
		{
			name:   "multiple tag filter matches task with all those tags (AND semantics)",
			filter: TaskFilter{Tags: []string{"work", "urgent"}},
			task:   &Task{Tags: []string{"work", "urgent"}},
			want:   true,
		},
		{
			name:   "multiple tag filter matches task with all filter tags plus more",
			filter: TaskFilter{Tags: []string{"work", "urgent"}},
			task:   &Task{Tags: []string{"work", "urgent", "important", "personal"}},
			want:   true,
		},
		{
			name:   "tag filter rejects task missing one required tag",
			filter: TaskFilter{Tags: []string{"work", "urgent"}},
			task:   &Task{Tags: []string{"work"}},
			want:   false,
		},
		{
			name:   "tag filter rejects task with no matching tags",
			filter: TaskFilter{Tags: []string{"work"}},
			task:   &Task{Tags: []string{"personal", "home"}},
			want:   false,
		},
		{
			name:   "tag filter rejects task with no tags",
			filter: TaskFilter{Tags: []string{"work"}},
			task:   &Task{Tags: []string{}},
			want:   false,
		},
		{
			name:   "tag filter rejects task with nil tags",
			filter: TaskFilter{Tags: []string{"work"}},
			task:   &Task{Tags: nil},
			want:   false,
		},
		{
			name:   "empty filter matches task with no tags",
			filter: TaskFilter{Tags: []string{}},
			task:   &Task{Tags: []string{}},
			want:   true,
		},
		{
			name:   "nil filter matches task with nil tags",
			filter: TaskFilter{Tags: nil},
			task:   &Task{Tags: nil},
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.task)
			if got != tt.want {
				t.Errorf("TaskFilter.Matches() = %v, want %v (filter tags: %v, task tags: %v)",
					got, tt.want, tt.filter.Tags, tt.task.Tags)
			}
		})
	}
}

func TestTaskFilter_Matches_DueDateRange(t *testing.T) {
	tests := []struct {
		name   string
		filter TaskFilter
		task   *Task
		want   bool
	}{
		{
			name:   "nil DueAfter and DueBefore matches any due date",
			filter: TaskFilter{DueAfter: nil, DueBefore: nil},
			task:   &Task{DueDate: &now},
			want:   true,
		},
		{
			name:   "DueAfter matches task due after that date",
			filter: TaskFilter{DueAfter: &yesterday},
			task:   &Task{DueDate: &now},
			want:   true,
		},
		{
			name:   "DueAfter matches task due exactly on that date (inclusive)",
			filter: TaskFilter{DueAfter: &now},
			task:   &Task{DueDate: &now},
			want:   true,
		},
		{
			name:   "DueAfter rejects task due before that date",
			filter: TaskFilter{DueAfter: &now},
			task:   &Task{DueDate: &yesterday},
			want:   false,
		},
		{
			name:   "DueBefore matches task due before that date",
			filter: TaskFilter{DueBefore: &tomorrow},
			task:   &Task{DueDate: &now},
			want:   true,
		},
		{
			name:   "DueBefore matches task due exactly on that date (inclusive)",
			filter: TaskFilter{DueBefore: &now},
			task:   &Task{DueDate: &now},
			want:   true,
		},
		{
			name:   "DueBefore rejects task due after that date",
			filter: TaskFilter{DueBefore: &now},
			task:   &Task{DueDate: &tomorrow},
			want:   false,
		},
		{
			name:   "DueAfter and DueBefore together define inclusive range",
			filter: TaskFilter{DueAfter: &yesterday, DueBefore: &tomorrow},
			task:   &Task{DueDate: &now},
			want:   true,
		},
		{
			name:   "date range matches task at lower boundary (inclusive)",
			filter: TaskFilter{DueAfter: &yesterday, DueBefore: &tomorrow},
			task:   &Task{DueDate: &yesterday},
			want:   true,
		},
		{
			name:   "date range matches task at upper boundary (inclusive)",
			filter: TaskFilter{DueAfter: &yesterday, DueBefore: &tomorrow},
			task:   &Task{DueDate: &tomorrow},
			want:   true,
		},
		{
			name:   "date range rejects task before range",
			filter: TaskFilter{DueAfter: &now, DueBefore: &nextWeek},
			task:   &Task{DueDate: &yesterday},
			want:   false,
		},
		{
			name:   "date range rejects task after range",
			filter: TaskFilter{DueAfter: &yesterday, DueBefore: &now},
			task:   &Task{DueDate: &tomorrow},
			want:   false,
		},
		{
			name:   "DueAfter filter rejects task with nil DueDate",
			filter: TaskFilter{DueAfter: &now},
			task:   &Task{DueDate: nil},
			want:   false,
		},
		{
			name:   "DueBefore filter rejects task with nil DueDate",
			filter: TaskFilter{DueBefore: &now},
			task:   &Task{DueDate: nil},
			want:   false,
		},
		{
			name:   "date range filter rejects task with nil DueDate",
			filter: TaskFilter{DueAfter: &yesterday, DueBefore: &tomorrow},
			task:   &Task{DueDate: nil},
			want:   false,
		},
		{
			name:   "nil date filters match task with nil DueDate",
			filter: TaskFilter{DueAfter: nil, DueBefore: nil},
			task:   &Task{DueDate: nil},
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.task)
			if got != tt.want {
				t.Errorf("TaskFilter.Matches() = %v, want %v (filter: after=%v, before=%v; task due: %v)",
					got, tt.want, tt.filter.DueAfter, tt.filter.DueBefore, tt.task.DueDate)
			}
		})
	}
}

func TestTaskFilter_Matches_MultipleFilters_AllMustMatch(t *testing.T) {
	workTag := "work"
	urgentTag := "urgent"

	tests := []struct {
		name   string
		filter TaskFilter
		task   *Task
		want   bool
	}{
		{
			name: "task matches all filters - status and tags",
			filter: TaskFilter{
				Status: &statusToday,
				Tags:   []string{"work"},
			},
			task: &Task{
				Status: StatusToday,
				Tags:   []string{"work", "urgent"},
			},
			want: true,
		},
		{
			name: "task matches all filters - status, tags, and due date",
			filter: TaskFilter{
				Status:    &statusToday,
				Tags:      []string{"work"},
				DueAfter:  &yesterday,
				DueBefore: &tomorrow,
			},
			task: &Task{
				Status:  StatusToday,
				Tags:    []string{"work", "urgent"},
				DueDate: &now,
			},
			want: true,
		},
		{
			name: "task fails when status doesn't match",
			filter: TaskFilter{
				Status: &statusPool,
				Tags:   []string{"work"},
			},
			task: &Task{
				Status: StatusToday,
				Tags:   []string{"work"},
			},
			want: false,
		},
		{
			name: "task fails when tags don't match",
			filter: TaskFilter{
				Status: &statusToday,
				Tags:   []string{"work", "urgent"},
			},
			task: &Task{
				Status: StatusToday,
				Tags:   []string{"work"},
			},
			want: false,
		},
		{
			name: "task fails when due date doesn't match",
			filter: TaskFilter{
				Status:    &statusToday,
				Tags:      []string{"work"},
				DueBefore: &yesterday,
			},
			task: &Task{
				Status:  StatusToday,
				Tags:    []string{"work"},
				DueDate: &now,
			},
			want: false,
		},
		{
			name: "task fails when multiple filters don't match",
			filter: TaskFilter{
				Status: &statusPool,
				Tags:   []string{"urgent"},
			},
			task: &Task{
				Status: StatusToday,
				Tags:   []string{"work"},
			},
			want: false,
		},
		{
			name: "task matches with some nil filters",
			filter: TaskFilter{
				Status:    &statusToday,
				Tags:      nil, // nil = match any
				DueAfter:  nil, // nil = no filter
				DueBefore: &tomorrow,
			},
			task: &Task{
				Status:  StatusToday,
				Tags:    []string{workTag, urgentTag},
				DueDate: &now,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.task)
			if got != tt.want {
				t.Errorf("TaskFilter.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskFilter_Matches_NoFilters_MatchesAll(t *testing.T) {
	tests := []struct {
		name   string
		filter TaskFilter
		task   *Task
		want   bool
	}{
		{
			name:   "empty filter matches task with all fields",
			filter: TaskFilter{},
			task: &Task{
				Status:  StatusToday,
				Tags:    []string{"work", "urgent"},
				DueDate: &now,
			},
			want: true,
		},
		{
			name:   "empty filter matches task with minimal fields",
			filter: TaskFilter{},
			task: &Task{
				Status: StatusPool,
			},
			want: true,
		},
		{
			name:   "empty filter matches task with nil fields",
			filter: TaskFilter{},
			task: &Task{
				Status:  StatusDone,
				Tags:    nil,
				DueDate: nil,
			},
			want: true,
		},
		{
			name: "filter with all nil/empty values matches any task",
			filter: TaskFilter{
				Status:    nil,
				Tags:      []string{},
				DueAfter:  nil,
				DueBefore: nil,
				Limit:     0,
			},
			task: &Task{
				Status:  StatusToday,
				Tags:    []string{"work"},
				DueDate: &now,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.task)
			if got != tt.want {
				t.Errorf("TaskFilter.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskFilter_Matches_TagEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		filter TaskFilter
		task   *Task
		want   bool
	}{
		{
			name:   "filter matches task with duplicate tags",
			filter: TaskFilter{Tags: []string{"work"}},
			task:   &Task{Tags: []string{"work", "work", "urgent"}},
			want:   true,
		},
		{
			name:   "filter with duplicate tags still uses AND semantics",
			filter: TaskFilter{Tags: []string{"work", "work"}},
			task:   &Task{Tags: []string{"work", "urgent"}},
			want:   true,
		},
		{
			name:   "filter matches with case-sensitive tag comparison",
			filter: TaskFilter{Tags: []string{"Work"}},
			task:   &Task{Tags: []string{"Work"}},
			want:   true,
		},
		{
			name:   "filter rejects with different case tags",
			filter: TaskFilter{Tags: []string{"Work"}},
			task:   &Task{Tags: []string{"work"}},
			want:   false,
		},
		{
			name:   "filter matches task with tags in different order",
			filter: TaskFilter{Tags: []string{"work", "urgent"}},
			task:   &Task{Tags: []string{"urgent", "work"}},
			want:   true,
		},
		{
			name:   "filter with empty string tag matches task with empty string tag",
			filter: TaskFilter{Tags: []string{""}},
			task:   &Task{Tags: []string{"", "work"}},
			want:   true,
		},
		{
			name:   "filter with empty string tag rejects task without empty string tag",
			filter: TaskFilter{Tags: []string{""}},
			task:   &Task{Tags: []string{"work"}},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.task)
			if got != tt.want {
				t.Errorf("TaskFilter.Matches() = %v, want %v (filter tags: %v, task tags: %v)",
					got, tt.want, tt.filter.Tags, tt.task.Tags)
			}
		})
	}
}

func TestTaskFilter_Matches_TimeEdgeCases(t *testing.T) {
	// Create times with different precision
	timeWithNanos := time.Date(2025, 11, 9, 12, 0, 0, 123456789, time.UTC)
	timeWithoutNanos := time.Date(2025, 11, 9, 12, 0, 0, 0, time.UTC)
	sameSecondDiffNanos := time.Date(2025, 11, 9, 12, 0, 0, 987654321, time.UTC)

	tests := []struct {
		name   string
		filter TaskFilter
		task   *Task
		want   bool
	}{
		{
			name:   "times with different nanoseconds are treated as different",
			filter: TaskFilter{DueAfter: &timeWithNanos},
			task:   &Task{DueDate: &timeWithoutNanos},
			want:   false,
		},
		{
			name:   "exact time comparison includes nanoseconds",
			filter: TaskFilter{DueAfter: &timeWithNanos, DueBefore: &timeWithNanos},
			task:   &Task{DueDate: &timeWithNanos},
			want:   true,
		},
		{
			name:   "different nanoseconds in same second compare correctly",
			filter: TaskFilter{DueAfter: &timeWithNanos},
			task:   &Task{DueDate: &sameSecondDiffNanos},
			want:   true,
		},
		{
			name: "inverted range (DueAfter > DueBefore) rejects all tasks",
			filter: TaskFilter{
				DueAfter:  &tomorrow,
				DueBefore: &yesterday,
			},
			task: &Task{DueDate: &now},
			want: false,
		},
		{
			name: "inverted range rejects task even at boundaries",
			filter: TaskFilter{
				DueAfter:  &tomorrow,
				DueBefore: &yesterday,
			},
			task: &Task{DueDate: &yesterday},
			want: false,
		},
		{
			name: "equal DueAfter and DueBefore creates single-instant range",
			filter: TaskFilter{
				DueAfter:  &now,
				DueBefore: &now,
			},
			task: &Task{DueDate: &now},
			want: true,
		},
		{
			name: "single-instant range rejects different time",
			filter: TaskFilter{
				DueAfter:  &now,
				DueBefore: &now,
			},
			task: &Task{DueDate: &tomorrow},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.task)
			if got != tt.want {
				t.Errorf("TaskFilter.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskFilter_LimitIsIgnoredByMatches(t *testing.T) {
	tests := []struct {
		name   string
		filter TaskFilter
		task   *Task
		want   bool
	}{
		{
			name: "Limit=0 does not affect matching",
			filter: TaskFilter{
				Status: &statusToday,
				Limit:  0,
			},
			task: &Task{Status: StatusToday},
			want: true,
		},
		{
			name: "Limit=1 does not affect matching",
			filter: TaskFilter{
				Status: &statusToday,
				Limit:  1,
			},
			task: &Task{Status: StatusToday},
			want: true,
		},
		{
			name: "Limit=100 does not affect matching",
			filter: TaskFilter{
				Status: &statusToday,
				Limit:  100,
			},
			task: &Task{Status: StatusToday},
			want: true,
		},
		{
			name: "negative Limit does not affect matching",
			filter: TaskFilter{
				Status: &statusToday,
				Limit:  -1,
			},
			task: &Task{Status: StatusToday},
			want: true,
		},
		{
			name: "Limit does not cause match to fail",
			filter: TaskFilter{
				Status: &statusPool,
				Limit:  5,
			},
			task: &Task{Status: StatusToday},
			want: false, // Fails because of status, not limit
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.task)
			if got != tt.want {
				t.Errorf("TaskFilter.Matches() = %v, want %v (Limit should be ignored, but got different result)",
					got, tt.want)
			}
		})
	}
}
