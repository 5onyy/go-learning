# Channels in Go

A **channel** is a typed conduit (pipe) that allows goroutines to **communicate and synchronize** by sending and receiving values. It's Go's primary mechanism for safe communication between concurrent goroutines — following the Go philosophy:

> *"Do not communicate by sharing memory; instead, share memory by communicating."*

---

## Creating Channels

```go
ch := make(chan int)        // unbuffered channel of ints
ch := make(chan string, 5)  // buffered channel with capacity 5
```

---

## Send & Receive

```go
ch <- value   // send value INTO the channel
value := <-ch // receive value FROM the channel
```

Both operations **block** until the other side is ready (for unbuffered channels).

---

## Unbuffered vs Buffered

### Unbuffered (synchronous)

```go
ch := make(chan int)
```

- Sender **blocks** until a receiver is ready
- Receiver **blocks** until a sender sends
- Guarantees synchronization at the point of communication

```go
go func() {
    ch <- 42  // blocks here until main receives
}()
fmt.Println(<-ch) // blocks here until goroutine sends
// Output: 42
```

### Buffered (asynchronous up to capacity)

```go
ch := make(chan int, 3)
```

- Sender only blocks when the buffer is **full**
- Receiver only blocks when the buffer is **empty**
- Decouples sender and receiver timing

```go
ch <- 1  // doesn't block (buffer has space)
ch <- 2
ch <- 3
ch <- 4  // BLOCKS — buffer is full

fmt.Println(<-ch) // 1
```

---

## Directional Channels

You can restrict a channel to send-only or receive-only in function signatures:

```go
func producer(ch chan<- int) {  // send-only
    ch <- 100
}

func consumer(ch <-chan int) {  // receive-only
    fmt.Println(<-ch)
}
```

This enforces correct usage at compile time.

---

## Closing a Channel

```go
close(ch)
```

- Only the **sender** should close a channel
- Receiving from a closed channel returns the zero value immediately
- Sending to a closed channel **panics**

```go
v, ok := <-ch
// ok is false if channel is closed and empty
```

---

## Ranging Over a Channel

```go
ch := make(chan int, 3)
ch <- 1
ch <- 2
ch <- 3
close(ch)  // must close, otherwise range blocks forever

for v := range ch {
    fmt.Println(v)  // 1, 2, 3
}
```

---

## `select` Statement

`select` lets a goroutine wait on **multiple channels** at once — like a switch for channels:

```go
select {
case msg := <-ch1:
    fmt.Println("from ch1:", msg)
case msg := <-ch2:
    fmt.Println("from ch2:", msg)
case ch3 <- "hello":
    fmt.Println("sent to ch3")
default:
    fmt.Println("no channel ready")  // non-blocking
}
```

---

## Common Patterns

### Done / Quit Signal

```go
done := make(chan struct{})  // struct{} uses zero bytes

go func() {
    // ... do work ...
    close(done)  // signal completion
}()

<-done  // wait for goroutine to finish
```

### Fan-out (one sender, many receivers)

```go
for i := 0; i < numWorkers; i++ {
    go worker(jobs)  // all read from the same channel
}
```

### Pipeline

```go
func generate(nums ...int) <-chan int { ... }
func square(in <-chan int) <-chan int { ... }

// Chain stages together
out := square(generate(2, 3, 4))
```

---

## Summary Table

| Feature        | Unbuffered                      | Buffered                          |
|----------------|---------------------------------|-----------------------------------|
| Sync point     | Yes (sender & receiver meet)    | No (up to capacity)               |
| Blocks sender  | Until receiver ready            | Until buffer full                 |
| Blocks receiver| Until sender sends              | Until buffer non-empty            |
| Use case       | Strict synchronization          | Decoupled producers/consumers     |

Channels, combined with goroutines, are the foundation of Go's concurrency model and make safe concurrent programming much more approachable than traditional lock-based approaches.
