# Makefile for rendering the book from Markdown.
#
# The Markdown files are the single source of truth. This Makefile just renders
# them into a print-ready artifact with pandoc (already installed).
#
# PDF output additionally needs a LaTeX engine. The lightest modern option:
#     brew install tectonic
# Then `make pdf` works. EPUB and HTML need only pandoc (no extra install).

# Chapters in reading order. Add new chapters here as the project grows.
CHAPTERS = \
	00-preface.md \
	01-foundations.md \
	02-leader-election.md \
	03-stage0-scaffolding.md \
	appendix-go-for-c.md \
	appendix-bug-checklist.md \
	references.md

META = metadata.yaml
OUT  = raftkv-book

.PHONY: all pdf epub html clean

all: html

# Polished PDF. Requires a LaTeX engine (see note above).
pdf: $(META) $(CHAPTERS)
	pandoc $(META) $(CHAPTERS) -o $(OUT).pdf --pdf-engine=tectonic

# E-reader friendly. Needs only pandoc.
epub: $(META) $(CHAPTERS)
	pandoc $(META) $(CHAPTERS) -o $(OUT).epub

# Self-contained styled HTML; open in a browser and Print > Save as PDF.
html: $(META) $(CHAPTERS)
	pandoc $(META) $(CHAPTERS) -o $(OUT).html --standalone --embed-resources

clean:
	rm -f $(OUT).pdf $(OUT).epub $(OUT).html
