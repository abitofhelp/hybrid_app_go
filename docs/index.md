# Hybrid_App_Go Documentation Index

**Version:** 1.0.2  
**Date:** November 28, 2025  
**SPDX-License-Identifier:** BSD-3-Clause  
**License File:** See the LICENSE file in the project root.  
**Copyright:** © 2025 Michael Gardner, A Bit of Help, Inc.  
**Status:** Released  

---

## Welcome

Welcome to the **Hybrid_App_Go** documentation. This Go 1.23+ application starter demonstrates professional hexagonal architecture with functional programming principles, static dependency injection via generics, and railway-oriented error handling.

---

## Quick Navigation

### Getting Started

- 🚀 **[Quick Start Guide](quick_start.md)** - Get up and running in minutes
  - Installation instructions
  - First build and run
  - Understanding the architecture
  - Making your first change
  - Running tests

### Formal Documentation

- 📋 **[Software Requirements Specification (SRS)](formal/software_requirements_specification.md)** - Complete requirements
  - Functional requirements (FR-01 through FR-12)
  - Non-functional requirements (NFR-01 through NFR-06)
  - System constraints
  - Test coverage mapping

- 🏗️ **[Software Design Specification (SDS)](formal/software_design_specification.md)** - Architecture and design
  - 5-layer hexagonal architecture
  - Static dependency injection via generics
  - Railway-oriented programming patterns
  - Component relationships
  - Data flow diagrams
  - Design patterns used

- 🧪 **[Software Test Guide](formal/software_test_guide.md)** - Testing documentation
  - Test organization (unit/integration/e2e)
  - Running tests (make test, make test-all)
  - Test framework documentation
  - Coverage procedures
  - Writing new tests

---

## Architecture Overview

Hybrid_App_Go implements a **5-layer hexagonal architecture** (also known as Ports and Adapters or Clean Architecture):

### Layer Structure

```
┌─────────────────────────────────────────────┐
│  Bootstrap                                  │  Composition Root (wiring)
├─────────────────────────────────────────────┤
│  Presentation                               │  Driving Adapters (CLI)
├─────────────────────────────────────────────┤
│  Application                                │  Use Cases + Ports
├─────────────────────────────────────────────┤
│  Infrastructure                             │  Driven Adapters (Console Writer)
├─────────────────────────────────────────────┤
│  Domain                                     │  Business Logic (ZERO dependencies)
└─────────────────────────────────────────────┘
```

### Key Principles

1. **Domain Isolation**: Domain layer has zero external dependencies
2. **Presentation Boundary**: Presentation layer cannot access Domain directly (uses application/error re-exports)
3. **Static Dispatch**: Dependency injection via generics (compile-time, zero overhead)
4. **Railway-Oriented**: Result monads for error handling (no panics across boundaries)
5. **Multi-Module Workspace**: go.work manages separate go.mod per layer

---

## Visual Documentation

### UML Diagrams

Located in `diagrams/` directory:

- **layer_dependencies.svg** - Shows 5-layer dependency flow
- **application_error_pattern.svg** - Re-export pattern for Presentation isolation
- **package_structure.svg** - Actual package hierarchy
- **error_handling_flow.svg** - Error propagation through layers
- **static_dispatch.svg** - Generic vs interface comparison

All diagrams are generated from PlantUML sources (.puml files).

---

## Project Statistics

### Code Metrics (v1.0.0)

- **Go Source Files**: 20 (.go)
- **Test Files**: Unit, Integration, E2E suites
- **Architecture Layers**: 5 (Domain, Application, Infrastructure, Presentation, Bootstrap)
- **Build Targets**: 20+ Makefile targets
- **Dependencies**: testify (test only), ZERO in domain layer

### Test Coverage

- **Unit Tests**: 42 assertions (Domain layer)
- **Integration Tests**: 21 tests (CLI binary execution)
- **E2E Tests**: 10 tests (Full system verification)
- **Test Framework**: Custom lightweight framework + testify

### Code Quality

- **Compiler Warnings**: 0
- **Architecture Validation**: Enforced by arch_guard.py
- **Go Version**: 1.23+

---

## Key Features

### Static Dependency Injection

Uses **generics** instead of interfaces for dependency injection:

```go
// Port definition (interface constraint)
type WriterPort interface {
    Write(ctx context.Context, message string) domerr.Result[model.Unit]
}

// Generic use case with static dispatch
type GreetUseCase[W WriterPort] struct {
    writer W
}

// Wiring in Bootstrap (compile-time resolution)
consoleWriter := adapter.NewConsoleWriter()
uc := usecase.NewGreetUseCase[*adapter.ConsoleWriter](consoleWriter)
```

**Benefits**:
- Zero runtime overhead (no vtable lookups)
- Full inlining potential
- Compile-time type safety
- All method calls devirtualized

### Railway-Oriented Programming

All errors propagate via **Result monad** (no panics across boundaries):

```go
// Use case returns Result[Unit]
func (uc *GreetUseCase[W]) Execute(ctx context.Context, cmd command.GreetCommand) domerr.Result[model.Unit] {
    personResult := valueobject.CreatePerson(cmd.Name())

    if personResult.IsError() {
        return domerr.Err[model.Unit](personResult.ErrorInfo())
    }

    person := personResult.Value()
    return uc.writer.Write(ctx, person.GreetingMessage())
}
```

