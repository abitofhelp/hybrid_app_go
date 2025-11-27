# Go Enterprise Starter with Hybrid DDD/Clean/Hexagonal Architecture

**Version:** 1.0.1  
**Date:** November 26, 2025  
**Copyright:** © 2025 Michael Gardner, A Bit of Help, Inc.  
**License:** BSD-3-Clause  

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

## Error Handling: Railway-Oriented Programming

**NO PANICS across layer boundaries.** All errors propagate via domain Result monad:

```go
// Domain defines custom Result[T] monad (ZERO external module dependencies)
import (
    domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
    "github.com/abitofhelp/hybrid_app_go/application/model"
    "github.com/abitofhelp/hybrid_app_go/domain/valueobject"
)

// Usage Pattern
func Execute(cmd GreetCommand) domerr.Result[model.Unit] {
    personResult := valueobject.CreatePerson(cmd.Name)

    if personResult.IsError() {
        return domerr.Err[model.Unit](personResult.ErrorInfo())
    }

    person := personResult.Value()
    return writer(person.GreetingMessage())
}
```

**Error Flow:**
1. **Domain:** Validates, returns `Err` variant if invalid
2. **Application:** Orchestrates, propagates errors upward
3. **Infrastructure:** Catches panics at boundaries, converts to `Err` via recovery pattern
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
testify v1.11.1    # Testing assertions (test module only, NOT domain)
```

**Note:** Domain layer has ZERO external module dependencies. Custom Result/Option monads are implemented in `domain/error/result.go` and `domain/valueobject/option.go`.

## Module Boundaries

**Enforced by go.mod dependencies:**

- **domain**: ZERO external module dependencies (custom Result/Option types)
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

### 4. Static Dispatch via Generics

**Pattern:** Generic types with interface constraints
**Wiring:** Bootstrap instantiates with concrete type parameters
**Benefit:** Compile-time resolution, methods devirtualized (zero runtime overhead)

### 5. Concurrency-Ready Pattern

This starter is **concurrency-ready** without implementing actual goroutines. The patterns are in place for when you need them:

**Context Propagation:**
```go
// Use case accepts context for cancellation/timeout
func (uc *GreetUseCase) Execute(ctx context.Context, cmd GreetCommand) domerr.Result[model.Unit]

// Infrastructure checks context before I/O
select {
case <-ctx.Done():
    return domerr.Err[model.Unit](apperr.NewInfrastructureError(
        fmt.Sprintf("write cancelled: %v", ctx.Err())))
default:
    // proceed with operation
}
```

**Panic Recovery at Boundaries:**
```go
// Infrastructure adapters recover panics and convert to Result errors
func (cw *ConsoleWriter) Write(ctx context.Context, message string) (result domerr.Result[model.Unit]) {
    defer func() {
        if r := recover(); r != nil {
            result = domerr.Err[model.Unit](domerr.NewInfrastructureError(
                fmt.Sprintf("write panicked: %v", r)))
        }
    }()
    // ... perform I/O
    return domerr.Ok(model.Unit{})
}
```

**When You Add Goroutines:**
- Pass `ctx` to all goroutines for cancellation signaling
- Use `ctx.Done()` channel in `select` statements
- Map `ctx.Err()` to `InfrastructureError` at boundaries
- No "spawn-and-forget" goroutines (always handle lifecycle)
- Use channels or `sync.WaitGroup` for coordination

**Example Extension (not in starter):**
```go
// Background monitor pattern (add when needed)
func StartMonitor(ctx context.Context, events chan<- Event) {
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()
        for {
            select {
            case <-ctx.Done():
                return // graceful shutdown
            case <-ticker.C:
                events <- checkHealth()
            }
        }
    }()
}
```

### 6. Go 1.23 Features

- **Workspaces** (`go.work` for multi-module projects)
- **Generics** (custom domain Result[T], Option[T] types)
- **Type parameters** (used in domain monads)

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
2. **Result Monads:** All fallible operations return `domerr.Result[T]`
3. **No Panics:** Errors are values, not thrown (recovery patterns for panic conversion)
4. **Module Boundaries:** Compiler-enforced via go.mod
5. **Static Dispatch:** Generic types with interface constraints for zero-overhead DI
6. **Table-Driven Tests:** Using testify assertions (test module, NOT domain)

## Comparison with Ada Version

| Aspect                  | Ada (Original)              | Go (This Port)                     |
|-------------------------|-----------------------------|------------------------------------|
| **Error Handling**      | Domain.Error.Result monad   | domain/error Result[T] monad       |
| **Dependency Injection**| Generic instantiation       | Static dispatch via generics       |
| **String Handling**     | Bounded strings             | Regular strings (GC handles it)    |
| **Memory Model**        | Stack allocation            | Stack + GC                         |
| **Polymorphism**        | Compile-time (generics)     | Compile-time (generics)            |
| **Module Boundaries**   | GPR project dependencies    | go.mod dependencies                |
| **Contracts**           | Pre/Post aspects            | Comments + assertions              |

## Project Status

✅ **Completed:**
- Multi-module workspace structure with go.work
- Custom domain Result/Option monads (ZERO external module dependencies)
- Static dispatch via generics (zero-overhead DI)
- Application.Error re-export pattern
- Module boundary enforcement via go.mod
- Comprehensive Makefile automation
- All layers ported from Ada to Go
- Functioning CLI application
- Context propagation for cancellation/timeout support
- Panic recovery at infrastructure boundaries
- Concurrency-ready patterns (documented, ready for extension)

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

**Prerequisites:**
- Python 3.8+
- Source template (this repository)
- Target directory (will be created)

The script auto-detects Go vs Ada projects and applies language-specific transformations.

## Learning Resources

- [Go Workspaces](https://go.dev/doc/tutorial/workspaces)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Railway-Oriented Programming](https://fsharpforfunandprofit.com/rop/)

## License

BSD-3-Clause - See LICENSE file in project root.

## Author

Michael Gardner  
A Bit of Help, Inc.  
https://github.com/abitofhelp
