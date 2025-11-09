# togo Project Task List

**Project**: togo - A beautiful, keyboard-driven CLI/TUI task manager
**Architecture**: TUI-first using Bubble Tea, Bubbles, and Lip Gloss (no Cobra)
**Total Tasks**: 87
**Status**: Ready for implementation

---

## Task Status Legend

- ⏸️ **Pending**: Not yet started
- 🔄 **In Progress**: Currently being worked on
- ✅ **Completed**: Finished and verified

---

## Phase 1: Core Domain & Storage Foundation (Tasks 1-31)

**Goal**: Build solid foundation with domain model, encrypted storage, and service layer
**Estimated Duration**: 1 week
**Key Deliverables**: Tested domain model, encrypted JSON storage, TaskService with all operations

### Project Setup (Tasks 1-2)

#### Task 1: Create project directory structure per architecture plan

**Status**: ✅ Completed
**Layer**: Project Setup
**Dependencies**: None
**Complexity**: Simple
**Description**:

- Create directory structure: `cmd/togo/`, `pkg/tui/`, `internal/{model,storage,encryption,service,sync}/`
- Initialize `go.mod` with module path
- Add initial dependencies: `github.com/charmbracelet/bubbletea`, `bubbles`, `lipgloss`
- Create placeholder `main.go` and `README.md`

**Acceptance Criteria**:

- All directories created
- `go.mod` initialized with correct dependencies
- Project builds with `go build ./...`

---

#### Task 2: Set up testing infrastructure and CI configuration

**Status**: ✅ Completed
**Layer**: Project Setup
**Dependencies**: Task 1
**Complexity**: Simple
**Description**:

- Look at existing CI config in .github/
- move the existing tests to a tests/ directory and get them working in that Directory

**Acceptance Criteria**:

- CI runs on push/PR
- All quality checks pass

---

### Domain Model Layer - Value Objects (Tasks 3-8)

#### Task 3: Implement TaskID value object with UUID generation

**Status**: ✅ Completed
**Layer**: Domain Model
**Dependencies**: Task 1
**Complexity**: Simple
**File**: `internal/model/task_id.go`

**Description**:

- Create TaskID type (string wrapper)
- Implement `NewTaskID()` factory with UUID v4 generation
- Implement `ParseTaskID(string)` with validation
- Implement `String()` method

**Test Cases**:

- `TestNewTaskID_GeneratesValidUUID`
- `TestParseTaskID_ValidUUID_Success`
- `TestParseTaskID_InvalidString_ReturnsError`
- `TestTaskID_String_ReturnsUnderlyingValue`

**Acceptance Criteria**:

- TaskID type defined
- UUID generation works
- Validation prevents invalid UUIDs
- All tests passing with >85% coverage

---

#### Task 4: Implement TaskStatus value object with validation

**Status**: ⏸️ Pending
**Layer**: Domain Model
**Dependencies**: Task 1
**Complexity**: Simple
**File**: `internal/model/task_status.go`

**Description**:

- Create TaskStatus type (string)
- Define constants: `StatusPool`, `StatusToday`, `StatusDone`
- Implement `Valid()` method
- Support JSON marshaling/unmarshaling

**Test Cases**:

- `TestTaskStatus_Valid_AllValidStatuses`
- `TestTaskStatus_Valid_InvalidStatus_ReturnsFalse`
- `TestTaskStatus_String_ReturnsValue`
- `TestTaskStatus_JSONRoundTrip`

**Acceptance Criteria**:

- TaskStatus type with three constants
- Validation works correctly
- JSON serialization works

---

#### Task 5: Define domain errors in errors.go

**Status**: ⏸️ Pending
**Layer**: Domain Model
**Dependencies**: Task 1
**Complexity**: Simple
**File**: `internal/model/errors.go`

**Description**:

- Define sentinel errors using `errors.New()`
- Create `ValidationError` struct with Field and Reason
- Define: `ErrTaskNotFound`, `ErrInvalidStatus`, `ErrInvalidStateTransition`, `ErrEmptyTitle`, `ErrDuplicateTaskID`

**Test Cases**:

- `TestValidationError_Error_FormatsCorrectly`
- `TestDomainErrors_AreUnique`

**Acceptance Criteria**:

- All domain errors defined
- ValidationError implements error interface
- Errors follow Go conventions (lowercase, no punctuation)

---

#### Task 6: Implement TaskFilter value object for queries

**Status**: ⏸️ Pending
**Layer**: Domain Model
**Dependencies**: Task 4
**Complexity**: Moderate
**File**: `internal/model/task_filter.go`

**Description**:

- Create TaskFilter struct with optional fields (Status, Tags, DueDate range, Limit)
- Implement `Matches(*Task) bool` for in-memory filtering
- Use pointer types for optional filters

**Test Cases**:

- `TestTaskFilter_Matches_StatusFilter`
- `TestTaskFilter_Matches_TagsFilter`
- `TestTaskFilter_Matches_DueDateRange`
- `TestTaskFilter_Matches_MultipleFilters_AllMustMatch`
- `TestTaskFilter_Matches_NoFilters_MatchesAll`

**Acceptance Criteria**:

- TaskFilter supports all filter criteria
- Matches() evaluates correctly
- Table-driven tests passing

---

### Domain Model Layer - Task Entity (Tasks 7-10)

#### Task 7: Implement Task entity core with factory function

**Status**: ⏸️ Pending
**Layer**: Domain Model
**Dependencies**: Tasks 3, 4, 5
**Complexity**: Moderate
**File**: `internal/model/task.go`

**Description**:

- Create Task struct with all fields: ID, Title, Notes, CreatedAt, DueDate, Status, CompletedAt, DeferredCount, Tags
- Implement `NewTask(title string, tags []string)` factory
- Validate title is non-empty (trimmed)
- Initialize with StatusPool and current timestamp
- Add JSON struct tags

**Test Cases**:

- `TestNewTask_ValidTitle_Success`
- `TestNewTask_EmptyTitle_ReturnsError`
- `TestNewTask_WhitespaceTitle_ReturnsError`
- `TestNewTask_SetsDefaultValues`
- `TestTask_JSONRoundTrip`

**Acceptance Criteria**:

- Task struct fully defined
- Factory enforces invariants
- JSON marshaling works
- All tests passing

---

#### Task 8: Implement Task state transition methods (Pick, Defer, Complete)

**Status**: ⏸️ Pending
**Layer**: Domain Model
**Dependencies**: Task 7
**Complexity**: Moderate
**File**: `internal/model/task.go`

**Description**:

- Implement `Pick()`: pool→today, today→today (idempotent), done→error
- Implement `Defer()`: today→pool (increment DeferredCount), pool→pool (idempotent), done→error
- Implement `Complete()`: any→done (set CompletedAt), done→idempotent
- Return domain errors for invalid transitions

**Test Cases**:

- `TestTask_Pick_FromPool_Success`
- `TestTask_Pick_FromToday_Idempotent`
- `TestTask_Pick_FromDone_ReturnsError`
- `TestTask_Defer_FromToday_IncrementsCount`
- `TestTask_Defer_FromPool_Idempotent`
- `TestTask_Defer_FromDone_ReturnsError`
- `TestTask_Complete_SetsCompletedAt`
- `TestTask_Complete_AlreadyDone_Idempotent`

**Acceptance Criteria**:

- All three methods implemented
- State transitions validated
- DeferredCount increments correctly
- Table-driven tests passing

---

#### Task 9: Implement Task validation with Validate() method

**Status**: ⏸️ Pending
**Layer**: Domain Model
**Dependencies**: Task 8
**Complexity**: Simple
**File**: `internal/model/task.go`

**Description**:

- Implement `Validate()` method
- Check: valid ID, valid status, non-empty title, CompletedAt set iff status==done, DeferredCount >= 0
- Return ValidationError with context

**Test Cases**:

