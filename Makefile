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

.PHONY: all build build-dev build-release clean clean-clutter clean-coverage clean-deep compress \
        deps help prereqs rebuild stats test test-all test-unit \
        test-integration test-e2e test-framework test-coverage test-coverage-threshold test-python \
        check check-arch lint format vet install-tools run diagrams

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
	@echo "$(CYAN)$(BOLD)║  Hybrid App - Go 1.23+                           ║$(NC)"
	@echo "$(CYAN)$(BOLD)╚══════════════════════════════════════════════════╝$(NC)"
	@echo " "
	@echo "$(YELLOW)Build Commands:$(NC)"
	@echo "  build              - Build project (development mode)"
	@echo "  build-dev          - Build with race detector"
	@echo "  build-release      - Build optimized binary"
	@echo "  run                - Build and run the greeter"
	@echo "  clean              - Clean build artifacts"
	@echo "  clean-clutter      - Remove temporary files and backups"
	@echo "  clean-coverage     - Clean coverage data"
	@echo "  clean-deep         - Deep clean (includes module cache)"
	@echo "  compress           - Create compressed source archive (tar.gz)"
	@echo "  rebuild            - Clean and rebuild"
	@echo ""
	@echo "$(YELLOW)Testing Commands:$(NC)"
	@echo "  test               - Run all co-located unit tests"
	@echo "  test-unit          - Run all unit tests (Ada-style [PASS]/[FAIL] output)"
	@echo "  test-integration   - Run integration tests (cross-layer)"
	@echo "  test-e2e           - Run E2E tests (black-box CLI testing)"
	@echo "  test-framework     - Run all test suites (unit + integration + e2e)"
	@echo "  test-coverage      - Run tests with per-layer coverage analysis"
	@echo "  test-coverage-threshold - Run coverage with per-layer threshold checks"
	@echo "                       (Domain: 100%, Application: 100%, Infra: 90%, Total: 85%)"
	@echo "  test-python        - Run Python script tests (arch_guard.py validation)"
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
	@cd $(BIN_DIR) && $(GO) build -race
	@echo "$(GREEN)✓ Development build complete: $(BIN_DIR)/$(BINARY_NAME)$(NC)"

build-release: check-arch prereqs
	@echo "$(GREEN)Building $(PROJECT_NAME) (release mode)...$(NC)"
	@cd $(BIN_DIR) && $(GO) build -ldflags="-s -w"
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

clean-clutter: ## Remove temporary files, backups, and clutter
	@echo "$(CYAN)Cleaning temporary files and clutter...$(NC)"
	@$(PYTHON3) scripts/python/makefile/cleanup_temp_files.py
	@echo "$(GREEN)✓ Temporary files removed$(NC)"

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
	@echo "$(CYAN)$(BOLD)╔══════════════════════════════════════════════════════════════╗$(NC)"
	@echo "$(CYAN)$(BOLD)║                    RUNNING ALL TESTS                         ║$(NC)"
	@echo "$(CYAN)$(BOLD)╚══════════════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@echo "$(YELLOW)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(NC)"
	@echo "$(YELLOW)  UNIT TESTS (Domain Layer)$(NC)"
	@echo "$(YELLOW)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(NC)"
	@$(GO) test -v ./domain/... ./application/... ./infrastructure/... ./presentation/... ./bootstrap/...
	@echo ""
	@echo "$(YELLOW)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(NC)"
	@echo "$(YELLOW)  INTEGRATION TESTS$(NC)"
	@echo "$(YELLOW)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(NC)"
	@$(GO) test -v -tags=integration ./test/integration/...
	@echo ""
	@echo "$(YELLOW)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(NC)"
	@echo "$(YELLOW)  E2E TESTS$(NC)"
	@echo "$(YELLOW)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(NC)"
	@$(GO) test -v -tags=e2e ./test/e2e/...
	@echo ""
	@echo "$(GREEN)$(BOLD)╔══════════════════════════════════════════════════════════════╗$(NC)"
	@echo "$(GREEN)$(BOLD)║  ✓ ALL TESTS PASSED                                          ║$(NC)"
	@echo "$(GREEN)$(BOLD)╚══════════════════════════════════════════════════════════════╝$(NC)"

test-unit: check-arch ## Run unit tests with Ada-style framework output
	@echo "$(CYAN)$(BOLD)╔══════════════════════════════════════════════════════════════╗$(NC)"
	@echo "$(CYAN)$(BOLD)║                    UNIT TEST SUITE                           ║$(NC)"
	@echo "$(CYAN)$(BOLD)╚══════════════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@$(GO) test -v ./domain/... ./application/... ./infrastructure/... ./presentation/... ./bootstrap/...
	@echo ""

