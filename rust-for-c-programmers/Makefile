# Makefile for rendering "Rust for C Programmers" from Markdown.
#
# The Markdown files are the single source of truth. This Makefile only renders
# them into a print-ready artifact with pandoc.
#
#   make html   -> self-contained HTML (open in a browser; Print > Save as PDF)
#   make epub   -> e-reader friendly file
#   make pdf    -> polished PDF (needs a LaTeX engine: `brew install tectonic`)
#
# DIAGRAMS: chapters use ```mermaid fenced blocks. GitHub, VS Code, and many
# Markdown viewers render these automatically, so you need NOTHING extra to read
# the book on screen. To turn Mermaid into pictures inside the PDF/HTML/EPUB,
# install the filter once:
#
#   make mermaid-deps      # installs mermaid-filter locally via npm (needs Node)
#
# If the filter is not installed, the book still renders fine; the Mermaid blocks
# just appear as plain code listings (every Mermaid diagram is paired with an
# ASCII diagram or prose, so no information is lost).

PANDOC ?= pandoc

# Chapters in reading order. Keep this list in sync with README.md.
CHAPTERS = \
	00-preface.md \
	01-why-rust.md \
	02-toolchain-and-cargo.md \
	03-program-structure.md \
	04-variables-types-constants.md \
	05-control-flow.md \
	06-functions-and-closures.md \
	07-ownership-and-moves.md \
	08-borrowing-and-references.md \
	09-lifetimes.md \
	10-slices-and-strings.md \
	11-structs-and-methods.md \
	12-enums-and-pattern-matching.md \
	13-error-handling.md \
	14-generics.md \
	15-traits.md \
	16-collections-and-iterators.md \
	17-smart-pointers.md \
	18-memory-without-gc.md \
	19-threads-and-concurrency.md \
	20-channels-and-shared-state.md \
	21-async-await.md \
	22-cargo-crates-and-workspaces.md \
	23-testing.md \
	24-tooling.md \
	25-unsafe-and-ffi.md \
	26-macros.md \
	27-idioms-and-style.md \
	28-gotchas-for-c-programmers.md \
	29-interview-questions.md \
	30-capstone-projects.md \
	appendix-c-to-rust-cheatsheet.md \
	appendix-glossary.md \
	references.md

META = metadata.yaml
OUT  = rust-for-c-programmers

# Use the mermaid filter only if it is available (on PATH or installed locally).
MERMAID_BIN := $(shell command -v mermaid-filter 2>/dev/null || ([ -x node_modules/.bin/mermaid-filter ] && echo node_modules/.bin/mermaid-filter))
FILTERS     := $(if $(MERMAID_BIN),-F $(MERMAID_BIN),)

.PHONY: all pdf epub html clean mermaid-deps

all: html

pdf: $(META) $(CHAPTERS)
	$(PANDOC) $(FILTERS) $(META) $(CHAPTERS) -o $(OUT).pdf --pdf-engine=tectonic

epub: $(META) $(CHAPTERS)
	$(PANDOC) $(FILTERS) $(META) $(CHAPTERS) -o $(OUT).epub

html: $(META) $(CHAPTERS)
	$(PANDOC) $(FILTERS) $(META) $(CHAPTERS) -o $(OUT).html --standalone --embed-resources --toc

# One-time: install the Mermaid->image filter locally (needs Node + npm).
mermaid-deps:
	npm install mermaid-filter

clean:
	rm -f $(OUT).pdf $(OUT).epub $(OUT).html
