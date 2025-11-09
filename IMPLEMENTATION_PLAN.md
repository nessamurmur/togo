# Implementation Plan — togo (Bullet Journal TUI)

This document captures the implementation plan for a TUI-first todo app that simulates basic bullet journaling. It uses libraries from the Charmbracelet ecosystem (Bubble Tea, Bubbles, Lipgloss) and aims to keep core logic UI-agnostic. Optional CLI commands are provided for quick operations via Go's standard library flag package.

## Goals

- Add tasks to a pool of tasks
- Pick tasks for "Today"
- Defer tasks back into the pool
- Remove tasks
- Generate a weekly report of completed tasks (every Monday morning)
- Support syncing an encrypted todo list (user-manageable keys)

## High-level approach

- Keep a small, well-tested core model (Task + business logic) independent from the UI.
- Build a beautiful, interactive TUI as the primary interface (Bubble Tea).
- Provide optional CLI commands for quick operations (standard library flag parsing).
- Store data locally in an encrypted payload; provide sync adapters (git-backed being the MVP).

## Technology choices

- **Primary UI**: `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/bubbles`, `github.com/charmbracelet/lipgloss`
- **CLI parsing**: Go standard library `flag` package (lightweight, no external dependencies)
- **Encryption**: `filippo.io/age` for modern, well-maintained encryption; for passphrase mode derive a symmetric key with `scrypt` or `argon2id`.
- **Storage (MVP)**: encrypted JSON file at `$XDG_DATA_HOME/togo/todo.json.age` (easy to implement and inspect). Optionally swap to `bbolt` for robust local DB with encrypted payloads if needed.
- **Sync (MVP)**: git-backed remote that stores only the encrypted blob (`todo.json.age`) allowing push/pull via normal git remotes.

## Data model

Task (initial fields):

- `ID string` (uuid)
- `Title string`
- `Notes string`
- `CreatedAt time.Time`
- `DueDate *time.Time` (optional)
- `Status string` enum: `pool | today | done`
- `CompletedAt *time.Time` (optional)
- `DeferredCount int`
- `Tags []string`

Storage file shape:

```json
{
  "meta": { ... },
  "tasks": [ ... ]
}
```

The entire JSON payload will be encrypted before being written to disk or pushed to a remote.

## Primary Interface: TUI (Bubble Tea)

The TUI is the primary interface for togo, providing a beautiful, interactive experience:

### Main Views

1. **List View (Tab-Based Navigation)**
   - Three tabs: Pool / Today / Done
   - Navigate with Tab/Shift+Tab
   - Keyboard-driven selection (j/k or arrow keys)
   - Real-time filtering with /
   - Visual status indicators (icons and colors)

2. **Add Task View**
   - Text input for title (bubbles/textinput)
   - Optional tags input (comma-separated)
   - Optional notes input
   - Enter to save, Esc to cancel

3. **Task Detail View**
   - Show full task information
   - Edit notes (bubbles/textarea)
   - View metadata (created, deferred count, due date)
   - Quick actions available

### Keybindings

**Navigation:**
- `j` / `k` or `↑` / `↓` — navigate list
- `tab` / `shift+tab` — switch between Pool/Today/Done tabs
- `/` — filter tasks
- `?` — toggle help

**Task Operations:**
- `a` — add new task
- `enter` — pick task (Pool → Today) or view details
- `x` — complete task
- `d` — defer task (Today → Pool)
- `delete` — remove task (with confirmation)
- `e` — edit task notes

**System Operations:**
- `s` — sync (push/pull with visual feedback)
- `r` — generate report
- `q` / `ctrl+c` — quit

### Components Used

- **bubbles/list** — Main task lists with custom delegate for styling
- **bubbles/textinput** — Single-line inputs (title, tags)
- **bubbles/textarea** — Multi-line notes editing
- **bubbles/spinner** — Loading states during sync
- **bubbles/help** — Contextual help footer
- **lipgloss** — All styling, layouts, and color schemes

### Visual Design

- Status-based colors (Pool: blue, Today: pink, Done: green)
- Rounded borders and boxes
- Responsive layout adapts to terminal width
- Tab navigation bar at top
- Status bar showing counts (Pool: 5 | Today: 3 | Done: 12)
- Help footer with key bindings
- Smooth transitions and animations

## Optional CLI Commands

For quick operations without launching TUI:

- `todo` — launch interactive TUI (default)
- `todo add "title" [-n notes] [-t tag,...]` — add a task to the pool
- `todo list [--pool|--today|--done] [--tag x]` — list tasks (text output)
- `todo pick <id>` — move task to `today`
- `todo defer <id>` — move task back to `pool`
- `todo complete <id>` — mark done
- `todo remove <id>` — delete task
- `todo report [--since 7d] [--format text|json]` — show completed tasks
- `todo sync <push|pull|status>` — sync encrypted storage
- `todo help` — show usage

CLI commands use Go's standard library `flag` package for parsing.

## Storage & encryption