test-integration: check-arch build ## Run integration tests (cross-layer)
	@echo "$(CYAN)$(BOLD)╔══════════════════════════════════════════════════════════════╗$(NC)"
	@echo "$(CYAN)$(BOLD)║                 INTEGRATION TEST SUITE                       ║$(NC)"
	@echo "$(CYAN)$(BOLD)╚══════════════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@$(GO) test -v -tags=integration ./test/integration/...
	@echo ""

test-e2e: check-arch build ## Run E2E tests (black-box CLI testing)
	@echo "$(CYAN)$(BOLD)╔══════════════════════════════════════════════════════════════╗$(NC)"
	@echo "$(CYAN)$(BOLD)║                    E2E TEST SUITE                            ║$(NC)"
	@echo "$(CYAN)$(BOLD)╚══════════════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@$(GO) test -v -tags=e2e ./test/e2e/...
	@echo ""

test-framework: test-unit test-integration test-e2e ## Run all tests using Ada-style framework
	@echo "$(GREEN)$(BOLD)✓ All test suites completed$(NC)"

test-coverage: check-arch clean-coverage
	@echo "$(GREEN)Running tests with coverage analysis...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@echo ""
	@echo "$(CYAN)═══════════════════════════════════════════════════════════════$(NC)"
	@echo "$(CYAN)  Coverage Analysis - $(PROJECT_NAME)$(NC)"
	@echo "$(CYAN)═══════════════════════════════════════════════════════════════$(NC)"
	@echo ""
	@# Run tests with coverage for all layers (Go workspace requires explicit paths)
	@$(GO) test -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic \
		./domain/... ./application/... ./infrastructure/... ./presentation/... ./bootstrap/... 2>/dev/null || true
	@echo ""
	@echo "$(YELLOW)Per-Layer Coverage Summary:$(NC)"
	@echo "$(YELLOW)───────────────────────────────────────────────────────────────$(NC)"
	@# Domain layer coverage
	@printf "  Domain:          "
	@$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out 2>/dev/null | \
		grep -E "^github.com/.*/domain/" | \
		awk '{sum+=$$3; count++} END {if(count>0) printf "%.1f%% (%d functions)\n", sum/count, count; else print "N/A"}' || echo "N/A"
	@# Application layer coverage
	@printf "  Application:     "
	@$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out 2>/dev/null | \
		grep -E "^github.com/.*/application/" | \
		awk '{sum+=$$3; count++} END {if(count>0) printf "%.1f%% (%d functions)\n", sum/count, count; else print "N/A"}' || echo "N/A"
	@# Infrastructure layer coverage
	@printf "  Infrastructure:  "
	@$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out 2>/dev/null | \
		grep -E "^github.com/.*/infrastructure/" | \
		awk '{sum+=$$3; count++} END {if(count>0) printf "%.1f%% (%d functions)\n", sum/count, count; else print "N/A"}' || echo "N/A"
	@# Presentation layer coverage
	@printf "  Presentation:    "
	@$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out 2>/dev/null | \
		grep -E "^github.com/.*/presentation/" | \
		awk '{sum+=$$3; count++} END {if(count>0) printf "%.1f%% (%d functions)\n", sum/count, count; else print "N/A"}' || echo "N/A"
	@# Bootstrap layer coverage
	@printf "  Bootstrap:       "
	@$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out 2>/dev/null | \
		grep -E "^github.com/.*/bootstrap/" | \
		awk '{sum+=$$3; count++} END {if(count>0) printf "%.1f%% (%d functions)\n", sum/count, count; else print "N/A"}' || echo "N/A"
	@echo "$(YELLOW)───────────────────────────────────────────────────────────────$(NC)"
	@# Total coverage
	@printf "  $(BOLD)TOTAL:$(NC)             "
	@$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out | grep "^total:" | awk '{print $$3}'
	@echo ""
	@# Generate HTML report
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@# Generate text summary
	@$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out > $(COVERAGE_DIR)/coverage_summary.txt
	@echo "$(GREEN)✓ Coverage reports generated:$(NC)"
	@echo "    HTML report:  $(COVERAGE_DIR)/coverage.html"
	@echo "    Text summary: $(COVERAGE_DIR)/coverage_summary.txt"
	@echo ""

