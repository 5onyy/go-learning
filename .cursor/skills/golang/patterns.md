# Go Patterns Reference

## Testing Patterns

### Test File Organization

- Test files live next to the code: `user.go` → `user_test.go`
- Use `package foo_test` for black-box testing (tests only exported API)
- Use `package foo` for white-box testing (tests can access unexported)

### Table-Driven Tests (Full Pattern)

```go
func TestParseURL(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *URL
        wantErr bool
    }{
        {
            name:  "valid http",
            input: "http://example.com",
            want:  &URL{Scheme: "http", Host: "example.com"},
        },
        {
            name:    "empty string",
            input:   "",
            wantErr: true,
        },
        {
            name:    "missing scheme",
            input:   "example.com",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseURL(tt.input)
            if tt.wantErr {
                if err == nil {
                    t.Fatal("expected error, got nil")
                }
                return
            }
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got %+v, want %+v", got, tt.want)
            }
        })
    }
}
```

### Test Helpers

Use `t.Helper()` so failure messages report the caller's line:

```go
func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

### Testdata Directory

Store test fixtures in a `testdata/` directory (ignored by `go build`):

```go
func TestProcess(t *testing.T) {
    input, err := os.ReadFile("testdata/input.json")
    if err != nil {
        t.Fatal(err)
    }
    // ...
}
```

### Golden Files

Compare output against "golden" reference files:

```go
func TestRender(t *testing.T) {
    got := render(input)

    golden := filepath.Join("testdata", t.Name()+".golden")

    if *update {
        os.WriteFile(golden, []byte(got), 0644)
        return
    }

    want, _ := os.ReadFile(golden)
    if got != string(want) {
        t.Errorf("output mismatch; run with -update to update golden files")
    }
}
```

### Testing HTTP Handlers

```go
func TestGetUser(t *testing.T) {
    req := httptest.NewRequest("GET", "/users/123", nil)
    w := httptest.NewRecorder()

    getUser(w, req)

    resp := w.Result()
    if resp.StatusCode != http.StatusOK {
        t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
    }

    body, _ := io.ReadAll(resp.Body)
    // assert body contents...
}
```

### Benchmarks

```go
func BenchmarkProcess(b *testing.B) {
    data := loadTestData()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        Process(data)
    }
}
```

Run: `go test -bench=BenchmarkProcess -benchmem ./...`

### Subtests for Setup/Teardown

```go
func TestDatabase(t *testing.T) {
    db := setupTestDB(t)

    t.Run("Create", func(t *testing.T) {
        err := db.Create(&User{Name: "Alice"})
        if err != nil {
            t.Fatal(err)
        }
    })

    t.Run("Read", func(t *testing.T) {
        user, err := db.Get("Alice")
        if err != nil {
            t.Fatal(err)
        }
        if user.Name != "Alice" {
            t.Errorf("got %q, want %q", user.Name, "Alice")
        }
    })
}
```

### t.Cleanup for Resource Teardown

```go
func setupTestDB(t *testing.T) *DB {
    db, err := NewDB(":memory:")
    if err != nil {
        t.Fatal(err)
    }
    t.Cleanup(func() {
        db.Close()
    })
    return db
}
```

---

## Concurrency Patterns

### Worker Pool

Process items concurrently with bounded parallelism:

```go
func processAll(items []Item, workers int) []Result {
    jobs := make(chan Item, len(items))
    results := make(chan Result, len(items))

    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for item := range jobs {
                results <- process(item)
            }
        }()
    }

    // Send jobs
    for _, item := range items {
        jobs <- item
    }
    close(jobs)

    // Wait and collect
    go func() {
        wg.Wait()
        close(results)
    }()

    var out []Result
    for r := range results {
        out = append(out, r)
    }
    return out
}
```

### Fan-Out / Fan-In

```go
func fanOut(ctx context.Context, input <-chan int, workers int) <-chan int {
    results := make(chan int)
    var wg sync.WaitGroup

    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for n := range input {
                select {
                case results <- process(n):
                case <-ctx.Done():
                    return
                }
            }
        }()
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}
```

### errgroup for Concurrent Error Handling

```go
import "golang.org/x/sync/errgroup"

