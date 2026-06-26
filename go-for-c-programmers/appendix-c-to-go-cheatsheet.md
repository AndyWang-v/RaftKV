# Appendix A — C to Go Cheat Sheet

A fast lookup table for translating C into Go. Keep it open while you work. For the
full story behind each row, follow the chapter pointers. Code is shown as `c` and `go`.

## Types

| C | Go | Note |
|---|---|---|
| `char` | `byte` (alias for `uint8`) | C `char` signedness is platform-defined; Go `byte` is unsigned. |
| `signed char` | `int8` | |
| `unsigned char` | `byte` / `uint8` | |
| `short` / `unsigned short` | `int16` / `uint16` | |
| `int` / `unsigned int` | `int32` / `uint32`, or Go `int` | C `int` is usually 32-bit; Go `int` is the machine word (32 or 64-bit). |
| `long` / `long long` | `int64` (pick a fixed size) | C widths vary by platform; Go's are fixed. |
| `unsigned long long` | `uint64` | |
| `size_t` | `int` | Go uses signed `int` for lengths and indices. |
| `ptrdiff_t` | `int` | |
| `intptr_t` / `uintptr_t` | `uintptr` | pointer-sized integer; rarely needed. |
| `float` / `double` | `float32` / `float64` | `float64` is Go's default float. |
| `_Bool` / `bool` | `bool` | not an integer; no implicit `int`↔`bool`. |
| `void *` | `any` or `unsafe.Pointer` | `any` for a type-safe value; `unsafe.Pointer` for raw bytes. |
| `T *` | `*T` | pointer; **no pointer arithmetic**. |
| `T name[N]` | `[N]T` | fixed array; a **value** (copied on assignment). |
| `T *` used as buffer | `[]T` (slice) | the everyday dynamic array (ptr, len, cap). |
| `struct {...}` | `struct {...}` | similar layout; Go adds methods and visibility. |
| `union {...}` | (none) | use a struct, an `interface`, or `unsafe` to overlay. |
| `enum {...}` | `type X int` + `const` / `iota` | a distinct named type. |
| `char *` (text) | `string` (immutable) / `[]byte` (mutable) | `string` is UTF-8 bytes. |
| code point | `rune` (alias for `int32`) | one Unicode character. |
| `int (*)(int)` | `func(int) int` | functions are first-class values. |
| `_Complex` | `complex64` / `complex128` | built in. |

> **Rule of thumb.** Reach for the fixed-size type (`int32`, `uint64`) when you are
> matching a C ABI or a wire format; use plain `int` for ordinary counts and indices.

## Declarations

| Thing | C | Go |
|---|---|---|
| Variable, typed | `int x = 5;` | `var x int = 5` |
| Variable, inferred | `int x = 5;` | `x := 5` (inside a function) |
| Uninitialized | `int x;` (garbage) | `var x int` (zero value `0`) |
| Multiple | `int a = 1, b = 2;` | `a, b := 1, 2` |
| Constant | `#define N 10` / `const int N = 10;` | `const N = 10` |
| Enum group | `enum { A, B, C };` | `const ( A = iota; B; C )` |
| Fixed array | `int a[10];` | `var a [10]int` |
| Array literal | `int a[] = {1, 2, 3};` | `a := [...]int{1, 2, 3}` |
| Slice literal | (manual) | `s := []int{1, 2, 3}` |
| Pointer | `int *p;` | `var p *int` |
| Address / deref | `&x` / `*p` | `&x` / `*p` (same) |
| Function | `int add(int a, int b) {...}` | `func add(a, b int) int {...}` |
| Function pointer | `int (*fp)(int, int);` | `var fp func(int, int) int` |
| Struct type | `struct P { int x, y; };` | `type P struct { X, Y int }` |
| Distinct type | `typedef int Id;` (no teeth) | `type Id int` (no implicit conversion) |
| True alias | (same `typedef`) | `type Id = int` (Go 1.9+) |
| File-private global | `static int x;` | `var x int` (lowercase = unexported) |
| Function-local persistent | `static int n;` inside func | a package-level `var n int` |
| Declared elsewhere | `extern int x;` | `import` the package; use `pkg.X` |
| Public symbol | non-`static` in a header | Capitalized name (`Foo`) |
| Private symbol | `static` | lowercase name (`foo`) |

See Chapter 3 — Program Structure: Packages, Imports, and Visibility and Chapter 4 —
Types, Variables, and Constants.

## Operators and keywords that differ

