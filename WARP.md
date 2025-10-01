# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

`daemonic` is a Go template/foundation for building robust daemon services. It implements a clean shutdown pattern with signal handling, structured logging, and graceful service lifecycle management.

## Key Architecture

The project follows a service interface pattern where all long-running services implement:
```go
type Service interface {
    Run(ctx context.Context) error
    Shutdown(ctx context.Context) error
}
```

### Core Components

- **Main application** (`cmd/daemonic/main.go`): Provides the daemon framework with signal handling, graceful shutdown, and structured logging
- **Service interface**: Abstract interface that all services must implement
- **Example service**: `TickerService` demonstrates the pattern with a simple ticker that logs timestamps

### Key Features

- **Signal handling**: Responds to SIGINT and SIGTERM for graceful shutdown
- **Structured logging**: Uses Go's `slog` package with JSON output
- **Context cancellation**: Proper context propagation for shutdown coordination
- **Shutdown timeout**: 30-second graceful shutdown window
- **Error handling**: Clean error propagation and exit codes

## Development Commands

### Build
```bash
go build ./cmd/daemonic
```

### Test
```bash
just test
```

### Lint and Format
```bash
# Format code
gofmt -w .

# Vet code
go vet ./...

# If golangci-lint is available
golangci-lint run
```

### Run
```bash
./daemonic
```

### Dependencies
```bash
go mod tidy
```

## Development Patterns

### Adding New Services

1. Implement the `Service` interface with `Run()` and `Shutdown()` methods
2. Replace `TickerService` initialization in `main()` with your service
3. Ensure proper context handling for graceful shutdown
4. Use structured logging with `slog` for consistency

### Service Implementation Guidelines

- **Run method**: Should block and handle the main service loop, responding to `ctx.Done()` for shutdown
- **Shutdown method**: Should clean up resources and complete within the shutdown timeout
- **Error handling**: Return errors from `Run()` for unrecoverable failures; log recoverable errors
- **Context awareness**: Always respect context cancellation signals

## Project Structure

- `cmd/daemonic/`: Main application entry point
- `int/`: Empty directory, likely intended for internal packages
- Binary output: `./daemonic` (gitignored: `./bin`)

## Notes

- This is a new repository with no commits yet
- Uses Go 1.24.2
- Minimal dependencies (only standard library)
- Designed as a template/starting point for daemon applications