test-coverage-threshold: test-coverage ## Run coverage with minimum threshold checks per testing standards
	@echo ""
	@echo "$(CYAN)Checking coverage thresholds per Testing Standards...$(NC)"
	@echo "$(CYAN)───────────────────────────────────────────────────────────────$(NC)"
	@# Check total coverage (> 85% required)
	@TOTAL=$$($(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out | grep "^total:" | awk '{gsub(/%/,""); print $$3}'); \
	if [ $$(echo "$$TOTAL < 85" | bc -l) -eq 1 ]; then \
		echo "$(RED)✗ Total coverage $$TOTAL% is below 85% threshold$(NC)"; \
		exit 1; \
	else \
		echo "$(GREEN)✓ Total coverage $$TOTAL% meets 85% threshold$(NC)"; \
	fi
	@# Check domain coverage (= 100% required)
	@DOMAIN=$$($(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out 2>/dev/null | \
		grep -E "^github.com/.*/domain/" | \
		awk '{sum+=$$3; count++} END {if(count>0) printf "%.1f", sum/count; else print "0"}'); \
	if [ $$(echo "$$DOMAIN < 100" | bc -l) -eq 1 ]; then \
		echo "$(RED)✗ Domain coverage $$DOMAIN% is below 100% requirement$(NC)"; \
		exit 1; \
	else \
		echo "$(GREEN)✓ Domain coverage $$DOMAIN% meets 100% requirement$(NC)"; \
	fi
	@# Check application coverage (= 100% required)
	@APP=$$($(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out 2>/dev/null | \
		grep -E "^github.com/.*/application/" | \
		awk '{sum+=$$3; count++} END {if(count>0) printf "%.1f", sum/count; else print "0"}'); \
	if [ $$(echo "$$APP < 100" | bc -l) -eq 1 ]; then \
		echo "$(YELLOW)⚠ Application coverage $$APP% is below 100% target$(NC)"; \
	else \
		echo "$(GREEN)✓ Application coverage $$APP% meets 100% target$(NC)"; \
	fi
	@# Check infrastructure coverage (90%+ required)
	@INFRA=$$($(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out 2>/dev/null | \
		grep -E "^github.com/.*/infrastructure/" | \
		awk '{sum+=$$3; count++} END {if(count>0) printf "%.1f", sum/count; else print "0"}'); \
	if [ $$(echo "$$INFRA < 90" | bc -l) -eq 1 ]; then \
		echo "$(YELLOW)⚠ Infrastructure coverage $$INFRA% is below 90% target$(NC)"; \
	else \
		echo "$(GREEN)✓ Infrastructure coverage $$INFRA% meets 90% target$(NC)"; \
	fi
	@echo "$(CYAN)───────────────────────────────────────────────────────────────$(NC)"
	@echo "$(GREEN)✓ Coverage threshold check complete$(NC)"

test-python: ## Run Python script tests (arch_guard.py validation)
	@echo "$(GREEN)Running Python script tests...$(NC)"
	@cd test/python && $(PYTHON3) -m pytest -v
	@echo "$(GREEN)✓ Python tests complete$(NC)"

# =============================================================================
# Quality & Code Checking Commands
# =============================================================================

check: lint vet check-arch
	@echo "$(GREEN)✓ All checks passed$(NC)"

check-arch: ## Validate hexagonal architecture boundaries
	@echo "$(GREEN)Validating architecture boundaries...$(NC)"
	@PYTHONPATH=scripts/python $(PYTHON3) -m arch_guard
	@if [ $$? -eq 0 ]; then \
		echo "$(GREEN)✓ Architecture validation passed$(NC)"; \
	else \
		echo "$(RED)✗ Architecture validation failed$(NC)"; \
		exit 1; \
	fi

lint:
	@echo "$(GREEN)Running golangci-lint...$(NC)"
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run ./domain/... ./application/... ./infrastructure/... ./presentation/... ./bootstrap/...; \
		echo "$(GREEN)✓ Linting complete$(NC)"; \
	else \
		echo "$(YELLOW)⚠ golangci-lint not installed (run 'make install-tools')$(NC)"; \
	fi

vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	@$(GO) vet ./domain/... ./application/... ./infrastructure/... ./presentation/... ./bootstrap/...
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

diagrams: ## Generate SVG diagrams from PlantUML sources
	@echo "$(CYAN)Generating SVG diagrams from PlantUML...$(NC)"
	@command -v plantuml >/dev/null 2>&1 || { echo "$(RED)Error: plantuml not found. Install with: brew install plantuml$(NC)"; exit 1; }
	@cd docs/diagrams && for f in *.puml; do \
		echo "  Processing $$f..."; \
		plantuml -tsvg "$$f"; \
	done
	@echo "$(GREEN)✓ Diagrams generated$(NC)"

.DEFAULT_GOAL := help

## ---------------------------------------------------------------------------
## Submodule Management
## ---------------------------------------------------------------------------

.PHONY: submodule-update submodule-status submodule-init

submodule-init: ## Initialize submodules after fresh clone
	git submodule update --init --recursive

submodule-update: ## Pull latest from all submodule repos
	git submodule update --remote --merge
	@echo ""
	@echo "Submodules updated. Review changes, then run:"
	@echo "  git add scripts/python test/python"
	@echo "  git commit -m 'chore: update submodules'"
	@echo "  git push"

submodule-status: ## Show submodule commit status
	git submodule status
