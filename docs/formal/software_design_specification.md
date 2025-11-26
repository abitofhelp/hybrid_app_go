# Software Design Specification (SDS)

**Project:** Hybrid_App_Go - Go 1.22+ Application Starter
**Version:** 1.0.0
**Date:** November 25, 2025
**SPDX-License-Identifier:** BSD-3-Clause
**License File:** See the LICENSE file in the project root.
**Copyright:** (c) 2025 Michael Gardner, A Bit of Help, Inc.
**Status:** Released

---

## 1. Introduction

### 1.1 Purpose

This Software Design Specification (SDS) describes the architectural design and detailed design of Hybrid_App_Go, a professional Go 1.22+ application starter demonstrating hexagonal architecture with functional programming principles.

### 1.2 Scope

This document covers:
- Architectural patterns and design decisions
- 5-layer organization and dependencies
- Key components and their responsibilities
- Data flow and error handling strategies
- Design patterns employed
- Static dependency injection implementation via generics

---

## 2. Architectural Design

### 2.1 Architecture Style

Hybrid_App_Go uses **Hexagonal Architecture** (Ports and Adapters / Clean Architecture).

**Benefits**:
- Clear separation of concerns
- Testable business logic (pure functions)
- Swappable infrastructure (adapters)
- Compiler-enforced boundaries
- Educational value for learning clean architecture

### 2.2 Layer Organization

```
+-------------------------------------------------------------+
|  Bootstrap                                                  |
|  (Composition Root - Wires Everything)                      |
|  - Generic instantiations with concrete types               |
|  - Port-to-adapter binding                                  |
|  - Application entry point                                  |
+-------------------------------------------------------------+
                          |
+-------------------------------------------------------------+
|  Presentation                                               |
|  (Driving Adapters - User Interfaces)                       |
|  - CLI commands                                             |
|  - Argument parsing                                         |
|  - Error message formatting                                 |
|  - Depends on: Application ONLY (not Domain)                |
+-------------------------------------------------------------+
                          |
+-------------------------------------------------------------+
|  Application                                                |
|  (Use Cases + Ports)                                        |
|  - Use case orchestration                                   |
|  - Inbound ports (use case interfaces)                      |
|  - Outbound ports (infrastructure interfaces)               |
|  - Commands (input DTOs)                                    |
|  - Models (output DTOs)                                     |
|  - application/error (re-exports domain/error)              |
|  - Depends on: Domain                                       |
+-------------------------------------------------------------+
                     |              ^
+----------------------+    +------------------------------+
|  Infrastructure      |    |  Domain                      |
|  (Driven Adapters)   |    |  (Business Logic)            |
|  - Console Writer    |    |  - Value Objects (Person)    |
|  - Adapts external   |    |  - Error types + Result[T]   |
|  - Panic -> Result   |    |  - Pure functions            |
|  - Depends on:       |    |  - ZERO dependencies         |
|    App + Domain      |    |                              |
+----------------------+    +------------------------------+
```

### 2.3 Layer Responsibilities

#### Domain Layer
- **Purpose**: Pure business logic, no external dependencies
- **Components**:
  - Value Objects: `domain/valueobject/person.go`
  - Error types: `domain/error/error.go`
  - Result monad: `domain/error/result.go`
- **Rules**:
  - Immutable value objects
  - Validation in constructors
  - Pure functions only (no side effects)
  - No infrastructure dependencies

#### Application Layer
- **Purpose**: Orchestrate domain logic, define port interfaces
- **Components**:
  - Use Cases: `application/usecase/greet.go`
  - Commands: `application/command/greet.go`
  - Models: `application/model/unit.go`
  - Inbound Ports: `application/port/inward/greet.go`
  - Outbound Ports: `application/port/outward/writer.go`
  - Error Re-export: `application/error/error.go`
- **Rules**:
  - Stateless use cases
  - Depends on Domain only
  - Defines interfaces for infrastructure
  - Generic over port types (static dispatch)

#### Infrastructure Layer
- **Purpose**: Implement technical concerns, adapt external systems
- **Components**:
  - Adapters: `infrastructure/adapter/consolewriter.go`
  - Implements outbound ports
- **Rules**:
  - Implements Application port interfaces
  - Catches panics at boundaries via defer/recover
  - Converts panics to Result errors
  - Depends on Application + Domain
  - Supports context.Context for cancellation

#### Presentation Layer
- **Purpose**: User interface implementation
- **Components**:
  - CLI Commands: `presentation/cli/command/greet.go`
  - Argument parsing
  - Error formatting
- **Rules**:
  - **Cannot access Domain directly**
  - Uses application/error (not domain/error)
  - Uses application/model (not Domain entities)
  - Depends on Application ONLY
  - Generic over use case types (static dispatch)

