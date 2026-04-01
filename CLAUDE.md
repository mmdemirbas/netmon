# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
task              # Build binary for current platform → bin/
task build        # Build binary for current platform → bin/
task build-all    # Cross-compile: macOS arm64/amd64, Windows, Linux
task test         # Run Go tests
task lint         # Run golangci-lint
task quality      # Tests with race+coverage, then lint
task clean        # Remove build artifacts
task run          # Run in the foreground
task start        # Run in the background (logs to log/netmon.log)
task stop         # Stop the background process
task status       # Check if netmon is running
task install      # Install as system service (requires sudo)
task uninstall    # Uninstall system service (requires sudo)
```

Run a single test:
```bash
go test . -run TestName
```

## Architecture

**netmon** is a network latency and speed monitor. A single binary periodically runs speed tests, stores results in a local SQLite database, and serves a web UI with visual charts.

Single `package main`, flat file layout — no sub-packages.

### Core Files

| File | Responsibility |
|---|---|
| `main.go` | Entry point: CLI flags, service lifecycle (install/start/stop), program init |
| `collector.go` | Periodic speed test runner; writes results to the database |
| `database.go` | SQLite schema, init, and query functions |
| `server.go` | HTTP server: serves embedded static web UI and JSON API for chart data |
| `network-name.go` | Platform-specific active network interface name detection |

Static web UI lives in `static/` and is embedded into the binary via `//go:embed`.

### Data Flow

**Collection:** `startCollector` ticker fires → `collect()` → speedtest-go measures latency/download/upload → row inserted into SQLite.

**Display:** Browser → HTTP → `server.go` queries SQLite → JSON → Chart.js renders charts.

### Persistence

- `data/netmon.db` — SQLite database (configurable via `-db-file`)
- `bin/` — compiled binary
- `log/` — log file when running via `task start`

### Dependencies

- `kardianos/service` — cross-platform system service management (launchd/systemd/Windows Service)
- `mattn/go-sqlite3` — SQLite driver (**requires CGO**)
- `showwin/speedtest-go` — speed test client

### CGO Requirement

`go-sqlite3` requires CGO. `CGO_ENABLED=1` is the default and must not be disabled. Cross-compilation to a different OS requires the corresponding C cross-compiler.

### Configuration

All configuration via CLI flags — no config file:

| Flag | Default | Description |
|---|---|---|
| `-db-file` | `data/netmon.db` | SQLite database path |
| `-interval` | `5m` | Measurement interval |
| `-port` | `9898` | HTTP port |
| `-service` | — | Service control: `install`, `start`, `stop`, `uninstall` |

### Known Issues

- `network-name.go`: macOS SSID detection uses deprecated `airport` binary — marked `FIXME`.
- `main.go`: Suspected memory leak and missing graceful shutdown for the collector goroutine — marked `FIXME`/`TODO`.
- `main.go`: `fmt.Errorf` return value silently discarded in `Start` — errors from `run()` are lost.
