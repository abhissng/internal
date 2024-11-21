# Makefile to run tests

# Variable for the Bash script
RUN_INITIAL_BUILD_SCRIPT := ./build.sh

# Default target
all: run_build_checks

# Run tests
run_build_checks:
	@echo "Running build checks"
	go mod tidy
	@bash $(RUN_INITIAL_BUILD_SCRIPT)

# Clean up any generated files (e.g., test reports)
clean:
	@echo "Cleaning up..."
	@rm -f *.log

# Help command to display available targets
help:
	@echo "Makefile for running build checks Go code"
	@echo "Available targets:"
	@echo "  all        - Run static and security tests"
	@echo "  run_build_checks  - Run build checks"
	@echo "  clean      - Clean up generated files"
	@echo "  help       - Show this help message"