- MVP: JSON payload encrypted with `age` before writing to disk. Metadata contains:
  - encryption mode (age pubkey vs passphrase)
  - salt (if passphrase-derived key)
  - last modified timestamp
- Provide utility commands to `todo key generate` and `todo key import` for managing age keys or passphrases.

## Sync strategy (MVP: git-backed)

- The local repo stores only the encrypted file `todo.json.age`.
- Sync adapter operations:
  - `status`: check if local and remote hashes differ
  - `pull`: `git pull` -> decrypt -> merge -> write decrypted to memory -> re-encrypt and save locally
  - `push`: encrypt local payload -> `git add/commit/push`
- Conflict handling (initial): last-write-wins with a user-visible conflict file; later add interactive merge UI.

## Weekly report automation

- Provide a `todo report --since 7d` command that lists/exports completed tasks in the last 7 days.
- Document a cron/systemd-timer that runs every Monday at a preferred time and either prints to stdout, sends email, or writes a local report file.

Cron example (runs every Monday at 08:00):

```cron
0 8 * * MON /usr/bin/todo report --since 7d --format text > $HOME/todo/last_week_report.txt
```

Optionally add SMTP/email integration to mail the report.

## Testing strategy

- Unit tests:
  - Task business logic: pick/defer/complete/remove
  - Storage interface: Load/Save (with in-memory FS or temp files)
  - Report generation: filter completed tasks by date range
- TUI tests: Bubble Tea `Update/View` tests that send `tea.KeyMsg` messages and assert model state/view strings
- Integration tests:
  - Sync: use a temporary git repo to simulate push/pull cycle

CI should run `gofmt`, `go vet`, and `go test -v ./...` (we already set up a workflow for that).

## Project layout

```
cmd/todo/
├── main.go                # Entry point (CLI dispatcher + TUI launcher)
└── commands/
    ├── add.go             # CLI: add command
    ├── list.go            # CLI: list command
    ├── pick.go            # CLI: pick command
    ├── defer.go           # CLI: defer command
    ├── complete.go        # CLI: complete command
    ├── remove.go          # CLI: remove command
    ├── report.go          # CLI: report command
    └── sync.go            # CLI: sync command

internal/
├── model/                 # Domain model (core business logic)
│   ├── task.go            # Task entity + behavior
│   ├── task_test.go
│   ├── task_collection.go
│   ├── task_status.go
│   ├── task_id.go
│   ├── task_filter.go
│   └── errors.go
│
├── service/               # Application services
│   ├── task_service.go    # Orchestrates task operations
│   ├── task_service_test.go
│   ├── report_service.go
│   └── report_service_test.go
│
├── storage/               # Storage adapters
│   ├── storage.go         # Storage interface
│   ├── json_storage.go    # Encrypted JSON implementation
│   ├── json_storage_test.go
│   └── memory_storage.go  # In-memory (testing)
│
├── encryption/            # Encryption adapters
│   ├── encryptor.go       # Encryptor interface
│   ├── age_encryptor.go   # Age implementation
│   └── age_encryptor_test.go
│
└── sync/                  # Sync adapters
    ├── sync.go            # SyncAdapter interface
    ├── git_sync.go        # Git-backed sync
    └── git_sync_test.go

pkg/tui/                   # TUI implementation (Bubble Tea)
├── model.go               # Main Bubble Tea model
├── model_test.go
├── update.go              # Update logic (message handling)
├── view.go                # View rendering
├── commands.go            # Tea commands (async operations)
├── styles.go              # Lip Gloss styles
├── delegate.go            # List item delegate
├── keybindings.go         # Key binding definitions
└── messages.go            # Custom message types

scripts/
└── cron-report.sh         # Example cron script

README.md
ARCHITECTURE.md
IMPLEMENTATION_PLAN.md
CHARM_RESEARCH.md
```

## Milestones & estimated timeline

### Phase 1: Core Domain & Storage (Week 1)
**Focus**: Foundation and business logic

- Scaffold project structure
- Implement `Task` entity with state transitions (pick/defer/complete)
- Implement `TaskCollection` aggregate
- Implement value objects (TaskID, TaskStatus, TaskFilter)
- Implement `TaskService` orchestration layer
- Implement encrypted JSON storage with Age
- Unit tests for domain model (85%+ coverage)
- Integration tests for storage

**Deliverables**: Core domain model tested and working, encrypted persistence

### Phase 2: Interactive TUI (Week 2)
**Focus**: Beautiful, functional TUI

- Implement Bubble Tea main model (Model-Update-View)
- Implement tab-based navigation (Pool/Today/Done)
- Implement task list view with bubbles/list
  - Custom delegate for styling
  - Status-based colors and icons
  - Filtering support
- Implement add task view (bubbles/textinput)
- Implement task operations (pick/defer/complete/remove)
- Implement Lip Gloss styling system
  - Define color palette
  - Create reusable styles
  - Responsive layouts
- Implement help system (bubbles/help)
- Unit tests for TUI Update logic