| C | Go | Note |
|---|---|---|
| `p->field` | `p.field` | Go auto-dereferences; **no `->`**. |
| `(*p).field` | `p.field` | same result. |
| `cond ? a : b` | (none) | **no ternary**; use `if`/`else` or a tiny helper. |
| `x++`, `++x` (as a value) | `x++` (a **statement**, not an expression) | no `++x`; cannot use inside an expression. |
| `a, b` (comma operator) | (none) | use separate statements. |
| `if (x = f())` (assign in expr) | `if x := f(); x != 0 {}` | assignment is not an expression; use the init clause. |
| `(int)x` (cast) | `int(x)` (conversion) | the type is used like a function. |
| `&`, `\|`, `^`, `<<`, `>>` | same | `^` is also unary NOT. |
| `~x` (bitwise NOT) | `^x` | unary `^`. |
| (no operator) | `x &^ y` | bit clear (AND NOT). |
| `sizeof x` | `unsafe.Sizeof(x)` | see the Memory table. |
| `goto label;` | `goto label` | exists, but rare; prefer `defer` and labeled `break`. |

## Control flow

| C | Go | Note |
|---|---|---|
| `while (c) { }` | `for c { }` | one loop keyword. |
| `do { } while (c);` | `for { ...; if !c { break } }` | no `do`/`while`. |
| `for (i = 0; i < n; i++)` | `for i := 0; i < n; i++ { }` | no parentheses; braces required. |
| count `0..n-1` | `for i := range n { }` | integer range, Go 1.22+. |
| `for (;;)` | `for { }` | infinite loop. |
| walk an array | `for i, v := range s { }` | index and value. |
| `switch` falls through | `switch` **breaks** after each case | use `fallthrough` to continue to the next case. |
| `switch` on integers only | `switch` on any comparable type | cases may be expressions; `switch { }` replaces an if/else chain. |
| `if (c) { }` | `if c { }` | braces required; init clause `if x := f(); c { }`. |
| `break` / `continue` | same, plus labels | `break Outer`, `continue Outer`. |
| `goto cleanup;` | `defer cleanup()` | the idiomatic replacement. |

See Chapter 5 — Control Flow.

## Memory

| C | Go | Note |
|---|---|---|
| `malloc(sizeof *p)` | `p := new(T)` or `p := &T{}` | returns a zeroed `*T`. |
| `malloc(n * sizeof(T))` | `s := make([]T, n)` | slice of `n` zero values. |
| `calloc(n, sizeof(T))` | `make([]T, n)` | already zeroed (zero value). |
| `realloc(p, m)` | `s = append(s, ...)` / `slices.Grow(s, m)` | append grows the backing array. |
| `free(p)` | (nothing) | the garbage collector reclaims it. |
| `sizeof(T)` | `unsafe.Sizeof(T{})` | size in bytes, at compile time. |
| `sizeof(a)/sizeof(a[0])` | `len(a)` | length is part of the value. |
| choose stack vs heap | (compiler chooses) | escape analysis; returning `&local` is **safe**. |
| `memcpy(d, s, n)` | `copy(d, s)` | copies `min(len(d), len(s))` elements. |
| `memmove(d, s, n)` | `copy(d, s)` | `copy` handles overlap. |
| `memset(p, 0, n)` | `clear(s)` | also empties a map; or rely on the zero value. |
| array decays to a pointer | slice keeps len and cap | no decay; access is bounds-checked. |

See Chapter 7 — Pointers and Chapter 17 — Memory and the Garbage Collector.

## Strings

| C | Go | Note |
|---|---|---|
| `char *` / `char[]` | `string` (immutable) / `[]byte` (mutable) | `string` is UTF-8 bytes. |
| `strlen(s)` | `len(s)` | byte count, O(1); not the number of characters. |
| character count | `utf8.RuneCountInString(s)` | counts runes (code points). |
| `s[i]` | `s[i]` (a `byte`) | one byte; `for _, r := range s` gives runes. |
| `strcmp(a, b) == 0` | `a == b` | `==` compares contents. |
| ordering | `a < b`, `strings.Compare(a, b)` | lexical byte order. |
| `strcpy(d, s)` | `d = s` / `copy(dst, src)` for `[]byte` | string assignment is a cheap copy of the header. |
| `strcat(a, b)` | `a + b` | for loops, use a builder. |
| build incrementally | `strings.Builder` | avoids O(n²) reallocation. |
| `strchr` / `strstr` | `strings.IndexByte` / `strings.Index` | returns index or `-1`. |
| `strtok` | `strings.Split` / `strings.Fields` | returns a `[]string`. |
| `toupper` / case | `strings.ToUpper` | Unicode-aware. |
| `printf(fmt, ...)` | `fmt.Printf(fmt, ...)` | format verbs checked by `go vet`. |
| `sprintf(buf, ...)` | `s := fmt.Sprintf(...)` | returns a new string; no buffer sizing. |
| `snprintf(...)` | `fmt.Sprintf(...)` | no overflow to guard. |
| `fprintf(stderr, ...)` | `fmt.Fprintf(os.Stderr, ...)` | writes to any `io.Writer`. |
| `puts(s)` | `fmt.Println(s)` | adds a newline. |
| `atoi(s)` / `strtol` | `strconv.Atoi(s)` / `strconv.ParseInt` | returns `(value, error)`. |

