// Project 1: a small `wc`-like counter.
// Counts lines, words, and bytes for each file given on the command line,
// or for standard input when no files are given.
//
// Flags:
//   -l   print the line count
//   -w   print the word count
//   -c   print the byte count
// With no flags, all three are printed (like `wc`).

use std::env;
use std::fs::File;
use std::io::{self, BufRead, BufReader, Write};
use std::process::ExitCode;

/// The three counts we collect for one input.
#[derive(Debug, Default, Clone, Copy)]
struct Counts {
    lines: u64,
    words: u64,
    bytes: u64,
}

impl Counts {
    /// Add another set of counts into this one (used for the total line).
    fn add(&mut self, other: &Counts) {
        self.lines += other.lines;
        self.words += other.words;
        self.bytes += other.bytes;
    }
}

/// Which columns the user asked for.
#[derive(Debug, Clone, Copy)]
struct Flags {
    lines: bool,
    words: bool,
    bytes: bool,
}

impl Flags {
    /// When the user gives no flag at all, show every column.
    fn or_default(self) -> Flags {
        if !self.lines && !self.words && !self.bytes {
            Flags {
                lines: true,
                words: true,
                bytes: true,
            }
        } else {
            self
        }
    }
}

/// Read everything from `reader` and count lines, words, and bytes.
fn count<R: BufRead>(mut reader: R) -> io::Result<Counts> {
    let mut counts = Counts::default();
    let mut line = Vec::new();

    // read_until keeps the trailing '\n', so byte counts match `wc`.
    loop {
        line.clear();
        let n = reader.read_until(b'\n', &mut line)?;
        if n == 0 {
            break; // end of input
        }
        counts.bytes += n as u64;
        if line.ends_with(b"\n") {
            counts.lines += 1;
        }
        // Words are runs of non-whitespace bytes, split on ASCII whitespace.
        let text = String::from_utf8_lossy(&line);
        counts.words += text.split_whitespace().count() as u64;
    }
    Ok(counts)
}

/// Print one result row, only the columns the flags asked for.
fn print_row(out: &mut impl Write, c: &Counts, flags: Flags, label: &str) -> io::Result<()> {
    if flags.lines {
        write!(out, "{:>8}", c.lines)?;
    }
    if flags.words {
        write!(out, "{:>8}", c.words)?;
    }
    if flags.bytes {
        write!(out, "{:>8}", c.bytes)?;
    }
    if label.is_empty() {
        writeln!(out)
    } else {
        writeln!(out, " {label}")
    }
}

/// Parse arguments into flags and the list of file paths.
fn parse_args(args: Vec<String>) -> Result<(Flags, Vec<String>), String> {
    let mut flags = Flags {
        lines: false,
        words: false,
        bytes: false,
    };
    let mut paths = Vec::new();

    for arg in args {
        if arg == "--" {
            continue;
        }
        if arg.starts_with('-') && arg.len() > 1 {
            for ch in arg.chars().skip(1) {
                match ch {
                    'l' => flags.lines = true,
                    'w' => flags.words = true,
                    'c' => flags.bytes = true,
                    other => return Err(format!("unknown flag: -{other}")),
                }
            }
        } else {
            paths.push(arg);
        }
    }
    Ok((flags, paths))
}

fn run() -> Result<(), String> {
    let raw: Vec<String> = env::args().skip(1).collect();
    let (flags, paths) = parse_args(raw)?;
    let flags = flags.or_default();

    let stdout = io::stdout();
    let mut out = stdout.lock();

    if paths.is_empty() {
        // No files: read standard input.
        let stdin = io::stdin();
        let counts = count(stdin.lock()).map_err(|e| format!("stdin: {e}"))?;
        print_row(&mut out, &counts, flags, "").map_err(|e| e.to_string())?;
        return Ok(());
    }

    let mut total = Counts::default();
    let mut had_error = false;
    let many = paths.len() > 1;

    for path in &paths {
        match open_and_count(path) {
            Ok(counts) => {
                total.add(&counts);
                print_row(&mut out, &counts, flags, path).map_err(|e| e.to_string())?;
            }
            Err(e) => {
                // Report the error but keep going with the other files.
                eprintln!("linecount: {path}: {e}");
                had_error = true;
            }
        }
    }

    if many {
        print_row(&mut out, &total, flags, "total").map_err(|e| e.to_string())?;
    }

    if had_error {
        Err("one or more files could not be read".to_string())
    } else {
        Ok(())
    }
}

/// Open a file and count it. A helper so `?` can do the error plumbing.
fn open_and_count(path: &str) -> io::Result<Counts> {
    let file = File::open(path)?;
    // BufReader avoids one syscall per byte, like a stdio FILE* buffer in C.
    count(BufReader::new(file))
}

fn main() -> ExitCode {
    match run() {
        Ok(()) => ExitCode::SUCCESS,
        Err(msg) => {
            eprintln!("linecount: {msg}");
            ExitCode::FAILURE
        }
    }
}
