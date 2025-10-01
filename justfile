# Default target shows available commands
default:
    @just --list

# Build the project
build *args="./cmd/daemonic":
    go build {{args}}

# Clean build artifacts
clean *args="./...":
    go clean {{args}}
    rm -f daemonic

# Run tests
test *args="./...":
    go test {{args}}

# Build and run the daemon
run *args="": build
    ./daemonic {{args}}

# Distribution target (no-op for now)
dist *args="./...":
    @echo "Distribution target not implemented yet (args: {{args}})"
