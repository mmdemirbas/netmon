# Contributing to netmon

## Getting Started

1. Install [Go](https://go.dev/) 1.26+ and [Task](https://taskfile.dev/).
2. Clone the repo and build:

```bash
git clone https://github.com/mmdemirbas/netmon.git
cd netmon
task build
```

**Note:** go-sqlite3 requires CGO. Make sure a C compiler (GCC or Clang) is available.

## Development Workflow

```bash
task build      # Build to bin/netmon
task test       # Run all tests
task lint       # Run golangci-lint
task quality    # Tests with race+coverage, then lint
task run        # Run in foreground
task clean      # Remove build artifacts
```

## Submitting Changes

1. Create a branch from `main`.
2. Make your changes. Keep the scope focused — one feature or fix per PR.
3. Ensure `task quality` passes.
4. Open a pull request with a clear description.

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Keep functions focused and small
- Add tests for new functionality
- Use table-driven tests where appropriate

## Reporting Issues

Open an issue with:
- What you did (command, flags, input)
- What happened vs. what you expected
- OS, Go version (`go version`), and C compiler version (`gcc --version`)
