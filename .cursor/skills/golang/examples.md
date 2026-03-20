# Go Code Examples

Complete, runnable examples for common Go tasks.

## Example 1: CLI Tool That Reads Files

A command-line tool that counts word frequency in files.

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
    "strings"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintf(os.Stderr, "usage: %s <file> [file...]\n", os.Args[0])
        os.Exit(1)
    }

    counts := make(map[string]int)
    for _, filename := range os.Args[1:] {
        if err := countWords(filename, counts); err != nil {
            fmt.Fprintf(os.Stderr, "error processing %s: %v\n", filename, err)
            continue
        }
    }

    printTopWords(counts, 10)
}

func countWords(filename string, counts map[string]int) error {
    f, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("opening file: %w", err)
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    scanner.Split(bufio.ScanWords)
    for scanner.Scan() {
        word := strings.ToLower(scanner.Text())
        counts[word]++
    }
    return scanner.Err()
}

func printTopWords(counts map[string]int, n int) {
    type wordCount struct {
        word  string
        count int
    }

    var wcs []wordCount
    for w, c := range counts {
        wcs = append(wcs, wordCount{w, c})
    }
    sort.Slice(wcs, func(i, j int) bool {
        return wcs[i].count > wcs[j].count
    })

    if n > len(wcs) {
        n = len(wcs)
    }
    for _, wc := range wcs[:n] {
        fmt.Printf("%6d %s\n", wc.count, wc.word)
    }
}
```

---

## Example 2: HTTP JSON API

A simple REST API with in-memory storage.

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"
)

type Todo struct {
    ID   string `json:"id"`
    Text string `json:"text"`
    Done bool   `json:"done"`
}

type TodoStore struct {
    mu    sync.RWMutex
    todos map[string]Todo
    next  int
}

func NewTodoStore() *TodoStore {
    return &TodoStore{todos: make(map[string]Todo)}
}

func (s *TodoStore) List() []Todo {
    s.mu.RLock()
    defer s.mu.RUnlock()
    result := make([]Todo, 0, len(s.todos))
    for _, t := range s.todos {
        result = append(result, t)
    }
    return result
}

func (s *TodoStore) Add(text string) Todo {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.next++
    id := fmt.Sprintf("%d", s.next)
    t := Todo{ID: id, Text: text}
    s.todos[id] = t
    return t
}

func main() {
    store := NewTodoStore()
    mux := http.NewServeMux()

    mux.HandleFunc("GET /todos", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(store.List())
    })

    mux.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Text string `json:"text"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid JSON", http.StatusBadRequest)
            return
        }
        todo := store.Add(req.Text)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(todo)
    })

    log.Println("listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

---

## Example 3: Concurrent URL Fetcher

Fetch multiple URLs concurrently and report results.

```go
package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

type FetchResult struct {
    URL      string
    Status   int
    Bytes    int64
    Duration time.Duration
    Err      error
}

func fetch(ctx context.Context, url string) FetchResult {
    start := time.Now()
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return FetchResult{URL: url, Err: err, Duration: time.Since(start)}
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return FetchResult{URL: url, Err: err, Duration: time.Since(start)}
    }
    defer resp.Body.Close()

    n, _ := io.Copy(io.Discard, resp.Body)
    return FetchResult{
        URL:      url,
        Status:   resp.StatusCode,
        Bytes:    n,
        Duration: time.Since(start),
    }
}

func fetchAll(urls []string, timeout time.Duration) []FetchResult {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    results := make(chan FetchResult, len(urls))

    for _, url := range urls {
        go func(url string) {
            results <- fetch(ctx, url)
        }(url)
    }

    var out []FetchResult
    for range urls {
        out = append(out, <-results)
    }
    return out
}

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintf(os.Stderr, "usage: %s <url> [url...]\n", os.Args[0])
        os.Exit(1)
    }

    results := fetchAll(os.Args[1:], 30*time.Second)
    for _, r := range results {
        if r.Err != nil {
            fmt.Printf("ERR  %s: %v\n", r.URL, r.Err)
        } else {
            fmt.Printf("%-4d %7d bytes  %v  %s\n", r.Status, r.Bytes, r.Duration.Round(time.Millisecond), r.URL)
        }
    }
}
```

---

## Example 4: Table-Driven Tests with Edge Cases

```go
package calc

import (
    "math"
    "testing"
)

func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

func TestDivide(t *testing.T) {
    tests := []struct {
        name    string
        a, b    float64
        want    float64
        wantErr bool
    }{
        {name: "simple", a: 10, b: 2, want: 5},
        {name: "fraction", a: 1, b: 3, want: 1.0 / 3.0},
        {name: "negative", a: -10, b: 2, want: -5},
        {name: "both negative", a: -10, b: -2, want: 5},
        {name: "zero numerator", a: 0, b: 5, want: 0},
        {name: "divide by zero", a: 10, b: 0, wantErr: true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.a, tt.b)
            if tt.wantErr {
                if err == nil {
                    t.Fatal("expected error, got nil")
                }
                return
            }
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if math.Abs(got-tt.want) > 1e-9 {
                t.Errorf("Divide(%g, %g) = %g, want %g", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

---

## Example 5: Reading and Writing JSON Config

```go
package config

import (
    "encoding/json"
    "fmt"
    "os"
)

type Config struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    LogLevel string `json:"log_level"`
    Database struct {
        DSN         string `json:"dsn"`
        MaxConns    int    `json:"max_conns"`
        IdleTimeout string `json:"idle_timeout"`
    } `json:"database"`
}

func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("reading config file: %w", err)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parsing config: %w", err)
    }

    if cfg.Port == 0 {
        cfg.Port = 8080
    }
    if cfg.LogLevel == "" {
        cfg.LogLevel = "info"
    }

    return &cfg, nil
}

func Save(path string, cfg *Config) error {
    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return fmt.Errorf("marshaling config: %w", err)
    }
    return os.WriteFile(path, data, 0644)
}
```

---

## Example 6: Middleware Stack

```go
package middleware

import (
    "log"
    "net/http"
    "time"
)

func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

func Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic: %v", err)
                http.Error(w, "internal server error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func CORS(origin string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusNoContent)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## Example 7: Using Context for Cancellation

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func slowOperation(ctx context.Context) (string, error) {
    select {
    case <-time.After(3 * time.Second):
        return "done", nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    result, err := slowOperation(ctx)
    if err != nil {
        fmt.Printf("operation failed: %v\n", err)
        return
    }
    fmt.Println(result)
}
```
