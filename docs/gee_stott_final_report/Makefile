# Makefile for Gee Stott Final Report (IEEE Format)
# Requires: pdflatex (part of texlive)

# Variables
MAIN = report
OUTPUT = gee_stott_final_report.pdf
FIGURES = $(wildcard figures/*.png)

# Default target
.PHONY: all
all: $(OUTPUT)

# Build the PDF
$(OUTPUT): $(MAIN).tex $(FIGURES)
	@echo "Building IEEE format report..."
	pdflatex $(MAIN).tex
	pdflatex $(MAIN).tex
	@mv $(MAIN).pdf $(OUTPUT)
	@echo "Report built successfully: $(OUTPUT)"

# Clean generated files
.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	rm -f $(MAIN).pdf $(OUTPUT) $(MAIN).aux $(MAIN).log $(MAIN).out $(MAIN).bbl $(MAIN).blg
	@echo "Clean complete."

# View the PDF (requires a PDF viewer)
.PHONY: view
view: $(OUTPUT)
	@if command -v xdg-open > /dev/null 2>&1; then \
		xdg-open $(OUTPUT); \
	elif command -v open > /dev/null 2>&1; then \
		open $(OUTPUT); \
	elif command -v evince > /dev/null 2>&1; then \
		evince $(OUTPUT) & \
	elif command -v okular > /dev/null 2>&1; then \
		okular $(OUTPUT) & \
	else \
		echo "No PDF viewer found. Please open $(OUTPUT) manually."; \
	fi

# Check if pdflatex is installed
.PHONY: check
check:
	@echo "Checking for required tools..."
	@command -v pdflatex >/dev/null 2>&1 || { echo "Error: pdflatex is not installed. Install texlive."; exit 1; }
	@echo "pdflatex version: $$(pdflatex --version | head -n 1)"
	@echo "All required tools are installed."

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all (default)  - Build the PDF report"
	@echo "  clean          - Remove generated PDF and auxiliary files"
	@echo "  view           - Build and open the PDF in a viewer"
	@echo "  check          - Verify required tools are installed"
	@echo "  help           - Show this help message"
