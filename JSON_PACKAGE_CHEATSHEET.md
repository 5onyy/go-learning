# Go `encoding/json` Package — Complete Cheatsheet

> From struct tags to HTTP APIs — everything about working with JSON in Go.

```go
import "encoding/json"
```

---

## Table of Contents

**Part 1 — Fundamentals**
1. [JSON ↔ Go Type Mapping](#1-json--go-type-mapping)
2. [Struct Tags](#2-struct-tags)
3. [Marshal (Go → JSON)](#3-marshal-go--json)
4. [Unmarshal (JSON → Go)](#4-unmarshal-json--go)
5. [How Field Matching Works](#5-how-field-matching-works)
6. [Handling Unknown / Dynamic JSON](#6-handling-unknown--dynamic-json)
7. [json.RawMessage — Defer Decoding](#7-jsonrawmessage--defer-decoding)

**Part 2 — Streaming with Encoder / Decoder**
8. [json.Decoder — Read from a Stream](#8-jsondecoder--read-from-a-stream)
9. [json.Encoder — Write to a Stream](#9-jsonencoder--write-to-a-stream)
10. [Decoder vs Unmarshal — When to Use Which](#10-decoder-vs-unmarshal--when-to-use-which)

**Part 3 — HTTP Use Cases**
11. [GET — Fetch JSON from an API](#11-get--fetch-json-from-an-api)
12. [POST — Send JSON to an API](#12-post--send-json-to-an-api)
13. [Building a JSON REST Server](#13-building-a-json-rest-server)
14. [Error Responses](#14-error-responses)
15. [Real-World Example — GitHub API Client](#15-real-world-example--github-api-client)

**Part 4 — Advanced**
16. [Custom Marshal / Unmarshal](#16-custom-marshal--unmarshal)
17. [Nested & Embedded Structs](#17-nested--embedded-structs)
18. [Common Gotchas](#18-common-gotchas)

---

# Part 1 — Fundamentals

## 1. JSON ↔ Go Type Mapping

| JSON type | Go type (Marshal) | Go type (Unmarshal into `any`) |
|---|---|---|
| `"string"` | `string` | `string` |
| `123` | `int`, `float64`, etc. | `float64` (always!) |
| `true` / `false` | `bool` | `bool` |
| `null` | `nil` pointer, nil slice/map | `nil` |
| `[1, 2, 3]` | `[]T`, `[N]T` | `[]any` |
| `{"k": "v"}` | `struct`, `map[string]T` | `map[string]any` |

The important surprise: **all JSON numbers become `float64`** when decoded into `any`/`map[string]any`.

---

## 2. Struct Tags

Struct tags control how fields map to/from JSON. They are read by `encoding/json` via reflection at runtime.

### Syntax

```go
type User struct {
    Name     string    `json:"name"`
    Email    string    `json:"email,omitempty"`
    Age      int       `json:"age"`
    Password string    `json:"-"`
    IsAdmin  bool      `json:"is_admin,omitempty"`
    Score    float64   `json:"score,string"`
    JoinDate time.Time `json:"join_date"`
}
```

### All Tag Options

| Tag | Effect on Marshal | Effect on Unmarshal |
|---|---|---|
| `json:"name"` | JSON key = `"name"` | Matches JSON key `"name"` |
| `json:"-"` | Always excluded | Always ignored |
| `json:"name,omitempty"` | Excluded if zero value | Normal matching |
| `json:",omitempty"` | Uses field name as key, omit if zero | Normal matching |
| `json:",string"` | Encodes number/bool as JSON string `"123"` | Decodes JSON string `"123"` into number |
| No tag | Uses exact field name as key | Case-insensitive match on field name |

### What Counts as "Zero Value" for `omitempty`

| Type | Zero value (omitted) |
|---|---|
| `string` | `""` |
| `int`, `float64` | `0`, `0.0` |
| `bool` | `false` |
| Pointer | `nil` |
| Slice, Map | `nil` (NOT empty `[]` or `{}`) |
| Struct | Never omitted (even if all fields are zero) |

Gotcha: `omitempty` does **not** omit empty slices `[]T{}` or empty maps — only `nil` ones. And it **never** omits a struct field, even if the struct itself is all zeros.

---

## 3. Marshal (Go → JSON)

### Basic

```go
type Product struct {
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

p := Product{Name: "Widget", Price: 9.99}

data, err := json.Marshal(p)
fmt.Println(string(data))
// {"name":"Widget","price":9.99}
```

### Pretty Print

```go
data, err := json.MarshalIndent(p, "", "  ")
// {
//   "name": "Widget",
//   "price": 9.99
// }

// First arg: the value
// Second arg: prefix for each line (usually "")
// Third arg: indent string (usually "  " or "\t")
```

### Marshal a Map

```go
m := map[string]any{
    "name": "Alice",
    "age":  30,
    "tags": []string{"admin", "user"},
}
data, _ := json.Marshal(m)
// {"age":30,"name":"Alice","tags":["admin","user"]}
// Note: map keys are sorted alphabetically in output
```

### Marshal a Slice

```go
items := []string{"apple", "banana", "cherry"}
data, _ := json.Marshal(items)
// ["apple","banana","cherry"]
```

### What Can Be Marshaled

| Type | Marshals to |
|---|---|
| `struct` | `{}` object |
| `map[string]T` | `{}` object (keys sorted) |
| `[]T`, `[N]T` | `[]` array |
| `string` | `"string"` |
| `int`, `float64` | number |
| `bool` | `true` / `false` |
| `nil` pointer | `null` |
| `nil` slice/map | `null` |
| `[]byte` | Base64-encoded string |
| `time.Time` | RFC 3339 string `"2026-03-24T..."` |

---

## 4. Unmarshal (JSON → Go)

### Into a Struct

```go
jsonStr := `{"name":"Gadget","price":19.99}`

var p Product
err := json.Unmarshal([]byte(jsonStr), &p)

fmt.Println(p.Name)    // "Gadget"
fmt.Println(p.Price)   // 19.99
```

### Into a Map (When You Don't Know the Type)

```go
var m map[string]any
json.Unmarshal([]byte(jsonStr), &m)

fmt.Println(m["name"])    // "Gadget" (string)
fmt.Println(m["price"])   // 19.99 (float64 — always float64!)
```

### Into `any` (Fully Generic)

```go
var v any
json.Unmarshal([]byte(`[1, "two", true, null]`), &v)

arr := v.([]any)
fmt.Println(arr[0])    // float64(1)
fmt.Println(arr[1])    // "two"
fmt.Println(arr[2])    // true
fmt.Println(arr[3])    // nil
```

### Partial Decode — Only Extract What You Need

```go
// GitHub returns 30+ fields per issue, but we only care about a few
type Issue struct {
    Number int       `json:"number"`
    Title  string    `json:"title"`
    State  string    `json:"state"`
    // All other JSON fields are silently ignored
}
```

---

## 5. How Field Matching Works

When decoding JSON into a struct, `encoding/json` matches each JSON key to a struct field using this priority:

```
1. Exact json tag match          →  `json:"html_url"` matches "html_url"
2. Exact field name match        →  Title matches "Title"
3. Case-insensitive field name   →  Title matches "title", "TITLE", "tItLe"
4. No match found                →  JSON field is silently IGNORED
```

Missing JSON fields leave the struct field at its zero value. No error.

```go
type User struct {
    Name  string `json:"name"`     // priority 1: matches "name"
    Email string                    // priority 2/3: matches "Email", "email", "EMAIL"
    Age   int    `json:"age"`
}

// This JSON has extra fields and is missing "age":
jsonStr := `{"name":"Alice","email":"a@b.com","role":"admin","active":true}`

var u User
json.Unmarshal([]byte(jsonStr), &u)
// u.Name  = "Alice"
// u.Email = "a@b.com"
// u.Age   = 0          ← zero value (missing from JSON)
// "role" and "active"  ← silently ignored
```

---

## 6. Handling Unknown / Dynamic JSON

### Decode into `map[string]any`

```go
var data map[string]any
json.Unmarshal(rawJSON, &data)

// Access with type assertions
name := data["name"].(string)
age := int(data["age"].(float64))   // numbers are always float64!
```

### Preserve Integer Precision with `UseNumber`

```go
dec := json.NewDecoder(bytes.NewReader(rawJSON))
dec.UseNumber()

var data map[string]any
dec.Decode(&data)

age := data["age"].(json.Number)
n, _ := age.Int64()       // 30 (precise integer)
f, _ := age.Float64()     // 30.0
s := age.String()         // "30"
```

### Reject Unknown Fields (Strict Mode)

```go
dec := json.NewDecoder(bytes.NewReader(rawJSON))
dec.DisallowUnknownFields()

var u User
err := dec.Decode(&u)
// Returns error if JSON has keys not matching any struct field
```

---

## 7. json.RawMessage — Defer Decoding

`json.RawMessage` stores raw JSON bytes without parsing. Decode the outer structure first, then decode the inner payload based on a discriminator field.

```go
type Event struct {
    Type    string          `json:"type"`
    Payload json.RawMessage `json:"payload"`
}

rawJSON := `{"type":"user_created","payload":{"name":"Alice","age":30}}`

var event Event
json.Unmarshal([]byte(rawJSON), &event)
// event.Type = "user_created"
// event.Payload = []byte(`{"name":"Alice","age":30}`)  ← still raw bytes

switch event.Type {
case "user_created":
    var u User
    json.Unmarshal(event.Payload, &u)
case "order_placed":
    var o Order
    json.Unmarshal(event.Payload, &o)
}
```

Also useful for **re-marshaling** — preserving JSON you don't need to inspect:

```go
type Wrapper struct {
    ID      int             `json:"id"`
    Config  json.RawMessage `json:"config"`   // pass through unchanged
}
```

---

# Part 2 — Streaming with Encoder / Decoder

## 8. json.Decoder — Read from a Stream

`json.Decoder` reads JSON from any `io.Reader` (HTTP body, file, stdin) **without buffering the entire input into memory first**.

```go
// From HTTP response
resp, _ := http.Get("https://api.example.com/users")
defer resp.Body.Close()

var users []User
json.NewDecoder(resp.Body).Decode(&users)
```

```go
// From a file
f, _ := os.Open("data.json")
defer f.Close()

var data Config
json.NewDecoder(f).Decode(&data)
```

```go
// From stdin
var input Command
json.NewDecoder(os.Stdin).Decode(&input)
```

### Decoder Options

```go
dec := json.NewDecoder(reader)
dec.UseNumber()                  // numbers stay as json.Number, not float64
dec.DisallowUnknownFields()      // error on unrecognized JSON keys

var v MyStruct
err := dec.Decode(&v)
```

### Decode Multiple JSON Values (NDJSON / JSON Lines)

```go
dec := json.NewDecoder(reader)

for dec.More() {
    var record Record
    if err := dec.Decode(&record); err != nil {
        break
    }
    process(record)
}
```

---

## 9. json.Encoder — Write to a Stream

`json.Encoder` writes JSON directly to any `io.Writer`.

```go
// To stdout
json.NewEncoder(os.Stdout).Encode(user)

// To a file
f, _ := os.Create("output.json")
defer f.Close()
json.NewEncoder(f).Encode(data)

// To HTTP response
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

### Encoder Options

```go
enc := json.NewEncoder(w)
enc.SetIndent("", "  ")         // pretty print
enc.SetEscapeHTML(false)         // don't escape <, >, & in strings
enc.Encode(v)
```

Note: `Encode` appends a `\n` after the JSON. `Marshal` does not.

---

## 10. Decoder vs Unmarshal — When to Use Which

```
json.Unmarshal([]byte, &v)         — data already in memory
json.NewDecoder(reader).Decode(&v) — data in a stream (io.Reader)
```

| | `Unmarshal` | `Decoder` |
|---|---|---|
| Input | `[]byte` | `io.Reader` |
| Memory | Must load all bytes first | Streams — lower peak memory |
| Multiple values | No | Yes (`dec.More()` loop) |
| Options | None | `UseNumber`, `DisallowUnknownFields` |
| Best for | Small, in-memory JSON | HTTP bodies, files, stdin |

```
Unmarshal flow:
  io.ReadAll(resp.Body) → []byte → json.Unmarshal → struct
  ⚠ 2 copies in memory: raw bytes + parsed struct

Decoder flow:
  resp.Body → json.NewDecoder → .Decode → struct
  ✓ 1 copy: just the parsed struct
```

Same for output:

| | `Marshal` | `Encoder` |
|---|---|---|
| Output | `[]byte` | `io.Writer` |
| Best for | Build JSON in memory, then use it | Write directly to HTTP/file/stdout |

---

# Part 3 — HTTP Use Cases

## 11. GET — Fetch JSON from an API

### Minimal Example

```go
func fetchUser(id int) (*User, error) {
    url := fmt.Sprintf("https://api.example.com/users/%d", id)

    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("fetch user: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status: %s", resp.Status)
    }

    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, fmt.Errorf("decode user: %w", err)
    }
    return &user, nil
}
```

### With Custom Headers / Timeout

```go
func fetchWithAuth(url, token string) (*Data, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Accept", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("status %d", resp.StatusCode)
    }

    var data Data
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, fmt.Errorf("decode: %w", err)
    }
    return &data, nil
}
```

---

## 12. POST — Send JSON to an API

### Using json.NewEncoder with a Pipe

```go
func createUser(user User) (*User, error) {
    body, err := json.Marshal(user)
    if err != nil {
        return nil, err
    }

    resp, err := http.Post(
        "https://api.example.com/users",
        "application/json",
        bytes.NewReader(body),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var created User
    json.NewDecoder(resp.Body).Decode(&created)
    return &created, nil
}
```

### Using bytes.Buffer (Alternative)

```go
func createUser(user User) (*User, error) {
    var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(user); err != nil {
        return nil, err
    }

    resp, err := http.Post(
        "https://api.example.com/users",
        "application/json",
        &buf,
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var created User
    json.NewDecoder(resp.Body).Decode(&created)
    return &created, nil
}
```

### PUT / PATCH / DELETE with Custom Request

```go
func updateUser(id int, updates map[string]any) error {
    body, _ := json.Marshal(updates)

    req, err := http.NewRequest("PATCH",
        fmt.Sprintf("https://api.example.com/users/%d", id),
        bytes.NewReader(body),
    )
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("update failed: %s", resp.Status)
    }
    return nil
}
```

---

## 13. Building a JSON REST Server

### Read JSON Request → Write JSON Response

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Decode request body
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
        return
    }

    // 2. Validate
    if req.Name == "" || req.Email == "" {
        http.Error(w, `{"error":"name and email required"}`, http.StatusBadRequest)
        return
    }

    // 3. Do business logic (e.g., save to DB)
    id := 42

    // 4. Send JSON response
    resp := CreateUserResponse{ID: id, Name: req.Name, Email: req.Email}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(resp)
}
```

### Helper Functions (DRY)

```go
func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func readJSON(r *http.Request, v any) error {
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    return dec.Decode(v)
}

// Limit body size to prevent abuse
func readJSONSafe(r *http.Request, v any, maxBytes int64) error {
    r.Body = http.MaxBytesReader(nil, r.Body, maxBytes)
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    return dec.Decode(v)
}
```

### Usage in Handlers

```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    user := User{Name: "Alice", Email: "alice@example.com"}
    writeJSON(w, http.StatusOK, user)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := readJSON(r, &req); err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }
    // ...
}
```

---

## 14. Error Responses

### Consistent Error Format

```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, APIError{Code: status, Message: msg})
}

// Usage
writeError(w, 404, "user not found")
// {"code":404,"message":"user not found"}
```

### Parse API Error from Response

```go
func checkResponse(resp *http.Response) error {
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        return nil
    }

    var apiErr APIError
    if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
        return fmt.Errorf("HTTP %d", resp.StatusCode)
    }
    return fmt.Errorf("API error %d: %s", apiErr.Code, apiErr.Message)
}
```

---

## 15. Real-World Example — GitHub API Client

The pattern used in your `chapter4/github` code:

```go
// 1. Define types for the JSON you care about (ignore the rest)
type IssueSearchResult struct {
    TotalCount int      `json:"total_count"`
    Items      []*Issue
}

type Issue struct {
    Number    int
    HTMLURL   string    `json:"html_url"`
    Title     string
    State     string
    User      *User
    CreatedAt time.Time `json:"created_at"`
    Body      string
}

type User struct {
    Login   string
    HTMLURL string `json:"html_url"`
}

// 2. Make HTTP request
func SearchIssues(terms []string) (*IssueSearchResult, error) {
    q := url.QueryEscape(strings.Join(terms, " "))
    resp, err := http.Get(issueURL + "?q=" + q)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("search failed: %s", resp.Status)
    }

    // 3. Decode directly from response body stream
    var result IssueSearchResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    return &result, nil
}
```

Key takeaways from this pattern:
- **Define only the fields you need** — GitHub returns 50+ fields, your struct has 7
- **Use `json.NewDecoder`** — stream from `resp.Body` without buffering
- **Use struct tags** for snake_case → CamelCase mapping (`json:"html_url"`)
- **Nested structs** decode automatically (`User` inside `Issue`)
- **Always check status code** before attempting to decode

---

# Part 4 — Advanced

## 16. Custom Marshal / Unmarshal

Implement `json.Marshaler` or `json.Unmarshaler` for full control.

### Custom Marshaler

```go
type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {
    s := fmt.Sprintf(`"%d"`, time.Time(t).Unix())
    return []byte(s), nil
}

type Event struct {
    Name string    `json:"name"`
    At   Timestamp `json:"at"`
}

// Marshals "at" as Unix timestamp string: {"name":"deploy","at":"1711324800"}
```

### Custom Unmarshaler

```go
func (t *Timestamp) UnmarshalJSON(data []byte) error {
    var unix int64
    // data might be `"1711324800"` (string) or `1711324800` (number)
    if err := json.Unmarshal(data, &unix); err != nil {
        var s string
        if err := json.Unmarshal(data, &s); err != nil {
            return err
        }
        unix, _ = strconv.ParseInt(s, 10, 64)
    }
    *t = Timestamp(time.Unix(unix, 0))
    return nil
}
```

### Enum with Custom JSON

```go
type Status int

const (
    StatusActive  Status = iota
    StatusPending
    StatusClosed
)

var statusNames = map[Status]string{
    StatusActive:  "active",
    StatusPending: "pending",
    StatusClosed:  "closed",
}

var statusValues = map[string]Status{
    "active":  StatusActive,
    "pending": StatusPending,
    "closed":  StatusClosed,
}

func (s Status) MarshalJSON() ([]byte, error) {
    return json.Marshal(statusNames[s])
}

func (s *Status) UnmarshalJSON(data []byte) error {
    var name string
    if err := json.Unmarshal(data, &name); err != nil {
        return err
    }
    val, ok := statusValues[name]
    if !ok {
        return fmt.Errorf("unknown status: %q", name)
    }
    *s = val
    return nil
}
```

---

## 17. Nested & Embedded Structs

### Nested (Separate JSON Object)

```go
type Address struct {
    City    string `json:"city"`
    Country string `json:"country"`
}

type User struct {
    Name    string  `json:"name"`
    Address Address `json:"address"`
}

// JSON: {"name":"Alice","address":{"city":"Hanoi","country":"Vietnam"}}
```

### Embedded (Flattened into Parent)

```go
type User struct {
    Name string `json:"name"`
    Address                        // no json tag, no field name
}

// JSON: {"name":"Alice","city":"Hanoi","country":"Vietnam"}
// Address fields are promoted — flattened into User's JSON
```

### Pointer Fields (null vs absent)

```go
type Update struct {
    Name  *string `json:"name,omitempty"`   // nil → omitted, "" → "name":""
    Score *int    `json:"score,omitempty"`   // nil → omitted, 0 → "score":0
}

// Distinguish "field not provided" (nil) from "field set to zero" (pointer to zero)
name := "Alice"
update := Update{Name: &name, Score: nil}
// {"name":"Alice"}  — Score omitted (nil pointer)
```

---

## 18. Common Gotchas

### 1. Numbers Are Always `float64` in `any`

```go
var m map[string]any
json.Unmarshal([]byte(`{"age":30}`), &m)

m["age"]          // float64(30), NOT int
int(m["age"].(float64))  // convert manually

// Fix: use UseNumber() or decode into a struct with typed fields
```

### 2. `nil` Slice vs Empty Slice in JSON

```go
type R struct {
    Tags []string `json:"tags"`
}

json.Marshal(R{Tags: nil})       // {"tags":null}
json.Marshal(R{Tags: []string{}}) // {"tags":[]}
```

APIs usually expect `[]`, not `null`. Initialize slices if needed.

### 3. `omitempty` Doesn't Omit Empty Structs

```go
type Inner struct{ X int }
type Outer struct {
    Inner Inner `json:"inner,omitempty"`
}

json.Marshal(Outer{})
// {"inner":{"X":0}}  ← NOT omitted! Structs are never "empty" for omitempty

// Fix: use a pointer
type Outer2 struct {
    Inner *Inner `json:"inner,omitempty"`
}
json.Marshal(Outer2{})
// {}  ← nil pointer IS omitted
```

### 4. `Encode` Adds a Trailing Newline, `Marshal` Doesn't

```go
json.NewEncoder(w).Encode(v)   // writes: {"name":"Alice"}\n
json.Marshal(v)                 // returns: {"name":"Alice"}
```

### 5. Unexported Fields Are Invisible

```go
type Config struct {
    host string   // lowercase → unexported → IGNORED by encoding/json
    Port int      // uppercase → exported → included
}

json.Marshal(Config{host: "localhost", Port: 8080})
// {"Port":8080}  ← host is missing
```

### 6. `time.Time` Uses RFC 3339

```go
t := time.Now()
json.Marshal(t)
// "2026-03-24T15:04:05.999999999+07:00"

// If your API uses a different format, implement custom MarshalJSON/UnmarshalJSON
```

---

## Quick Reference Card

```
MARSHAL       json.Marshal(v)          → []byte, error
              json.MarshalIndent(v,p,i) → []byte, error (pretty)
UNMARSHAL     json.Unmarshal([]byte,&v) → error
ENCODE        json.NewEncoder(w).Encode(v)   — write to io.Writer
DECODE        json.NewDecoder(r).Decode(&v)  — read from io.Reader

TAGS          `json:"name"`             custom key name
              `json:"-"`                skip always
              `json:",omitempty"`       skip if zero value
              `json:",string"`          number/bool as JSON string

DECODER OPTS  dec.UseNumber()           keep numbers as json.Number
              dec.DisallowUnknownFields()  reject extra keys

TYPES         json.RawMessage           defer decoding (raw bytes)
              json.Number               string-form number (precise)
              json.Marshaler            interface — custom marshal
              json.Unmarshaler          interface — custom unmarshal

HTTP PATTERN  GET:  http.Get → json.NewDecoder(resp.Body).Decode(&v)
              POST: json.Marshal → bytes.NewReader → http.Post
              SERVER: json.NewDecoder(r.Body).Decode(&req)
                      json.NewEncoder(w).Encode(resp)
```

## Decision Flowchart

```
Do I have []byte or io.Reader?
├─ []byte in memory → json.Unmarshal / json.Marshal
└─ io.Reader (HTTP body, file, stdin) → json.NewDecoder / json.NewEncoder

Do I know the structure?
├─ YES, I have a struct → decode into struct (best: type-safe, fast)
├─ PARTIALLY (wrapper known, payload varies) → json.RawMessage
└─ NO, fully dynamic → map[string]any (need type assertions)

Do I need precise integers?
├─ YES → decoder.UseNumber() → json.Number → .Int64()
└─ NO → default float64 is fine
```
