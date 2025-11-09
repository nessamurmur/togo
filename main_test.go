package main

import (
    "strings"
    "testing"

    tea "github.com/charmbracelet/bubbletea"
)

// helper to build rune-based key messages used in tests (e.g. "j", "k", "q", " ")
func keyMsg(s string) tea.KeyMsg {
    return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func TestInitializeModel(t *testing.T) {
    m := initializeModel()
    if len(m.choices) != 3 {
        t.Fatalf("expected 3 choices, got %d", len(m.choices))
    }
    if len(m.selected) != 0 {
        t.Fatalf("expected selected to be empty, got %v", m.selected)
    }
}

func TestNavigationBounds(t *testing.T) {
    m := initializeModel()

    // at the top, pressing 'k' (up) should not move the cursor
    nm, _ := m.Update(keyMsg("k"))
    got := nm.(model)
    if got.cursor != 0 {
        t.Fatalf("expected cursor 0 after up at top, got %d", got.cursor)
    }

    // press 'j' (down) twice to move to the last item
    nm2, _ := got.Update(keyMsg("j"))
    got2 := nm2.(model)
    if got2.cursor != 1 {
        t.Fatalf("expected cursor 1 after one down, got %d", got2.cursor)
    }

    nm3, _ := got2.Update(keyMsg("j"))
    got3 := nm3.(model)
    if got3.cursor != 2 {
        t.Fatalf("expected cursor 2 after two downs, got %d", got3.cursor)
    }

    // one more down should not advance past last index
    nm4, _ := got3.Update(keyMsg("j"))
    got4 := nm4.(model)
    if got4.cursor != 2 {
        t.Fatalf("expected cursor to remain 2 at bottom, got %d", got4.cursor)
    }
}

func TestToggleSelection(t *testing.T) {
    m := initializeModel()

    // toggle select the first item using space
    nm, _ := m.Update(keyMsg(" "))
    got := nm.(model)
    if _, ok := got.selected[0]; !ok {
        t.Fatalf("expected item 0 to be selected")
    }

    // toggle again to deselect
    nm2, _ := got.Update(keyMsg(" "))
    got2 := nm2.(model)
    if _, ok := got2.selected[0]; ok {
        t.Fatalf("expected item 0 to be deselected")
    }
}

func TestQuitCommand(t *testing.T) {
    m := initializeModel()

    _, cmd := m.Update(keyMsg("q"))
    if cmd == nil {
        t.Fatalf("expected non-nil command when pressing 'q' to quit")
    }
}

func TestViewRendering(t *testing.T) {
    m := initializeModel()
    // move cursor to second item and select third
    nm, _ := m.Update(keyMsg("j")) // cursor -> 1
    got := nm.(model)
    // select third (move down then press space)
    nm2, _ := got.Update(keyMsg("j")) // cursor -> 2
    got2 := nm2.(model)
    nm3, _ := got2.Update(keyMsg(" ")) // select index 2
    got3 := nm3.(model)

    view := got3.View()

    // expect cursor marker on the last line for cursor==2
    if !strings.Contains(view, "> [" ) {
        t.Fatalf("expected view to contain a cursor marker '>' somewhere; got:\n%s", view)
    }

    // expect the selected marker 'x' for the third item
    if !strings.Contains(view, "[x] Dream") {
        t.Fatalf("expected view to show selected item for 'Dream'; got:\n%s", view)
    }
}
