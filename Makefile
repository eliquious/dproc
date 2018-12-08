# Simple Makefile for running tests and building website

docs:
	@echo "Building docs..."
	@echo "Running hugo on landing page..."
	@cd site/landing_page && hugo -d "../../docs"
	@echo "Running hugo on articles..."
	@cd site/articles && hugo -d "../../docs/articles"
	@echo "Done."

serve: docs
	@cd docs && python -m SimpleHTTPServer
