# Go Language Cheatsheet

> A concise reference covering core Go concepts, operators, data structures, and I/O.

---

## Table of Contents
1. [Basic Go](#1-basic-go)
2. [Arithmetic Operators](#2-arithmetic-operators)
3. [Standard Data Structures](#3-standard-data-structures)
4. [Standard Input / I/O](#4-standard-input--io)

---

## 1. Basic Go

### Program Structure

```go
package main          // every executable must be in package main

import "fmt"          // import standard packages

func main() {         // entry point
    fmt.Println("Hello, Go!")
}
```

### Variables & Constants

```go
// Explicit declaration
var name string = "Alice"
var age  int    = 30

// Short declaration (inside functions only)
city := "Hanoi"

// Multiple variables
var x, y int = 1, 2

// Zero values (default when not initialized)
var i int     // 0
var f float64 // 0.0
var b bool    // false
var s string  // ""

// Constants
const Pi      = 3.14159
const MaxSize = 100

// Typed constants
const Greeting string = "Hello"

// iota — auto-incrementing constants
type Direction int
const (
    North Direction = iota // 0
    East                   // 1
    South                  // 2
    West                   // 3
)
```

### Basic Types

| Type        | Size      | Example              |
|-------------|-----------|----------------------|
| `bool`      | 1 byte    | `true`, `false`      |
| `int`       | platform  | `-42`, `0`, `100`    |
| `int8/16/32/64` | fixed | `int32(42)`          |
| `uint`      | platform  | `uint(42)`           |
| `float32`   | 4 bytes   | `3.14`               |
| `float64`   | 8 bytes   | `3.141592653589793`  |
| `complex64` | 8 bytes   | `1 + 2i`             |
| `complex128`| 16 bytes  | `1.5 + 2.5i`         |
| `string`    | variable  | `"hello"`            |
| `byte`      | alias u8  | `'A'`                |
| `rune`      | alias i32 | `'あ'`               |

### Type Conversion

```go
var i int   = 42
var f float64 = float64(i)   // explicit cast — Go never coerces implicitly
var u uint  = uint(f)

// String ↔ number (use strconv, not casting)
import "strconv"
s := strconv.Itoa(42)          // int  → "42"
n, err := strconv.Atoi("42")   // "42" → int  (may fail)

// []byte ↔ string
b := []byte("hello")
s2 := string(b)
```

### Control Flow

```go
// if / else if / else
if x > 0 {
    fmt.Println("positive")
} else if x == 0 {
    fmt.Println("zero")
} else {
    fmt.Println("negative")
}

// if with init statement
if v, err := someFunc(); err != nil {
    fmt.Println("error:", err)
} else {
    fmt.Println("value:", v)
}

// for — Go's only loop keyword
for i := 0; i < 5; i++ { }           // C-style
for i < 10 { i++ }                    // while-style
for { break }                         // infinite loop

// range
nums := []int{10, 20, 30}
for idx, val := range nums {
    fmt.Println(idx, val)
}
for _, val := range nums { }          // ignore index
for idx := range nums { }             // index only

// switch
switch day {
case "Mon", "Tue", "Wed", "Thu", "Fri":
    fmt.Println("Weekday")
case "Sat", "Sun":
    fmt.Println("Weekend")
default:
    fmt.Println("Unknown")
}

// switch with no condition (acts like if-else chain)
switch {
case x < 0:  fmt.Println("neg")
case x == 0: fmt.Println("zero")
default:     fmt.Println("pos")
}

// defer — runs at function return, LIFO order
defer fmt.Println("world")
fmt.Println("hello")  // prints: hello\nworld

// goto (rare)
goto Label
Label:
    fmt.Println("jumped here")
```

### Functions

```go
// Basic
func add(a, b int) int { return a + b }

// Multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 { return 0, fmt.Errorf("division by zero") }
    return a / b, nil
}

// Named return values
func minMax(nums []int) (min, max int) {
    min, max = nums[0], nums[0]
    for _, v := range nums {
        if v < min { min = v }
        if v > max { max = v }
    }
    return  // bare return
}

// Variadic
func sum(nums ...int) int {
    total := 0
    for _, n := range nums { total += n }
    return total
}
sum(1, 2, 3)
sum(nums...)   // spread a slice

// First-class functions
double := func(n int) int { return n * 2 }
apply := func(f func(int) int, v int) int { return f(v) }

// Closures
func counter() func() int {
    n := 0
    return func() int { n++; return n }
}
c := counter()
c() // 1
c() // 2
```

### Pointers

```go
x := 42
p := &x        // p is *int; holds address of x
fmt.Println(*p) // dereference → 42
*p = 99         // modify through pointer → x is now 99

// new() allocates zeroed memory
p2 := new(int)  // *int pointing to 0
*p2 = 7

// Go has NO pointer arithmetic
```

### Structs & Methods

```go
type Point struct {
    X, Y float64
}

// Constructor pattern (no built-in constructors)
func NewPoint(x, y float64) Point { return Point{x, y} }

// Value receiver (works on a copy)
func (p Point) String() string {
    return fmt.Sprintf("(%g, %g)", p.X, p.Y)
}

// Pointer receiver (can mutate original)
func (p *Point) Scale(factor float64) {
    p.X *= factor
    p.Y *= factor
}

pt := Point{3, 4}
pt.Scale(2)             // Go auto-takes address
fmt.Println(pt.String()) // (6, 8)

// Anonymous / embedded structs
type Circle struct {
    Point               // embeds all Point fields & methods
    Radius float64
}
c := Circle{Point: Point{1, 2}, Radius: 5}
fmt.Println(c.X) // promoted field
```

### Interfaces

```go
type Shape interface {
    Area() float64
    Perimeter() float64
}

type Rect struct{ W, H float64 }
func (r Rect) Area() float64      { return r.W * r.H }
func (r Rect) Perimeter() float64 { return 2 * (r.W + r.H) }

var s Shape = Rect{3, 4}   // implicit satisfaction — no "implements" keyword

// Empty interface / any
func printAny(v interface{}) { fmt.Println(v) }
func printAny2(v any) { fmt.Println(v) }   // 'any' is alias since Go 1.18

// Type assertion
val, ok := s.(Rect)     // safe form
val2 := s.(Rect)        // panics if wrong type

// Type switch
switch t := v.(type) {
case int:    fmt.Println("int:", t)
case string: fmt.Println("string:", t)
default:     fmt.Println("other:", t)
}
```

### Error Handling

```go
// errors package
import "errors"

var ErrNotFound = errors.New("not found")

// fmt.Errorf with wrapping (%w)
err := fmt.Errorf("open config: %w", ErrNotFound)

// Unwrap / check
errors.Is(err, ErrNotFound)          // true (walks the chain)
errors.As(err, &target)              // extract concrete type

// Custom error type
type ValidationError struct {
    Field   string
    Message string
}
func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// panic / recover (use sparingly — prefer errors)
defer func() {
    if r := recover(); r != nil {
        fmt.Println("recovered:", r)
    }
}()
panic("something went very wrong")
```

### Goroutines & Channels (intro)

```go
go func() { fmt.Println("runs concurrently") }()  // goroutine

ch := make(chan int)          // unbuffered channel
go func() { ch <- 42 }()     // send in goroutine
v := <-ch                    // receive blocks until value arrives

bch := make(chan string, 3)   // buffered channel (size 3)
bch <- "a"                   // non-blocking if buffer not full

// Select — multiplex on channels
select {
case msg := <-ch1:
    fmt.Println("ch1:", msg)
case ch2 <- 99:
    fmt.Println("sent to ch2")
case <-time.After(1 * time.Second):
    fmt.Println("timeout")
}

// sync.WaitGroup
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        fmt.Println(n)
    }(i)
}
wg.Wait()
```

---

## 2. Arithmetic Operators

### Basic Operators

| Operator | Name           | Example           | Result |
|----------|----------------|-------------------|--------|
| `+`      | Addition       | `5 + 3`           | `8`    |
| `-`      | Subtraction    | `5 - 3`           | `2`    |
| `*`      | Multiplication | `5 * 3`           | `15`   |
| `/`      | Division       | `7 / 2`           | `3` (int) |
| `%`      | Modulo         | `7 % 2`           | `1`    |

```go
a, b := 7, 2

fmt.Println(a + b)   // 9
fmt.Println(a - b)   // 5
fmt.Println(a * b)   // 14
fmt.Println(a / b)   // 3  — integer division truncates toward zero
fmt.Println(a % b)   // 1

// Float division
fa, fb := 7.0, 2.0
fmt.Println(fa / fb) // 3.5
```

### Increment / Decrement

```go
x := 5
x++   // x = 6   (statement, NOT expression — no ++x in Go)
x--   // x = 5
// y := x++  → compile error
```

### Assignment Operators

```go
x := 10
x += 3   // x = 13
x -= 2   // x = 11
x *= 4   // x = 44
x /= 4   // x = 11
x %= 3   // x = 2
```

### Bitwise Operators

| Operator | Name         | Example      | Result (binary)         |
|----------|--------------|--------------|-------------------------|
| `&`      | AND          | `5 & 3`      | `0101 & 0011 = 0001` → 1 |
| `\|`     | OR           | `5 \| 3`     | `0101 \| 0011 = 0111` → 7 |
| `^`      | XOR          | `5 ^ 3`      | `0101 ^ 0011 = 0110` → 6 |
| `^`      | NOT (unary)  | `^5`         | bitwise complement       |
| `&^`     | AND NOT      | `5 &^ 3`     | `0101 &^ 0011 = 0100` → 4 |
| `<<`     | Left shift   | `1 << 3`     | `8`                      |
| `>>`     | Right shift  | `16 >> 2`    | `4`                      |

```go
a, b := 5, 3
fmt.Println(a & b)   // 1
fmt.Println(a | b)   // 7
fmt.Println(a ^ b)   // 6
fmt.Println(^a)      // -6 (two's complement)
fmt.Println(a &^ b)  // 4  (clear bits)
fmt.Println(1 << 4)  // 16
fmt.Println(32 >> 2) // 8
```

### math Package — Advanced Operations

```go
import "math"

math.Abs(-3.5)          // 3.5
math.Sqrt(16)           // 4.0
math.Cbrt(27)           // 3.0   (cube root)
math.Pow(2, 10)         // 1024.0
math.Exp(1)             // e ≈ 2.718
math.Log(math.E)        // 1.0   (natural log)
math.Log2(8)            // 3.0
math.Log10(100)         // 2.0
math.Ceil(1.2)          // 2.0
math.Floor(1.9)         // 1.0
math.Round(1.5)         // 2.0
math.Trunc(3.9)         // 3.0
math.Mod(7.5, 2.0)      // 1.5   (float modulo)
math.Hypot(3, 4)        // 5.0   (√(a²+b²))
math.Sin(math.Pi / 2)   // 1.0
math.Cos(0)             // 1.0
math.Tan(math.Pi / 4)   // ~1.0
math.Min(3.0, 5.0)      // 3.0
math.Max(3.0, 5.0)      // 5.0
math.Inf(1)             // +Inf
math.IsNaN(math.NaN())  // true
math.IsInf(math.Inf(1), 1) // true

// Integer min/max (Go 1.21+)
import "cmp"
cmp.Min(3, 5)   // 3
cmp.Max(3, 5)   // 5
```

### big.Int / big.Float — Arbitrary Precision

```go
import "math/big"

// big.Int — arbitrary precision integers
a := new(big.Int).SetInt64(1_000_000_000)
b := new(big.Int).SetInt64(1_000_000_000)
product := new(big.Int).Mul(a, b)
fmt.Println(product)  // 1000000000000000000

// Factorial of 100
n := big.NewInt(100)
result := new(big.Int).MulRange(1, n.Int64())
fmt.Println(result)

// big.Float — arbitrary precision floating point
pi, _, _ := big.ParseFloat("3.14159265358979323846264338327950288", 10, 256, big.ToNearestEven)
fmt.Println(pi.Text('f', 30))  // 30 decimal places
```

### Operator Precedence (highest → lowest)

```
5  (highest)   *   /   %   <<  >>  &   &^
4              +   -   |   ^
3              ==  !=  <   <=  >   >=
2              &&
1  (lowest)    ||
```

```go
fmt.Println(2 + 3*4)        // 14  (not 20)
fmt.Println((2 + 3) * 4)    // 20
fmt.Println(1 + 2<<3)       // 17  (1 + 16)
fmt.Println(true || false && false) // true  (&& binds tighter)
```

---

## 3. Standard Data Structures

### Array

Fixed-size, value type (copying an array copies all elements).

```go
// Declaration
var a [5]int                      // [0 0 0 0 0]
b := [3]string{"foo", "bar", "baz"}
c := [...]int{1, 2, 3, 4}         // compiler infers length → [4]int

// Access
a[0] = 42
fmt.Println(len(a))   // 5

// Iterate
for i, v := range b {
    fmt.Println(i, v)
}

// Multidimensional
matrix := [2][3]int{{1, 2, 3}, {4, 5, 6}}
fmt.Println(matrix[1][2]) // 6
```

### Slice

Dynamic-size view into an underlying array. The most-used sequence type.

```go
// Create
s := []int{1, 2, 3}
s2 := make([]int, 5)         // len=5, cap=5, zero-filled
s3 := make([]int, 3, 10)     // len=3, cap=10

// Slice of slice (shares memory!)
a := []int{0, 1, 2, 3, 4}
s4 := a[1:3]   // [1 2]   (low inclusive, high exclusive)
s5 := a[:3]    // [0 1 2]
s6 := a[2:]    // [2 3 4]
s7 := a[:]     // full copy-like view

// Append (may allocate new backing array)
s = append(s, 4, 5)
s = append(s, []int{6, 7}...)  // spread another slice

// Copy (copies min(len(dst), len(src)) elements)
dst := make([]int, len(s))
n := copy(dst, s)

// len vs cap
fmt.Println(len(s3), cap(s3))   // 3 10

// 2D slice
grid := make([][]int, 3)
for i := range grid {
    grid[i] = make([]int, 4)
}

// Delete element at index i (order preserved)
i := 2
s = append(s[:i], s[i+1:]...)

// Delete element (order NOT preserved — fast)
s[i] = s[len(s)-1]
s = s[:len(s)-1]

// Sort
import "sort"
sort.Ints(s)
sort.Strings(strSlice)
sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
sort.SliceStable(...)   // preserves order of equal elements

// Search (binary search on sorted slice)
idx := sort.SearchInts(s, 3)   // index of first 3
```

### Map

Hash map (unordered key-value pairs).

```go
// Create
m := map[string]int{"a": 1, "b": 2}
m2 := make(map[string]int)
m3 := make(map[string]int, 100)  // hint initial capacity

// CRUD
m["c"] = 3                  // insert / update
v := m["a"]                 // get (returns zero value if missing)
v, ok := m["z"]             // ok=false if key absent
delete(m, "b")              // delete key

// Iterate (order is randomized each run)
for k, v := range m {
    fmt.Println(k, v)
}

// Nested map
nested := map[string]map[string]int{
    "group1": {"x": 1, "y": 2},
}

// Map of slices
graph := make(map[string][]string)
graph["a"] = append(graph["a"], "b", "c")

// Check existence
if val, exists := m["key"]; exists {
    fmt.Println(val)
}
```

### String

Immutable byte sequence (UTF-8 encoded).

```go
s := "Hello, 世界"

// Length
len(s)                     // byte count, not rune count
utf8.RuneCountInString(s)  // rune (character) count

// Indexing gives bytes, range gives runes
for i, r := range s { fmt.Println(i, string(r)) }

// strings package
import "strings"
strings.Contains(s, "Hello")         // true
strings.HasPrefix(s, "Hello")        // true
strings.HasSuffix(s, "界")           // true
strings.Index(s, "世")               // byte position
strings.Count(s, "l")               // 2
strings.ToUpper(s)
strings.ToLower(s)
strings.TrimSpace("  hi  ")         // "hi"
strings.Trim(s, "!")                // trim chars from both ends
strings.TrimLeft / TrimRight
strings.Replace(s, "l", "L", -1)   // -1 = all occurrences
strings.Split("a,b,c", ",")        // []string{"a","b","c"}
strings.Join([]string{"a","b"}, "-") // "a-b"
strings.Repeat("ab", 3)            // "ababab"
strings.Fields("foo bar  baz")     // split by whitespace

// Builder (efficient string concatenation)
var sb strings.Builder
sb.WriteString("Hello")
sb.WriteRune(',')
sb.WriteString(" Go!")
result := sb.String()              // "Hello, Go!"
```

### Struct (as data structure)

Already covered in §1. Key patterns:

```go
// Comparable structs can be map keys
type Point struct{ X, Y int }
visited := map[Point]bool{}
visited[Point{1, 2}] = true

// Slice of structs
type Person struct{ Name string; Age int }
people := []Person{{"Alice", 30}, {"Bob", 25}}
sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})
```

### Stack & Queue (using slices)

```go
// Stack (LIFO)
stack := []int{}
stack = append(stack, 1)           // push
top := stack[len(stack)-1]         // peek
stack = stack[:len(stack)-1]       // pop

// Queue (FIFO)
queue := []int{}
queue = append(queue, 1)           // enqueue
front := queue[0]                  // peek
queue = queue[1:]                  // dequeue (may leave memory; see ring buffer for perf)
```

### Heap / Priority Queue

```go
import "container/heap"

// Min-heap of ints
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }   // change to > for max-heap
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

h := &IntHeap{3, 1, 4, 1, 5}
heap.Init(h)
heap.Push(h, 2)
min := heap.Pop(h).(int)   // 1
```

### Linked List

```go
import "container/list"

l := list.New()
l.PushBack(1)
l.PushBack(2)
l.PushFront(0)

for e := l.Front(); e != nil; e = e.Next() {
    fmt.Println(e.Value)   // 0 1 2
}

l.Remove(l.Front())
fmt.Println(l.Len())   // 2
```

### Ring Buffer

```go
import "container/ring"

r := ring.New(5)          // ring of 5 elements
for i := 0; i < r.Len(); i++ {
    r.Value = i
    r = r.Next()
}
r.Do(func(v any) { fmt.Println(v) })  // 0 1 2 3 4
```

### sync.Map (concurrent map)

```go
import "sync"

var sm sync.Map
sm.Store("key", 42)
v, ok := sm.Load("key")
sm.LoadOrStore("key2", 99)
sm.Delete("key")
sm.Range(func(k, v any) bool {
    fmt.Println(k, v)
    return true   // return false to stop
})
```

---

## 4. Standard Input / I/O

### fmt — Formatted I/O

```go
import "fmt"

// Read from stdin
var name string
fmt.Scan(&name)                       // reads one word
fmt.Scanln(&name)                     // reads until newline
fmt.Scanf("%s %d", &name, &age)       // formatted read

// Read multiple values
var a, b int
fmt.Scan(&a, &b)   // separated by whitespace

// Print to stdout
fmt.Print("no newline")
fmt.Println("with newline")
fmt.Printf("name=%s age=%d\n", name, age)

// Sprint / Sprintf — return formatted string
s := fmt.Sprintf("Hello, %s! (age %d)", name, age)
msg := fmt.Sprint("a", "b", "c")       // "a b c" (spaces between)

// Errorf
err := fmt.Errorf("failed at step %d: %w", 3, someErr)
```

### bufio.Scanner — Line-by-Line Input

```go
import (
    "bufio"
    "os"
    "strings"
)

// Read lines from stdin
scanner := bufio.NewScanner(os.Stdin)
for scanner.Scan() {                  // false on EOF or error
    line := scanner.Text()
    fmt.Println("got:", line)
}
if err := scanner.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "error:", err)
}

// Increase buffer for long lines (default 64 KiB)
scanner.Buffer(make([]byte, 1<<20), 1<<20)

// Scan words
scanner.Split(bufio.ScanWords)
for scanner.Scan() {
    fmt.Println(scanner.Text())
}

// Scan from a string
r := strings.NewReader("one two three")
ws := bufio.NewScanner(r)
ws.Split(bufio.ScanWords)
for ws.Scan() { fmt.Println(ws.Text()) }
```

### bufio.Reader — Fine-grained Input

```go
reader := bufio.NewReader(os.Stdin)

// Read until delimiter (includes delimiter)
line, err := reader.ReadString('\n')
line = strings.TrimRight(line, "\n")

// ReadLine (low-level, avoids allocations)
lineBytes, isPrefix, err := reader.ReadLine()

// Read a single byte
b, err := reader.ReadByte()

// Read a single rune (UTF-8 aware)
r, size, err := reader.ReadRune()

// Peek without consuming
peeked, _ := reader.Peek(4)   // look at next 4 bytes
```

### os — Standard Streams & Environment

```go
import "os"

os.Stdin   // *os.File  — standard input
os.Stdout  // *os.File  — standard output
os.Stderr  // *os.File  — standard error

// Write to stderr
fmt.Fprintln(os.Stderr, "error occurred")

// Program arguments
args := os.Args            // []string, args[0] = program name

// Environment variables
home := os.Getenv("HOME")
os.Setenv("MY_VAR", "value")
all := os.Environ()        // []string of "KEY=VALUE"

// Exit
os.Exit(1)   // exits immediately; deferred funcs do NOT run
```

### Reading Files

```go
import (
    "bufio"
    "io"
    "os"
)

// ── Simplest: read entire file at once ──────────────────────
data, err := os.ReadFile("file.txt")   // []byte  (Go 1.16+)
if err != nil { panic(err) }
content := string(data)

// ── Open → read → close ────────────────────────────────────
f, err := os.Open("file.txt")          // read-only
if err != nil { panic(err) }
defer f.Close()

// Read into a fixed buffer
buf := make([]byte, 1024)
n, err := f.Read(buf)                  // reads up to len(buf) bytes
fmt.Printf("read %d bytes\n", n)

// ── Read line-by-line efficiently ──────────────────────────
f2, _ := os.Open("file.txt")
defer f2.Close()

scanner := bufio.NewScanner(f2)
for scanner.Scan() {
    fmt.Println(scanner.Text())
}

// ── io.ReadAll — read everything from any Reader ───────────
f3, _ := os.Open("file.txt")
defer f3.Close()
all, _ := io.ReadAll(f3)

// ── bufio.NewReader ─────────────────────────────────────────
f4, _ := os.Open("file.txt")
defer f4.Close()
br := bufio.NewReader(f4)
for {
    line, err := br.ReadString('\n')
    line = strings.TrimRight(line, "\r\n")
    if line != "" { fmt.Println(line) }
    if err == io.EOF { break }
    if err != nil { panic(err) }
}
```

### Writing Files

```go
// ── Simplest: write in one call ─────────────────────────────
os.WriteFile("out.txt", []byte("hello\n"), 0644)  // Go 1.16+

// ── Create / truncate then write ────────────────────────────
f, _ := os.Create("out.txt")      // creates or truncates
defer f.Close()
f.WriteString("line 1\n")
fmt.Fprintln(f, "line 2")

// ── Append to existing file ─────────────────────────────────
f2, _ := os.OpenFile("out.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
defer f2.Close()
fmt.Fprintln(f2, "appended line")

// ── Buffered write (flush for performance) ──────────────────
f3, _ := os.Create("out.txt")
defer f3.Close()
bw := bufio.NewWriter(f3)
fmt.Fprintln(bw, "buffered line")
bw.Flush()   // IMPORTANT: flush to disk
```

### strings.NewReader / bytes.Buffer — In-memory I/O

```go
import (
    "bytes"
    "strings"
    "io"
)

// Read from a string as if it were a file
r := strings.NewReader("hello world\n")
io.Copy(os.Stdout, r)

// bytes.Buffer — read/write in memory
var buf bytes.Buffer
buf.WriteString("Hello")
buf.WriteByte(',')
buf.WriteString(" World!")
fmt.Println(buf.String())      // Hello, World!
line, _ := buf.ReadString('!')  // reads until '!'
```

### io.Reader / io.Writer — Interface Composition

```go
import "io"

// Copy between any Reader and Writer
src, _ := os.Open("src.txt")
dst, _ := os.Create("dst.txt")
n, err := io.Copy(dst, src)       // copies all bytes

// TeeReader — read and simultaneously write to another writer
import "io"
tee := io.TeeReader(src, dst)     // reads from src, writes to dst simultaneously
io.ReadAll(tee)

// LimitReader — cap how many bytes can be read
limited := io.LimitReader(src, 512)  // at most 512 bytes

// MultiReader — chain multiple readers
combined := io.MultiReader(r1, r2, r3)

// MultiWriter — fan-out to multiple writers
w := io.MultiWriter(os.Stdout, logFile)
fmt.Fprintln(w, "goes to both stdout and file")

// Pipe — synchronous in-memory pipe
pr, pw := io.Pipe()
go func() {
    fmt.Fprintln(pw, "piped data")
    pw.Close()
}()
io.Copy(os.Stdout, pr)
```

### CSV Reading / Writing

```go
import "encoding/csv"

// Read CSV
f, _ := os.Open("data.csv")
defer f.Close()
reader := csv.NewReader(f)
reader.Comma = ','             // default; change for TSV: '\t'

records, err := reader.ReadAll()    // [][]string
// or read row-by-row:
for {
    record, err := reader.Read()    // []string
    if err == io.EOF { break }
    if err != nil { panic(err) }
    fmt.Println(record)
}

// Write CSV
out, _ := os.Create("out.csv")
defer out.Close()
writer := csv.NewWriter(out)
writer.Write([]string{"name", "age"})
writer.Write([]string{"Alice", "30"})
writer.WriteAll([][]string{{"Bob", "25"}, {"Carol", "28"}})
writer.Flush()
```

### JSON Reading / Writing

```go
import "encoding/json"

// Marshal (struct → JSON bytes)
type User struct {
    Name  string `json:"name"`
    Email string `json:"email,omitempty"`
    Age   int    `json:"age"`
}
u := User{Name: "Alice", Age: 30}
data, err := json.Marshal(u)
data, err  = json.MarshalIndent(u, "", "  ")   // pretty print

// Unmarshal (JSON bytes → struct)
jsonStr := `{"name":"Bob","age":25}`
var u2 User
err = json.Unmarshal([]byte(jsonStr), &u2)

// Streaming encode/decode (for files / HTTP bodies)
enc := json.NewEncoder(os.Stdout)
enc.SetIndent("", "  ")
enc.Encode(u)

dec := json.NewDecoder(os.Stdin)
var u3 User
err = dec.Decode(&u3)

// Generic JSON with map
var generic map[string]any
json.Unmarshal([]byte(jsonStr), &generic)
```

### Flag Parsing (CLI args)

```go
import "flag"

name    := flag.String("name", "World", "your name")
verbose := flag.Bool("verbose", false, "enable verbose mode")
port    := flag.Int("port", 8080, "port number")
flag.Parse()

// Usage: ./app -name Alice -verbose -port 9000
fmt.Println("Hello,", *name)
remaining := flag.Args()   // non-flag arguments
```

---

## Quick Reference Card

```
VARIABLES     var x int = 5  |  x := 5  |  const C = 10
TYPES         int float64 bool string byte rune complex128
ZERO VALUES   0  0.0  false  ""  0  0  0+0i
FUNCTIONS     func f(a, b int) (int, error)
LOOPS         for init; cond; post  |  for cond  |  for range
SWITCH        no break needed; fallthrough to continue
DEFER         LIFO, runs at function exit
GOROUTINE     go f()
CHANNEL       ch := make(chan T [, buf])  |  ch <- v  |  v := <-ch
ERRORS        if err != nil { return err }
INTERFACES    implicit (no "implements"), satisfied by method set
GENERICS      func F[T any](v T) T  (Go 1.18+)
POINTERS      &x (address)  |  *p (dereference)  |  new(T)
SLICES        make([]T, len, cap)  |  append  |  copy
MAPS          make(map[K]V)  |  m[k]  |  delete(m, k)
FILE READ     os.ReadFile  |  bufio.Scanner  |  io.ReadAll
FILE WRITE    os.WriteFile  |  os.Create  |  bufio.Writer + Flush
```
