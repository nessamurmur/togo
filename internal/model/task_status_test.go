package model

import (
	"encoding/json"
	"testing"
)

// TestTaskStatus_Valid_AllValidStatuses verifies that all three valid statuses
// (pool, today, done) are recognized as valid by the Valid() method.
func TestTaskStatus_Valid_AllValidStatuses(t *testing.T) {
	tests := []struct {
		name   string
		status TaskStatus
		want   bool
	}{
		{
			name:   "StatusPool is valid",
			status: StatusPool,
			want:   true,
		},
		{
			name:   "StatusToday is valid",
			status: StatusToday,
			want:   true,
		},
		{
			name:   "StatusDone is valid",
			status: StatusDone,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.Valid()
			if got != tt.want {
				t.Errorf("TaskStatus.Valid() = %v, want %v for status %q", got, tt.want, tt.status)
			}
		})
	}
}

// TestTaskStatus_Valid_InvalidStatus_ReturnsFalse verifies that invalid status
// values are correctly identified as invalid by the Valid() method.
func TestTaskStatus_Valid_InvalidStatus_ReturnsFalse(t *testing.T) {
	tests := []struct {
		name   string
		status TaskStatus
		want   bool
	}{
		{
			name:   "empty string is invalid",
			status: TaskStatus(""),
			want:   false,
		},
		{
			name:   "random string is invalid",
			status: TaskStatus("invalid"),
			want:   false,
		},
		{
			name:   "uppercase POOL is invalid",
			status: TaskStatus("POOL"),
			want:   false,
		},
		{
			name:   "uppercase TODAY is invalid",
			status: TaskStatus("TODAY"),
			want:   false,
		},
		{
			name:   "uppercase DONE is invalid",
			status: TaskStatus("DONE"),
			want:   false,
		},
		{
			name:   "mixed case Pool is invalid",
			status: TaskStatus("Pool"),
			want:   false,
		},
		{
			name:   "misspelled poo is invalid",
			status: TaskStatus("poo"),
			want:   false,
		},
		{
			name:   "whitespace with valid status is invalid",
			status: TaskStatus(" pool "),
			want:   false,
		},
		{
			name:   "numeric value is invalid",
			status: TaskStatus("123"),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.Valid()
			if got != tt.want {
				t.Errorf("TaskStatus.Valid() = %v, want %v for status %q", got, tt.want, tt.status)
			}
		})
	}
}

// TestTaskStatus_String_ReturnsValue verifies that the String() method returns
// the correct string representation of the TaskStatus.
func TestTaskStatus_String_ReturnsValue(t *testing.T) {
	tests := []struct {
		name   string
		status TaskStatus
		want   string
	}{
		{
			name:   "StatusPool returns 'pool'",
			status: StatusPool,
			want:   "pool",
		},
		{
			name:   "StatusToday returns 'today'",
			status: StatusToday,
			want:   "today",
		},
		{
			name:   "StatusDone returns 'done'",
			status: StatusDone,
			want:   "done",
		},
		{
			name:   "custom invalid status returns its value",
			status: TaskStatus("invalid"),
			want:   "invalid",
		},
		{
			name:   "empty status returns empty string",
			status: TaskStatus(""),
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.String()
			if got != tt.want {
				t.Errorf("TaskStatus.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestTaskStatus_JSONRoundTrip verifies that TaskStatus can be marshaled to JSON
// and unmarshaled back correctly, preserving the original value.
func TestTaskStatus_JSONRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		status   TaskStatus
		wantJSON string
	}{
		{
			name:     "StatusPool marshals to 'pool'",
			status:   StatusPool,
			wantJSON: `"pool"`,
		},
		{
			name:     "StatusToday marshals to 'today'",
			status:   StatusToday,
			wantJSON: `"today"`,
		},
		{
			name:     "StatusDone marshals to 'done'",
			status:   StatusDone,
			wantJSON: `"done"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			jsonData, err := json.Marshal(tt.status)
			if err != nil {
				t.Fatalf("failed to marshal TaskStatus: %v", err)
			}

			gotJSON := string(jsonData)
			if gotJSON != tt.wantJSON {
				t.Errorf("json.Marshal() = %q, want %q", gotJSON, tt.wantJSON)
			}

			// Test unmarshaling
			var unmarshaled TaskStatus
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Fatalf("failed to unmarshal TaskStatus: %v", err)
			}

			if unmarshaled != tt.status {
				t.Errorf("after unmarshal = %q, want %q", unmarshaled, tt.status)
			}
		})
	}
}

// TestTaskStatus_JSONRoundTrip_InStruct verifies that TaskStatus works correctly
// when embedded in a struct during JSON marshaling/unmarshaling.
func TestTaskStatus_JSONRoundTrip_InStruct(t *testing.T) {
	type TestTask struct {
		ID     string     `json:"id"`
		Status TaskStatus `json:"status"`
		Title  string     `json:"title"`
	}

	tests := []struct {
		name     string
		task     TestTask
		wantJSON string
	}{
		{
			name: "task with pool status",
			task: TestTask{
				ID:     "123",
				Status: StatusPool,
				Title:  "Test Task",
			},
			wantJSON: `{"id":"123","status":"pool","title":"Test Task"}`,
		},
		{
			name: "task with today status",
			task: TestTask{
				ID:     "456",
				Status: StatusToday,
				Title:  "Another Task",
			},
			wantJSON: `{"id":"456","status":"today","title":"Another Task"}`,
		},
		{
			name: "task with done status",
			task: TestTask{
				ID:     "789",
				Status: StatusDone,
				Title:  "Completed Task",
			},
			wantJSON: `{"id":"789","status":"done","title":"Completed Task"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			jsonData, err := json.Marshal(tt.task)
			if err != nil {
				t.Fatalf("failed to marshal task: %v", err)
			}

			gotJSON := string(jsonData)
			if gotJSON != tt.wantJSON {
				t.Errorf("json.Marshal() = %q, want %q", gotJSON, tt.wantJSON)
			}

			// Test unmarshaling
			var unmarshaled TestTask
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Fatalf("failed to unmarshal task: %v", err)
			}

			if unmarshaled.Status != tt.task.Status {
				t.Errorf("after unmarshal status = %q, want %q", unmarshaled.Status, tt.task.Status)
			}
			if unmarshaled.ID != tt.task.ID {
				t.Errorf("after unmarshal ID = %q, want %q", unmarshaled.ID, tt.task.ID)
			}
			if unmarshaled.Title != tt.task.Title {
				t.Errorf("after unmarshal Title = %q, want %q", unmarshaled.Title, tt.task.Title)
			}
		})
	}
}

// TestTaskStatus_Constants_HaveCorrectValues verifies that the constant values
// are set to the expected string literals.
func TestTaskStatus_Constants_HaveCorrectValues(t *testing.T) {
	tests := []struct {
		name     string
		constant TaskStatus
		want     string
	}{
		{
			name:     "StatusPool constant value",
			constant: StatusPool,
			want:     "pool",
		},
		{
			name:     "StatusToday constant value",
			constant: StatusToday,
			want:     "today",
		},
		{
			name:     "StatusDone constant value",
			constant: StatusDone,
			want:     "done",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(tt.constant)
			if got != tt.want {
				t.Errorf("constant value = %q, want %q", got, tt.want)
			}
		})
	}
}
