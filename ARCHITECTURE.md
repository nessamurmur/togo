# togo - Architecture & Design Document

**Version:** 1.0
**Last Updated:** 2025-11-09
**Status:** Active Development

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Architectural Principles](#architectural-principles)
3. [System Architecture](#system-architecture)
4. [Domain Model](#domain-model)
5. [Package Structure](#package-structure)
6. [Core Abstractions](#core-abstractions)
7. [Design Patterns](#design-patterns)
8. [Testing Strategy](#testing-strategy)
9. [Error Handling](#error-handling)
10. [Security & Encryption](#security--encryption)
11. [Concurrency Patterns](#concurrency-patterns)
12. [Go Best Practices](#go-best-practices)
13. [Code Organization Guidelines](#code-organization-guidelines)
14. [Integration Points](#integration-points)
15. [Quality Standards](#quality-standards)

---

## Executive Summary

**togo** is a command-line and terminal user interface (TUI) application for task management inspired by bullet journaling. The system is designed around clean architecture principles with a clear separation between domain logic, infrastructure concerns, and presentation layers.

### Key Design Goals

- **UI-Agnostic Core**: Business logic independent of UI implementation
- **Security First**: End-to-end encryption for all persisted data
- **Testability**: Comprehensive test coverage with clear testing boundaries
- **Maintainability**: SOLID principles and idiomatic Go patterns
- **Extensibility**: Easy to add new storage backends, sync adapters, and UI modes

### Technology Stack

- **Primary UI**: Bubble Tea (TUI framework), Bubbles (components), Lip Gloss (styling)
- **CLI Parsing**: Go standard library `flag` package (lightweight, no external dependencies)
- **Encryption**: filippo.io/age with scrypt/argon2id for passphrase mode
- **Storage**: Encrypted JSON (MVP), with optional bbolt migration path
- **Sync**: Git-backed remote (MVP), extensible to cloud services
- **Testing**: Standard library testing with table-driven tests

---

## Architectural Principles

### 1. SOLID Principles

#### Single Responsibility Principle (SRP)
Each package, type, and function has one well-defined reason to change:
- **Domain model** (`internal/model`) only changes when business rules change
- **Storage** (`internal/storage`) only changes when persistence requirements change
- **Sync** (`internal/sync`) only changes when synchronization logic changes
- **UI** (`pkg/tui`, `cmd/todo`) only changes when user interaction requirements change

#### Open/Closed Principle (OCP)
Design for extension without modification:
- Use **interfaces** for storage, sync, and encryption adapters
- New storage backends can be added without modifying existing code
- New sync providers can be plugged in through the sync adapter interface
- Report formatters can be added without changing report generation logic

#### Liskov Substitution Principle (LSP)
All implementations of an interface must be truly substitutable:
- Any `Storage` implementation must honor the same contracts and guarantees
- Any `SyncAdapter` must provide consistent push/pull semantics
- Any `Encryptor` must produce decryptable output

#### Interface Segregation Principle (ISP)
Create focused, client-specific interfaces:
- `TaskRepository` interface for task CRUD operations
- `Encryptor` interface for encryption/decryption only
- `SyncAdapter` interface separate from storage concerns
- `Reporter` interface for generating reports

#### Dependency Inversion Principle (DIP)
Depend on abstractions, not concretions:
- High-level domain logic depends on `Storage` interface, not JSON implementation
- CLI and TUI depend on task service abstractions, not concrete implementations
- Sync operations depend on abstract storage and encryptor interfaces

### 2. Domain-Driven Design (DDD)

#### Ubiquitous Language
Core terminology used consistently across code, tests, and documentation:
- **Task**: A unit of work with lifecycle states
- **Pool**: Collection of tasks not yet scheduled for today
- **Today**: Tasks picked for current focus
- **Defer**: Move a task back to the pool (intentional postponement)
- **Complete**: Mark a task as done
- **Pick**: Select a task from pool for today's focus
- **Report**: Aggregate view of completed tasks over time

#### Bounded Contexts

1. **Task Management Context** (`internal/model`)
   - Core domain: Task lifecycle, state transitions, business rules
   - Entities: Task
   - Value Objects: TaskID, TaskStatus, TaskFilter
   - Domain Services: TaskService (orchestrates operations)

2. **Storage Context** (`internal/storage`)
   - Responsibility: Persistence and retrieval of encrypted task data
   - Aggregates: TaskCollection
   - Infrastructure: JSON serialization, encryption, file I/O

3. **Sync Context** (`internal/sync`)
   - Responsibility: Synchronization with remote storage
   - Services: GitSyncAdapter
   - Operations: Push, Pull, Status, Merge

4. **Presentation Context** (`pkg/tui`, `cmd/todo`)
   - Responsibility: User interaction through TUI (primary) and CLI (optional)
   - No domain logic, only presentation and input handling
   - TUI: Bubble Tea Model-Update-View pattern
   - CLI: Simple command dispatcher with flag parsing

#### Tactical Patterns

- **Entities**: Task (has identity via UUID, mutable state)
- **Value Objects**: TaskID, TaskStatus (immutable, equality by value)
- **Aggregates**: TaskCollection (consistency boundary for task operations)
- **Repositories**: TaskRepository interface (abstract persistence)
- **Domain Services**: TaskService (coordinates task operations, enforces invariants)
- **Factory**: NewTask, NewTaskFromJSON (controlled construction)

### 3. Hexagonal Architecture (Ports & Adapters)

```
┌─────────────────────────────────────────────────────┐
│                   Presentation                       │
│  ┌──────────────┐              ┌──────────────┐    │
│  │  CLI (flag)  │              │  TUI (Bubble │    │
│  │  Optional    │              │  Tea) Primary│    │
│  └──────┬───────┘              └──────┬────────┘    │
└─────────┼──────────────────────────────┼────────────┘
          │         Inbound Ports        │
          ▼                              ▼
┌─────────────────────────────────────────────────────┐
│              Application Services                    │
│  ┌───────────────────────────────────────────────┐ │
│  │           TaskService (Orchestration)         │ │
│  │  - AddTask, PickTask, DeferTask, Complete     │ │
│  │  - ListTasks, RemoveTask, GenerateReport      │ │
│  └───────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
          │                              │
          ▼                              ▼
┌─────────────────────────────────────────────────────┐
│                  Domain Model                        │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────┐ │
│  │    Task     │  │  TaskStatus  │  │  TaskID    │ │
│  │  (Entity)   │  │  (Value Obj) │  │ (Value Obj)│ │
│  └─────────────┘  └──────────────┘  └────────────┘ │
└─────────────────────────────────────────────────────┘
          │         Outbound Ports       │
          ▼                              ▼
┌─────────────────────────────────────────────────────┐
│                    Adapters                          │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────┐ │
│  │Storage (JSON)│  │Sync (Git)    │  │Encryptor  │ │
│  │Adapter       │  │Adapter       │  │(Age)      │ │
│  └──────────────┘  └──────────────┘  └───────────┘ │
└─────────────────────────────────────────────────────┘
```

**Benefits:**
- Domain logic isolated from infrastructure
- Easy to swap implementations (e.g., JSON → bbolt, Git → S3)
- Testable with mock adapters
- Clear dependency flow (inward dependencies only)

---

## System Architecture

### High-Level Component Diagram

```
┌────────────────────────────────────────────────────────────┐
│                         User Layer                         │
│                                                            │
│   ┌──────────────┐                   ┌──────────────┐     │
│   │  CLI Mode    │                   │   TUI Mode   │     │
│   │  (optional)  │                   │  (primary)   │     │
│   │   flag pkg   │                   │ (Bubble Tea) │     │
│   └──────┬───────┘                   └──────┬───────┘     │
│          │                                  │             │
└──────────┼──────────────────────────────────┼─────────────┘
           │                                  │
           └─────────┬────────────────────────┘
                     │
           ┌─────────▼───────────┐
           │  main.go Dispatcher │
           │  (no args = TUI)    │
           └─────────┬───────────┘
                     │
           ┌─────────▼───────────┐
           │   TaskService       │ ◄────── Application Service Layer
           │  (orchestration)    │
           └─────────┬───────────┘
                     │
      ┌──────────────┼──────────────┐
      │              │              │
      ▼              ▼              ▼
┌──────────┐  ┌──────────┐  ┌──────────┐
│  Task    │  │TaskRepo  │  │Reporter  │ ◄── Domain Layer
│ (Entity) │  │Interface │  │Interface │
└──────────┘  └────┬─────┘  └────┬─────┘
                   │             │
      ┌────────────┼─────────────┘
      │            │
      ▼            ▼
┌──────────┐  ┌──────────┐  ┌──────────┐
│JSONStore │  │Encryptor │  │GitSync   │ ◄── Infrastructure Layer
│ (Adapter)│  │ (Age)    │  │ (Adapter)│
└──────────┘  └──────────┘  └──────────┘
      │            │              │
      └────────────┼──────────────┘
                   │
                   ▼
          ┌────────────────┐
          │  File System   │
          │  (encrypted)   │
          └────────────────┘
```

### Data Flow

#### Adding a Task (CLI - Optional Mode)
1. User executes: `todo add "Buy groceries" -t shopping`
2. main.go dispatcher parses command using flag package
3. handleAdd() creates AddTaskRequest
4. TaskService.AddTask() called with title and tags
5. New Task entity created with UUID, timestamp, status="pool"
6. TaskRepository.Save() persists task collection
7. JSONStorage encrypts payload with Age
8. Encrypted file written to `$XDG_DATA_HOME/togo/todo.json.age`
9. Success message printed with Lip Gloss styling

#### Picking a Task (TUI - Primary Mode)
1. User launches `todo` (no args) → TUI starts
2. User navigates list with j/k, presses Enter on a pool task
3. Bubble Tea Update() processes tea.KeyMsg
4. Update() returns pickTaskCmd (tea.Cmd)
5. Command executes asynchronously: TaskService.PickTask(taskID)
6. Domain model validates state transition (pool → today)
7. Task.Status updated, validation passes
8. TaskRepository.Save() called
9. Command returns taskPickedMsg
10. Update() handles taskPickedMsg, refreshes task list
11. View() re-renders with updated task list and status message

#### Syncing Tasks (TUI)
1. User presses 's' key in TUI
2. Update() returns syncCmd with spinner
3. Spinner starts animating
4. SyncAdapter.Push() executes asynchronously
5. Reads encrypted local file
6. Git adapter: `git add todo.json.age && git commit && git push`
7. syncCompleteMsg sent back to Update()
8. Spinner stops, status bar shows "Sync complete" or error
9. View() re-renders with updated status

---

## Domain Model

### Task Entity

```go
// Task represents a unit of work in the bullet journal system.
// It is the core entity with identity (ID) and lifecycle management.
type Task struct {
    // Immutable identity
    ID          TaskID        `json:"id"`
    CreatedAt   time.Time     `json:"created_at"`

    // Mutable attributes
    Title       string        `json:"title"`
    Notes       string        `json:"notes,omitempty"`
    Status      TaskStatus    `json:"status"`
    Tags        []string      `json:"tags,omitempty"`

    // Optional temporal data
    DueDate     *time.Time    `json:"due_date,omitempty"`
    CompletedAt *time.Time    `json:"completed_at,omitempty"`

    // Metadata
    DeferredCount int         `json:"deferred_count"`
}
```

**Invariants:**
- ID must be a valid UUID
- Status must be one of: pool, today, done
- CompletedAt must be set when Status == done
- CompletedAt must be nil when Status != done
- DeferredCount >= 0
- Title must not be empty

**Behavior Methods:**
```go
func (t *Task) Pick() error
func (t *Task) Defer() error
func (t *Task) Complete() error
func (t *Task) Validate() error
func (t *Task) MatchesFilter(filter TaskFilter) bool
```

### TaskStatus Value Object

```go
type TaskStatus string

const (
    StatusPool  TaskStatus = "pool"
    StatusToday TaskStatus = "today"
    StatusDone  TaskStatus = "done"
)

func (s TaskStatus) Valid() bool
func (s TaskStatus) String() string
```

**Rationale:** Type-safe enumeration prevents invalid status values. Immutable by design.

### TaskID Value Object

```go
type TaskID string

func NewTaskID() TaskID
func ParseTaskID(s string) (TaskID, error)
func (id TaskID) Valid() bool
func (id TaskID) String() string
```

**Rationale:** Encapsulates UUID generation and validation. Prevents invalid IDs from being created.

### TaskFilter Value Object

```go
type TaskFilter struct {
    Status    *TaskStatus
    Tags      []string
    DueAfter  *time.Time
    DueBefore *time.Time
    Limit     int
}

func (f TaskFilter) Matches(t *Task) bool
```

**Rationale:** Encapsulates filtering logic, reusable across CLI and TUI.

### TaskCollection Aggregate

```go
type TaskCollection struct {
    tasks map[TaskID]*Task
    meta  Metadata
}

type Metadata struct {
    Version        string
    LastModified   time.Time
    EncryptionMode string
    Salt           string
}

func (tc *TaskCollection) Add(t *Task) error
func (tc *TaskCollection) Get(id TaskID) (*Task, error)
func (tc *TaskCollection) Remove(id TaskID) error
func (tc *TaskCollection) Find(filter TaskFilter) []*Task
func (tc *TaskCollection) All() []*Task
```

**Rationale:** Provides consistency boundary for task operations. Ensures no duplicate IDs, enforces collection-level invariants.

---

## Package Structure

### Recommended Layout

```
togo/
├── cmd/
│   └── todo/
│       ├── main.go              # Entry point & dispatcher
│       └── commands/
│           ├── add.go           # CLI: add command handler
│           ├── list.go          # CLI: list command handler
│           ├── pick.go          # CLI: pick command handler
│           ├── defer.go         # CLI: defer command handler
│           ├── complete.go      # CLI: complete command handler
│           ├── remove.go        # CLI: remove command handler
│           ├── report.go        # CLI: report command handler
│           ├── sync.go          # CLI: sync command handler
│           └── help.go          # CLI: help text
│
├── internal/
│   ├── model/                   # Domain model (core business logic)
│   │   ├── task.go              # Task entity + behavior
│   │   ├── task_test.go         # Task unit tests
│   │   ├── task_collection.go   # Aggregate root
│   │   ├── task_collection_test.go
│   │   ├── task_status.go       # Value object
│   │   ├── task_id.go           # Value object
│   │   ├── task_filter.go       # Value object
│   │   └── errors.go            # Domain errors
│   │
│   ├── service/                 # Application services
│   │   ├── task_service.go      # Orchestrates task operations
│   │   ├── task_service_test.go
│   │   ├── report_service.go    # Generates reports
│   │   └── report_service_test.go
│   │
│   ├── storage/                 # Storage adapters
│   │   ├── storage.go           # Storage interface (port)
│   │   ├── json_storage.go      # Encrypted JSON implementation
│   │   ├── json_storage_test.go
│   │   └── memory_storage.go    # In-memory (for testing)
│   │
│   ├── encryption/              # Encryption adapters
│   │   ├── encryptor.go         # Encryptor interface (port)
│   │   ├── age_encryptor.go     # Age implementation
│   │   ├── age_encryptor_test.go
│   │   └── noop_encryptor.go    # No-op (for testing)
│   │
│   ├── sync/                    # Sync adapters
│   │   ├── sync.go              # SyncAdapter interface (port)
│   │   ├── git_sync.go          # Git-backed sync
│   │   ├── git_sync_test.go
│   │   └── noop_sync.go         # No-op sync (for testing)
│   │
│   └── config/                  # Configuration management
│       ├── config.go            # Config struct and loading
│       ├── config_test.go
│       └── defaults.go          # Default values
│
├── pkg/                         # Public packages (TUI)
│   └── tui/
│       ├── model.go             # Main Bubble Tea model
│       ├── model_test.go        # Model tests
│       ├── update.go            # Update logic (message handling)
│       ├── view.go              # View rendering
│       ├── commands.go          # Tea commands (async operations)
│       ├── styles.go            # Lip Gloss styles
│       ├── delegate.go          # List item delegate
│       ├── keybindings.go       # Key binding definitions
│       └── messages.go          # Custom message types
│
├── pkg/reporter/                # Report formatting (optional)
│   ├── text_formatter.go
│   ├── json_formatter.go
│   └── formatter.go             # Formatter interface
│
├── scripts/
│   └── cron-report.sh           # Example cron script
│
├── go.mod
├── go.sum
├── README.md
├── ARCHITECTURE.md              # This document
└── IMPLEMENTATION_PLAN.md
```

### Package Dependency Rules

**Allowed Dependencies (inward only):**
```
cmd/todo → internal/service → internal/model
                            → internal/storage
                            → internal/sync
                            → internal/encryption
pkg/tui  → internal/service → internal/model

internal/storage → internal/model
                 → internal/encryption

internal/sync → internal/storage
              → internal/encryption

internal/encryption → (external libs only)
internal/model → (stdlib only)
```

**Forbidden Dependencies:**
- `internal/model` must NOT import `internal/storage`, `internal/sync`, or `cmd`
- `internal/storage` must NOT import `internal/sync` or `cmd`
- No circular dependencies between packages

**Rationale:** Enforces hexagonal architecture with inward-pointing dependencies. Domain model remains pure and testable.

---

## Core Abstractions

### 1. TaskRepository Interface (Port)

```go
// TaskRepository defines the contract for task persistence.
// This is the primary port for storage adapters.
type TaskRepository interface {
    // Load reads all tasks from storage.
    // Returns empty collection if no data exists (not an error).
    Load(ctx context.Context) (*TaskCollection, error)

    // Save persists the entire task collection.
    // Must be atomic: either all tasks save or none.
    Save(ctx context.Context, collection *TaskCollection) error

    // Close releases any resources held by the repository.
    Close() error
}
```

**Implementation Guarantees:**
- Load must be idempotent
- Save must be atomic (all-or-nothing)
- Concurrent calls to Save are safe (last write wins)
- Decryption errors are propagated to caller

### 2. Encryptor Interface (Port)

```go
// Encryptor handles encryption and decryption of data at rest.
type Encryptor interface {
    // Encrypt takes plaintext and returns ciphertext.
    Encrypt(plaintext []byte) ([]byte, error)

    // Decrypt takes ciphertext and returns plaintext.
    // Returns error if decryption fails (wrong key, corrupted data).
    Decrypt(ciphertext []byte) ([]byte, error)
}
```

**Implementation Requirements:**
- Must use authenticated encryption (AEAD)
- Must handle key derivation (for passphrase mode)
- Must be deterministic for same input + key
- Must fail safely on tampering

### 3. SyncAdapter Interface (Port)

```go
// SyncAdapter defines the contract for remote synchronization.
type SyncAdapter interface {
    // Status checks if local and remote are in sync.
    // Returns SyncStatus indicating ahead/behind/synced/diverged.
    Status(ctx context.Context) (SyncStatus, error)

    // Pull fetches remote data and merges with local.
    // Returns conflict if merge cannot be automatic.
    Pull(ctx context.Context) (*PullResult, error)

    // Push uploads local data to remote.
    // Returns error if remote has diverged.
    Push(ctx context.Context) error
}

type SyncStatus int

const (
    SyncStatusSynced SyncStatus = iota
    SyncStatusAhead
    SyncStatusBehind
    SyncStatusDiverged
)

type PullResult struct {
    Conflict bool
    Message  string
}
```

**Implementation Requirements:**
- Must handle network failures gracefully
- Must detect conflicts and report them
- Must not lose data on conflict (preserve both versions)
- Must work with encrypted blobs (no plaintext on remote)

### 4. TaskService Interface (Application Service)

```go
// TaskService orchestrates task operations and enforces business rules.
// This is the primary interface for CLI and TUI layers.
type TaskService interface {
    // AddTask creates a new task in the pool.
    AddTask(ctx context.Context, req AddTaskRequest) (*Task, error)

    // ListTasks retrieves tasks matching filter.
    ListTasks(ctx context.Context, filter TaskFilter) ([]*Task, error)

    // PickTask moves a task from pool to today.
    PickTask(ctx context.Context, id TaskID) error

    // DeferTask moves a task from today back to pool.
    DeferTask(ctx context.Context, id TaskID) error

    // CompleteTask marks a task as done.
    CompleteTask(ctx context.Context, id TaskID) error

    // RemoveTask deletes a task permanently.
    RemoveTask(ctx context.Context, id TaskID) error

    // GenerateReport produces a report of completed tasks.
    GenerateReport(ctx context.Context, opts ReportOptions) (*Report, error)
}

type AddTaskRequest struct {
    Title   string
    Notes   string
    Tags    []string
    DueDate *time.Time
}

type ReportOptions struct {
    Since  time.Time
    Until  time.Time
    Format string // "text" or "json"
}
```

**Rationale:** Service layer decouples UI from domain model, handles transaction boundaries, coordinates multiple operations.

---

## Design Patterns

### 1. Repository Pattern

**Problem:** Domain model needs persistence without coupling to storage technology.

**Solution:** Define `TaskRepository` interface in domain layer, implement in infrastructure layer.

**Example:**
```go
// Domain layer defines the interface
type TaskRepository interface {
    Load(ctx context.Context) (*TaskCollection, error)
    Save(ctx context.Context, collection *TaskCollection) error
}

// Infrastructure layer provides implementation
type JSONStorage struct {
    path      string
    encryptor Encryptor
}

func (s *JSONStorage) Load(ctx context.Context) (*TaskCollection, error) {
    // Implementation details
}
```

### 2. Adapter Pattern

**Problem:** Multiple storage backends, encryption methods, and sync providers.

**Solution:** Define ports (interfaces), implement adapters for each technology.

**Example:**
```go
// Port
type Encryptor interface {
    Encrypt([]byte) ([]byte, error)
    Decrypt([]byte) ([]byte, error)
}

// Adapter 1: Age encryption
type AgeEncryptor struct { /* ... */ }

// Adapter 2: No-op (testing)
type NoopEncryptor struct { /* ... */ }
```

### 3. Factory Pattern

**Problem:** Complex object construction with invariants.

**Solution:** Provide factory functions that enforce validation.

**Example:**
```go
// NewTask creates a task with required fields validated.
func NewTask(title string, tags []string) (*Task, error) {
    if strings.TrimSpace(title) == "" {
        return nil, ErrEmptyTitle
    }

    return &Task{
        ID:        NewTaskID(),
        Title:     title,
        Tags:      tags,
        Status:    StatusPool,
        CreatedAt: time.Now(),
    }, nil
}
```

### 4. Strategy Pattern

**Problem:** Multiple report formats (text, JSON) without coupling.

**Solution:** Define `Formatter` interface, swap implementations based on user flag.

**Example:**
```go
type Formatter interface {
    Format(report *Report) (string, error)
}

type TextFormatter struct{}
type JSONFormatter struct{}

// Usage
var formatter Formatter
if opts.Format == "json" {
    formatter = &JSONFormatter{}
} else {
    formatter = &TextFormatter{}
}
output := formatter.Format(report)
```

### 5. Command Pattern (Bubble Tea)

**Problem:** User actions in TUI need to trigger async operations.

**Solution:** Use Bubble Tea's `tea.Cmd` pattern for side effects.

**Example:**
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            // Return command to perform async operation
            return m, m.pickTaskCmd(m.selectedTaskID)
        }
    case taskPickedMsg:
        // Handle result of async operation
        m.status = "Task picked successfully"
        return m, nil
    }
    return m, nil
}

func (m model) pickTaskCmd(id TaskID) tea.Cmd {
    return func() tea.Msg {
        err := m.service.PickTask(context.Background(), id)
        return taskPickedMsg{err: err}
    }
}
```

### 6. Table-Driven Tests

**Problem:** Testing multiple scenarios with similar structure.

**Solution:** Use Go's table-driven test pattern.

**Example:**
```go
func TestTask_Pick(t *testing.T) {
    tests := []struct {
        name      string
        initial   TaskStatus
        wantErr   bool
        wantStatus TaskStatus
    }{
        {
            name:       "pick from pool succeeds",
            initial:    StatusPool,
            wantErr:    false,
            wantStatus: StatusToday,
        },
        {
            name:       "pick from today is idempotent",
            initial:    StatusToday,
            wantErr:    false,
            wantStatus: StatusToday,
        },
        {
            name:       "pick from done fails",
            initial:    StatusDone,
            wantErr:    true,
            wantStatus: StatusDone,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            task := &Task{Status: tt.initial}
            err := task.Pick()

            if (err != nil) != tt.wantErr {
                t.Errorf("Pick() error = %v, wantErr %v", err, tt.wantErr)
            }

            if task.Status != tt.wantStatus {
                t.Errorf("Status = %v, want %v", task.Status, tt.wantStatus)
            }
        })
    }
}
```

---

## Bubble Tea Architecture

### Overview

Bubble Tea is the primary UI framework for togo, providing a beautiful, interactive terminal experience. It follows the Elm Architecture (Model-Update-View) pattern, ensuring predictable state management and easy testability.

### Core Pattern: Model-Update-View

```go
// Model: Application state
type Model struct {
    // State
    currentView viewState      // Which view is active
    activeTab   tab            // Pool/Today/Done
    allTasks    []*model.Task  // All loaded tasks

    // Components (Bubbles)
    list     list.Model        // Task list component
    addInput textinput.Model   // Add task input
    spinner  spinner.Model     // Loading spinner
    help     help.Model        // Help component

    // Dependencies
    service  *service.TaskService

    // UI State
    width    int
    height   int
    statusMsg string
}

// Update: Message handler (pure function of state)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle keyboard input
        return m.handleKeyPress(msg)
    case tasksLoadedMsg:
        // Handle async operation result
        m.allTasks = msg.tasks
        m.refreshList()
        return m, nil
    case tea.WindowSizeMsg:
        // Handle terminal resize
        m.width = msg.Width
        m.height = msg.Height
        return m, nil
    }
    return m, nil
}

// View: Render state to string (pure function)
func (m Model) View() string {
    switch m.currentView {
    case listView:
        return m.viewList()
    case addView:
        return m.viewAdd()
    case detailView:
        return m.viewDetail()
    }
    return ""
}
```

### State Management

**View States** (viewState enum):
- `listView`: Main task list with tabs (primary view)
- `addView`: Add new task modal
- `detailView`: View/edit task details

**Tab States** (tab enum):
- `poolTab`: Unscheduled tasks
- `todayTab`: Tasks picked for today
- `doneTab`: Completed tasks

**State Transitions**:
```go
// Switching views
m.currentView = addView
return m, textinput.Blink  // Must call Init() equivalent

// Switching tabs
m.activeTab = (m.activeTab + 1) % 3
m.refreshList()  // Filter tasks for new tab
```

### Custom Messages

Define custom message types for async operations:

```go
// Data loading
type tasksLoadedMsg struct {
    tasks []*model.Task
}

// Task operations
type taskOperationMsg struct {
    operation string  // "picked", "deferred", "completed", etc.
    err       error
}

// Sync operations
type syncCompleteMsg struct {
    status string  // "pushed", "pulled", "synced"
    err    error
}

// Spinner ticks
type spinner.TickMsg struct{}
```

### Commands (Async Operations)

Commands are functions that return messages. They execute side effects outside the main event loop:

```go
// Load tasks from storage
func loadTasks(svc *service.TaskService) tea.Cmd {
    return func() tea.Msg {
        ctx := context.Background()
        tasks, err := svc.ListTasks(ctx, model.TaskFilter{})
        if err != nil {
            return taskOperationMsg{err: err}
        }
        return tasksLoadedMsg{tasks: tasks}
    }
}

// Pick a task
func (m Model) pickTaskCmd(taskID model.TaskID) tea.Cmd {
    return func() tea.Msg {
        err := m.service.PickTask(context.Background(), taskID)
        return taskOperationMsg{
            operation: "picked",
            err:       err,
        }
    }
}

// Batch multiple commands
return m, tea.Batch(
    m.spinner.Tick,
    loadTasks(m.service),
)
```

### Component Integration

**Bubbles Components Used:**

1. **list.Model** - Main task list
   ```go
   // Create with custom delegate
   delegate := newTaskDelegate()
   l := list.New([]list.Item{}, delegate, width, height)
   l.Title = "Pool"
   l.SetFilteringEnabled(true)

   // Update
   m.list, cmd = m.list.Update(msg)

   // Render
   listView := m.list.View()
   ```

2. **textinput.Model** - Single-line inputs
   ```go
   ti := textinput.New()
   ti.Placeholder = "Task title..."
   ti.Focus()
   ti.CharLimit = 200

   // In Update
   m.addInput, cmd = m.addInput.Update(msg)
   ```

3. **textarea.Model** - Multi-line notes
   ```go
   ta := textarea.New()
   ta.Placeholder = "Task notes..."
   ta.SetValue(task.Notes)
   ta.Focus()
   ```

4. **spinner.Model** - Loading states
   ```go
   s := spinner.New()
   s.Spinner = spinner.Dot
   s.Style = lipgloss.NewStyle().Foreground(primaryColor)

   // Must tick in Update
   m.spinner, cmd = m.spinner.Update(msg)
   ```

5. **help.Model** - Keybinding help
   ```go
   h := help.New()
   helpView := h.View(m.keys)  // keys is keyMap
   ```

### Lip Gloss Styling

**Color Palette**:
```go
var (
    primaryColor   = lipgloss.Color("205") // Pink
    secondaryColor = lipgloss.Color("99")  // Purple
    poolColor      = lipgloss.Color("63")  // Blue
    todayColor     = lipgloss.Color("212") // Pink (lighter)
    doneColor      = lipgloss.Color("42")  // Green
    mutedColor     = lipgloss.Color("241") // Gray
    errorColor     = lipgloss.Color("196") // Red
)
```

**Style Definitions**:
```go
var (
    // Text styles
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(primaryColor).
        MarginBottom(1)

    // Layout styles
    containerStyle = lipgloss.NewStyle().
        Padding(1, 2)

    boxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(primaryColor).
        Padding(1, 2)

    // Tab styles
    tabStyle = lipgloss.NewStyle().
        Padding(0, 2).
        Background(lipgloss.Color("238")).
        Foreground(lipgloss.Color("248"))

    activeTabStyle = lipgloss.NewStyle().
        Padding(0, 2).
        Background(primaryColor).
        Foreground(lipgloss.Color("0")).
        Bold(true)

    // Status styles
    poolStyle = lipgloss.NewStyle().
        Foreground(poolColor)

    todayStyle = lipgloss.NewStyle().
        Foreground(todayColor).
        Bold(true)

    doneStyle = lipgloss.NewStyle().
        Foreground(doneColor).
        Strikethrough(true)
)
```

**Layout Composition**:
```go
func (m Model) viewList() string {
    // Build sections
    header := titleStyle.Render("togo")
    tabs := m.renderTabs()
    listView := m.list.View()
    statusBar := m.renderStatusBar()
    helpView := helpStyle.Render(m.help.View(m.keys))

    // Compose vertically
    content := lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        tabs,
        "",  // Spacer
        listView,
        "",
        statusBar,
        helpView,
    )

    // Apply container
    return containerStyle.Render(content)
}

func (m Model) renderTabs() string {
    tabs := []string{"Pool", "Today", "Done"}
    rendered := make([]string, len(tabs))

    for i, name := range tabs {
        style := tabStyle
        if i == int(m.activeTab) {
            style = activeTabStyle
        }
        rendered[i] = style.Render(name)
    }

    // Join horizontally
    return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}
```

**Responsive Design**:
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

        // Update component sizes
        m.list.SetSize(msg.Width-4, msg.Height-10)

        // Adapt layout for narrow terminals
        if m.width < 80 {
            m.compactMode = true
        }

        return m, nil
    }
}
```

### Key Bindings

**Define Key Map**:
```go
type keyMap struct {
    Up       key.Binding
    Down     key.Binding
    Tab      key.Binding
    Pick     key.Binding
    Complete key.Binding
    Defer    key.Binding
    Add      key.Binding
    Remove   key.Binding
    Sync     key.Binding
    Help     key.Binding
    Quit     key.Binding
}

// Implement help.KeyMap interface
func (k keyMap) ShortHelp() []key.Binding {
    return []key.Binding{k.Up, k.Down, k.Pick, k.Add, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {k.Up, k.Down, k.Tab, k.Pick},
        {k.Complete, k.Defer, k.Add, k.Remove},
        {k.Sync, k.Help, k.Quit},
    }
}

var keys = keyMap{
    Up: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("↑/k", "move up"),
    ),
    Pick: key.NewBinding(
        key.WithKeys("enter"),
        key.WithHelp("enter", "pick task"),
    ),
    // ... etc
}
```

**Handle Key Presses**:
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Global keys (work in all views)
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "?":
            m.help.ShowAll = !m.help.ShowAll
            return m, nil
        }

        // View-specific routing
        switch m.currentView {
        case listView:
            return m.handleListKeys(msg)
        case addView:
            return m.handleAddKeys(msg)
        }
    }
    return m, nil
}

func (m Model) handleListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "tab":
        m.activeTab = (m.activeTab + 1) % 3
        m.refreshList()
        return m, nil
    case "a":
        m.currentView = addView
        m.addInput.Focus()
        return m, textinput.Blink
    case "enter":
        return m, m.pickSelectedTask()
    // ... etc
    }
    return m, nil
}
```

### Multi-View Pattern

**Pattern: State-Based Routing**:
```go
type viewState int

const (
    listView viewState = iota
    addView
    detailView
)

type Model struct {
    currentView viewState
    // ... other fields
}

// Update routes to view-specific handlers
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Global handlers first
        // ...

        // Then route by view
        switch m.currentView {
        case listView:
            return m.updateList(msg)
        case addView:
            return m.updateAdd(msg)
        case detailView:
            return m.updateDetail(msg)
        }
    }
    return m, nil
}

