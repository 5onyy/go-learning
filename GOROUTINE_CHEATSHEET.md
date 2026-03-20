# Goroutines in Go

A **goroutine** is a lightweight, concurrently executing function managed by the Go runtime. You launch one with the `go` keyword:

```go
go someFunction()
```

That's it — `someFunction` now runs **concurrently** alongside everything else.

---

## Goroutines vs OS Threads

Goroutines are often compared to threads, but they are fundamentally different:

| | Goroutine | OS Thread |
|---|---|---|
| Stack size | ~2KB (grows dynamically) | ~1–8MB (fixed) |
| Creation cost | Very cheap | Expensive |
| Managed by | Go runtime | Operating system |
| Scheduling | Cooperative + preemptive (Go scheduler) | OS scheduler |
| Typical count | Millions possible | Thousands max |

Go multiplexes many goroutines onto a smaller number of OS threads automatically — this is the **M:N scheduler**.

---

## Basic Example

```go
func sayHello() {
    fmt.Println("Hello!")
}

func main() {
    go sayHello()         // runs concurrently
    fmt.Println("World")  // main continues immediately
    time.Sleep(time.Second) // wait, otherwise main exits before goroutine runs
}
```

Without the `Sleep`, the program may exit before `sayHello` ever runs — because when `main` returns, **all goroutines are killed**.

---

## Anonymous Goroutines

Very common pattern — launching an inline function:

```go
go func() {
    fmt.Println("I'm a goroutine")
}()  // note the () to call it immediately
```

---

## Synchronization

Goroutines run independently, so you need a way to coordinate them. Go provides two main tools:

### 1. Channels (communicate)

```go
ch := make(chan string)

go func() {
    ch <- "done"  // send result
}()

msg := <-ch  // main blocks here until goroutine sends
fmt.Println(msg)
```

### 2. `sync.WaitGroup` (just wait)

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        fmt.Println(n)
    }(i)
}

wg.Wait()  // blocks until all 5 goroutines call Done()
```

---

## Preventing Data Races

**Data races** happen when two goroutines access the same memory simultaneously and at least one is writing. Go gives you several tools to prevent this:

### Channels (the Go-idiomatic way)

Only one goroutine owns the data at a time — ownership is transferred through the channel:

```go
ch := make(chan string)

go func() {
    result := doWork()  // goroutine owns result
    ch <- result        // transfers ownership to receiver
}()

msg := <-ch  // main now owns it, goroutine is done with it
```

### `sync.Mutex` — mutual exclusion lock

When goroutines **must** share a variable, use a mutex to ensure only one can access it at a time:

```go
var mu sync.Mutex
counter := 0

for i := 0; i < 100; i++ {
    go func() {
        mu.Lock()
        counter++    // only one goroutine at a time
        mu.Unlock()
    }()
}
```

Without the mutex, `counter++` is a **data race** — read, increment, write are three separate operations that can interleave.

### `sync/atomic` — for simple counters/flags

```go
var counter int64

go func() {
    atomic.AddInt64(&counter, 1)  // safe increment
}()
```

---

## The Classic Trap: Closure Variable Capture

```go
// WRONG — all goroutines share the same `i`
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i)  // captures i by reference — data race!
    }()
}
// likely prints: 5 5 5 5 5
```

```go
// CORRECT — pass i as an argument (each goroutine gets its own copy)
for i := 0; i < 5; i++ {
    go func(n int) {
        fmt.Println(n)  // n is a local copy
    }(i)
}
// prints: 0 1 2 3 4 (in some order)
```

---

## Detecting Data Races

Go has a built-in **race detector**:

```bash
go run -race main.go
go test -race ./...
```

It will report exactly which goroutines are racing and on which line.

---

## Real-world Example: Concurrent HTTP Fetcher

```go
func main() {
    start := time.Now()
    ch := make(chan string)

    for _, url := range os.Args[1:] {
        go fetch(url, ch)   // launch one goroutine per URL
    }
    for range os.Args[1:] {
        fmt.Println(<-ch)   // collect one result per goroutine
    }

    fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
    start := time.Now()
    resp, err := http.Get(url)
    if err != nil {
        ch <- fmt.Sprintf("while getting %s: %v", url, err)
        return
    }
    defer resp.Body.Close()
    nbytes, err := io.Copy(io.Discard, resp.Body)
    if err != nil {
        ch <- fmt.Sprintf("while reading %s: %v", url, err)
        return
    }
    secs := time.Since(start).Seconds()
    ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
```

```
main goroutine          goroutine 1 (url1)       goroutine 2 (url2)
─────────────           ──────────────────        ──────────────────
launch goroutine 1  →   http.Get(url1)...
launch goroutine 2  →                             http.Get(url2)...
wait on <-ch        ←   ch <- "0.23s  5432 url1"
wait on <-ch        ←                             ch <- "0.51s 12300 url2"
print total time
```

Without goroutines, total time = sum of all fetches. With goroutines, total time ≈ the **slowest** single fetch.

---

## Synchronization Tool Summary

| Tool | Use when |
|---|---|
| **Channels** | Passing data between goroutines (preferred) |
| **`sync.Mutex`** | Multiple goroutines share a variable |
| **`sync.WaitGroup`** | Waiting for goroutines to finish |
| **`sync/atomic`** | Simple counters or flags |
| **Local variables** | Each goroutine works on its own data (best — no sharing at all) |

---

## Key Rules

1. `go f()` — launches `f` as a goroutine; returns immediately
2. Goroutines are **not** OS threads — they're much cheaper
3. When `main` returns, **all goroutines stop** regardless of their state
4. Goroutines don't have return values — use **channels** to get results back
5. Sharing memory between goroutines without synchronization causes **data races** — use channels or `sync` primitives
