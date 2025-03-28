# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."

	@go build -o bin/main src/cmd/main.go

	@echo "Done! Binary now is in bin/main"

# Run the application
run:
	@echo "Running..."

	@go run src/cmd/main.go

# Create container
up:
	@if docker compose up --build -d; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build -d; \
	fi

# Shutdown container
down:
	@if docker compose down; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."

	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."

	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f bin/main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
		air; \
		echo "Watching...";\
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
			echo "Watching...";\
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi

.PHONY: all build run test clean watch up down itest