**Deliverables**: Fully functional TUI with all core operations

### Phase 3: Sync & CLI Commands (Week 3)
**Focus**: Sync infrastructure and optional CLI

- Implement git-backed sync adapter
  - Push/pull operations
  - Status checking
  - Basic conflict detection
- Implement CLI dispatcher (flag-based)
  - Command routing (no args = launch TUI)
  - Implement all CLI commands (add/list/pick/defer/complete/remove)
  - Help text generation
- Add spinner for loading states (bubbles/spinner)
- Sync integration in TUI (s key)
- Integration tests for sync operations
- Lip Gloss styling for CLI output

**Deliverables**: Working sync, both TUI and CLI functional

### Phase 4: Polish & Advanced Features (Week 4)
**Focus**: Reports, refinements, documentation

- Implement report generation
  - ReportService with date range filtering
  - Text and JSON formatters
  - TUI report view
  - CLI report command
- Implement task detail view in TUI
  - Notes editing (bubbles/textarea)
  - Metadata display
- Implement task filtering in TUI (/ key)
- Add confirmation dialogs for destructive operations
- Cron script for weekly reports
- Documentation updates
- Performance testing and optimization
- Final polish on visual design

## Edge cases & risks

- Sync conflicts are complex — start simple (last-write-wins) and improve with interactive resolution later.
- Encryption key management: losing keys/passphrase means data loss — document backup and recovery thoroughly.
- Merge strategy: use stable `ID`-based merging to avoid duplicated tasks.
- Offline-first behavior: ensure sync operations are explicit and warned when remote is unreachable.

## Acceptance criteria (per feature)

- add: task persisted and appears in `list --pool`
- pick: status changes to `today` and is persisted
- defer: status changes to `pool` and `DeferredCount` increments
- complete: `CompletedAt` set and `status==done` stored
- remove: task deleted with confirmation
- report: `todo report --since 7d` prints tasks completed in the last 7 days
- sync: `todo sync pull` and `todo sync push` transfer the encrypted blob reliably between local copy and remote

## Implementation Notes

### TUI Architecture Patterns

**Model-Update-View Pattern** (Bubble Tea core):
- **Model**: Application state (tasks, current tab, selected item, etc.)
- **Update**: Message handler that modifies state and returns commands
- **View**: Pure function that renders state to string

**State Management**:
- Use view state enum (listView, addView, detailView)
- Use tab enum (poolTab, todayTab, doneTab)
- Keep all tasks in memory, filter by current tab
- Refresh list on tab switch

**Custom Messages**:
```go
type tasksLoadedMsg struct { tasks []*model.Task }
type taskOperationMsg struct { message string; err error }
type syncCompleteMsg struct { status string; err error }
```

**Commands** (async operations):
- Always return `tea.Cmd` for side effects
- Handle errors in custom messages
- Show spinner during long operations
- Update status bar on completion

**Component Integration**:
- bubbles/list handles its own Update/View
- Delegate list.Update to component when in list view
- Always call Init() when switching views
- Propagate WindowSizeMsg to all components

### Lip Gloss Styling Strategy

**Color Palette**:
- Primary: Pink (#d7005f / color 205)
- Secondary: Purple (#875fff / color 99)
- Pool: Blue (#5f87ff / color 63)
- Today: Pink (#ff87d7 / color 212)
- Done: Green (#00d75f / color 42)
- Muted: Gray (color 241)
- Error: Red (color 196)

**Style Composition**:
```go
containerStyle.Render(
  lipgloss.JoinVertical(
    lipgloss.Left,
    header,
    tabs,
    content,
    statusBar,
    help,
  )
)
```

**Responsive Design**:
- Use WindowSizeMsg to track terminal dimensions
- Adapt layout for narrow terminals (<80 cols)
- Cap max width for readability (100 cols)

### CLI Command Dispatcher Pattern

**Main Entry Point**:
```go
func main() {
  if len(os.Args) == 1 {
    launchTUI()  // No args = launch TUI
    return
  }

  command := os.Args[1]
  switch command {
  case "add": handleAdd(os.Args[2:])
  case "list": handleList(os.Args[2:])
  // ... etc
  case "help": printHelp()
  default: fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
  }
}
```

**Flag Parsing Per Command**:
```go
func handleAdd(args []string) {
  fs := flag.NewFlagSet("add", flag.ExitOnError)
  notes := fs.String("n", "", "Task notes")
  tags := fs.String("t", "", "Comma-separated tags")
  fs.Parse(args)

  // ... implementation
}
```

### Testing Strategy

**TUI Tests**:
- Test Update logic by sending messages
- Test View rendering by inspecting output strings
- Mock TaskService for isolation
- Test keybindings and state transitions

**Integration Tests**:
- Use temp directories for storage tests
- Test full TUI workflow (add → pick → complete)
- Test sync with temp git repo
- Verify encryption/decryption round-trip

**Table-Driven Tests**:
- Task state transitions (pick/defer/complete)
- Filter matching logic
- Command parsing and validation
