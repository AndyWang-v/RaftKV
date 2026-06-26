# References and Further Reading

Where to go next. These are the resources this book draws on and the ones most useful
to a C programmer learning Go. Links are to the official or canonical page; where no
stable link exists, the resource is named so you can search for it.

## Official documentation

- **The Go website** — https://go.dev — the home for downloads, documentation, and the
  blog. Start here.
- **A Tour of Go** — https://go.dev/tour/ — an interactive, in-browser introduction.
  The fastest way to get hands-on with the syntax.
- **Effective Go** — https://go.dev/doc/effective_go — how to write idiomatic Go.
  Essential reading once the syntax feels comfortable.
- **The Go Programming Language Specification** — https://go.dev/ref/spec — the precise
  language definition. Short and readable, like the C standard but far smaller.
- **Standard library reference (pkg.go.dev)** — https://pkg.go.dev — searchable docs for
  the standard library and every public module. Your daily reference.
- **The Go Blog** — https://go.dev/blog/ — official articles on language features,
  releases, and design decisions.
- **Go Code Review Comments** — https://go.dev/wiki/CodeReviewComments — a checklist of
  common review feedback. A concise list of "how Go is supposed to look."
- **Google Go Style Guide** — https://google.github.io/styleguide/go/ — Google's
  internal style guide, made public. Opinionated and practical.
- **The Go Memory Model** — https://go.dev/ref/mem — the exact rules for when one
  goroutine's writes are visible to another. The Go equivalent of the C11 memory model.
- **The Go FAQ** — https://go.dev/doc/faq — answers to "why does Go do it this way?",
  including many questions a C or C++ programmer will ask.

## Books

- **The Go Programming Language**, Alan A. A. Donovan and Brian W. Kernighan —
  https://www.gopl.io — the definitive book, often called "the K&R of Go." Kernighan
  co-wrote *The C Programming Language*, so the style will feel familiar.
- **Learning Go**, Jon Bodner (O'Reilly) — a thorough, modern introduction that covers
  generics and recent versions. A good cover-to-cover first book.
- **Concurrency in Go**, Katherine Cox-Buday (O'Reilly) — a focused, deep treatment of
  goroutines, channels, and concurrency patterns.
- **100 Go Mistakes and How to Avoid Them**, Teiva Harsanyi (Manning) — a catalog of
  real traps and their fixes. Excellent for moving past "writing C in Go."

## Practice

- **Go by Example** — https://gobyexample.com — short, annotated programs for each
  language feature. Great for "show me the code" lookups.
- **Exercism Go track** — https://exercism.org/tracks/go — graded exercises with human
  mentorship, free.
- **Gophercises** — https://gophercises.com — small, practical coding projects that
  build real skills.

## Tools

- **gopls** — https://pkg.go.dev/golang.org/x/tools/gopls — the official language server
  that powers editor features (completion, jump-to-definition, refactoring).
- **staticcheck** — https://staticcheck.dev — a powerful static analyzer that catches
  bugs and suspicious code beyond `go vet`.
- **golangci-lint** — https://golangci-lint.run — a fast runner that bundles many
  linters with one config. The common choice for CI.
- **Delve** — https://github.com/go-delve/delve — the Go debugger (`dlv`), with editor
  integrations. The rough equivalent of `gdb` for Go.
- **govulncheck** — https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck — scans your
  code and dependencies for known vulnerabilities, reporting only those you actually
  call.

## Community

- **r/golang** — https://www.reddit.com/r/golang/ — active discussion, news, and
  questions.
- **Gophers Slack** — https://invite.slack.golangbridge.org — a large community
  workspace with channels for beginners, specific packages, and more.
- **Go Time podcast** — https://changelog.com/gotime — interviews and discussion on Go
  practice, tooling, and community.

## Especially for C programmers

These four pieces explain the parts of Go that surprise C programmers most. Read them
early.

- **Go Slices: usage and internals** — https://go.dev/blog/slices-intro — what a slice
  is and how it sits on top of an array. Directly addresses the C "array vs pointer"
  mental model.
- **Arrays, slices (and strings): the mechanics of 'append'** —
  https://go.dev/blog/slices — exactly how `append` grows a slice and why the aliasing
  traps happen. Read it before you trust `append`.
- **Share Memory By Communicating** — https://go.dev/blog/codelab-share — the idea
  behind channels, contrasted with the shared-memory-plus-locks model you know from
  pthreads.
- **Effective Go** — https://go.dev/doc/effective_go — the single best document for
  turning correct-but-C-like code into idiomatic Go.