- `TestTask_Validate_ValidTask`
- `TestTask_Validate_InvalidID_ReturnsError`
- `TestTask_Validate_InvalidStatus_ReturnsError`
- `TestTask_Validate_EmptyTitle_ReturnsError`
- `TestTask_Validate_CompletedAtMismatch_ReturnsError`
- `TestTask_Validate_NegativeDeferredCount_ReturnsError`

**Acceptance Criteria**:

- Validate() checks all invariants
- ValidationError provides context
- All tests passing

---

### Domain Model Layer - Aggregate (Tasks 10-11)

#### Task 10: Implement TaskCollection aggregate with CRUD operations

**Status**: ⏸️ Pending
**Layer**: Domain Model
**Dependencies**: Tasks 6, 9
**Complexity**: Moderate
**File**: `internal/model/task_collection.go`

**Description**:

- Create TaskCollection struct with internal `map[TaskID]*Task`
- Add Metadata struct (Version, LastModified, EncryptionMode, Salt)
- Implement: `Add()`, `Get()`, `Remove()`, `Find(TaskFilter)`, `All()`
- All() returns slice sorted by CreatedAt (newest first)
- Add() enforces no duplicate IDs

**Test Cases**:

- `TestTaskCollection_Add_NewTask_Success`
- `TestTaskCollection_Add_DuplicateID_ReturnsError`
- `TestTaskCollection_Get_ExistingTask_Success`
- `TestTaskCollection_Get_NonExistentTask_ReturnsError`
- `TestTaskCollection_Remove_ExistingTask_Success`
- `TestTaskCollection_Remove_NonExistentTask_ReturnsError`
- `TestTaskCollection_Find_FilterByStatus`
- `TestTaskCollection_Find_FilterByTags`
- `TestTaskCollection_All_ReturnsSortedByCreatedAt`

**Acceptance Criteria**:

- TaskCollection with all CRUD operations
- Find() respects all filter criteria
- All tests passing with >85% coverage

---

#### Task 11: Write comprehensive unit tests for domain model

**Status**: ⏸️ Pending
**Layer**: Domain Model
**Dependencies**: Tasks 3-10
**Complexity**: Moderate
**File**: `internal/model/*_test.go`

**Description**:

- Ensure all domain model components have >85% test coverage
- Use table-driven tests for comprehensive scenarios
- Test edge cases and error conditions
- Test invariant enforcement

**Acceptance Criteria**:

- All domain model tests passing
- Coverage >85% for internal/model package
- Edge cases covered

---

### Infrastructure Layer - Encryption (Tasks 12-14)

#### Task 12: Define Encryptor interface

**Status**: ⏸️ Pending
**Layer**: Infrastructure (Interface)
**Dependencies**: Task 1
**Complexity**: Simple
**File**: `internal/encryption/encryptor.go`

**Description**:

- Define interface: `Encrypt([]byte) ([]byte, error)` and `Decrypt([]byte) ([]byte, error)`
- Document AEAD requirement
- Document thread-safety expectations
- Document error cases

**Acceptance Criteria**:

- Encryptor interface defined
- Comprehensive package documentation
- Contracts documented

---

#### Task 13: Implement AgeEncryptor with passphrase mode

**Status**: ⏸️ Pending
**Layer**: Infrastructure
**Dependencies**: Task 12
**Complexity**: Moderate
**File**: `internal/encryption/age_encryptor.go`

**Description**:

- Implement AgeEncryptor using `filippo.io/age`
- Support passphrase mode with scrypt derivation
- Enforce minimum passphrase length (12 characters)
- Use scryptN = 16384 work factor
- Implement Encrypt() and Decrypt()

**Test Cases**:

- `TestAgeEncryptor_NewWithPassphrase_TooShort_ReturnsError`
- `TestAgeEncryptor_NewWithPassphrase_ValidLength_Success`
- `TestAgeEncryptor_Encrypt_Decrypt_RoundTrip`
- `TestAgeEncryptor_Decrypt_WrongPassphrase_ReturnsError`
- `TestAgeEncryptor_Decrypt_CorruptedData_ReturnsError`
- `TestAgeEncryptor_Encrypt_DifferentOutputEachTime`

**Acceptance Criteria**:

- AgeEncryptor implements Encryptor interface
- Passphrase mode works correctly
- All tests passing
- Robust error handling

---

#### Task 14: Implement NoopEncryptor for testing

**Status**: ⏸️ Pending
**Layer**: Infrastructure
**Dependencies**: Task 12
**Complexity**: Simple
**File**: `internal/encryption/noop_encryptor.go`

**Description**:

- Implement NoopEncryptor that returns plaintext
- Document as TESTING ONLY
- Use in storage layer tests

**Test Cases**:

- `TestNoopEncryptor_Encrypt_ReturnsPlaintext`
- `TestNoopEncryptor_Decrypt_ReturnsPlaintext`

**Acceptance Criteria**:

- NoopEncryptor implements Encryptor
- Clearly documented as testing-only
- All tests passing

---

### Infrastructure Layer - Storage (Tasks 15-23)

#### Task 15: Define TaskRepository interface

**Status**: ⏸️ Pending
**Layer**: Infrastructure (Interface)
**Dependencies**: Task 10
**Complexity**: Simple
**File**: `internal/storage/repository.go`

**Description**:

- Define interface: `Load(ctx) (*TaskCollection, error)`, `Save(ctx, *TaskCollection) error`, `Close() error`
- Document idempotency and atomicity guarantees
- Document empty collection behavior
- Add context support for cancellation

**Acceptance Criteria**:

- TaskRepository interface defined
- Comprehensive documentation
- Contracts clearly specified

---

#### Task 16: Implement MemoryStorage adapter for testing

**Status**: ⏸️ Pending
**Layer**: Infrastructure
**Dependencies**: Task 15
**Complexity**: Simple
**File**: `internal/storage/memory_storage.go`

**Description**:

- Implement in-memory storage (field holding TaskCollection)
- Implement Load(), Save(), Close()
- Add Reset() method for test cleanup
- Thread-safe if needed (likely not for tests)

**Test Cases**:

- `TestMemoryStorage_SaveAndLoad_Success`
- `TestMemoryStorage_Load_BeforeSave_ReturnsEmpty`
- `TestMemoryStorage_Reset_ClearsData`

**Acceptance Criteria**:

- MemoryStorage implements TaskRepository
- Works correctly in memory
- Used in service layer tests

---

#### Task 17: Implement XDG Base Directory support

**Status**: ⏸️ Pending
**Layer**: Infrastructure
**Dependencies**: Task 1
**Complexity**: Simple
**File**: `internal/storage/paths.go`

**Description**:

- Implement function to get data directory
- Check `$XDG_DATA_HOME`, fallback to `~/.local/share`
- Create `togo/` subdirectory
- Return full path to `todo.json.age`

**Test Cases**:

- `TestGetDataPath_XDGSet_UsesXDG`
- `TestGetDataPath_XDGNotSet_UsesFallback`

**Acceptance Criteria**:

- XDG Base Directory spec followed
- Fallback works correctly
- All tests passing

---

#### Task 18: Implement JSONStorage file I/O operations

**Status**: ⏸️ Pending
**Layer**: Infrastructure
**Dependencies**: Tasks 15, 17
**Complexity**: Moderate
**File**: `internal/storage/json_storage.go`

**Description**:

- Implement JSONStorage struct with path and encryptor fields
- Implement Load(): read file, return empty collection if not exists
- Implement Save(): write file atomically
- Create directory if needed (os.MkdirAll)
- Constructor: `NewJSONStorage(encryptor Encryptor)`

**Test Cases**:

- `TestJSONStorage_Load_FileNotExists_ReturnsEmptyCollection`
- `TestJSONStorage_Load_FileExists_ReturnsCollection`
- `TestJSONStorage_Save_CreatesFile`
- `TestJSONStorage_Save_CreatesDirectoryIfNeeded`
- `TestJSONStorage_Save_OverwritesExistingFile`
- Use `t.TempDir()` for all tests

