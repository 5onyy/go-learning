# Go Methods — Complete Cheatsheet

> From method declarations to embedding, method values, and the `Stringer` interface — everything about methods on Go types.

---

## Table of Contents

1. [What Is a Method?](#1-what-is-a-method)
2. [Method Declarations](#2-method-declarations)
3. [Value Receiver vs Pointer Receiver](#3-value-receiver-vs-pointer-receiver)
4. [Methods on Named Types](#4-methods-on-named-types)
5. [Implicit Conversions (Addressability)](#5-implicit-conversions-addressability)
6. [Nil Receivers](#6-nil-receivers)
7. [Struct Embedding & Method Promotion](#7-struct-embedding--method-promotion)
8. [Methods on Anonymous Structs](#8-methods-on-anonymous-structs)
9. [Method Values & Method Expressions](#9-method-values--method-expressions)
10. [Encapsulation](#10-encapsulation)
11. [The `String()` Method (Stringer Interface)](#11-the-string-method-stringer-interface)
12. [Common Patterns & Recipes](#12-common-patterns--recipes)

---

## 1. What Is a Method?

A method is a **function with a receiver** — a special parameter that binds the function to a type. It lets you call `value.Method()` instead of `Function(value)`.

```go
type Point struct{ X, Y float64 }

// Function — standalone
func Distance(p, q Point) float64 {
    return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// Method — attached to Point
func (p Point) Distance(q Point) float64 {
    return math.Hypot(q.X-p.X, q.Y-p.Y)
}
```

```go
p := Point{1, 2}
q := Point{4, 6}

Distance(p, q)     // function call
p.Distance(q)      // method call — p is the receiver
```

The parameter before the function name is called the **receiver**. Go convention: use a **short (1–2 letter) abbreviation** of the type name, not `this` or `self`.

---

## 2. Method Declarations

### Syntax

```go
func (receiver ReceiverType) MethodName(params) returnType {
    // body
}
```

### Rules

| Rule | Detail |
|------|--------|
| Named types only | Methods can only be declared on named types defined in the **same package** |
| No primitives | Cannot attach methods to `int`, `string`, etc. directly |
| Same package | The method and the type must be in the same package |
| Receiver name | Convention: 1–2 letter abbreviation of the type (e.g., `p` for `Point`) |
| Consistent receiver name | Use the same receiver name across all methods of a type |

### Methods on Named Slice Types

You cannot attach methods to built-in types, but you can create a named type and attach methods to it:

```go
type Path []Point

func (path Path) Distance() float64 {
    sum := 0.0
    for i := range path {
        if i > 0 {
            sum += path[i-1].Distance(path[i])
        }
    }
    return sum
}
```

```go
perim := Path{{1, 1}, {5, 1}, {5, 4}, {1, 1}}
fmt.Println(perim.Distance())    // 12
```

---

## 3. Value Receiver vs Pointer Receiver

### Value Receiver — Operates on a Copy

```go
func (p Point) Distance(q Point) float64 {
    return math.Hypot(q.X-p.X, q.Y-p.Y)
}
```

- Receiver is **copied** — mutations inside the method don't affect the caller
- Safe for concurrent use (no shared state)

### Pointer Receiver — Operates on the Original

```go
func (p *Point) ScaleBy(factor float64) {
    p.X *= factor
    p.Y *= factor
}
```

- Receiver is a **pointer** — mutations are visible to the caller
- Avoids copying large structs
- Required when the method must modify the receiver

### When to Use Which?

```
Need to modify the receiver?
├── Yes → Pointer receiver
└── No
    ├── Large struct (avoid copy overhead)? → Pointer receiver
    ├── Consistency (other methods use pointer)? → Pointer receiver
    └── Small, immutable value? → Value receiver
```

### The Consistency Rule

> If **any** method of a type has a pointer receiver, then **all** methods of that type should have a pointer receiver — even ones that don't strictly need it.

This avoids confusion about whether a value or pointer is needed and prevents subtle bugs.

### Summary Table

| Aspect | Value receiver `(p Point)` | Pointer receiver `(p *Point)` |
|--------|---------------------------|------------------------------|
| Gets a | Copy of the value | Pointer to the original |
| Can modify receiver? | No (only the copy) | Yes |
| Can call on value? | Yes | Yes (implicit `&` taken) |
| Can call on pointer? | Yes (implicit `*` deref) | Yes |
| Nil receiver possible? | No | Yes |
| Best for | Small immutable types | Mutable types, large structs |

---

## 4. Methods on Named Types

You can define methods on **any named type** in your package — not just structs.

### Named Numeric Type

```go
type Celsius float64

func (c Celsius) String() string {
    return fmt.Sprintf("%.1f°C", c)
}
```

### Named Slice Type

```go
type Path []Point

func (path Path) Distance() float64 { ... }
```

### Named Map Type

```go
// From the standard library: net/url
type Values map[string][]string

func (v Values) Get(key string) string {
    if vs := v[key]; len(vs) > 0 {
        return vs[0]
    }
    return ""
}

func (v Values) Add(key, value string) {
    v[key] = append(v[key], value)
}
```

```go
m := url.Values{"lang": {"en"}}
m.Add("item", "1")
m.Add("item", "2")

fmt.Println(m.Get("lang"))    // en
fmt.Println(m.Get("q"))       // "" (zero value)
fmt.Println(m.Get("item"))    // 1 (first value)
fmt.Println(m["item"])         // [1 2] (direct map access)
```

### Restriction: No Methods on Unnamed Types or Types from Other Packages

```go
// ✗ Cannot define methods on built-in types
func (s string) Reverse() string { ... }     // compile error

// ✗ Cannot define methods on types from other packages
func (t time.Time) Hello() { ... }           // compile error

// ✓ Create a named type alias and add methods to that
type MyString string
func (s MyString) Reverse() string { ... }   // OK
```

---

## 5. Implicit Conversions (Addressability)

Go will implicitly take the address of a value or dereference a pointer so you can call methods conveniently.

### Calling a Pointer-Receiver Method on a Value

If the variable is **addressable**, the compiler inserts `&` automatically:

```go
p := Point{1, 2}
p.ScaleBy(2)             // compiler rewrites to: (&p).ScaleBy(2)
fmt.Println(p)           // {2 4} — p was modified
```

### Calling a Value-Receiver Method on a Pointer

The compiler automatically dereferences:

```go
pptr := &Point{1, 2}
d := pptr.Distance(q)    // compiler rewrites to: (*pptr).Distance(q)
```

### When Implicit `&` Does NOT Work

A value must be **addressable** — it must live in a variable. Temporary values are NOT addressable:

```go
Point{1, 2}.ScaleBy(2)   // ✗ compile error — cannot take address of Point literal
```

```go
// ✓ Assign to a variable first
p := Point{1, 2}
p.ScaleBy(2)

// ✓ Or use a pointer directly
(&Point{1, 2}).ScaleBy(2) // Technically works but rarely used
```

### Summary

| Call | Actual | Works? |
|------|--------|--------|
| `value.ValueMethod()` | Direct call | Yes |
| `value.PtrMethod()` | `(&value).PtrMethod()` | Yes, if addressable |
| `pointer.PtrMethod()` | Direct call | Yes |
| `pointer.ValueMethod()` | `(*pointer).ValueMethod()` | Yes |
| `literal.PtrMethod()` | — | No (not addressable) |

---

## 6. Nil Receivers

Methods with pointer receivers can be called on `nil` values. This enables nil-safe methods:

```go
type IntList struct {
    Value int
    Next  *IntList
}

func (l *IntList) Sum() int {
    if l == nil {
        return 0          // nil-safe: treat nil as empty list
    }
    return l.Value + l.Next.Sum()
}
```

```go
var list *IntList
fmt.Println(list.Sum())   // 0 — no panic
```

### Nil Map Receiver

The `url.Values` type (a named `map`) demonstrates nil receiver pitfalls:

```go
m := url.Values{"lang": {"en"}}
m = nil

fmt.Println(m.Get("item"))   // "" — reading nil map returns zero value (safe)
m.Add("item", "3")           // PANIC — writing to nil map
```

A nil map can be **read** but not **written**. Methods that write must guard against nil.

---

## 7. Struct Embedding & Method Promotion

### Embedding Promotes Methods

When you embed a type in a struct, its methods are **promoted** — you can call them as if they belong to the outer struct.

```go
type Point struct{ X, Y float64 }

func (p *Point) Distance(q Point) float64 {
    return math.Hypot(p.X-q.X, p.Y-q.Y)
}

func (p *Point) ScaleBy(factor float64) {
    p.X *= factor
    p.Y *= factor
}

type ColoredPoint struct {
    Point                        // embedded — NOT inheritance
    Color color.RGBA
}
```

```go
var cp ColoredPoint
cp.X = 1                        // access embedded field directly
cp.Point.Y = 2                  // or via the explicit field name
cp.ScaleBy(2)                   // call promoted method directly
```

### Embedding Is NOT Inheritance

`ColoredPoint` is **not** a `Point`. You cannot use a `ColoredPoint` where a `Point` is expected:

```go
p := ColoredPoint{Point{1, 4}, red}
q := ColoredPoint{Point{5, 4}, blue}

p.Distance(q.Point)             // ✓ must pass q.Point, not q
p.Distance(q)                   // ✗ compile error — q is not a Point
```

Think of it as: "a `ColoredPoint` **has a** `Point`" — not "a `ColoredPoint` **is a** `Point`."

### Multiple Embedding

A struct can embed multiple types. If methods conflict (same name), you must disambiguate:

```go
type A struct{}
func (A) Hello() string { return "A" }

type B struct{}
func (B) Hello() string { return "B" }

type C struct {
    A
    B
}

var c C
c.Hello()        // ✗ compile error — ambiguous
c.A.Hello()      // ✓ "A"
c.B.Hello()      // ✓ "B"
```

### Embedding Depth and Promotion Priority

Methods at a shallower depth take precedence:

```go
type Inner struct{}
func (Inner) Greet() string { return "inner" }

type Middle struct{ Inner }
func (Middle) Greet() string { return "middle" }

type Outer struct{ Middle }

var o Outer
o.Greet()        // "middle" — Middle.Greet shadows Inner.Greet
```

---

## 8. Methods on Anonymous Structs

Embedding enables methods on anonymous (unnamed) struct types. This is useful for one-off utility types:

```go
var cache = struct {
    sync.Mutex                   // embeds Lock() and Unlock()
    mapping map[string]string
}{
    mapping: make(map[string]string),
}

func lookUp(key string) string {
    cache.Lock()                 // promoted from sync.Mutex
    v := cache.mapping[key]
    cache.Unlock()
    return v
}
```

Use this pattern when you need a quick, one-off struct with mutex protection and don't want to define a named type.

---

## 9. Method Values & Method Expressions

### Method Values — Bind a Method to a Receiver

A method value is a function that **captures** a specific receiver. You can store it and call it later without the receiver:

```go
p := Point{1, 2}
q := Point{4, 6}

distanceFromP := p.Distance     // method value: binds p as receiver
fmt.Println(distanceFromP(q))   // 5 — equivalent to p.Distance(q)
```

Useful for passing methods as callbacks:

```go
type Rocket struct{ /* ... */ }
func (r *Rocket) Launch() { /* ... */ }

r := new(Rocket)
time.AfterFunc(10*time.Second, r.Launch)   // pass method value as callback
```

### Method Expressions — Unbound Method as Function

A method expression produces a function where the **first parameter is the receiver**:

```go
distance := Point.Distance      // method expression: func(Point, Point) float64
fmt.Println(distance(p, q))     // equivalent to p.Distance(q)

scaleBy := (*Point).ScaleBy     // method expression: func(*Point, float64)
scaleBy(&p, 2)                  // equivalent to p.ScaleBy(2)
```

Useful when you need to select a method dynamically:

```go
type Point struct{ X, Y float64 }
func (p Point) Add(q Point) Point { return Point{p.X + q.X, p.Y + q.Y} }
func (p Point) Sub(q Point) Point { return Point{p.X - q.X, p.Y - q.Y} }

func TranslatePath(path Path, offset Point, op func(Point, Point) Point) Path {
    result := make(Path, len(path))
    for i, p := range path {
        result[i] = op(p, offset)
    }
    return result
}

TranslatePath(path, Point{1, 2}, Point.Add)   // move right & up
TranslatePath(path, Point{1, 2}, Point.Sub)   // move left & down
```

### Comparison: Method Value vs Method Expression

| | Method Value | Method Expression |
|-|-------------|-------------------|
| Syntax | `receiver.Method` | `Type.Method` |
| Result type | `func(params) returns` | `func(Receiver, params) returns` |
| Receiver | Bound (captured at creation) | Unbound (passed as first arg) |
| Use case | Callbacks, event handlers | Dynamic method selection |

---

## 10. Encapsulation

Go controls visibility through **capitalization**, not access modifiers:

| Naming | Visibility | Meaning |
|--------|-----------|---------|
| `ExportedField` | Public | Accessible from any package |
| `unexportedField` | Package-private | Accessible only within the package |

### Encapsulation Through Unexported Fields

```go
// package bitset

type BitSet struct {
    n    int       // unexported — hidden from other packages
    data []uint64  // unexported — hidden from other packages
}

func (bs *BitSet) Len() int { return bs.n }   // exported getter
func (bs *BitSet) Add(x uint64) { ... }       // exported method
```

### Benefits of Encapsulation

1. **Callers can't modify internals directly** — must use methods
2. **Implementation can change** without breaking callers (e.g., change `data` from `[]uint64` to `[]uint32`)
3. **Invariants are maintained** — e.g., `n` always equals the actual count

### Getter Naming Convention

Go does **not** use `GetX()` — just name the getter after the field:

```go
// ✗ Java-style — not idiomatic Go
func (bs *BitSet) GetLen() int { return bs.n }

// ✓ Go-style
func (bs *BitSet) Len() int { return bs.n }
```

---

## 11. The `String()` Method (Stringer Interface)

Implementing `String() string` satisfies the `fmt.Stringer` interface, which `fmt.Println` and friends use automatically:

```go
type BitSet struct { ... }

func (bs *BitSet) String() string {
    var buf bytes.Buffer
    buf.WriteByte('{')
    for i, item := range bs.data {
        if item == 0 {
            continue
        }
        for j := 0; j < 64; j++ {
            if (item>>j)&1 != 0 {
                if buf.Len() > 1 {
                    buf.WriteByte(' ')
                }
                fmt.Fprintf(&buf, "%d", 64*i+j)
            }
        }
    }
    buf.WriteByte('}')
    return buf.String()
}
```

```go
var bs BitSet
bs.Add(5)
bs.Add(100)
bs.Add(64)
fmt.Println(&bs)               // {5 64 100}
```

**Important:** If `String()` has a pointer receiver (`*BitSet`), pass a pointer to `fmt.Println` to trigger it.

---

## 12. Common Patterns & Recipes

### Builder Pattern with Method Chaining

```go
type Query struct {
    table  string
    wheres []string
    limit  int
}

func (q *Query) From(table string) *Query {
    q.table = table
    return q
}

func (q *Query) Where(cond string) *Query {
    q.wheres = append(q.wheres, cond)
    return q
}

func (q *Query) Limit(n int) *Query {
    q.limit = n
    return q
}

// Usage:
q := new(Query).From("users").Where("age > 18").Limit(10)
```

### Functional Options Constructor

```go
type Server struct {
    addr    string
    timeout time.Duration
    logger  *log.Logger
}

type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) { s.timeout = d }
}

func WithLogger(l *log.Logger) Option {
    return func(s *Server) { s.logger = l }
}

func NewServer(addr string, opts ...Option) *Server {
    s := &Server{addr: addr, timeout: 30 * time.Second}
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage:
srv := NewServer(":8080", WithTimeout(10*time.Second), WithLogger(myLogger))
```

### Set Operations via Methods (BitSet Example)

```go
var a, b BitSet
a.Add(1); a.Add(2); a.Add(3)
b.Add(2); b.Add(3); b.Add(4)

a.Union(&b)               // a = {1 2 3 4}
a.Intersect(&b)           // a = {2 3 4}
a.Difference(&b)          // a = {}
```

### Copy Method (Deep Clone)

```go
func (bs *BitSet) Copy() *BitSet {
    if bs == nil {
        return nil
    }
    return &BitSet{
        data: slices.Clone(bs.data),
        n:    bs.n,
    }
}
```

---

## Quick Reference Card

```
DECLARATION     func (receiver Type) Method(params) returns { ... }
RECEIVER        Value: (p Point)    Pointer: (p *Point)
CALL            value.Method()      pointer.Method()
PROMOTION       Embedded fields promote their methods to the outer struct
METHOD VALUE    f := value.Method   → bound function, no receiver needed
METHOD EXPR     f := Type.Method    → unbound function, receiver is first arg
VISIBILITY      Uppercase = exported, lowercase = package-private
STRINGER        func (t Type) String() string → used by fmt.Println

RULES:
- Methods only on named types defined in the same package
- Pick ONE receiver type per type (value or pointer) — don't mix
- If any method needs pointer receiver → make all pointer receivers
- Nil pointer receivers are valid — guard with nil check
- Embedding promotes methods, but is NOT inheritance (has-a, not is-a)
```

## Value vs Pointer Receiver Decision

```
                  ┌──────────────────────────────────┐
                  │   Does the method mutate state?   │
                  └─────────┬──────────┬─────────────┘
                         YES│          │NO
                            ▼          ▼
                     Use *pointer    Is the struct large?
                                    ┌──┴──┐
                                   YES    NO
                                    │      │
                               *pointer   Do other methods
                                          use *pointer?
                                         ┌──┴──┐
                                        YES    NO
                                         │      │
                                    *pointer   value receiver OK
```

## Method Set Rules (Determines Interface Satisfaction)

| Type | Method set includes |
|------|-------------------|
| Value `T` | Only value-receiver methods `func (t T)` |
| Pointer `*T` | Both value-receiver AND pointer-receiver methods |

This is why: if an interface requires a pointer-receiver method, only `*T` (not `T`) satisfies it.

```go
type Stringer interface { String() string }

type MyType struct{}
func (m *MyType) String() string { return "hello" }

var _ Stringer = MyType{}     // ✗ compile error — MyType doesn't have String()
var _ Stringer = &MyType{}    // ✓ *MyType has String()
```
