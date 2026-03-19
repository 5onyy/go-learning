# Go `io` Package Cheatsheet

> Source: [pkg.go.dev/io@go1.26.1](https://pkg.go.dev/io@go1.26.1)  
> The `io` package provides basic interfaces to I/O primitives. It wraps low-level operations into shared public interfaces.

```go
import "io"
```

---

## Table of Contents
1. [Sentinel Values](#1-sentinel-values)
2. [Core Interfaces](#2-core-interfaces)
3. [Composite Interfaces](#3-composite-interfaces)
4. [Functions — Reading](#4-functions--reading)
5. [Functions — Writing](#5-functions--writing)
6. [Functions — Copying](#6-functions--copying)
7. [Special Readers](#7-special-readers)
8. [Pipe](#8-pipe)
9. [SectionReader](#9-sectionreader)
10. [NopCloser & Discard](#10-nopcloser--discard)
11. [Quick Reference Table](#11-quick-reference-table)

---

## 1. Sentinel Values

### Errors

| Variable            | Value                                        | When returned                                      |
|---------------------|----------------------------------------------|----------------------------------------------------|
| `io.EOF`            | `errors.New("EOF")`                          | Graceful end of input (not an error to propagate)  |
| `io.ErrUnexpectedEOF` | `"unexpected EOF"`                         | EOF in the middle of a fixed-size structure        |
| `io.ErrShortWrite`  | `"short write"`                              | Write accepted fewer bytes than given              |
| `io.ErrShortBuffer` | `"short buffer"`                             | Buffer too small for the required read             |
| `io.ErrNoProgress`  | `"multiple Read calls return no data or error"` | Broken Reader implementation                    |
| `io.ErrClosedPipe`  | `"io: read/write on closed pipe"`            | R/W on a closed `io.Pipe`                         |

### Seek Constants

```go
io.SeekStart   = 0  // relative to start of file
io.SeekCurrent = 1  // relative to current offset
io.SeekEnd     = 2  // relative to end of file
```

```go
// Example: seek to 10 bytes before the end of a file
f, _ := os.Open("file.txt")
offset, err := f.Seek(-10, io.SeekEnd)
```

---

## 2. Core Interfaces

These are the building blocks of all I/O in Go. Any type implementing these can
be used with `io` functions.

### `io.Reader`

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

- Reads up to `len(p)` bytes into `p`.
- Returns `(0, io.EOF)` when no more data.
- **Never assume `n == len(p)`** — loop or use helper functions.

```go
// Manual read loop
buf := make([]byte, 512)
for {
    n, err := r.Read(buf)
    if n > 0 {
        process(buf[:n])   // process BEFORE checking err
    }
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
}
```

### `io.Writer`

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

- Writes `len(p)` bytes from `p`.
- Returns error if `n < len(p)`.

```go
// Write a byte slice
n, err := w.Write([]byte("hello, world\n"))
```

### `io.Closer`

```go
type Closer interface {
    Close() error
}
```

```go
f, _ := os.Open("file.txt")
defer f.Close()   // always defer Close
```

### `io.Seeker`

```go
type Seeker interface {
    Seek(offset int64, whence int) (int64, error)
}
```

```go
f, _ := os.Open("file.txt")
// Jump to byte 100 from start
pos, err := f.Seek(100, io.SeekStart)

// Get current position
cur, _ := f.Seek(0, io.SeekCurrent)

// Get file size
size, _ := f.Seek(0, io.SeekEnd)
```

### `io.ReaderAt`

```go
type ReaderAt interface {
    ReadAt(p []byte, off int64) (n int, err error)
}
```

- Reads at a specific byte offset **without moving** the underlying seek pointer.
- Supports **parallel reads** from different offsets.

```go
f, _ := os.Open("file.txt")
buf := make([]byte, 20)
n, err := f.ReadAt(buf, 50)   // read 20 bytes starting at offset 50
```

### `io.WriterAt`

```go
type WriterAt interface {
    WriteAt(p []byte, off int64) (n int, err error)
}
```

```go
f, _ := os.OpenFile("file.bin", os.O_RDWR, 0644)
f.WriteAt([]byte{0xFF, 0xFE}, 4)   // overwrite bytes at offset 4
```

### `io.ReaderFrom` / `io.WriterTo`

```go
type ReaderFrom interface { ReadFrom(r Reader) (n int64, err error) }
type WriterTo   interface { WriteTo(w Writer)  (n int64, err error) }
```

- `io.Copy` automatically uses these if implemented — more efficient (avoids tmp buffer).

---

## 3. Composite Interfaces

Convenience combinations — a single type variable can satisfy multiple roles.

| Interface          | Embeds                        | Common implementors       |
|--------------------|-------------------------------|---------------------------|
| `ReadWriter`       | `Reader` + `Writer`           | `bytes.Buffer`            |
| `ReadCloser`       | `Reader` + `Closer`           | `http.Response.Body`      |
| `WriteCloser`      | `Writer` + `Closer`           | `os.File`, gzip writer    |
| `ReadWriteCloser`  | `Reader` + `Writer` + `Closer`| `net.Conn`                |
| `ReadSeeker`       | `Reader` + `Seeker`           | `bytes.Reader`, `os.File` |
| `WriteSeeker`      | `Writer` + `Seeker`           | `os.File`                 |
| `ReadSeekCloser`   | `Reader` + `Seeker` + `Closer`| `os.File`                 |
| `ReadWriteSeeker`  | `Reader` + `Writer` + `Seeker`| `os.File`                 |

```go
// Accept any readable+closeable source (e.g. file or HTTP body)
func processBody(rc io.ReadCloser) {
    defer rc.Close()
    data, _ := io.ReadAll(rc)
    fmt.Println(string(data))
}

// Works with both:
processBody(os.Stdin)
processBody(resp.Body)
```

---

## 4. Functions — Reading

### `io.ReadAll` *(go1.16+)*

```go
func ReadAll(r Reader) ([]byte, error)
```

Reads everything from `r` until EOF. Returns all bytes.

```go
import (
    "io"
    "strings"
    "fmt"
)

r := strings.NewReader("Go is great!")
data, err := io.ReadAll(r)
if err != nil { log.Fatal(err) }
fmt.Println(string(data))  // "Go is great!"

// Common: read HTTP response body
resp, _ := http.Get("https://example.com")
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
```

### `io.ReadFull`

```go
func ReadFull(r Reader, buf []byte) (n int, err error)
```

Reads **exactly** `len(buf)` bytes. Returns `ErrUnexpectedEOF` if source ends early.

```go
// Read a fixed 16-byte header
header := make([]byte, 16)
n, err := io.ReadFull(r, header)
if err == io.ErrUnexpectedEOF {
    fmt.Println("file too short")
} else if err != nil {
    log.Fatal(err)
}
// n is guaranteed to be 16 when err == nil
```

### `io.ReadAtLeast`

```go
func ReadAtLeast(r Reader, buf []byte, min int) (n int, err error)
```

Reads at least `min` bytes into `buf`. Returns `ErrUnexpectedEOF` if fewer available.

```go
buf := make([]byte, 64)
n, err := io.ReadAtLeast(r, buf, 10)   // need at least 10 bytes
fmt.Printf("read %d bytes\n", n)
```

---

## 5. Functions — Writing

### `io.WriteString`

```go
func WriteString(w Writer, s string) (n int, err error)
```

Writes a string to any `io.Writer` (avoids `[]byte` conversion if `w` implements `StringWriter`).

```go
// Write to stdout
io.WriteString(os.Stdout, "Hello, World!\n")

// Write to a file
f, _ := os.Create("out.txt")
defer f.Close()
io.WriteString(f, "line one\n")
io.WriteString(f, "line two\n")
```

---

## 6. Functions — Copying

### `io.Copy`

```go
func Copy(dst Writer, src Reader) (written int64, err error)
```

Copies all bytes from `src` to `dst` until EOF. The most commonly used copy function.

```go
// Copy a file
src, _ := os.Open("src.txt")
defer src.Close()
dst, _ := os.Create("dst.txt")
defer dst.Close()

n, err := io.Copy(dst, src)
fmt.Printf("copied %d bytes\n", n)

// Copy HTTP body to stdout
resp, _ := http.Get("https://example.com")
defer resp.Body.Close()
io.Copy(os.Stdout, resp.Body)

// Copy between strings
r := strings.NewReader("hello")
var buf bytes.Buffer
io.Copy(&buf, r)
fmt.Println(buf.String())  // "hello"
```

### `io.CopyN`

```go
func CopyN(dst Writer, src Reader, n int64) (written int64, err error)
```

Copies exactly `n` bytes. Returns `EOF` error if source runs out early.

```go
// Copy only the first 1024 bytes
n, err := io.CopyN(dst, src, 1024)
if err != nil && err != io.EOF {
    log.Fatal(err)
}
```

### `io.CopyBuffer` *(go1.5+)*

```go
func CopyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error)
```

Like `Copy` but uses a caller-supplied buffer — avoids allocation in hot paths.

```go
buf := make([]byte, 32*1024)   // reuse 32 KiB buffer
for _, path := range files {
    src, _ := os.Open(path)
    dst, _ := os.Create(filepath.Join("out", filepath.Base(path)))
    io.CopyBuffer(dst, src, buf)   // no alloc each iteration
    src.Close()
    dst.Close()
}
```

---

## 7. Special Readers

### `io.LimitReader`

```go
func LimitReader(r Reader, n int64) Reader
```

Returns a `Reader` that stops with EOF after `n` bytes. Prevents reading more than allowed.

```go
// Only read the first 512 bytes of a file
f, _ := os.Open("big.log")
defer f.Close()
limited := io.LimitReader(f, 512)
data, _ := io.ReadAll(limited)
fmt.Println(len(data))  // at most 512

// Protect against oversized HTTP bodies
body := io.LimitReader(resp.Body, 10<<20)   // cap at 10 MiB
io.ReadAll(body)
```

### `io.MultiReader`

```go
func MultiReader(readers ...Reader) Reader
```

Concatenates multiple `Reader`s — reads them sequentially as if they were one stream.

```go
// Prepend a header to a file body
header := strings.NewReader("===HEADER===\n")
f, _ := os.Open("body.txt")
defer f.Close()

combined := io.MultiReader(header, f)
io.Copy(os.Stdout, combined)   // prints header + file contents

// Chain multiple byte buffers
r := io.MultiReader(
    bytes.NewReader([]byte("part1")),
    bytes.NewReader([]byte("part2")),
    bytes.NewReader([]byte("part3")),
)
data, _ := io.ReadAll(r)
fmt.Println(string(data))  // "part1part2part3"
```

### `io.TeeReader`

```go
func TeeReader(r Reader, w Writer) Reader
```

Returns a `Reader` that **simultaneously writes** to `w` everything it reads from `r`.
Useful for logging, hashing, or inspecting data in-flight.

```go
// Log all bytes read from a connection while processing them
var log bytes.Buffer
tee := io.TeeReader(conn, &log)

// Process the data (also writes to log)
io.Copy(processor, tee)

// log now contains everything that was processed
fmt.Println("logged:", log.String())

// Compute hash while reading
h := sha256.New()
tee := io.TeeReader(resp.Body, h)
io.ReadAll(tee)
fmt.Printf("SHA256: %x\n", h.Sum(nil))
```

---

## 8. Pipe

```go
func Pipe() (*PipeReader, *PipeWriter)
```

Creates a **synchronous in-memory pipe**. Writes block until a reader consumes the data
(no internal buffer). Ideal for connecting components that expect a `Reader` or `Writer`.

```go
pr, pw := io.Pipe()

// Writer goroutine
go func() {
    fmt.Fprintln(pw, "hello through pipe")
    pw.Close()   // signals EOF to reader
}()

// Reader (blocking — waits for data)
data, _ := io.ReadAll(pr)
fmt.Println(string(data))   // "hello through pipe\n"
```

```go
// Real-world: stream gzip encoding through an HTTP request
pr, pw := io.Pipe()
go func() {
    gw := gzip.NewWriter(pw)
    io.Copy(gw, dataSource)
    gw.Close()
    pw.Close()
}()
http.Post(url, "application/gzip", pr)
```

```go
// Close with a custom error
pr, pw := io.Pipe()
go func() {
    pw.CloseWithError(errors.New("upstream failure"))
}()
_, err := io.ReadAll(pr)
fmt.Println(err)   // "upstream failure"
```

---

## 9. SectionReader

```go
func NewSectionReader(r ReaderAt, off int64, n int64) *SectionReader
```

Presents a **window** (slice) of a `ReaderAt` — read/seek within `[off, off+n)` without
loading the whole source. Common for parsing binary file formats.

```go
f, _ := os.Open("archive.bin")
defer f.Close()

// Read only bytes 100..199 (100 bytes)
section := io.NewSectionReader(f, 100, 100)

data := make([]byte, 20)
section.Read(data)                  // reads from offset 100
section.Seek(50, io.SeekStart)      // seek to offset 150 in the file
section.ReadAt(data, 80)            // reads at offset 180 in the file

fmt.Println("section size:", section.Size())  // 100

// Outer returns original reader + bounds (go1.22+)
orig, off, size := section.Outer()
fmt.Println(orig, off, size)
```

---

## 10. NopCloser & Discard

### `io.NopCloser` *(go1.16+)*

```go
func NopCloser(r Reader) ReadCloser
```

Wraps a `Reader` with a **no-op `Close`** method. Use when an API expects a
`ReadCloser` but you only have a `Reader`.

```go
// strings.NewReader has no Close — wrap it
rc := io.NopCloser(strings.NewReader("test body"))
processBody(rc)   // function expects io.ReadCloser

// Useful in tests
func makeBody(s string) io.ReadCloser {
    return io.NopCloser(strings.NewReader(s))
}
```

### `io.Discard`

```go
var Discard Writer
```

A `Writer` that **silently discards** all bytes written to it. Useful for draining
readers you don't care about (avoids nil pointer issues).

```go
// Drain an HTTP body to allow connection reuse
resp, _ := http.Get(url)
io.Copy(io.Discard, resp.Body)
resp.Body.Close()

// Benchmark a reader without caring about content
n, _ := io.Copy(io.Discard, bigReader)
fmt.Printf("drained %d bytes\n", n)

// Suppress output in tests
fmt.Fprintln(io.Discard, "silenced")
```

---

## 11. Quick Reference Table

| Function / Type          | Signature / Description                                              | Most common use                          |
|--------------------------|----------------------------------------------------------------------|------------------------------------------|
| `io.ReadAll`             | `(r Reader) → []byte`                                                | Read entire reader into memory           |
| `io.ReadFull`            | `(r Reader, buf []byte) → n, err`                                    | Read exactly N bytes                     |
| `io.ReadAtLeast`         | `(r Reader, buf []byte, min int) → n, err`                           | Read at least min bytes                  |
| `io.WriteString`         | `(w Writer, s string) → n, err`                                      | Write a string to any Writer             |
| `io.Copy`                | `(dst Writer, src Reader) → written, err`                            | Stream data between Reader and Writer    |
| `io.CopyN`               | `(dst, src, n int64) → written, err`                                 | Copy exactly N bytes                     |
| `io.CopyBuffer`          | `(dst, src, buf []byte) → written, err`                              | Copy with reusable buffer (performance)  |
| `io.LimitReader`         | `(r Reader, n int64) → Reader`                                       | Cap max bytes readable                   |
| `io.MultiReader`         | `(readers ...Reader) → Reader`                                       | Concatenate multiple sources             |
| `io.MultiWriter`         | `(writers ...Writer) → Writer`                                       | Fan-out writes to multiple targets       |
| `io.TeeReader`           | `(r Reader, w Writer) → Reader`                                      | Read and simultaneously copy elsewhere   |
| `io.Pipe`                | `() → (*PipeReader, *PipeWriter)`                                    | In-memory channel between goroutines     |
| `io.NewSectionReader`    | `(r ReaderAt, off, n int64) → *SectionReader`                        | Read a byte range from a ReaderAt        |
| `io.NopCloser`           | `(r Reader) → ReadCloser`                                            | Wrap Reader as ReadCloser (no-op close)  |
| `io.Discard`             | `Writer` — silently drops all writes                                 | Drain / suppress output                  |
| `io.EOF`                 | `error` — end of input                                               | Check with `err == io.EOF`               |
| `io.SeekStart/Current/End` | `int` constants for `Seek` whence                                  | File positioning                         |

### `io.MultiWriter` Example

```go
func MultiWriter(writers ...Writer) Writer
```

```go
// Write to both a file and stdout simultaneously
f, _ := os.Create("output.log")
defer f.Close()

mw := io.MultiWriter(os.Stdout, f)
fmt.Fprintln(mw, "goes to both stdout and file")
io.WriteString(mw, "also this line\n")

// Compute hash while writing to disk
h := sha256.New()
mw2 := io.MultiWriter(destFile, h)
io.Copy(mw2, srcFile)
fmt.Printf("SHA256: %x\n", h.Sum(nil))
```