**Acceptance Criteria**:

- JSONStorage implements TaskRepository
- File I/O works correctly
- All tests passing with temp directories

---

#### Task 19: Integrate encryption into JSONStorage

**Status**: ⏸️ Pending
**Layer**: Infrastructure
**Dependencies**: Tasks 13, 18
**Complexity**: Moderate
**File**: `internal/storage/json_storage.go`

**Description**:

- Save flow: JSON marshal → Encrypt → Write file
- Load flow: Read file → Decrypt → JSON unmarshal
- Wrap errors with context
- Handle decryption failures gracefully

**Test Cases**:

- `TestJSONStorage_Save_DataIsEncrypted`
- `TestJSONStorage_Load_DecryptsCorrectly`
- `TestJSONStorage_Load_WrongKey_ReturnsError`
- `TestJSONStorage_RoundTrip_WithRealEncryptor`
- `TestJSONStorage_EncryptedFileContainsNoPlaintext`

**Acceptance Criteria**:

- Encryption/decryption fully integrated
- No plaintext in stored files
- All tests passing
- Round-trip test with real Age encryptor

---

#### Task 20: Write storage layer integration tests

**Status**: ⏸️ Pending
**Layer**: Integration Test
**Dependencies**: Tasks 13, 19
**Complexity**: Moderate
**File**: `internal/storage/integration_test.go`

**Description**:

- Test JSONStorage with real AgeEncryptor
- Test full round-trip: create tasks → save → load → verify
- Test encrypted file does not contain plaintext
- Use temp directory for files

**Test Cases**:

- `TestIntegration_JSONStorage_RoundTrip`
- `TestIntegration_JSONStorage_EncryptedFileNotPlaintext`
- `TestIntegration_JSONStorage_WrongPassphrase_Fails`

**Acceptance Criteria**:

- Integration tests passing
- Uses real components (no mocks)
- Temp directory cleanup

---

### Service Layer (Tasks 21-31)

#### Task 21: Define TaskService interface

**Status**: ⏸️ Pending
**Layer**: Service
**Dependencies**: Task 10
**Complexity**: Simple
**File**: `internal/service/task_service.go`

**Description**:

- Define interface with methods: AddTask, ListTasks, GetTask, PickTask, DeferTask, CompleteTask, RemoveTask
- Define request/response types (AddTaskRequest, etc.)
- Document transaction boundaries
- Document validation responsibilities

**Acceptance Criteria**:

- TaskService interface defined
- Request/response types defined
- Comprehensive documentation

---

#### Task 22: Implement TaskService - AddTask operation

**Status**: ⏸️ Pending
**Layer**: Service
**Dependencies**: Tasks 15, 21
**Complexity**: Moderate
**File**: `internal/service/task_service.go`

**Description**:

- Validate AddTaskRequest (title non-empty, reasonable length)
- Load collection from repository
- Create new Task using factory
- Add to collection
- Save collection
- Wrap errors with context at each step

**Test Cases**:

- `TestTaskService_AddTask_ValidRequest_Success`
- `TestTaskService_AddTask_EmptyTitle_ReturnsError`
- `TestTaskService_AddTask_TitleTooLong_ReturnsError`
- `TestTaskService_AddTask_SaveFails_ReturnsError`
- Use mock/memory repository

**Acceptance Criteria**:

- AddTask method implemented
- Request validation thorough
- All tests passing with mock repository

---

#### Task 23: Implement TaskService - ListTasks operation

**Status**: ⏸️ Pending
**Layer**: Service
**Dependencies**: Tasks 21, 22
**Complexity**: Simple
**File**: `internal/service/task_service.go`

**Description**:

- Load collection from repository
- Use TaskCollection.Find() with provided filter
- Return task slice
- Handle empty results gracefully

**Test Cases**:

- `TestTaskService_ListTasks_NoFilter_ReturnsAll`
- `TestTaskService_ListTasks_FilterByStatus`
- `TestTaskService_ListTasks_FilterByTags`
- `TestTaskService_ListTasks_EmptyCollection_ReturnsEmpty`

**Acceptance Criteria**:

- ListTasks method implemented
- Filtering works correctly
- All tests passing

---

#### Task 24: Implement TaskService - PickTask operation

**Status**: ⏸️ Pending
**Layer**: Service
**Dependencies**: Task 23
**Complexity**: Moderate
**File**: `internal/service/task_service.go`

**Description**:

- Load collection
- Get task by ID (handle not found)
- Call task.Pick() (domain validation)
- Save collection
- Transaction-like behavior: if any step fails, return error

**Test Cases**:

- `TestTaskService_PickTask_Success`
- `TestTaskService_PickTask_TaskNotFound_ReturnsError`
- `TestTaskService_PickTask_InvalidStateTransition_ReturnsError`
- `TestTaskService_PickTask_SaveFails_ReturnsError`

**Acceptance Criteria**:

- PickTask method implemented
- Delegates to domain model
- All tests passing
- Comprehensive error handling

---

#### Task 25: Implement TaskService - DeferTask operation

**Status**: ⏸️ Pending
**Layer**: Service
**Dependencies**: Task 24
**Complexity**: Moderate
**File**: `internal/service/task_service.go`

**Description**:

- Follow same pattern as PickTask
- Call task.Defer()
- Ensure DeferredCount increments only when moving from today→pool

**Test Cases**:

- `TestTaskService_DeferTask_Success`
- `TestTaskService_DeferTask_TaskNotFound_ReturnsError`
- `TestTaskService_DeferTask_IncrementsCount`

**Acceptance Criteria**:

- DeferTask method implemented
- Domain model method used correctly
- All tests passing

---

#### Task 26: Implement TaskService - CompleteTask operation

**Status**: ⏸️ Pending
**Layer**: Service
**Dependencies**: Task 25
**Complexity**: Moderate
**File**: `internal/service/task_service.go`

**Description**:

- Follow same pattern as PickTask/DeferTask
- Call task.Complete()
- Ensure CompletedAt timestamp set correctly

**Test Cases**:

- `TestTaskService_CompleteTask_Success`
- `TestTaskService_CompleteTask_SetsCompletedAt`

**Acceptance Criteria**:

- CompleteTask method implemented
- All tests passing

---

#### Task 27: Implement TaskService - RemoveTask operation

**Status**: ⏸️ Pending
**Layer**: Service
**Dependencies**: Task 26
**Complexity**: Simple
**File**: `internal/service/task_service.go`

**Description**:

- Load collection
- Remove task by ID
- Save collection
- Return error if task not found

**Test Cases**:

- `TestTaskService_RemoveTask_Success`
- `TestTaskService_RemoveTask_TaskNotFound_ReturnsError`
- `TestTaskService_RemoveTask_TaskIsRemoved`

**Acceptance Criteria**:

- RemoveTask method implemented
- All tests passing
- Error handling correct

---

#### Task 28: Implement TaskService - GetTask operation

**Status**: ⏸️ Pending
**Layer**: Service
**Dependencies**: Task 27
**Complexity**: Simple
**File**: `internal/service/task_service.go`

**Description**:

- Load collection
- Get task by ID
- Return task or error if not found

**Test Cases**:

- `TestTaskService_GetTask_Success`
- `TestTaskService_GetTask_NotFound_ReturnsError`

**Acceptance Criteria**:

- GetTask method implemented
- All tests passing

---

#### Task 29: Write comprehensive service layer unit tests

**Status**: ⏸️ Pending
**Layer**: Service
**Dependencies**: Tasks 22-28
**Complexity**: Moderate
**File**: `internal/service/*_test.go`

**Description**:

- Ensure all service methods have comprehensive tests
- Use mock/memory repository for isolation
- Test error cases and edge conditions
- Achieve >85% coverage

