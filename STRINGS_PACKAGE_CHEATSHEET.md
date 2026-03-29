# Go `strings` Package Cheatsheet

> Complete reference for the `strings` standard library package — searching, transforming, splitting, joining, and building strings.

```go
import "strings"
```

---

## Table of Contents
1. [Searching & Matching](#1-searching--matching)
2. [Case Conversion](#2-case-conversion)
3. [Trimming](#3-trimming)
4. [Splitting & Joining](#4-splitting--joining)
5. [Replacing](#5-replacing)
6. [Prefix & Suffix](#6-prefix--suffix)
7. [Comparison](#7-comparison)
8. [Repeating & Counting](#8-repeating--counting)
9. [Mapping & Transforming](#9-mapping--transforming)
10. [strings.Builder — Efficient Concatenation](#10-stringsbuilder--efficient-concatenation)
11. [strings.Reader — Read from a String](#11-stringsreader--read-from-a-string)
12. [strings.Replacer — Multi-Pattern Replace](#12-stringsreplacer--multi-pattern-replace)
13. [Common Patterns & Recipes](#13-common-patterns--recipes)

---

## 1. Searching & Matching

### Contains / ContainsAny / ContainsRune

```go
strings.Contains("seafood", "foo")       // true  — substring match
strings.Contains("seafood", "bar")       // false

strings.ContainsAny("failure", "aeiou")  // true  — any single char in set
strings.ContainsAny("crwth", "aeiou")    // false

strings.ContainsRune("hello", 'e')       // true  — single rune match
strings.ContainsRune("hello", 'x')       // false

// ContainsFunc (Go 1.21+) — match by predicate
strings.ContainsFunc("hello123", unicode.IsDigit)  // true
```

### Index / LastIndex

```go
strings.Index("go gopher", "go")          // 0   — first occurrence (byte index)
strings.Index("go gopher", "xyz")         // -1  — not found

strings.LastIndex("go gopher go", "go")   // 10  — last occurrence

strings.IndexByte("golang", 'l')          // 2   — first byte match
strings.IndexRune("Hello, 世界", '世')     // 7   — first rune match (byte offset)

strings.IndexAny("chicken", "aeiou")      // 2   — first char matching any in set
strings.LastIndexAny("chicken", "aeiou")  // 5

strings.IndexFunc("Hello123", unicode.IsDigit)      // 5 — first rune satisfying func
strings.LastIndexFunc("Hello123", unicode.IsLetter)  // 4
```

---

## 2. Case Conversion

```go
strings.ToUpper("hello")       // "HELLO"
strings.ToLower("HELLO")       // "hello"
strings.ToTitle("hello world") // "HELLO WORLD" (Unicode title case — uppercases everything)

// Title-case each word (deprecated in favor of golang.org/x/text/cases)
strings.Title("hello world")   // "Hello World"  (deprecated since Go 1.18)

// Locale-aware case conversion
strings.ToUpperSpecial(unicode.TurkishCase, "i")  // "İ" (dotted I)
strings.ToLowerSpecial(unicode.TurkishCase, "I")  // "ı" (dotless i)
strings.ToTitleSpecial(unicode.TurkishCase, "i")  // "İ"
```

---

## 3. Trimming

### Whitespace

```go
strings.TrimSpace("  hello  \n")     // "hello"      — removes leading & trailing whitespace
```

### Specific Characters

```go
strings.Trim("***hello***", "*")     // "hello"      — removes chars from BOTH ends
strings.TrimLeft("***hello", "*")    // "hello"      — left end only
strings.TrimRight("hello***", "*")   // "hello"      — right end only
```

### Substrings (Prefix / Suffix)

```go
strings.TrimPrefix("HelloWorld", "Hello")  // "World"     — removes exact prefix
strings.TrimPrefix("HelloWorld", "Bye")    // "HelloWorld" — no match, returns original

strings.TrimSuffix("main.go", ".go")       // "main"
strings.TrimSuffix("main.go", ".rs")       // "main.go"
```

### Predicate

```go
strings.TrimFunc("¡¡Hello!!", func(r rune) bool {
    return !unicode.IsLetter(r) && !unicode.IsNumber(r)
})
// "Hello"

strings.TrimLeftFunc("123abc", unicode.IsDigit)   // "abc"
strings.TrimRightFunc("abc123", unicode.IsDigit)  // "abc"
```

### TrimPrefix vs TrimLeft — Important Distinction

```go
// TrimPrefix removes an exact prefix string once
strings.TrimPrefix("aabcaa", "aa")   // "bcaa"

// TrimLeft removes individual characters from a set repeatedly
strings.TrimLeft("aabcaa", "a")      // "bcaa"  — looks same here, but...
strings.TrimLeft("abba", "ab")       // ""      — removes all leading 'a' or 'b'
strings.TrimPrefix("abba", "ab")     // "ba"    — removes "ab" prefix once
```

---

## 4. Splitting & Joining

### Split

```go
strings.Split("a,b,c", ",")           // ["a" "b" "c"]
strings.Split("abc", "")              // ["a" "b" "c"]  — split every char
strings.Split("a,b,,c", ",")          // ["a" "b" "" "c"] — empty strings preserved

strings.SplitN("a,b,c,d", ",", 2)    // ["a" "b,c,d"]   — at most 2 pieces
strings.SplitN("a,b,c", ",", -1)     // ["a" "b" "c"]   — -1 means no limit (same as Split)
```

### SplitAfter — Keep the Separator

```go
strings.SplitAfter("a,b,c", ",")      // ["a," "b," "c"]  — separator stays with left part
strings.SplitAfterN("a,b,c", ",", 2)  // ["a," "b,c"]
```

### Fields — Split by Whitespace

```go
strings.Fields("  foo   bar  baz  ")   // ["foo" "bar" "baz"]  — splits on any whitespace run

// Custom field separator
strings.FieldsFunc("foo1bar2baz", unicode.IsDigit)  // ["foo" "bar" "baz"]
```

### Join

```go
parts := []string{"usr", "local", "bin"}
strings.Join(parts, "/")              // "usr/local/bin"

words := []string{"Hello", "World"}
strings.Join(words, " ")             // "Hello World"

// Join empty slice
strings.Join([]string{}, ",")        // ""
```

### Cut (Go 1.18+) — Split Around First Separator

```go
before, after, found := strings.Cut("user@host", "@")
// before="user", after="host", found=true

before, after, found = strings.Cut("noatsign", "@")
// before="noatsign", after="", found=false

// CutPrefix / CutSuffix (Go 1.20+)
after, found := strings.CutPrefix("Hello, World", "Hello, ")
// after="World", found=true

before, found := strings.CutSuffix("file.tar.gz", ".gz")
// before="file.tar", found=true
```

---

## 5. Replacing

```go
// Replace n occurrences (-1 = all)
strings.Replace("oink oink oink", "oink", "moo", 2)   // "moo moo oink"
strings.Replace("oink oink oink", "oink", "moo", -1)  // "moo moo moo"

// ReplaceAll (shorthand for n=-1, Go 1.12+)
strings.ReplaceAll("oink oink", "oink", "moo")        // "moo moo"
```

---

## 6. Prefix & Suffix

```go
strings.HasPrefix("Gopher", "Go")     // true
strings.HasPrefix("Gopher", "go")     // false  — case-sensitive

strings.HasSuffix("main.go", ".go")   // true
strings.HasSuffix("main.go", ".rs")   // false
```

---

## 7. Comparison

```go
// Case-sensitive (use == operator directly)
"abc" == "abc"   // true
"abc" == "ABC"   // false

// Case-insensitive comparison (no allocation, faster than ToLower+==)
strings.EqualFold("Go", "go")         // true
strings.EqualFold("ABC", "abc")       // true
strings.EqualFold("ß", "ss")          // false  (not linguistic equivalence)

// Lexicographic comparison
strings.Compare("a", "b")             // -1   (a < b)
strings.Compare("b", "a")             // +1   (b > a)
strings.Compare("a", "a")             //  0   (equal)
// Note: strings.Compare exists for symmetry with bytes.Compare;
//       idiomatic Go prefers using ==, <, > directly.
```

---

## 8. Repeating & Counting

```go
strings.Repeat("na", 4)              // "nananana"
strings.Repeat("-", 40)              // "----------------------------------------"
strings.Repeat("abc", 0)             // ""

strings.Count("cheese", "e")         // 3   — non-overlapping occurrences
strings.Count("five", "")            // 5   — empty string → len(s) + 1
```

---

## 9. Mapping & Transforming

### Map — Apply Function to Every Rune

```go
// ROT13
rot13 := strings.Map(func(r rune) rune {
    switch {
    case r >= 'A' && r <= 'Z':
        return 'A' + (r-'A'+13)%26
    case r >= 'a' && r <= 'z':
        return 'a' + (r-'a'+13)%26
    }
    return r
}, "Hello World")
// "Uryyb Jbeyq"

// Remove all digits
noDigits := strings.Map(func(r rune) rune {
    if unicode.IsDigit(r) {
        return -1   // -1 drops the rune
    }
    return r
}, "h3ll0 w0rld")
// "hll wrld"

// Convert non-ASCII to '?'
ascii := strings.Map(func(r rune) rune {
    if r > 127 {
        return '?'
    }
    return r
}, "café")
// "caf?"
```

---

## 10. strings.Builder — Efficient Concatenation

`strings.Builder` minimizes memory allocations when building strings incrementally. It's the idiomatic way to concatenate in a loop.

```go
var sb strings.Builder

sb.WriteString("Hello")        // append a string
sb.WriteByte(',')              // append a single byte
sb.WriteRune(' ')              // append a single rune (UTF-8 safe)
sb.WriteString("World!")

result := sb.String()          // "Hello, World!"
fmt.Println(sb.Len())          // 13  — current byte length

sb.Reset()                     // clear and reuse
```

### Builder vs Concatenation — Why It Matters

```go
// BAD: O(n²) — creates a new string each iteration
func joinBad(parts []string) string {
    var result string
    for _, p := range parts {
        result += p + ","
    }
    return result
}

// GOOD: O(n) — single allocation with Builder
func joinGood(parts []string) string {
    var sb strings.Builder
    for i, p := range parts {
        if i > 0 {
            sb.WriteByte(',')
        }
        sb.WriteString(p)
    }
    return sb.String()
}

// BEST: use strings.Join when you just need a separator
func joinBest(parts []string) string {
    return strings.Join(parts, ",")
}
```

### Preallocate with Grow

```go
var sb strings.Builder
sb.Grow(1024)                 // pre-allocate 1024 bytes to avoid reallocs
for _, line := range lines {
    sb.WriteString(line)
    sb.WriteByte('\n')
}
```

---

## 11. strings.Reader — Read from a String

`strings.Reader` implements `io.Reader`, `io.ReaderAt`, `io.Seeker`, `io.WriterTo`, and `io.ByteScanner`, letting you use a string anywhere an `io.Reader` is expected.

```go
r := strings.NewReader("Hello, World!")

fmt.Println(r.Len())          // 13  — unread bytes
fmt.Println(r.Size())         // 13  — original length

// Use as io.Reader
buf := make([]byte, 5)
n, _ := r.Read(buf)           // n=5, buf="Hello"

// Seek back to start
r.Seek(0, io.SeekStart)

// Copy to a writer
io.Copy(os.Stdout, r)         // prints "Hello, World!"

// Common use: pass a string to any function expecting io.Reader
resp := http.NewRequest("POST", url, strings.NewReader(`{"key":"value"}`))
```

---

## 12. strings.Replacer — Multi-Pattern Replace

`strings.Replacer` is optimized for applying multiple replacement rules at once. It's safe for concurrent use.

```go
// Pairs: old1, new1, old2, new2, ...
r := strings.NewReplacer(
    "&", "&amp;",
    "<", "&lt;",
    ">", "&gt;",
    `"`, "&quot;",
)

escaped := r.Replace(`<div class="box">Hello & World</div>`)
// "&lt;div class=&quot;box&quot;&gt;Hello &amp; World&lt;/div&gt;"

// Also supports writing to an io.Writer
r.WriteString(os.Stdout, "<b>bold</b>")
// &lt;b&gt;bold&lt;/b&gt;
```

### Template-Style Replacement

```go
r := strings.NewReplacer(
    "{name}", "Alice",
    "{role}", "Admin",
    "{org}",  "Acme Corp",
)

msg := r.Replace("Hello {name}, you are a {role} at {org}.")
// "Hello Alice, you are a Admin at Acme Corp."
```

---

## 13. Common Patterns & Recipes

### Check If String Is Empty

```go
if s == "" { }           // idiomatic — fast, clear
if len(s) == 0 { }       // also fine, equivalent
```

### Reverse a String (rune-safe)

```go
func reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
reverse("Hello, 世界")  // "界世 ,olleH"
```

### Remove All Whitespace

```go
noSpaces := strings.Map(func(r rune) rune {
    if unicode.IsSpace(r) {
        return -1
    }
    return r
}, "  h e l l o  ")
// "hello"
```

### Check If String Contains Only Digits

```go
isNumeric := strings.IndexFunc(s, func(r rune) bool {
    return !unicode.IsDigit(r)
}) == -1 && s != ""
```

### Word Count

```go
count := len(strings.Fields("The quick brown fox"))  // 4
```

### Truncate to Max Length (rune-safe)

```go
func truncate(s string, maxRunes int) string {
    runes := []rune(s)
    if len(runes) > maxRunes {
        return string(runes[:maxRunes]) + "…"
    }
    return s
}
```

### Parse Key=Value Pairs

```go
line := "host=localhost port=5432 user=admin"
for _, field := range strings.Fields(line) {
    key, value, ok := strings.Cut(field, "=")
    if ok {
        fmt.Printf("%s → %s\n", key, value)
    }
}
```

### Pad a String

```go
func padLeft(s string, width int, pad byte) string {
    if len(s) >= width {
        return s
    }
    return strings.Repeat(string(pad), width-len(s)) + s
}

func padRight(s string, width int, pad byte) string {
    if len(s) >= width {
        return s
    }
    return s + strings.Repeat(string(pad), width-len(s))
}

padLeft("42", 5, '0')    // "00042"
padRight("hi", 10, '.')  // "hi........"
```

---

## Quick Reference Card

```
SEARCH        Contains  ContainsAny  ContainsRune  ContainsFunc
FIND INDEX    Index  LastIndex  IndexByte  IndexRune  IndexAny  IndexFunc
PREFIX/SUFFIX HasPrefix  HasSuffix  TrimPrefix  TrimSuffix  CutPrefix  CutSuffix
CASE          ToUpper  ToLower  ToTitle  EqualFold
TRIM          TrimSpace  Trim  TrimLeft  TrimRight  TrimFunc
SPLIT         Split  SplitN  SplitAfter  Fields  FieldsFunc  Cut
JOIN          Join
REPLACE       Replace  ReplaceAll  Map  NewReplacer
BUILD         Builder{WriteString, WriteByte, WriteRune, String, Reset, Grow}
READ          NewReader → io.Reader from a string
MISC          Repeat  Count  Compare  Clone (Go 1.20+)
```

## Performance Tips

| Task | Approach | Why |
|------|----------|-----|
| Loop concatenation | `strings.Builder` | O(n) vs O(n²) for `+=` |
| Case-insensitive compare | `strings.EqualFold` | No allocation (unlike `ToLower` + `==`) |
| Multi-pattern replace | `strings.NewReplacer` | Optimized trie-based matching |
| Known size | `Builder.Grow(n)` | Avoids reallocation |
| Simple join | `strings.Join` | Single allocation, cleanest API |
| Check prefix/suffix | `HasPrefix` / `HasSuffix` | O(len(prefix)), no allocation |
