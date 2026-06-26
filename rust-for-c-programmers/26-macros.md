# Chapter 26 — Macros

> **What you'll learn.** What macros are and why Rust has them, how
> `macro_rules!` writes code from patterns, why the macros you already use
> (`println!`, `vec!`, `assert!`) *must* be macros, and a tour of procedural
> macros like `#[derive(...)]`. Throughout, how all of this beats the C
> preprocessor.

## What a macro is, and why a language needs one

A **macro** is code that writes code. The compiler runs the macro at compile
time, the macro produces more source code, and that produced code is then
compiled as if you had typed it yourself. This is called **metaprogramming**:
a program that generates a program.

You already know one macro system: the **C preprocessor**. In C you write
`#define`, `#include`, and `#ifdef`, and a separate pass edits your text before
the real compiler ever sees it.

Why have macros at all? A normal function cannot do some things:

- **Variadic, type-checked interfaces.** `println!("{} {}", a, b)` takes any
  number of arguments of different types and checks the format string against
  them *at compile time*. A normal Rust function cannot take a variable number
  of arguments like C's `printf`.
- **New syntax (small DSLs).** A **DSL** (domain-specific language) is a mini
  language for one job. `vec![1, 2, 3]` and `vec![0; 100]` are tiny DSLs for
  building a vector.
- **Removing boilerplate.** Generating repetitive code (trait implementations,
  table entries) without copy-paste.

