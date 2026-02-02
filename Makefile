.PHONY: help build frontend backend install uninstall clean run dev test

# Default target
help:
	@echo "GoScouter - Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make build      - Build both frontend and backend"
	@echo "  make frontend   - Build frontend only"
	@echo "  make backend    - Build backend only"
	@echo "  make install    - Install goscouter system-wide"
	@echo "  make uninstall  - Remove goscouter from system"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make run        - Build and run goscouter"
	@echo "  make dev        - Run frontend in dev mode"
	@echo "  make test       - Run tests"

# Build everything
build: frontend backend

# Build frontend
frontend:
	@echo "📦 Building frontend..."
	@cd frontend && npm install && npm run build
	@echo "✓ Frontend built"

# Build backend
backend:
	@echo "🔨 Building backend..."
	@go build -o goscouter .
	@echo "✓ Backend built"

# Install system-wide
install:
	@./install.sh

# Uninstall
uninstall:
	@echo "🗑️  Uninstalling goscouter..."
	@sudo rm -f /usr/local/bin/goscouter
	@rm -rf $(HOME)/.goscouter
	@echo "✓ Uninstalled"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f goscouter
	@rm -rf frontend/out
	@rm -rf frontend/.next
	@echo "✓ Cleaned"

# Build and run
run: build
	@./goscouter run

# Run frontend in development mode
dev:
	@cd frontend && npm run dev

# Run tests
test:
	@go test ./...
