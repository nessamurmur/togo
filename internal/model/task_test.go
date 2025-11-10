package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestNewTask_ValidTitle_Success verifies that NewTask successfully creates
// a task when provided with a valid, non-empty title.
func TestNewTask_ValidTitle_Success(t *testing.T) {
	// Arrange
	title := "Buy groceries"
	tags := []string{"personal"}

	// Act
	task, err := NewTask(title, tags)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if task == nil {
		t.Fatal("expected task to be non-nil")
	}
	if task.Title != title {
		t.Errorf("expected title %q, got %q", title, task.Title)
	}
	if task.ID.IsEmpty() {
		t.Error("expected non-empty task ID")
	}
	if len(task.Tags) != 1 || task.Tags[0] != "personal" {
		t.Errorf("expected tags [\"personal\"], got %v", task.Tags)
	}
}

// TestNewTask_EmptyTitle_ReturnsError verifies that NewTask returns
// ErrEmptyTitle when given an empty string as the title.
func TestNewTask_EmptyTitle_ReturnsError(t *testing.T) {
	// Arrange
	title := ""
	tags := []string{}

	// Act
	task, err := NewTask(title, tags)

	// Assert
	if err != ErrEmptyTitle {
		t.Errorf("expected error %v, got %v", ErrEmptyTitle, err)
	}
	if task != nil {
		t.Errorf("expected nil task, got %+v", task)
	}
}

// TestNewTask_WhitespaceTitle_ReturnsError verifies that NewTask returns
// ErrEmptyTitle when given a title that contains only whitespace characters.
func TestNewTask_WhitespaceTitle_ReturnsError(t *testing.T) {
	tests := []struct {
		name  string
		title string
	}{
		{
			name:  "single space",
			title: " ",
		},
		{
			name:  "multiple spaces",
			title: "   ",
		},
		{
			name:  "tab character",
			title: "\t",
		},
		{
			name:  "newline character",
			title: "\n",
		},
		{
			name:  "mixed whitespace",
			title: " \t\n ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			task, err := NewTask(tt.title, nil)

			// Assert
			if err != ErrEmptyTitle {
				t.Errorf("expected error %v, got %v", ErrEmptyTitle, err)
			}
			if task != nil {
				t.Errorf("expected nil task, got %+v", task)
			}
		})
	}
}

// TestNewTask_SetsDefaultValues verifies that NewTask initializes all
// fields to their correct default values when creating a new task.
func TestNewTask_SetsDefaultValues(t *testing.T) {
	// Arrange
	title := "Test task"
	tags := []string{"work"}
	beforeCreation := time.Now()

	// Act
	task, err := NewTask(title, tags)
	afterCreation := time.Now()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify ID is generated and valid
	if task.ID.IsEmpty() {
		t.Error("expected non-empty ID")
	}
	// Verify ID is a valid UUID
	_, err = uuid.Parse(task.ID.String())
	if err != nil {
		t.Errorf("expected valid UUID for ID, got %v", err)
	}

	// Verify CreatedAt is set to current time (within reasonable range)
	if task.CreatedAt.Before(beforeCreation) || task.CreatedAt.After(afterCreation) {
		t.Errorf("expected CreatedAt between %v and %v, got %v",
			beforeCreation, afterCreation, task.CreatedAt)
	}

	// Verify Title is set correctly
	if task.Title != title {
		t.Errorf("expected Title %q, got %q", title, task.Title)
	}

	// Verify Status defaults to StatusPool
	if task.Status != StatusPool {
		t.Errorf("expected Status %q, got %q", StatusPool, task.Status)
	}

	// Verify Notes defaults to empty string
	if task.Notes != "" {
		t.Errorf("expected empty Notes, got %q", task.Notes)
	}

	// Verify DeferredCount defaults to 0
	if task.DeferredCount != 0 {
		t.Errorf("expected DeferredCount 0, got %d", task.DeferredCount)
	}

	// Verify DueDate defaults to nil
	if task.DueDate != nil {
		t.Errorf("expected nil DueDate, got %v", task.DueDate)
	}

	// Verify CompletedAt defaults to nil
	if task.CompletedAt != nil {
		t.Errorf("expected nil CompletedAt, got %v", task.CompletedAt)
	}

	// Verify Tags are set correctly
	if len(task.Tags) != 1 || task.Tags[0] != "work" {
		t.Errorf("expected Tags [\"work\"], got %v", task.Tags)
	}
}