// View routes to view-specific renderers
func (m Model) View() string {
    switch m.currentView {
    case listView:
        return m.viewList()
    case addView:
        return m.viewAdd()
    case detailView:
        return m.viewDetail()
    }
    return "Unknown view"
}
```

**Critical Rule: Always Call Init() When Switching**:
```go
// WRONG - component won't initialize
case "a":
    m.currentView = addView
    return m, nil

// CORRECT - component properly initialized
case "a":
    m.currentView = addView
    m.addInput.Focus()
    return m, textinput.Blink  // Init equivalent
```

### Tab-Based Navigation

**Pattern: Single List, Multiple Filters**:
```go
type tab int

const (
    poolTab tab = iota
    todayTab
    doneTab
)

func (m *Model) refreshList() {
    var filtered []*model.Task
    var title string

    // Filter by active tab
    switch m.activeTab {
    case poolTab:
        title = "Pool"
        for _, t := range m.allTasks {
            if t.Status == model.StatusPool {
                filtered = append(filtered, t)
            }
        }
    case todayTab:
        title = "Today"
        for _, t := range m.allTasks {
            if t.Status == model.StatusToday {
                filtered = append(filtered, t)
            }
        }
    case doneTab:
        title = "Done"
        for _, t := range m.allTasks {
            if t.Status == model.StatusDone {
                filtered = append(filtered, t)
            }
        }
    }

    // Convert to list items
    items := make([]list.Item, len(filtered))
    for i, task := range filtered {
        items[i] = taskItem{task: task}
    }

    m.list.SetItems(items)
    m.list.Title = title
}
```

### Custom List Delegate

**Implement list.ItemDelegate**:
```go
type taskDelegate struct{}

