# Testing Patterns

## Table-Driven Tests

The standard pattern for testing multiple cases.

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 1, 2, 3},
        {"zero", 0, 0, 0},
        {"negative", -1, -2, -3},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

## Test Helpers

Use `t.Helper()` so failures report the caller's line number.

```go
func assertError(t *testing.T, got, want error) {
    t.Helper()
    if !errors.Is(got, want) {
        t.Errorf("got error %v, want %v", got, want)
    }
}

func assertEqual[T comparable](t *testing.T, got, want T) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

## HTTP Handler Testing with httptest

```go
func TestGetUser(t *testing.T) {
    handler := http.HandlerFunc(GetUserHandler)

    req := httptest.NewRequest("GET", "/users/123", nil)
    w := httptest.NewRecorder()

    handler.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
    }

    var user User
    if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
        t.Fatalf("decode response: %v", err)
    }
}
```

## Test Server for External APIs

```go
func TestFetchData(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, `{"id": 1, "name": "test"}`)
    }))
    defer server.Close()

    result, err := FetchData(server.URL)
    if err != nil {
        t.Fatalf("FetchData: %v", err)
    }
    assertEqual(t, result.Name, "test")
}
```

## Testdata Directory

Place test fixtures in a `testdata/` directory (ignored by `go build`).

```go
func TestParseConfig(t *testing.T) {
    data, err := os.ReadFile("testdata/config.json")
    if err != nil {
        t.Fatalf("read testdata: %v", err)
    }

    cfg, err := ParseConfig(data)
    if err != nil {
        t.Fatalf("ParseConfig: %v", err)
    }
    assertEqual(t, cfg.Port, 8080)
}
```

## Testing with Interfaces (Mocks)

```go
// Define a minimal interface in the test file
type mockStore struct {
    users map[string]*User
}

func (m *mockStore) GetUser(id string) (*User, error) {
    u, ok := m.users[id]
    if !ok {
        return nil, ErrNotFound
    }
    return u, nil
}

func TestService(t *testing.T) {
    store := &mockStore{
        users: map[string]*User{"1": {ID: "1", Name: "Alice"}},
    }
    svc := NewService(store)

    user, err := svc.GetUser("1")
    if err != nil {
        t.Fatalf("GetUser: %v", err)
    }
    assertEqual(t, user.Name, "Alice")
}
```

## Essential Test Commands

```bash
go test ./...              # Run all tests
go test -race ./...        # Detect races
go test -cover ./...       # Coverage summary
go test -coverprofile=c.out ./... && go tool cover -html=c.out  # HTML report
go test -run TestName ./pkg/...  # Run specific test
go test -v ./...           # Verbose output
go test -short ./...       # Skip long tests (check testing.Short())
go test -count=1 ./...     # Disable test caching
```