// TestNewTask_TrimsTitleWhitespace verifies that NewTask trims leading and
// trailing whitespace from the title while preserving internal whitespace.
func TestNewTask_TrimsTitleWhitespace(t *testing.T) {
	tests := []struct {
		name          string
		inputTitle    string
		expectedTitle string
	}{
		{
			name:          "leading space",
			inputTitle:    " Buy milk",
			expectedTitle: "Buy milk",
		},
		{
			name:          "trailing space",
			inputTitle:    "Buy milk ",
			expectedTitle: "Buy milk",
		},
		{
			name:          "both leading and trailing spaces",
			inputTitle:    "  Buy milk  ",
			expectedTitle: "Buy milk",
		},
		{
			name:          "leading tab",
			inputTitle:    "\tBuy milk",
			expectedTitle: "Buy milk",
		},
		{
			name:          "mixed whitespace",
			inputTitle:    " \t Buy milk \n ",
			expectedTitle: "Buy milk",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			task, err := NewTask(tt.inputTitle, nil)

			// Assert
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if task.Title != tt.expectedTitle {
				t.Errorf("expected title %q, got %q", tt.expectedTitle, task.Title)
			}
		})
	}
}

// TestNewTask_TitleWithInternalWhitespace_Preserved verifies that NewTask
// preserves internal whitespace within the title while trimming edges.
func TestNewTask_TitleWithInternalWhitespace_Preserved(t *testing.T) {
	tests := []struct {
		name          string
		inputTitle    string
		expectedTitle string
	}{
		{
			name:          "single space between words",
			inputTitle:    "Buy some milk",
			expectedTitle: "Buy some milk",
		},
		{
			name:          "multiple spaces between words",
			inputTitle:    "Buy  some   milk",
			expectedTitle: "Buy  some   milk",
		},
		{
			name:          "tab between words",
			inputTitle:    "Buy\tsome\tmilk",
			expectedTitle: "Buy\tsome\tmilk",
		},
		{
			name:          "mixed internal whitespace with edge trimming",
			inputTitle:    "  Buy  some   milk  ",
			expectedTitle: "Buy  some   milk",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			task, err := NewTask(tt.inputTitle, nil)

			// Assert
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if task.Title != tt.expectedTitle {
				t.Errorf("expected title %q, got %q", tt.expectedTitle, task.Title)
			}
		})
	}
}

// TestNewTask_MultipleTags verifies that NewTask correctly handles multiple
// tags and stores them in the task.
func TestNewTask_MultipleTags(t *testing.T) {
	// Arrange
	title := "Complete project report"
	tags := []string{"work", "urgent", "documentation"}

	// Act
	task, err := NewTask(title, tags)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(task.Tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(task.Tags))
	}
	expectedTags := []string{"work", "urgent", "documentation"}
	for i, expectedTag := range expectedTags {
		if task.Tags[i] != expectedTag {
			t.Errorf("expected tag[%d] = %q, got %q", i, expectedTag, task.Tags[i])
		}
	}
}

// TestNewTask_EmptyTagsSlice verifies that NewTask stores nil for tags
// when given an empty slice, to ensure proper JSON omitempty behavior.
func TestNewTask_EmptyTagsSlice(t *testing.T) {
	tests := []struct {
		name string
		tags []string
	}{
		{
			name: "nil tags",
			tags: nil,
		},
		{
			name: "empty slice",
			tags: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			task, err := NewTask("Test task", tt.tags)

			// Assert
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if task.Tags != nil {
				t.Errorf("expected nil tags for empty slice, got %v", task.Tags)
			}
		})
	}
}