See Chapter 8 — Arrays, Slices, and Strings.

## C header to Go package

| C header / function | Purpose | Go package(s) |
|---|---|---|
| `<stdio.h>` | print, files, buffered I/O | `fmt`, `os`, `bufio` |
| `<stdlib.h>` `atoi`, `strtol` | parse numbers | `strconv` |
| `<stdlib.h>` `malloc`, `free` | memory | built in (`make`, `new`) + the GC |
| `<stdlib.h>` `qsort`, `bsearch` | sort, search | `sort`, `slices` |
| `<stdlib.h>` `rand` | randomness | `math/rand/v2`, `crypto/rand` |
| `<stdlib.h>` `getenv`, `exit` | environment, exit | `os.Getenv`, `os.Exit` |
| `<string.h>` | string and memory ops | `strings`, `bytes` |
| `<ctype.h>` | character classes | `unicode` |
| `<math.h>` | math functions | `math` |
| `<time.h>` | clocks and dates | `time` |
| `<pthread.h>` | threads, mutexes, condvars | goroutines (`go`), `sync`, channels |
| `<errno.h>` | error codes | `error` values, `errors` (`syscall.Errno` for raw codes) |
| `<assert.h>` | assertions | none built in; use tests and `panic` |
| `<signal.h>` | signals | `os/signal` |
| `<unistd.h>`, `<fcntl.h>` | files, syscalls | `os`, `io`, `syscall` |
| `<stdint.h>` | fixed-width integers | built in (`int32`, `uint64`, ...) |
| `<stdbool.h>` | booleans | built in (`bool`) |
| `<stdarg.h>` | variadic args | variadic parameters `...T` |

See Chapter 20 — The Standard Library Tour.

## Idioms

**Error handling** (return codes / `errno` → `error` values):

```c
int fd = open(path, O_RDONLY);
if (fd < 0) { /* inspect errno */ return -1; }
/* use fd */
close(fd);
```

```go
f, err := os.Open(path)
if err != nil {
	return err
}
defer f.Close()
// use f
```

**Cleanup** (`goto cleanup` → `defer`, which runs in last-in-first-out order):

```c
FILE *a = fopen(pa, "r");
if (!a) goto done;
FILE *b = fopen(pb, "r");
if (!b) goto close_a;
/* work */
fclose(b);
close_a:
fclose(a);
done:
return;
```

```go
a, err := os.Open(pa)
if err != nil {
	return err
}
defer a.Close()
b, err := os.Open(pb)
if err != nil {
	return err
}
defer b.Close()
// work; a and b close automatically, b before a
```

**A struct of function pointers → an interface** (a vtable you build by hand → one the
compiler builds and type-checks):

```c
typedef struct {
	int  (*Read)(void *self, char *buf, int n);
	int  (*Close)(void *self);
} Reader;
```

```go
type Reader interface {
	Read(p []byte) (n int, err error)
	Close() error
}
```

Any type with `Read` and `Close` methods satisfies `Reader` automatically; there is no
`implements` keyword and no manual wiring. See Chapter 11 — Interfaces.

## The `go` command

| Task | Command |
|---|---|
| Run the package here | `go run .` |
| Build a binary | `go build -o app .` |
| Build everything | `go build ./...` |
| Install a tool | `go install example.com/cmd@latest` |
| Run tests | `go test ./...` |
| Verbose tests | `go test -v ./...` |
| Tests with the race detector | `go test -race ./...` |
| Coverage | `go test -cover ./...` |
| Benchmarks | `go test -bench=. -benchmem` |
| Format in place | `gofmt -w .` (or `go fmt ./...`) |
| Static checks | `go vet ./...` |
| New module | `go mod init example.com/app` |
| Tidy dependencies | `go mod tidy` |
| Upgrade one dependency | `go get example.com/pkg@latest` |
| List dependencies | `go list -m all` |
| Read docs | `go doc fmt.Println` |
| Show settings | `go env` |
| Cross-compile | `GOOS=linux GOARCH=arm64 go build .` |
| Scan for vulnerabilities | `govulncheck ./...` |

See Chapter 2 — Installing Go and the `go` Command.

## Handy one-liners

```go
b, err := os.ReadFile("f.txt")          // read a whole file
err = os.WriteFile("f.txt", b, 0o644)   // write a whole file
in, _ := io.ReadAll(os.Stdin)           // read all of stdin
parts := strings.Split(s, ",")          // split on a separator
s = strings.Join(parts, ",")            // join with a separator
n, err := strconv.Atoi("42")            // string -> int
s = strconv.Itoa(42)                    // int -> string
s = fmt.Sprintf("%s=%d", name, n)       // format into a string
slices.Sort(nums)                       // sort a slice in place
ok := slices.Contains(nums, x)          // membership test
keys := slices.Collect(maps.Keys(m))    // map keys -> slice (Go 1.23+)
time.Sleep(100 * time.Millisecond)      // pause
```
