# Hybrid_App_Go Quick Start Guide

**Version:** 1.0.1  
**Date:** November 26, 2025  
**SPDX-License-Identifier:** BSD-3-Clause  
**License File:** See the LICENSE file in the project root.  
**Copyright:** © 2025 Michael Gardner, A Bit of Help, Inc.  
**Status:** Released  

---

## Table of Contents

- [Installation](#installation)
- [First Build](#first-build)
- [Running the Application](#running-the-application)
- [Understanding the Architecture](#understanding-the-architecture)
- [Making Your First Change](#making-your-first-change)
- [Running Tests](#running-tests)
- [Build Targets](#build-targets)
- [Common Issues](#common-issues)
- [Next Steps](#next-steps)

---

## Installation

### Prerequisites

- **Go**: Version 1.23+ (workspace and generics support required)
- **Make**: GNU Make for build automation
- **Python 3**: For architecture validation (optional)
- **Java 11+**: For PlantUML diagram generation (optional)

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/abitofhelp/hybrid_app_go.git
cd hybrid_app_go

# Build with Make
make build

# Or build directly with Go
go build -o cmd/greeter/greeter ./cmd/greeter
```

### Verify Installation

```bash
# Check that the executable was built
ls -lh cmd/greeter/greeter

# Run the application
./cmd/greeter/greeter World
# Output: Hello, World!
```

**Success!** You've built your first hexagonal architecture application in Go.

---

## First Build

The project uses Make for building:

### Using Make (Recommended)

```bash
# Development build (with race detector)
make build

# Or explicit development mode
make build-dev

# Release build (optimized)
make build-release
```

### Using Go Directly

```bash
# Development build
go build -race -o cmd/greeter/greeter ./cmd/greeter

# Release build
go build -ldflags="-s -w" -o cmd/greeter/greeter ./cmd/greeter
```

**Build Output:**
- Executable: `cmd/greeter/greeter`

---

## Running the Application

The Hybrid_App_Go starter includes a simple greeter application demonstrating all architectural layers:

### Basic Usage

```bash
# Greet a person
./cmd/greeter/greeter Alice
# Output: Hello, Alice!

# Name with spaces (use quotes)
./cmd/greeter/greeter "Bob Smith"
# Output: Hello, Bob Smith!

# Show usage
./cmd/greeter/greeter
# Output: Usage: greeter <name>
# Exit code: 1
```

### Error Handling Example

```bash
# Empty name triggers validation error
./cmd/greeter/greeter ""
# Output: Error: Person name cannot be empty
# Exit code: 1
```

**Key Points:**
- All errors return via Result monad (no panics across boundaries)
- Exit code 0 = success, 1 = error
- Validation happens in Domain layer
- Errors propagate through Application to Presentation

---

## Understanding the Architecture

Hybrid_App_Go demonstrates **5-layer hexagonal architecture**:

### Layer Overview

```
┌─────────────────────────────────────────────┐
│  Bootstrap (Composition Root)               │  ← Wires everything together
├─────────────────────────────────────────────┤
│  Presentation (CLI)                         │  ← User interface (depends on Application ONLY)
├─────────────────────────────────────────────┤
│  Application (Use Cases + Ports)            │  ← Orchestration layer
├─────────────────────────────────────────────┤
│  Infrastructure (Adapters)                  │  ← Technical implementations
├─────────────────────────────────────────────┤
│  Domain (Business Logic)                    │  ← Pure business rules (ZERO dependencies)
└─────────────────────────────────────────────┘
```

### Key Architectural Principles

1. **Domain has zero dependencies** - Pure business logic
2. **Presentation cannot access Domain** - Must use Application layer re-exports
3. **Static dependency injection** - Via generics (compile-time wiring)
4. **Railway-oriented programming** - Result monads for error handling
5. **Multi-module workspace** - go.work manages separate go.mod per layer

### Request Flow Example

```
User Input ("Alice")
    ↓
presentation/adapter/cli/command.GreetCommand (parses input)
    ↓
application/usecase.GreetUseCase (validates via Domain)
    ↓
domain/valueobject.Person (business rules)
    ↓
infrastructure/adapter.ConsoleWriter (output)
    ↓
Result[Unit] (success or error)
    ↓
Exit Code (0 or 1)
```

---

## Making Your First Change

Let's modify the greeting message:

### Step 1: Locate the Domain Logic

```bash
# Open the Person value object
# File: domain/valueobject/person.go
```

### Step 2: Modify Greeting Format

Find the `GreetingMessage()` method and modify it. The business logic is pure and has no dependencies.

### Step 3: Rebuild and Test

```bash
# Rebuild
make rebuild

# Run tests to ensure nothing broke
make test-all

# Test manually
./cmd/greeter/greeter Alice
```

**Best Practice**: Always run tests after making changes.

---

## Running Tests

Hybrid_App_Go includes comprehensive testing:

### Test Organization

- **Unit Tests** (42 assertions): Domain layer logic
- **Integration Tests** (21 tests): CLI binary execution
- **E2E Tests** (10 tests): Full system verification

### Run All Tests

```bash
# Run entire test suite
make test-all

# Expected output:
# ########################################
# ###                                  ###
# ###   ALL TESTS PASSED               ###
# ###                                  ###
# ########################################
```

### Run Specific Test Suites

```bash
# Unit tests only (fast)
make test

# Integration tests only
make test-integration

# E2E tests only
make test-e2e
```

**Test Framework**: Custom lightweight framework in `domain/test/` plus testify for assertions.

---

## Build Targets

### Building

```bash
make build              # Development build (default)
make build-dev          # Explicit development mode
make build-release      # Release build (optimized)
make rebuild            # Clean and rebuild
```

### Testing

```bash
make test               # Run unit tests
make test-all           # Run all tests (unit + integration + e2e)
make test-integration   # Integration tests only
make test-e2e           # E2E tests only
```

### Quality & Architecture

```bash
make check-arch         # Validate architecture boundaries
make fmt                # Format code with gofmt
make lint               # Run golangci-lint
```

### Cleaning

```bash
make clean              # Clean build artifacts
```

### Utilities

```bash
make run NAME=Alice     # Build and run with argument
make help               # Show all available targets
```

---

## Common Issues

### Q: Build fails with "go.work not found"

**A:** Ensure you're in the project root directory:

```bash
pwd
# Should show: .../hybrid_app_go
```

### Q: "package not found" errors

**A:** Sync the workspace:

```bash
go work sync
```

### Q: Architecture validation warnings appear

**A:** The `make check-arch` target validates layer boundaries:

```bash
# View architecture validation
make check-arch
```

Violations indicate forbidden imports (e.g., Presentation importing Domain).

### Q: Tests fail with "binary not found"

**A:** Build the application first:

```bash
make build
make test-all
```

### Q: How do I run a single test?

**A:** Use go test directly:

```bash
# Run specific test
go test -v -run TestGreeter_ValidName ./test/integration/...

# Run with build tag
go test -v -tags=integration ./test/integration/...
```

---

## Next Steps

### Explore the Architecture

- **[Software Design Specification](formal/software_design_specification.md)** - Deep dive into architecture
- **[Architecture Diagrams](diagrams/)** - Visual documentation
- **[Layer Dependencies](diagrams/layer_dependencies.svg)** - See dependency flow

### Read the Source Code

Start with the wiring in Bootstrap:

```bash
# See how all layers are wired together
cat bootstrap/cli/cli.go
```

Then explore each layer:

```bash
# Domain (pure business logic)
ls domain/

# Application (use cases and ports)
ls application/

# Infrastructure (adapters)
ls infrastructure/

# Presentation (CLI)
ls presentation/

# Bootstrap (composition root)
ls bootstrap/
```

### Study the Test Suite

```bash
# See how tests are organized
ls -R test/

# Read test framework
cat domain/test/test_framework.go
```

### Understand Static Dispatch

```go
// Generic type with interface constraint
type GreetUseCase[W outbound.WriterPort] struct {
    writer W
}

// Concrete type known at compile time
uc := usecase.NewGreetUseCase[*adapter.ConsoleWriter](writer)

// Method call is statically dispatched (no vtable)
uc.Execute(ctx, cmd)
```

### Add Your Own Use Case

Follow the pattern:

1. **Domain**: Create value objects/entities (`domain/valueobject/`)
2. **Application**: Define command, use case, ports (`application/`)
3. **Infrastructure**: Implement adapters (`infrastructure/adapter/`)
4. **Presentation**: Create CLI command (`presentation/adapter/cli/command/`)
5. **Bootstrap**: Wire everything together (`bootstrap/cli/`)
6. **Tests**: Add unit/integration/e2e tests (`test/`)

---

## Documentation Index

- 📖 **[Main Documentation Hub](index.md)** - All documentation links
- 📋 **[Software Requirements Specification](formal/software_requirements_specification.md)** - Requirements
- 🏗️ **[Software Design Specification](formal/software_design_specification.md)** - Architecture
- 🧪 **[Software Test Guide](formal/software_test_guide.md)** - Testing guide

---

## Support

For questions or issues:

- 🐛 **Issues**: [GitHub Issues](https://github.com/abitofhelp/hybrid_app_go/issues)
- 📖 **Documentation**: See `docs/` directory

---

## License

Hybrid_App_Go is licensed under the BSD-3-Clause License.
Copyright © 2025 Michael Gardner, A Bit of Help, Inc.

See [LICENSE](../LICENSE) for full license text.