#### Bootstrap Layer
- **Purpose**: Composition root (dependency injection)
- **Components**:
  - `bootstrap/cli/cli.go` - wires everything together
  - Generic instantiations with concrete types
  - Port-adapter binding
- **Rules**:
  - Only layer that knows about all others
  - Static wiring (compile-time)
  - No business logic

---

## 3. Detailed Design

### 3.1 Domain Layer Design

#### Value Objects

**domain/valueobject/person.go**:
```go
type Person struct {
    name string  // private, immutable
}

func CreatePerson(name string) domerr.Result[Person] {
    // Validation
    if len(name) == 0 {
        return domerr.Err[Person](domerr.NewValidationError("name cannot be empty"))
    }
    if len(name) > MaxNameLength {
        return domerr.Err[Person](domerr.NewValidationError("name exceeds maximum length"))
    }
    return domerr.Ok(Person{name: name})
}

func (p Person) GetName() string { return p.name }
func (p Person) GreetingMessage() string { return "Hello, " + p.name + "!" }
```

**Design Decisions**:
- Immutable (unexported field, no setters)
- Validation in `CreatePerson` factory function
- Returns `Result[Person]`
- Pure methods (no side effects)

#### Error Handling

**domain/error/result.go**:
```go
type ErrorKind int

const (
    ValidationError ErrorKind = iota
    InfrastructureError
)

type ErrorType struct {
    Kind    ErrorKind
    Message string
}

type Result[T any] struct {
    value   T
    err     ErrorType
    isError bool
}

func Ok[T any](value T) Result[T]
func Err[T any](err ErrorType) Result[T]

func (r Result[T]) IsOk() bool
func (r Result[T]) IsError() bool
func (r Result[T]) Value() T         // Pre: IsOk()
func (r Result[T]) ErrorInfo() ErrorType  // Pre: IsError()
func (r Result[T]) UnwrapOr(defaultVal T) T
```

**Design Decisions**:
- Generic Result[T] for type-safe error handling
- Value semantics (not pointer receiver)
- No panics thrown from accessors (use preconditions)
- UnwrapOr for convenience with defaults

### 3.2 Application Layer Design

#### Use Cases

**application/usecase/greet.go** (Generic):
```go
// GreetUseCase is generic over WriterPort for static dispatch
type GreetUseCase[W outward.WriterPort] struct {
    writer W
}

func NewGreetUseCase[W outward.WriterPort](writer W) *GreetUseCase[W] {
    return &GreetUseCase[W]{writer: writer}
}

func (uc *GreetUseCase[W]) Execute(ctx context.Context, cmd command.GreetCommand) domerr.Result[model.Unit] {
    // 1. Extract name from DTO
    name := cmd.GetName()

    // 2. Validate and create Person (domain validation)
    personResult := valueobject.CreatePerson(name)
    if personResult.IsError() {
        return domerr.Err[model.Unit](personResult.ErrorInfo())
    }

    // 3. Generate greeting message
    person := personResult.Value()
    message := person.GreetingMessage()

    // 4. Write via output port (STATIC DISPATCH)
    return uc.writer.Write(ctx, message)
}
```

**Design Decisions**:
- Generic over WriterPort (static dispatch)
- Writer type W known at instantiation
- Railway-oriented error handling
- Pure orchestration (no business logic)
- Context support for cancellation

#### Application.Error Re-Export Pattern

**Problem**: Presentation cannot access Domain directly

**Solution**:
```go
// application/error/error.go
import domerr "github.com/abitofhelp/hybrid_app_go/domain/error"

// Type aliases - zero overhead
type ErrorType = domerr.ErrorType
type ErrorKind = domerr.ErrorKind
type Result[T any] = domerr.Result[T]

// Re-export error kinds
var (
    ValidationError     = domerr.ValidationError
    InfrastructureError = domerr.InfrastructureError
)
```

**Benefits**:
- Zero overhead (type aliases)
- Maintains boundary (Presentation -> Application only)
- Type-safe (compile-time verification)

### 3.3 Infrastructure Layer Design

**infrastructure/adapter/consolewriter.go**:
```go
type ConsoleWriter struct {
    w io.Writer
}

func NewConsoleWriter() *ConsoleWriter {
    return &ConsoleWriter{w: os.Stdout}
}

func NewWriter(w io.Writer) *ConsoleWriter {
    return &ConsoleWriter{w: w}
}

func (cw *ConsoleWriter) Write(ctx context.Context, message string) (result domerr.Result[model.Unit]) {
    // Panic recovery at infrastructure boundary
    defer func() {
        if r := recover(); r != nil {
            result = domerr.Err[model.Unit](
                domerr.NewInfrastructureError(fmt.Sprintf("panic recovered: %v", r)))
        }
    }()

    // Check context cancellation
    select {
    case <-ctx.Done():
        return domerr.Err[model.Unit](
            domerr.NewInfrastructureError("write cancelled: " + ctx.Err().Error()))
    default:
    }

    // Perform I/O
    _, err := fmt.Fprintln(cw.w, message)
    if err != nil {
        return domerr.Err[model.Unit](
            domerr.NewInfrastructureError("write failed: " + err.Error()))
    }

    return domerr.Ok(model.UnitValue)
}
```

