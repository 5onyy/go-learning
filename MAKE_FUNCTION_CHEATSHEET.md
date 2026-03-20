# Go `make` Built-in Function

`make` is a built-in function in Go used to **allocate and initialize** three specific types:

- **slices**
- **maps**
- **channels**

These are the only types that require initialization before use because they have internal data structures (pointers, lengths, capacities, etc.) that must be set up.

---

## Why `make` and not `new`?

- `new(T)` allocates zeroed memory and returns a `*T` (pointer). It does **not** initialize the internal structure.
- `make(T, ...)` returns an **initialized** (not zeroed) value of type `T` (not a pointer). It's required because slices, maps, and channels need internal bookkeeping set up before they can be used.

For example, a `map` declared as `var m map[string]int` is `nil` — writing to it will panic. You need `make` to create a usable map.

---

## Signatures

```go
// Slice: make([]T, length, capacity)
make([]T, len)          // capacity defaults to len
make([]T, len, cap)

// Map: make(map[K]V, initialCapacity)
make(map[K]V)           // no initial capacity hint
make(map[K]V, cap)      // pre-allocate space for ~cap entries

// Channel: make(chan T, bufferSize)
make(chan T)             // unbuffered channel
make(chan T, cap)        // buffered channel with capacity cap
```

---

## 1. Slices

```go
// Creates a slice of 5 ints, all initialized to 0
s := make([]int, 5)
// s = [0, 0, 0, 0, 0], len=5, cap=5

// Creates a slice with length 0 but capacity 10
s := make([]int, 0, 10)
// s = [], len=0, cap=10
// You can append up to 10 elements before a reallocation happens
```

**When to use:** When you know the size upfront and want to avoid repeated allocations from `append`.

---

## 2. Maps

```go
// Creates an empty, ready-to-use map
m := make(map[string]int)
m["age"] = 30  // works fine

// With capacity hint (optimization, not a hard limit)
m := make(map[string]int, 100)
// Pre-allocates space for ~100 entries to reduce rehashing
```

**Key point:** Without `make`, a map is `nil` and writing to it panics:

```go
var m map[string]int   // m is nil
m["key"] = 1           // PANIC: assignment to entry in nil map
```

---

## 3. Channels

```go
// Unbuffered channel — sender blocks until receiver is ready
ch := make(chan string)

// Buffered channel — can hold up to 5 values before blocking
ch := make(chan string, 5)
```

---

## `make` vs Composite Literals

For slices and maps, you can also use **composite literals** as a shorthand:

```go
// These are equivalent for an empty, initialized slice:
s := make([]int, 0)
s := []int{}

// These are equivalent for an empty, initialized map:
m := make(map[string]int)
m := map[string]int{}
```

Use `make` when you need to specify a **length or capacity**. Use literals when you want to initialize with values or just need an empty collection.

---

## Summary Table

| Type    | Syntax                        | Returns         | Zero value (without make)  |
|---------|-------------------------------|-----------------|----------------------------|
| Slice   | `make([]T, len, cap)`         | `[]T`           | `nil` (unusable for indexing) |
| Map     | `make(map[K]V, cap)`          | `map[K]V`       | `nil` (panics on write)    |
| Channel | `make(chan T, cap)`            | `chan T`         | `nil` (blocks forever)     |

The key takeaway: **use `make` whenever you need a slice, map, or channel that is ready to use**, especially when you want to control the initial size or buffer capacity.
