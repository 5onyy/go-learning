# Go Slices — Complete Cheatsheet

> From initialization to the `slices` package — everything about Go's most-used data structure.

---

## Table of Contents

**Part 1 — Slice Fundamentals**

1. [What Is a Slice?](#1-what-is-a-slice)
2. [Creating Slices](#2-creating-slices)
3. [Slicing (Sub-Slices)](#3-slicing-sub-slices)
4. [Built-in: `append](#4-built-in-append)`
5. [Built-in: `copy](#5-built-in-copy)`
6. [Built-in: `clear` (Go 1.21+)](#6-built-in-clear-go-121)
7. [Iterating](#7-iterating)
8. [Passing Slices to Functions](#8-passing-slices-to-functions)
9. [Multi-Dimensional Slices](#9-multi-dimensional-slices)
10. [Nil vs Empty Slice](#10-nil-vs-empty-slice)
11. [Common Manual Operations](#11-common-manual-operations)

**Part 2 — `slices` Package (Go 1.21+)**

1. [Searching](#12-searching)
2. [Sorting](#13-sorting)
3. [Comparing](#14-comparing)
4. [Inserting & Deleting](#15-inserting--deleting)
5. [Growing & Clipping](#16-growing--clipping)
6. [Cloning & Compacting](#17-cloning--compacting)
7. [Min / Max](#18-min--max)
8. [Replacing & Concatenating](#19-replacing--concatenating)
9. [Reversing](#20-reversing)
10. [Iterators (Go 1.23+)](#21-iterators-go-123)
11. [Common Patterns & Recipes](#22-common-patterns--recipes)

---

# Part 1 — Slice Fundamentals

## 1. What Is a Slice?

A slice is a **dynamically-sized view** into a contiguous block of memory (an underlying array). Under the hood, every slice is a 3-field struct:

```
┌──────────┬─────┬─────┐
│  pointer │ len │ cap │
└────┬─────┴─────┴─────┘
     │
     ▼
   backing array: [e0, e1, e2, e3, e4, _, _, _]
                   ◄──── len ────►◄─ unused ─►
                   ◄──────────── cap ──────────►
```

- **pointer** — address of the first element the slice can see
- **len** — number of elements the slice currently holds
- **cap** — number of elements from the pointer to the end of the backing array

```go
s := make([]int, 3, 5)
fmt.Println(len(s))   // 3
fmt.Println(cap(s))   // 5
```

Key insight: **a slice does not own data** — it's a window into an array. Multiple slices can share the same backing array.

---

## 2. Creating Slices

### Literal

```go
s := []int{1, 2, 3}                  // len=3, cap=3
names := []string{"alice", "bob"}     // len=2, cap=2
empty := []int{}                      // len=0, cap=0 (NOT nil)
```

### `make` — Control Length and Capacity

```go
s1 := make([]int, 5)                 // len=5, cap=5, all zeros: [0 0 0 0 0]
s2 := make([]int, 0, 10)             // len=0, cap=10, empty but pre-allocated
s3 := make([]int, 3, 10)             // len=3, cap=10: [0 0 0]
```

Use `make([]T, 0, n)` when you know the final size but want to `append` into it.
Use `make([]T, n)` when you want to index directly (`s[i] = val`).

### From an Array

```go
arr := [5]int{10, 20, 30, 40, 50}

s := arr[:]                           // full slice, shares arr's memory
s2 := arr[1:4]                        // [20 30 40], len=3, cap=4
```

### `var` Declaration (nil slice)

```go
var s []int                           // nil, len=0, cap=0
s = append(s, 1)                      // works fine — append handles nil
```

### From Another Slice

```go
original := []int{0, 1, 2, 3, 4}
sub := original[1:3]                  // [1 2] — shares backing array!
```

---

## 3. Slicing (Sub-Slices)

### Syntax: `s[low:high]`

```go
s := []int{0, 1, 2, 3, 4}

s[1:3]    // [1 2]       — indices 1, 2 (high is exclusive)
s[:3]     // [0 1 2]     — from start
s[2:]     // [2 3 4]     — to end
s[:]      // [0 1 2 3 4] — full slice
```

### Memory Is Shared!

```go
a := []int{0, 1, 2, 3, 4}
b := a[1:3]                 // b = [1, 2]

b[0] = 99                   // BOTH slices see this change
fmt.Println(a)              // [0 99 2 3 4]
fmt.Println(b)              // [99 2]
```

### Capacity of a Sub-Slice

```go
s := []int{0, 1, 2, 3, 4}   // len=5, cap=5
b := s[1:3]                  // len=2, cap=4 (from index 1 to end of backing array)
```

The capacity extends from the start of the sub-slice to the end of the original backing array:

```
s:  [0, 1, 2, 3, 4]
         ↑        ↑
b start  b[0:2]   cap boundary
         ◄─ len ─►
         ◄──── cap ────►
```

### Full Slice Expression: `s[low:high:max]` — Limit Capacity

```go
a := []int{0, 1, 2, 3, 4}
b := a[1:3:4]                // b = [1, 2], len=2, cap=3 (capped at index 4)

// WHY: prevents b from overwriting a[4] when you append to b
b = append(b, 99)            // len=3, cap=3 — fills remaining capacity
b = append(b, 100)           // cap exceeded → NEW backing array allocated (safe!)
```

This is a defensive pattern to prevent accidental aliasing bugs.

---

## 4. Built-in: `append`

### Basic Usage

```go
s := []int{1, 2, 3}

s = append(s, 4)               // [1 2 3 4]
s = append(s, 5, 6, 7)         // [1 2 3 4 5 6 7]  — multiple values
s = append(s, []int{8, 9}...)  // [1 2 3 4 5 6 7 8 9]  — spread a slice
```

**Golden rule: always reassign** `s = append(s, ...)` — `append` may return a new slice header.

### How `append` Works Internally

```
Has capacity?
├── YES → write into existing array, return header with len+1
└── NO  → allocate new (larger) array, copy old data, write new element,
           return header with new pointer, new len, new cap
```

```go
s := make([]int, 3, 5)        // len=3, cap=5

s = append(s, 4)              // fits → same array, len=4, cap=5
s = append(s, 5)              // fits → same array, len=5, cap=5
s = append(s, 6)              // DOESN'T fit → new array, len=6, cap=10 (approx 2x growth)
```

### The Aliasing Trap

```go
a := make([]int, 3, 5)        // [0 0 0], cap=5
a[0], a[1], a[2] = 1, 2, 3

b := a[:]                     // b shares a's backing array
b = append(b, 4)              // enough cap → writes to SHARED array

fmt.Println(a[:cap(a)])        // [1 2 3 4 _]  — a's backing array was modified!
```

Fix: use a full slice expression `b := a[:len(a):len(a)]` to force a new allocation on append.

### Why You Must Reassign

```go
func myAppend(s []int, v int) []int {
    return append(s, v)
}

s := make([]int, 3, 10)       // len=3, cap=10
s[0], s[1], s[2] = 1, 2, 3

fmt.Println(myAppend(s, 4))   // [1 2 3 4]  — returned slice has len=4
fmt.Println(s)                 // [1 2 3]    — s still has len=3!
fmt.Println(s[:4])             // [1 2 3 4]  — the 4 IS in the array, s just can't see it

s = myAppend(s, 4)             // ✓ now s has len=4
fmt.Println(s)                 // [1 2 3 4]
```

The caller's slice header (pointer, len, cap) is passed **by value**. `append` inside the function returns a new header with updated `len`, but the caller's `s` is unchanged unless you reassign.

---

## 5. Built-in: `copy`

```go
src := []int{1, 2, 3, 4, 5}
dst := make([]int, 3)

n := copy(dst, src)            // n=3 (copies min(len(dst), len(src)))
fmt.Println(dst)               // [1 2 3]
```

### Copy Into Middle of Slice

```go
s := make([]int, 5)
copy(s[2:], []int{7, 8, 9})   // s = [0 0 7 8 9]
```

### Copy Between Overlapping Slices (Safe)

```go
s := []int{0, 1, 2, 3, 4}
copy(s[1:], s[0:4])           // shift right: [0 0 1 2 3]
```

### Clone a Slice (Independent Copy)

```go
original := []int{1, 2, 3}

// Method 1: make + copy
clone := make([]int, len(original))
copy(clone, original)

// Method 2: append to nil
clone2 := append([]int(nil), original...)

// Method 3: slices.Clone (Go 1.21+)
clone3 := slices.Clone(original)
```

---

## 6. Built-in: `clear` (Go 1.21+)

Sets all elements to their zero value without changing length or capacity.

```go
s := []int{1, 2, 3, 4, 5}
clear(s)
fmt.Println(s)                 // [0 0 0 0 0], len=5, cap unchanged
```

---

## 7. Iterating

### `for` with Index

```go
s := []string{"a", "b", "c"}

for i := 0; i < len(s); i++ {
    fmt.Println(i, s[i])
}
```

### `range` — Index + Value

```go
for i, v := range s {
    fmt.Println(i, v)          // 0 a, 1 b, 2 c
}
```

### `range` — Index Only

```go
for i := range s {
    fmt.Println(i)             // 0, 1, 2
}
```

### `range` — Value Only

```go
for _, v := range s {
    fmt.Println(v)             // a, b, c
}
```

### Range Copies Values

```go
type Point struct{ X, Y int }
points := []Point{{1, 2}, {3, 4}}

for _, p := range points {
    p.X = 99                   // modifies the COPY, not the slice element
}
fmt.Println(points)            // [{1 2} {3 4}]  — unchanged!

// Fix: use index
for i := range points {
    points[i].X = 99
}
```

---

## 8. Passing Slices to Functions

### Slices Are Passed by Value — But the Header Contains a Pointer

```go
func double(s []int) {
    for i := range s {
        s[i] *= 2              // modifies the SHARED backing array
    }
}

nums := []int{1, 2, 3}
double(nums)
fmt.Println(nums)              // [2 4 6]  — modified!
```

The slice header is copied, but both copies point to the **same backing array**. Mutations through indexing are visible to the caller.

### But `append` Inside a Function Won't Update Caller's `len`

```go
func addElement(s []int, v int) {
    s = append(s, v)           // updates local copy's len, not caller's
}

nums := []int{1, 2, 3}
addElement(nums, 4)
fmt.Println(nums)              // [1 2 3]  — unchanged!
```

Fix: **return the new slice** or use a **pointer to the slice**:

```go
// Option 1: return
func addElement(s []int, v int) []int {
    return append(s, v)
}
nums = addElement(nums, 4)

// Option 2: pointer to slice (less idiomatic)
func addElement(s *[]int, v int) {
    *s = append(*s, v)
}
```

### Summary: What Can a Function Do to Your Slice?


| Operation                              | Visible to caller?            |
| -------------------------------------- | ----------------------------- |
| `s[i] = val` (modify existing element) | Yes — shared backing array    |
| `append(s, val)` (without return)      | No — caller's `len` unchanged |
| `s = append(s, val)` + return          | Yes — if caller reassigns     |


---

## 9. Multi-Dimensional Slices

### 2D Slice (Slice of Slices)

```go
rows, cols := 3, 4

grid := make([][]int, rows)
for i := range grid {
    grid[i] = make([]int, cols)        // each row is independently allocated
}

grid[1][2] = 42
```

### Literal 2D

```go
matrix := [][]int{
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9},
}
```

### Jagged (Rows of Different Lengths)

```go
triangle := [][]int{
    {1},
    {1, 1},
    {1, 2, 1},
    {1, 3, 3, 1},
}
```

### Single Allocation (Performance)

```go
rows, cols := 3, 4
data := make([]int, rows*cols)         // one contiguous block
grid := make([][]int, rows)
for i := range grid {
    grid[i] = data[i*cols : (i+1)*cols]
}
```

---

## 10. Nil vs Empty Slice

```go
var s1 []int                  // nil — s1 == nil → true
s2 := []int{}                 // empty — s2 == nil → false
s3 := make([]int, 0)          // empty — s3 == nil → false
```


|                  | nil          | empty        |
| ---------------- | ------------ | ------------ |
| `len(s)`         | 0            | 0            |
| `cap(s)`         | 0            | 0            |
| `s == nil`       | `true`       | `false`      |
| `append(s, 1)`   | works        | works        |
| `range s`        | 0 iterations | 0 iterations |
| JSON marshal     | `null`       | `[]`         |
| `fmt.Println(s)` | `[]`         | `[]`         |


Use `nil` for "no value / unset". Use `[]T{}` when you need an empty but non-null collection (e.g., JSON APIs).

---

## 11. Common Manual Operations

### Delete at Index (Order Preserved)

```go
s := []int{0, 1, 2, 3, 4}
i := 2
s = append(s[:i], s[i+1:]...)         // [0 1 3 4]
```

### Delete at Index (Order NOT Preserved — Fast)

```go
s := []int{0, 1, 2, 3, 4}
i := 2
s[i] = s[len(s)-1]                   // swap with last
s = s[:len(s)-1]                      // shrink
// [0 1 4 3]
```

### Insert at Index

```go
s := []int{0, 1, 3, 4}
i := 2
s = append(s[:i+1], s[i:]...)        // make room
s[i] = 2                             // insert value
// [0 1 2 3 4]

// Or more readable (Go 1.21+):
s = slices.Insert(s, i, 2)
```

### Stack (LIFO)

```go
stack := []int{}
stack = append(stack, 1)              // push
stack = append(stack, 2)              // push
top := stack[len(stack)-1]            // peek → 2
stack = stack[:len(stack)-1]          // pop
```

### Queue (FIFO)

```go
queue := []int{}
queue = append(queue, 1)              // enqueue
queue = append(queue, 2)              // enqueue
front := queue[0]                     // peek → 1
queue = queue[1:]                     // dequeue (leaks memory over time; see container/ring)
```

### Reverse

```go
s := []int{1, 2, 3, 4, 5}
for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
    s[i], s[j] = s[j], s[i]
}
// [5 4 3 2 1]

// Or (Go 1.21+):
slices.Reverse(s)
```

### Filter In-Place (No Allocation)

```go
s := []int{1, 2, 3, 4, 5, 6}
n := 0
for _, v := range s {
    if v%2 == 0 {
        s[n] = v
        n++
    }
}
s = s[:n]                              // [2 4 6]
```

---

# Part 2 — `slices` Package (Go 1.21+)

```go
import "slices"
```

Generic functions for searching, sorting, comparing, and transforming slices. Prefer this over the older `sort` package.

---

## 12. Searching

### Contains / ContainsFunc

```go
nums := []int{1, 2, 3, 4, 5}

slices.Contains(nums, 3)               // true
slices.Contains(nums, 9)               // false

slices.ContainsFunc(nums, func(n int) bool {
    return n > 4
})  // true
```

### Index / IndexFunc

```go
names := []string{"alice", "bob", "carol"}

slices.Index(names, "bob")             // 1
slices.Index(names, "dave")            // -1

slices.IndexFunc(names, func(s string) bool {
    return len(s) > 4
})  // 0  ("alice" has length 5)
```

### BinarySearch / BinarySearchFunc

Works on **sorted** slices. Returns the index and whether the value was found.

```go
sorted := []int{1, 3, 5, 7, 9}

idx, found := slices.BinarySearch(sorted, 5)    // idx=2, found=true
idx, found = slices.BinarySearch(sorted, 4)     // idx=2, found=false (insertion point)

// Custom comparison
type User struct{ Name string; Age int }
users := []User{{"Alice", 25}, {"Bob", 30}, {"Carol", 35}}

idx, found = slices.BinarySearchFunc(users, User{"Bob", 30}, func(a, b User) int {
    return cmp.Compare(a.Name, b.Name)
})
// idx=1, found=true
```

---

## 13. Sorting

### Sort / SortFunc / SortStableFunc

```go
nums := []int{5, 3, 1, 4, 2}
slices.Sort(nums)                      // [1 2 3 4 5]

// Custom comparison
type Task struct{ Name string; Priority int }
tasks := []Task{{"C", 3}, {"A", 1}, {"B", 2}}

slices.SortFunc(tasks, func(a, b Task) int {
    return cmp.Compare(a.Priority, b.Priority)
})
// [{A 1} {B 2} {C 3}]

// Stable sort — preserves order of equal elements
slices.SortStableFunc(tasks, func(a, b Task) int {
    return cmp.Compare(a.Priority, b.Priority)
})
```

### IsSorted / IsSortedFunc

```go
slices.IsSorted([]int{1, 2, 3, 4})     // true
slices.IsSorted([]int{1, 3, 2, 4})     // false
```

### Multi-Field Sorting

```go
import "cmp"

slices.SortFunc(users, func(a, b User) int {
    if c := cmp.Compare(a.Age, b.Age); c != 0 {
        return c
    }
    return cmp.Compare(a.Name, b.Name)   // tie-break by name
})
```

---

## 14. Comparing

### Equal / EqualFunc

```go
slices.Equal([]int{1, 2, 3}, []int{1, 2, 3})    // true
slices.Equal([]int{1, 2, 3}, []int{1, 2, 4})    // false
slices.Equal([]int{1, 2}, []int{1, 2, 3})        // false (different length)
slices.Equal[[]int](nil, []int{})                 // true  (nil and empty are equal)
```

### Compare / CompareFunc

Lexicographic comparison.

```go
slices.Compare([]int{1, 2, 3}, []int{1, 2, 4})   // -1  (3 < 4)
slices.Compare([]int{1, 2, 3}, []int{1, 2, 3})   //  0
slices.Compare([]int{1, 2}, []int{1, 2, 3})       // -1  (shorter < longer)
```

---

## 15. Inserting & Deleting

### Insert

```go
s := []int{1, 2, 5, 6}
s = slices.Insert(s, 2, 3, 4)          // [1 2 3 4 5 6]
```

### Delete / DeleteFunc

```go
s := []int{1, 2, 3, 4, 5}
s = slices.Delete(s, 1, 3)             // [1 4 5]  — delete indices [1, 3)

s = []int{1, 2, 3, 4, 5, 6}
s = slices.DeleteFunc(s, func(n int) bool {
    return n%2 == 0
})
// [1 3 5]
```

### Remove Single Element

```go
s := []int{10, 20, 30, 40}
s = slices.Delete(s, 2, 3)             // [10 20 40]
```

---

## 16. Growing & Clipping

### Grow — Pre-allocate Extra Capacity

```go
s := []int{1, 2, 3}                    // len=3, cap=3
s = slices.Grow(s, 100)                // len=3, cap≥103
```

### Clip — Release Unused Capacity

```go
s := make([]int, 3, 1000)              // len=3, cap=1000
s = slices.Clip(s)                     // len=3, cap=3
// Equivalent to s[:len(s):len(s)]
```

---

## 17. Cloning & Compacting

### Clone — Shallow Copy

```go
original := []int{1, 2, 3}
cloned := slices.Clone(original)

cloned[0] = 99
fmt.Println(original[0])               // 1 — unaffected
```

### Compact / CompactFunc — Remove Consecutive Duplicates

Works like Unix `uniq`. Sort first for full dedup.

```go
s := []int{1, 1, 2, 2, 2, 3, 3, 1}
s = slices.Compact(s)                   // [1 2 3 1] — only adjacent dups removed

// Full dedup: sort + compact
s = []int{3, 1, 2, 1, 3, 2}
slices.Sort(s)                          // [1 1 2 2 3 3]
s = slices.Compact(s)                   // [1 2 3]
```

---

## 18. Min / Max

```go
nums := []int{5, 3, 8, 1, 9, 2}

slices.Min(nums)                        // 1
slices.Max(nums)                        // 9

// Custom comparison
type Product struct{ Name string; Price float64 }
products := []Product{{"A", 9.99}, {"B", 4.99}, {"C", 14.99}}

cheapest := slices.MinFunc(products, func(a, b Product) int {
    return cmp.Compare(a.Price, b.Price)
})
// {B 4.99}

// Panics on empty slice!
```

---

## 19. Replacing & Concatenating

### Replace — Replace a Range

```go
s := []int{1, 2, 3, 4, 5}
s = slices.Replace(s, 1, 3, 20, 30, 40)  // [1 20 30 40 4 5]
```

### Concat (Go 1.22+)

```go
result := slices.Concat([]int{1, 2}, []int{3, 4}, []int{5, 6})
// [1 2 3 4 5 6]  — always returns new slice
```

---

## 20. Reversing

```go
s := []int{1, 2, 3, 4, 5}
slices.Reverse(s)                       // [5 4 3 2 1]  — in-place

// Reverse a sub-slice
s = []int{1, 2, 3, 4, 5}
slices.Reverse(s[1:4])                  // [1 4 3 2 5]
```

---

## 21. Iterators (Go 1.23+)

### All / Values / Backward

```go
s := []string{"a", "b", "c"}

for i, v := range slices.All(s) {       // index + value
    fmt.Println(i, v)
}

for v := range slices.Values(s) {       // values only
    fmt.Println(v)
}

for i, v := range slices.Backward(s) {  // reverse order
    fmt.Println(i, v)                   // 2 c, 1 b, 0 a
}
```

### Collect / Sorted / Chunk / Repeat

```go
// Collect — materialize iterator into slice
seq := slices.Values([]int{1, 2, 3})
s := slices.Collect(seq)               // []int{1, 2, 3}

// Sorted — collect and sort an iterator
sorted := slices.Sorted(slices.Values([]int{3, 1, 2}))  // [1 2 3]

// Chunk — split into groups
for chunk := range slices.Chunk([]int{1, 2, 3, 4, 5}, 2) {
    fmt.Println(chunk)
}
// [1 2]
// [3 4]
// [5]

// Repeat — repeat a slice N times
r := slices.Repeat([]int{1, 2}, 3)     // [1 2 1 2 1 2]
```

---

## 22. Common Patterns & Recipes

### Filter

```go
func filter[T any](s []T, keep func(T) bool) []T {
    result := make([]T, 0, len(s))
    for _, v := range s {
        if keep(v) {
            result = append(result, v)
        }
    }
    return result
}

// Or with DeleteFunc (removes non-matching, in-place):
evens := slices.Clone(nums)
evens = slices.DeleteFunc(evens, func(n int) bool { return n%2 != 0 })
```

### Map / Transform

```go
func mapSlice[T, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s {
        result[i] = f(v)
    }
    return result
}

strs := mapSlice([]int{1, 2, 3}, strconv.Itoa)  // ["1" "2" "3"]
```

### Unique (Deduplicate)

```go
func unique[T cmp.Ordered](s []T) []T {
    s = slices.Clone(s)
    slices.Sort(s)
    return slices.Compact(s)
}

unique([]int{3, 1, 2, 1, 3})  // [1 2 3]
```

### Flatten (2D → 1D)

```go
func flatten[T any](ss [][]T) []T {
    return slices.Concat(ss...)
}

flatten([][]int{{1, 2}, {3}, {4, 5}})  // [1 2 3 4 5]
```

### Rotate Left by N

```go
func rotateLeft[T any](s []T, n int) {
    n = n % len(s)
    slices.Reverse(s[:n])
    slices.Reverse(s[n:])
    slices.Reverse(s)
}

s := []int{0, 1, 2, 3, 4, 5}
rotateLeft(s, 2)                       // [2 3 4 5 0 1]
```

### Partition

```go
func partition[T any](s []T, pred func(T) bool) (yes, no []T) {
    for _, v := range s {
        if pred(v) {
            yes = append(yes, v)
        } else {
            no = append(no, v)
        }
    }
    return
}

evens, odds := partition([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 })
// evens=[2 4], odds=[1 3 5]
```

---

## Quick Reference Card

```
CREATION      []T{}  make([]T, len, cap)  arr[:]  var s []T (nil)
BUILT-INS     append  copy  clear  len  cap  make
SLICING       s[low:high]  s[low:high:max]  s[:]
ITERATE       for i, v := range s  |  for i := range s  |  for _, v := range s

slices PACKAGE (Go 1.21+):
SEARCH        Contains  ContainsFunc  Index  IndexFunc  BinarySearch  BinarySearchFunc
SORT          Sort  SortFunc  SortStableFunc  IsSorted  IsSortedFunc
COMPARE       Equal  EqualFunc  Compare  CompareFunc
INSERT/DELETE Insert  Delete  DeleteFunc  Replace
GROW/CLIP     Grow  Clip
COPY          Clone  Concat (1.22+)  Repeat (1.23+)
COMPACT       Compact  CompactFunc
MIN/MAX       Min  MinFunc  Max  MaxFunc
REVERSE       Reverse
ITERATORS     All  Values  Backward  Chunk  Collect  Sorted  SortedFunc  AppendSeq (1.23+)
```

## `slices` vs Old `sort` Package


| Task           | Old (`sort`)                             | New (`slices`)                             |
| -------------- | ---------------------------------------- | ------------------------------------------ |
| Sort ints      | `sort.Ints(s)`                           | `slices.Sort(s)`                           |
| Sort with func | `sort.Slice(s, func(i,j int) bool{...})` | `slices.SortFunc(s, func(a,b T) int{...})` |
| Binary search  | `sort.SearchInts(s, v)`                  | `slices.BinarySearch(s, v)`                |
| Is sorted?     | `sort.IntsAreSorted(s)`                  | `slices.IsSorted(s)`                       |
| Stable sort    | `sort.SliceStable(...)`                  | `slices.SortStableFunc(...)`               |


**Prefer `slices`**: generic (type-safe), faster (no interface boxing), comparison func receives elements not indices.

## Performance Tips


| Task               | Approach                                            | Why                                  |
| ------------------ | --------------------------------------------------- | ------------------------------------ |
| Known final size   | `make([]T, 0, n)` or `slices.Grow`                  | Avoids repeated reallocation         |
| Prevent aliasing   | `a[low:high:max]` (full slice expr)                 | Limits cap so append can't overwrite |
| Release memory     | `slices.Clip(s)`                                    | Shrinks cap to len                   |
| Deep copy          | `slices.Clone(s)`                                   | New backing array                    |
| Bulk delete        | `slices.DeleteFunc`                                 | Single pass, O(n)                    |
| Check membership   | `slices.Contains` for small, `map[T]bool` for large | O(n) vs O(1) lookup                  |
| Sorted search      | `slices.BinarySearch`                               | O(log n) vs O(n) for `Contains`      |
| Build large string | `strings.Builder`, not `+=` in loop                 | O(n) vs O(n²)                        |
| 2D perf-critical   | Single `make` + manual row slicing                  | Cache-friendly, one allocation       |


