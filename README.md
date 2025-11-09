# togo

A small example Bubble Tea application that demonstrates a selectable list.

This repository contains a simple Bubble Tea model defined in `main.go`. The app renders a list of choices and allows navigating with `j`/`k` or `up`/`down`, toggling selection with `enter` or space, and quitting with `q` or `ctrl+c`.

## Requirements

- Go 1.24+
- Modules enabled (this repo includes `go.mod`)

## Run the app

```bash
# run the program
go run main.go
```

## Tests

Unit tests were added to exercise the model logic (no TTY required). They call the model's `Update` and `View` directly to keep tests fast and deterministic.

Covered cases:

- InitializeModel: ensures the default choices and empty selection map
- Navigation bounds: `up`/`down` don't move the cursor out of range
- Toggle selection: `enter`/space toggles items in the `selected` map
- Quit command: `q` and `ctrl+c` produce a non-nil quit command
- View rendering: the view contains the expected cursor marker `>` and selection marker `x`

Run tests:

```bash
# run all tests
go test ./...

# verbose
go test -v ./...

# coverage report (opens HTML)
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
```

## Notes

- Tests are unit-level and avoid spawning a real `tea.Program`/TTY. If you want end-to-end/integration tests that run the full program and renderer, I can add one using `tea.Program` options (e.g., `WithOutput`/`WithoutRenderer`).

If you want the README expanded with contribution notes or CI test instructions, tell me what CI provider you use and I can add a simple workflow.