func fetchAll(ctx context.Context, urls []string) ([]string, error) {
    g, ctx := errgroup.WithContext(ctx)
    results := make([]string, len(urls))

    for i, url := range urls {
        i, url := i, url
        g.Go(func() error {
            body, err := fetch(ctx, url)
            if err != nil {
                return err
            }
            results[i] = body
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}
```

### Select with Timeout

```go
select {
case result := <-ch:
    fmt.Println(result)
case <-time.After(5 * time.Second):
    fmt.Println("timed out")
case <-ctx.Done():
    fmt.Println("cancelled")
}
```

### Mutex for Shared State

```go
type SafeCounter struct {
    mu sync.Mutex
    v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.v[key]++
}

func (c *SafeCounter) Get(key string) int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.v[key]
}
```

### sync.Once for One-Time Initialization

```go
var (
    instance *DB
    once     sync.Once
)

func GetDB() *DB {
    once.Do(func() {
        instance = connectDB()
    })
    return instance
}
```

---

## HTTP Patterns

### JSON Request/Response

```go
func createUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }

    user, err := svc.CreateUser(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

### Middleware Chaining

```go
func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

// Usage
handler := chain(mux, logging, recovery, cors)
http.ListenAndServe(":8080", handler)
```

### Graceful Shutdown

```go
func main() {
    srv := &http.Server{Addr: ":8080", Handler: mux}

    go func() {
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("shutdown error: %v", err)
    }
}
```

### HTTP Client with Timeout

```go
client := &http.Client{
    Timeout: 10 * time.Second,
}

resp, err := client.Get("https://api.example.com/data")
if err != nil {
    return fmt.Errorf("fetching data: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
    return fmt.Errorf("unexpected status: %s", resp.Status)
}
```

---

## JSON Patterns

### Struct Tags

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    Password  string    `json:"-"` // never serialized
}
```

### Common Tag Options

| Tag | Effect |
|-----|--------|
| `json:"name"` | Use "name" as JSON key |
| `json:"name,omitempty"` | Omit if zero value |
| `json:"-"` | Always omit |
| `json:",string"` | Encode number as JSON string |

### Custom Marshal/Unmarshal

```go
type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
    return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
    var s string
    if err := json.Unmarshal(b, &s); err != nil {
        return err
    }
    dur, err := time.ParseDuration(s)
    if err != nil {
        return err
    }
    *d = Duration(dur)
    return nil
}
```

---

## File I/O Patterns

### Read Entire File

```go
data, err := os.ReadFile("config.json")
if err != nil {
    return fmt.Errorf("reading config: %w", err)
}
```

### Write File

```go
err := os.WriteFile("output.txt", []byte(content), 0644)
```

### Read Line by Line

```go
file, err := os.Open("data.txt")
if err != nil {
    return err
}
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    // process line
}
if err := scanner.Err(); err != nil {
    return fmt.Errorf("scanning: %w", err)
}
```

### io.Reader / io.Writer Composition

```go
// Chain readers/writers together
r := io.LimitReader(resp.Body, 1<<20) // limit to 1MB
data, err := io.ReadAll(r)

// Write to multiple destinations
w := io.MultiWriter(file, os.Stdout)
fmt.Fprintln(w, "logged to file and stdout")

// Copy between reader and writer
n, err := io.Copy(dst, src)
```

---

## Functional Options Pattern

For configurable constructors with clean API:

```go
type Server struct {
    port    int
    timeout time.Duration
    logger  *log.Logger
}

type Option func(*Server)

func WithPort(port int) Option {
    return func(s *Server) { s.port = port }
}

func WithTimeout(d time.Duration) Option {
    return func(s *Server) { s.timeout = d }
}

func NewServer(opts ...Option) *Server {
    s := &Server{
        port:    8080,
        timeout: 30 * time.Second,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
srv := NewServer(WithPort(9090), WithTimeout(10*time.Second))
```

---

## Dependency Injection

Prefer constructor injection via interfaces:

```go
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
}

type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    return s.repo.FindByID(ctx, id)
}
```

This makes testing easy — pass a mock implementation:

```go
type mockRepo struct {
    users map[string]*User
}

func (m *mockRepo) FindByID(_ context.Context, id string) (*User, error) {
    u, ok := m.users[id]
    if !ok {
        return nil, ErrNotFound
    }
    return u, nil
}

func TestGetUser(t *testing.T) {
    repo := &mockRepo{users: map[string]*User{"1": {Name: "Alice"}}}
    svc := NewUserService(repo)

    user, err := svc.GetUser(context.Background(), "1")
    // assert...
}
```
