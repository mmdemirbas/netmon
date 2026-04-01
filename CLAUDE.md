# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
just build           # build for current platform → bin/netmon
just build-all       # cross-compile for linux/darwin/windows (amd64, arm64)
just run             # build must precede run; runs with -interval 30s override
just stop clean build run  # full dev cycle
```

CGO is required (`go-sqlite3` uses cgo). Set `CGO_ENABLED=1` if the environment disables it.

There are no tests in this project yet.

## Architecture

Single `package main`, flat file layout — no sub-packages:

| File | Responsibility |
|---|---|
| `main.go` | CLI flags, `kardianos/service` lifecycle, wires collector + server |
| `collector.go` | Periodic ticker → `speedtest-go` → `saveMetric` |
| `database.go` | SQLite init, `Metrics` struct, `saveMetric` / `getAllMetrics` |
| `server.go` | HTTP server; embeds `static/` at compile time; `/metrics` JSON endpoint |
| `network-name.go` | OS-specific Wi-Fi SSID detection via shell commands |

**Data flow:** `startCollector` ticks → `collect()` runs a speed test → `saveMetric` writes a row to SQLite → `handleMetrics` reads all rows and serves them as JSON → `static/script.js` renders Chart.js charts.

**Service integration:** `kardianos/service` wraps the program so it can run interactively or as a system service (launchd/systemd/Windows Service) without code changes.

**Static assets** (`static/index.html`, `static/script.js`) are embedded into the binary via `//go:embed static`.

## Known Issues (tracked in code)

- `network-name.go`: macOS SSID detection is broken (uses deprecated `airport` binary). Marked `FIXME`.
- `main.go`: Suspected memory leak and missing graceful shutdown for the collector goroutine. Marked `FIXME` / `TODO`.
- `main.go:119`: `fmt.Errorf` return value silently discarded in `Start` — errors from `run()` are lost.
