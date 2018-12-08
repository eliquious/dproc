# Simple Makefile for running tests and building website
.PHONY: docs

docs: FORCE
	@echo "Building docs..."
	@echo "Running hugo on landing page..."
	@cd site/landing_page && hugo --minify -d "../../docs"
	@echo "Running hugo on articles..."
	@cd site/articles && hugo --minify -d "../../docs/articles"
	@echo "Done."

serve: docs
	@cd docs && static -debug

FORCE: