# Chapter 21 — Testing

> **What you'll learn.** How Go builds testing into the language and the `go`
> tool: `_test.go` files, table-driven tests with subtests, benchmarks (including
> the new `b.Loop`), fuzzing, example functions that double as documentation,
> setup/teardown helpers, golden files, and test doubles — all with no external
> framework.

## C has no standard test framework; Go does

In C, testing is bring-your-own. You sprinkle `assert()`, write a `main` that
calls your checks, or pull in a third-party framework (CUnit, Check, Unity). There
is no single `make test` everyone agrees on, and the framework you pick is rarely
the one the next project uses.

Go puts testing **in the box**. The `testing` package and the `go test` command
ship with the toolchain. Every Go project tests the same way, so you can drop
into any codebase and run `go test ./...` immediately (see Chapter 2 — Installing
Go and the `go` Command).

| Concept | C | Go |
|---|---|---|
| Test framework | Third party (CUnit, Check) or `assert.h` | Built in: `testing` + `go test` |
| Where tests live | Separate files/dirs you wire up | `*_test.go` beside the code |
| How to run | Custom `make test` target | `go test ./...` |
| Benchmarks | Hand-rolled timers | `func BenchmarkXxx`, `go test -bench` |
| Fuzzing | Separate tool (libFuzzer, AFL) | Built in: `func FuzzXxx`, `go test -fuzz` |
| Coverage | `gcov`/`lcov` setup | `go test -cover` |

## Tests live beside the code

A test file ends in **`_test.go`** and sits in the **same directory** as the code
it tests. The `go build` command ignores these files, so test code never ships in
your binary; `go test` compiles them in.

```
strutil/
├── reverse.go        package strutil   (the code)
├── reverse_test.go   package strutil   (white-box tests)
└── export_test.go    package strutil_test (black-box tests, optional)
```

A test file may use one of two package clauses:

- **`package strutil`** — the *same* package as the code. The test can see
  unexported names (white-box testing). This is the common choice.
- **`package strutil_test`** — a separate test-only package that *imports*
  `strutil`. The test sees only the exported API, exactly as a real user would
  (black-box testing). Go allows this one exception to the "one package per
  directory" rule from Chapter 3 — Program Structure.

Here is the code we will test throughout the chapter:

```go
// file: reverse.go
package strutil

// Reverse returns s with its runes (Unicode code points) in reverse order.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
```

## Your first test

A test is a function named **`TestXxx`** (the name after `Test` must start with an
uppercase letter) that takes one argument, `*testing.T`.

```go
// file: reverse_test.go
package strutil

import "testing"

func TestReverse(t *testing.T) {
	got := Reverse("hello")
	want := "olleh"
	if got != want {
		t.Errorf("Reverse(%q) = %q, want %q", "hello", got, want)
	}
}
```

There are no assertion macros. You write a normal `if`, and if the result is
wrong you call a method on `t` to report it. Run it:

```sh
go test            # run tests in the current package
go test ./...      # run every test in the module
```

`t` gives you two ways to report a failure, and the difference matters:

| Method | Marks test failed? | Keeps going? | Use when |
|---|---|---|---|
| `t.Error` / `t.Errorf` | yes | **yes** | one of several independent checks |
| `t.Fatal` / `t.Fatalf` | yes | **no** (stops this test) | a precondition failed; later code would panic |
| `t.Log` / `t.Logf` | no | yes | extra detail (shown only with `-v` or on failure) |

> **Rule of thumb.** Use `t.Fatal` when continuing makes no sense — for example
> after `err != nil` from a setup step, because the next line would dereference a
> `nil`. Use `t.Error` to record a failed comparison and let the rest of the
> checks run, so you see *all* the problems at once.

A **helper** is a function that contains the actual checks. Call `t.Helper()`
inside it so failures are reported at the *caller's* line, not deep inside the
helper:

```go
func assertEqual(t *testing.T, got, want string) {
	t.Helper() // report failures at the caller, not here
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

## Table-driven tests with subtests

The idiomatic Go test is **table-driven**: a slice of test cases (the "table")
looped over, with each case run as a **subtest** via `t.Run`. This is the single
most important pattern in this chapter.

```go
func TestReverse(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: ""},
		{name: "ascii", in: "hello", want: "olleh"},
		{name: "palindrome", in: "level", want: "level"},
		{name: "unicode", in: "héllo", want: "olléh"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Reverse(tc.in)
			if got != tc.want {
				t.Errorf("Reverse(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
```

Why this pattern wins:

- **Adding a case is one line** in the table — no copy-pasted test function.
- `t.Run(name, ...)` gives each case its **own name**, so a failure says exactly
  which one broke: `--- FAIL: TestReverse/unicode`.
- You can run a single case: `go test -run TestReverse/unicode`.

> **C vs Go.** In C you would write a loop over an array of structs and `printf`
> the failures yourself, then remember to set the exit code. `t.Run` gives you
> named, individually-runnable, independently-failing subtests for free, and the
> framework tracks pass/fail.

## Running tests

The flags you will use constantly:

```sh
go test ./...                 # all packages
go test -v ./...              # verbose: show each test and its logs
go test -run TestReverse      # only tests whose name matches the regexp
go test -run TestReverse/ascii # only that subtest
go test -race ./...           # build with the race detector (Chapter 22)
go test -count=1 ./...         # disable the test cache (force a real run)
```

> **Watch out.** Go **caches** test results. If nothing the test depends on has
> changed, `go test` prints `(cached)` and does not re-run. That is great for
> speed but confusing when a test depends on something outside the build (a clock,
> a file, the network). Add **`-count=1`** to force a fresh run; it is the
> standard "turn off caching" idiom.

### Coverage

Coverage shows which lines your tests exercise.

```sh
go test -cover ./...                      # print a percentage per package
go test -coverprofile=cover.out ./...     # write a detailed profile
go tool cover -func=cover.out             # per-function coverage in the terminal
go tool cover -html=cover.out             # open an annotated, colorized view
```

> **C vs Go.** This replaces the `gcc --coverage` + `gcov`/`lcov` dance. No extra
> compiler flags, no separate tool to install — coverage is part of `go test`.

## Benchmarks

A benchmark is a function named **`BenchmarkXxx`** taking `*testing.B`. The
framework runs your code many times and reports nanoseconds per operation.

The modern form (Go 1.24) uses **`b.Loop()`**:

```go
func BenchmarkReverse(b *testing.B) {
	for b.Loop() { // runs the body enough times to get a stable measurement
		Reverse("the quick brown fox")
	}
}
```

`b.Loop()` is preferred because it keeps the loop count internal, prevents the
compiler from optimizing the call away, and excludes one-time setup automatically.
You will still see the older `b.N` form in existing code:

```go
func BenchmarkReverseOld(b *testing.B) {
	for i := 0; i < b.N; i++ { // the framework chooses b.N
		Reverse("the quick brown fox")
	}
}
```

If a benchmark needs expensive setup, do it first and then call `b.ResetTimer()`
so the setup is not counted:

```go
func BenchmarkProcess(b *testing.B) {
	data := buildLargeInput() // not part of what we measure
	b.ResetTimer()
	for b.Loop() {
		Process(data)
	}
}
```

Run benchmarks (they do **not** run during a normal `go test`):

```sh
go test -bench=.              # run all benchmarks
go test -bench=Reverse -benchmem  # also report allocations per operation
```

Sample output (`-benchmem` adds the last two columns):

```
BenchmarkReverse-8   28412641   42.11 ns/op   16 B/op   1 allocs/op
```

That reads: ran ~28 million times, 42 ns each, 16 bytes and 1 allocation per call.

## Fuzzing

A **fuzz test** feeds your code a flood of automatically generated inputs to find
cases you did not think of. It is built into `go test` (Go 1.18). A fuzz target is
named **`FuzzXxx`** and takes `*testing.F`.

```go
import "unicode/utf8"

func FuzzReverse(f *testing.F) {
	f.Add("hello") // seed corpus: example inputs to start from
	f.Add("héllo")
	f.Add("")

	f.Fuzz(func(t *testing.T, s string) {
		rev := Reverse(s)
		doubled := Reverse(rev)
		// Property 1: reversing twice returns the original.
		if s != doubled {
			t.Errorf("Reverse(Reverse(%q)) = %q", s, doubled)
		}
		// Property 2: a valid UTF-8 string stays valid after reversing.
		if utf8.ValidString(s) && !utf8.ValidString(rev) {
			t.Errorf("Reverse(%q) produced invalid UTF-8: %q", s, rev)
		}
	})
}
```

You do not assert exact outputs (you cannot know them for random input). Instead
you assert **properties** that must always hold. Run it:

```sh
go test -fuzz=FuzzReverse              # fuzz until it finds a failure or you stop
go test -fuzz=FuzzReverse -fuzztime=30s # fuzz for 30 seconds, then stop
```

When fuzzing finds a crashing input, it **saves that input** under
`testdata/fuzz/` and turns it into a permanent regression test. Without `-fuzz`,
`go test` still runs the seed corpus as ordinary test cases.

> **C vs Go.** In C you bolt on libFuzzer or AFL, write a separate harness, and
> build with special flags. Go fuzzing is just another function in your
> `_test.go` file and one command-line flag.

## Example functions are tests *and* documentation

A function named **`ExampleXxx`** with a trailing `// Output:` comment is compiled,
run, and checked — and it also appears on the package's documentation page. One
piece of code serves as a test and a usage example.

```go
import "fmt"

func ExampleReverse() {
	fmt.Println(Reverse("hello"))
	// Output: olleh
}
```

`go test` runs it and fails if the printed text does not match the `// Output:`
block exactly. If the order of lines is not deterministic, use
`// Unordered output:` instead. An example with no `// Output:` comment is
compiled (so it cannot rot) but not run.

> **Watch out.** The match is **exact**, including whitespace. A stray trailing
> space or a missing newline makes the example fail. This strictness is the point:
> the docs cannot drift from the code.

## Setup and teardown

Most tests need no global setup. When you do, use these tools instead of a
hand-written `main`.

**`TestMain`** runs once for the whole package, around all its tests. Use it for
expensive shared setup (a database container, a temp directory):

```go
func TestMain(m *testing.M) {
	// setup before any test runs
	code := m.Run() // runs all Test/Benchmark/Fuzz functions; returns the exit code
	// teardown after all tests finish
	os.Exit(code)
}
```

**Per-test helpers** clean up automatically — no manual teardown to forget:

```go
func TestSaveConfig(t *testing.T) {
	dir := t.TempDir()                          // temp dir, deleted after the test
	path := filepath.Join(dir, "config.json")
	t.Setenv("APP_ENV", "test")                 // env var, restored after the test
	t.Cleanup(func() { closeSomething() })       // runs at the end, even on failure

	if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	// ... assertions ...
}
```

- `t.TempDir()` returns a unique directory removed when the test ends.
- `t.Setenv(k, v)` sets an environment variable and restores the old value after.
- `t.Cleanup(fn)` registers teardown that runs even if the test fails — like
  `defer`, but tied to the test's lifetime.

## `testdata/` and golden files

The directory named **`testdata/`** is ignored by the Go tools, so it is the
standard home for test fixtures: sample inputs and expected outputs. A **golden
file** holds the expected output of a function; the test compares against it and a
`-update` flag rewrites it when the output legitimately changes.

```go
var update = flag.Bool("update", false, "update golden files")

func TestRender(t *testing.T) {
	got := Render(sampleInput)
	golden := filepath.Join("testdata", "render.golden")

	if *update { // run `go test -update` to regenerate
		if err := os.WriteFile(golden, got, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("output does not match %s (run go test -update to refresh)", golden)
	}
}
```

## Test doubles via interfaces (no mocking framework)

Go does not need a mocking framework. Because interfaces are implicit (Chapter 11
— Interfaces), you write a small fake that satisfies the interface your code
depends on. Accept an interface, pass a fake in the test.

```go
// The code depends on an interface, not a concrete database.
type Store interface {
	Get(id string) (string, error)
}

func Greet(s Store, id string) (string, error) {
	name, err := s.Get(id)
	if err != nil {
		return "", err
	}
	return "Hello, " + name, nil
}
```

```go
// The test supplies a fake Store — no real database needed.
type fakeStore struct{ data map[string]string }

func (f fakeStore) Get(id string) (string, error) {
	v, ok := f.data[id]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func TestGreet(t *testing.T) {
	s := fakeStore{data: map[string]string{"1": "Alice"}}
	got, err := Greet(s, "1")
	if err != nil {
		t.Fatal(err)
	}
	if got != "Hello, Alice" {
		t.Errorf("got %q", got)
	}
}
```

### Testing HTTP with `net/http/httptest`

For HTTP code, `net/http/httptest` spins up a real server on a random local port,
or records a handler's response in memory.

```go
func TestPingServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "pong")
		}))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if strings.TrimSpace(string(body)) != "pong" {
		t.Errorf("body = %q, want pong", body)
	}
}
```

To test a handler without any network, use `httptest.NewRecorder` (it implements
`http.ResponseWriter`) and `httptest.NewRequest`:

```go
func TestPingHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/ping", nil)
	rec := httptest.NewRecorder()

	pingHandler(rec, req) // call the handler directly

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}
```

## Parallel tests

Call `t.Parallel()` at the top of a test to run it alongside other parallel tests.
The runner pauses such tests, then runs them together, which speeds up an I/O-heavy
suite.

```go
func TestSlowA(t *testing.T) {
	t.Parallel()
	// ... independent work ...
}
```

> **Watch out.** Parallel tests run **at the same time**, so they must not share
> mutable state (the same temp file, the same global variable). If two parallel
> tests touch one resource, you get flaky failures or a data race. Give each test
> its own `t.TempDir()` and avoid package-level mutable variables.

## A note on `testify`

The most popular third-party test library is **`testify`**, which adds assertion
helpers (`assert.Equal(t, want, got)`, `require.NoError(t, err)`) and mock
support. It is widely used and fine to adopt. But you do not *need* it: the
standard `if got != want { t.Errorf(...) }` style is idiomatic, dependency-free,
and what the Go standard library itself uses. Learn the stdlib way first.

## Key takeaways

- Testing is built in: `*_test.go` files beside the code, `func TestXxx(t
  *testing.T)`, and `go test`. No framework to install.
- Use `package foo` for white-box tests (see unexported names) or `package
  foo_test` for black-box tests (exported API only).
- `t.Error` records a failure and continues; `t.Fatal` stops the test. Mark
  helpers with `t.Helper()`.
- Prefer **table-driven tests** with `t.Run` subtests: easy to extend, named,
  individually runnable.
- Benchmarks use `func BenchmarkXxx(b *testing.B)` with `for b.Loop()` (Go 1.24);
  fuzz tests use `func FuzzXxx(f *testing.F)` and assert *properties*.
- `ExampleXxx` with `// Output:` is a test and documentation at once.
- `TestMain`, `t.Cleanup`, `t.TempDir`, and `t.Setenv` handle setup/teardown;
  `testdata/` holds fixtures and golden files.
- Fake dependencies with small interface implementations; use `net/http/httptest`
  for HTTP.

## Watch out (gotchas for C programmers)

- **Test results are cached.** Use `-count=1` to force a re-run when a test
  depends on something outside the build.
- **`t.Error` vs `t.Fatal`.** `Fatal` only stops the *current* test/goroutine;
  prefer it after failed setup so later code does not panic.
- **Parallel tests share nothing safely.** `t.Parallel()` tests run concurrently;
  isolate their state or you get flaky races.
- **Example output must match exactly.** Trailing spaces and newlines count.
- **Benchmarks need `-bench`.** A plain `go test` skips them, so a broken
  benchmark can hide until you ask for it.
- **Fuzz corpus failures are saved to `testdata/fuzz/`.** Commit them — they
  become permanent regression tests.

## Interview questions

**Q: How do you write and run a basic test in Go?**
A: Put a function `func TestXxx(t *testing.T)` in a file ending in `_test.go` in
the same package. Inside, call the code, compare the result, and report failures
with `t.Errorf` (continue) or `t.Fatalf` (stop). Run it with `go test ./...`.
There are no assertion macros; you use ordinary `if` statements.

**Q: What is a table-driven test and why is it the preferred style?**
A: It is a slice of test-case structs (inputs and expected outputs) that you loop
over, running each case as a subtest with `t.Run(name, ...)`. It minimizes
duplication, names each case so failures are precise, lets you run one case with
`-run`, and makes adding a case a one-line change.

**Q: What is the difference between `t.Error` and `t.Fatal`?**
A: Both mark the test as failed. `t.Error`/`Errorf` records the failure and lets
the test keep running, so you see all problems. `t.Fatal`/`Fatalf` records the
failure and stops the current test immediately (it calls `runtime.Goexit`). Use
`Fatal` when continuing would panic, such as after a setup error.

**Q: Why might `go test` say "(cached)" and how do you force a re-run?**
A: Go caches test results keyed on the test binary and its inputs; if nothing
changed it reuses the previous result. To force a real run, pass `-count=1`, which
disables the cache. This matters for tests that read a clock, the filesystem, or
the network.

**Q: How does Go's fuzzing work, and what do you assert in a fuzz test?**
A: You write `func FuzzXxx(f *testing.F)`, add seed inputs with `f.Add`, and call
`f.Fuzz` with a function that takes `*testing.T` plus the fuzzed arguments. Go
generates many inputs and runs the body. You assert *properties* that must always
hold (for example, reversing twice yields the original), not exact outputs.
Crashing inputs are saved as regression tests.

**Q: How do you test code that depends on a database or an HTTP service without
the real thing?**
A: Depend on an interface, then pass a small fake that implements it in the test —
Go's implicit interfaces make this trivial, so no mocking framework is required.
For HTTP, use `net/http/httptest` to run a real local server or to record a
handler's response with `httptest.NewRecorder`.

## Try it

1. Create `strutil/reverse.go` with the `Reverse` function and a
   `reverse_test.go` with the table-driven `TestReverse`. Run `go test -v`.
2. Add `BenchmarkReverse` using `for b.Loop()` and run `go test -bench=. -benchmem`.
3. Add `FuzzReverse` and run `go test -fuzz=FuzzReverse -fuzztime=15s`. Try
   introducing a bug (reverse bytes instead of runes) and watch the fuzzer catch a
   multi-byte UTF-8 input.
