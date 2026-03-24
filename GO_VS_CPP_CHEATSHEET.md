# Go for C++ Developers — Concept Cheatsheet

A side-by-side comparison of Go and C++ concepts for beginners transitioning between the two languages.

---

## Table of Contents

1. [Basic Syntax](#1-basic-syntax)
2. [Variables & Types](#2-variables--types)
3. [Pointers](#3-pointers)
4. [Strings](#4-strings)
5. [Arrays & Slices vs Vectors](#5-arrays--slices-vs-vectors)
6. [Maps vs std::map](#6-maps-vs-stdmap)
7. [Structs & Methods](#7-structs--methods)
8. [Interfaces vs Abstract Classes](#8-interfaces-vs-abstract-classes)
9. [Error Handling vs Exceptions](#9-error-handling-vs-exceptions)
10. [Goroutines vs Threads](#10-goroutines-vs-threads)
11. [Channels vs Queues/Condition Variables](#11-channels-vs-queuescondition-variables)
12. [Memory Management](#12-memory-management)
13. [Packages vs Headers/Namespaces](#13-packages-vs-headersnamespaces)
14. [Control Flow](#14-control-flow)
15. [Closures & Function Types](#15-closures--function-types)
16. [Defer vs RAII/Destructors](#16-defer-vs-raiidestructors)
17. [Key Differences Summary](#17-key-differences-summary)

---

## 1. Basic Syntax

### Hello World

**C++:**
```cpp
#include <iostream>

int main() {
    std::cout << "Hello, World!" << std::endl;
    return 0;
}
```

**Go:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

Key differences:
- No semicolons in Go (compiler inserts them)
- No `return 0` needed — `main` returns void
- `package main` + `func main()` is the entry point
- `import` instead of `#include`

---

## 2. Variables & Types

### Declaration

**C++:**
```cpp
int x = 10;
auto y = 3.14;          // type inferred
const int MAX = 100;
```

**Go:**
```go
var x int = 10
y := 3.14               // short declaration, type inferred
const MAX = 100
```

### Type Comparison

| Go | C++ | Notes |
|---|---|---|
| `int`, `int8`, `int16`, `int32`, `int64` | `int`, `int8_t`, `int16_t`, `int32_t`, `int64_t` | Go's `int` is platform-sized |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `unsigned int`, `uint8_t`, ... | Same idea |
| `float32`, `float64` | `float`, `double` | Go has no `long double` |
| `bool` | `bool` | Same |
| `string` | `std::string` | Go strings are immutable |
| `byte` | `uint8_t` / `char` | Alias for `uint8` |
| `rune` | `char32_t` | Alias for `int32` (Unicode code point) |

### Zero Values

In Go, every variable gets a default **zero value** — no uninitialized garbage:

| Type | Go Zero Value | C++ |
|---|---|---|
| `int` | `0` | undefined (local) |
| `float64` | `0.0` | undefined (local) |
| `bool` | `false` | undefined (local) |
| `string` | `""` | `""` (if `std::string`) |
| pointer | `nil` | `nullptr` |

---

## 3. Pointers

Go has pointers but **no pointer arithmetic** — safer than C++.

**C++:**
```cpp
int x = 42;
int* p = &x;
std::cout << *p;     // dereference: 42
p++;                 // pointer arithmetic — allowed
```

**Go:**
```go
x := 42
p := &x
fmt.Println(*p)      // dereference: 42
// p++               // COMPILE ERROR — no pointer arithmetic
```

### Accessing struct members through pointer

**C++:**
```cpp
struct Point { int x; int y; };
Point p = {1, 2};
Point* ptr = &p;
ptr->x;              // arrow operator for pointers
(*ptr).x;            // or explicit dereference
```

**Go:**
```go
type Point struct { X, Y int }
p := Point{1, 2}
ptr := &p
ptr.X                // dot operator works for BOTH values and pointers
(*ptr).X             // also works, but nobody writes this
```

Go removed `->` entirely — `.` auto-dereferences pointers.

### new vs make

**C++:**
```cpp
int* p = new int(42);        // heap allocation
std::vector<int>* v = new std::vector<int>();
delete p;                     // manual cleanup
delete v;
```

**Go:**
```go
p := new(int)                 // allocates zeroed int, returns *int
*p = 42

s := make([]int, 5)          // allocates + initializes a slice
m := make(map[string]int)    // allocates + initializes a map
ch := make(chan int)          // allocates + initializes a channel
// no delete/free — garbage collector handles it
```

- `new(T)` → allocates zeroed memory, returns `*T`
- `make(T)` → only for slices, maps, channels — returns initialized `T`

---

## 4. Strings

**C++:**
```cpp
std::string s = "hello";
s[0] = 'H';                  // mutable
s += " world";               // concatenation
s.length();                   // length
s.substr(0, 5);               // substring
```

**Go:**
```go
s := "hello"
// s[0] = 'H'                // COMPILE ERROR — strings are immutable
s += " world"                 // creates a new string
len(s)                        // length in bytes
s[0:5]                        // substring via slice syntax
```

Go strings are **immutable**. To mutate, convert to `[]byte`:
```go
b := []byte(s)
b[0] = 'H'
s = string(b)
```

---

## 5. Arrays & Slices vs Vectors

### Fixed-size Arrays

**C++:**
```cpp
int arr[5] = {1, 2, 3, 4, 5};
std::array<int, 5> arr = {1, 2, 3, 4, 5};
```

**Go:**
```go
arr := [5]int{1, 2, 3, 4, 5}  // fixed size, rarely used directly
```

### Slices ≈ std::vector

This is the more common comparison:

**C++:**
```cpp
std::vector<int> v = {1, 2, 3};
v.push_back(4);              // append
v.size();                     // length
v.capacity();                 // capacity
```

**Go:**
```go
s := []int{1, 2, 3}
s = append(s, 4)             // append returns a new slice
len(s)                        // length
cap(s)                        // capacity
```

### Slicing

**C++:**
```cpp
// No built-in slicing, use iterators or spans
std::span<int> sub(v.begin() + 1, v.begin() + 3);
```

**Go:**
```go
s := []int{0, 10, 20, 30, 40}
sub := s[1:3]                // [10, 20] — shares underlying array
```

---

## 6. Maps vs std::map

**C++:**
```cpp
std::unordered_map<std::string, int> m;
m["age"] = 30;
m.count("age");                    // check existence (0 or 1)
auto it = m.find("age");
if (it != m.end()) { ... }
m.erase("age");
```

**Go:**
```go
m := make(map[string]int)
m["age"] = 30
val, ok := m["age"]               // ok is true if key exists
if ok { ... }
delete(m, "age")
```

Go maps are hash maps (like `std::unordered_map`), not tree maps.

---

## 7. Structs & Methods

### Structs

**C++:**
```cpp
class Person {
public:
    std::string name;
    int age;

    void greet() {
        std::cout << "Hi, I'm " << name << std::endl;
    }
};

Person p{"Alice", 30};
p.greet();
```

**Go:**
```go
type Person struct {
    Name string
    Age  int
}

func (p Person) Greet() {
    fmt.Printf("Hi, I'm %s\n", p.Name)
}

p := Person{Name: "Alice", Age: 30}
p.Greet()
```

Key differences:
- Go has **no classes** — only structs
- Methods are defined **outside** the struct with a **receiver** `(p Person)`
- No `public`/`private` keywords — **uppercase = exported, lowercase = unexported**

### Value vs Pointer Receiver

**C++:**
```cpp
class Counter {
    int count = 0;
public:
    void increment() { count++; }        // modifies the object
    int getCount() const { return count; } // doesn't modify
};
```

**Go:**
```go
type Counter struct {
    count int
}

func (c *Counter) Increment() { c.count++ }  // pointer receiver — modifies original
func (c Counter) GetCount() int { return c.count } // value receiver — works on a copy
```

Rule of thumb: use a **pointer receiver** when you need to modify the struct (like non-const methods in C++).

---

## 8. Interfaces vs Abstract Classes

This is a major philosophical difference.

**C++ (explicit — must inherit):**
```cpp
class Writer {
public:
    virtual void Write(const std::string& data) = 0; // pure virtual
    virtual ~Writer() = default;
};

class FileWriter : public Writer {    // explicit inheritance
public:
    void Write(const std::string& data) override { ... }
};
```

**Go (implicit — no "implements" keyword):**
```go
type Writer interface {
    Write(data []byte) (int, error)
}

type FileWriter struct { ... }

func (f *FileWriter) Write(data []byte) (int, error) { ... }
// FileWriter automatically satisfies Writer — no declaration needed
```

In Go, if a type has the right methods, it satisfies the interface **automatically** (duck typing). No need for inheritance, `extends`, or `implements`.

---

## 9. Error Handling vs Exceptions

**C++ (exceptions):**
```cpp
try {
    auto file = openFile("data.txt");    // might throw
    auto data = readAll(file);           // might throw
} catch (const FileError& e) {
    std::cerr << e.what() << std::endl;
} catch (const std::exception& e) {
    std::cerr << e.what() << std::endl;
}
```

**Go (explicit error return):**
```go
file, err := os.Open("data.txt")
if err != nil {
    log.Fatal(err)
}

data, err := io.ReadAll(file)
if err != nil {
    log.Fatal(err)
}
```

Go has **no exceptions** (no try/catch). Errors are values returned by functions. You check them after every call. This is verbose but explicit — you always know where errors can occur.

Go does have `panic`/`recover` for truly unrecoverable situations, but they're not used for normal control flow.

---

## 10. Goroutines vs Threads

**C++ (OS threads — heavy):**
```cpp
#include <thread>

void doWork(int id) {
    std::cout << "Worker " << id << std::endl;
}

int main() {
    std::thread t1(doWork, 1);
    std::thread t2(doWork, 2);
    t1.join();  // wait for thread
    t2.join();
}
```

**Go (goroutines — lightweight):**
```go
func doWork(id int) {
    fmt.Println("Worker", id)
}

func main() {
    go doWork(1)   // just add "go"
    go doWork(2)
    time.Sleep(time.Second)  // crude wait (use sync.WaitGroup in real code)
}
```

| | Goroutine | C++ Thread |
|---|---|---|
| Stack size | ~2KB (grows) | ~1–8MB (fixed) |
| Creation cost | Very cheap | Expensive |
| Typical count | Millions | Thousands |
| Managed by | Go runtime | OS |

---

## 11. Channels vs Queues/Condition Variables

In C++, thread communication is manual with mutexes and condition variables. Go uses **channels**.

**C++ (manual synchronization):**
```cpp
std::mutex mtx;
std::condition_variable cv;
std::queue<std::string> q;

// Producer
{
    std::lock_guard<std::mutex> lock(mtx);
    q.push("data");
}
cv.notify_one();

// Consumer
{
    std::unique_lock<std::mutex> lock(mtx);
    cv.wait(lock, [&]{ return !q.empty(); });
    auto val = q.front();
    q.pop();
}
```

**Go (channels — built-in):**
```go
ch := make(chan string)

// Producer (goroutine)
go func() {
    ch <- "data"      // send
}()

// Consumer
val := <-ch            // receive (blocks until data available)
```

All that mutex/condition_variable boilerplate is built into the channel.

---

## 12. Memory Management

**C++:**
```cpp
int* p = new int(42);    // manual allocation
delete p;                 // manual deallocation
// or use smart pointers:
auto p = std::make_unique<int>(42);  // auto cleanup
auto p = std::make_shared<int>(42);  // reference counted
```

**Go:**
```go
p := new(int)            // heap allocated
*p = 42
// no delete — garbage collector handles it
```

| Feature | C++ | Go |
|---|---|---|
| Memory management | Manual or smart pointers | Garbage collected |
| Stack vs heap | You decide | Compiler decides (escape analysis) |
| Dangling pointers | Possible | Impossible |
| Memory leaks | Possible | Possible (e.g. goroutine leaks) |
| Destructors | Yes (`~Class`) | No (use `defer`) |
| RAII | Yes | No (use `defer`) |

---

## 13. Packages vs Headers/Namespaces

**C++:**
```cpp
// math_utils.h
#ifndef MATH_UTILS_H
#define MATH_UTILS_H
namespace math_utils {
    int add(int a, int b);
}
#endif

// math_utils.cpp
#include "math_utils.h"
namespace math_utils {
    int add(int a, int b) { return a + b; }
}

// main.cpp
#include "math_utils.h"
math_utils::add(1, 2);
```

**Go:**
```go
// mathutils/mathutils.go
package mathutils

func Add(a, b int) int { return a + b }  // uppercase = exported

// main.go
package main
import "myproject/mathutils"

mathutils.Add(1, 2)
```

No header files, no include guards, no forward declarations, no namespaces. Just packages with exported (uppercase) and unexported (lowercase) identifiers.

---

## 14. Control Flow

### If

**C++:**
```cpp
if (x > 0) { ... }
```

**Go:**
```go
if x > 0 { ... }           // no parentheses required
if err := doThing(); err != nil { ... }  // init statement in if
```

### For (Go has only `for` — no `while`)

**C++:**
```cpp
for (int i = 0; i < 10; i++) { ... }
while (true) { ... }
for (auto& v : vec) { ... }
```

**Go:**
```go
for i := 0; i < 10; i++ { ... }
for { ... }                          // infinite loop (replaces while(true))
for _, v := range slice { ... }      // range-based loop
```

### Switch

**C++:**
```cpp
switch (x) {
    case 1: doA(); break;     // must break explicitly
    case 2: doB(); break;
    default: doC();
}
```

**Go:**
```go
switch x {
case 1:
    doA()                     // auto-breaks (no fallthrough by default)
case 2:
    doB()
default:
    doC()
}
```

---

## 15. Closures & Function Types

**C++ (lambdas):**
```cpp
auto add = [](int a, int b) -> int { return a + b; };
add(1, 2);

int x = 10;
auto adder = [x](int n) { return x + n; };  // captures x by value
```

**Go (closures):**
```go
add := func(a, b int) int { return a + b }
add(1, 2)

x := 10
adder := func(n int) int { return x + n }  // captures x by reference
```

Key difference: Go closures capture by **reference** by default, C++ lambdas capture by value or reference depending on the capture clause.

---

## 16. Defer vs RAII/Destructors

**C++ (RAII — destructor runs when object goes out of scope):**
```cpp
{
    std::ifstream file("data.txt");   // opened
    // use file...
}   // destructor closes file automatically
```

**Go (defer — runs when function returns):**
```go
func readFile() {
    file, err := os.Open("data.txt")
    if err != nil { log.Fatal(err) }
    defer file.Close()    // scheduled to run when readFile() returns

    // use file...
}   // file.Close() runs here
```

`defer` is Go's replacement for RAII. It schedules cleanup right next to the resource acquisition, but the cleanup runs at function exit.

Multiple defers run in **LIFO order** (stack):
```go
defer fmt.Println("first")
defer fmt.Println("second")
defer fmt.Println("third")
// Output: third, second, first
```

---

## 17. Key Differences Summary

| Concept | C++ | Go |
|---|---|---|
| Paradigm | OOP + multi-paradigm | Procedural + concurrent |
| Classes | Yes (inheritance) | No (structs + interfaces) |
| Inheritance | Yes (`class B : public A`) | No (composition instead) |
| Polymorphism | Virtual functions | Interfaces (implicit) |
| Templates/Generics | Templates | Generics (since Go 1.18) |
| Exceptions | `try`/`catch`/`throw` | No — return `error` values |
| Memory | Manual / smart pointers | Garbage collected |
| Concurrency | `std::thread` + mutexes | Goroutines + channels |
| Pointer arithmetic | Yes | No |
| Operator overloading | Yes | No |
| Function overloading | Yes | No |
| Default parameters | Yes | No |
| Header files | Yes | No |
| Build system | CMake, Make, etc. | `go build` (built-in) |
| Unused imports/vars | Warning | **Compile error** |
| Semicolons | Required | Auto-inserted |
| Arrow operator `->` | Yes (pointers) | No (`.` works for both) |
