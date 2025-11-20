# Hybrid_App_Ada Documentation Index

**Version:** 1.0.0
**Date:** November 18, 2025
**SPDX-License-Identifier:** BSD-3-Clause
**License File:** See the LICENSE file in the project root.
**Copyright:** © 2025 Michael Gardner, A Bit of Help, Inc.
**Status:** Released

---

## Welcome

Welcome to the **Hybrid_App_Ada** documentation. This Ada 2022 application starter demonstrates professional hexagonal architecture with functional programming principles, static dependency injection, and railway-oriented error handling.

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
  - Running tests (make test-unit, test-integration, test-all)
  - Test framework documentation
  - Coverage procedures
  - Writing new tests

### Development Planning

- 🗺️ **[Roadmap](roadmap.md)** - Future development plans
  - v1.0.0 achievements
  - v1.1.0 planned features (Q1 2026)
  - v1.2.0 advanced patterns (Q2 2026)
  - v2.0.0 enterprise features (Q3 2026)
  - Contributing to roadmap

---

## Architecture Overview

Hybrid_App_Ada implements a **5-layer hexagonal architecture** (also known as Ports and Adapters or Clean Architecture):

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
2. **Presentation Boundary**: Presentation layer cannot access Domain directly (uses Application.Error, Application.Model re-exports)
3. **Static Dispatch**: Dependency injection via generics (compile-time, zero overhead)
4. **Railway-Oriented**: Result monads for error handling (no exceptions across boundaries)
5. **Single Project**: No aggregate projects, easy Alire deployment

---

## Visual Documentation

### UML Diagrams

Located in `diagrams/` directory:

- **01_layer_dependencies.svg** - Shows 5-layer dependency flow
- **02_application_error_pattern.svg** - Re-export pattern for Presentation isolation
- **03_package_structure.svg** - Actual package hierarchy
- **04_error_handling_flow.svg** - Error propagation through layers
- **05_static_vs_dynamic_dispatch.svg** - Generic vs interface comparison

All diagrams are generated from PlantUML sources (.puml files).

---

## Project Statistics

### Code Metrics (v1.0.0)

- **Ada Specification Files**: 26 (.ads)
- **Ada Implementation Files**: 11 (.adb)
- **Test Files**: 11 (5 unit, 3 integration, 2 e2e, 1 runner)
- **Architecture Layers**: 5 (Domain, Application, Infrastructure, Presentation, Bootstrap)
- **Build Targets**: 30+ Makefile targets
- **Dependencies**: functional ^1.0.0, aunit ^24.0.0, gnatcov ^22.0.1

### Test Coverage

- **Total Tests**: 82 (100% passing)
  - Unit Tests: 48
  - Integration Tests: 26
  - E2E Tests: 8
- **Test Framework**: Custom lightweight framework (test/common/test_framework.{ads,adb})
- **Coverage Analysis**: GNATcoverage support

### Code Quality

- **Compiler Warnings**: 0
- **Style Violations**: 0
- **Architecture Validation**: Enforced by arch_guard.py
- **Aspect Syntax**: 100% (no obsolescent pragmas)
- **Ada Version**: Ada 2022

---

## Key Features

### Static Dependency Injection

Uses **generics** instead of interfaces for dependency injection:

```ada
-- Port definition (generic function signature)
generic
   with function Writer (Message : String) return Result;
package Application.Usecase.Greet is
   function Execute (...) return Result;
end Application.Usecase.Greet;

-- Wiring in Bootstrap (compile-time resolution)
package Greet_UseCase is new Application.Usecase.Greet
  (Writer => Console_Writer.Write);
```

**Benefits**:
- Zero runtime overhead (no vtable lookups)
- Full inlining potential
- Stack allocation (no heap)
- Compile-time type safety

**Trade-off**: Fixed at compile time (perfect for most applications)

### Railway-Oriented Programming

All errors propagate via **Result monad** (no exceptions):

```ada
-- Use case returns Result[Unit, Error]
function Execute (Cmd : Greet_Command) return Unit_Result is
   Person_Result : constant Person_Result.Result :=
      Domain.Value_Object.Person.Create(Cmd.Name);
begin
   if Person_Result.Is_Error then
      return Unit_Result.Error(Person_Result.Error_Info);
   end if;

   -- Continue on success track...
   return Writer(Greeting_Message(Person_Result.Value));
end Execute;
```

**Benefits**:
- Explicit error handling (compiler-enforced)
- No unexpected control flow
- Composable error types
- SPARK-verifiable

### Application.Error Re-Export Pattern

**Problem**: Presentation cannot access Domain directly
**Solution**: Application re-exports Domain.Error for Presentation

```ada
-- Application.Error (zero-overhead facade)
package Application.Error is
   subtype Error_Type is Domain.Error.Error_Type;
   subtype Error_Kind is Domain.Error.Error_Kind;
   package Error_Strings renames Domain.Error.Error_Strings;
end Application.Error;
```

This maintains clean boundaries while allowing Presentation to handle errors.

---

