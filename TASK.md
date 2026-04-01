# TASK.md

Outstanding work items, in priority order.
Items resolved in previous sessions are omitted.

---

## P1 — Correctness / code quality

### Fix log message capitalization inconsistency

Several log messages start with a capital letter (pre-existing style), while
newer messages follow Go convention (lowercase). The mix causes visual noise and
makes grep/parsing harder. All should be lowercase.

Affected calls (non-exhaustive):
- `collector.go`: "Error getting network name", "Error testing server", "Error saving online data"
- `main.go`: "Failed to run service"
- `server.go`: "Error loading data", "Error encoding metrics", "Error writing response"

---

### Implement or remove the stale CLI-flags TODO

`main.go:46` — `// TODO: Honor cli flags during service installation or start maybe`

When installing the service via `-service install`, the configured `-port`,
`-interval`, and `-db-file` flags are not baked into the service definition.
The service always starts with defaults. Either pass them as service arguments
on install, or remove the TODO and document the limitation in README.

---

## P2 — Improvements

### Simplify `saveMetric` — remove unnecessary transaction

`saveMetric` wraps a single `INSERT` in `Begin / Prepare / Exec / Commit`.
A transaction adds no value for a single write and makes the code harder to
read. Replace with a direct `db.Exec` call.

---

### Suppress or handle logger return values

`service.Logger` methods return `error`, but every call site discards the
return value silently. Either suppress explicitly with `_` (acceptable for
logging) or add a golangci-lint `errcheck` exclusion so the linter does not
flag these in future.

---

### Expand test coverage

Current tests cover: `initDatabase`, `saveMetric`, `getMetricsSince`,
`closeDatabase`, and four `handleMetrics` HTTP scenarios.

Missing coverage:
- **`collector.go`** — `collect()` wires together network name lookup,
  speedtest, and DB write. Needs the external `speedtest.FetchServers` call
  abstracted behind an interface to be testable.
- **`network-name.go`** — platform branches are untestable as written because
  they exec OS binaries. Consider an injectable exec function or build-tag
  isolated unit tests that mock the command output.

---

## P3 — Documentation / polish

### Update README

The README is stale since the module rename and recent feature additions:
- Module path changed from `netmon` to `github.com/mmdemirbas/netmon`.
- `/metrics` now defaults to the last 24 hours (was: all time).
- UI has a time range selector (1 h / 6 h / 24 h / 7 d / 30 d / all).
- Chart.js is bundled in the binary — no CDN dependency.
- Development pre-requisites and command list need updating.
