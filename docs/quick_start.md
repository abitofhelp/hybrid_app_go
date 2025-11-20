# Hybrid_App_Ada Quick Start Guide

**Version:** 1.0.0
**Date:** November 18, 2025
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

- **GNAT Compiler**: GNAT FSF 13+ or GNAT Pro (Ada 2022 support required)
- **Alire**: Version 2.0+ (Ada package manager)
- **Java 11+**: For PlantUML diagram generation (optional)
- **Python 3**: For architecture validation and tooling

### Using Alire (Recommended)

```bash
# Clone the repository
git clone https://github.com/abitofhelp/hybrid_app_ada.git
cd hybrid_app_ada

# Build with Alire (automatically fetches dependencies)
alr build

# Or use make
make build
```

### Verify Installation

```bash
# Check that the executable was built
ls -lh bin/greeter

# Run the application
./bin/greeter World
# Output: Hello, World!
```

**Success!** You've built your first hexagonal architecture application in Ada 2022.

---

## First Build

The project uses both Alire and Make for building:

### Using Make (Recommended)

```bash
# Development build (with debug symbols)
make build

# Or explicit development mode
make build-dev

# Optimized build (O2)
make build-opt

# Release build
make build-release
```

### Using Alire Directly

```bash
# Development build
alr build --development

# Release build
alr build --release
```

**Build Output:**
- Executable: `bin/greeter`
- Object files: `obj/`
- Library artifacts: `lib/`

---

## Running the Application

The Hybrid_App_Ada starter includes a simple greeter application demonstrating all architectural layers:

### Basic Usage

```bash
# Greet a person
./bin/greeter Alice
# Output: Hello, Alice!

# Name with spaces (use quotes)
./bin/greeter "Bob Smith"
# Output: Hello, Bob Smith!

# Show usage
./bin/greeter
# Output: Usage: greeter <name>
# Exit code: 1
```

### Error Handling Example

```bash
# Empty name triggers validation error
./bin/greeter ""
# Output: Error: Name cannot be empty
# Exit code: 1
```

**Key Points:**
- All errors return via Result monad (no exceptions)
- Exit code 0 = success, 1 = error
- Validation happens in Domain layer
- Errors propagate through Application to Presentation

---

## Understanding the Architecture

Hybrid_App_Ada demonstrates **5-layer hexagonal architecture**:

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
2. **Presentation cannot access Domain** - Must use Application layer
3. **Static dependency injection** - Via generics (compile-time wiring)
4. **Railway-oriented programming** - Result monads for error handling
5. **Single-project structure** - Easy to deploy via Alire

### Request Flow Example

```
User Input ("Alice")
    ↓
Presentation.CLI.Command.Greet (parses input)
    ↓
Application.UseCase.Greet (validates via Domain)
    ↓
Domain.Value_Object.Person (business rules)
    ↓
Infrastructure.Adapter.Console_Writer (output)
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
# File: src/domain/value_object/domain-value_object-person.adb
```

### Step 2: Modify Greeting Format

Find the greeting logic and modify it. The business logic is pure and has no dependencies.

### Step 3: Rebuild and Test

```bash
# Rebuild
make rebuild

# Run tests to ensure nothing broke
make test-all

# Test manually
./bin/greeter Alice
```

**Best Practice**: Always run tests after making changes. This project has 82 tests ensuring correctness.

---

## Running Tests

Hybrid_App_Ada includes comprehensive testing:

### Test Organization

- **Unit Tests** (48 tests): Domain and Application logic
- **Integration Tests** (26 tests): Cross-layer interactions
- **E2E Tests** (8 tests): Full system via CLI

### Run All Tests

```bash
# Run entire test suite
make test-all

# Expected output:
# Running all test executables...
#
# ########################################
# ###                                  ###
# ###   ALL TEST SUITES: SUCCESS      ###
# ###   All tests passed!              ###
# ###                                  ###
# ########################################
```

### Run Specific Test Suites

```bash
# Unit tests only (fast)
make test-unit

# Integration tests only
make test-integration

# E2E tests only
make test-e2e
```

### Test Coverage

Hybrid_App_Ada supports code coverage analysis using GNATcoverage.

#### First-Time Setup

Before running coverage analysis for the first time, you need to build the GNATcoverage runtime library:

```bash
# Build the coverage runtime (one-time setup)
make build-coverage-runtime

# This will:
# - Locate your GNATcoverage installation
# - Build the runtime library from sources
# - Install it to external/gnatcov_rts/
```

**Note**: This step is automatically performed the first time you run `make test-coverage`, but you can run it explicitly if needed.

#### Running Coverage Analysis

```bash
# Run tests with GNATcoverage analysis
make test-coverage

# View coverage report
# Coverage reports generated in coverage/ directory
# - coverage/index.html - HTML coverage report
# - coverage/*.xcov - Detailed coverage files
```

#### Cleaning Coverage Data

```bash
# Remove coverage artifacts
make clean-coverage
```

**Coverage Runtime**: The runtime is built from your GNATcoverage installation sources and cached in `external/gnatcov_rts/`. This ensures reproducible builds across different environments.

