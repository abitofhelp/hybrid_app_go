# hybrid_app_go - Go Application with Strict Module Boundaries

**Version:** 1.0.0  
**Date:** November 20, 2025  
**Copyright:** © 2025 Michael Gardner, A Bit of Help, Inc.  
**License:** BSD-3-Clause

## Overview

A **professional Go application** demonstrating **hybrid DDD/Clean/Hexagonal architecture** with **strict module boundaries** enforced via Go workspaces and **functional programming** principles using the `samber/mo` library for Result monads.

This is a **desktop/enterprise application template** showcasing:
- **5-Layer Hexagonal Architecture** (Domain, Application, Infrastructure, Presentation, Bootstrap)
- **Strict Module Boundaries** via go.work and separate go.mod per layer
- **Function Injection Dependency Injection** (lightweight Go pattern)
- **Railway-Oriented Programming** with Result monads (no panics across boundaries)
- **Presentation Isolation** pattern (only Domain is shareable across apps)
- **Multi-Module Workspace** (compiler-enforced boundaries)

## Architecture

### Module Structure

**Strict boundaries enforced by Go modules:**

```
hybrid_app_go/
├── go.work                          # Workspace definition (manages all modules)
├── domain/                          # Module: Pure business logic (ZERO dependencies)
│   └── go.mod                       # Depends only on samber/mo
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

**Go (Function Injection)**:
```go
type WriterFunc func(message string) mo.Result[Unit]

type GreetUseCase struct {
    writer WriterFunc
}

func NewGreetUseCase(writer WriterFunc) *GreetUseCase {
    return &GreetUseCase{writer: writer}
}

func (uc *GreetUseCase) Execute(cmd GreetCommand) mo.Result[Unit] {
    // Use uc.writer(message)
}
```

**Wiring in Bootstrap:**
```go
// Step 1: Wire Infrastructure → Port
consoleWriter := adapter.NewConsoleWriter()

// Step 2: Wire Use Case → Port
greetUseCase := usecase.NewGreetUseCase(consoleWriter)

// Step 3: Wire Command → Use Case
greetCommand := command.NewGreetCommand(greetUseCase.Execute)

// Step 4: Run
return greetCommand.Run(os.Args)
```

**Benefits:**
- ✅ **Zero runtime overhead** (direct function calls)
- ✅ **Type-safe** (verified at compile time)
- ✅ **Functional composition** (functions passed as dependencies)
- ✅ **Lightweight** (no reflection, no interfaces needed)

## Error Handling: Railway-Oriented Programming

**NO PANICS across layer boundaries.** All errors propagate via mo.Result monad:

```go
// Domain defines Result[T] monad using samber/mo
import "github.com/samber/mo"

// Usage Pattern
func Execute(cmd GreetCommand) mo.Result[Unit] {
    personResult := valueobject.CreatePerson(cmd.Name)
    
    if personResult.IsError() {
        return mo.Err[Unit](personResult.Error())
    }
    
    person := personResult.MustGet()
    return writer(person.GreetingMessage())
}
```

**Error Flow:**
1. **Domain:** Validates, returns `Err` variant if invalid
2. **Application:** Orchestrates, propagates errors upward
3. **Infrastructure:** Catches panics at boundaries, converts to `Err` with `mo.Try`
4. **Presentation:** Pattern matches on `ErrorKind`, displays user-friendly messages

## Building

### Prerequisites

- **Go 1.23+** (for workspace support)
- **golangci-lint** (optional, for linting)

### Build Commands

```bash
# Build the project
make build

# Run the application
make run NAME="Alice"

# Run specific targets
./bin/greeter Alice

# Clean artifacts
make clean

# Run tests
make test

# Format code
make fmt

# Run linter
make lint
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

## Exit Codes

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

## Dependencies

Managed by Go modules (`go.mod` per module):