func (d taskDelegate) Height() int  { return 2 }
func (d taskDelegate) Spacing() int { return 1 }
func (d taskDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
    return nil
}

func (d taskDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
    i, ok := item.(taskItem)
    if !ok {
        return
    }

    task := i.task

    // Choose style based on status
    var icon string
    var style lipgloss.Style

    switch task.Status {
    case model.StatusPool:
        icon = "○"
        style = poolStyle
    case model.StatusToday:
        icon = "◉"
        style = todayStyle
    case model.StatusDone:
        icon = "✓"
        style = doneStyle
    }

    // Highlight if selected
    if index == m.Index() {
        style = selectedItemStyle
        icon = "▶"
    }

    // Render title and description
    title := style.Render(fmt.Sprintf("%s %s", icon, task.Title))
    desc := subtleStyle.Render(i.Description())

    fmt.Fprintf(w, "%s\n%s", title, desc)
}
```

**Implement list.Item**:
```go
type taskItem struct {
    task *model.Task
}

func (i taskItem) FilterValue() string {
    return i.task.Title
}

func (i taskItem) Title() string {
    return i.task.Title
}

func (i taskItem) Description() string {
    parts := []string{}

    if len(i.task.Tags) > 0 {
        parts = append(parts, fmt.Sprintf("Tags: %s", strings.Join(i.task.Tags, ", ")))
    }

    if i.task.DueDate != nil {
        parts = append(parts, fmt.Sprintf("Due: %s", i.task.DueDate.Format("2006-01-02")))
    }

    return strings.Join(parts, " | ")
}
```

### Testing Bubble Tea Applications

**Test Update Logic**:
```go
func TestModel_PickTask(t *testing.T) {
    // Setup
    mockService := &MockTaskService{}
    model := NewModel(mockService)
    model.allTasks = []*model.Task{
        {ID: "123", Title: "Test", Status: model.StatusPool},
    }
    model.refreshList()

    // Action: simulate Enter key
    msg := tea.KeyMsg{Type: tea.KeyEnter}
    newModel, cmd := model.Update(msg)

    // Assert: command returned
    assert.NotNil(t, cmd)

    // Execute command
    resultMsg := cmd()

    // Handle result
    finalModel, _ := newModel.(Model).Update(resultMsg)

    // Verify
    assert.Equal(t, 1, mockService.PickTaskCallCount)
}
```

**Test View Rendering**:
```go
func TestModel_ViewList(t *testing.T) {
    model := NewModel(nil)
    model.currentView = listView
    model.activeTab = poolTab

    view := model.View()

    // Assert expected content
    assert.Contains(t, view, "Pool")
    assert.Contains(t, view, "togo")
}
```

**Test Tab Switching**:
```go
func TestModel_TabSwitch(t *testing.T) {
    model := NewModel(nil)
    model.activeTab = poolTab

    // Simulate Tab key
    msg := tea.KeyMsg{Type: tea.KeyTab}
    newModel, _ := model.Update(msg)

    m := newModel.(Model)
    assert.Equal(t, todayTab, m.activeTab)
}
```

### Best Practices

**1. Separation of Concerns**:
- Keep business logic in TaskService
- TUI only handles presentation and input
- Use commands for all side effects

**2. State Immutability**:
- Update returns new model (Go's copy semantics)
- Don't mutate model in View
- View is pure function of state

**3. Error Handling**:
- Return errors in custom messages
- Display in status bar or error view
- Never panic in Update or View

**4. Performance**:
- View is called frequently - keep it fast
- Avoid heavy computation in View
- Pre-compute in Update, cache in model

**5. Responsiveness**:
- Use spinners for long operations
- Provide immediate feedback for key presses
- Update status bar to show progress

---

## Testing Strategy

### Testing Pyramid

```
        ┌──────────────┐
        │   E2E Tests  │  (Few, critical user journeys)
        │  Integration │
        └──────────────┘
       ┌────────────────┐
       │ Integration    │  (Adapter integration, TUI flows)
       │    Tests       │
       └────────────────┘
     ┌──────────────────────┐
     │    Unit Tests        │  (Most tests here, fast and isolated)
     └──────────────────────┘
