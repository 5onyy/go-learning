# Go `net/http` Package Cheatsheet

> Source: [pkg.go.dev/net/http](https://pkg.go.dev/net/http)  
> The `net/http` package provides HTTP client and server implementations.

```go
import "net/http"
```

---

## Table of Contents

1. [Quick Start](#1-quick-start)
2. [HTTP Client — Making Requests](#2-http-client--making-requests)
3. [HTTP Server — Handling Requests](#3-http-server--handling-requests)
4. [Request & Response Objects](#4-request--response-objects)
5. [Headers](#5-headers)
6. [ServeMux — Routing (go1.22+)](#6-servemux--routing-go122)
7. [Middleware Pattern](#7-middleware-pattern)
8. [Serving Static Files](#8-serving-static-files)
9. [Cookies](#9-cookies)
10. [Custom Client & Transport](#10-custom-client--transport)
11. [Graceful Shutdown](#11-graceful-shutdown)
12. [Status Codes & Method Constants](#12-status-codes--method-constants)
13. [Quick Reference Table](#13-quick-reference-table)

---

## 1. Quick Start

### Minimal Server

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, World!")
    })
    http.ListenAndServe(":8080", nil)
}
```

### Minimal Client

```go
resp, err := http.Get("https://example.com")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
fmt.Println(string(body))
```

---

## 2. HTTP Client — Making Requests

### `http.Get`

Simple GET request using the default client.

```go
resp, err := http.Get("https://api.example.com/users")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

fmt.Println("Status:", resp.StatusCode)
body, _ := io.ReadAll(resp.Body)
fmt.Println(string(body))
```

### `http.Post`

POST with a body and content type.

```go
jsonBody := `{"name": "Alice", "age": 30}`
resp, err := http.Post(
    "https://api.example.com/users",
    "application/json",
    strings.NewReader(jsonBody),
)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
fmt.Println(string(body))
```

### `http.PostForm`

POST with form-encoded data.

```go
resp, err := http.PostForm("https://example.com/login", url.Values{
    "username": {"alice"},
    "password": {"secret123"},
})
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()
```

### `http.NewRequest` + `Client.Do`

Full control — custom method, headers, body.

```go
// Create request
body := strings.NewReader(`{"title": "Go Book"}`)
req, err := http.NewRequest(http.MethodPut, "https://api.example.com/books/1", body)
if err != nil {
    log.Fatal(err)
}

// Set headers
req.Header.Set("Content-Type", "application/json")
req.Header.Set("Authorization", "Bearer my-token")

// Send
client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

fmt.Println("Status:", resp.Status)
```

### `http.NewRequestWithContext` — with Timeout/Cancel

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.example.com/slow", nil)
resp, err := http.DefaultClient.Do(req)
if err != nil {
    log.Fatal(err)  // returns error if timeout exceeded
}
defer resp.Body.Close()
```

### `http.Head`

Fetch only headers (no body).

```go
resp, err := http.Head("https://example.com/file.zip")
if err != nil {
    log.Fatal(err)
}
fmt.Println("Content-Length:", resp.Header.Get("Content-Length"))
```

---

## 3. HTTP Server — Handling Requests

### `http.HandleFunc` — Register Handler Function

```go
http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Query().Get("name"))
})
```

### `http.Handle` — Register Handler Interface

```go
type healthHandler struct{}

func (h healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, `{"status": "ok"}`)
}

http.Handle("/health", healthHandler{})
```

### `http.ListenAndServe`

Start HTTP server on given address.

```go
// nil means use http.DefaultServeMux
log.Fatal(http.ListenAndServe(":8080", nil))
```

### `http.ListenAndServeTLS`

Start HTTPS server.

```go
log.Fatal(http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil))
```

### `http.Server` — Full Control

```go
srv := &http.Server{
    Addr:              ":8080",
    Handler:           mux,              // your ServeMux or router
    ReadTimeout:       10 * time.Second,
    ReadHeaderTimeout: 5 * time.Second,
    WriteTimeout:      15 * time.Second,
    IdleTimeout:       60 * time.Second,
    MaxHeaderBytes:    1 << 20,          // 1 MB
}
log.Fatal(srv.ListenAndServe())
```

---

## 4. Request & Response Objects

### `http.Request` — Key Fields

```go
func handler(w http.ResponseWriter, r *http.Request) {
    r.Method          // "GET", "POST", etc.
    r.URL.Path        // "/api/users"
    r.URL.Query()     // url.Values — query params
    r.Header          // http.Header (map[string][]string)
    r.Body            // io.ReadCloser — request body
    r.Host            // "example.com"
    r.RemoteAddr      // "192.168.1.1:54321"
    r.ContentLength   // body length
    r.Context()       // context.Context (cancelled when client disconnects)
}
```

### Reading Request Body (JSON)

```go
func createUser(w http.ResponseWriter, r *http.Request) {
    var user struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }

    // Limit body size to prevent abuse
    r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB

    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }
    fmt.Fprintf(w, "Created: %s", user.Name)
}
```

### Reading Form Data

```go
func loginHandler(w http.ResponseWriter, r *http.Request) {
    // For POST forms
    r.ParseForm()
    username := r.FormValue("username")
    password := r.FormValue("password")

    // For query params
    page := r.URL.Query().Get("page")
}
```

### File Upload

```go
func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // Max 10 MB in memory
    r.ParseMultipartForm(10 << 20)

    file, header, err := r.FormFile("avatar")
    if err != nil {
        http.Error(w, "No file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    fmt.Fprintf(w, "Uploaded: %s (%d bytes)", header.Filename, header.Size)

    // Save to disk
    dst, _ := os.Create("/uploads/" + header.Filename)
    defer dst.Close()
    io.Copy(dst, file)
}
```

### Writing Responses

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Plain text
    fmt.Fprintln(w, "Hello!")

    // Set status code (must call before Write)
    w.WriteHeader(http.StatusCreated)  // 201

    // JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

    // Error response
    http.Error(w, "Something went wrong", http.StatusInternalServerError)

    // Redirect
    http.Redirect(w, r, "/new-url", http.StatusFound) // 302
}
```

### `http.Response` — Key Fields (Client Side)

```go
resp, _ := http.Get(url)

resp.StatusCode     // 200
resp.Status         // "200 OK"
resp.Header         // http.Header
resp.Body           // io.ReadCloser — MUST close
resp.ContentLength  // body length (-1 if unknown)
resp.Cookies()      // []*http.Cookie
resp.Location()     // *url.URL (from Location header)
```

---

## 5. Headers

```go
// Header is just map[string][]string
type Header map[string][]string
```

### Server — Set Response Headers

```go
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")  // replace
    w.Header().Add("X-Custom", "value1")                // append
    w.Header().Add("X-Custom", "value2")                // now has 2 values
    w.Header().Del("X-Custom")                          // remove all values
}
```

### Client — Set Request Headers

```go
req, _ := http.NewRequest("GET", url, nil)
req.Header.Set("Authorization", "Bearer token123")
req.Header.Set("Accept", "application/json")
req.Header.Set("User-Agent", "MyApp/1.0")
```

### Read Headers

```go
// Get first value
ct := r.Header.Get("Content-Type")       // "application/json"

// Get all values
vals := r.Header.Values("Accept")         // ["text/html", "application/json"]

// Check existence
if r.Header.Get("Authorization") == "" {
    http.Error(w, "Unauthorized", 401)
}
```

---

## 6. ServeMux — Routing (go1.22+)

Go 1.22 introduced powerful pattern matching with method and path wildcards.

### Basic Patterns

```go
mux := http.NewServeMux()

// Exact path
mux.HandleFunc("/", homeHandler)

// Path prefix (trailing slash matches subtree)
mux.HandleFunc("/api/", apiHandler)

// Method + path
mux.HandleFunc("GET /users", listUsers)
mux.HandleFunc("POST /users", createUser)

// Path wildcards
mux.HandleFunc("GET /users/{id}", getUser)
mux.HandleFunc("DELETE /users/{id}", deleteUser)

// Catch-all wildcard
mux.HandleFunc("GET /files/{path...}", serveFile)

// Exact path only (not subtree)
mux.HandleFunc("GET /{$}", rootOnly)  // matches "/" but not "/foo"

// Host-specific
mux.HandleFunc("api.example.com/v1/", apiV1Handler)

http.ListenAndServe(":8080", mux)
```

### `r.PathValue` — Extract Wildcards

```go
mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "User ID: %s", id)
})

mux.HandleFunc("GET /files/{path...}", func(w http.ResponseWriter, r *http.Request) {
    filepath := r.PathValue("path")
    fmt.Fprintf(w, "File: %s", filepath)
})
```

### Complete REST API Example

```go
mux := http.NewServeMux()

mux.HandleFunc("GET /api/posts",       listPosts)
mux.HandleFunc("POST /api/posts",      createPost)
mux.HandleFunc("GET /api/posts/{id}",  getPost)
mux.HandleFunc("PUT /api/posts/{id}",  updatePost)
mux.HandleFunc("DELETE /api/posts/{id}", deletePost)

srv := &http.Server{
    Addr:    ":8080",
    Handler: mux,
}
log.Fatal(srv.ListenAndServe())
```

---

## 7. Middleware Pattern

Middleware wraps a handler to add cross-cutting behavior (logging, auth, etc.).

### Writing Middleware

```go
// Middleware signature: takes Handler, returns Handler
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### Applying Middleware

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /api/data", dataHandler)

// Chain: logging → auth → mux
handler := loggingMiddleware(authMiddleware(mux))
http.ListenAndServe(":8080", handler)
```

### Built-in: `http.StripPrefix`

```go
// Strip /static/ prefix before passing to file server
fs := http.FileServer(http.Dir("./public"))
http.Handle("/static/", http.StripPrefix("/static/", fs))
```

### Built-in: `http.TimeoutHandler`

```go
slowHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    time.Sleep(10 * time.Second)
    fmt.Fprintln(w, "done")
})

// Limit handler execution to 3 seconds
http.Handle("/slow", http.TimeoutHandler(slowHandler, 3*time.Second, "timeout!"))
```

### Built-in: `http.MaxBytesHandler` *(go1.18+)*

```go
// Limit request body to 1 MB for all requests to this handler
http.Handle("/upload", http.MaxBytesHandler(uploadHandler, 1<<20))
```

---

## 8. Serving Static Files

### `http.FileServer`

```go
// Serve files from ./static at /static/
fs := http.FileServer(http.Dir("./static"))
http.Handle("/static/", http.StripPrefix("/static/", fs))

// Serve files from current directory at root
http.Handle("/", http.FileServer(http.Dir(".")))
```

### `http.ServeFile` — Single File

```go
http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/favicon.ico")
})
```

### `http.FileServerFS` *(go1.22+)* — Using `fs.FS`

```go
//go:embed static
var staticFS embed.FS

http.Handle("/", http.FileServerFS(staticFS))
```

---

## 9. Cookies

### Set Cookie

```go
func loginHandler(w http.ResponseWriter, r *http.Request) {
    http.SetCookie(w, &http.Cookie{
        Name:     "session_id",
        Value:    "abc123",
        Path:     "/",
        MaxAge:   3600,            // 1 hour
        HttpOnly: true,            // no JS access
        Secure:   true,            // HTTPS only
        SameSite: http.SameSiteLaxMode,
    })
}
```

### Read Cookie

```go
func handler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_id")
    if err == http.ErrNoCookie {
        http.Error(w, "No session", http.StatusUnauthorized)
        return
    }
    fmt.Println("Session:", cookie.Value)

    // Get all cookies
    for _, c := range r.Cookies() {
        fmt.Println(c.Name, c.Value)
    }
}
```

### Delete Cookie

```go
http.SetCookie(w, &http.Cookie{
    Name:   "session_id",
    Value:  "",
    Path:   "/",
    MaxAge: -1,  // delete immediately
})
```

---

## 10. Custom Client & Transport

### `http.Client` with Timeout

```go
client := &http.Client{
    Timeout: 10 * time.Second,  // total timeout (connect + read)
}

resp, err := client.Get("https://api.example.com/data")
```

### Custom `http.Transport`

```go
transport := &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    IdleConnTimeout:     90 * time.Second,
    TLSHandshakeTimeout: 10 * time.Second,
    DisableCompression:  false,
}