**Acceptance Criteria**:

- All service tests passing
- Coverage >85% for internal/service package

---

#### Task 30: Write service layer integration tests

**Status**: ⏸️ Pending
**Layer**: Integration Test
**Dependencies**: Tasks 20, 29
**Complexity**: Moderate
**File**: `internal/service/integration_test.go`

**Description**:

- Test TaskService with real storage (MemoryStorage or JSONStorage)
- Test complete workflows: add → pick → defer → complete
- Verify persistence across service method calls

**Test Cases**:

- `TestIntegration_TaskService_AddAndList`
- `TestIntegration_TaskService_PickDeferComplete_Workflow`
- `TestIntegration_TaskService_Remove_WorksCorrectly`

**Acceptance Criteria**:

- Integration tests covering main workflows
- Tests use realistic storage
- All tests passing

---

#### Task 31: Phase 1 verification - Run all tests and verify coverage

**Status**: ⏸️ Pending
**Layer**: Quality Assurance
**Dependencies**: Tasks 1-30
**Complexity**: Simple

**Description**:

- Run `go test ./...` and verify all tests pass
- Run coverage report and verify >85% coverage for all packages
- Run `go vet` and `golangci-lint`
- Verify CI pipeline passes

**Acceptance Criteria**:

- All tests passing
- Coverage >85%
- No lint errors
- CI green

---

## Phase 2: Interactive TUI (Tasks 32-54)

**Goal**: Build beautiful, functional TUI with Bubble Tea, Bubbles, and Lip Gloss
**Estimated Duration**: 1 week
**Key Deliverables**: Fully functional TUI with tab navigation, task operations, styling

### TUI Foundation (Tasks 32-37)

#### Task 32: Create TUI package structure and Model struct

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 21
**Complexity**: Moderate
**File**: `pkg/tui/model.go`

**Description**:

- Create Model struct with: TaskService, tasks list, view state, active tab, selected index
- Define view state enum: listView, addView, detailView, reportView
- Define tab enum: poolTab, todayTab, doneTab
- Include bubbles components as fields (list.Model, textinput.Model, textarea.Model, spinner.Model, help.Model)
- Add window dimensions fields

**Test Cases**:

- `TestModel_Initialization_DefaultState`
- `TestModel_Initialization_InjectsDependencies`

**Acceptance Criteria**:

- Model struct defined
- Enums for states and tabs
- All fields properly typed

---

#### Task 33: Define custom message types for async operations

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 32
**Complexity**: Simple
**File**: `pkg/tui/messages.go`

**Description**:

- Define `tasksLoadedMsg` with tasks slice and error
- Define `taskOperationMsg` with operation string and error
- Define `syncCompleteMsg` with status and error
- Use specific types for type-safe handling

**Acceptance Criteria**:

- All custom message types defined
- Well-documented purpose
- Type-safe message handling

---

#### Task 34: Create Lip Gloss style system and color palette

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 32
**Complexity**: Simple
**File**: `pkg/tui/styles.go`

**Description**:

- Define color constants (primary, secondary, pool, today, done, muted, error)
- Define base styles (title, container, box, tab, activeTab, statusBar)
- Define status-specific styles (poolStyle, todayStyle, doneStyle)
- Follow architecture's color palette

**Acceptance Criteria**:

- styles.go with all style definitions
- Colors defined as lipgloss.Color
- Styles composed with lipgloss.NewStyle()
- Commented with usage examples

---

#### Task 35: Implement key bindings system with help.KeyMap

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 32
**Complexity**: Simple
**File**: `pkg/tui/keybindings.go`

**Description**:

- Define keyMap struct with all bindings
- Implement help.KeyMap interface (ShortHelp, FullHelp)
- Follow keybindings from implementation plan
- Group keys logically (navigation, operations, system)

**Test Cases**:

- `TestKeyMap_ShortHelp_ReturnsEssentialKeys`
- `TestKeyMap_FullHelp_ReturnsAllKeys`

**Acceptance Criteria**:

- keyMap implements help.KeyMap
- All keys documented with help text
- Logical grouping

---

#### Task 36: Implement custom list delegate for task rendering

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Tasks 7, 34
**Complexity**: Moderate
**File**: `pkg/tui/delegate.go`

**Description**:

- Implement list.ItemDelegate interface (Height, Spacing, Update, Render)
- Create taskItem wrapper implementing list.Item (FilterValue, Title, Description)
- Render with status-based icons (○ pool, ◉ today, ✓ done)
- Apply status-based styles
- Show tags and due date in description

**Test Cases**:

- `TestTaskItem_FilterValue_ReturnsTitle`
- `TestTaskItem_Description_IncludesTags`
- `TestTaskItem_Description_IncludesDueDate`
- `TestTaskDelegate_Height_ReturnsCorrectValue`

**Acceptance Criteria**:

- taskDelegate and taskItem implemented
- All interface methods implemented
- Visual styling matches design spec

---

#### Task 37: Implement Model Init() method

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Tasks 32-36
**Complexity**: Moderate
**File**: `pkg/tui/model.go`

**Description**:

- Create and configure bubbles/list with custom delegate
- Create bubbles/help component
- Set initial view state to listView
- Set initial tab to poolTab
- Return command to load tasks asynchronously

**Test Cases**:

- `TestModel_Init_ReturnsLoadCommand`
- `TestModel_Init_InitializesComponents`

**Acceptance Criteria**:

- Init() method implemented
- All components initialized
- Returns tea.Cmd to load initial data

---

### TUI Commands (Tasks 38-42)

#### Task 38: Implement loadTasksCmd async command

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Tasks 33, 37
**Complexity**: Simple
**File**: `pkg/tui/commands.go`

**Description**:

- Return tea.Cmd that calls TaskService.ListTasks()
- Handle errors gracefully
- Return tasksLoadedMsg with result or error

**Test Cases**:

- `TestLoadTasksCmd_Success_ReturnsTasksLoadedMsg`
- `TestLoadTasksCmd_Error_ReturnsMessageWithError`

**Acceptance Criteria**:

- loadTasksCmd function implemented
- Returns tea.Cmd
- Error handling robust

---

#### Task 39: Implement pickTaskCmd async command

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 38
**Complexity**: Simple
**File**: `pkg/tui/commands.go`

**Description**:

- Return tea.Cmd that calls TaskService.PickTask()
- Return taskOperationMsg with result

**Test Cases**:

- `TestPickTaskCmd_Success_ReturnsSuccessMsg`
- `TestPickTaskCmd_Error_ReturnsErrorMsg`

**Acceptance Criteria**:

- pickTaskCmd implemented
- All tests passing

---

#### Task 40: Implement deferTaskCmd async command

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 39
**Complexity**: Simple
**File**: `pkg/tui/commands.go`

**Description**:

- Return tea.Cmd that calls TaskService.DeferTask()
- Return taskOperationMsg with result

**Test Cases**:

- `TestDeferTaskCmd_Success_UpdatesList`

**Acceptance Criteria**:

- deferTaskCmd implemented
- All tests passing

---

#### Task 41: Implement completeTaskCmd async command

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 40
**Complexity**: Simple
**File**: `pkg/tui/commands.go`

**Description**:

- Return tea.Cmd that calls TaskService.CompleteTask()
- Return taskOperationMsg with result

**Test Cases**:

- `TestCompleteTaskCmd_Success_UpdatesList`

**Acceptance Criteria**:

- completeTaskCmd implemented
- All tests passing

---

#### Task 42: Implement addTaskCmd and removeTaskCmd async commands

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 41
**Complexity**: Simple
**File**: `pkg/tui/commands.go`

**Description**:

- Implement addTaskCmd calling TaskService.AddTask()
- Implement removeTaskCmd calling TaskService.RemoveTask()
- Both return taskOperationMsg

**Test Cases**:

- `TestAddTaskCmd_ValidInput_Success`
- `TestRemoveTaskCmd_Success_RemovesTask`

**Acceptance Criteria**:

- Both commands implemented
- All tests passing

---

### TUI Update Logic (Tasks 43-48)

#### Task 43: Implement Update() handler for global keys

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Tasks 35, 37
**Complexity**: Moderate
**File**: `pkg/tui/update.go`

**Description**:

- Handle quit keys (q, ctrl+c)
- Handle help toggle (?)
- Handle WindowSizeMsg for responsive layout
- Route messages to view-specific handlers
- Handle tasksLoadedMsg to update model state

**Test Cases**:

- `TestUpdate_QuitKey_ReturnsQuitCmd`
- `TestUpdate_HelpToggle_TogglesHelpState`
- `TestUpdate_WindowSizeMsg_UpdatesDimensions`
- `TestUpdate_TasksLoadedMsg_UpdatesTaskList`

**Acceptance Criteria**:

- Update() method implemented
- Global keys handled correctly
- Message routing works

---

#### Task 44: Implement tab switching logic with filtering

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 43
**Complexity**: Moderate
**File**: `pkg/tui/update.go`

**Description**:

- Handle Tab key to cycle through tabs
- Handle Shift+Tab to cycle backwards
- Implement refreshList() to filter tasks by active tab
- Convert tasks to list.Item slice
- Update list.Model with filtered items

**Test Cases**:

- `TestUpdate_TabKey_SwitchesToNextTab`
- `TestUpdate_ShiftTabKey_SwitchesToPreviousTab`
- `TestRefreshList_PoolTab_ShowsOnlyPoolTasks`
- `TestRefreshList_TodayTab_ShowsOnlyTodayTasks`
- `TestRefreshList_DoneTab_ShowsOnlyDoneTasks`

**Acceptance Criteria**:

- Tab switching works
- refreshList() filters correctly
- All tests passing

---

#### Task 45: Implement task operation handlers (pick, defer, complete)

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Tasks 39-41, 44
**Complexity**: Moderate
**File**: `pkg/tui/update.go`

**Description**:

- Handle Enter key in Pool tab → pick task
- Handle 'd' key in Today tab → defer task
- Handle 'x' key → complete task
- Handle taskOperationMsg result
- Update status bar with success/error
- Refresh list after operation

**Test Cases**:

- `TestUpdate_EnterKey_PoolTab_CallsPickTask`
- `TestUpdate_DKey_TodayTab_CallsDeferTask`
- `TestUpdate_XKey_CallsCompleteTask`
- `TestUpdate_TaskOperationMsg_UpdatesStatusBar`

**Acceptance Criteria**:

- All operation keys handled
- Commands called asynchronously
- Status bar updated
- All tests passing

---

#### Task 46: Implement remove task with confirmation dialog

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Tasks 42, 45
**Complexity**: Moderate
**File**: `pkg/tui/update.go`

**Description**:

- Handle Delete key → enter confirmation state
- Show y/n prompt in status bar
- Track task ID to remove
- Handle 'y' → call removeTaskCmd
- Handle 'n' or Esc → cancel
- Clear confirmation state after operation

**Test Cases**:

- `TestUpdate_DeleteKey_EntersConfirmationState`
- `TestUpdate_ConfirmRemove_YKey_CallsRemoveTask`
- `TestUpdate_ConfirmRemove_NKey_CancelsRemoval`

**Acceptance Criteria**:

- Delete key starts confirmation
- y/n handling correct
- All tests passing

---

#### Task 47: Implement responsive layout handling

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 43
**Complexity**: Moderate
**File**: `pkg/tui/update.go`

**Description**:

- Handle WindowSizeMsg to track dimensions
- Adjust list height based on terminal height
- Adjust list width based on terminal width
- Set compactMode flag for narrow terminals (<80 cols)
- Propagate size to bubbles components

**Test Cases**:

- `TestUpdate_WindowSizeMsg_UpdatesListSize`
- `TestUpdate_NarrowTerminal_SetsCompactMode`
- `TestUpdate_WideTerminal_UnsetCompactMode`

**Acceptance Criteria**:

- WindowSizeMsg handled
- List size adjusts dynamically
- Compact mode implemented

---

#### Task 48: Implement add task view with text inputs

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Tasks 42, 45
**Complexity**: Moderate
**File**: `pkg/tui/update.go`

**Description**:

- Create textinput.Model for title and tags
- Handle 'a' key → switch to addView
- Focus title input on view switch
- Handle Esc → cancel (return to list)
- Handle Enter → save task
- Parse tags (comma-separated)
- Clear inputs after success

**Test Cases**:

- `TestUpdate_AKey_SwitchesToAddView`
- `TestUpdate_AddView_EscKey_ReturnsToListView`
- `TestUpdate_AddView_EnterKey_CallsAddTask`
- `TestUpdate_AddTaskSuccess_ClearsInputs`

**Acceptance Criteria**:

- addView state implemented
- Text inputs work
- Save/cancel work correctly

---

### TUI View Rendering (Tasks 49-54)

#### Task 49: Implement list view rendering with tabs

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Tasks 34, 36, 44
**Complexity**: Moderate
**File**: `pkg/tui/view.go`

**Description**:

- Implement View() method for list view
- Compose layout with lipgloss.JoinVertical
- Sections: header, tabs, list, status bar, help
- Use styles from styles.go
- Render tabs with active tab highlighted

**Test Cases**:

- `TestView_ListViewActive_RendersCorrectly`
- `TestView_ListViewActive_ShowsActiveTab`
- `TestView_ListViewActive_IncludesStatusBar`

**Acceptance Criteria**:

- View() method implemented
- Returns composed layout string
- All sections included
- Visual design matches spec

---

#### Task 50: Implement renderTabs() helper

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 49
**Complexity**: Simple
**File**: `pkg/tui/view.go`

**Description**:

- Render three tabs: Pool, Today, Done
- Apply activeTabStyle to active tab
- Apply tabStyle to inactive tabs
- Use lipgloss.JoinHorizontal
- Show task counts in each tab

**Test Cases**:

- `TestRenderTabs_PoolActive_HighlightsPool`
- `TestRenderTabs_ShowsTaskCounts`

**Acceptance Criteria**:

- renderTabs() method implemented
- Active tab visually distinct
- All tests passing

---

#### Task 51: Implement status bar rendering

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 49
**Complexity**: Simple
**File**: `pkg/tui/view.go`

**Description**:

- Show counts: "Pool: 5 | Today: 3 | Done: 12"
- Show status message (last operation result)
- Style with muted color for counts
- Success/error colors for messages

**Test Cases**:

- `TestRenderStatusBar_ShowsCounts`
- `TestRenderStatusBar_ShowsMessage`
- `TestRenderStatusBar_SuccessMessage_GreenColor`
- `TestRenderStatusBar_ErrorMessage_RedColor`

**Acceptance Criteria**:

- renderStatusBar() implemented
- Shows counts and messages
- Styled correctly

---

#### Task 52: Implement add task view rendering

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Tasks 48, 49
**Complexity**: Moderate
**File**: `pkg/tui/view.go`

**Description**:

- Implement viewAdd() for add task view
- Show title input and tags input
- Show instructions
- Use Lip Gloss boxes and styling
- Esc/Enter hints

**Test Cases**:

- `TestView_AddView_ShowsInputs`
- `TestView_AddView_ShowsInstructions`

**Acceptance Criteria**:

- viewAdd() method implemented
- All elements shown
- Visual design matches spec

---

#### Task 53: Integrate help component

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Tasks 35, 49
**Complexity**: Simple
**File**: `pkg/tui/view.go`

**Description**:

- Create help.Model in initialization
- Toggle full help with '?' key
- Render help in footer (short help by default)
- Use help.View() with keyMap

**Test Cases**:

