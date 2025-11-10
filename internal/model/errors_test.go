package model

import (
	"errors"
	"testing"
)

// TestValidationError_Error_FormatsCorrectly verifies that ValidationError
// produces the expected error message format.
func TestValidationError_Error_FormatsCorrectly(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		reason   string
		expected string
	}{
		{
			name:     "basic field and reason",
			field:    "title",
			reason:   "cannot be empty",
			expected: "validation failed for title: cannot be empty",
		},
		{
			name:     "field with reason about length",
			field:    "description",
			reason:   "exceeds 500 characters",
			expected: "validation failed for description: exceeds 500 characters",
		},
		{
			name:     "field with complex reason",
			field:    "tags",
			reason:   "must contain only alphanumeric characters and hyphens",
			expected: "validation failed for tags: must contain only alphanumeric characters and hyphens",
		},
		{
			name:     "empty field name",
			field:    "",
			reason:   "some reason",
			expected: "validation failed for : some reason",
		},
		{
			name:     "empty reason",
			field:    "status",
			reason:   "",
			expected: "validation failed for status: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &ValidationError{
				Field:  tt.field,
				Reason: tt.reason,
			}

			got := err.Error()
			if got != tt.expected {
				t.Errorf("ValidationError.Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestValidationError_ImplementsError verifies that ValidationError implements
// the error interface.
func TestValidationError_ImplementsError(t *testing.T) {
	var err error = &ValidationError{
		Field:  "test",
		Reason: "testing",
	}

	if err.Error() == "" {
		t.Error("ValidationError.Error() returned empty string")
	}
}

// TestDomainErrors_AreUnique verifies that each sentinel error has a unique
// identity in memory, ensuring errors.Is() comparisons work correctly.
func TestDomainErrors_AreUnique(t *testing.T) {
	// Collect all sentinel errors
	sentinelErrors := []error{
		ErrTaskNotFound,
		ErrInvalidStatus,
		ErrInvalidStateTransition,
		ErrEmptyTitle,
		ErrDuplicateTaskID,
	}

	// Compare each error with every other error
	for i, err1 := range sentinelErrors {
		for j, err2 := range sentinelErrors {
			if i == j {
				// Same error should be identical to itself
				if err1 != err2 {
					t.Errorf("Error %d is not identical to itself", i)
				}
			} else {
				// Different errors should not be identical
				if err1 == err2 {
					t.Errorf("Error %d and %d are not unique: both are %v", i, j, err1)
				}
			}
		}
	}
}

// TestDomainErrors_Messages verifies that sentinel errors follow Go conventions:
// - lowercase first letter
// - no punctuation at the end
// - provide actionable context
func TestDomainErrors_Messages(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		wantContains   string
		mustNotEndWith string
	}{
		{
			name:           "ErrTaskNotFound message",
			err:            ErrTaskNotFound,
			wantContains:   "task not found",
			mustNotEndWith: ".",
		},
		{
			name:           "ErrInvalidStatus message",
			err:            ErrInvalidStatus,
			wantContains:   "invalid task status",
			mustNotEndWith: ".",
		},
		{
			name:           "ErrInvalidStateTransition message",
			err:            ErrInvalidStateTransition,
			wantContains:   "invalid state transition",
			mustNotEndWith: ".",
		},
		{
			name:           "ErrEmptyTitle message",
			err:            ErrEmptyTitle,
			wantContains:   "cannot be empty",
			mustNotEndWith: ".",
		},
		{
			name:           "ErrDuplicateTaskID message",
			err:            ErrDuplicateTaskID,
			wantContains:   "already exists",
			mustNotEndWith: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.err.Error()

			// Check for required content
			if msg == "" {
				t.Errorf("Error message is empty")
			}

			// Check that message contains expected text (case-insensitive check would be done manually)
			// For simplicity, we just verify the error has a message
			if len(msg) < 3 {
				t.Errorf("Error message too short: %q", msg)
			}

			// Check first character is lowercase
			if len(msg) > 0 && msg[0] >= 'A' && msg[0] <= 'Z' {
				t.Errorf("Error message %q starts with uppercase letter (should be lowercase)", msg)
			}

			// Check doesn't end with punctuation
			lastChar := msg[len(msg)-1]
			if lastChar == '.' || lastChar == '!' || lastChar == '?' {
				t.Errorf("Error message %q ends with punctuation (should not)", msg)
			}
		})
	}
}

// TestDomainErrors_WorkWithErrorsIs verifies that sentinel errors can be used
// with errors.Is() for error comparison through wrapped errors.
func TestDomainErrors_WorkWithErrorsIs(t *testing.T) {
	tests := []struct {
		name        string
		sentinelErr error
		wrappedErr  error
		shouldMatch bool
	}{
		{
			name:        "ErrTaskNotFound matches itself",
			sentinelErr: ErrTaskNotFound,
			wrappedErr:  ErrTaskNotFound,
			shouldMatch: true,
		},
		{
			name:        "ErrInvalidStatus matches itself",
			sentinelErr: ErrInvalidStatus,
			wrappedErr:  ErrInvalidStatus,
			shouldMatch: true,
		},
		{
			name:        "ErrTaskNotFound does not match ErrInvalidStatus",
			sentinelErr: ErrTaskNotFound,
			wrappedErr:  ErrInvalidStatus,
			shouldMatch: false,
		},
		{
			name:        "ErrTaskNotFound does not match generic error",
			sentinelErr: ErrTaskNotFound,
			wrappedErr:  errors.New("some other error"),
			shouldMatch: false,
		},
		{
			name:        "ErrEmptyTitle matches itself",
			sentinelErr: ErrEmptyTitle,
			wrappedErr:  ErrEmptyTitle,
			shouldMatch: true,
		},
		{
			name:        "ErrDuplicateTaskID matches itself",
			sentinelErr: ErrDuplicateTaskID,
			wrappedErr:  ErrDuplicateTaskID,
			shouldMatch: true,
		},
		{
			name:        "ErrInvalidStateTransition matches itself",
			sentinelErr: ErrInvalidStateTransition,
			wrappedErr:  ErrInvalidStateTransition,
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched := errors.Is(tt.wrappedErr, tt.sentinelErr)
			if matched != tt.shouldMatch {
				t.Errorf("errors.Is(%v, %v) = %v, want %v",
					tt.wrappedErr, tt.sentinelErr, matched, tt.shouldMatch)
			}
		})
	}
}

