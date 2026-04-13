# Go `flag` Package Cheatsheet

## What Is `flag`?

The `flag` package lets you define **command-line flags** so users can pass options when running your program:

```bash
./myapp -name "Alice" -age 30 -verbose
```

---

## Basic Usage

### 1. Define Flags → 2. Parse → 3. Use Values

```go
package main

import (
    "flag"
    "fmt"
)

func main() {
    // Step 1: Define flags (returns pointers)
    name := flag.String("name", "World", "who to greet")
    age := flag.Int("age", 0, "your age")
    verbose := flag.Bool("verbose", false, "enable verbose output")

    // Step 2: Parse command-line arguments (MUST call before using flags)
    flag.Parse()

    // Step 3: Use values (dereference pointers with *)
    fmt.Printf("Hello, %s!\n", *name)

    if *verbose {
        fmt.Printf("Age: %d\n", *age)
    }
}
```

```bash
$ go run main.go
Hello, World!

$ go run main.go -name Alice -age 30 -verbose
Hello, Alice!
Age: 30
```

---

## Two Styles of Defining Flags

### Style 1: `flag.Type()` — Returns a Pointer

```go
name := flag.String("name", "default", "description")
// name is *string, use *name to read the value
```

### Style 2: `flag.TypeVar()` — Binds to an Existing Variable

```go
var name string
flag.StringVar(&name, "name", "default", "description")
// name is a plain string, use it directly
```

Both styles work the same way. Use whichever feels clearer.

---

## Available Built-in Flag Types

| Function | Type | Example |
|---|---|---|
| `flag.String` | `*string` | `-msg "hello"` |
| `flag.Int` | `*int` | `-count 5` |
| `flag.Int64` | `*int64` | `-size 1024` |
| `flag.Float64` | `*float64` | `-rate 3.14` |
| `flag.Bool` | `*bool` | `-verbose` (true if present) |
| `flag.Duration` | `*time.Duration` | `-timeout 5s` |

### `flag.Duration` — A Handy Built-in

Parses human-readable durations like `300ms`, `1.5h`, `2m30s`:

```go
period := flag.Duration("period", 1*time.Second, "sleep period")
flag.Parse()
time.Sleep(*period)
```

```bash
$ ./sleep -period 50ms
$ ./sleep -period 2m30s
$ ./sleep -period 1.5h
```

---

## The Three Parts of Every Flag

```go
flag.String("name", "default", "description")
//          ^^^^^^   ^^^^^^^^^   ^^^^^^^^^^^^^
//          flag name  default    help text shown with -help
```

- **Name**: what the user types after `-` (e.g., `-name`)
- **Default**: value used if the user doesn't provide the flag
- **Usage**: text shown when user runs with `-help`

---

## Automatic `-help` Flag

Every program using `flag` gets `-help` for free:

```bash
$ ./myapp -help
Usage of ./myapp:
  -age int
        your age
  -name string
        who to greet (default "World")
  -verbose
        enable verbose output
```

---

## Positional Arguments (Non-Flag Args)

After `flag.Parse()`, use `flag.Args()` to get remaining non-flag arguments:

```go
flag.Parse()
fmt.Println(flag.Args()) // arguments after all flags
```

```bash
$ ./myapp -verbose file1.txt file2.txt
# flag.Args() returns ["file1.txt", "file2.txt"]
```

---

## Custom Flag Types with `flag.Value`

For types beyond the built-ins, implement the `flag.Value` interface:

```go
type Value interface {
    String() string    // format the flag value for display
    Set(string) error  // parse the string and update the value
}
```

### Example: A Temperature Flag

```go
type Celsius float64

type celsiusFlag struct{ Celsius }

// String formats the value for -help output
func (f *celsiusFlag) String() string {
    return fmt.Sprintf("%g°C", f.Celsius)
}

// Set parses the user's input like "100C" or "212F"
func (f *celsiusFlag) Set(s string) error {
    var unit string
    var value float64
    fmt.Sscanf(s, "%f%s", &value, &unit)
    switch unit {
    case "C", "°C":
        f.Celsius = Celsius(value)
        return nil
    case "F", "°F":
        f.Celsius = Celsius((value - 32) * 5 / 9)
        return nil
    }
    return fmt.Errorf("invalid temperature %q", s)
}
```

Register it with `flag.Var`:

```go
func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
    f := celsiusFlag{value}
    flag.CommandLine.Var(&f, name, usage)
    return &f.Celsius
}
```

Use it:

```go
temp := CelsiusFlag("temp", 20.0, "the temperature")
flag.Parse()
fmt.Println(*temp)
```

```bash
$ ./tempflag
20°C
$ ./tempflag -temp 212F
100°C
$ ./tempflag -temp -18C
-18°C
```

---

## Common Mistakes

| Mistake | Fix |
|---|---|
| Forgetting `flag.Parse()` | Always call it before reading any flag value |
| Reading `name` instead of `*name` | `flag.String` returns a pointer — dereference with `*` |
| Defining flags inside `main` after `Parse` | Define all flags **before** calling `flag.Parse()` |
| Using `os.Args` directly | Let `flag` handle parsing; use `flag.Args()` for positional args |

---

## Quick Reference

```go
// Define
name := flag.String("name", "default", "help text")

// Or bind to existing variable
var name string
flag.StringVar(&name, "name", "default", "help text")

// Parse (call once, before using any flag)
flag.Parse()

// Use
fmt.Println(*name)       // pointer style
fmt.Println(name)        // variable style

// Remaining non-flag arguments
args := flag.Args()

// Custom types
flag.Var(&myCustomFlag, "flagname", "help text")
```
