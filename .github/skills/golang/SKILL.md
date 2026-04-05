---
name: golang
description: "Idiomatic Go patterns and best practices. Use when: writing new Go code, designing Go packages, choosing between concurrency patterns, structuring error handling, or applying Go idioms."
argument-hint: "Describe the Go code you're writing (e.g., HTTP handler, worker pool, CLI tool)"
---

# Go Development Patterns

Guide for writing idiomatic, robust Go code. Covers principles, decision points, and links to detailed references.

## When to Use

- Writing new Go functions, types, or packages
- Choosing between concurrency approaches
- Structuring error handling
- Designing interfaces or package layouts

## Core Principles

| Principle | Rule |
|-----------|------|
| Simplicity | Clear and direct over clever. Code should be obvious |
| Zero values | Design types so zero value is immediately usable |
| Interfaces | Accept interfaces, return concrete structs |
| Errors | Treat errors as values; always wrap with context using `%w` |
| Return early | Handle errors first, keep the happy path unindented |
| Dependency injection | Pass dependencies via constructors, not package-level state |

## Procedure: Writing New Go Code

### 1. Package & File Setup

- Use short, lowercase package names — no underscores, no camelCase
- One package per directory; `package main` only in `cmd/` entry points
- Follow [standard project layout](./references/project-layout.md):
  `cmd/`, `internal/`, `pkg/`, `api/`, `testdata/`

### 2. Define Types

- Make zero values useful (avoid nil maps/slices that require init)
- Use functional options for configurable constructors (`WithTimeout`, `WithLogger`)
- Prefer composition via embedding over inheritance
- Pick **one** receiver type (value or pointer) per type — don't mix
- See [interfaces and structs reference](./references/interfaces-and-structs.md)

### 3. Error Handling

- Always wrap errors with context: `fmt.Errorf("load config %s: %w", path, err)`
- Use sentinel errors (`var ErrNotFound = errors.New(...)`) for expected conditions
- Use custom error types for errors carrying structured data
- Check with `errors.Is` / `errors.As`, never string comparison
- Never ignore errors with `_` unless explicitly documented why
- See [error handling reference](./references/error-handling.md)

### 4. Concurrency (if needed)

**Decision tree:**

```
Need concurrent work?
├── No → Don't use goroutines
└── Yes
    ├── Fan-out/fan-in with fixed workers → Worker pool (sync.WaitGroup)
    ├── Multiple independent tasks, fail on first error → errgroup
    ├── Pipeline stages → Channels with context cancellation
    ├── Need timeout/cancellation → context.WithTimeout + select
    ├── Limit concurrent access to a resource → Semaphore (chan struct{} or x/sync/semaphore)
    ├── Rate-limit outgoing calls → x/time/rate.Limiter or token-bucket channel
    └── Broadcast event to multiple listeners → sync.Cond or close(chan)
```

- Always pass `context.Context` as the first parameter
- Use buffered channels or `select` with `ctx.Done()` to prevent goroutine leaks
- See [concurrency reference](./references/concurrency.md)

### 5. Performance Considerations

- Preallocate slices when size is known: `make([]T, 0, n)`
- Use `strings.Builder` or `strings.Join` instead of `+=` in loops
- Use `sync.Pool` for frequently allocated short-lived objects
- See [performance reference](./references/performance.md)

### 6. Testing

- Use table-driven tests with `t.Run` for subtests
- Use `t.Helper()` in test helpers for correct line reporting
- Use `httptest` for handler and HTTP client testing
- Place fixtures in `testdata/` (ignored by `go build`)
- Mock dependencies via interfaces, not frameworks
- See [testing reference](./references/testing.md)

### 7. Format & Validate

```bash
gofmt -w .          # Format
go vet ./...        # Static analysis
go test -race ./... # Tests with race detector
```

- See [tooling reference](./references/tooling.md)

## Quick Reference: Go Idioms

| Idiom | Description |
|-------|-------------|
| Accept interfaces, return structs | Functions accept interface params, return concrete types |
| Errors are values | Treat errors as first-class values, not exceptions |
| Don't communicate by sharing memory | Use channels for coordination between goroutines |
| Make the zero value useful | Types should work without explicit initialization |
| A little copying > a little dependency | Avoid unnecessary external dependencies |
| Clear is better than clever | Prioritize readability over cleverness |
| gofmt is everyone's friend | Always format with gofmt/goimports |
| Return early | Handle errors first, keep happy path unindented |

## Anti-Patterns

- Naked returns in long functions — always name what you return
- `panic` for control flow — use error returns instead
- `context.Context` in structs — pass as first function parameter
- Mixing value and pointer receivers on the same type
- Package-level mutable state / `init()` for setup — use constructors
- Ignoring errors with `_` without a comment explaining why