**Design Decisions**:
- Panic recovery via defer/recover
- Context cancellation support
- Implements WriterPort interface
- Converts all errors to Result
- io.Writer for testability

### 3.4 Presentation Layer Design

**presentation/cli/command/greet.go** (Generic):
```go
// GreetCommand is generic over GreetPort for static dispatch
type GreetCommand[UC inward.GreetPort] struct {
    useCase UC
}

func NewGreetCommand[UC inward.GreetPort](useCase UC) *GreetCommand[UC] {
    return &GreetCommand[UC]{useCase: useCase}
}

func (c *GreetCommand[UC]) Run(args []string) int {
    // Parse arguments
    if len(args) != 2 {
        fmt.Fprintln(os.Stderr, "Usage: greeter <name>")
        return 1
    }

    // Create command DTO
    cmd := command.NewGreetCommand(args[1])
    ctx := context.Background()

    // Execute use case (STATIC DISPATCH)
    result := c.useCase.Execute(ctx, cmd)

    if result.IsError() {
        errInfo := result.ErrorInfo()
        switch errInfo.Kind {
        case apperr.ValidationError:
            fmt.Fprintf(os.Stderr, "Error: Please provide a valid name.\n")
        case apperr.InfrastructureError:
            fmt.Fprintf(os.Stderr, "Error: A system error occurred.\n")
        }
        return 1
    }

    return 0
}
```

**Design Decisions**:
- Generic over GreetPort (receives use case interface)
- Returns exit code (0=success, 1=error)
- Uses application/error (not domain/error)
- Pattern matches on ErrorKind for user-friendly messages

### 3.5 Bootstrap Design

**bootstrap/cli/cli.go**:
```go
func Run(args []string) int {
    // Step 1: Create Infrastructure adapter
    consoleWriter := adapter.NewConsoleWriter()

    // Step 2: Instantiate Use Case with concrete writer type
    // STATIC DISPATCH: GreetUseCase knows *ConsoleWriter at compile time
    greetUseCase := usecase.NewGreetUseCase[*adapter.ConsoleWriter](consoleWriter)

    // Step 3: Instantiate Command with concrete use case type
    // STATIC DISPATCH continues through the chain
    greetCommand := command.NewGreetCommand[*usecase.GreetUseCase[*adapter.ConsoleWriter]](greetUseCase)

    // Step 4: Run the application
    return greetCommand.Run(args)
}
```

**Design Decisions**:
- All wiring in one place
- Static instantiation (compile-time)
- Concrete type parameters at each level
- Clear dependency flow
- No runtime overhead

---

## 4. Design Patterns

### 4.1 Railway-Oriented Programming

**Pattern**: Result[T] monad for error handling
**Purpose**: Avoid panics, explicit error paths
**Implementation**: `domain/error/result.go`

```
Success Track:  Ok(Value) -> Continue -> Ok(Result)
Error Track:    Err(Error) -> Propagate -> Err(Error)
```

### 4.2 Hexagonal Architecture (Ports and Adapters)

**Pattern**: Decouple business logic from technical details
**Ports**: Interfaces defined in Application layer
**Adapters**: Implementations in Infrastructure/Presentation

### 4.3 Static Dependency Injection via Generics

**Pattern**: Generic types with interface constraints
**Wiring**: Bootstrap instantiates all generics with concrete types
**Benefits**: Compile-time resolution, zero overhead, type-safe

```go
// Generic definition with interface constraint
type GreetUseCase[W WriterPort] struct { writer W }

// Instantiation with concrete type
uc := NewGreetUseCase[*ConsoleWriter](writer)
```

### 4.4 Application Service Re-Export

**Pattern**: Facade pattern for layer boundaries
**Purpose**: Presentation cannot access Domain
**Implementation**: application/error re-exports domain/error

### 4.5 Value Object Pattern

**Pattern**: Immutable domain primitives with validation
**Implementation**: `domain/valueobject/person.go`
**Benefits**: Type safety, validated at construction, immutable

---

## 5. Data Flow

### 5.1 Request Flow (Success Path)

```
User: ./greeter Alice
    |
main (cmd/greeter/main.go): calls bootstrap.Run(os.Args)
    |
bootstrap.Run: wires dependencies (generic instantiation), calls Presentation
    |
presentation.GreetCommand.Run: parses args, creates command.GreetCommand DTO
    |
application.GreetUseCase.Execute: validates, orchestrates (STATIC DISPATCH)
    |
domain.CreatePerson("Alice"): validates -> Ok(Person)
    |
Person.GreetingMessage(): returns "Hello, Alice!"
    |
infrastructure.ConsoleWriter.Write: outputs "Hello, Alice!" (STATIC DISPATCH)
    |
Returns: Ok(Unit) -> exit code 0
```