```

### 1. Unit Tests

**Scope:** Individual functions and methods in isolation.

**Coverage Requirements:**
- 80%+ coverage for domain model (`internal/model`)
- 70%+ coverage for services (`internal/service`)
- 60%+ coverage for adapters (`internal/storage`, `internal/sync`)

**Testing Approach:**
```go
// Test domain logic in isolation
func TestTask_Defer(t *testing.T) {
    task := &Task{Status: StatusToday, DeferredCount: 0}

    err := task.Defer()

    assert.NoError(t, err)
    assert.Equal(t, StatusPool, task.Status)
    assert.Equal(t, 1, task.DeferredCount)
}

// Test with table-driven approach for multiple scenarios
func TestTaskCollection_Add(t *testing.T) {
    tests := []struct {
        name    string
        task    *Task
        wantErr bool
    }{
        {"valid task", &Task{ID: "123", Title: "Test"}, false},
        {"nil task", nil, true},
        {"empty title", &Task{ID: "456", Title: ""}, true},
    }
    // ... test execution
}
```

### 2. Integration Tests

**Scope:** Multiple components working together, real adapters with test fixtures.

**Test Cases:**
- Storage: Write encrypted file, read back, verify data integrity
- Sync: Push to temp git repo, pull from another instance, verify merge
- TUI: Send key sequences, verify model state changes

**Example:**
```go
func TestJSONStorage_RoundTrip(t *testing.T) {
    // Use temp directory for test
    tmpDir := t.TempDir()

    // Real encryptor with test passphrase
    encryptor, _ := NewAgeEncryptor("test-passphrase")

    // Real storage adapter
    storage := NewJSONStorage(filepath.Join(tmpDir, "test.json.age"), encryptor)

    // Create test data
    collection := &TaskCollection{
        tasks: map[TaskID]*Task{
            "1": {ID: "1", Title: "Test Task", Status: StatusPool},
        },
    }

    // Save
    err := storage.Save(context.Background(), collection)
    require.NoError(t, err)

    // Verify file exists and is encrypted
    data, _ := os.ReadFile(filepath.Join(tmpDir, "test.json.age"))
    assert.NotContains(t, string(data), "Test Task") // Should be encrypted

    // Load
    loaded, err := storage.Load(context.Background())
    require.NoError(t, err)

    // Verify data integrity
    assert.Equal(t, collection.tasks["1"].Title, loaded.tasks["1"].Title)
}
```

### 3. TUI Tests

**Approach:** Test Bubble Tea Update and View functions directly.

**Example:**
```go
func TestTUI_PickTask(t *testing.T) {
    // Setup: model with mock service
    mockService := &MockTaskService{}
    model := NewTUIModel(mockService)
    model.tasks = []*Task{
        {ID: "1", Title: "Test", Status: StatusPool},
    }
    model.cursor = 0

    // Action: simulate Enter key press
    msg := tea.KeyMsg{Type: tea.KeyEnter}
    newModel, cmd := model.Update(msg)

    // Assert: command returned (async operation)
    assert.NotNil(t, cmd)

    // Execute command and handle result
    resultMsg := cmd()
    finalModel, _ := newModel.(model).Update(resultMsg)

    // Verify service was called
    assert.Equal(t, 1, mockService.PickTaskCallCount)

    // Verify view reflects change
    view := finalModel.(model).View()
    assert.Contains(t, view, "today")
}
```

### 4. End-to-End Tests

**Scope:** Full user workflows using compiled binary.

**Critical Journeys:**
1. Add task via CLI, verify via TUI
2. Pick task in TUI, defer via CLI, complete via TUI, verify in report
3. Sync push/pull cycle with conflict resolution

**Implementation:**
```bash
#!/bin/bash
# e2e_test.sh

