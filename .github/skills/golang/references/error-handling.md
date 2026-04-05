# Error Handling Patterns

## Error Wrapping with Context

Always add context when wrapping errors so the call chain is visible in logs.

```go
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("load config %s: %w", path, err)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parse config %s: %w", path, err)
    }

    return &cfg, nil
}
```

## Custom Error Types

Use when errors need to carry structured data beyond a message.

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}
```

## Sentinel Errors

Use for expected, well-known conditions that callers check against.

```go
var (
    ErrNotFound     = errors.New("resource not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrInvalidInput = errors.New("invalid input")
)
```

## Checking Errors with errors.Is and errors.As

```go
func HandleError(err error) {
    // Check for specific sentinel error
    if errors.Is(err, sql.ErrNoRows) {
        log.Println("No records found")
        return
    }

    // Check for error type
    var validationErr *ValidationError
    if errors.As(err, &validationErr) {
        log.Printf("Validation error on field %s: %s",
            validationErr.Field, validationErr.Message)
        return
    }

    log.Printf("Unexpected error: %v", err)
}
```

## Never Ignore Errors

```go
// Bad: Ignoring error with blank identifier
result, _ := doSomething()

// Good: Handle or propagate
result, err := doSomething()
if err != nil {
    return err
}

// Acceptable: Best-effort cleanup with comment
_ = writer.Close() // Best-effort cleanup, error logged elsewhere
```