- `TestUpdate_QuestionMarkKey_TogglesHelp`
- `TestView_ShowsShortHelp_ByDefault`
- `TestView_ShowsFullHelp_WhenToggled`

**Acceptance Criteria**:

- help.Model integrated
- Toggle works
- Help shown in footer

---

#### Task 54: Write comprehensive TUI unit tests

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Tasks 32-53
**Complexity**: Complex
**File**: `pkg/tui/*_test.go`

**Description**:

- Test Update() handlers with mock service
- Test View() output contains expected strings
- Test state transitions
- Test all key handlers
- Use table-driven tests for key sequences

**Test Cases**:

- All test cases from previous TUI tasks
- Additional edge cases

**Acceptance Criteria**:

- All TUI tests passing
- Coverage >80% for pkg/tui
- Mock service used for isolation

---

## Phase 3: Sync Capabilities & CLI Commands (Tasks 55-71)

**Goal**: Implement Git-based sync and complete CLI interface
**Estimated Duration**: 1 week
**Key Deliverables**: Working sync infrastructure, full CLI command suite

### Sync Infrastructure (Tasks 55-63)

#### Task 55: Define SyncAdapter interface

**Status**: ⏸️ Pending
**Layer**: Infrastructure (Interface)
**Dependencies**: Task 1
**Complexity**: Simple
**File**: `internal/sync/adapter.go`

**Description**:

- Define interface: Status(), Pull(), Push()
- Define SyncStatus enum (synced, ahead, behind, diverged)
- Define PullResult struct (conflict bool, message string)
- Document contract: works with encrypted blobs only

**Acceptance Criteria**:

- SyncAdapter interface defined
- SyncStatus and PullResult types defined
- Comprehensive documentation

---

#### Task 56: Implement GitSync adapter - Status check

**Status**: ⏸️ Pending
**Layer**: Infrastructure
**Dependencies**: Task 55
**Complexity**: Moderate
**File**: `internal/sync/git_sync.go`

**Description**:

- Use exec.CommandContext for `git status --porcelain` and `git rev-list`
- Check if local has unpushed commits (ahead)
- Check if remote has unpulled commits (behind)
- Check for divergence
- Return SyncStatus enum

**Test Cases**:

- `TestGitSync_Status_Synced`
- `TestGitSync_Status_Ahead`
- `TestGitSync_Status_Behind`
- `TestGitSync_Status_Diverged`
- `TestGitSync_Status_NoGitRepo_ReturnsError`
- Use temp git repos for testing

**Acceptance Criteria**:

- Status() method implemented
- All sync states detected
- All tests passing

---

#### Task 57: Implement GitSync adapter - Push operation

**Status**: ⏸️ Pending
**Layer**: Infrastructure
**Dependencies**: Task 56
**Complexity**: Moderate
**File**: `internal/sync/git_sync.go`

**Description**:

- Stage with `git add todo.json.age`
- Commit with `git commit -m "Update tasks"`
- Push with `git push origin main`
- Handle errors: no changes, push rejected, network failure
- Idempotent (no changes is success)

**Test Cases**:

- `TestGitSync_Push_Success`
- `TestGitSync_Push_NoChanges_Success`
- `TestGitSync_Push_RemoteDiverged_ReturnsError`
- `TestGitSync_Push_NetworkFailure_ReturnsError`

**Acceptance Criteria**:

- Push() method implemented
- Git commands executed correctly
- All tests passing

---

#### Task 58: Implement GitSync adapter - Pull operation

**Status**: ⏸️ Pending
**Layer**: Infrastructure
**Dependencies**: Task 57
**Complexity**: Complex
**File**: `internal/sync/git_sync.go`

**Description**:

- Run `git pull origin main`
- Handle merge conflicts gracefully
- If conflict: return PullResult with conflict=true
- If success: return PullResult with conflict=false
- Detect if merge is required

**Test Cases**:

- `TestGitSync_Pull_Success_NoConflict`
- `TestGitSync_Pull_AlreadyUpToDate`
- `TestGitSync_Pull_Conflict_ReturnsConflictResult`
- `TestGitSync_Pull_NetworkFailure_ReturnsError`

**Acceptance Criteria**:

- Pull() method implemented
- Handles fast-forward merges
- Detects conflicts
- All tests passing

---

#### Task 59: Implement conflict resolution strategy (last-write-wins)

**Status**: ⏸️ Pending
**Layer**: Infrastructure
**Dependencies**: Task 58
**Complexity**: Moderate
**File**: `internal/sync/git_sync.go`

**Description**:

- On conflict, prefer remote version
- Store conflicted local version in backup file
- Return PullResult indicating conflict was resolved
- Document behavior clearly

**Test Cases**:

- `TestGitSync_Pull_ConflictResolution_PrefersRemote`
- `TestGitSync_Pull_ConflictResolution_BacksUpLocal`
- `TestGitSync_Pull_ConflictResolution_ReturnsMessage`

**Acceptance Criteria**:

- Conflict resolution implemented
- Local backup created
- User notified
- All tests passing

---

#### Task 60: Implement syncCmd for TUI

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Tasks 55-59
**Complexity**: Moderate
**File**: `pkg/tui/commands.go`

**Description**:

- Implement syncCmd (async) that calls SyncAdapter
- Return syncCompleteMsg with result
- Handle both push and pull operations

**Test Cases**:

- `TestSyncCmd_Success_ReturnsCompleteMsg`
- `TestSyncCmd_Error_ReturnsErrorMsg`

**Acceptance Criteria**:

- syncCmd implemented
- All tests passing

---

#### Task 61: Integrate sync into TUI with spinner

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 60
**Complexity**: Moderate
**File**: `pkg/tui/update.go`, `pkg/tui/view.go`

**Description**:

- Handle 's' key to trigger sync
- Show spinner during sync operation
- Handle syncCompleteMsg with result
- Display status: "Syncing...", "Sync complete", "Sync failed: error"
- Handle conflict message if PullResult.Conflict==true

**Test Cases**:

- `TestUpdate_SKey_TriggersSyncCmd`
- `TestUpdate_SyncCompleteMsg_UpdatesStatusBar`
- `TestUpdate_SyncCompleteMsg_StopsSpinner`

**Acceptance Criteria**:

- 's' key handled
- Spinner shown during operation
- Status bar updated with result

---

#### Task 62: Write sync layer unit tests

**Status**: ⏸️ Pending
**Layer**: Infrastructure
**Dependencies**: Tasks 56-59
**Complexity**: Moderate
**File**: `internal/sync/*_test.go`

**Description**:

- Ensure all sync methods have comprehensive tests
- Use temp git repos
- Test all sync states and operations
- Test conflict resolution

**Acceptance Criteria**:

- All sync tests passing
- Coverage >85% for internal/sync package

---

#### Task 63: Write sync integration tests

**Status**: ⏸️ Pending
**Layer**: Integration Test
**Dependencies**: Task 62
**Complexity**: Complex
**File**: `internal/sync/integration_test.go`

**Description**:

- Test full sync cycle: push → pull → conflict resolution
- Use real git repos with remotes
- Test divergence scenarios
- Verify data integrity

**Test Cases**:

- `TestIntegration_Sync_PushPull_Workflow`
- `TestIntegration_Sync_Conflict_Resolution`
- `TestIntegration_Sync_Divergence_Handling`

**Acceptance Criteria**:

- Integration tests covering sync workflows
- Tests use real git repos
- All tests passing

---

### CLI Commands (Tasks 64-71)

#### Task 64: Implement main entry point with CLI dispatcher

**Status**: ⏸️ Pending
**Layer**: Presentation (CLI)
**Dependencies**: Tasks 21, 32
**Complexity**: Moderate
**File**: `cmd/togo/main.go`

**Description**:

- No args: launch TUI
- Args: parse command and route to handlers
- Use switch statement for routing
- Commands: add, list, pick, defer, complete, remove, report, sync, help
- Initialize dependencies (storage, service) before routing