client := &http.Client{
    Transport: transport,
    Timeout:   30 * time.Second,
}
```

### Reuse Clients (Important!)

```go
// GOOD: create once, reuse everywhere
var httpClient = &http.Client{Timeout: 10 * time.Second}

func fetchData(url string) ([]byte, error) {
    resp, err := httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}
```

### Always Close & Drain Response Body

```go
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()

// Even if you don't need the body — drain it for connection reuse
io.Copy(io.Discard, resp.Body)
```

---

## 11. Graceful Shutdown

```go
srv := &http.Server{Addr: ":8080", Handler: mux}

// Start server in a goroutine
go func() {
    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("listen: %v", err)
    }
}()

// Wait for interrupt signal
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
log.Println("Shutting down...")

// Give active connections 5 seconds to finish
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := srv.Shutdown(ctx); err != nil {
    log.Fatalf("forced shutdown: %v", err)
}
log.Println("Server exited")
```

---

## 12. Status Codes & Method Constants

### HTTP Methods

```go
http.MethodGet     // "GET"
http.MethodPost    // "POST"
http.MethodPut     // "PUT"
http.MethodPatch   // "PATCH"
http.MethodDelete  // "DELETE"
http.MethodHead    // "HEAD"
http.MethodOptions // "OPTIONS"
```

### Common Status Codes

```go
// 2xx Success
http.StatusOK                  // 200
http.StatusCreated             // 201
http.StatusNoContent           // 204

