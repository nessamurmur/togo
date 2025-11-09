package main

import (
	"io"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TestProgramIntegration runs the program (without a renderer), sends
// navigation + toggle + quit key messages, and asserts the returned model
// reflects the actions.
func TestProgramIntegration(t *testing.T) {
	p := tea.NewProgram(initializeModel(), tea.WithoutRenderer(), tea.WithOutput(io.Discard))

	done := make(chan struct{})
	var returnedModel tea.Model
	var runErr error
	go func() {
		returnedModel, runErr = p.Run()
		close(done)
	}()

	// send messages: down (move cursor to index 1), space (toggle selection), then q (quit)
	// small sleeps allow messages to be processed by the program's Update loop.
	p.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	time.Sleep(20 * time.Millisecond)
	p.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	time.Sleep(20 * time.Millisecond)
	p.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})

	select {
	case <-done:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("program did not exit in time")
	}

	if runErr != nil {
		t.Fatalf("program returned error: %v", runErr)
	}

	m, ok := returnedModel.(model)
	if !ok {
		t.Fatalf("returned model has unexpected type")
	}

	// after sending 'j' then space, cursor should be 1 and selected should contain 1
	if m.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", m.cursor)
	}

	if _, ok := m.selected[1]; !ok {
		t.Fatalf("expected item 1 to be selected")
	}
}