**Acceptance Criteria**:

- main.go with dispatcher logic
- TUI launched when no args
- Command routing works
- Unknown commands handled gracefully

---

#### Task 65: Implement CLI command - add

**Status**: ⏸️ Pending
**Layer**: Presentation (CLI)
**Dependencies**: Task 64
**Complexity**: Simple
**File**: `cmd/togo/commands/add.go`

**Description**:

- Flag parsing: -n (notes), -t (tags)
- First positional arg is title (required)
- Call TaskService.AddTask()
- Print success message with Lip Gloss styling

**Test Cases**:

- `TestHandleAdd_ValidInput_Success`
- `TestHandleAdd_NoTitle_ReturnsError`
- `TestHandleAdd_WithTags_Success`

**Acceptance Criteria**:

- handleAdd() function implemented
- Flag parsing works
- Styled output with Lip Gloss

---

#### Task 66: Implement CLI command - list

**Status**: ⏸️ Pending
**Layer**: Presentation (CLI)
**Dependencies**: Task 65
**Complexity**: Moderate
**File**: `cmd/togo/commands/list.go`

**Description**:

- Flags: --pool, --today, --done, --tag (filter)
- Build TaskFilter from flags
- Call TaskService.ListTasks()
- Format output as table with Lip Gloss

**Test Cases**:

- `TestHandleList_NoFilter_ListsAll`
- `TestHandleList_FilterByStatus`
- `TestHandleList_FilterByTag`

**Acceptance Criteria**:

- handleList() function implemented
- Filtering works
- Output formatted nicely

---

#### Task 67: Implement CLI commands - pick, defer, complete

**Status**: ⏸️ Pending
**Layer**: Presentation (CLI)
**Dependencies**: Task 66
**Complexity**: Simple
**File**: `cmd/togo/commands/{pick,defer,complete}.go`

**Description**:

- Parse task ID from first positional arg
- Call respective TaskService method
- Print success/error message

**Test Cases**:

- `TestHandlePick_ValidID_Success`
- `TestHandleDefer_ValidID_Success`
- `TestHandleComplete_ValidID_Success`
- Error cases for each

**Acceptance Criteria**:

- All three commands implemented
- All tests passing

---

#### Task 68: Implement CLI command - remove

**Status**: ⏸️ Pending
**Layer**: Presentation (CLI)
**Dependencies**: Task 67
**Complexity**: Simple
**File**: `cmd/togo/commands/remove.go`

**Description**:

- Parse ID, call RemoveTask()
- Print warning about permanent deletion
- No confirmation in CLI

**Test Cases**:

- `TestHandleRemove_ValidID_Success`
- `TestHandleRemove_TaskNotFound_ReturnsError`

**Acceptance Criteria**:

- handleRemove() implemented
- Warns user
- All tests passing

---

#### Task 69: Implement CLI command - sync

**Status**: ⏸️ Pending
**Layer**: Presentation (CLI)
**Dependencies**: Tasks 55-59, 68
**Complexity**: Moderate
**File**: `cmd/togo/commands/sync.go`

**Description**:

- Subcommands: push, pull, status
- Call appropriate SyncAdapter method
- Print status results clearly
- Handle conflicts in pull

**Test Cases**:

- `TestHandleSync_Push_Success`
- `TestHandleSync_Pull_Success`
- `TestHandleSync_Status_PrintsStatus`
- `TestHandleSync_Pull_Conflict_NotifiesUser`

**Acceptance Criteria**:

- handleSync() implemented
- All subcommands work
- All tests passing

---

#### Task 70: Implement CLI command - help

**Status**: ⏸️ Pending
**Layer**: Presentation (CLI)
**Dependencies**: Task 64
**Complexity**: Simple
**File**: `cmd/togo/commands/help.go`

**Description**:

- Print usage information for all commands
- Include examples
- Style with Lip Gloss
- Document TUI keybindings

**Test Cases**:

- `TestHandleHelp_PrintsUsage`
- `TestHandleHelp_IncludesAllCommands`

**Acceptance Criteria**:

- handleHelp() implemented
- Comprehensive help text
- Nicely formatted

---

#### Task 71: Write CLI integration tests

**Status**: ⏸️ Pending
**Layer**: Integration Test
**Dependencies**: Tasks 64-70
**Complexity**: Moderate
**File**: `cmd/togo/integration_test.go`

**Description**:

- Use exec.Command to run compiled binary
- Test workflows: add → list → pick → complete
- Verify output contains expected strings
- Use temp directory for data storage

**Test Cases**:

- `TestIntegration_CLI_AddAndList`
- `TestIntegration_CLI_PickTask_Workflow`
- `TestIntegration_CLI_Sync_Workflow`

**Acceptance Criteria**:

- Integration tests using compiled binary
- Tests cover main workflows
- All tests passing

---

## Phase 4: Polish & Advanced Features (Tasks 72-87)

**Goal**: Add reports, advanced features, polish, and comprehensive documentation
**Estimated Duration**: 1 week
**Key Deliverables**: Report generation, refined UX, comprehensive docs

### Report Generation (Tasks 72-76)

#### Task 72: Implement ReportService core logic

**Status**: ⏸️ Pending
**Layer**: Service
**Dependencies**: Task 21
**Complexity**: Moderate
**File**: `internal/service/report_service.go`

**Description**:

- Define Report struct with sections and metadata
- Define ReportOptions (since, until, format)
- Load all tasks
- Filter by CompletedAt in date range
- Group by date (daily breakdown)
- Calculate statistics

**Test Cases**:

- `TestReportService_Generate_FiltersByDateRange`
- `TestReportService_Generate_GroupsByDate`
- `TestReportService_Generate_CalculatesStatistics`

**Acceptance Criteria**:

- ReportService implemented
- GenerateReport() method works
- All tests passing

---

#### Task 73: Implement text report formatter

**Status**: ⏸️ Pending
**Layer**: Service
**Dependencies**: Task 72
**Complexity**: Simple
**File**: `internal/service/formatters/text_formatter.go`

**Description**:

- Implement Formatter interface
- Format report as readable text
- Sections: summary, daily breakdown, statistics
- Use Lip Gloss for styling

**Test Cases**:

- `TestTextFormatter_Format_IncludesAllSections`
- `TestTextFormatter_Format_FormatsReadably`

**Acceptance Criteria**:

- TextFormatter implements Formatter
- Format() returns text string
- Output is human-readable

---

#### Task 74: Implement JSON report formatter

**Status**: ⏸️ Pending
**Layer**: Service
**Dependencies**: Task 72
**Complexity**: Simple
**File**: `internal/service/formatters/json_formatter.go`

**Description**:

- Implement Formatter interface
- Marshal Report struct to JSON
- Pretty-print with indentation

**Test Cases**:

- `TestJSONFormatter_Format_ValidJSON`
- `TestJSONFormatter_Format_IncludesAllData`

**Acceptance Criteria**:

- JSONFormatter implements Formatter
- Format() returns valid JSON
- All tests passing

---

#### Task 75: Implement CLI command - report

**Status**: ⏸️ Pending
**Layer**: Presentation (CLI)
**Dependencies**: Tasks 73, 74
**Complexity**: Moderate
**File**: `cmd/togo/commands/report.go`

**Description**:

- Flags: --since (duration), --until (date), --format (text|json)
- Parse duration/date strings
- Call ReportService.GenerateReport()
- Choose formatter based on flag
- Print formatted report

**Test Cases**:

- `TestHandleReport_DefaultOptions_Success`
- `TestHandleReport_CustomDateRange_FiltersCorrectly`
- `TestHandleReport_JSONFormat_OutputsJSON`

**Acceptance Criteria**:

- handleReport() implemented
- Flag parsing works
- Both formatters work

---