**Benefits**:
- Explicit error handling (compiler-enforced)
- No unexpected control flow
- Composable error types

### Application.Error Re-Export Pattern

**Problem**: Presentation cannot access Domain directly
**Solution**: Application re-exports Domain.Error for Presentation

```go
// application/error/error.go (zero-overhead type aliases)
import domerr "github.com/abitofhelp/hybrid_app_go/domain/error"

type ErrorType = domerr.ErrorType
type ErrorKind = domerr.ErrorKind
type Result[T any] = domerr.Result[T]

var ValidationError = domerr.ValidationError
var InfrastructureError = domerr.InfrastructureError
```

This maintains clean boundaries while allowing Presentation to handle errors.

---

## Directory Structure

```
hybrid_app_go/
├── domain/                    # Pure business logic
│   ├── go.mod                 # ZERO external dependencies
│   ├── error/                 # Result monad, error types
│   ├── valueobject/           # Immutable value objects
│   └── test/                  # Test framework
├── application/               # Use cases + ports
│   ├── go.mod                 # Depends only on domain
│   ├── command/               # Input DTOs
│   ├── error/                 # Re-exports for Presentation
│   ├── model/                 # Output DTOs (Unit)
│   ├── port/
│   │   ├── inbound/           # Use case interfaces
│   │   └── outbound/          # Infrastructure interfaces
│   └── usecase/               # Use case orchestration
├── infrastructure/            # Adapters (driven)
│   ├── go.mod                 # Depends on application + domain
│   └── adapter/               # Console writer
├── presentation/              # Adapters (driving)
│   ├── go.mod                 # Depends only on application
│   └── adapter/cli/command/   # CLI commands
├── bootstrap/                 # Composition root
│   ├── go.mod                 # Depends on all layers
│   └── cli/                   # CLI wiring
├── cmd/greeter/               # Main entry point
│   ├── go.mod                 # Depends only on bootstrap
│   └── main.go                # 3 lines
├── test/
│   ├── integration/           # CLI integration tests
│   └── e2e/                   # Full system tests
├── docs/
│   ├── formal/                # SRS, SDS, Test Guide
│   ├── diagrams/              # UML diagrams
│   ├── quick_start.md         # Getting started
│   └── index.md               # This file
├── scripts/
│   └── arch_guard/            # Architecture validation
├── go.work                    # Workspace definition
├── Makefile                   # Build automation
└── README.md                  # Project overview
```

---

## Build System

### Make Targets

**Building**:
```bash
make build              # Development build
make build-release      # Release build
make rebuild            # Clean + build
```

**Testing**:
```bash
make test               # Run unit tests
make test-all           # Run all tests (unit + integration + e2e)
make test-integration   # Integration tests only
make test-e2e           # E2E tests only
```

**Quality**:
```bash
make check-arch         # Architecture validation
make fmt                # Format code
make lint               # Run linter
```

**Utilities**:
```bash
make clean              # Clean artifacts
make run NAME=Alice     # Run application
make help               # Show all targets
```

---

## Learning Path

### For Beginners

1. **Start Here**: [Quick Start Guide](quick_start.md)
2. **Understand Architecture**: [Software Design Specification](formal/software_design_specification.md)
3. **Run Tests**: `make test-all`
4. **Explore Code**: Start with `bootstrap/cli/cli.go`
5. **Read Examples**: Study how layers are wired together

### For Experienced Developers

1. **Architecture Patterns**: See [SDS - Design Patterns](formal/software_design_specification.md)
2. **Static DI Deep Dive**: See diagrams/static_dispatch.svg
3. **Railway-Oriented Programming**: See diagrams/error_handling_flow.svg
4. **Add Use Case**: Follow pattern in existing code

---

## Dependencies

### Runtime Dependencies

- **None**: Domain layer has zero external dependencies

### Development Dependencies

- **testify** (v1.9.0): Testing assertions (test modules only)

### Build Requirements

- **Go**: 1.23+ (workspace and generics support)
- **Make**: GNU Make for build automation
- **Python 3**: For architecture validation (arch_guard.py)
- **Java 11+**: For PlantUML diagram generation (optional)

---

## Documentation Updates

All documentation is maintained for v1.0.0 release:

- **Copyright**: © 2025 Michael Gardner, A Bit of Help, Inc.
- **License**: BSD-3-Clause
- **Version**: 1.0.0
- **Date**: November 25, 2025
- **Status**: Released

For documentation issues or suggestions, please file an issue on GitHub.

---

## Support and Contributing

### Getting Help

- 🐛 **Issues**: [GitHub Issues](https://github.com/abitofhelp/hybrid_app_go/issues)
- 📖 **Documentation**: This directory

### Contributing

We welcome contributions! See:

- Code style enforced by architecture validation
- Run `make test-all` before submitting

---

## License

Hybrid_App_Go is licensed under the **BSD-3-Clause License**.

Copyright © 2025 Michael Gardner, A Bit of Help, Inc.

See [LICENSE](../LICENSE) for full license text.

---

## Project Links

- **GitHub**: https://github.com/abitofhelp/hybrid_app_go
- **Author**: Michael Gardner (https://github.com/abitofhelp)
- **Company**: A Bit of Help, Inc.

---

**Last Updated**: November 25, 2025