## Directory Structure

```
hybrid_app_ada/
├── src/
│   ├── domain/                 # Pure business logic
│   │   ├── error/              # Result monad, error types
│   │   └── value_object/       # Immutable value objects
│   ├── application/            # Use cases + ports
│   │   ├── command/            # Input DTOs
│   │   ├── error/              # Re-exports for Presentation
│   │   ├── model/              # Output DTOs
│   │   ├── port/               # Port interfaces (inward/outward)
│   │   └── usecase/            # Use case orchestration
│   ├── infrastructure/         # Adapters (driven)
│   │   └── adapter/            # Console writer, etc.
│   ├── presentation/           # Adapters (driving)
│   │   └── cli/                # CLI commands
│   ├── bootstrap/              # Composition root
│   │   └── cli/                # CLI wiring
│   └── greeter.adb             # Main (3 lines)
├── test/
│   ├── unit/                   # Domain + Application tests
│   ├── integration/            # Cross-layer tests
│   ├── e2e/                    # Full system tests
│   └── common/                 # Test framework
├── docs/
│   ├── formal/                 # SRS, SDS, Test Guide
│   ├── diagrams/               # UML diagrams
│   ├── quick_start.md          # Getting started
│   ├── roadmap.md              # Future plans
│   └── index.md                # This file
├── scripts/
│   └── arch_guard.py           # Architecture validation
├── hybrid_app_ada.gpr           # Main project file
├── alire.toml                  # Alire manifest
├── Makefile                    # Build automation
└── README.md                   # Project overview
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
make test-all           # Run all tests
make test-unit          # Unit tests only
make test-integration   # Integration tests
make test-e2e           # E2E tests
make test-coverage      # With coverage analysis
```

**Quality**:
```bash
make check              # Static analysis
make check-arch         # Architecture validation
make stats              # Project statistics
```

**Utilities**:
```bash
make clean              # Clean artifacts
make deps               # Show dependencies
make help               # Show all targets
```

See [Quick Start Guide](quick_start.md#build-targets) for complete list.

---

## Learning Path

### For Beginners

1. **Start Here**: [Quick Start Guide](quick_start.md)
2. **Understand Architecture**: [Software Design Specification](formal/software_design_specification.md)
3. **Run Tests**: `make test-all`
4. **Explore Code**: Start with `src/bootstrap/cli/bootstrap-cli.adb`
5. **Read Examples**: Study how layers are wired together

### For Experienced Developers

1. **Architecture Patterns**: See [SDS - Design Patterns](formal/software_design_specification.md#design-patterns)
2. **Static DI Deep Dive**: See diagrams/05_static_vs_dynamic_dispatch.svg
3. **Railway-Oriented Programming**: See diagrams/04_error_handling_flow.svg
4. **Add Use Case**: Follow pattern in existing code
5. **Contribute**: See roadmap for future enhancements

---

## Dependencies

### Runtime Dependencies

- **functional** (^1.0.0): Result/Option/Either monads for functional error handling

### Development Dependencies

- **aunit** (^24.0.0): Unit testing framework (not used - we have custom framework)
- **gnatcov** (^22.0.1): Coverage analysis tool

### Build Requirements

- **GNAT**: FSF 13+ or GNAT Pro (Ada 2022 support)
- **Alire**: 2.0+ package manager
- **Make**: GNU Make for build automation
- **Python 3**: For tooling (arch_guard.py, etc.)
- **Java 11+**: For PlantUML diagram generation (optional)

---

## Documentation Updates

All documentation is maintained for v1.0.0 release:

- **Copyright**: © 2025 Michael Gardner, A Bit of Help, Inc.
- **License**: BSD-3-Clause
- **Version**: 1.0.0
- **Date**: November 18, 2025
- **Status**: Released

For documentation issues or suggestions, please file an issue on GitHub.

---

## Support and Contributing

### Getting Help

- 📧 **Email**: support@abitofhelp.com
- 🐛 **Issues**: [GitHub Issues](https://github.com/abitofhelp/hybrid_app_ada/issues)
- 📖 **Documentation**: This directory
- 💬 **Discussions**: [GitHub Discussions](https://github.com/abitofhelp/hybrid_app_ada/discussions)

### Contributing

We welcome contributions! See:

- [CONTRIBUTING.md](../CONTRIBUTING.md) - Contribution guidelines
- [Roadmap](roadmap.md) - Future development plans
- Code style enforced by architecture validation

---

## License

Hybrid_App_Ada is licensed under the **BSD-3-Clause License**.

Copyright © 2025 Michael Gardner, A Bit of Help, Inc.

See [LICENSE](../LICENSE) for full license text.

---

## Project Links

- **GitHub**: https://github.com/abitofhelp/hybrid_app_ada
- **Alire**: https://alire.ada.dev (search for "hybrid_app_ada")
- **Author**: Michael Gardner (https://github.com/abitofhelp)
- **Company**: A Bit of Help, Inc. (https://abitofhelp.com)

---

**Last Updated**: November 18, 2025