set -e

# Setup
export TOGO_DATA_DIR="$(mktemp -d)"
trap "rm -rf $TOGO_DATA_DIR" EXIT

# Build binary
go build -o /tmp/todo cmd/todo/main.go

# Test workflow
/tmp/todo add "Buy milk" -t shopping
/tmp/todo add "Write code" -t work

# Verify list
OUTPUT=$(/tmp/todo list --pool)
echo "$OUTPUT" | grep "Buy milk" || exit 1
echo "$OUTPUT" | grep "Write code" || exit 1

# Pick and complete
TASK_ID=$(echo "$OUTPUT" | head -1 | awk '{print $1}')
/tmp/todo pick "$TASK_ID"
/tmp/todo complete "$TASK_ID"

# Verify report
/tmp/todo report --since 1d | grep "Buy milk" || exit 1

echo "E2E tests passed!"
```

### 5. Mock Strategy

**When to Mock:**
- External dependencies (file system, network) in unit tests
- Slow operations in unit tests
- Non-deterministic operations (time, UUID generation)

**When NOT to Mock:**
- Integration tests (use real adapters with test fixtures)
- Domain model tests (test real objects)

**Example Mock:**
```go
type MockTaskRepository struct {
    LoadFunc func(ctx context.Context) (*TaskCollection, error)
    SaveFunc func(ctx context.Context, collection *TaskCollection) error
}

func (m *MockTaskRepository) Load(ctx context.Context) (*TaskCollection, error) {
    if m.LoadFunc != nil {
        return m.LoadFunc(ctx)
    }
    return &TaskCollection{tasks: make(map[TaskID]*Task)}, nil
}

// Usage in test
func TestTaskService_AddTask(t *testing.T) {
    repo := &MockTaskRepository{
        LoadFunc: func(ctx context.Context) (*TaskCollection, error) {
            return &TaskCollection{tasks: make(map[TaskID]*Task)}, nil
        },
        SaveFunc: func(ctx context.Context, c *TaskCollection) error {
            return nil // Success
        },
    }

    service := NewTaskService(repo)
    task, err := service.AddTask(context.Background(), AddTaskRequest{Title: "Test"})

    assert.NoError(t, err)
    assert.NotNil(t, task)
}
```

### Test Coverage Goals

| Package                | Target Coverage | Priority |
|------------------------|-----------------|----------|
| `internal/model`       | 85%+            | Critical |
| `internal/service`     | 80%+            | Critical |
| `internal/storage`     | 70%+            | High     |
| `internal/encryption`  | 75%+            | High     |
| `internal/sync`        | 65%+            | High     |
| `pkg/tui`              | 60%+            | Medium   |
| `cmd/todo`             | 50%+            | Medium   |

**Running Tests:**
```bash
# All tests with coverage
go test -v -cover ./...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Race detection
go test -race ./...

# Integration tests only (tagged)
go test -v -tags=integration ./...
```

---

## Error Handling

### Error Handling Philosophy

1. **Explicit over Implicit**: All errors must be explicitly handled, never silently ignored
2. **Context-Rich Errors**: Wrap errors with context at each layer
3. **Fail Fast**: Validate inputs early, return errors immediately
4. **User-Friendly Messages**: Translate technical errors to actionable messages at presentation layer

### Domain Errors

Define domain-specific error types in `internal/model/errors.go`:

```go
package model

import "errors"

var (
    // ErrTaskNotFound indicates a task ID does not exist.
    ErrTaskNotFound = errors.New("task not found")

    // ErrInvalidStatus indicates an invalid task status value.
    ErrInvalidStatus = errors.New("invalid task status")

    // ErrInvalidStateTransition indicates a state transition is not allowed.
    ErrInvalidStateTransition = errors.New("invalid state transition")

    // ErrEmptyTitle indicates a task title is empty or whitespace-only.
    ErrEmptyTitle = errors.New("task title cannot be empty")

    // ErrDuplicateTaskID indicates a task with the same ID already exists.
    ErrDuplicateTaskID = errors.New("task with this ID already exists")
)