// 3xx Redirection
http.StatusMovedPermanently    // 301
http.StatusFound               // 302
http.StatusTemporaryRedirect   // 307
http.StatusPermanentRedirect   // 308

// 4xx Client Error
http.StatusBadRequest          // 400
http.StatusUnauthorized        // 401
http.StatusForbidden           // 403
http.StatusNotFound            // 404
http.StatusMethodNotAllowed    // 405
http.StatusConflict            // 409
http.StatusUnprocessableEntity // 422
http.StatusTooManyRequests     // 429

// 5xx Server Error
http.StatusInternalServerError // 500
http.StatusBadGateway          // 502
http.StatusServiceUnavailable  // 503
http.StatusGatewayTimeout      // 504
```

### `http.StatusText`

```go
fmt.Println(http.StatusText(404))  // "Not Found"
fmt.Println(http.StatusText(201))  // "Created"
```

### `http.Error` — Quick Error Response

```go
http.Error(w, "resource not found", http.StatusNotFound)
http.Error(w, "internal error", http.StatusInternalServerError)
```

### `http.Redirect`

```go
http.Redirect(w, r, "/login", http.StatusFound)            // 302
http.Redirect(w, r, "/new", http.StatusMovedPermanently)    // 301
```

---

## 13. Quick Reference Table

### Client Side

| Function / Method                | Use case                        |
|----------------------------------|---------------------------------|
| `http.Get(url)`                  | Simple GET request              |
| `http.Post(url, ct, body)`       | POST with body + content type   |
| `http.PostForm(url, values)`     | POST form-encoded data          |
| `http.Head(url)`                 | Get headers only                |
| `http.NewRequest(method, url, body)` | Build custom request        |
| `http.NewRequestWithContext(...)` | Request with timeout/cancel    |
| `client.Do(req)`                 | Execute custom request          |

### Server Side

| Function / Type                  | Use case                         |
|----------------------------------|----------------------------------|
| `http.HandleFunc(pattern, fn)`   | Register handler function        |
| `http.Handle(pattern, handler)`  | Register `Handler` interface     |
| `http.ListenAndServe(addr, h)`   | Start HTTP server                |
| `http.ListenAndServeTLS(...)`    | Start HTTPS server               |
| `http.NewServeMux()`             | Create custom router             |
| `http.Server{}`                  | Server with full configuration   |
| `srv.Shutdown(ctx)`              | Graceful shutdown                |
| `srv.Close()`                    | Immediate shutdown               |

### Response Helpers

| Function                         | Use case                         |
|----------------------------------|----------------------------------|
| `http.Error(w, msg, code)`       | Send error response              |
| `http.Redirect(w, r, url, code)` | Redirect client                 |
| `http.NotFound(w, r)`            | Reply with 404                   |
| `http.ServeFile(w, r, path)`     | Serve a single file              |
| `http.SetCookie(w, cookie)`      | Set a response cookie            |
| `http.StatusText(code)`          | "OK", "Not Found", etc.          |
| `http.DetectContentType(data)`   | Sniff MIME type from bytes       |

### Handler Wrappers

| Function                         | Use case                         |
|----------------------------------|----------------------------------|
| `http.StripPrefix(prefix, h)`    | Remove URL prefix                |
| `http.TimeoutHandler(h, d, msg)` | Limit handler execution time     |
| `http.MaxBytesHandler(h, n)`     | Limit request body size          |
| `http.FileServer(root)`          | Serve directory listing          |
| `http.RedirectHandler(url, code)`| Always redirect to URL           |

### Key Interfaces

```go
// Handler — implement this to handle requests
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// HandlerFunc — adapter to use plain functions as handlers
type HandlerFunc func(ResponseWriter, *Request)

// ResponseWriter — write response back to client
type ResponseWriter interface {
    Header() Header
    Write([]byte) (int, error)
    WriteHeader(statusCode int)
}
```
