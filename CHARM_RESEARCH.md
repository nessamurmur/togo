# Replacing Cobra with Charm Libraries: Comprehensive Research

**Author**: Charm Expert Agent
**Date**: 2025-11-09
**Purpose**: Guide for replacing Cobra with Bubble Tea, Bubbles, and Lip Gloss for togo CLI/TUI application

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Architectural Comparison](#architectural-comparison)
3. [Command-Line Argument Handling](#command-line-argument-handling)
4. [Multi-View Application Structure](#multi-view-application-structure)
5. [Bubbles Components Reference](#bubbles-components-reference)
6. [Lip Gloss Styling Patterns](#lip-gloss-styling-patterns)
7. [Complete Implementation Example](#complete-implementation-example)
8. [Migration Strategy](#migration-strategy)
9. [Best Practices](#best-practices)
10. [Trade-offs and Considerations](#trade-offs-and-considerations)

---

## Executive Summary

### Key Findings

**Bubble Tea does NOT replace Cobra's command-line parsing** - they serve different purposes:

- **Cobra**: Command routing, flag parsing, subcommand hierarchy (`todo add`, `todo list`)
- **Bubble Tea**: Interactive TUI framework (full-screen terminal applications)

### Recommended Hybrid Approach

For the togo application, you have two architectural options:

#### Option 1: Cobra + Bubble Tea (Recommended)
```
CLI Layer: Cobra handles argument parsing and routing
             ↓
Mode Decision: --tui flag OR tui subcommand
             ↓
  ┌──────────┴──────────┐
  ↓                     ↓
CLI Mode            TUI Mode
(Quick ops)     (Interactive Bubble Tea)
```

**Best for**: Professional CLI tools with both quick commands and interactive modes

#### Option 2: Pure Bubble Tea (Alternative)
```
main.go: Parse args with flag/pflag package
             ↓
If args exist → Execute command directly (non-interactive)
If no args → Launch Bubble Tea TUI
```

**Best for**: Primarily TUI applications with optional CLI shortcuts

### What You Actually Replace

To eliminate Cobra completely, you need:

1. **Standard library `flag` or `pflag`** - Parse command-line arguments
2. **Manual command dispatcher** - Route commands to handlers
3. **Bubble Tea** - Handle interactive TUI mode
4. **Bubbles components** - Build TUI interface (lists, inputs, etc.)
5. **Lip Gloss** - Style terminal output beautifully

---

## Architectural Comparison

### Traditional Cobra Architecture

```
$ todo add "Buy milk" --tag shopping

                    main.go
                       ↓
                  Cobra Root Command
                       ↓
        ┌──────────────┼──────────────┐
        ↓              ↓              ↓
    AddCommand    ListCommand    PickCommand
        ↓              ↓              ↓
            TaskService (shared)
```

**Cobra provides:**
- Automatic help generation (`todo --help`)
- Subcommand routing (`todo add`, `todo list`)
- Flag parsing and validation (`--tag`, `--format`)
- Shell completion
- Error handling

### Pure Bubble Tea Architecture

```
$ todo                    $ todo add "Buy milk"

    main.go                       main.go
       ↓                             ↓
Parse args (empty)           Parse args (has command)
       ↓                             ↓
Launch Bubble Tea            Execute command directly
       ↓                             ↓
┌──────────────┐              Print result and exit
│  TUI Model   │
│  ┌────────┐  │
│  │ Views  │  │
│  ├────────┤  │
│  │ Pool   │  │
│  │ Today  │  │
│  │ Done   │  │
│  └────────┘  │
└──────────────┘
```

**Pure Bubble Tea provides:**
- Interactive multi-view interface
- Real-time updates and animations
- Keyboard-driven navigation
- Beautiful styled output
- Component-based UI (lists, inputs, etc.)

### Hybrid Architecture (Recommended for togo)

```
$ todo add "Task"              $ todo tui / $ todo

     main.go                        main.go
        ↓                              ↓
   Cobra Parser                   Cobra Parser
        ↓                              ↓
   AddCommand                      TuiCommand
        ↓                              ↓
   TaskService                    Bubble Tea App
        ↓                              ↓
   Print & Exit              ┌─────────────────┐
                            │   TUI Model     │
                            │  ┌───────────┐  │
                            │  │Task List  │  │
                            │  │Add Form   │  │
                            │  │Help View  │  │
                            │  └───────────┘  │
                            └─────────────────┘
```

---

## Command-Line Argument Handling

### Without Cobra: Standard Library Approach

```go
// File: cmd/todo/main.go
package main

import (
    "flag"
    "fmt"
    "os"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
    "togo/internal/service"
    "togo/pkg/tui"
)

func main() {
    // If no arguments, launch TUI
    if len(os.Args) == 1 {
        launchTUI()
        return
    }

    // Parse command
    command := os.Args[1]

    switch command {
    case "add":
        handleAdd(os.Args[2:])
    case "list":
        handleList(os.Args[2:])
    case "pick":
        handlePick(os.Args[2:])
    case "defer":
        handleDefer(os.Args[2:])
    case "complete":
        handleComplete(os.Args[2:])
    case "remove":
        handleRemove(os.Args[2:])
    case "report":
        handleReport(os.Args[2:])
    case "sync":
        handleSync(os.Args[2:])
    case "tui":
        launchTUI()
    case "help", "--help", "-h":
        printHelp()
    default:
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
        fmt.Fprintf(os.Stderr, "Run 'todo help' for usage.\n")
        os.Exit(1)
    }
}

// handleAdd implements: todo add "title" [-n notes] [-t tag1,tag2]
func handleAdd(args []string) {
    fs := flag.NewFlagSet("add", flag.ExitOnError)
    notes := fs.String("n", "", "Task notes")
    tags := fs.String("t", "", "Comma-separated tags")

    if err := fs.Parse(args); err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
        os.Exit(1)
    }

    if fs.NArg() == 0 {
        fmt.Fprintf(os.Stderr, "Usage: todo add \"title\" [-n notes] [-t tags]\n")
        os.Exit(1)
    }

    title := fs.Arg(0)

    // Parse tags
    var tagList []string
    if *tags != "" {
        tagList = strings.Split(*tags, ",")
        for i := range tagList {
            tagList[i] = strings.TrimSpace(tagList[i])
        }
    }

    // Create request
    req := service.AddTaskRequest{
        Title: title,
        Notes: *notes,
        Tags:  tagList,
    }

    // Execute command
    svc := initService()
    task, err := svc.AddTask(context.Background(), req)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error adding task: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("✓ Task added: %s (ID: %s)\n", task.Title, task.ID)
}

// handleList implements: todo list [--pool|--today|--done] [--tag x]
func handleList(args []string) {
    fs := flag.NewFlagSet("list", flag.ExitOnError)
    pool := fs.Bool("pool", false, "Show pool tasks")
    today := fs.Bool("today", false, "Show today tasks")
    done := fs.Bool("done", false, "Show completed tasks")
    tag := fs.String("tag", "", "Filter by tag")

    if err := fs.Parse(args); err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
        os.Exit(1)
    }

    // Build filter
    filter := model.TaskFilter{
        Tags: []string{},
    }

    if *tag != "" {
        filter.Tags = append(filter.Tags, *tag)
    }

    // Determine status filter
    if *pool {
        status := model.StatusPool
        filter.Status = &status
    } else if *today {
        status := model.StatusToday
        filter.Status = &status
    } else if *done {
        status := model.StatusDone
        filter.Status = &status
    }

    // Execute query
    svc := initService()
    tasks, err := svc.ListTasks(context.Background(), filter)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error listing tasks: %v\n", err)
        os.Exit(1)
    }

    // Format output with Lip Gloss
    printTaskList(tasks)
}

func launchTUI() {
    svc := initService()
    model := tui.NewModel(svc)

    p := tea.NewProgram(
        model,
        tea.WithAltScreen(),       // Use alternate screen buffer
        tea.WithMouseCellMotion(), // Enable mouse support
    )

    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
        os.Exit(1)
    }
}

func printHelp() {
    help := `togo - A bullet journal-inspired task manager

USAGE:
    todo                              Launch interactive TUI
    todo <command> [arguments]        Execute command directly

COMMANDS:
    add "title" [-n notes] [-t tags]  Add a new task to pool
    list [--pool|--today|--done]      List tasks
    pick <id>                         Pick task for today
    defer <id>                        Defer task back to pool
    complete <id>                     Mark task as done
    remove <id>                       Delete task
    report [--since 7d]               Generate completion report
    sync <push|pull|status>           Sync with remote
    tui                               Launch interactive TUI
    help                              Show this help

FLAGS:
    -h, --help                        Show help for any command

EXAMPLES:
    todo add "Buy groceries" -t shopping
    todo list --today
    todo pick abc123
    todo report --since 7d

For interactive mode, run: todo
`
    fmt.Print(help)
}

func initService() *service.TaskService {
    // Initialize storage, encryption, service
    // (Implementation details omitted for brevity)
    return nil
}
```

### Alternative: Using pflag (More Cobra-like)

```go
import (
    "github.com/spf13/pflag"
)

func handleAdd(args []string) {
    fs := pflag.NewFlagSet("add", pflag.ExitOnError)
    notes := fs.StringP("notes", "n", "", "Task notes")
    tags := fs.StringSliceP("tags", "t", []string{}, "Tags")
    dueDate := fs.StringP("due", "d", "", "Due date (YYYY-MM-DD)")

    fs.Parse(args)

    if fs.NArg() == 0 {
        fmt.Fprintf(os.Stderr, "Usage: todo add \"title\" [flags]\n")
        fs.PrintDefaults()
        os.Exit(1)
    }

    // ... rest of implementation
}
```

**Advantages of pflag:**
- POSIX-style flag parsing (--flag vs -flag)
- Short flags (-n) and long flags (--notes)
- Better compatibility if migrating from Cobra
- More familiar to users of standard CLI tools

---

## Multi-View Application Structure

### Core Pattern: State-Based View Switching

Bubble Tea applications with multiple views use a **state machine pattern**:

1. Main model contains a **state enum** (current view)
2. **Update** function routes messages based on state
3. **View** function renders based on state
4. Each "view" can be a separate sub-model or just rendering logic

### Pattern 1: Simple Boolean State

**Best for**: 2-3 views

```go
// File: pkg/tui/model.go
package tui

import (
    tea "github.com/charmbracelet/bubbletea"
)

type viewState int

const (
    viewList viewState = iota
    viewAdd
    viewHelp
)

type Model struct {
    // Current view
    currentView viewState

    // Shared state
    service *service.TaskService
    tasks   []*model.Task

    // View-specific state
    cursor       int           // For list view
    addInput     textinput.Model // For add view
    statusMsg    string

    width  int
    height int
}

func NewModel(svc *service.TaskService) Model {
    ti := textinput.New()
    ti.Placeholder = "Enter task title..."
    ti.Focus()

    return Model{
        currentView: viewList,
        service:     svc,
        tasks:       []model.Task{},
        addInput:    ti,
    }
}

func (m Model) Init() tea.Cmd {
    return loadTasks(m.service)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.KeyMsg:
        // Global keybindings (work in all views)
        switch msg.String() {
        case "ctrl+c", "q":
            if m.currentView != viewAdd {
                return m, tea.Quit
            }
        case "?":
            m.currentView = viewHelp
            return m, nil
        case "esc":
            if m.currentView == viewHelp {
                m.currentView = viewList
                return m, nil
            }
        }

        // Route to view-specific handler
        switch m.currentView {
        case viewList:
            return m.updateList(msg)
        case viewAdd:
            return m.updateAdd(msg)
        case viewHelp:
            return m.updateHelp(msg)
        }

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        return m, nil

    case tasksLoadedMsg:
        m.tasks = msg.tasks
        return m, nil

    case taskAddedMsg:
        if msg.err != nil {
            m.statusMsg = fmt.Sprintf("Error: %v", msg.err)
        } else {
            m.statusMsg = "Task added!"
            m.currentView = viewList
            m.addInput.SetValue("")
        }
        return m, loadTasks(m.service)
    }

    return m, nil
}

// View routing
func (m Model) View() string {
    switch m.currentView {
    case viewList:
        return m.viewList()
    case viewAdd:
        return m.viewAdd()
    case viewHelp:
        return m.viewHelp()
    default:
        return "Unknown view"
    }
}

// List view update handler
func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "up", "k":
        if m.cursor > 0 {
            m.cursor--
        }
    case "down", "j":
        if m.cursor < len(m.tasks)-1 {
            m.cursor++
        }
    case "a":
        m.currentView = viewAdd
        return m, textinput.Blink
    case "enter":
        // Pick task
        if len(m.tasks) > 0 {
            return m, m.pickTask(m.tasks[m.cursor].ID)
        }
    case "d":
        // Defer task
        if len(m.tasks) > 0 {
            return m, m.deferTask(m.tasks[m.cursor].ID)
        }
    case "x":
        // Complete task
        if len(m.tasks) > 0 {
            return m, m.completeTask(m.tasks[m.cursor].ID)
        }
    }
    return m, nil
}

// Add view update handler
func (m Model) updateAdd(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "enter":
        title := m.addInput.Value()
        if title != "" {
            return m, m.addTask(title)
        }
    case "esc":
        m.currentView = viewList
        m.addInput.SetValue("")
        return m, nil
    }

    // Update text input
    var cmd tea.Cmd
    m.addInput, cmd = m.addInput.Update(msg)
    return m, cmd
}

// Help view update handler
func (m Model) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    if msg.String() == "esc" || msg.String() == "?" {
        m.currentView = viewList
    }
    return m, nil
}
```

### Pattern 2: Sub-Model Delegation

**Best for**: Complex applications with 4+ views, each view has significant logic

```go
// File: pkg/tui/model.go
package tui

type viewState int

const (
    listView viewState = iota
    addView
    reportView
    configView
)

type Model struct {
    currentView viewState

    // Sub-models (each is a tea.Model)
    listModel   *ListModel
    addModel    *AddModel
    reportModel *ReportModel
    configModel *ConfigModel

    // Shared dependencies
    service *service.TaskService
    width   int
    height  int
}

func NewModel(svc *service.TaskService) Model {
    return Model{
        currentView: listView,
        service:     svc,
        listModel:   NewListModel(svc),
        addModel:    NewAddModel(svc),
        reportModel: NewReportModel(svc),
        configModel: NewConfigModel(svc),
    }
}

func (m Model) Init() tea.Cmd {
    // Initialize current view
    return m.listModel.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {

    case tea.KeyMsg:
        // Global navigation
        if msg.String() == "ctrl+c" {
            return m, tea.Quit
        }

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        // Propagate size to all sub-models
        m.listModel.SetSize(msg.Width, msg.Height)
        m.addModel.SetSize(msg.Width, msg.Height)
        m.reportModel.SetSize(msg.Width, msg.Height)
        m.configModel.SetSize(msg.Width, msg.Height)
        return m, nil

    case switchViewMsg:
        // Custom message for view switching
        m.currentView = msg.view

        // Call Init() on new view (important!)
        switch msg.view {
        case listView:
            cmd = m.listModel.Init()
        case addView:
            cmd = m.addModel.Init()
        case reportView:
            cmd = m.reportModel.Init()
        case configView:
            cmd = m.configModel.Init()
        }
        return m, cmd
    }

    // Delegate to active sub-model
    switch m.currentView {
    case listView:
        newModel, cmd := m.listModel.Update(msg)
        m.listModel = newModel.(*ListModel)
        return m, cmd
    case addView:
        newModel, cmd := m.addModel.Update(msg)
        m.addModel = newModel.(*AddModel)
        return m, cmd
    case reportView:
        newModel, cmd := m.reportModel.Update(msg)
        m.reportModel = newModel.(*ReportModel)
        return m, cmd
    case configView:
        newModel, cmd := m.configModel.Update(msg)
        m.configModel = newModel.(*ConfigModel)
        return m, cmd
    }

    return m, cmd
}

func (m Model) View() string {
    // Delegate rendering to active sub-model
    switch m.currentView {
    case listView:
        return m.listModel.View()
    case addView:
        return m.addModel.View()
    case reportView:
        return m.reportModel.View()
    case configView:
        return m.configModel.View()
    default:
        return "Unknown view"
    }
}

// Custom message for switching views
type switchViewMsg struct {
    view viewState
}

func switchView(view viewState) tea.Cmd {
    return func() tea.Msg {
        return switchViewMsg{view: view}
    }
}
```

### Pattern 3: Tab-Based Navigation

**Best for**: Applications with parallel views (like Pool/Today/Done)

```go
// File: pkg/tui/model.go
package tui

type tab int

const (
    poolTab tab = iota
    todayTab
    doneTab
)

type Model struct {
    activeTab tab
    tabs      []string

    // Shared task list
    allTasks []*model.Task

    // List component (reused across tabs)
    list     list.Model

    service  *service.TaskService
}

func NewModel(svc *service.TaskService) Model {
    // Create list delegate
    delegate := list.NewDefaultDelegate()

    // Create list
    l := list.New([]list.Item{}, delegate, 0, 0)
    l.Title = "Pool"
    l.SetShowStatusBar(false)
    l.SetFilteringEnabled(true)

    return Model{
        activeTab: poolTab,
        tabs:      []string{"Pool", "Today", "Done"},
        list:      l,
        service:   svc,
    }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit

        case "tab":
            m.activeTab = (m.activeTab + 1) % 3
            m.refreshList()
            return m, nil

        case "shift+tab":
            m.activeTab = (m.activeTab - 1 + 3) % 3
            m.refreshList()
            return m, nil
        }

    case tasksLoadedMsg:
        m.allTasks = msg.tasks
        m.refreshList()
        return m, nil
    }

    // Delegate to list component
    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}

func (m *Model) refreshList() {
    // Filter tasks based on active tab
    var filtered []*model.Task
    var title string

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

func (m Model) View() string {
    // Render tabs
    tabsView := m.renderTabs()

    // Render list
    listView := m.list.View()

    // Combine
    return lipgloss.JoinVertical(
        lipgloss.Left,
        tabsView,
        listView,
    )
}

func (m Model) renderTabs() string {
    var tabs []string

    for i, tab := range m.tabs {
        style := tabStyle
        if i == int(m.activeTab) {
            style = activeTabStyle
        }
        tabs = append(tabs, style.Render(tab))
    }

    return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}
```

### Key Principle: ALWAYS Call Init() When Switching

**Critical Rule**: When switching to a new view/model, you MUST call its `Init()` method:

```go
// ❌ WRONG - animations and background tasks won't work
case switchViewMsg:
    m.currentView = msg.view
    return m, nil

// ✅ CORRECT - properly initializes new view
case switchViewMsg:
    m.currentView = msg.view
    var cmd tea.Cmd
    switch msg.view {
    case listView:
        cmd = m.listModel.Init()
    case addView:
        cmd = m.addModel.Init()
    }
    return m, cmd
```

**Why**: The `Init()` method returns commands (like starting spinners, loading data, etc.) that must be executed. Without calling `Init()`, these background tasks won't start.

---

## Bubbles Components Reference

### Complete Component List

| Component | Purpose | Key Features |
|-----------|---------|--------------|
| **List** | Browse items with filtering | Pagination, fuzzy search, help, status |
| **Text Input** | Single-line text entry | Unicode, paste, scrolling, validation |
| **Text Area** | Multi-line text entry | Unicode, paste, vertical scroll, wrapping |
| **Table** | Tabular data display | Columns, rows, sorting, scrolling |
| **Progress** | Progress indication | Solid/gradient fills, percentage, custom |
| **Spinner** | Loading indicator | Multiple styles, custom frames |
| **Paginator** | Pagination UI | Dots or numbers, customizable |
| **Viewport** | Scrollable content | Mouse wheel, keybindings, wrapping |
| **File Picker** | File system navigation | Directory browsing, extension filter |
| **Timer** | Countdown timer | Custom format, events on completion |
| **Stopwatch** | Count-up timer | Precise timing, lap tracking |
| **Help** | Keybinding help | Auto-generated, single/multi-line |
| **Key** | Keybinding management | Remapping, help text, enabled state |

### Recommended Components for togo

#### 1. List Component (Primary View)

**Use for**: Displaying Pool/Today/Done task lists

```go
import (
    "github.com/charmbracelet/bubbles/list"
    tea "github.com/charmbracelet/bubbletea"
)

// Task item implements list.Item interface
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
    // Show tags and due date
    desc := ""
    if len(i.task.Tags) > 0 {
        desc += fmt.Sprintf("Tags: %s", strings.Join(i.task.Tags, ", "))
    }
    if i.task.DueDate != nil {
        desc += fmt.Sprintf(" | Due: %s", i.task.DueDate.Format("2006-01-02"))
    }
    return desc
}

// Custom delegate for styling
type taskDelegate struct{}

func (d taskDelegate) Height() int { return 2 }
func (d taskDelegate) Spacing() int { return 1 }
func (d taskDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d taskDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
    i, ok := item.(taskItem)
    if !ok {
        return
    }

    task := i.task

    // Style based on status
    var style lipgloss.Style
    var icon string

    switch task.Status {
    case model.StatusPool:
        style = poolStyle
        icon = "○"
    case model.StatusToday:
        style = todayStyle
        icon = "◉"
    case model.StatusDone:
        style = doneStyle
        icon = "✓"
    }

    // Highlight if selected
    if index == m.Index() {
        style = style.Copy().Foreground(lipgloss.Color("212"))
    }

    // Render
    title := style.Render(fmt.Sprintf("%s %s", icon, task.Title))
    desc := subtleStyle.Render(i.Description())

    fmt.Fprintf(w, "%s\n%s", title, desc)
}

// Creating the list
func NewTaskList() list.Model {
    delegate := taskDelegate{}

    l := list.New([]list.Item{}, delegate, 80, 20)
    l.Title = "Tasks"
    l.SetShowStatusBar(true)
    l.SetFilteringEnabled(true)
    l.Styles.Title = titleStyle

    // Custom key bindings
    l.AdditionalShortHelpKeys = func() []key.Binding {
        return []key.Binding{
            key.NewBinding(
                key.WithKeys("enter"),
                key.WithHelp("enter", "pick task"),
            ),
            key.NewBinding(
                key.WithKeys("x"),
                key.WithHelp("x", "complete"),
            ),
        }
    }

    return l
}
```

#### 2. Text Input Component (Quick Add)

**Use for**: Adding tasks with single-line input

```go
import "github.com/charmbracelet/bubbles/textinput"

type AddModel struct {
    titleInput textinput.Model
    tagsInput  textinput.Model
    focusIndex int
}

func NewAddModel() AddModel {
    ti := textinput.New()
    ti.Placeholder = "Task title..."
    ti.Focus()
    ti.CharLimit = 200
    ti.Width = 50

    tags := textinput.New()
    tags.Placeholder = "Tags (comma-separated)..."
    tags.CharLimit = 100
    tags.Width = 50

    return AddModel{
        titleInput: ti,
        tagsInput:  tags,
        focusIndex: 0,
    }
}

func (m AddModel) Init() tea.Cmd {
    return textinput.Blink
}

func (m AddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            m.focusIndex = (m.focusIndex + 1) % 2
            if m.focusIndex == 0 {
                m.titleInput.Focus()
                m.tagsInput.Blur()
            } else {
                m.titleInput.Blur()
                m.tagsInput.Focus()
            }
            return m, nil

        case "enter":
            // Submit form
            return m, m.submitTask()

        case "esc":
            return m, switchView(listView)
        }
    }

    // Update focused input
    var cmd tea.Cmd
    if m.focusIndex == 0 {
        m.titleInput, cmd = m.titleInput.Update(msg)
    } else {
        m.tagsInput, cmd = m.tagsInput.Update(msg)
    }

    return m, cmd
}

func (m AddModel) View() string {
    return fmt.Sprintf(
        "%s\n\n%s\n\n%s\n\n%s\n\n%s",
        titleStyle.Render("Add New Task"),
        m.titleInput.View(),
        m.tagsInput.View(),
        subtleStyle.Render("Tab: next field | Enter: save | Esc: cancel"),
        m.statusMsg,
    )
}
```

#### 3. Text Area Component (Notes/Details)

**Use for**: Editing task notes or details

```go
import "github.com/charmbracelet/bubbles/textarea"

type NotesModel struct {
    textarea textarea.Model
    task     *model.Task
}

func NewNotesModel(task *model.Task) NotesModel {
    ta := textarea.New()
    ta.Placeholder = "Enter task notes..."
    ta.SetValue(task.Notes)
    ta.Focus()
    ta.CharLimit = 5000

    return NotesModel{
        textarea: ta,
        task:     task,
    }
}

func (m NotesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+s":
            return m, m.saveNotes()
        case "esc":
            return m, switchView(listView)
        }
    }

    var cmd tea.Cmd
    m.textarea, cmd = m.textarea.Update(msg)
    return m, cmd
}

func (m NotesModel) View() string {
    return fmt.Sprintf(
        "%s\n\n%s\n\n%s",
        titleStyle.Render(fmt.Sprintf("Notes: %s", m.task.Title)),
        m.textarea.View(),
        subtleStyle.Render("Ctrl+S: save | Esc: cancel"),
    )
}
```

#### 4. Spinner Component (Loading States)

**Use for**: Indicating sync operations or loading tasks

```go
import "github.com/charmbracelet/bubbles/spinner"

type Model struct {
    spinner  spinner.Model
    loading  bool
    tasks    []*model.Task
}

func NewModel() Model {
    s := spinner.New()
    s.Spinner = spinner.Dot
    s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

    return Model{
        spinner: s,
        loading: true,
    }
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(
        m.spinner.Tick,
        loadTasks(),
    )
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tasksLoadedMsg:
        m.tasks = msg.tasks
        m.loading = false
        return m, nil

    case spinner.TickMsg:
        var cmd tea.Cmd
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd
    }

    return m, nil
}

func (m Model) View() string {
    if m.loading {
        return fmt.Sprintf("\n\n   %s Loading tasks...\n\n", m.spinner.View())
    }

    // Render task list
    return m.renderTasks()
}
```

#### 5. Help Component (Keybinding Reference)

**Use for**: Showing available commands

```go
import (
    "github.com/charmbracelet/bubbles/help"
    "github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
    Up       key.Binding
    Down     key.Binding
    Pick     key.Binding
    Complete key.Binding
    Defer    key.Binding
    Add      key.Binding
    Quit     key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
    return []key.Binding{k.Up, k.Down, k.Pick, k.Add, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {k.Up, k.Down, k.Pick, k.Complete},
        {k.Defer, k.Add, k.Quit},
    }
}

var keys = keyMap{
    Up: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("↑/k", "move up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("down", "j"),
        key.WithHelp("↓/j", "move down"),
    ),
    Pick: key.NewBinding(
        key.WithKeys("enter"),
        key.WithHelp("enter", "pick task"),
    ),
    Complete: key.NewBinding(
        key.WithKeys("x"),
        key.WithHelp("x", "complete"),
    ),
    Defer: key.NewBinding(
        key.WithKeys("d"),
        key.WithHelp("d", "defer"),
    ),
    Add: key.NewBinding(
        key.WithKeys("a"),
        key.WithHelp("a", "add task"),
    ),
    Quit: key.NewBinding(
        key.WithKeys("q", "ctrl+c"),
        key.WithHelp("q", "quit"),
    ),
}

type Model struct {
    help   help.Model
    keys   keyMap
    // ... other fields
}

func NewModel() Model {
    return Model{
        help: help.New(),
        keys: keys,
    }
}

func (m Model) View() string {
    // Main content
    content := m.renderContent()

    // Help footer
    helpView := m.help.View(m.keys)

    return lipgloss.JoinVertical(
        lipgloss.Left,
        content,
        "\n",
        helpView,
    )
}
```

---

## Lip Gloss Styling Patterns

### Core Styling Concepts

Lip Gloss uses a **declarative, CSS-like API** for terminal styling:

```go
import "github.com/charmbracelet/lipgloss"

// Define styles (like CSS classes)
var (
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FAFAFA")).
        Background(lipgloss.Color("#7D56F4")).
        Padding(0, 1).
        MarginTop(1)

    subtleStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("241"))

    boxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("63")).
        Padding(1, 2).
        Width(50)
)

// Use styles
func render() string {
    title := titleStyle.Render("My App")
    subtitle := subtleStyle.Render("A cool TUI application")

    content := lipgloss.JoinVertical(
        lipgloss.Left,
        title,
        subtitle,
    )

    return boxStyle.Render(content)
}
```

### Complete Styling System for togo

```go
// File: pkg/tui/styles.go
package tui

import (
    "github.com/charmbracelet/lipgloss"
)

// Color palette
var (
    primaryColor   = lipgloss.Color("205") // Pink
    secondaryColor = lipgloss.Color("99")  // Purple
    accentColor    = lipgloss.Color("212") // Light pink

    poolColor  = lipgloss.Color("63")  // Blue
    todayColor = lipgloss.Color("212") // Pink
    doneColor  = lipgloss.Color("42")  // Green

    mutedColor = lipgloss.Color("241") // Gray
    errorColor = lipgloss.Color("196") // Red
)

// Text styles
var (
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(primaryColor).
        MarginBottom(1)

    subtitleStyle = lipgloss.NewStyle().
        Foreground(secondaryColor).
        Italic(true)

    helpStyle = lipgloss.NewStyle().
        Foreground(mutedColor).
        MarginTop(1)

    errorStyle = lipgloss.NewStyle().
        Foreground(errorColor).
        Bold(true)
)

// Status-based styles
var (
    poolStyle = lipgloss.NewStyle().
        Foreground(poolColor)

    todayStyle = lipgloss.NewStyle().
        Foreground(todayColor).
        Bold(true)

    doneStyle = lipgloss.NewStyle().
        Foreground(doneColor).
        Strikethrough(true)
)

// Layout styles
var (
    containerStyle = lipgloss.NewStyle().
        Padding(1, 2)

    boxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(primaryColor).
        Padding(1, 2).
        MarginBottom(1)

    tabStyle = lipgloss.NewStyle().
        Padding(0, 2).
        Background(lipgloss.Color("238")).
        Foreground(lipgloss.Color("248"))

    activeTabStyle = lipgloss.NewStyle().
        Padding(0, 2).
        Background(primaryColor).
        Foreground(lipgloss.Color("#000000")).
        Bold(true)
)

// List item styles
var (
    selectedItemStyle = lipgloss.NewStyle().
        Foreground(accentColor).
        Bold(true).
        PaddingLeft(2)

    normalItemStyle = lipgloss.NewStyle().
        PaddingLeft(4)

    taskIconStyle = lipgloss.NewStyle().
        Width(2)
)

// Responsive sizing
func (m Model) boxWidth() int {
    width := m.width - 4 // Account for padding
    if width > 100 {
        width = 100 // Max width for readability
    }
    return width
}

// Adaptive colors (light/dark terminal)
var adaptiveStyle = lipgloss.NewStyle().
    Foreground(lipgloss.AdaptiveColor{
        Light: "#000000",
        Dark:  "#FFFFFF",
    })

// Gradient backgrounds
var gradientStyle = lipgloss.NewStyle().
    Background(lipgloss.Color("62")).
    Foreground(lipgloss.Color("230"))

// Complex layout composition
func (m Model) viewList() string {
    // Header
    header := titleStyle.Render("togo - Tasks")

    // Tabs
    tabs := m.renderTabs()

    // Task list
    taskList := m.renderTasks()

    // Status bar
    status := m.renderStatusBar()

    // Help
    help := helpStyle.Render(m.help.View(m.keys))

    // Compose layout
    content := lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        tabs,
        "",
        taskList,
        "",
        status,
        help,
    )

    // Apply container
    return containerStyle.
        Width(m.boxWidth()).
        Render(content)
}
```

### Advanced Layout Patterns

#### 1. Side-by-Side Layout

```go
func renderSideBySide(left, right string) string {
    leftBox := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1).
        Width(40).
        Render(left)

    rightBox := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1).
        Width(40).
        Render(right)

    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        leftBox,
        rightBox,
    )
}
```

#### 2. Three-Column Layout

```go
func renderThreeColumns(pool, today, done []string) string {
    poolCol := renderColumn("Pool", pool, poolColor)
    todayCol := renderColumn("Today", today, todayColor)
    doneCol := renderColumn("Done", done, doneColor)

    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        poolCol,
        todayCol,
        doneCol,
    )
}

func renderColumn(title string, items []string, color lipgloss.Color) string {
    header := lipgloss.NewStyle().
        Foreground(color).
        Bold(true).
        Render(title)

    itemsStr := strings.Join(items, "\n")

    content := lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        "",
        itemsStr,
    )

    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(color).
        Padding(1).
        Width(30).
        Height(20).
        Render(content)
}
```

#### 3. Responsive Layout

```go
func (m Model) View() string {
    // Switch layout based on terminal width
    if m.width < 80 {
        // Narrow: single column
        return m.renderSingleColumn()
    } else if m.width < 120 {
        // Medium: two columns
        return m.renderTwoColumns()
    } else {
        // Wide: three columns
        return m.renderThreeColumns()
    }
}
```

#### 4. Tables

```go
import "github.com/charmbracelet/lipgloss/table"

func renderTaskTable(tasks []*model.Task) string {
    t := table.New().
        Border(lipgloss.NormalBorder()).
        BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
        Headers("ID", "Title", "Status", "Tags")

    for _, task := range tasks {
        t.Row(
            task.ID.String()[:8], // Short ID
            task.Title,
            string(task.Status),
            strings.Join(task.Tags, ", "),
        )
    }

    return t.Render()
}
```

---

## Complete Implementation Example

Here's a complete, production-ready example for togo:

```go
// File: pkg/tui/model.go
package tui

import (
    "context"
    "fmt"
    "strings"

    "github.com/charmbracelet/bubbles/help"
    "github.com/charmbracelet/bubbles/key"
    "github.com/charmbracelet/bubbles/list"
    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"

    "togo/internal/model"
    "togo/internal/service"
)

// ============================================================================
// Model Definition
// ============================================================================

type viewState int

const (
    listView viewState = iota
    addView
    detailView
)

type tab int

const (
    poolTab tab = iota
    todayTab
    doneTab
)

type Model struct {
    // Dependencies
    service *service.TaskService

    // State
    currentView viewState
    activeTab   tab

    // Data
    allTasks []*model.Task

    // Components
    list      list.Model
    addInput  textinput.Model
    spinner   spinner.Model
    help      help.Model
    keys      keyMap

    // UI state
    loading    bool
    statusMsg  string
    width      int
    height     int
    quitting   bool
}

// ============================================================================
// Initialization
// ============================================================================

func NewModel(svc *service.TaskService) Model {
    // Initialize list
    delegate := newTaskDelegate()
    l := list.New([]list.Item{}, delegate, 0, 0)
    l.Title = "Pool"
    l.SetShowStatusBar(false)
    l.SetFilteringEnabled(true)
    l.Styles.Title = titleStyle

    // Initialize text input
    ti := textinput.New()
    ti.Placeholder = "Task title..."
    ti.CharLimit = 200

    // Initialize spinner
    s := spinner.New()
    s.Spinner = spinner.Dot
    s.Style = lipgloss.NewStyle().Foreground(primaryColor)

    return Model{
        service:     svc,
        currentView: listView,
        activeTab:   poolTab,
        list:        l,
        addInput:    ti,
        spinner:     s,
        help:        help.New(),
        keys:        keys,
        loading:     true,
    }
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(
        m.spinner.Tick,
        loadTasks(m.service),
    )
}

// ============================================================================
// Update
// ============================================================================

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    // Window size
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.list.SetSize(msg.Width-4, msg.Height-10)
        return m, nil

    // Tasks loaded
    case tasksLoadedMsg:
        m.allTasks = msg.tasks
        m.loading = false
        m.statusMsg = fmt.Sprintf("Loaded %d tasks", len(msg.tasks))
        m.refreshList()
        return m, nil

    // Task operation result
    case taskOperationMsg:
        if msg.err != nil {
            m.statusMsg = errorStyle.Render(fmt.Sprintf("Error: %v", msg.err))
        } else {
            m.statusMsg = fmt.Sprintf("✓ %s", msg.message)
        }
        return m, loadTasks(m.service)

    // Spinner tick
    case spinner.TickMsg:
        var cmd tea.Cmd
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd

    // Keyboard input
    case tea.KeyMsg:
        // Global keys
        switch msg.String() {
        case "ctrl+c", "q":
            if m.currentView != addView {
                m.quitting = true
                return m, tea.Quit
            }
        case "?":
            m.help.ShowAll = !m.help.ShowAll
            return m, nil
        }

        // View-specific routing
        switch m.currentView {
        case listView:
            return m.updateList(msg)
        case addView:
            return m.updateAdd(msg)
        }
    }

    // Delegate to list if in list view
    if m.currentView == listView {
        var cmd tea.Cmd
        m.list, cmd = m.list.Update(msg)
        return m, cmd
    }

    return m, nil
}

func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {

    // Tab navigation
    case "tab":
        m.activeTab = (m.activeTab + 1) % 3
        m.refreshList()
        return m, nil

    case "shift+tab":
        m.activeTab = (m.activeTab - 1 + 3) % 3
        m.refreshList()
        return m, nil

    // Task operations
    case "a":
        m.currentView = addView
        m.addInput.Focus()
        return m, textinput.Blink

    case "enter":
        // Pick task
        return m, m.pickSelectedTask()

    case "x":
        // Complete task
        return m, m.completeSelectedTask()

    case "d":
        // Defer task
        return m, m.deferSelectedTask()

    case "delete", "backspace":
        // Remove task
        return m, m.removeSelectedTask()
    }

    return m, nil
}

func (m Model) updateAdd(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {

    case "enter":
        title := strings.TrimSpace(m.addInput.Value())
        if title != "" {
            m.currentView = listView
            m.addInput.SetValue("")
            m.addInput.Blur()
            return m, m.addTask(title)
        }

    case "esc":
        m.currentView = listView
        m.addInput.SetValue("")
        m.addInput.Blur()
        return m, nil
    }

    var cmd tea.Cmd
    m.addInput, cmd = m.addInput.Update(msg)
    return m, cmd
}

// ============================================================================
// View
// ============================================================================

func (m Model) View() string {
    if m.quitting {
        return "\n  See you later! 👋\n\n"
    }

    if m.loading {
        return fmt.Sprintf(
            "\n\n   %s Loading tasks...\n\n",
            m.spinner.View(),
        )
    }

    switch m.currentView {
    case listView:
        return m.viewList()
    case addView:
        return m.viewAdd()
    default:
        return "Unknown view"
    }
}

func (m Model) viewList() string {
    // Header
    header := titleStyle.Render("togo")

    // Tabs
    tabs := m.renderTabs()

    // List
    listView := m.list.View()

    // Status bar
    status := m.renderStatusBar()

    // Help
    helpView := helpStyle.Render(m.help.View(m.keys))

    // Compose
    content := lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        tabs,
        "",
        listView,
        "",
        status,
        helpView,
    )

    return containerStyle.Render(content)
}

func (m Model) viewAdd() string {
    title := titleStyle.Render("Add New Task")
    input := m.addInput.View()
    help := helpStyle.Render("Enter: save | Esc: cancel")

    content := lipgloss.JoinVertical(
        lipgloss.Left,
        title,
        "",
        input,
        "",
        help,
    )

    return boxStyle.
        Width(60).
        Height(10).
        Render(content)
}

func (m Model) renderTabs() string {
    tabs := []string{"Pool", "Today", "Done"}
    var rendered []string

    for i, tab := range tabs {
        style := tabStyle
        if i == int(m.activeTab) {
            style = activeTabStyle
        }
        rendered = append(rendered, style.Render(tab))
    }

    return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

func (m Model) renderStatusBar() string {
    poolCount := 0
    todayCount := 0
    doneCount := 0

    for _, t := range m.allTasks {
        switch t.Status {
        case model.StatusPool:
            poolCount++
        case model.StatusToday:
            todayCount++
        case model.StatusDone:
            doneCount++
        }
    }

    stats := fmt.Sprintf(
        "Pool: %d | Today: %d | Done: %d",
        poolCount, todayCount, doneCount,
    )

    left := subtleStyle.Render(stats)
    right := subtleStyle.Render(m.statusMsg)

    gap := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 4
    if gap < 0 {
        gap = 0
    }

    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        left,
        strings.Repeat(" ", gap),
        right,
    )
}

// ============================================================================
// Helpers
// ============================================================================

func (m *Model) refreshList() {
    // Filter tasks by active tab
    var filtered []*model.Task
    var title string

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

func (m Model) getSelectedTask() *model.Task {
    item, ok := m.list.SelectedItem().(taskItem)
    if !ok {
        return nil
    }
    return item.task
}

// ============================================================================
// Commands (Async Operations)
// ============================================================================

type tasksLoadedMsg struct {
    tasks []*model.Task
}

type taskOperationMsg struct {
    message string
    err     error
}

func loadTasks(svc *service.TaskService) tea.Cmd {
    return func() tea.Msg {
        filter := model.TaskFilter{}
        tasks, err := svc.ListTasks(context.Background(), filter)
        if err != nil {
            return taskOperationMsg{err: err}
        }
        return tasksLoadedMsg{tasks: tasks}
    }
}

func (m Model) addTask(title string) tea.Cmd {
    return func() tea.Msg {
        req := service.AddTaskRequest{Title: title}
        _, err := m.service.AddTask(context.Background(), req)
        return taskOperationMsg{
            message: "Task added",
            err:     err,
        }
    }
}

func (m Model) pickSelectedTask() tea.Cmd {
    task := m.getSelectedTask()
    if task == nil {
        return nil
    }

    return func() tea.Msg {
        err := m.service.PickTask(context.Background(), task.ID)
        return taskOperationMsg{
            message: "Task picked",
            err:     err,
        }
    }
}

func (m Model) completeSelectedTask() tea.Cmd {
    task := m.getSelectedTask()
    if task == nil {
        return nil
    }

    return func() tea.Msg {
        err := m.service.CompleteTask(context.Background(), task.ID)
        return taskOperationMsg{
            message: "Task completed",
            err:     err,
        }
    }
}

func (m Model) deferSelectedTask() tea.Cmd {
    task := m.getSelectedTask()
    if task == nil {
        return nil
    }

    return func() tea.Msg {
        err := m.service.DeferTask(context.Background(), task.ID)
        return taskOperationMsg{
            message: "Task deferred",
            err:     err,
        }
    }
}

func (m Model) removeSelectedTask() tea.Cmd {
    task := m.getSelectedTask()
    if task == nil {
        return nil
    }

    return func() tea.Msg {
        err := m.service.RemoveTask(context.Background(), task.ID)
        return taskOperationMsg{
            message: "Task removed",
            err:     err,
        }
    }
}

// ============================================================================
// List Item Implementation
// ============================================================================

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

    if i.task.DeferredCount > 0 {
        parts = append(parts, fmt.Sprintf("Deferred %d times", i.task.DeferredCount))
    }

    return strings.Join(parts, " | ")
}

// ============================================================================
// Key Bindings
// ============================================================================

type keyMap struct {
    Up       key.Binding
    Down     key.Binding
    Tab      key.Binding
    Pick     key.Binding
    Complete key.Binding
    Defer    key.Binding
    Add      key.Binding
    Remove   key.Binding
    Help     key.Binding
    Quit     key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
    return []key.Binding{k.Up, k.Down, k.Pick, k.Add, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {k.Up, k.Down, k.Tab, k.Pick},
        {k.Complete, k.Defer, k.Add, k.Remove},
        {k.Help, k.Quit},
    }
}

var keys = keyMap{
    Up: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("↑/k", "up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("down", "j"),
        key.WithHelp("↓/j", "down"),
    ),
    Tab: key.NewBinding(
        key.WithKeys("tab"),
        key.WithHelp("tab", "switch view"),
    ),
    Pick: key.NewBinding(
        key.WithKeys("enter"),
        key.WithHelp("enter", "pick"),
    ),
    Complete: key.NewBinding(
        key.WithKeys("x"),
        key.WithHelp("x", "complete"),
    ),
    Defer: key.NewBinding(
        key.WithKeys("d"),
        key.WithHelp("d", "defer"),
    ),
    Add: key.NewBinding(
        key.WithKeys("a"),
        key.WithHelp("a", "add"),
    ),
    Remove: key.NewBinding(
        key.WithKeys("delete", "backspace"),
        key.WithHelp("del", "remove"),
    ),
    Help: key.NewBinding(
        key.WithKeys("?"),
        key.WithHelp("?", "help"),
    ),
    Quit: key.NewBinding(
        key.WithKeys("q", "ctrl+c"),
        key.WithHelp("q", "quit"),
    ),
}
```

```go
// File: pkg/tui/delegate.go
package tui

import (
    "fmt"
    "io"

    "github.com/charmbracelet/bubbles/list"
    tea "github.com/charmbracelet/bubbletea"

    "togo/internal/model"
)

type taskDelegate struct{}

func newTaskDelegate() taskDelegate {
    return taskDelegate{}
}

func (d taskDelegate) Height() int                               { return 2 }
func (d taskDelegate) Spacing() int                              { return 1 }
func (d taskDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d taskDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
    i, ok := item.(taskItem)
    if !ok {
        return
    }

    task := i.task

    // Icon and style based on status
    var icon string
    var titleStyle lipgloss.Style

    switch task.Status {
    case model.StatusPool:
        icon = "○"
        titleStyle = poolStyle
    case model.StatusToday:
        icon = "◉"
        titleStyle = todayStyle
    case model.StatusDone:
        icon = "✓"
        titleStyle = doneStyle
    }

    // Highlight selected
    if index == m.Index() {
        titleStyle = selectedItemStyle
        icon = "▶"
    }

    // Render
    title := titleStyle.Render(fmt.Sprintf("%s %s", icon, task.Title))
    desc := subtleStyle.Render(i.Description())

    fmt.Fprintf(w, "%s\n%s", title, desc)
}
```

```go
// File: pkg/tui/styles.go
package tui

import "github.com/charmbracelet/lipgloss"

var (
    primaryColor = lipgloss.Color("205")
    poolColor    = lipgloss.Color("63")
    todayColor   = lipgloss.Color("212")
    doneColor    = lipgloss.Color("42")
    mutedColor   = lipgloss.Color("241")
    errorColor   = lipgloss.Color("196")
)

var (
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(primaryColor).
        MarginBottom(1)

    subtleStyle = lipgloss.NewStyle().
        Foreground(mutedColor)

    helpStyle = lipgloss.NewStyle().
        Foreground(mutedColor).
        MarginTop(1)

    errorStyle = lipgloss.NewStyle().
        Foreground(errorColor).
        Bold(true)

    poolStyle = lipgloss.NewStyle().
        Foreground(poolColor)

    todayStyle = lipgloss.NewStyle().
        Foreground(todayColor).
        Bold(true)

    doneStyle = lipgloss.NewStyle().
        Foreground(doneColor).
        Strikethrough(true)

    selectedItemStyle = lipgloss.NewStyle().
        Foreground(primaryColor).
        Bold(true)

    containerStyle = lipgloss.NewStyle().
        Padding(1, 2)

    boxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(primaryColor).
        Padding(1, 2)

    tabStyle = lipgloss.NewStyle().
        Padding(0, 2).
        Background(lipgloss.Color("238")).
        Foreground(lipgloss.Color("248"))

    activeTabStyle = lipgloss.NewStyle().
        Padding(0, 2).
        Background(primaryColor).
        Foreground(lipgloss.Color("0")).
        Bold(true)
)
```

---

## Migration Strategy

### Step-by-Step Migration from Cobra

#### Phase 1: Identify Cobra Usage

Audit current Cobra implementation:

```bash
# Find all Cobra command definitions
grep -r "cobra.Command" .

# Find all flag definitions
grep -r "cmd.Flags()" .

# List all subcommands
ls cmd/todo/commands/
```

#### Phase 2: Choose Architecture

**Decision Point**: Hybrid or Pure Bubble Tea?

**Hybrid (Recommended)**:
- Keep Cobra for CLI commands
- Add `tui` subcommand that launches Bubble Tea
- Users get both modes

**Pure Bubble Tea**:
- Replace Cobra with `flag` or `pflag`
- Implement manual command dispatcher
- Launch TUI when no args provided

#### Phase 3: Implement Gradually

**Week 1**: TUI Implementation
- Build Bubble Tea model with all views
- Integrate with existing TaskService
- Test interactive operations

**Week 2**: CLI Integration
- If hybrid: Add `tui` command to Cobra
- If pure: Implement command dispatcher
- Ensure both paths use same service layer

**Week 3**: Polish & Testing
- Add help documentation
- Style with Lip Gloss
- Test all workflows

#### Phase 4: User Migration

**If using hybrid**:
```bash
# Old: Cobra only
todo add "Task"
todo list

# New: Cobra + TUI
todo add "Task"      # Still works
todo list            # Still works
todo tui             # New interactive mode
```

**If going pure Bubble Tea**:
```bash
# Old: Cobra
todo add "Task"
todo list

# New: Pure Bubble Tea
todo add "Task"      # Works (via flag parser)
todo list            # Works (via flag parser)
todo                 # Launches TUI (new default)
```

---

## Best Practices

### 1. Separation of Concerns

**Keep TUI and business logic separate**:

```go
// ❌ BAD: Business logic in Update
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case tea.KeyMsg:
        if msg.String() == "x" {
            // DON'T: Inline business logic
            task := m.tasks[m.cursor]
            task.Status = model.StatusDone
            task.CompletedAt = time.Now()
            // Save to storage...
        }
}

// ✅ GOOD: Delegate to service layer
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case tea.KeyMsg:
        if msg.String() == "x" {
            return m, m.completeTask()
        }
}

func (m Model) completeTask() tea.Cmd {
    return func() tea.Msg {
        err := m.service.CompleteTask(ctx, m.tasks[m.cursor].ID)
        return taskCompletedMsg{err: err}
    }
}
```

### 2. Always Call Init() on View Changes

```go
// ✅ CORRECT
case switchViewMsg:
    m.currentView = msg.view
    return m, m.subModel.Init()  // MUST call Init()
```

### 3. Use Custom Messages for Communication

```go
// Define custom messages
type taskAddedMsg struct {
    task *model.Task
    err  error
}

type syncCompleteMsg struct {
    status string
    err    error
}

// Use in commands
func addTask(svc *service.TaskService, title string) tea.Cmd {
    return func() tea.Msg {
        task, err := svc.AddTask(context.Background(), AddTaskRequest{Title: title})
        return taskAddedMsg{task: task, err: err}
    }
}

// Handle in Update
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case taskAddedMsg:
        if msg.err != nil {
            m.statusMsg = fmt.Sprintf("Error: %v", msg.err)
        } else {
            m.tasks = append(m.tasks, msg.task)
            m.statusMsg = "Task added!"
        }
        return m, nil
    }
}
```

### 4. Handle Window Resize

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

        // Update all components
        m.list.SetSize(msg.Width-4, msg.Height-10)
        m.viewport.Width = msg.Width - 4
        m.viewport.Height = msg.Height - 8

        return m, nil
    }
}
```

### 5. Proper Error Handling in TUI

```go
// Show errors in status bar
func (m Model) handleError(err error) {
    m.statusMsg = errorStyle.Render(fmt.Sprintf("Error: %v", err))
}

// Or dedicated error view
type errorMsg struct {
    err error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case errorMsg:
        m.currentView = errorView
        m.errorText = msg.err.Error()
        return m, nil
}
```

### 6. Testability

```go
// Test Update logic
func TestModel_PickTask(t *testing.T) {
    svc := &MockTaskService{}
    m := NewModel(svc)
    m.tasks = []*model.Task{{ID: "123", Status: model.StatusPool}}

    // Simulate Enter key
    msg := tea.KeyMsg{Type: tea.KeyEnter}
    newModel, cmd := m.Update(msg)

    // Assert command returned
    assert.NotNil(t, cmd)

    // Execute command
    resultMsg := cmd()

    // Assert service called
    assert.Equal(t, 1, svc.PickTaskCallCount)
}

// Test View rendering
func TestModel_View(t *testing.T) {
    m := NewModel(nil)
    m.currentView = listView

    view := m.View()

    assert.Contains(t, view, "togo")
    assert.Contains(t, view, "Pool")
}
```

---

## Trade-offs and Considerations

### Cobra vs Pure Bubble Tea

| Aspect | Cobra | Pure Bubble Tea |
|--------|-------|-----------------|
| **CLI Commands** | ✅ Excellent | ⚠️ Manual implementation needed |
| **Help Generation** | ✅ Automatic | ⚠️ Must write manually |
| **Flag Parsing** | ✅ Rich & validated | ⚠️ Basic (use pflag) |
| **Subcommands** | ✅ Natural | ⚠️ Manual routing |
| **Interactive TUI** | ⚠️ Need integration | ✅ Native |
| **User Experience** | ✅ Familiar CLI | ✅ Beautiful interactive |
| **Binary Size** | ⚠️ +1-2MB | ✅ Smaller |
| **Learning Curve** | ✅ Well-documented | ⚠️ New paradigm |

### Recommendations

#### Use Hybrid (Cobra + Bubble Tea) if:
- You need robust CLI with many flags/subcommands
- Users expect traditional CLI behavior
- You want both quick commands AND interactive mode
- Professional/enterprise tool

#### Use Pure Bubble Tea if:
- TUI is the primary interface
- CLI is just convenience shortcuts
- You want minimal dependencies
- Personal/hobby project

### Performance Considerations

**Bubble Tea**:
- Efficient rendering (only redraws changed parts)
- Handles large lists well (with viewport)
- Mouse support adds minimal overhead
- Alt screen keeps terminal history clean

**Memory**:
- List component: ~100 bytes per item
- Typical usage: <10MB for 1000 tasks
- No concerns for togo use case

### Accessibility

**Bubble Tea**:
- ✅ Fully keyboard-driven
- ✅ Works with screen readers (plain text)
- ⚠️ Mouse support optional
- ⚠️ Colors may not work in all terminals (use adaptive colors)

**Best Practices**:
- Always provide keyboard alternatives
- Use `lipgloss.AdaptiveColor` for light/dark terminals
- Test with `NO_COLOR=1` environment variable
- Provide `--no-tui` flag for pure text output

---

## Conclusion

### Summary of Findings

1. **Bubble Tea does NOT replace Cobra** - they solve different problems
   - Cobra: CLI parsing and routing
   - Bubble Tea: Interactive TUI framework

2. **Recommended Architecture**: Hybrid approach
   - Use Cobra (or pflag) for CLI commands
   - Use Bubble Tea for interactive TUI mode
   - Share TaskService layer between both

3. **Key Components for togo**:
   - **List**: Task lists (Pool/Today/Done)
   - **Text Input**: Quick task addition
   - **Spinner**: Loading states during sync
   - **Help**: Keybinding reference

4. **Multi-View Pattern**: State-based routing
   - Simple boolean for 2-3 views
   - Sub-model delegation for complex apps
   - Always call Init() when switching

5. **Styling with Lip Gloss**: Define reusable styles
   - Status-based colors (pool/today/done)
   - Responsive layouts
   - Adaptive colors for accessibility

### Next Steps for Implementation

1. **Keep Cobra** for command routing (or implement minimal dispatcher)
2. **Build Bubble Tea TUI** using patterns from this document
3. **Integrate both** through shared service layer
4. **Test thoroughly** - both CLI and TUI paths
5. **Document keybindings** - help users navigate

### Additional Resources

- **Bubble Tea Examples**: https://github.com/charmbracelet/bubbletea/tree/main/examples
- **Bubbles Components**: https://github.com/charmbracelet/bubbles
- **Lip Gloss Docs**: https://github.com/charmbracelet/lipgloss
- **Real-world Apps**: Glow, Soft Serve, VHS (all by Charm)

---

**Files Generated**:
- `/var/home/nj/code/togo/CHARM_RESEARCH.md` (this document)

**Recommended Package Structure**:
```
togo/
├── cmd/todo/
│   ├── main.go              # Entry point (Cobra root or flag parser)
│   └── commands/
│       ├── add.go
│       ├── list.go
│       └── tui.go           # TUI command (launches Bubble Tea)
├── pkg/tui/
│   ├── model.go             # Main Bubble Tea model
│   ├── delegate.go          # List item rendering
│   ├── styles.go            # Lip Gloss styles
│   └── commands.go          # Tea commands (async ops)
└── internal/
    └── service/
        └── task_service.go  # Shared by CLI and TUI
```