// ValidationError wraps multiple validation failures.
type ValidationError struct {
    Field  string
    Reason string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Reason)
}
```

### Error Wrapping Pattern

Use `fmt.Errorf` with `%w` to wrap errors with context:

```go
func (s *JSONStorage) Load(ctx context.Context) (*TaskCollection, error) {
    data, err := os.ReadFile(s.path)
    if err != nil {
        if os.IsNotExist(err) {
            // Not an error, return empty collection
            return NewTaskCollection(), nil
        }
        return nil, fmt.Errorf("failed to read storage file %s: %w", s.path, err)
    }

    decrypted, err := s.encryptor.Decrypt(data)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt storage: %w", err)
    }

    var collection TaskCollection
    if err := json.Unmarshal(decrypted, &collection); err != nil {
        return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
    }

    return &collection, nil
}
```

### Error Checking with errors.Is and errors.As

```go
func (s *TaskService) RemoveTask(ctx context.Context, id TaskID) error {
    collection, err := s.repo.Load(ctx)
    if err != nil {
        return fmt.Errorf("load tasks: %w", err)
    }

    err = collection.Remove(id)
    if err != nil {
        if errors.Is(err, model.ErrTaskNotFound) {
            // Handle specific error
            return fmt.Errorf("cannot remove task: %w", err)
        }
        return err
    }

    if err := s.repo.Save(ctx, collection); err != nil {
        return fmt.Errorf("save tasks after removal: %w", err)
    }

    return nil
}
```

### Presentation Layer Error Translation

CLI and TUI must translate errors to user-friendly messages:

```go
func (c *RemoveCommand) Execute(args []string) error {
    err := c.service.RemoveTask(context.Background(), TaskID(args[0]))
    if err != nil {
        if errors.Is(err, model.ErrTaskNotFound) {
            return fmt.Errorf("task not found: %s (use 'todo list' to see available tasks)", args[0])
        }
        return fmt.Errorf("failed to remove task: %w", err)
    }

    fmt.Println("Task removed successfully")
    return nil
}
```

### Panic Policy

**Never panic in library code.** Panic only in these scenarios:
1. Unrecoverable programmer errors (e.g., nil pointer that violates contract)
2. Main function initialization failures (e.g., invalid config file format)

**Example of Acceptable Panic:**
```go
func main() {
    cfg, err := config.Load()
    if err != nil {
        // Cannot proceed without config
        panic(fmt.Sprintf("fatal: failed to load config: %v", err))
    }
    // ... rest of main
}
```

**Example of Unacceptable Panic:**
```go
// BAD: Never panic in library code
func (t *Task) Pick() error {
    if t == nil {
        panic("cannot pick nil task")
    }
    // ... implementation
}

// GOOD: Return error instead
func (t *Task) Pick() error {
    if t == nil {
        return errors.New("cannot pick nil task")
    }
    // ... implementation
}
```

### Validation Pattern

Validate inputs at API boundaries (service layer):

```go
func (s *TaskService) AddTask(ctx context.Context, req AddTaskRequest) (*Task, error) {
    // Validate request
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }

    // Business logic
    task, err := model.NewTask(req.Title, req.Tags)
    if err != nil {
        return nil, err
    }

    // ... rest of implementation
}

func (r *AddTaskRequest) Validate() error {
    if strings.TrimSpace(r.Title) == "" {
        return &ValidationError{Field: "title", Reason: "cannot be empty"}
    }
    if len(r.Title) > 500 {
        return &ValidationError{Field: "title", Reason: "exceeds 500 characters"}
    }
    return nil
}
```

---

## Security & Encryption

### Encryption Architecture

**Goals:**
1. Encrypt all data at rest
2. Support both passphrase and public key modes
3. Defend against tampering (authenticated encryption)
4. Enable secure remote sync (encrypted blob only)

### Age Encryption Integration

Use `filippo.io/age` for modern, well-maintained encryption.

**Passphrase Mode:**
```go
type AgeEncryptor struct {
    passphrase string
    scryptN    int // Work factor (default: 16384)
}

func NewAgeEncryptorWithPassphrase(passphrase string) (*AgeEncryptor, error) {
    if len(passphrase) < 12 {
        return nil, errors.New("passphrase must be at least 12 characters")
    }
    return &AgeEncryptor{passphrase: passphrase, scryptN: 16384}, nil
}

func (e *AgeEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
    // Derive key from passphrase using scrypt
    identity, err := age.NewScryptIdentity(e.passphrase)
    if err != nil {
        return nil, fmt.Errorf("create identity: %w", err)
    }

    out := &bytes.Buffer{}
    w, err := age.Encrypt(out, identity.Recipient())
    if err != nil {
        return nil, fmt.Errorf("create encryptor: %w", err)
    }

    if _, err := w.Write(plaintext); err != nil {
        return nil, fmt.Errorf("write plaintext: %w", err)
    }

    if err := w.Close(); err != nil {
        return nil, fmt.Errorf("close encryptor: %w", err)
    }

    return out.Bytes(), nil
}

func (e *AgeEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
    identity, err := age.NewScryptIdentity(e.passphrase)
    if err != nil {
        return nil, fmt.Errorf("create identity: %w", err)
    }

    r, err := age.Decrypt(bytes.NewReader(ciphertext), identity)
    if err != nil {
        return nil, fmt.Errorf("decrypt: %w", err)
    }

    plaintext, err := io.ReadAll(r)
    if err != nil {
        return nil, fmt.Errorf("read plaintext: %w", err)
    }

    return plaintext, nil
}
```

**Public Key Mode:**
```go
func NewAgeEncryptorWithKeys(publicKey, privateKey string) (*AgeEncryptor, error) {
    recipient, err := age.ParseX25519Recipient(publicKey)
    if err != nil {
        return nil, fmt.Errorf("parse public key: %w", err)
    }

    identity, err := age.ParseX25519Identity(privateKey)
    if err != nil {
        return nil, fmt.Errorf("parse private key: %w", err)
    }

    return &AgeEncryptor{recipient: recipient, identity: identity}, nil
}
```

### Key Management

**Configuration:**
```yaml
# $XDG_CONFIG_HOME/togo/config.yaml
encryption:
  mode: "passphrase"  # or "pubkey"
  # For passphrase mode: stored securely in system keyring
  # For pubkey mode:
  public_key_path: "~/.config/togo/keys/id_age.pub"
  private_key_path: "~/.config/togo/keys/id_age"
```

**CLI Commands:**
```bash
# Generate new age key pair
todo key generate

# Import existing age key
todo key import ~/.ssh/age_key.txt

# Change passphrase
todo key change-passphrase
```

### Security Best Practices

1. **Never Log Secrets**: Never log passphrases, keys, or decrypted data
2. **Secure Key Storage**: Use system keyring (keychain on macOS, Secret Service on Linux)
3. **Zero Sensitive Memory**: Overwrite sensitive byte slices after use
4. **Constant-Time Comparisons**: Use `subtle.ConstantTimeCompare` for passphrase checks
5. **Audit Dependencies**: Regularly audit `filippo.io/age` and crypto dependencies

**Example: Zero Sensitive Memory**
```go
func (e *AgeEncryptor) Close() error {
    if e.passphrase != "" {
        // Zero out passphrase in memory
        for i := range e.passphrase {
            e.passphrase = strings.Replace(e.passphrase, string(e.passphrase[i]), "0", -1)
        }
    }
    return nil
}
```

### Threat Model

**In Scope:**
- Protect data at rest from unauthorized access
- Protect data in transit (encrypted blob over git)
- Prevent tampering with task data

**Out of Scope (V1):**
- Protect against malicious code execution on user's machine
- Protect against compromised git remote (blob is encrypted)
- Multi-user access control (single-user design)

---

## Concurrency Patterns

### Concurrency Philosophy

- **Default to Sequential**: Only introduce concurrency where it provides clear value
- **Prefer Simplicity**: Avoid goroutines unless necessary for performance or responsiveness
- **Use Channels for Communication**: Share memory by communicating (Go proverb)
- **Always Handle Cancellation**: Respect context cancellation

### TUI Concurrency (Bubble Tea Commands)

Bubble Tea uses a message-passing model for concurrency:

```go
// Async operation as tea.Cmd
func (m model) saveTasksCmd(collection *TaskCollection) tea.Cmd {
    return func() tea.Msg {
        err := m.service.Save(context.Background(), collection)
        return tasksSavedMsg{err: err}
    }
}

