# Go Structs — Complete Cheatsheet

> From declaration to embedding, methods, tags, and real-world patterns.

---

## Table of Contents

1. [What Is a Struct?](#1-what-is-a-struct)
2. [Declaring Structs](#2-declaring-structs)
3. [Creating Struct Values](#3-creating-struct-values)
4. [Accessing & Modifying Fields](#4-accessing--modifying-fields)
5. [Pointers to Structs](#5-pointers-to-structs)
6. [Methods — Value vs Pointer Receiver](#6-methods--value-vs-pointer-receiver)
7. [Constructor Pattern](#7-constructor-pattern)
8. [Struct Embedding (Composition)](#8-struct-embedding-composition)
9. [Anonymous Structs](#9-anonymous-structs)
10. [Struct Comparison & Equality](#10-struct-comparison--equality)
11. [Struct Tags](#11-struct-tags)
12. [Structs and JSON](#12-structs-and-json)
13. [Structs as Map Keys](#13-structs-as-map-keys)
14. [Structs and Interfaces](#14-structs-and-interfaces)
15. [Common Patterns & Recipes](#15-common-patterns--recipes)

---

## 1. What Is a Struct?

A struct is a **composite type** that groups together zero or more named fields of different types. It's Go's primary tool for building custom data types (no classes, no inheritance).

```go
type User struct {
    Name  string
    Email string
    Age   int
}
```

Key properties:

- **Value type** — assigning a struct copies all fields
- **Fields are stored contiguously** in memory
- **Zero value is usable** — all fields are set to their zero values
- **No constructors, no inheritance** — composition via embedding instead

---

## 2. Declaring Structs

### Named Struct Type

```go
type Point struct {
    X, Y float64                  // multiple fields of same type
}

type Person struct {
    FirstName string
    LastName  string
    Age       int
    Active    bool
}
```

### Empty Struct — Zero Bytes

```go
type Signal struct{}              // size = 0 bytes

// Used as:
done := make(chan struct{})       // signal channel (no data, just events)
seen := map[string]struct{}{}    // set (more memory-efficient than map[string]bool)
seen["alice"] = struct{}{}
```

### Exported vs Unexported Fields

```go
type Config struct {
    Host     string     // Exported   — visible outside the package (uppercase)
    Port     int        // Exported
    timeout  int        // unexported — only visible within this package (lowercase)
    retries  int        // unexported
}
```

The same rule applies to the struct type name itself: `type Config` is exported, `type config` is not.

---

## 3. Creating Struct Values

### Zero Value (All Fields Default)

```go
var u User                        // {"" "" 0}
fmt.Println(u.Name)              // ""
fmt.Println(u.Age)               // 0
```

### Named Fields (Recommended)

```go
u := User{
    Name:  "Alice",
    Email: "alice@example.com",
    Age:   30,
}

// Order doesn't matter, omitted fields get zero value
u2 := User{
    Name: "Bob",
    // Email defaults to "", Age defaults to 0
}
```

### Positional Fields (Fragile — Avoid in Public APIs)

```go
p := Point{3.0, 4.0}            // must provide ALL fields in declaration order
```

This breaks if you add/reorder fields later. Prefer named fields.

### Pointer to Struct (Heap Allocation)

```go
u:= &User{
    Name:  "Carol",
    Email: "carol@example.com",
    Age:   25,
} 
// u is *User
```

### Using `new` (Rare — Gives Zero Value)

```go
u := new(User)                   // *User, all fields zeroed
u.Name = "Dave"
```

---

## 4. Accessing & Modifying Fields

### Dot Notation

```go
u := User{Name: "Alice", Age: 30}

fmt.Println(u.Name)              // "Alice"
u.Age = 31                       // modify directly
```

### Structs Are Value Types — Assignment Copies

```go
a := Point{1, 2}
b := a                            // COPY — b is independent

b.X = 99
fmt.Println(a.X)                 // 1  — a is unchanged
fmt.Println(b.X)                 // 99
```

This is different from slices/maps which share underlying data. If you need shared mutation, use pointers.

---

## 5. Pointers to Structs

### Why Use Pointers?

1. **Avoid copying** large structs
2. **Allow mutation** — modifying through a pointer changes the original
3. **Represent absence** — `nil` pointer means "no value"

### Basic Usage

```go
u := User{Name: "Alice", Age: 30}
p := &u                           // p is *User

p.Age = 31                        // Go auto-dereferences: same as (*p).Age = 31
fmt.Println(u.Age)               // 31 — original modified
```

### `nil` Pointer Check

```go
var p *User                       // nil
if p != nil {
    fmt.Println(p.Name)
}
// Accessing fields on nil pointer → panic
```

### When to Use Pointer vs Value


| Use **pointer** `*T` when...         | Use **value** `T` when...       |
| ------------------------------------ | ------------------------------- |
| Struct is large (>~64 bytes)         | Struct is small (few fields)    |
| Need to mutate the original          | Want an independent copy        |
| Need to represent "no value" (`nil`) | Always valid (zero value OK)    |
| Sharing between goroutines           | Local, stack-only use           |
| Methods need pointer receivers       | All methods use value receivers |


---

## 6. Methods — Value vs Pointer Receiver

### Value Receiver — Works on a Copy

```go
type Point struct{ X, Y float64 }

func (p Point) Distance() float64 {
    return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p Point) String() string {
    return fmt.Sprintf("(%g, %g)", p.X, p.Y)
}

pt := Point{3, 4}
fmt.Println(pt.Distance())       // 5
```

### Pointer Receiver — Mutates the Original

```go
func (p *Point) Scale(factor float64) {
    p.X *= factor
    p.Y *= factor
}

func (p *Point) Translate(dx, dy float64) {
    p.X += dx
    p.Y += dy
}

pt := Point{3, 4}
pt.Scale(2)                       // Go auto-takes address of pt
fmt.Println(pt)                  // {6, 8} — modified!
```

### The Rules

```
Value receiver:   func (t T) Method()
  - Gets a COPY of the struct
  - Cannot modify original
  - Can be called on both value and pointer

Pointer receiver: func (t *T) Method()
  - Gets a POINTER to the struct
  - CAN modify original
  - Can be called on both value and pointer (Go auto-takes address)
```

### Consistency Rule

If **any** method needs a pointer receiver, make **all** methods use pointer receivers. Mixing is allowed but discouraged.

```go
// GOOD: consistent pointer receivers
func (u *User) SetName(name string) { u.Name = name }
func (u *User) FullName() string    { return u.FirstName + " " + u.LastName }

// BAD: mixing receivers (confusing)
func (u User) FullName() string     { return u.FirstName + " " + u.LastName }
func (u *User) SetName(name string) { u.Name = name }
```

### Method Set (Matters for Interfaces)


| Receiver type | Method set                              |
| ------------- | --------------------------------------- |
| Value `T`     | Only value receiver methods             |
| Pointer `*T`  | Both value AND pointer receiver methods |


```go
type Stringer interface { String() string }
type Scaler interface { Scale(float64) }

var _ Stringer = Point{}          // ✓ String() has value receiver
var _ Scaler = Point{}            // ✗ Scale() has pointer receiver
var _ Scaler = &Point{}           // ✓ pointer has both method sets
```

---

## 7. Constructor Pattern

Go has no constructors. Use a `New*` function instead.

### Basic Constructor

```go
func NewUser(name, email string, age int) User {
    return User{
        Name:  name,
        Email: email,
        Age:   age,
    }
}
```

### Return Pointer (When Struct Is Large or Has Unexported Fields)

```go
func NewServer(addr string) *Server {
    return &Server{
        addr:    addr,
        timeout: 30 * time.Second,     // sensible default
        logger:  log.Default(),
    }
}
```

### Constructor with Validation

```go
func NewUser(name string, age int) (*User, error) {
    if name == "" {
        return nil, errors.New("name cannot be empty")
    }
    if age < 0 || age > 150 {
        return nil, fmt.Errorf("invalid age: %d", age)
    }
    return &User{Name: name, Age: age}, nil
}
```

### Functional Options Pattern (Flexible Config)

```go
type Server struct {
    addr    string
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

func WithLogger(l *log.Logger) Option {
    return func(s *Server) { s.logger = l }
}

func NewServer(addr string, opts ...Option) *Server {
    s := &Server{
        addr:    addr,
        port:    8080,                  // default
        timeout: 30 * time.Second,      // default
        logger:  log.Default(),         // default
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage — only specify what you want to override
srv := NewServer("localhost",
    WithPort(9090),
    WithTimeout(60*time.Second),
)
```

---

## 8. Struct Embedding (Composition)

Go has no inheritance. Instead, you **embed** one struct inside another to reuse fields and methods.

### Basic Embedding

```go
type Animal struct {
    Name string
}

func (a Animal) Speak() string {
    return a.Name + " makes a sound"
}

type Dog struct {
    Animal                          // embedded — no field name
    Breed string
}

d := Dog{
    Animal: Animal{Name: "Rex"},
    Breed:  "Labrador",
}

fmt.Println(d.Name)               // "Rex"       — promoted field
fmt.Println(d.Speak())            // "Rex makes a sound" — promoted method
fmt.Println(d.Animal.Name)        // "Rex"       — explicit access also works
```

### Embedding Promotes Fields and Methods

The embedded type's fields and methods are **promoted** to the outer type — you can access them directly without naming the embedded type.

### Override (Shadowing) a Promoted Method

```go
func (d Dog) Speak() string {
    return d.Name + " barks"
}

d := Dog{Animal: Animal{Name: "Rex"}}
fmt.Println(d.Speak())            // "Rex barks"       — Dog's method
fmt.Println(d.Animal.Speak())     // "Rex makes a sound" — Animal's method still accessible
```

### Multiple Embedding

```go
type Logger struct{}
func (l Logger) Log(msg string) { fmt.Println("LOG:", msg) }

type Metrics struct{}
func (m Metrics) Record(name string, val float64) { /* ... */ }

type Service struct {
    Logger                          // gets Log()
    Metrics                         // gets Record()
    Name string
}

svc := Service{Name: "UserService"}
svc.Log("started")
svc.Record("latency", 42.5)
```

### Embedding Pointers

```go
type Cache struct {
    *sync.Mutex                     // embed a pointer — gets Lock()/Unlock()
    data map[string]string
}

c := &Cache{
    Mutex: &sync.Mutex{},
    data:  make(map[string]string),
}
c.Lock()
c.data["key"] = "value"
c.Unlock()
```

### Ambiguity — When Two Embeddings Collide

```go
type A struct{}
func (A) Hello() string { return "A" }

type B struct{}
func (B) Hello() string { return "B" }

type C struct {
    A
    B
}

c := C{}
// c.Hello()            // ✗ compile error: ambiguous
fmt.Println(c.A.Hello())  // ✓ "A" — disambiguate explicitly
fmt.Println(c.B.Hello())  // ✓ "B"
```

---

## 9. Anonymous Structs

### Inline Struct (No Type Name)

```go
point := struct {
    X, Y int
}{10, 20}

fmt.Println(point.X)              // 10
```

### Useful for One-Off Data

```go
// Test table
tests := []struct {
    input    string
    expected int
}{
    {"hello", 5},
    {"", 0},
    {"日本語", 3},
}

for _, tt := range tests {
    if got := len([]rune(tt.input)); got != tt.expected {
        fmt.Printf("len(%q) = %d, want %d\n", tt.input, got, tt.expected)
    }
}
```

### Grouping Related Variables

```go
var config struct {
    Debug   bool
    Verbose bool
    MaxConn int
}
config.Debug = true
config.MaxConn = 100
```

---

## 10. Struct Comparison & Equality

### Comparable Structs (All Fields Comparable)

```go
type Point struct{ X, Y int }

a := Point{1, 2}
b := Point{1, 2}
c := Point{3, 4}

fmt.Println(a == b)               // true
fmt.Println(a == c)               // false
fmt.Println(a != c)               // true
```

### Non-Comparable Structs

If a struct contains **slices, maps, or functions**, it is NOT comparable with `==`:

```go
type Data struct {
    Values []int                   // slice → not comparable
}

a := Data{Values: []int{1, 2}}
b := Data{Values: []int{1, 2}}
// a == b                         // ✗ compile error!
```

Use `reflect.DeepEqual` for deep comparison (slow, avoid in hot paths):

```go
import "reflect"
reflect.DeepEqual(a, b)           // true
```

Or write a custom `Equal` method:

```go
func (d Data) Equal(other Data) bool {
    return slices.Equal(d.Values, other.Values)
}
```

---

## 11. Struct Tags

Struct tags are **string metadata** attached to fields. They don't affect the struct's behavior directly — they're read at runtime by packages like `encoding/json`, `encoding/xml`, database ORMs, validators, etc.

### Syntax

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age,omitempty" db:"user_age"`
}
```

### Reading Tags at Runtime

```go
import "reflect"

t := reflect.TypeOf(User{})
field, _ := t.FieldByName("Email")
fmt.Println(field.Tag.Get("json"))      // "email"
fmt.Println(field.Tag.Get("validate"))  // "required,email"
```

### Common Tag Formats


| Package         | Tag        | Example                         |
| --------------- | ---------- | ------------------------------- |
| `encoding/json` | `json`     | `json:"name,omitempty"`         |
| `encoding/xml`  | `xml`      | `xml:"name,attr"`               |
| `database/sql`  | `db`       | `db:"user_name"`                |
| `gorm`          | `gorm`     | `gorm:"column:name;primaryKey"` |
| `yaml`          | `yaml`     | `yaml:"name"`                   |
| `validate`      | `validate` | `validate:"required,min=1"`     |


---

## 12. Structs and JSON

### Marshal (Struct → JSON)

```go
type Product struct {
    Name     string  `json:"name"`
    Price    float64 `json:"price"`
    InStock  bool    `json:"in_stock"`
    Internal string  `json:"-"`              // always excluded from JSON
    Notes    string  `json:"notes,omitempty"` // excluded if zero value
}

p := Product{Name: "Widget", Price: 9.99, InStock: true, Internal: "secret"}
data, _ := json.Marshal(p)
fmt.Println(string(data))
// {"name":"Widget","price":9.99,"in_stock":true}

// Pretty print
data, _ = json.MarshalIndent(p, "", "  ")
```

### Unmarshal (JSON → Struct)

```go
jsonStr := `{"name":"Gadget","price":19.99,"in_stock":false}`
var p Product
err := json.Unmarshal([]byte(jsonStr), &p)

fmt.Println(p.Name)              // "Gadget"
fmt.Println(p.Price)             // 19.99
```

### JSON Tag Options


| Tag                     | Meaning                                                 |
| ----------------------- | ------------------------------------------------------- |
| `json:"name"`           | Use "name" as the JSON key                              |
| `json:"-"`              | Always skip this field                                  |
| `json:"name,omitempty"` | Skip if zero value (0, "", false, nil, empty slice/map) |
| `json:",omitempty"`     | Use field name as key, but skip if zero                 |
| `json:",string"`        | Encode number/bool as JSON string                       |


### Nested Structs

```go
type Address struct {
    City    string `json:"city"`
    Country string `json:"country"`
}

type User struct {
    Name    string  `json:"name"`
    Address Address `json:"address"`
}

// Produces: {"name":"Alice","address":{"city":"Hanoi","country":"Vietnam"}}

// Flatten with embedding:
type User2 struct {
    Name string `json:"name"`
    Address                        // embedded — fields promoted into User2's JSON
}
// Produces: {"name":"Alice","city":"Hanoi","country":"Vietnam"}
```

---

## 13. Structs as Map Keys

Comparable structs can be used as map keys — great for multi-dimensional lookups:

```go
type Point struct{ X, Y int }

visited := map[Point]bool{}
visited[Point{1, 2}] = true
visited[Point{3, 4}] = true

if visited[Point{1, 2}] {
    fmt.Println("already visited")
}
```

```go
type Edge struct{ From, To string }

weights := map[Edge]float64{
    {"A", "B"}: 1.5,
    {"B", "C"}: 2.0,
}
```

---

## 14. Structs and Interfaces

Go interfaces are satisfied **implicitly** — no `implements` keyword.

### Implement an Interface

```go
type Shape interface {
    Area() float64
    Perimeter() float64
}

type Rectangle struct{ W, H float64 }

func (r Rectangle) Area() float64      { return r.W * r.H }
func (r Rectangle) Perimeter() float64 { return 2 * (r.W + r.H) }

type Circle struct{ Radius float64 }

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

// Both satisfy Shape
shapes := []Shape{Rectangle{3, 4}, Circle{5}}
for _, s := range shapes {
    fmt.Printf("Area=%.2f  Perimeter=%.2f\n", s.Area(), s.Perimeter())
}
```

### Compile-Time Interface Check

```go
var _ Shape = Rectangle{}         // ✓ compile error if Rectangle doesn't satisfy Shape
var _ Shape = (*Circle)(nil)      // ✓ for pointer receiver methods
```

### Stringer Interface (Like `toString()`)

```go
type User struct{ Name string; Age int }

func (u User) String() string {
    return fmt.Sprintf("%s (age %d)", u.Name, u.Age)
}

u := User{"Alice", 30}
fmt.Println(u)                    // "Alice (age 30)" — fmt calls String() automatically
```

### Error Interface

```go
type NotFoundError struct {
    Resource string
    ID       string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s %q not found", e.Resource, e.ID)
}

func FindUser(id string) (*User, error) {
    return nil, &NotFoundError{Resource: "user", ID: id}
}
```

---

## 15. Common Patterns & Recipes

### Builder Pattern

```go
type QueryBuilder struct {
    table  string
    wheres []string
    limit  int
}

func NewQuery(table string) *QueryBuilder {
    return &QueryBuilder{table: table, limit: -1}
}

func (q *QueryBuilder) Where(cond string) *QueryBuilder {
    q.wheres = append(q.wheres, cond)
    return q                       // return self for chaining
}

func (q *QueryBuilder) Limit(n int) *QueryBuilder {
    q.limit = n
    return q
}

func (q *QueryBuilder) Build() string {
    query := "SELECT * FROM " + q.table
    for i, w := range q.wheres {
        if i == 0 {
            query += " WHERE " + w
        } else {
            query += " AND " + w
        }
    }
    if q.limit > 0 {
        query += fmt.Sprintf(" LIMIT %d", q.limit)
    }
    return query
}

// Usage
sql := NewQuery("users").Where("age > 18").Where("active = true").Limit(10).Build()
// "SELECT * FROM users WHERE age > 18 AND active = true LIMIT 10"
```

### Enum with Struct Methods

```go
type Color int

const (
    Red Color = iota
    Green
    Blue
)

func (c Color) String() string {
    return [...]string{"Red", "Green", "Blue"}[c]
}

fmt.Println(Red)                  // "Red"
```

### Immutable Struct (Unexported Fields + Getters)

```go
type Money struct {
    amount   int64                 // unexported — can't be modified from outside
    currency string
}

func NewMoney(amount int64, currency string) Money {
    return Money{amount: amount, currency: currency}
}

func (m Money) Amount() int64    { return m.amount }
func (m Money) Currency() string { return m.currency }
func (m Money) Add(other Money) Money {
    if m.currency != other.currency {
        panic("currency mismatch")
    }
    return Money{amount: m.amount + other.amount, currency: m.currency}
}
```

### Copy with Modification

```go
type Config struct {
    Host    string
    Port    int
    Timeout time.Duration
}

original := Config{Host: "localhost", Port: 8080, Timeout: 30 * time.Second}

modified := original              // copy all fields
modified.Port = 9090              // then override what you need
```

### Slice of Structs — Sort

```go
type Employee struct {
    Name   string
    Salary float64
}

employees := []Employee{
    {"Carol", 95000},
    {"Alice", 75000},
    {"Bob", 85000},
}

slices.SortFunc(employees, func(a, b Employee) int {
    return cmp.Compare(a.Salary, b.Salary)
})
// [{Alice 75000} {Bob 85000} {Carol 95000}]
```

### Struct ↔ Map Conversion

```go
func structToMap(u User) map[string]any {
    return map[string]any{
        "name":  u.Name,
        "email": u.Email,
        "age":   u.Age,
    }
}

// For generic conversion, use reflect or encoding/json round-trip:
func toMap(v any) map[string]any {
    data, _ := json.Marshal(v)
    var m map[string]any
    json.Unmarshal(data, &m)
    return m
}
```

---

## Quick Reference Card

```
DECLARE       type T struct { Field Type }
CREATE        T{}  T{field: val}  &T{}  new(T)
ACCESS        t.Field  p.Field (auto-deref)
METHODS       func (t T) M()   — value receiver (copy)
              func (t *T) M()  — pointer receiver (mutate)
EMBEDDING     type Outer struct { Inner }  — promotes fields & methods
TAGS          `json:"name,omitempty"  db:"col"`
ZERO VALUE    all fields zeroed: 0, "", false, nil
COMPARABLE    == works if ALL fields are comparable (no slices/maps/funcs)
INTERFACE     satisfied implicitly by matching method set
CONSTRUCTOR   func NewT(...) *T { return &T{...} }
```

## Value vs Pointer — Decision Flowchart

```
Should I use *T or T?
│
├─ Does any method need to modify the struct?
│  └─ YES → use *T (pointer receiver + pointer everywhere)
│
├─ Is the struct large (many fields / large arrays)?
│  └─ YES → use *T (avoid copying)
│
├─ Do you need nil to mean "absent"?
│  └─ YES → use *T
│
├─ Is the struct used concurrently?
│  └─ YES → use *T (shared state needs pointer)
│
└─ Otherwise → T is fine (simple, safe, stack-allocated)
```

## Memory Layout

```go
type Example struct {
    A bool      // 1 byte
    // 7 bytes padding (alignment for B)
    B int64     // 8 bytes
    C bool      // 1 byte
    // 3 bytes padding (alignment for D)
    D int32     // 4 bytes
}
// Total: 24 bytes (with padding)

// Reorder for less padding:
type ExampleOptimized struct {
    B int64     // 8 bytes
    D int32     // 4 bytes
    A bool      // 1 byte
    C bool      // 1 byte
    // 2 bytes padding
}
// Total: 16 bytes — saved 8 bytes!
```

Use `unsafe.Sizeof(Example{})` to check. Field ordering matters for memory-sensitive applications.