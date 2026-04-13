# Go Interfaces Cheatsheet (§7.1 – §7.3)

## 7.1 Interfaces as Contracts

### Concrete vs Interface Types

| | Concrete Type | Interface Type |
|---|---|---|
| **What you know** | Exactly what it *is* and what you can *do* with it | Only what it can *do* (its methods) |
| **Representation** | Exposes exact value layout and intrinsic operations | Hides representation and internal structure |
| **Examples** | `int`, `[]string`, `*os.File` | `io.Writer`, `fmt.Stringer` |

### The `io.Writer` Contract

`io.Writer` is the classic example of how interfaces decouple code:

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

**How `fmt` uses it:**

```go
// Fprintf writes to ANY io.Writer — it doesn't care what's behind it
func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)

// Printf just wraps Fprintf with os.Stdout (a *os.File)
func Printf(format string, args ...interface{}) (int, error) {
    return Fprintf(os.Stdout, format, args...)
}

// Sprintf wraps Fprintf with a *bytes.Buffer (in-memory)
func Sprintf(format string, args ...interface{}) string {
    var buf bytes.Buffer
    Fprintf(&buf, format, args...)
    return buf.String()
}
```

**Key insight:** `Fprintf` doesn't know or care whether it writes to a file, a network connection, or memory. It only knows the value can `Write`. This is the **contract**: the caller provides something that can `Write`, and `Fprintf` does its formatting job.

### Substitutability

> The freedom to substitute one type for another that satisfies the same interface is called **substitutability** — a hallmark of object-oriented programming.

### ByteCounter Example — Custom `io.Writer`

Any type with a matching `Write` method satisfies `io.Writer` automatically:

```go
type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
    *c += ByteCounter(len(p))
    return len(p), nil
}

var c ByteCounter
fmt.Fprintf(&c, "hello, %s", "Dolly")
fmt.Println(c) // "12" — counted the bytes written
```

### `fmt.Stringer` — Another Key Interface

```go
type Stringer interface {
    String() string
}
```

Any type with a `String() string` method controls how it prints with `fmt.Println`, `%v`, etc.

---

## 7.2 Interface Types

### Defining Interfaces

An interface type specifies **a set of methods** that a concrete type must have.

```go
// Single-method interfaces (Go naming convention: method name + "er")
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}
```

### Embedding Interfaces (Composition)

Combine existing interfaces into larger ones using **embedding**:

```go
// Embedding style (preferred — concise)
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

This is equivalent to writing all methods out explicitly:

```go
// Explicit style (same effect, more verbose)
type ReadWriter interface {
    Read(p []byte) (n int, err error)
    Write(p []byte) (n int, err error)
}
```

You can even mix both styles:

```go
type ReadWriter interface {
    Read(p []byte) (n int, err error)
    Writer
}
```

**All three declarations are identical.** Order of methods doesn't matter — only the **set of methods** matters.

---

## 7.3 Interface Satisfaction

### The Rule

> A type **satisfies** an interface if it possesses **all the methods** the interface requires.

```
*os.File     satisfies  io.Reader, Writer, Closer, ReadWriter, ReadWriteCloser
*bytes.Buffer satisfies  io.Reader, Writer, ReadWriter
*bytes.Buffer does NOT satisfy io.Closer  (no Close method)
```

Go programmers say: "`*bytes.Buffer` **is an** `io.Writer`" — meaning it satisfies the interface.

### Assignment Rules

You can only assign a value to an interface variable if the value's type satisfies that interface:

```go
var w io.Writer
w = os.Stdout           // ✅ *os.File has Write method
w = new(bytes.Buffer)    // ✅ *bytes.Buffer has Write method
w = time.Second          // ❌ time.Duration lacks Write method

var rwc io.ReadWriteCloser
rwc = os.Stdout          // ✅ *os.File has Read, Write, Close
rwc = new(bytes.Buffer)  // ❌ *bytes.Buffer lacks Close
```

This works between interfaces too:

```go
w = rwc                  // ✅ ReadWriteCloser has Write (superset)
rwc = w                  // ❌ Writer lacks Close (subset)
```

**Rule of thumb:** A "bigger" interface (more methods) can always be assigned to a "smaller" one, but not the reverse.

### Pointer vs Value Receivers — The Subtlety

For a type `T`, methods with receiver `T` belong to **both** `T` and `*T`, but methods with receiver `*T` belong to **only `*T`**.

```go
type IntSet struct { /* ... */ }
func (*IntSet) String() string  // pointer receiver

var s IntSet

// Calling the method:
_ = s.String()          // ✅ OK — s is a variable, compiler takes &s automatically
_ = IntSet{}.String()   // ❌ compile error — non-addressable value, can't take &

// Satisfying the interface:
var _ fmt.Stringer = &s  // ✅ *IntSet has String method
var _ fmt.Stringer = s   // ❌ IntSet (value) lacks String method
```

| Receiver type | Value `T` satisfies? | Pointer `*T` satisfies? |
|---|---|---|
| `func (t T) Method()` | ✅ Yes | ✅ Yes |
| `func (t *T) Method()` | ❌ No | ✅ Yes |

### The Empty Interface `interface{}`

An interface with **zero methods** is satisfied by **every type**:

```go
var any interface{}
any = true
any = 12.34
any = "hello"
any = map[string]int{"one": 1}
any = new(bytes.Buffer)
```

This is why `fmt.Println` can accept arguments of any type. But you **can't do anything** directly with the value inside — you need a **type assertion** (§7.10) to get it back out.

### Interface Hiding — Only Methods Shown

An interface **conceals** the concrete type's other methods:

```go
os.Stdout.Write([]byte("hello"))  // ✅ OK
os.Stdout.Close()                 // ✅ OK — *os.File has Close

var w io.Writer
w = os.Stdout
w.Write([]byte("hello"))          // ✅ OK — io.Writer exposes Write
w.Close()                         // ❌ compile error — io.Writer has no Close
```

An interface with **more methods** tells you more about what values it holds, but **demands more** from types that implement it.

### Compile-Time Interface Check (Assertion Pattern)

You can verify that a type satisfies an interface at compile time without allocating variables:

```go
// Assert *bytes.Buffer satisfies io.Writer at compile time
var _ io.Writer = (*bytes.Buffer)(nil)
```

This is useful for documentation and catching mistakes early — if `*bytes.Buffer` ever loses its `Write` method, this line will fail to compile.

### Implicit Satisfaction — Go's Distinctive Feature

Unlike Java or C# where you must declare `implements`, **Go interfaces are satisfied implicitly**:

- No `implements` keyword needed
- Simply having the right methods is enough
- You can create new interfaces satisfied by **existing types you don't control**
- This decouples packages elegantly

```
Java:    class MyWriter implements io.Writer { ... }  // explicit
Go:      // just define Write([]byte)(int, error) and you're an io.Writer
```

---

## Summary Table

| Concept | Key Point |
|---|---|
| Interface | Abstract type defined by a set of methods |
| Contract | Interface defines what caller must provide and what function guarantees |
| Substitutability | Any type satisfying an interface can be used where that interface is expected |
| Embedding | Compose interfaces by naming other interfaces inside |
| Satisfaction | A type satisfies an interface by having all its methods — no declaration needed |
| Pointer receiver | Only `*T` satisfies interfaces requiring `*T` methods, not `T` itself |
| Empty interface | `interface{}` is satisfied by all types — holds any value |
| Compile-time check | `var _ SomeInterface = (*MyType)(nil)` verifies satisfaction |
