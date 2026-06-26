// Command linecount is a small wc-like tool. It counts lines, words, and
// bytes from files named on the command line, or from standard input when no
// files are given.
//
// Usage:
//
//	linecount [-l] [-w] [-c] [file ...]
//
// With no flags it prints lines, words, and bytes, like wc. With one or more of
// -l, -w, -c it prints only the requested counts, in that fixed order.
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// counts holds the three totals for one input.
type counts struct {
	lines int64
	words int64
	bytes int64
}

// add accumulates another set of counts into c. It is used to build the grand
// total across several files.
func (c *counts) add(o counts) {
	c.lines += o.lines
	c.words += o.words
	c.bytes += o.bytes
}

// selection records which counts the user asked to see.
type selection struct {
	lines, words, bytes bool
}

// isSpace reports whether b is an ASCII whitespace byte. This is the same set
// C's isspace recognizes in the default "C" locale: space, tab, newline,
// vertical tab, form feed, and carriage return.
func isSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\v', '\f', '\r':
		return true
	default:
		return false
	}
}

// count reads everything from r and returns the line, word, and byte totals.
// A line is counted at each newline byte, and a word is a run of non-space
// bytes, exactly like wc. bufio.Reader buffers the underlying reads, so the
// byte-at-a-time loop is still fast.
func count(r io.Reader) (counts, error) {
	br := bufio.NewReader(r)
	var c counts
	inWord := false
	for {
		b, err := br.ReadByte()
		if err != nil {
			// io.EOF is the normal, expected end of input, not a failure.
			if errors.Is(err, io.EOF) {
				return c, nil
			}
			return c, err
		}
		c.bytes++
		if b == '\n' {
			c.lines++
		}
		if isSpace(b) {
			inWord = false
		} else if !inWord {
			inWord = true
			c.words++
		}
	}
}

// countFile opens name, counts it, and always closes the file.
func countFile(name string) (counts, error) {
	f, err := os.Open(name)
	if err != nil {
		// os.Open returns a *fs.PathError that already includes the file name
		// and the reason ("open foo: no such file or directory"), so we return
		// it unchanged.
		return counts{}, err
	}
	defer f.Close()
	return count(f)
}

// format renders the selected counts followed by an optional name. Counts are
// right-aligned in a fixed width so columns line up, similar to wc.
func format(c counts, sel selection, name string) string {
	var sb strings.Builder
	writeField := func(n int64) {
		if sb.Len() > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%7d", n)
	}
	if sel.lines {
		writeField(c.lines)
	}
	if sel.words {
		writeField(c.words)
	}
	if sel.bytes {
		writeField(c.bytes)
	}
	if name != "" {
		sb.WriteByte(' ')
		sb.WriteString(name)
	}
	return sb.String()
}

func main() {
	lineFlag := flag.Bool("l", false, "print the line count")
	wordFlag := flag.Bool("w", false, "print the word count")
	byteFlag := flag.Bool("c", false, "print the byte count")
	flag.Parse()

	// If the user named no counters, show all three (the wc default).
	if !*lineFlag && !*wordFlag && !*byteFlag {
		*lineFlag, *wordFlag, *byteFlag = true, true, true
	}
	sel := selection{lines: *lineFlag, words: *wordFlag, bytes: *byteFlag}

	files := flag.Args()

	// No file arguments: read standard input and print only the counts.
	if len(files) == 0 {
		c, err := count(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "linecount: stdin: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(format(c, sel, ""))
		return
	}

	// One or more files: count each, print its line, and keep a running total.
	// A bad file is reported on stderr; we continue and exit non-zero at the end.
	exit := 0
	var total counts
	for _, name := range files {
		c, err := countFile(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "linecount: %v\n", err)
			exit = 1
			continue
		}
		total.add(c)
		fmt.Println(format(c, sel, name))
	}
	if len(files) > 1 {
		fmt.Println(format(total, sel, "total"))
	}
	os.Exit(exit)
}
