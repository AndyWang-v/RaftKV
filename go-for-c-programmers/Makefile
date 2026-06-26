# Makefile for rendering "Go for C Programmers" from Markdown.
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
	01-why-go.md \
	02-toolchain-and-commands.md \
	03-program-structure.md \
	04-types-variables-constants.md \
	05-control-flow.md \
	06-functions.md \
	07-pointers.md \
	08-arrays-slices-strings.md \
	09-maps.md \
	10-structs-and-methods.md \
	11-interfaces.md \
	12-errors.md \
	13-goroutines-scheduler.md \
	14-channels-select.md \
	15-sync-and-context.md \
	16-concurrency-patterns.md \
	17-memory-and-gc.md \
	18-packages-and-modules.md \
	19-generics.md \
	20-standard-library-tour.md \
	21-testing.md \
	22-tooling.md \
	23-web-services.md \
	24-cli-tools.md \
	25-idioms-and-style.md \
	26-gotchas-for-c-programmers.md \
	27-interview-questions.md \
	28-capstone-projects.md \
	appendix-c-to-go-cheatsheet.md \
	appendix-glossary.md \
	references.md

META = metadata.yaml
OUT  = go-for-c-programmers

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