> **C vs Rust.** The C preprocessor does **raw text substitution**. It does not
> understand C: not types, not scopes, not even that your tokens form valid
> expressions. Rust macros are different. They operate on **structured tokens**
> (the compiler's view of the code), they respect scope, and they are
> **hygienic** (defined below). They cannot silently corrupt your program the
> way a careless `#define` can.

Rust has two kinds of macro:

- **Declarative macros**, written with `macro_rules!`. You match patterns of
  tokens and say what to expand them into. This is most of what you write.
- **Procedural macros**, which are Rust functions that take tokens and return
  tokens. They include `#[derive(...)]`, attribute macros, and function-like
  macros. They live in their own special crate.

## The C preprocessor, and its traps

First, a refresher on why C macros are dangerous, so you can see what Rust
fixes. Here is the classic broken `MAX`:

```c
#define MAX(a, b) ((a) > (b) ? (a) : (b))

int main(void) {
    int i = 5;
    int m = MAX(i++, 3);   /* i++ is evaluated TWICE: i becomes 7, not 6 */
    return m;
}
```

The preprocessor pasted `i++` into the body twice. There is no concept of
"evaluate the argument once" because there is no concept of an *expression* —
only text. Other classic traps:

```c
#define SQUARE(x) x * x
int a = SQUARE(2 + 3);     /* expands to 2 + 3 * 2 + 3 = 11, not 25 */

#define DOUBLE(x) ((x) + (x))
int n = 0;
int d = DOUBLE(n++);       /* n++ twice again: undefined behavior */
```

And the **hygiene** problem. **Hygiene** means a macro's own names cannot
accidentally collide with names at the call site. C macros are *not* hygienic:

```c
#define SWAP(a, b) { int tmp = a; a = b; b = tmp; }

int main(void) {
    int x = 1, tmp = 2;
    SWAP(x, tmp);          /* the macro's `tmp` clashes with your `tmp` */
    return x;              /* broken: the names collide */
}
```

Rust macros do not have any of these problems, as we will now see.

## Declarative macros with `macro_rules!`

A declarative macro looks like a `match`: you write **rules**, each with a
**matcher** (a pattern of tokens) on the left and a **transcriber** (the code to
produce) on the right.

Here is the smallest useful example — `square!`, which evaluates its argument
exactly once and is correctly parenthesized:

```rust
macro_rules! square {
    ($x:expr) => {
        $x * $x
    };
}

fn main() {
    let n = square!(2 + 3);   // expands to (2 + 3) * (2 + 3) = 25, correctly
    println!("{n}");
}
```

`$x:expr` is a **metavariable**. The `$x` is the name; `expr` is a **fragment
specifier** that says "match a whole expression here." Because Rust captured a
real expression — not raw text — `2 + 3` stays together as one unit. The
`SQUARE` bug from C cannot happen.

> **C vs Rust.** A C macro argument is text and you must defensively wrap it in
> parentheses. A Rust `$x:expr` is a complete, parsed expression. It already
> behaves as one value, so `square!(2 + 3)` is never re-parsed wrongly.

### Fragment specifiers (the `:kind` part)

The fragment specifier tells the macro what kind of token to match. The common
ones:

| Specifier | Matches | Example match |
|---|---|---|
| `expr` | an expression | `2 + 3`, `foo()`, `x` |
| `ident` | an identifier (a name) | `count`, `Widget` |
| `ty` | a type | `i32`, `Vec<u8>` |
| `literal` | a literal value | `42`, `"hi"`, `true` |
| `pat` | a pattern (for `match`/`let`) | `Some(x)`, `_` |
| `block` | a `{ ... }` block | `{ do_it(); }` |
| `stmt` | a statement | `let y = 1` |
| `path` | a path | `std::mem::swap` |
| `tt` | a single token tree (most flexible) | almost anything |

### Repetition: matching many arguments

To accept a variable number of arguments, you use a **repetition**. The syntax
is `$( ... )sep rep`, where `sep` is an optional separator and `rep` is `*`
(zero or more), `+` (one or more), or `?` (zero or one).

`$( $x:expr ),*` means "zero or more expressions, separated by commas." You then
repeat the expansion with the same `$( ... )*` shape. Here is a small `my_vec!`
that mimics the real `vec!`:

```rust
macro_rules! my_vec {
    // empty: my_vec![]
    () => {
        Vec::new()
    };
    // one or more elements: my_vec![a, b, c]
    ( $( $x:expr ),+ $(,)? ) => {{
        let mut v = Vec::new();
        $(
            v.push($x);
        )+
        v
    }};
}

fn main() {
    let a: Vec<i32> = my_vec![];
    let b = my_vec![1, 2, 3];
    let c = my_vec![10, 20,];     // trailing comma allowed by $(,)?
    println!("{a:?} {b:?} {c:?}");
}
```

Three things to notice:

- The matcher `$( $x:expr ),+` captures the list. The transcriber
  `$( v.push($x); )+` repeats once per captured expression. The repetition on
  the right must use the same metavariable as the left.
- `$(,)?` allows an optional trailing comma, like `vec![1, 2, 3,]`.
- The body uses **double braces** `{{ ... }}`. The outer braces are the macro's
  transcriber; the inner braces make a **block expression** so the macro expands
  to a single value (the `v` on the last line).

A `hashmap!` macro follows the same shape and shows a `key => value` separator:

```rust
use std::collections::HashMap;

macro_rules! hashmap {
    ( $( $key:expr => $val:expr ),* $(,)? ) => {{
        let mut m = HashMap::new();
        $(
            m.insert($key, $val);
        )*
        m
    }};
}

fn main() {
    let scores = hashmap!{
        "alice" => 10,
        "bob"   => 7,
    };
    println!("{:?}", scores.get("alice"));
}
```

### Hygiene: the macro's own names do not leak

Remember the broken C `SWAP` where `tmp` collided. Watch the Rust version. The
macro introduces a variable called `tmp`, and the caller *also* has a `tmp`.
They do **not** clash:

```rust
macro_rules! swap {
    ($a:expr, $b:expr) => {{
        let tmp = $a;
        $a = $b;
        $b = tmp;
    }};
}

fn main() {
    let mut x = 1;
    let mut tmp = 2;       // same name the macro uses internally
    swap!(x, tmp);         // works correctly; no collision
    println!("x={x} tmp={tmp}");   // x=2 tmp=1
}
```

The `tmp` inside the macro and the `tmp` in `main` are treated as different
names, because the compiler tracks where each identifier came from. This is
**hygiene**, and C macros do not have it.

> **Mental model.** Think of each macro expansion as getting its own private
> namespace painted a different color. The macro's `tmp` is "blue"; your `tmp`
> is "red." They never mix, even though they are spelled the same.

> **Watch out.** Hygiene is helpful but it has an edge: a macro normally cannot
> *invent* a name and hand it back for you to use, because that name is in the
> macro's color, not yours. If you need a name to cross the boundary, the caller
> must pass it in as an `ident`.

### Where macros are visible

A `macro_rules!` macro is in scope **after** the point where it is defined
(textually), which is unlike normal items. To use it across modules or export
it from a library, mark it `#[macro_export]`:

```rust
#[macro_export]
macro_rules! my_vec {
    // ...
    () => { Vec::new() };
}
```

This puts the macro at the crate root so other crates can `use your_crate::my_vec;`.

## The macros you already use (and why they must be macros)

You have been using macros since Chapter 1. Now you can see *why* each one is a
macro and not a function.

| Macro | Why it cannot be a plain function |
|---|---|
| `println!` / `print!` | Variable number of arguments; checks the format string against the arguments at compile time. |
| `format!` | Same as `println!`, but builds a `String` instead of printing. |
| `vec!` | Takes any number of elements, or the `vec![x; n]` form; builds a `Vec`. |
| `assert!` / `assert_eq!` | On failure, prints the **source text** of the failed expression and the file/line — only a macro can see that text. |
| `panic!` | Variadic, format-string aware, and records the call location. |
| `dbg!` | Prints the expression's source text *and* its value, then returns the value. |

The format-checking point is worth stressing. This is caught at **compile
time**, before the program runs:

```rust
// COMPILE ERROR: 1 positional argument in format string, but no arguments were given
fn main() {
    let _ = format!("value = {}");   // error: missing argument for `{}`
}
```

A C `printf` cannot do this in general; a wrong format specifier is at best a
warning and at worst undefined behavior at runtime. The Rust macro turns it into
a hard compile error.

```c
/* C: this compiles; it is undefined behavior at runtime */
printf("%d\n", "not a number");   /* wrong type for %d */
```

## Procedural macros (an overview — you will mostly *use* them)

**Procedural macros** ("proc macros") are the more powerful kind. A proc macro
is a Rust function that receives a stream of tokens and returns a new stream of
tokens. Because it is ordinary Rust code, it can do arbitrary work — parse the
input, inspect a struct's fields, and build new code.

Three forms exist:

- **Derive macros** — `#[derive(Debug, Clone, Serialize)]`. They generate a
  trait implementation for the type they are attached to. This is by far the
  most common form you will meet.
- **Attribute macros** — `#[tokio::main]`, `#[test]`, web route attributes like
  `#[get("/users")]`. They wrap or transform the item they decorate.
- **Function-like macros** — `sql!(SELECT * FROM users)`. They look like
  `macro_rules!` calls but are backed by a proc-macro function.

You will **use** proc macros far more often than you write them. The everyday
example is `serde`, which turns a plain struct into something that can be
converted to and from JSON with one line:

```rust
// needs: cargo add serde --features derive   and   cargo add serde_json
use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Debug)]
struct Config {
    name: String,
    retries: u32,
}

fn main() {
    let c = Config { name: "db".into(), retries: 3 };
    let json = serde_json::to_string(&c).unwrap();   // real code handles the error
    println!("{json}");                              // {"name":"db","retries":3}
}
```

The `#[derive(Serialize, Deserialize)]` line ran a proc macro that generated all
the serialization code for `Config` at compile time. You wrote none of it.

### How proc macros are built (just so you recognize it)

Two crates do the heavy lifting, and you will see them in any proc-macro
project:

- **`syn`** parses the incoming tokens into a syntax tree you can inspect.
- **`quote`** builds the outgoing tokens from a template, with `#var`
  interpolation.

A proc macro must live in **its own crate**, marked in `Cargo.toml`:

```toml
[lib]
proc-macro = true
```

This separate-crate rule exists because a proc macro is compiled and run *by the
compiler* while it compiles your program. It is a compiler plugin, so it cannot
sit in the same crate as the code that uses it.

> **Deep dive.** A minimal derive macro signature looks like this. You do not
> need to memorize it; just recognize the shape.
>
> ```rust
> // in a crate with proc-macro = true
> use proc_macro::TokenStream;
>
> #[proc_macro_derive(Hello)]
> pub fn derive_hello(input: TokenStream) -> TokenStream {
>     // parse `input` with syn, build output with quote, return tokens
>     TokenStream::new()   // (a real one returns generated code)
> }
> ```

## C macros vs Rust macros, side by side

| Feature | C preprocessor | `macro_rules!` | Rust proc macros |
|---|---|---|---|
| Operates on | raw text | structured tokens | structured tokens / syntax tree |
| Hygienic | no | yes | mostly (you control spans) |
| Understands types | no | partly (fragments) | yes (can inspect a struct) |
| Argument evaluated once | no (you must be careful) | yes (`$x:expr` is one value) | yes |
| Scoped / namespaced | no (global text) | yes (module-scoped) | yes |
| Can read a struct's fields | no | no | yes (derive macros do) |
| Where it lives | inline `#define` | inline or `#[macro_export]` | its own `proc-macro` crate |
| Debuggable | `cpp -E` shows text | `cargo expand` | `cargo expand` |

### The X-macro: a C trick Rust does better

C programmers use the **X-macro** pattern to keep a list and its generated code
in sync — for example, an enum and its name table:

```c
/* the single list */
#define COLORS  X(RED) X(GREEN) X(BLUE)

enum Color {
#define X(name) name,
    COLORS
#undef X
};

const char *color_name(enum Color c) {
    switch (c) {
#define X(name) case name: return #name;
        COLORS
#undef X
    }
    return "?";
}
```

It works, but it is fragile text trickery with `#define`/`#undef` dances. In
Rust you write a `macro_rules!` macro that builds both the enum and the name
function from one list — hygienically, with no `#undef`:

```rust
macro_rules! colors {
    ( $( $name:ident ),+ $(,)? ) => {
        #[derive(Debug, Clone, Copy)]
        enum Color { $( $name ),+ }

        impl Color {
            fn name(self) -> &'static str {
                match self {
                    $( Color::$name => stringify!($name) ),+
                }
            }
        }
    };
}

colors!(Red, Green, Blue);

fn main() {
    println!("{}", Color::Green.name());   // "Green"
}
```

`stringify!` is a built-in macro that turns tokens into a string literal — the
clean equivalent of C's `#name` stringizing operator. For a struct-by-struct
job (like generating one trait impl per type), a proc macro is even better,
because it can read the type's real fields.

## Seeing the expansion: `cargo expand`

Macros expand before type checking, so when one misbehaves the error can point
at generated code you never wrote. The fix is to **look at the expansion**. The
`cargo expand` tool (Chapter 24 — Tooling) prints your code with every macro
already expanded:

```sh
cargo install cargo-expand     # one-time
cargo expand                   # prints the whole crate, macros expanded
cargo expand main              # just one module/function
```

Run it on the `my_vec!` or `colors!` example above and you will see the plain
code the macro produced. This is your `cpp -E` for Rust, but it shows real
expanded Rust, not raw text.

> **Rule of thumb.** Reach for a macro only when a function or a generic cannot
> do the job. Macros are harder to read, harder to debug, and worse for IDE
> support. Variable argument counts, new syntax, and per-type code generation
> are good reasons. "I want a default argument" or "I want to avoid typing a
> generic" usually are not.

## Key takeaways

- A macro is code that writes code at compile time. Rust has **declarative**
  macros (`macro_rules!`) and **procedural** macros.
- Unlike C's text-substituting preprocessor, Rust macros work on **structured
  tokens**, are **hygienic** (their names do not collide with yours), and
  capture whole expressions so an argument is evaluated as one unit.
- `macro_rules!` matches **fragments** (`$x:expr`, `$x:ident`, `$x:ty`, ...) and
  supports **repetition** (`$( ... ),*`). Use double braces `{{ }}` to expand to
  a single block expression.
- `println!`, `format!`, `vec!`, `assert!`, `panic!`, and `dbg!` are macros
  because they are variadic and/or check the format string at compile time.
- Procedural macros power `#[derive(...)]`, attribute macros (`#[tokio::main]`),
  and function-like macros. They live in their own `proc-macro` crate and are
  built with `syn` and `quote`. You will mostly *use* them (e.g. serde).
- Use `cargo expand` to see what a macro produced when debugging.

## Watch out (gotchas for C programmers)

- **Hygiene is not C's behavior.** A Rust macro's internal `tmp` will not clash
  with yours, and a macro cannot leak a new variable name to the caller unless
  the caller passes the name in as an `ident`.
- **Macros expand before type checking.** A macro can produce code that *looks*
  fine but fails to type-check later, with errors that point at generated code.
  Run `cargo expand` to see the real source.
- **Procedural macros need a separate crate** (`proc-macro = true`). You cannot
  define a derive macro in the same crate that uses it.
- **`macro_rules!` is order-sensitive.** A macro is only usable after its
  definition, unlike normal functions and types. Use `#[macro_export]` to share
  it.
- **Prefer functions and generics when they suffice.** A macro is a last resort,
  not a first reach. Macros hurt readability, debugging, and tooling.
- **The `!` is part of the call.** `vec!`, `println!`, and friends always need
  the `!`; without it the compiler looks for a function of that name.

## Interview questions

**Q: How do Rust macros differ from C preprocessor macros?**
A: The C preprocessor does raw text substitution and understands nothing about
the language — not types, not scopes, not expressions. Rust macros operate on
structured tokens, are hygienic (their identifiers do not collide with the
caller's), and capture whole fragments such as expressions, so an argument is
evaluated once rather than pasted as text. This removes the classic C macro bugs
like double evaluation and name capture.

**Q: What does "hygienic" mean for a macro?**
A: It means names introduced inside the macro live in their own scope and cannot
accidentally collide with names at the call site. In C, a macro that declares
`int tmp;` breaks if the caller also has a `tmp`. In Rust the two `tmp`s are
treated as distinct, so the macro is safe to use anywhere.

**Q: What is the difference between a declarative macro and a procedural macro?**
A: A declarative macro (`macro_rules!`) matches token patterns and expands to
new code; it is good for variadic syntax and small DSLs. A procedural macro is a
Rust function that takes tokens and returns tokens; it can inspect a type's
structure and generate code from it. Proc macros power `#[derive(...)]`,
attribute macros, and function-like macros, and they live in their own
`proc-macro` crate.

**Q: Why is `println!` a macro instead of a function?**
A: It takes a variable number of arguments of different types and checks the
format string against those arguments at compile time. A normal Rust function
cannot be variadic like that, and only a macro can validate the format string
before the program runs, turning what would be a C runtime bug into a compile
error.

**Q: A macro you wrote produces a confusing type error. How do you debug it?**
A: Run `cargo expand`, which prints your source with all macros already
expanded. You can then read the generated Rust directly and see where it goes
wrong, since macros expand before type checking and the error often points at
code you did not type by hand.

## Try it

1. Write the `my_vec!` macro from this chapter and run `cargo expand` to see the
   `push` calls it generates.
2. Add a rule to `my_vec!` for the `vec![x; n]` form (hint: match
   `$elem:expr ; $count:expr` and call a loop or `vec!` internally).
3. Write the broken C `SQUARE(x) x * x` in your head, then confirm the Rust
   `square!($x:expr)` version returns 25 for `square!(2 + 3)`. Change `$x:expr`
   thinking about why text substitution would have failed.