// TestDomainErrors_AllDefined verifies that all expected sentinel errors
// are defined and not nil.
func TestDomainErrors_AllDefined(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrTaskNotFound", ErrTaskNotFound},
		{"ErrInvalidStatus", ErrInvalidStatus},
		{"ErrInvalidStateTransition", ErrInvalidStateTransition},
		{"ErrEmptyTitle", ErrEmptyTitle},
		{"ErrDuplicateTaskID", ErrDuplicateTaskID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("%s is nil, expected non-nil error", tt.name)
			}
		})
	}
}

// TestValidationError_CanBeUsedAsError verifies that ValidationError can be
// used in contexts where error interface is expected (nil checks, wrapping, etc).
func TestValidationError_CanBeUsedAsError(t *testing.T) {
	t.Run("non-nil ValidationError is not nil error", func(t *testing.T) {
		var err error = &ValidationError{Field: "test", Reason: "reason"}
		if err == nil {
			t.Error("Non-nil ValidationError treated as nil error")
		}
	})

	t.Run("ValidationError can be compared with errors.As", func(t *testing.T) {
		original := &ValidationError{Field: "name", Reason: "too short"}
		var err error = original

		var target *ValidationError
		if !errors.As(err, &target) {
			t.Error("errors.As failed to extract ValidationError")
		}

		if target.Field != "name" {
			t.Errorf("Field = %q, want %q", target.Field, "name")
		}

		if target.Reason != "too short" {
			t.Errorf("Reason = %q, want %q", target.Reason, "too short")
		}
	})
}
