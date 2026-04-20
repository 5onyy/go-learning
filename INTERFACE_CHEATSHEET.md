# Go Interfaces Cheatsheet (§7.1 – §7.6)

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

## 7.4 Parsing Flags with `flag.Value`

The `flag.Value` interface lets you define **custom command-line flag types**:

```go
package flag

type Value interface {
    String() string    // format for -help output
    Set(string) error  // parse user input and update the value
}
```

- `String()` formats the current value (used in help messages) — so every `flag.Value` is also a `fmt.Stringer`
- `Set()` is the **inverse of `String()`** — it parses a string and updates the flag value

### Example: Temperature Flag

A `celsiusFlag` that accepts both Celsius and Fahrenheit input:

```go
type celsiusFlag struct{ Celsius }

func (f *celsiusFlag) Set(s string) error {
    var unit string
    var value float64
    fmt.Sscanf(s, "%f%s", &value, &unit)
    switch unit {
    case "C", "°C":
        f.Celsius = Celsius(value)
        return nil
    case "F", "°F":
        f.Celsius = FToC(Fahrenheit(value))
        return nil
    }
    return fmt.Errorf("invalid temperature %q", s)
}
```

**Key trick:** `celsiusFlag` **embeds** `Celsius`, which already has a `String()` method. So it only needs to define `Set()` to satisfy `flag.Value`.

### Registering with `flag.Var`

```go
func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
    f := celsiusFlag{value}
    flag.CommandLine.Var(&f, name, usage)
    return &f.Celsius
}
```

Usage:

```bash
$ ./tempflag                  # 20°C  (default)
$ ./tempflag -temp -18C       # -18°C
$ ./tempflag -temp 212°F      # 100°C (converted)
$ ./tempflag -temp 273.15K    # error: invalid temperature
```

### Why the Help Message Shows °C for Default 20.0

The `-help` output calls `String()` on the flag value. Since `celsiusFlag` embeds `Celsius`, its `String()` method formats `20.0` as `"20°C"` — the `°C` comes from the `String` method, not from the raw `20.0` value.

---

## 7.5 Interface Values

An interface value has **two components** internally:

| Component | What it holds |
|---|---|
| **Dynamic type** | Type descriptor (e.g., `*os.File`) |
| **Dynamic value** | Actual value of that type (e.g., pointer to stdout) |

### Interface Value Lifecycle

```go
var w io.Writer       // type: nil,          value: nil    → nil interface
w = os.Stdout         // type: *os.File,     value: →stdout → non-nil
w = new(bytes.Buffer) // type: *bytes.Buffer, value: →buf   → non-nil
w = nil               // type: nil,          value: nil    → nil again
```

### Nil Interface vs Non-Nil

- **Nil interface:** both type and value are nil → `w == nil` is `true`
- Calling any method on a nil interface → **panic**

```go
var w io.Writer
w.Write([]byte("hello")) // panic: nil pointer dereference
```

### Dynamic Dispatch

When you call a method through an interface, the compiler uses **dynamic dispatch**:

1. Look up the method address from the type descriptor
2. Make an indirect call to that address
3. Pass a copy of the dynamic value as the receiver

```go
w = os.Stdout
w.Write([]byte("hello"))  // calls (*os.File).Write

w = new(bytes.Buffer)
w.Write([]byte("hello"))  // calls (*bytes.Buffer).Write
```

The same `w.Write()` call dispatches to completely different methods depending on the dynamic type.

### Comparing Interface Values

Interface values can be compared with `==` and `!=`:

```go
// Equal if:
// 1. Both are nil, OR
// 2. Same dynamic type AND equal dynamic values

var w1, w2 io.Writer
fmt.Println(w1 == w2) // true (both nil)
```

**Danger:** comparing interfaces whose dynamic type is **not comparable** (slice, map, function) causes a **panic**:

```go
var x interface{} = []int{1, 2, 3}
fmt.Println(x == x) // panic: comparing uncomparable type []int
```

### Inspecting Dynamic Type with `%T`

```go
var w io.Writer
fmt.Printf("%T\n", w) // "<nil>"

w = os.Stdout
fmt.Printf("%T\n", w) // "*os.File"

w = new(bytes.Buffer)
fmt.Printf("%T\n", w) // "*bytes.Buffer"
```

### 7.5.1 Caveat: An Interface Containing a Nil Pointer Is Non-Nil

This is the **most common trap** for Go beginners:

```go
const debug = false

func main() {
    var buf *bytes.Buffer           // nil pointer
    if debug {
        buf = new(bytes.Buffer)     // only allocated when debug=true
    }
    f(buf) // BUG: passes a non-nil interface containing a nil pointer
}

func f(out io.Writer) {
    if out != nil {                 // true! (interface has type *bytes.Buffer)
        out.Write([]byte("done!\n")) // panic: nil pointer dereference
    }
}
```