// Handle result in Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tasksSavedMsg:
        if msg.err != nil {
            m.status = fmt.Sprintf("Error: %v", msg.err)
        } else {
            m.status = "Tasks saved successfully"
        }
        return m, nil
    }
    return m, nil
}
```

### Sync Operations Concurrency

Sync operations may involve network I/O, so use context for cancellation:

```go
func (g *GitSyncAdapter) Push(ctx context.Context) error {
    // Check context before long operation
    if ctx.Err() != nil {
        return ctx.Err()
    }

    // Use exec.CommandContext for external commands
    cmd := exec.CommandContext(ctx, "git", "push", "origin", "main")
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("git push failed: %s: %w", output, err)
    }

    return nil
}
```

### Worker Pool for Report Generation (Future)

If report generation becomes CPU-intensive:

```go
func (r *ReportService) GenerateLargeReport(ctx context.Context, opts ReportOptions) (*Report, error) {
    tasks, err := r.repo.Load(ctx)
    if err != nil {
        return nil, err
    }

    // Fan-out pattern
    numWorkers := runtime.NumCPU()
    taskChunks := chunkTasks(tasks, numWorkers)

    results := make(chan *ReportSection, numWorkers)
    errors := make(chan error, numWorkers)

    var wg sync.WaitGroup
    for _, chunk := range taskChunks {
        wg.Add(1)
        go func(c []*Task) {
            defer wg.Done()
            section, err := r.processChunk(c, opts)
            if err != nil {
                errors <- err
                return
            }
            results <- section
        }(chunk)
    }

    // Wait and close channels
    go func() {
        wg.Wait()
        close(results)
        close(errors)
    }()

    // Collect results
    report := &Report{}
    for section := range results {
        report.Sections = append(report.Sections, section)
    }

    // Check for errors
    select {
    case err := <-errors:
        return nil, err
    default:
    }

    return report, nil
}
```

### Race Detector

Always run tests with race detector during development:

```bash
go test -race ./...
```

---

## Go Best Practices

### 1. Effective Go Principles

- **Gofmt**: Always format code with `gofmt` (enforced in CI)
- **Interfaces**: Accept interfaces, return structs
- **Errors**: Errors are values, handle them explicitly
- **Named Returns**: Use sparingly, only for documentation value
- **Defer**: Use for cleanup (file close, mutex unlock)

### 2. Idiomatic Error Handling

```go
// GOOD: Check error immediately, handle or wrap
func loadConfig() (*Config, error) {
    data, err := os.ReadFile("config.yaml")
    if err != nil {
        return nil, fmt.Errorf("read config: %w", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parse config: %w", err)
    }

    return &cfg, nil
}

// BAD: Nested if statements
func loadConfig() (*Config, error) {
    data, err := os.ReadFile("config.yaml")
    if err == nil {
        var cfg Config
        err = yaml.Unmarshal(data, &cfg)
        if err == nil {
            return &cfg, nil
        }
    }
    return nil, err
}
```

### 3. Interface Design

```go
// GOOD: Small, focused interface
type Saver interface {
    Save(context.Context, *TaskCollection) error
}

// GOOD: Compose interfaces
type Repository interface {
    Loader
    Saver
}

// BAD: Large interface (violates ISP)
type Storage interface {
    Save(context.Context, *TaskCollection) error
    Load(context.Context) (*TaskCollection, error)
    Delete(context.Context, TaskID) error
    List(context.Context, TaskFilter) ([]*Task, error)
    Count(context.Context) (int, error)
}
```

### 4. Struct Initialization

```go
// GOOD: Use field names for clarity
task := &Task{
    ID:        NewTaskID(),
    Title:     title,
    Status:    StatusPool,
    CreatedAt: time.Now(),
}

// BAD: Positional initialization (fragile)
task := &Task{NewTaskID(), title, "", time.Now(), nil, StatusPool, nil, 0, nil}
```

### 5. Context Usage

```go
// GOOD: Context as first parameter
func (s *TaskService) AddTask(ctx context.Context, req AddTaskRequest) (*Task, error) {
    // ... implementation
}

// GOOD: Check context cancellation
func (s *TaskService) LongOperation(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    // ... continue
}

// BAD: Context as field (violates Go conventions)
type TaskService struct {
    ctx context.Context // Don't do this
}
```

### 6. Error Message Conventions

```go
// GOOD: Lowercase, no punctuation, actionable context
return fmt.Errorf("failed to load tasks from %s: %w", path, err)

// BAD: Uppercase, punctuation, vague
return fmt.Errorf("Error: Failed to load tasks.")
```

### 7. Package Naming

```go
// GOOD: Short, lowercase, singular noun
package task
package storage
package sync

// BAD: Plural, generic, or stuttering
package tasks  // Plural
package utils  // Too generic
package taskmanager // Stutters with task.TaskManager
```

### 8. Avoid Package-Level State

```go
// BAD: Global mutable state
package config

var GlobalConfig *Config

func Load() error {
    GlobalConfig = &Config{...}
    return nil
}

// GOOD: Return value, let caller manage lifecycle
package config

func Load() (*Config, error) {
    return &Config{...}, nil
}
```

### 9. Use Table-Driven Tests

Already covered in [Testing Strategy](#testing-strategy), but critical for Go idioms.

### 10. Document Public APIs

```go
// Task represents a unit of work in the bullet journal system.
// Tasks have three possible states: pool (not scheduled), today (scheduled for current focus),
// and done (completed).
//
// Tasks are identified by a unique UUID and maintain creation timestamps.
// Optional due dates can be set to provide temporal constraints.
type Task struct {
    // ID is the unique identifier for this task (UUID v4).
    ID TaskID `json:"id"`

    // Title is the human-readable description of the task.
    // Must be non-empty after trimming whitespace.
    Title string `json:"title"`

    // Status indicates the current lifecycle state.
    // Valid values: pool, today, done.
    Status TaskStatus `json:"status"`
}

// Pick moves a task from the pool to today's focus list.
// Returns an error if the task is already done (invalid state transition).
func (t *Task) Pick() error {
    // ... implementation
}
```

---

## Code Organization Guidelines

### 1. File Naming Conventions

- **Implementation**: `task.go`, `storage.go`, `sync.go`
- **Tests**: `task_test.go`, `storage_test.go`
- **Interfaces**: Define in same file as implementation or in dedicated `interface.go`
- **Mocks**: `mock_storage_test.go` (test files can access unexported identifiers)

### 2. File Structure Template

```go
// Package comment (required for all packages)
// Package model implements the core domain logic for task management.
package model

import (
    // Standard library imports
    "context"
    "errors"
    "fmt"
    "time"

    // External imports (separated by blank line)
    "github.com/google/uuid"
)

// Constants
const (
    MaxTitleLength = 500
)

// Errors (domain errors defined at package level)
var (
    ErrTaskNotFound = errors.New("task not found")
)

// Type definitions
type Task struct {
    // ... fields
}

// Factory functions
func NewTask(title string) (*Task, error) {
    // ... implementation
}

// Methods (receiver methods grouped with their type)
func (t *Task) Pick() error {
    // ... implementation
}

func (t *Task) Defer() error {
    // ... implementation
}

// Helper functions (private, at end of file)
func validateTitle(title string) error {
    // ... implementation
}
```

### 3. Import Grouping

```go
import (
    // Group 1: Standard library
    "context"
    "errors"
    "fmt"
    "time"

    // Group 2: External dependencies
    "github.com/charmbracelet/bubbletea"
    "github.com/google/uuid"

    // Group 3: Internal packages (project-specific)
    "togo/internal/model"
    "togo/internal/storage"
)
```

### 4. Exported vs Unexported

**Export (PascalCase):**
- Public API functions and types
- Interface definitions
- Package-level errors meant for external use

**Unexported (camelCase):**
- Helper functions
- Internal implementation details
- Struct fields that should be accessed via methods

```go
// Exported: part of public API
type TaskService struct {
    repo TaskRepository
}

// Exported: factory function
func NewTaskService(repo TaskRepository) *TaskService {
    return &TaskService{repo: repo}
}

// Exported: public method
func (s *TaskService) AddTask(ctx context.Context, req AddTaskRequest) (*Task, error) {
    return s.addTaskImpl(ctx, req)
}

// Unexported: implementation detail
func (s *TaskService) addTaskImpl(ctx context.Context, req AddTaskRequest) (*Task, error) {
    // ... implementation
}
```

### 5. Package Documentation

Every package must have a package comment:

```go
// Package model implements the core domain model for togo task management.
//
// This package defines the Task entity and related value objects (TaskID, TaskStatus).
// It enforces business rules and state transitions for the task lifecycle:
//   - Tasks start in the "pool" (not scheduled)
//   - Tasks can be "picked" for today's focus
//   - Tasks can be "deferred" back to the pool
//   - Tasks can be "completed" (marked as done)
//
// All business logic is isolated from infrastructure concerns (storage, UI).
// This package has zero dependencies on external libraries.
package model
```

### 6. Directory Organization

```
internal/model/
├── task.go              # Task entity + methods
├── task_test.go         # Task unit tests
├── task_id.go           # TaskID value object
├── task_status.go       # TaskStatus value object
├── task_filter.go       # TaskFilter value object
├── task_collection.go   # TaskCollection aggregate
├── task_collection_test.go
└── errors.go            # Domain errors
```

**Rationale:** Group related files by entity/aggregate, keep test files alongside implementation.

---

## Integration Points

### 1. CLI to Service Layer

```
┌──────────────────┐
│  main.go         │
│  (dispatcher)    │
└────────┬─────────┘
         │
         │ Parses command + flags
         ▼
┌──────────────────┐
│  handleAdd()     │
│  (cmd/todo)      │
└────────┬─────────┘
         │
         │ Creates AddTaskRequest
         ▼
┌──────────────────┐
│  TaskService     │
│  (orchestration) │
└────────┬─────────┘
         │
         │ Domain operations
         ▼
┌──────────────────┐
│  Task Model      │
└──────────────────┘
```

**Example:**
```go
// cmd/todo/commands/add.go
func handleAdd(args []string, svc *service.TaskService) error {
    fs := flag.NewFlagSet("add", flag.ExitOnError)
    notes := fs.String("n", "", "Task notes")
    tags := fs.String("t", "", "Comma-separated tags")
    fs.Parse(args)

    if fs.NArg() == 0 {
        return errors.New("task title required")
    }

    req := service.AddTaskRequest{
        Title: fs.Arg(0),
        Notes: *notes,
        Tags:  parseTagsString(*tags),
    }

    task, err := svc.AddTask(context.Background(), req)
    if err != nil {
        return fmt.Errorf("failed to add task: %w", err)
    }

    // Pretty print with Lip Gloss
    fmt.Println(successStyle.Render(fmt.Sprintf(
        "✓ Task added: %s (ID: %s)",
        task.Title,
        task.ID,
    )))
    return nil
}
```

### 2. TUI to Service Layer

```
┌──────────────────┐
│  Bubble Tea      │
│  Model/Update    │
└────────┬─────────┘
         │
         │ KeyMsg event
         ▼
┌──────────────────┐
│  tea.Cmd         │  (async operation)
└────────┬─────────┘
         │
         │ Calls service
         ▼
┌──────────────────┐
│  TaskService     │
└────────┬─────────┘
         │
         │ Returns result
         ▼
┌──────────────────┐
│  Custom Msg      │  (taskPickedMsg)
└────────┬─────────┘
         │
         │ Handled in Update
         ▼
┌──────────────────┐
│  Model Updated   │
│  View Re-renders │
└──────────────────┘
```

**Example:**
```go
// pkg/tui/update.go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            return m, m.pickTaskCmd()
        }
    case taskPickedMsg:
        if msg.err != nil {
            m.status = fmt.Sprintf("Error: %v", msg.err)
        } else {
            m.status = "Task picked successfully"
            m.refreshTasks()
        }
    }
    return m, nil
}

