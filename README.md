# Enterprise Starter with Hybrid DDD/Clean/Hexagonal Architecture

[![License](https://img.shields.io/badge/license-BSD--3--Clause-blue.svg)](LICENSE) [![Go](https://img.shields.io/badge/Go-1.23+-00ADD8.svg)](https://go.dev)

**Version:** 1.0.1<br>
**Date:** 2025-11-29<br>
**SPDX-License-Identifier:** BSD-3-Clause<br>
**License File:** See the LICENSE file in the project root<br>
**Copyright:** © 2025 Michael Gardner, A Bit of Help, Inc.<br>
**Status:** Released

## Overview

A **professional Go application** demonstrating **hybrid DDD/Clean/Hexagonal architecture** with **strict module boundaries** enforced via Go workspaces and **functional programming** principles using custom **domain-level Result/Option monads** (ZERO external module dependencies in domain layer).

> **Starter Template:** This project serves as a **starter template for enterprise Go application development**. Use the included `scripts/brand_project/brand_project.py` script to generate a new project from this template with your own project name, module paths, and branding. See [Creating a New Project](#creating-a-new-project) below.

This is a **desktop/enterprise application template** showcasing:
- **5-Layer Hexagonal Architecture** (Domain, Application, Infrastructure, Presentation, Bootstrap)
- **Strict Module Boundaries** via go.work and separate go.mod per layer
- **Static Dispatch via Generics** (zero-overhead dependency injection)
- **Railway-Oriented Programming** with Result monads (no panics across boundaries)
- **Presentation Isolation** pattern (only Domain is shareable across apps)
- **Multi-Module Workspace** (compiler-enforced boundaries)

## Features

- ✅ Multi-module workspace structure with go.work
- ✅ Custom domain Result/Option monads (ZERO external module dependencies)
- ✅ Static dispatch via generics (zero-overhead DI)
- ✅ Application.Error re-export pattern
- ✅ Module boundary enforcement via go.mod
- ✅ Context propagation for cancellation/timeout support
- ✅ Panic recovery at infrastructure boundaries
- ✅ Concurrency-ready patterns (documented, ready for extension)
- ✅ Comprehensive Makefile automation

## Architecture

### Module Structure

**Strict boundaries enforced by Go modules:**

```
hybrid_app_go/
├── go.work                          # Workspace definition (manages all modules)
├── domain/                          # Module: Pure business logic (ZERO external module dependencies)
│   └── go.mod                       # ZERO external module dependencies - custom Result/Option types
├── application/                     # Module: Use cases and ports
│   └── go.mod                       # Depends ONLY on domain
├── infrastructure/                  # Module: Driven adapters
│   └── go.mod                       # Depends on application + domain
├── presentation/                    # Module: Driving adapters (CLI)
│   └── go.mod                       # Depends ONLY on application (NOT domain)
├── bootstrap/                       # Module: Composition root
│   └── go.mod                       # Depends on ALL modules
└── cmd/greeter/                     # Module: Main entry point
    └── go.mod                       # Depends only on bootstrap
```

### Key Architectural Rules

**Critical Boundary Rule:**
> **Presentation is the ONLY outer layer prevented from direct Domain access**

- ✅ **Infrastructure** CAN access `domain/*` (implements repositories, uses entities)
- ✅ **Application** depends on `domain/*` (orchestrates domain logic)
- ❌ **Presentation** CANNOT access `domain/*` (must use `application/error`, `application/model`, etc.)

**Why This Matters:**
- Domain is the **only shareable layer** across multiple applications
- Each app has its own Application/Infrastructure/Presentation/Bootstrap
- Prevents tight coupling between UI and business logic
- Allows multiple UIs (CLI, REST, GUI) to share the same Domain

**The Solution:** `application/error` re-exports `domain/error` types (zero overhead type aliases)

### Dependency Injection Pattern

**Go (Static Dispatch via Generics)**:
```go
import (
    "context"
    domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
    "github.com/abitofhelp/hybrid_app_go/application/model"
    "github.com/abitofhelp/hybrid_app_go/application/port/outbound"
)

// Port interface defines the contract
type WriterPort interface {
    Write(ctx context.Context, message string) domerr.Result[model.Unit]
}

// Generic use case with interface constraint
type GreetUseCase[W outbound.WriterPort] struct {
    writer W
}

func NewGreetUseCase[W outbound.WriterPort](writer W) *GreetUseCase[W] {
    return &GreetUseCase[W]{writer: writer}
}

func (uc *GreetUseCase[W]) Execute(ctx context.Context, cmd GreetCommand) domerr.Result[model.Unit] {
    // uc.writer.Write() is statically dispatched - compiler knows exact type
}
```

**Wiring in Bootstrap:**
```go
// Step 1: Create Infrastructure adapter (concrete type)
consoleWriter := adapter.NewConsoleWriter()

// Step 2: Instantiate Use Case with concrete type parameter
greetUseCase := usecase.NewGreetUseCase[*adapter.ConsoleWriter](consoleWriter)

// Step 3: Instantiate Command with concrete use case type
greetCommand := command.NewGreetCommand[*usecase.GreetUseCase[*adapter.ConsoleWriter]](greetUseCase)

// Step 4: Run - all method calls are statically dispatched
return greetCommand.Run(os.Args)
```

**Benefits:**
- ✅ **Zero runtime overhead** (no vtable lookups, methods devirtualized)
- ✅ **Type-safe** (verified at compile time)
- ✅ **Static dispatch** (compiler knows exact types)
- ✅ **Inlining potential** (optimizer can inline method calls)

## Quick Start

### Prerequisites

- **Go 1.23+** (for workspace support)
- **golangci-lint** (optional, for linting)

### Building

```bash
# Build the project
make build

# Clean artifacts
make clean

# Rebuild from scratch
make rebuild
```

### Running

```bash
# Run the application
make run NAME="Alice"

# Or run directly
./bin/greeter Alice
```

## Usage

```bash
# Greet a person
./bin/greeter Alice
# Output: Hello, Alice!

# Name with spaces
./bin/greeter "Bob Smith"
# Output: Hello, Bob Smith!

# No arguments (shows usage)
./bin/greeter
# Output: Usage: greeter <name>
# Exit code: 1

# Empty name (validation error)
./bin/greeter ""
# Output: Error: Person name cannot be empty
# Exit code: 1
```

### Exit Codes

- **0**: Success
- **1**: Failure (validation error, infrastructure error, or missing arguments)

## Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Unit tests only
make test-unit
```

**Test Structure:**
- **Unit tests**: Co-located with code (`*_test.go`)
- **Integration tests**: `test/integration/` with `//go:build integration` tag
- **E2E tests**: `test/e2e/` with `//go:build e2e` tag

## Documentation

- 📚 **[Go Workspaces](https://go.dev/doc/tutorial/workspaces)** - Multi-module workspace tutorial
- 🏗️ **[Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)** - Architecture pattern
- 🚂 **[Railway-Oriented Programming](https://fsharpforfunandprofit.com/rop/)** - Error handling pattern

## Code Standards

This project follows:
- **Go Language Standards** (`~/.claude/agents/go.md`)
- **Architecture Standards** (`~/.claude/agents/architecture.md`)
- **Functional Programming Standards** (`~/.claude/agents/functional.md`)

### Key Standards Applied

1. **SPDX Headers:** All `.go` files have SPDX license headers
2. **Result Monads:** All fallible operations return `domerr.Result[T]`
3. **No Panics:** Errors are values, not thrown (recovery patterns for panic conversion)
4. **Module Boundaries:** Compiler-enforced via go.mod
5. **Static Dispatch:** Generic types with interface constraints for zero-overhead DI
6. **Table-Driven Tests:** Using testify assertions (test module, NOT domain)

## Creating a New Project

This repository serves as a **starter template** for enterprise Go applications. Use the `brand_project.py` script to create a new project with your own branding:

```bash
# From the scripts directory
cd scripts
python3 -m brand_project \
    --old-project hybrid_app_go \
    --new-project my_awesome_app \
    --old-org abitofhelp \
    --new-org mycompany \
    --source /path/to/hybrid_app_go \
    --target /path/to/my_awesome_app
```

**What gets updated:**
- Project name throughout all files
- GitHub organization/username in module paths
- Copyright holder information
- All `go.mod` module paths
- Import statements in Go source files
- Documentation and README files

## Submodule Management

This project uses git submodules for shared Python tooling:

- `scripts/python` - Build, release, and architecture scripts
- `test/python` - Shared test fixtures and configuration

### Workflow

```
hybrid_python_scripts (source repo)
         │
         │ git push (manual)
         ▼
      GitHub
         │
         │ make submodule-update (in each consuming repo)
         ▼
┌─────────────────────────────────┐
│  1. Pull new submodule commit   │
│  2. Stage reference change      │
│  3. Commit locally              │
│  4. Push to remote              │
└─────────────────────────────────┘
```

### Commands

```bash
# After fresh clone
make submodule-init

# Pull latest from submodule repos
make submodule-update

# Check current submodule commits
make submodule-status
```

### Bulk Update (all repositories)

```bash
python3 ~/Python/src/github.com/abitofhelp/git/update_submodules.py

# Options:
#   --dry-run   Show what would happen without changes
#   --no-push   Update locally but do not push to remote
```

## Contributing

This project is not open to external contributions at this time.

## AI Assistance & Authorship

This project — including its source code, tests, documentation, and other deliverables — is designed, implemented, and maintained by human developers, with Michael Gardner as the Principal Software Engineer and project lead.

We use AI coding assistants (such as OpenAI GPT models and Anthropic Claude Code) as part of the development workflow to help with:

- drafting and refactoring code and tests,
- exploring design and implementation alternatives,
- generating or refining documentation and examples,
- and performing tedious and error-prone chores.

AI systems are treated as tools, not authors. All changes are reviewed, adapted, and integrated by the human maintainers, who remain fully responsible for the architecture, correctness, and licensing of this project.

## License

Copyright © 2025 Michael Gardner, A Bit of Help, Inc.

Licensed under the BSD-3-Clause License. See [LICENSE](LICENSE) for details.

## Author

Michael Gardner
A Bit of Help, Inc.
https://github.com/abitofhelp

## Project Status

**Status**: Production Ready (v1.0.1)

- ✅ Multi-module workspace structure with go.work
- ✅ Custom domain Result/Option monads (ZERO external module dependencies)
- ✅ Static dispatch via generics (zero-overhead DI)
- ✅ Application.Error re-export pattern
- ✅ Module boundary enforcement via go.mod
- ✅ Comprehensive Makefile automation
- ✅ All layers ported from Ada to Go
- ✅ Functioning CLI application
- ✅ Context propagation for cancellation/timeout support
- ✅ Panic recovery at infrastructure boundaries
- ✅ Concurrency-ready patterns (documented, ready for extension)