### 5.2 Error Flow

```
User: ./greeter ""
    |
domain.CreatePerson(""): validates empty string
    |
Returns: Err(ValidationError, "name cannot be empty")
    |
application.GreetUseCase: checks IsError() -> propagates
    |
presentation.GreetCommand: pattern matches ErrorKind
    |
Displays: "Error: Please provide a valid name."
    |
Returns: exit code 1
```

### 5.3 Context Cancellation Flow

```
User: ./greeter Alice (then Ctrl+C)
    |
Context is cancelled
    |
infrastructure.ConsoleWriter.Write: checks ctx.Done()
    |
Returns: Err(InfrastructureError, "write cancelled: context canceled")
    |
Propagates up through layers
    |
Returns: exit code 1
```

---

## 6. Concurrency Design

### 6.1 Thread Safety

- **Domain**: Stateless, pure functions -> inherently thread-safe
- **Application**: Stateless use cases -> thread-safe
- **Infrastructure**: Uses context for cancellation -> goroutine-safe
- **Presentation**: Single-threaded CLI (current design)

### 6.2 Context Support

Context is passed through all layers for:
- Cancellation propagation
- Deadline support
- Request-scoped values (future)

---

## 7. Performance Design

### 7.1 Zero Overhead Abstractions

- **Static Dispatch**: Generics compiled to direct calls (no vtable)
- **Result Monad**: Value types (stack allocation)
- **Bounded Strings**: No unbounded allocation in Domain
- **Pure Functions**: Compiler can optimize aggressively

### 7.2 Memory Management

- **Stack Allocation**: Result values, commands, models
- **Minimal Heap**: Only for io.Writer in ConsoleWriter
- **Garbage Collection**: Standard Go GC (minimal pressure)

---

## 8. Security Design

### 8.1 Input Validation

- All validation in Domain layer
- Early rejection of invalid inputs
- Type-safe boundaries (compiler-enforced)

### 8.2 Error Information

- No sensitive data in error messages
- Structured error types
- Safe for display to users

---

## 9. Build and Deployment

### 9.1 Project Structure

```
hybrid_app_go/
|-- go.work            # Workspace file
|-- go.mod             # Root module
|-- Makefile           # Build automation
|-- cmd/greeter/       # Application entry point
|-- domain/            # Domain layer module
|-- application/       # Application layer module
|-- infrastructure/    # Infrastructure layer module
|-- presentation/      # Presentation layer module
|-- bootstrap/         # Bootstrap layer module
|-- test/integration/  # CLI integration tests
|-- docs/              # Documentation
```

**Design Decision**: go.work workspace with separate modules per layer
**Benefits**: Clear boundaries, independent testing, modular deployment

---

## 10. Testing Strategy

### 10.1 Test Organization

```
test/
|-- integration/          # CLI integration tests (23 tests)
    |-- greet_flow_test.go

domain/
|-- error/result_test.go  # Result monad unit tests (19 assertions)
|-- valueobject/person_test.go  # Person unit tests (23 assertions)
```

### 10.2 Testing Approach

- **Unit**: Test Domain in isolation (pure functions)
- **Integration**: Test entire CLI via binary execution
- **No Mocks**: Integration tests run the actual binary

---

## 11. Static Dispatch vs Dynamic Dispatch

### 11.1 Dynamic Dispatch (Traditional Go)

```go
type GreetUseCase struct {
    writer WriterPort  // interface type
}

func (uc *GreetUseCase) Execute(...) {
    uc.writer.Write(...)  // vtable lookup at runtime
}
```

### 11.2 Static Dispatch (This Project)

```go
type GreetUseCase[W WriterPort] struct {
    writer W  // concrete type parameter
}

func (uc *GreetUseCase[W]) Execute(...) {
    uc.writer.Write(...)  // direct call, no vtable
}
```

### 11.3 Comparison

| Aspect | Dynamic Dispatch | Static Dispatch |
|--------|------------------|-----------------|
| Method Resolution | Runtime (vtable) | Compile-time |
| Runtime Overhead | Interface conversion | Zero |
| Inlining | Not possible | Possible |
| Type Safety | Runtime checks | Compile-time |
| Flexibility | Runtime swapping | Fixed at compile |
| Binary Size | Smaller | Larger (monomorphization) |

---

**Document Control**:
- Version: 1.0.0
- Last Updated: November 25, 2025
- Status: Released
- Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
- License: BSD-3-Clause
- SPDX-License-Identifier: BSD-3-Clause