func (m model) pickTaskCmd() tea.Cmd {
    taskID := m.tasks[m.cursor].ID
    return func() tea.Msg {
        err := m.service.PickTask(context.Background(), taskID)
        return taskPickedMsg{err: err}
    }
}
```

### 3. Service to Storage Layer

```
┌──────────────────┐
│  TaskService     │
└────────┬─────────┘
         │
         │ Calls Load/Save
         ▼
┌──────────────────┐
│  TaskRepository  │  (interface)
└────────┬─────────┘
         │
         │ Implementation
         ▼
┌──────────────────┐
│  JSONStorage     │
└────────┬─────────┘
         │
         │ Uses Encryptor
         ▼
┌──────────────────┐
│  AgeEncryptor    │
└────────┬─────────┘
         │
         │ Writes to disk
         ▼
┌──────────────────┐
│  File System     │
└──────────────────┘
```

**Example:**
```go
// internal/service/task_service.go
func (s *TaskService) AddTask(ctx context.Context, req AddTaskRequest) (*Task, error) {
    // Load current collection
    collection, err := s.repo.Load(ctx)
    if err != nil {
        return nil, fmt.Errorf("load tasks: %w", err)
    }

    // Create new task (domain logic)
    task, err := model.NewTask(req.Title, req.Tags)
    if err != nil {
        return nil, fmt.Errorf("create task: %w", err)
    }

    // Add to collection
    if err := collection.Add(task); err != nil {
        return nil, fmt.Errorf("add to collection: %w", err)
    }

    // Persist
    if err := s.repo.Save(ctx, collection); err != nil {
        return nil, fmt.Errorf("save tasks: %w", err)
    }

    return task, nil
}
```

### 4. Sync to Storage Layer

```
┌──────────────────┐
│  SyncCommand     │
└────────┬─────────┘
         │
         │ Calls Push/Pull
         ▼
┌──────────────────┐
│  GitSyncAdapter  │
└────────┬─────────┘
         │
         │ Reads encrypted file
         ▼
┌──────────────────┐
│  JSONStorage     │
└────────┬─────────┘
         │
         │ Git operations
         ▼
┌──────────────────┐
│  exec.Command    │  (git CLI)
└──────────────────┘
```

---

## Quality Standards

### 1. Code Review Checklist

Before merging any PR, verify:

- [ ] All tests pass (`go test ./...`)
- [ ] Race detector passes (`go test -race ./...`)
- [ ] Code is formatted (`gofmt -s -w .`)
- [ ] No vet warnings (`go vet ./...`)
- [ ] Coverage meets minimum thresholds
- [ ] Public APIs are documented
- [ ] Errors are wrapped with context
- [ ] No exported unexported dependencies (e.g., `internal` packages used in `pkg`)
- [ ] SOLID principles respected
- [ ] No panics in library code
- [ ] Context cancellation respected

### 2. Continuous Integration

`.github/workflows/ci.yml`:

```yaml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'

      - name: Format check
        run: |
          if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
            echo "Code not formatted. Run 'gofmt -s -w .'"
            exit 1
          fi

      - name: Vet
        run: go vet ./...

      - name: Test
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Coverage
        run: |
          go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' | \
            awk '{if ($1 < 70) {print "Coverage too low: " $1 "%"; exit 1}}'
```

### 3. Pre-Commit Hooks

`.git/hooks/pre-commit`:

```bash
#!/bin/bash

# Format code
echo "Running gofmt..."
gofmt -s -w .

# Vet code
echo "Running go vet..."
go vet ./...
if [ $? -ne 0 ]; then
    echo "go vet failed"
    exit 1
fi

# Run tests
echo "Running tests..."
go test -short ./...
if [ $? -ne 0 ]; then
    echo "Tests failed"
    exit 1
fi

echo "Pre-commit checks passed!"
```

### 4. Documentation Standards

- **README.md**: User-facing documentation (installation, usage, examples)
- **ARCHITECTURE.md**: This document (architectural decisions, patterns)
- **IMPLEMENTATION_PLAN.md**: Implementation roadmap and milestones
- **GoDoc**: All exported identifiers must have comments
- **Inline Comments**: Complex logic should have explanatory comments

### 5. Performance Benchmarks

For critical paths (e.g., report generation):

```go
func BenchmarkGenerateReport(b *testing.B) {
    service := setupBenchmarkService()
    opts := ReportOptions{Since: time.Now().Add(-7 * 24 * time.Hour)}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := service.GenerateReport(context.Background(), opts)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

Run benchmarks:
```bash
go test -bench=. -benchmem ./...
```

### 6. Dependency Management

- Pin dependency versions in `go.mod`
- Run `go mod tidy` regularly
- Audit dependencies for security: `go list -m all | nancy sleuth`
- Prefer standard library over third-party when possible

---

## Appendix: Decision Log

### ADR-001: Encrypted JSON vs bbolt

**Decision**: Use encrypted JSON for MVP.

**Rationale:**
- Simpler implementation (fewer moving parts)
- Easier to inspect and debug (decrypt manually)
- Sufficient for expected dataset size (<10k tasks)
- Can migrate to bbolt later if performance degrades

**Trade-offs:**
- Load/save requires full deserialization (O(n) with task count)
- No indexed queries (must filter in memory)
- Acceptable for V1 with small datasets

### ADR-002: Git vs Cloud Sync

**Decision**: Git-backed sync for MVP.

**Rationale:**
- Users already familiar with git workflows
- No third-party service dependencies
- Works with any git remote (GitHub, GitLab, self-hosted)
- Encrypted blob ensures privacy even on public repos

**Trade-offs:**
- Requires git installation
- Manual conflict resolution initially
- Can add cloud sync adapters later (S3, Dropbox)

### ADR-003: TUI-First vs CLI-First

**Decision**: TUI-first approach with optional CLI commands.

**Rationale:**
- TUI provides best user experience for task management
- Interactive interface suits the bullet journal metaphor
- CLI remains available for quick operations and scripting
- No args launches TUI (default behavior)
- Charm libraries (Bubble Tea, Bubbles, Lip Gloss) provide excellent TUI framework

**Trade-offs:**
- TUI has steeper learning curve for development (mitigated by comprehensive CHARM_RESEARCH.md)
- TUI requires terminal with ANSI color support
- Binary size slightly larger than pure CLI
- Benefits far outweigh costs for this use case

### ADR-004: No Cobra Dependency

**Decision**: Use Go standard library `flag` package for CLI parsing, no Cobra.

**Rationale:**
- Eliminates external dependency for simple CLI parsing needs
- TUI is the primary interface, CLI is secondary
- Flag package sufficient for basic command routing
- Reduces binary size and dependency complexity
- User explicitly requested no Cobra

**Trade-offs:**
- Manual command routing required (simple switch statement)
- No automatic help generation (write manually)
- No shell completion out-of-box
- Trade-offs acceptable given TUI-first approach

**Implementation Pattern:**
```go
func main() {
    if len(os.Args) == 1 {
        launchTUI()  // No args = default to TUI
        return
    }

    command := os.Args[1]
    switch command {
    case "add": handleAdd(os.Args[2:])
    case "list": handleList(os.Args[2:])
    // ... etc
    }
}
```

---

## Conclusion

This architecture document defines a clean, maintainable, and extensible foundation for the togo project. By adhering to SOLID principles, Domain-Driven Design, Go best practices, and embracing a TUI-first approach with Charm libraries, we ensure:

1. **Beautiful User Experience**: Bubble Tea provides an interactive, visually appealing interface that makes task management delightful
2. **Testability**: Clear boundaries and pure functions (Model-Update-View) enable comprehensive testing
3. **Maintainability**: Separation of concerns makes code easy to understand and modify
4. **Extensibility**: Interfaces and adapters allow adding new features without breaking existing code
5. **Idiomatic Go**: Following Go conventions ensures the codebase feels natural to Go developers
6. **No Unnecessary Dependencies**: Using standard library for CLI parsing keeps the dependency graph minimal

**Key Architectural Decisions:**

- **TUI-First**: Bubble Tea as primary interface, with optional CLI for scripting
- **No Cobra**: Standard library `flag` package for lightweight CLI parsing
- **Charm Ecosystem**: Bubbles components and Lip Gloss styling for consistent, beautiful UI
- **Hexagonal Architecture**: Clean separation between domain, application, and infrastructure layers
- **Encrypted Storage**: Age encryption for all persisted data

**Next Steps:**

1. Review and approve this architecture document
2. Begin implementation following the phases in IMPLEMENTATION_PLAN.md:
   - Phase 1: Core domain and encrypted storage
   - Phase 2: Interactive TUI with Bubble Tea
   - Phase 3: Sync infrastructure and optional CLI
   - Phase 4: Polish, reports, and advanced features
3. Continuously update this document as architectural decisions evolve

**Maintainers:**

- Review this document quarterly
- Update when making significant architectural changes
- Keep ADR section current with new decisions
- Ensure new contributors read and understand this document before contributing

---

**Document Version History:**

| Version | Date       | Author | Changes                     |
|---------|------------|--------|-----------------------------|
| 1.0     | 2025-11-09 | Claude | Initial architecture design |