**What happened:**

| | type | value | `== nil`? |
|---|---|---|---|
| **Nil interface** | nil | nil | `true` |
| **Interface with nil pointer** | `*bytes.Buffer` | nil | `false` |

The interface has a **dynamic type** (`*bytes.Buffer`), so `out != nil` is `true` even though the pointer inside is nil.

**The fix:** declare `buf` as the interface type so it stays truly nil when unassigned:

```go
var buf io.Writer               // nil interface (not *bytes.Buffer)
if debug {
    buf = new(bytes.Buffer)
}
f(buf)                           // OK: buf is a true nil interface
```

---

## 7.6 Sorting with `sort.Interface`

`sort.Sort` works with **any type** that implements these three methods:

```go
package sort

type Interface interface {
    Len() int
    Less(i, j int) bool  // i, j are indices
    Swap(i, j int)
}
```

The sort algorithm only needs: how many elements, how to compare two, and how to swap two.

### Sorting Strings

```go
type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

sort.Sort(StringSlice(names))
// or simply:
sort.Strings(names)
```

### Sorting Custom Structs

Define a new type per sort order. Only `Less` changes — `Len` and `Swap` are always the same for slices:

```go
type byArtist []*Track
func (x byArtist) Len() int           { return len(x) }
func (x byArtist) Less(i, j int) bool { return x[i].Artist < x[j].Artist }
func (x byArtist) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

sort.Sort(byArtist(tracks))

type byYear []*Track
func (x byYear) Len() int           { return len(x) }
func (x byYear) Less(i, j int) bool { return x[i].Year < x[j].Year }
func (x byYear) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

sort.Sort(byYear(tracks))
```

### Reverse Sorting with `sort.Reverse`

No need for a separate type — just wrap:

```go
sort.Sort(sort.Reverse(byArtist(tracks)))
```

**How `sort.Reverse` works** — it uses **composition** (embedding):

```go
package sort

type reverse struct{ Interface }  // embeds sort.Interface

func (r reverse) Less(i, j int) bool { return r.Interface.Less(j, i) }  // flipped!

func Reverse(data Interface) Interface { return reverse{data} }
```

- `Len` and `Swap` come from the embedded `Interface` for free
- Only `Less` is overridden with flipped indices `(j, i)` instead of `(i, j)`

### `customSort` — Avoid Defining a New Type Per Sort Order

Combine the slice with a comparison function:

```go
type customSort struct {
    t    []*Track
    less func(x, y *Track) bool
}

func (x customSort) Len() int           { return len(x.t) }
func (x customSort) Less(i, j int) bool { return x.less(x.t[i], x.t[j]) }
func (x customSort) Swap(i, j int)      { x.t[i], x.t[j] = x.t[j], x.t[i] }
```

Multi-tier sort (Title → Year → Length):

```go
sort.Sort(customSort{tracks, func(x, y *Track) bool {
    if x.Title != y.Title {
        return x.Title < y.Title
    }
    if x.Year != y.Year {
        return x.Year < y.Year
    }
    if x.Length != y.Length {
        return x.Length < y.Length
    }
    return false
}})
```

### Convenience Functions

The `sort` package provides shortcuts for common types:

| Slice type | Sort | Check sorted |
|---|---|---|
| `[]int` | `sort.Ints(a)` | `sort.IntsAreSorted(a)` |
| `[]string` | `sort.Strings(a)` | `sort.StringsAreSorted(a)` |
| `[]float64` | `sort.Float64s(a)` | `sort.Float64sAreSorted(a)` |

Reverse any of them:

```go
sort.Sort(sort.Reverse(sort.IntSlice(values)))
```

### `sort.IsSorted` — Check Without Sorting

Tests whether a sequence is already sorted (at most n-1 comparisons):

```go
values := []int{3, 1, 4, 1}
fmt.Println(sort.IntsAreSorted(values)) // false
sort.Ints(values)
fmt.Println(sort.IntsAreSorted(values)) // true
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
| `flag.Value` | Implement `String()` + `Set(string) error` for custom flag types |
| Interface value | Two components: dynamic type + dynamic value |
| Nil trap | Interface with nil pointer is **non-nil** — use interface type for optional params |
| Dynamic dispatch | Method calls through interfaces are resolved at runtime via type descriptor |
| `sort.Interface` | `Len()`, `Less(i,j)`, `Swap(i,j)` — sort anything |
| `sort.Reverse` | Wraps a `sort.Interface` with flipped `Less` indices via embedding |
| `customSort` | Struct with slice + comparison func — avoids one type per sort order |