**Test Framework**: Custom lightweight framework (no AUnit dependency) located in `test/common/test_framework.{ads,adb}`

---

## Build Targets

### Building

```bash
make build              # Development build (default)
make build-dev          # Explicit development mode
make build-opt          # Optimized build (O2)
make build-release      # Release build
make rebuild            # Clean and rebuild
```

### Testing

```bash
make test                    # Run all tests (alias for test-all)
make test-all                # Run entire test suite
make test-unit               # Unit tests only
make test-integration        # Integration tests only
make test-e2e                # E2E tests only
make test-coverage           # Tests with coverage analysis
make build-coverage-runtime  # Build GNATcoverage runtime (one-time setup)
```

### Quality & Architecture

```bash
make check              # Run static analysis
make check-arch         # Validate architecture boundaries
make stats              # Show project statistics
```

### Cleaning

```bash
make clean              # Clean build artifacts (fast rebuild)
make clean-deep         # Deep clean (includes dependencies - slow rebuild)
make clean-coverage     # Clean coverage data
make clean-clutter      # Remove temp files and backups
```

### Utilities

```bash
make deps               # Show dependency information
make prereqs            # Verify prerequisites
make refresh            # Refresh Alire dependencies
make compress           # Create source archive (tar.gz)
```

---

## Common Issues

### Q: Build fails with "Ada_2022 not supported"

**A:** You need GNAT FSF 13+ or GNAT Pro. Check your compiler version:

```bash
gnatls -v
```

### Q: "functional" dependency not found

**A:** Alire should fetch dependencies automatically. Try:

```bash
alr update
alr build
```

### Q: Architecture validation warnings appear

**A:** The `make check-arch` target validates layer boundaries. Warnings indicate potential violations:

```bash
# View architecture validation
make check-arch
```

These are warnings (not errors) - the build will still succeed.

### Q: Where are the test executables?

**A:** Test executables are in `test/bin/`:

```bash
ls -lh test/bin/
# Output:
# unit_runner
# integration_runner
# e2e_runner
```

### Q: How do I run a single test?

**A:** Execute the test runner directly:

```bash
# Run specific test runner
./test/bin/unit_runner
./test/bin/integration_runner
./test/bin/e2e_runner
```

---

## Next Steps

### Explore the Architecture

- **[Software Design Specification](formal/software_design_specification.md)** - Deep dive into architecture
- **[Architecture Diagrams](diagrams/)** - Visual documentation
- **[Layer Dependencies](diagrams/01_layer_dependencies.svg)** - See dependency flow

### Read the Source Code

Start with the wiring in Bootstrap:

```bash
# See how all layers are wired together
cat src/bootstrap/cli/bootstrap-cli.adb
```

Then explore each layer:

```bash
# Domain (pure business logic)
ls src/domain/

# Application (use cases and ports)
ls src/application/

# Infrastructure (adapters)
ls src/infrastructure/

# Presentation (CLI)
ls src/presentation/

# Bootstrap (composition root)
ls src/bootstrap/
```

### Study the Test Suite

```bash
# See how tests are organized
ls -R test/

# Read test framework
cat test/common/test_framework.ads
```

### Understand Error Handling

- **[Error Handling Strategy](guides/error_handling_strategy.md)** - Railway-oriented programming guide
- See how Result monads replace exceptions
- Study error propagation patterns

### Learn Dependency Injection

- **[Static vs Dynamic Dispatch](diagrams/05_static_vs_dynamic_dispatch.svg)** - Generic-based DI
- See Bootstrap wiring examples
- Understand compile-time polymorphism

### Add Your Own Use Case

Follow the pattern:

1. **Domain**: Create value objects/entities (`src/domain/`)
2. **Application**: Define command, use case, ports (`src/application/`)
3. **Infrastructure**: Implement adapters (`src/infrastructure/`)
4. **Presentation**: Create CLI command (`src/presentation/`)
5. **Bootstrap**: Wire everything together (`src/bootstrap/`)
6. **Tests**: Add unit/integration/e2e tests (`test/`)

### Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for:
- Code style guidelines
- Testing requirements
- Pull request process
- Architecture rules

---

## Documentation Index

- 📖 **[Main Documentation Hub](index.md)** - All documentation links
- 📋 **[Software Requirements Specification](formal/software_requirements_specification.md)** - Requirements
- 🏗️ **[Software Design Specification](formal/software_design_specification.md)** - Architecture
- 🧪 **[Software Test Guide](formal/software_test_guide.md)** - Testing guide
- 🗺️ **[Roadmap](roadmap.md)** - Future development plans

---

## Support

For questions, issues, or contributions:

- 📧 **Email**: support@abitofhelp.com
- 🐛 **Issues**: GitHub Issues
- 📖 **Documentation**: See `docs/` directory
- 💬 **Discussions**: GitHub Discussions

---

## License

Hybrid_App_Ada is licensed under the BSD-3-Clause License.
Copyright © 2025 Michael Gardner, A Bit of Help, Inc.

See [LICENSE](../LICENSE) for full license text.