```
samber/mo v1.13.0  # Result/Option/Either monads
testify v1.11.1    # Testing assertions
```

## Module Boundaries

**Enforced by go.mod dependencies:**

- **domain**: NO dependencies (except mo)
- **application**: domain ONLY
- **infrastructure**: application + domain
- **presentation**: application ONLY (NOT domain)
- **bootstrap**: ALL modules
- **cmd/greeter**: bootstrap ONLY

**Compiler enforces these rules** - attempting to import forbidden packages results in build errors.

## Key Design Patterns

### 1. Minimal Entry Point

**Main (cmd/greeter/main.go)** - Only 3 lines:

```go
func main() {
    exitCode := cli.Run(os.Args)
    os.Exit(exitCode)
}
```

### 2. Result Monad Pattern

**Railway-Oriented Programming:**
- Ok track: Successful computation continues
- Error track: Error propagates (short-circuit)
- Forces explicit error handling at compile time
- No panics thrown across boundaries

### 3. Application.Error Re-export Pattern

**Problem:** Presentation cannot access Domain directly  
**Solution:** Application re-exports Domain types for Presentation use  
**Implementation:** Type aliases and variable re-exports (zero overhead)

### 4. Function Injection

**Pattern:** Functions passed as dependencies  
**Wiring:** Bootstrap injects all functions  
**Benefit:** Compile-time resolution (zero runtime cost)

### 5. Go 1.23 Features

- **Workspaces** (`go.work` for multi-module projects)
- **Generics** (via samber/mo for Result[T], Option[T])
- **Type parameters** (used in mo library)

## Workspace Management

This project uses Go workspaces to manage multiple modules:

```bash
# Sync workspace (after pulling changes)
go work sync

# Add a new module to workspace
go work use ./new-module

# Check workspace status
go work edit -print
```

## Standards Compliance

This project follows:
- **Go Language Standards** (`~/.claude/agents/go.md`)
- **Architecture Standards** (`~/.claude/agents/architecture.md`)
- **Functional Programming Standards** (`~/.claude/agents/functional.md`)

### Key Standards Applied

1. **SPDX Headers:** All `.go` files have SPDX license headers
2. **Result Monads:** All fallible operations return `mo.Result[T]`
3. **No Panics:** Errors are values, not thrown (mo.Try for panic conversion)
4. **Module Boundaries:** Compiler-enforced via go.mod
5. **Function Injection:** Lightweight dependency injection
6. **Table-Driven Tests:** Using testify assertions

## Comparison with Ada Version

| Aspect                  | Ada (Original)              | Go (This Port)                     |
|-------------------------|-----------------------------|------------------------------------|
| **Error Handling**      | Domain.Error.Result monad   | mo.Result[T] monad                 |
| **Dependency Injection**| Generic instantiation       | Function injection                 |
| **String Handling**     | Bounded strings             | Regular strings (GC handles it)    |
| **Memory Model**        | Stack allocation            | Stack + GC                         |
| **Polymorphism**        | Compile-time (generics)     | Compile-time (function types)      |
| **Module Boundaries**   | GPR project dependencies    | go.mod dependencies                |
| **Contracts**           | Pre/Post aspects            | Comments + assertions              |

## Project Status

✅ **Completed:**
- Multi-module workspace structure with go.work
- Result monad error handling (mo.Result[T])
- Function injection dependency injection
- Application.Error re-export pattern
- Module boundary enforcement via go.mod
- Comprehensive Makefile automation
- All layers ported from Ada to Go
- Functioning CLI application

## Learning Resources

- [Go Workspaces](https://go.dev/doc/tutorial/workspaces)
- [samber/mo Library](https://github.com/samber/mo)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Railway-Oriented Programming](https://fsharpforfunandprofit.com/rop/)

## License

BSD-3-Clause - See LICENSE file in project root.

## Author

Michael Gardner  
A Bit of Help, Inc.  
https://github.com/abitofhelp