#### Task 76: Implement TUI report view

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Tasks 73, 54
**Complexity**: Moderate
**File**: `pkg/tui/view.go`, `pkg/tui/update.go`

**Description**:

- Add reportView to view state enum
- Handle 'r' key to switch to report view
- Async command to generate report
- Display formatted report text
- Allow scrolling with viewport
- Esc key to return to list view

**Test Cases**:

- `TestUpdate_RKey_SwitchesToReportView`
- `TestReportView_DisplaysReport`
- `TestReportView_EscKey_ReturnsToListView`

**Acceptance Criteria**:

- reportView state implemented
- 'r' key generates and displays report
- Scrolling works

---

### Advanced TUI Features (Tasks 77-80)

#### Task 77: Implement task detail view - display

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 54
**Complexity**: Moderate
**File**: `pkg/tui/view.go`, `pkg/tui/update.go`

**Description**:

- Add detailView to view state enum
- Handle Enter key in Today/Done tabs to view details
- Display all task fields with Lip Gloss boxes
- Show available actions
- Esc key to return to list

**Test Cases**:

- `TestUpdate_EnterKey_TodayTab_ShowsDetails`
- `TestDetailView_DisplaysAllFields`
- `TestDetailView_EscKey_ReturnsToListView`

**Acceptance Criteria**:

- detailView state implemented
- viewDetail() method renders task info
- All fields displayed

---

#### Task 78: Implement task detail view - edit notes

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 77
**Complexity**: Moderate
**File**: `pkg/tui/update.go`

**Description**:

- Handle 'e' key in detail view to edit notes
- Show textarea.Model with current notes
- Handle Esc to cancel editing
- Handle Ctrl+S to save changes
- Update task via service

**Test Cases**:

- `TestDetailView_EKey_EntersEditMode`
- `TestDetailView_EditMode_EscKey_CancelsEdit`
- `TestDetailView_EditMode_SaveKey_UpdatesNotes`

**Acceptance Criteria**:

- Notes editing implemented
- Textarea component integrated
- Save/cancel work correctly

---

#### Task 79: Implement task filtering with '/' key

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Task 54
**Complexity**: Moderate
**File**: `pkg/tui/update.go`

**Description**:

- Handle '/' key to enter filter mode
- Show textinput for filter query
- Filter tasks by title (substring) and tags
- Update list dynamically as user types
- Esc to clear filter
- Show filter indicator in status bar

**Test Cases**:

- `TestUpdate_SlashKey_EntersFilterMode`
- `TestFilterMode_FiltersTasksByTitle`
- `TestFilterMode_FiltersTasksByTag`
- `TestFilterMode_EscKey_ClearsFilter`

**Acceptance Criteria**:

- Filter mode implemented
- Real-time filtering works
- Filter indicator shown

---

#### Task 80: Polish confirmation dialogs and error messages

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: Tasks 46, 77-79
**Complexity**: Simple
**File**: `pkg/tui/view.go`

**Description**:

- Review all confirmation dialogs
- Ensure consistent styling
- Improve error message clarity
- Add context to errors

**Acceptance Criteria**:

- All dialogs reviewed
- Consistent styling
- Clear error messages

---

### Documentation & Final Polish (Tasks 81-87)

#### Task 81: Create cron script for weekly reports

**Status**: ⏸️ Pending
**Layer**: Documentation/Scripts
**Dependencies**: Task 75
**Complexity**: Simple
**File**: `scripts/cron-report.sh`

**Description**:

- Bash script that runs `todo report --since 7d`
- Includes crontab example for Monday 8am
- Document how to customize

**Acceptance Criteria**:

- Script created and executable
- Crontab example included
- Documented in README

---

#### Task 82: Write comprehensive README documentation

**Status**: ⏸️ Pending
**Layer**: Documentation
**Dependencies**: All features (Tasks 1-81)
**Complexity**: Moderate
**File**: `README.md`

**Description**:

- Installation instructions
- Quick start guide
- TUI keybindings reference
- CLI command reference with examples
- Encryption setup instructions
- Sync setup instructions
- Weekly report automation guide
- Troubleshooting section

**Acceptance Criteria**:

- README.md with all sections
- Examples for all commands
- Clear, beginner-friendly
- Includes architecture reference

---

#### Task 83: Performance testing and benchmarks

**Status**: ⏸️ Pending
**Layer**: Testing
**Dependencies**: All core features
**Complexity**: Moderate
**File**: `internal/service/benchmark_test.go`

**Description**:

- Benchmark report generation with large datasets
- Benchmark encryption/decryption
- Benchmark task filtering
- Use Go's testing.B framework
- Document performance characteristics

**Test Cases**:

- `BenchmarkGenerateReport_1000Tasks`
- `BenchmarkEncryptDecrypt_LargePayload`
- `BenchmarkTaskFilter_LargeCollection`

**Acceptance Criteria**:

- Benchmarks written and passing
- Performance documented
- No critical bottlenecks

---

#### Task 84: Final visual design polish

**Status**: ⏸️ Pending
**Layer**: Presentation (TUI)
**Dependencies**: All TUI features
**Complexity**: Simple

**Description**:

- Review color consistency across views
- Ensure borders and spacing are consistent
- Test in both light and dark terminals
- Ensure high contrast for accessibility
- Polish animations and transitions

**Acceptance Criteria**:

- Visual design review complete
- Consistency across all views
- Works well in various terminals
- Accessibility considerations met

---

#### Task 85: End-to-end testing

**Status**: ⏸️ Pending
**Layer**: Integration Test
**Dependencies**: All features
**Complexity**: Complex
**File**: `test/e2e_test.go`

**Description**:

- Test complete user workflows
- Test sync cycle: push → pull → conflict resolution
- Test TUI and CLI interoperability
- Use real components
- Verify data integrity

**Test Cases**:

- `TestE2E_CompleteWorkflow_TUI`
- `TestE2E_CompleteWorkflow_CLI`
- `TestE2E_SyncCycle_PushPull`
- `TestE2E_ConflictResolution`
- `TestE2E_TUI_CLI_Interop`

**Acceptance Criteria**:

- E2E tests covering all major workflows
- Tests use real components
- All tests passing

---

#### Task 86: Final code quality review

**Status**: ⏸️ Pending
**Layer**: Quality Assurance
**Dependencies**: Tasks 1-85
**Complexity**: Simple

**Description**:

- Run `go fmt` on all code
- Run `go vet` and fix issues
- Run `golangci-lint` and fix issues
- Ensure all tests passing
- Verify coverage >85% overall
- Review and improve error messages

**Acceptance Criteria**:

- All code formatted
- No vet or lint errors
- All tests passing
- Coverage >85%
- Error messages user-friendly

---

#### Task 87: Project completion verification

**Status**: ⏸️ Pending
**Layer**: Quality Assurance
**Dependencies**: Tasks 1-86
**Complexity**: Simple

**Description**:

- Verify all 87 tasks completed
- Run full test suite
- Build binary and test manually
- Verify CI pipeline passes
- Verify all documentation complete
- Tag v1.0.0 release

**Acceptance Criteria**:

- All tasks completed
- All tests passing
- Binary builds and works
- CI green
- Documentation complete
- v1.0.0 tagged

---

## Summary

**Total Tasks**: 87
**Phase 1 (Foundation)**: 31 tasks
**Phase 2 (TUI)**: 23 tasks
**Phase 3 (Sync & CLI)**: 17 tasks
**Phase 4 (Polish)**: 16 tasks

**Estimated Timeline**: 4 weeks (1 week per phase)

**Current Status**: All tasks pending, ready to begin Phase 1

---

## Next Steps

1. Begin with **Task 1**: Create project directory structure
2. Work through Phase 1 sequentially
3. Mark tasks as "in_progress" when starting
4. Mark tasks as "completed" when acceptance criteria met
5. Move to Phase 2 only after Phase 1 complete

---

*Last Updated: 2025-11-09*