// TestNewTask_DefensiveCopyOfTags verifies that NewTask makes a defensive
// copy of the tags slice so that external modifications don't affect the task.
func TestNewTask_DefensiveCopyOfTags(t *testing.T) {
	// Arrange
	title := "Test task"
	originalTags := []string{"tag1", "tag2"}
	tags := make([]string, len(originalTags))
	copy(tags, originalTags)

	// Act
	task, err := NewTask(title, tags)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Modify the original slice
	tags[0] = "modified"
	tags[1] = "changed"

	// Assert - task's tags should be unchanged
	if len(task.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(task.Tags))
	}
	if task.Tags[0] != "tag1" {
		t.Errorf("expected tag[0] = %q, got %q (defensive copy failed)", "tag1", task.Tags[0])
	}
	if task.Tags[1] != "tag2" {
		t.Errorf("expected tag[1] = %q, got %q (defensive copy failed)", "tag2", task.Tags[1])
	}
}

// TestNewTask_GeneratesUniqueIDs verifies that NewTask generates unique
// UUIDs for each task, ensuring no ID collisions.
func TestNewTask_GeneratesUniqueIDs(t *testing.T) {
	// Arrange
	title := "Test task"
	numTasks := 100

	// Act - create multiple tasks
	tasks := make([]*Task, numTasks)
	for i := 0; i < numTasks; i++ {
		task, err := NewTask(title, nil)
		if err != nil {
			t.Fatalf("failed to create task %d: %v", i, err)
		}
		tasks[i] = task
	}

	// Assert - verify all IDs are unique
	idSet := make(map[string]bool)
	for i, task := range tasks {
		idStr := task.ID.String()
		if idSet[idStr] {
			t.Errorf("duplicate ID found at task %d: %s", i, idStr)
		}
		idSet[idStr] = true
	}

	if len(idSet) != numTasks {
		t.Errorf("expected %d unique IDs, got %d", numTasks, len(idSet))
	}
}

// TestTask_JSONRoundTrip verifies that a Task can be marshaled to JSON
// and unmarshaled back correctly, preserving all field values.
func TestTask_JSONRoundTrip(t *testing.T) {
	// Arrange
	now := time.Now().UTC().Round(time.Millisecond) // Round for JSON precision
	dueDate := now.Add(24 * time.Hour)
	completedAt := now.Add(1 * time.Hour)

	task := &Task{
		ID:            TaskID(uuid.New()),
		CreatedAt:     now,
		Title:         "Test task with all fields",
		Notes:         "Some important notes",
		Status:        StatusToday,
		Tags:          []string{"work", "urgent"},
		DueDate:       &dueDate,
		CompletedAt:   &completedAt,
		DeferredCount: 3,
	}

	// Act - Marshal to JSON
	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("failed to marshal task: %v", err)
	}

	// Act - Unmarshal back to Task
	var unmarshaled Task
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal task: %v", err)
	}

	// Assert - Verify all fields match
	if unmarshaled.ID.String() != task.ID.String() {
		t.Errorf("expected ID %q, got %q", task.ID.String(), unmarshaled.ID.String())
	}
	if !unmarshaled.CreatedAt.Equal(task.CreatedAt) {
		t.Errorf("expected CreatedAt %v, got %v", task.CreatedAt, unmarshaled.CreatedAt)
	}
	if unmarshaled.Title != task.Title {
		t.Errorf("expected Title %q, got %q", task.Title, unmarshaled.Title)
	}
	if unmarshaled.Notes != task.Notes {
		t.Errorf("expected Notes %q, got %q", task.Notes, unmarshaled.Notes)
	}
	if unmarshaled.Status != task.Status {
		t.Errorf("expected Status %q, got %q", task.Status, unmarshaled.Status)
	}
	if len(unmarshaled.Tags) != len(task.Tags) {
		t.Fatalf("expected %d tags, got %d", len(task.Tags), len(unmarshaled.Tags))
	}
	for i := range task.Tags {
		if unmarshaled.Tags[i] != task.Tags[i] {
			t.Errorf("expected Tags[%d] = %q, got %q", i, task.Tags[i], unmarshaled.Tags[i])
		}
	}
	if unmarshaled.DueDate == nil || !unmarshaled.DueDate.Equal(*task.DueDate) {
		t.Errorf("expected DueDate %v, got %v", task.DueDate, unmarshaled.DueDate)
	}
	if unmarshaled.CompletedAt == nil || !unmarshaled.CompletedAt.Equal(*task.CompletedAt) {
		t.Errorf("expected CompletedAt %v, got %v", task.CompletedAt, unmarshaled.CompletedAt)
	}
	if unmarshaled.DeferredCount != task.DeferredCount {
		t.Errorf("expected DeferredCount %d, got %d", task.DeferredCount, unmarshaled.DeferredCount)
	}
}

