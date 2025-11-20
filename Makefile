# =============================================================================
# Project Makefile
# =============================================================================
# Project: hybrid_app_go
# Purpose: Hexagonal architecture demonstration with port/adapter pattern
#
# This Makefile provides:
#   - Build targets (build, clean, rebuild)
#   - Test infrastructure (test, test-coverage)
#   - Format/check targets (format, lint, stats)
#   - Documentation generation (docs)
#   - Development tools (check-arch, ci)
# =============================================================================

PROJECT_NAME := hybrid_app_go
BINARY_NAME := greeter

.PHONY: all build build-dev build-release clean clean-coverage clean-deep compress \
        deps help prereqs rebuild stats test test-all test-unit test-domain \
        check check-arch lint format vet install-tools run

# =============================================================================
# Colors for Output
# =============================================================================

GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
BLUE := \033[0;34m
CYAN := \033[0;36m
BOLD := \033[1m
NC := \033[0m

# =============================================================================
# Tool Paths
# =============================================================================

GO := go
GOFMT := gofmt
GOLINT := golangci-lint
PYTHON3 := python3

# =============================================================================
# Directories
# =============================================================================

BIN_DIR := cmd/greeter
COVERAGE_DIR := coverage
MAKEFILE_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

# =============================================================================
# Default Target
# =============================================================================

all: build

# =============================================================================
# Help Target
# =============================================================================

help: ## Display this help message
	@echo "$(CYAN)$(BOLD)╔══════════════════════════════════════════════════╗$(NC)"
	@echo "$(CYAN)$(BOLD)║  Hybrid App - Go 1.21+                           ║$(NC)"
	@echo "$(CYAN)$(BOLD)╚══════════════════════════════════════════════════╝$(NC)"
	@echo " "
	@echo "$(YELLOW)Build Commands:$(NC)"
	@echo "  build              - Build project (development mode)"
	@echo "  build-dev          - Build with race detector"
	@echo "  build-release      - Build optimized binary"
	@echo "  run                - Build and run the greeter"
	@echo "  clean              - Clean build artifacts"
	@echo "  clean-coverage     - Clean coverage data"
	@echo "  clean-deep         - Deep clean (includes module cache)"
	@echo "  compress           - Create compressed source archive (tar.gz)"
	@echo "  rebuild            - Clean and rebuild"
	@echo ""
	@echo "$(YELLOW)Testing Commands:$(NC)"
	@echo "  test               - Run all tests"
	@echo "  test-unit          - Run unit tests only"
	@echo "  test-domain        - Run domain layer tests"
	@echo "  test-coverage      - Run tests with coverage analysis"
	@echo ""
	@echo "$(YELLOW)Quality & Architecture Commands:$(NC)"
	@echo "  check              - Run all checks (lint + vet + arch)"
	@echo "  check-arch         - Validate hexagonal architecture boundaries"
	@echo "  lint               - Run golangci-lint"
	@echo "  vet                - Run go vet"
	@echo "  format             - Format all Go code"
	@echo "  stats              - Display project statistics by layer"
	@echo ""
	@echo "$(YELLOW)Utility Commands:$(NC)"
	@echo "  deps               - Show dependency information"
	@echo "  prereqs            - Verify prerequisites are satisfied"
	@echo "  install-tools      - Install development tools (golangci-lint)"
	@echo ""
	@echo "$(YELLOW)Workflow Shortcuts:$(NC)"
	@echo "  all                - Build project (default)"

# =============================================================================
# Build Commands
# =============================================================================

prereqs:
	@echo "$(GREEN)Checking prerequisites...$(NC)"
	@command -v $(GO) >/dev/null 2>&1 || { echo "$(RED)✗ go not found$(NC)"; exit 1; }
	@command -v $(PYTHON3) >/dev/null 2>&1 || { echo "$(RED)✗ python3 not found$(NC)"; exit 1; }
	@echo "$(GREEN)✓ All prerequisites satisfied$(NC)"

build: build-dev

build-dev: check-arch prereqs
	@echo "$(GREEN)Building $(PROJECT_NAME) (development mode)...$(NC)"
	@cd $(BIN_DIR) && $(GO) build -race -v
	@echo "$(GREEN)✓ Development build complete: $(BIN_DIR)/$(BINARY_NAME)$(NC)"

build-release: check-arch prereqs
	@echo "$(GREEN)Building $(PROJECT_NAME) (release mode)...$(NC)"
	@cd $(BIN_DIR) && $(GO) build -ldflags="-s -w" -v
	@echo "$(GREEN)✓ Release build complete: $(BIN_DIR)/$(BINARY_NAME)$(NC)"

run: build
	@echo "$(GREEN)Running $(BINARY_NAME)...$(NC)"
	@cd $(BIN_DIR) && ./$(BINARY_NAME) World

clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -f $(BIN_DIR)/$(BINARY_NAME)
	@$(GO) clean -cache -testcache
	@find . -name "*.test" -delete 2>/dev/null || true
	@find . -name "*.out" -delete 2>/dev/null || true
	@echo "$(GREEN)✓ Build artifacts cleaned$(NC)"

clean-deep: clean
	@echo "$(YELLOW)Deep cleaning ALL artifacts including module cache...$(NC)"
	@$(GO) clean -modcache
	@echo "$(GREEN)✓ Deep clean complete (next build will download modules)$(NC)"

clean-coverage:
	@echo "$(YELLOW)Cleaning coverage artifacts...$(NC)"
	@rm -rf $(COVERAGE_DIR) 2>/dev/null || true
	@find . -name "coverage.txt" -delete 2>/dev/null || true
	@find . -name "coverage.html" -delete 2>/dev/null || true
	@echo "$(GREEN)✓ Coverage artifacts cleaned$(NC)"

compress:
	@echo "$(CYAN)Creating compressed source archive...$(NC)"
	@tar -czvf "$(PROJECT_NAME).tar.gz" \
		--exclude="$(PROJECT_NAME).tar.gz" \
		--exclude='.git' \
		--exclude='vendor' \
		--exclude='bin' \
		--exclude='.build' \
		--exclude='coverage' \
		--exclude='.DS_Store' \
		--exclude='*.test' \
		--exclude='*.out' \
		.
	@echo "$(GREEN)✓ Archive created: $(PROJECT_NAME).tar.gz$(NC)"

rebuild: clean build

# =============================================================================
# Testing Commands
# =============================================================================

test: test-all

test-all: check-arch
	@echo "$(GREEN)Running all tests...$(NC)"
	@$(GO) test -v ./...
	@echo "$(GREEN)✓ All tests passed$(NC)"

test-unit: check-arch
	@echo "$(GREEN)Running unit tests...$(NC)"
	@$(GO) test -v -short ./...
	@echo "$(GREEN)✓ Unit tests passed$(NC)"

test-domain:
	@echo "$(GREEN)Running domain layer tests...$(NC)"
	@cd domain && $(GO) test -v ./...
	@echo "$(GREEN)✓ Domain tests passed$(NC)"

test-coverage: check-arch
	@echo "$(GREEN)Running tests with coverage analysis...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@$(GO) test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out
	@echo "$(GREEN)✓ Coverage report generated: $(COVERAGE_DIR)/coverage.html$(NC)"

# =============================================================================
# Quality & Code Checking Commands
# =============================================================================

check: lint vet check-arch
	@echo "$(GREEN)✓ All checks passed$(NC)"

check-arch: ## Validate hexagonal architecture boundaries
	@echo "$(GREEN)Validating architecture boundaries...$(NC)"
	@$(PYTHON3) scripts/makefile/arch_guard.py
	@if [ $$? -eq 0 ]; then \
		echo "$(GREEN)✓ Architecture validation passed$(NC)"; \
	else \
		echo "$(RED)✗ Architecture validation failed$(NC)"; \
		exit 1; \
	fi

lint:
	@echo "$(GREEN)Running golangci-lint...$(NC)"
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run ./...; \
		echo "$(GREEN)✓ Linting complete$(NC)"; \
	else \
		echo "$(YELLOW)⚠ golangci-lint not installed (run 'make install-tools')$(NC)"; \
	fi

vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	@$(GO) vet ./...
	@echo "$(GREEN)✓ Vet complete$(NC)"

format:
	@echo "$(GREEN)Formatting Go code...$(NC)"
	@$(GOFMT) -w -s .
	@echo "$(GREEN)✓ Code formatting complete$(NC)"

# =============================================================================
# Development Commands
# =============================================================================

stats:
	@echo "$(CYAN)$(BOLD)Project Statistics for $(PROJECT_NAME)$(NC)"
	@echo "$(YELLOW)════════════════════════════════════════$(NC)"
	@echo ""
	@echo "Go Source Files by Layer:"
	@echo "  Domain:          $$(find domain -name "*.go" ! -name "*_test.go" 2>/dev/null | wc -l | tr -d ' ')"
	@echo "  Application:     $$(find application -name "*.go" ! -name "*_test.go" 2>/dev/null | wc -l | tr -d ' ')"
	@echo "  Infrastructure:  $$(find infrastructure -name "*.go" ! -name "*_test.go" 2>/dev/null | wc -l | tr -d ' ')"
	@echo "  Presentation:    $$(find presentation -name "*.go" ! -name "*_test.go" 2>/dev/null | wc -l | tr -d ' ')"
	@echo "  Bootstrap:       $$(find bootstrap -name "*.go" ! -name "*_test.go" 2>/dev/null | wc -l | tr -d ' ')"
	@echo ""
	@echo "Test Files:"
	@echo "  Unit tests:      $$(find . -name "*_test.go" 2>/dev/null | wc -l | tr -d ' ')"
	@echo ""
	@echo "Lines of Code:"
	@find domain application infrastructure presentation bootstrap -name "*.go" ! -name "*_test.go" 2>/dev/null | \
	  xargs wc -l 2>/dev/null | tail -1 | awk '{printf "  Total: %d lines\n", $$1}' || echo "  Total: 0 lines"
	@echo ""
	@echo "Build Artifacts:"
	@if [ -f "$(BIN_DIR)/$(BINARY_NAME)" ]; then \
		echo "  Binary: $$(ls -lh $(BIN_DIR)/$(BINARY_NAME) 2>/dev/null | awk '{print $$5}')"; \
	else \
		echo "  No binary found (run 'make build')"; \
	fi

# =============================================================================
# Utility Targets
# =============================================================================

deps: ## Display project dependencies
	@echo "$(CYAN)Go module dependencies:$(NC)"
	@$(GO) list -m all
	@echo ""
	@echo "$(CYAN)Module graph:$(NC)"
	@$(GO) mod graph | head -20

install-tools: ## Install development tools
	@echo "$(CYAN)Installing development tools...$(NC)"
	@echo "  Installing golangci-lint..."
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✓ Tool installation complete$(NC)"

.DEFAULT_GOAL := help
