package model

import (
	"testing"
)

func TestNewTaskID_GeneratesValidUUID(t *testing.T) {
	taskID := NewTaskID()

	if taskID.IsEmpty() {
		t.Fatalf("expected NewTaskID to generate a non-empty UUID")
	}
}

func TestParseTaskID_ValidUUID(t *testing.T) {
	idStr := "550e8400-e29b-41d4-a716-446655440000"
	taskID, err := ParseTaskID(idStr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if taskID.String() != idStr {
		t.Fatalf("expected TaskID string to be %s, got %s", idStr, taskID.String())
	}
}

func TestParseTaskID_InvalidUUID(t *testing.T) {
	invalidIDStr := "invalid-uuid-string"
	_, err := ParseTaskID(invalidIDStr)
	if err == nil {
		t.Fatalf("expected error for invalid UUID string, got nil")
	}
}

func TestIsEmpty(t *testing.T) {
	emptyID := TaskID{}
	if !emptyID.IsEmpty() {
		t.Fatalf("expected empty TaskID to be empty")
	}

	nonEmptyID := NewTaskID()
	if nonEmptyID.IsEmpty() {
		t.Fatalf("expected non-empty TaskID to not be empty")
	}
}

func TestNotEquals_DifferentIDs(t *testing.T) {
	taskID1 := NewTaskID()
	taskID2 := NewTaskID()

	if taskID1.Equals(taskID2) {
		t.Fatalf("expected TaskIDs to not be equal")
	}

	if !taskID1.NotEquals(taskID2) {
		t.Fatalf("expected TaskIDs to be unequal")
	}
}

func TestParseTaskID_EmptyUUID(t *testing.T) {
	emptyIDStr := "00000000-0000-0000-0000-000000000000"
	taskID, err := ParseTaskID(emptyIDStr)
	if err != nil {
		t.Fatalf("expected no error for empty UUID string, got %v", err)
	}

	if !taskID.IsEmpty() {
		t.Fatalf("expected parsed TaskID to be empty")
	}
}