// TestTask_JSONMarshal_OmitsEmptyFields verifies that optional fields
// with omitempty tags are excluded from JSON when they are nil or empty.
func TestTask_JSONMarshal_OmitsEmptyFields(t *testing.T) {
	// Arrange - create minimal task with only required fields
	now := time.Now().UTC().Round(time.Millisecond)
	task := &Task{
		ID:            TaskID(uuid.New()),
		CreatedAt:     now,
		Title:         "Minimal task",
		Status:        StatusPool,
		Notes:         "", // Empty, should be omitted
		Tags:          nil,
		DueDate:       nil,
		CompletedAt:   nil,
		DeferredCount: 0,
	}

	// Act
	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("failed to marshal task: %v", err)
	}

	// Assert - Unmarshal to map to check field presence
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	// Verify omitempty fields are not present
	omitFields := []string{"notes", "tags", "due_date", "completed_at"}
	for _, field := range omitFields {
		if _, exists := jsonMap[field]; exists {
			t.Errorf("expected field %q to be omitted, but it was present", field)
		}
	}

	// Verify required fields are present
	requiredFields := []string{"id", "created_at", "title", "status", "deferred_count"}
	for _, field := range requiredFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("expected field %q to be present, but it was omitted", field)
		}
	}
}

// TestTask_JSONMarshal_IncludesNonEmptyFields verifies that optional fields
// are included in JSON when they have non-empty values.
func TestTask_JSONMarshal_IncludesNonEmptyFields(t *testing.T) {
	// Arrange - create task with all optional fields populated
	now := time.Now().UTC().Round(time.Millisecond)
	dueDate := now.Add(24 * time.Hour)
	completedAt := now.Add(1 * time.Hour)

	task := &Task{
		ID:            TaskID(uuid.New()),
		CreatedAt:     now,
		Title:         "Complete task",
		Notes:         "Important notes",
		Status:        StatusDone,
		Tags:          []string{"work"},
		DueDate:       &dueDate,
		CompletedAt:   &completedAt,
		DeferredCount: 2,
	}

	// Act
	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("failed to marshal task: %v", err)
	}

	// Assert - Unmarshal to map to check field presence
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	// Verify all fields are present
	allFields := []string{"id", "created_at", "title", "notes", "status", "tags", "due_date", "completed_at", "deferred_count"}
	for _, field := range allFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("expected field %q to be present, but it was omitted", field)
		}
	}

	// Verify specific values
	if jsonMap["notes"] != "Important notes" {
		t.Errorf("expected notes %q, got %v", "Important notes", jsonMap["notes"])
	}
	if jsonMap["status"] != "done" {
		t.Errorf("expected status %q, got %v", "done", jsonMap["status"])
	}
	if jsonMap["deferred_count"] != float64(2) {
		t.Errorf("expected deferred_count 2, got %v", jsonMap["deferred_count"])
	}
}
