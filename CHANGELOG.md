# Changelog

**Project:** Hybrid_App_Go - Go 1.22+ Application Starter
**SPDX-License-Identifier:** BSD-3-Clause
**License File:** See the LICENSE file in the project root.
**Copyright:** (c) 2025 Michael Gardner, A Bit of Help, Inc.

All notable changes to Hybrid_App_Go will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2025-11-25

_First stable release - Professional Go 1.22+ application starter template demonstrating hybrid DDD/Clean/Hexagonal architecture with functional programming principles._

### Added

#### Architecture
- **5-Layer Hexagonal Architecture**: Domain, Application, Infrastructure, Presentation, Bootstrap
- **Static Dispatch via Generics**: Zero-overhead dependency injection using Go generics
- **Port Abstraction**: Inward (GreetPort) and Outward (WriterPort) port interfaces
- **Architecture Guard**: Python script (arch_guard.py) for automated boundary validation
- **Context Support**: context.Context for cancellation and deadline propagation

#### Domain Layer
- **Result[T] Monad**: Railway-oriented programming for explicit error handling
  - `Ok[T]()`/`Err[T]()` constructors
  - `IsOk()`/`IsError()` predicates
  - `Value()`/`ErrorInfo()` accessors with precondition checks
  - `UnwrapOr()`, `UnwrapOrElse()`, `Expect()` safe alternatives
  - `Map()`, `MapTo()`, `AndThen()`, `AndThenTo()` functional operations
  - `MapError()`, `Fallback()`, `FallbackWith()` error handling
  - `Recover()`, `RecoverWith()` error recovery
  - `Tap()` for side effects
- **Person Value Object**: Immutable with validation via smart constructor
- **Error Types**: ValidationError and InfrastructureError with descriptive messages

#### Application Layer
- **GreetUseCase[W WriterPort]**: Generic use case with static dispatch
- **GreetCommand DTO**: Immutable command for crossing layer boundaries
- **Unit Model**: Type-safe void/unit representation
- **Error Re-export**: application/error re-exports domain/error for Presentation layer

#### Infrastructure Layer
- **ConsoleWriter**: Adapter implementing WriterPort
  - Panic recovery via defer/recover
  - Context cancellation support
  - io.Writer injection for testability
- `NewWriter(io.Writer)`: Factory for custom output targets
- `NewConsoleWriter()`: Convenience factory for stdout
- `NewStderrWriter()`: Factory for stderr output

#### Presentation Layer
- **GreetCommand[UC GreetPort]**: Generic CLI command with static dispatch
- **Exit code mapping**: 0 for success, 1 for error
- **User-friendly error messages**: Pattern matching on ErrorKind

#### Bootstrap Layer
- **Composition Root**: Single `Run()` function wiring all dependencies
- **Generic Instantiation**: Concrete types resolved at compile time
- **Static Dispatch Chain**: All method calls devirtualized

#### Testing
- **Domain Unit Tests**: 42 assertions in 2 test functions
  - Result[T] monad tests (19 assertions)
  - Person value object tests (23 assertions)
- **CLI Integration Tests**: 23 tests running actual binary
  - Valid input scenarios
  - Invalid input scenarios (empty, too long, wrong arg count)
  - Edge cases (whitespace, unicode, special characters)
  - Table-driven tests for comprehensive coverage
- **Test Infrastructure**: testify assertions, colored output helpers

#### Documentation
- **Formal Documentation**:
  - Software Requirements Specification (SRS)
  - Software Design Specification (SDS)
  - Software Test Guide (STG)
- **PlantUML Diagrams** (5 diagrams with SVG):
  - `layer_dependencies.puml` - 5-layer architecture
  - `package_structure.puml` - Package organization
  - `error_handling_flow.puml` - Railway-oriented flow
  - `static_dispatch.puml` - Static vs dynamic dispatch
  - `application_error_pattern.puml` - Error re-export pattern
- **Guides**:
  - Quick Start Guide
  - Architecture Mapping Guide
  - Ports Mapping Guide
- **Source Code Documentation**: Comprehensive godoc for all packages

#### Build System
- **Makefile**: Comprehensive build automation
- **go.work**: Workspace for multi-module project
- **Separate go.mod per layer**: Clear dependency boundaries

### Architecture Patterns

- **Static Dependency Injection**: Generic structs with interface constraints (compile-time DI)
- **Result Monad**: Railway-oriented error handling (no panics across boundaries)
- **Presentation Isolation**: Presentation uses application/error, not domain/error
- **Minimal Entry Point**: main() delegates to bootstrap.Run() (1 line)
- **Ports & Adapters**: Application defines ports, Infrastructure implements adapters

### Technical Details

```go
// Static dispatch - concrete types known at compile time
consoleWriter := adapter.NewConsoleWriter()
uc := usecase.NewGreetUseCase[*adapter.ConsoleWriter](consoleWriter)
cmd := command.NewGreetCommand[*usecase.GreetUseCase[*adapter.ConsoleWriter]](uc)

// All calls are statically dispatched - no vtable lookup
exitCode := cmd.Run(args)
```

```go
// Context support for cancellation
func (uc *GreetUseCase[W]) Execute(ctx context.Context, cmd command.GreetCommand) domerr.Result[model.Unit] {
    // Context flows through to infrastructure
    return uc.writer.Write(ctx, message)
}
```

```go
// Railway-oriented error handling
result := valueobject.CreatePerson(name)
if result.IsError() {
    return domerr.Err[model.Unit](result.ErrorInfo())
}
person := result.Value()
```

### Compatibility

- **Go Version**: 1.22+ (generics with type inference)
- **Platforms**: Linux, macOS, BSD, Windows
- **Dependencies**: testify (test only)

### Notes

This is the initial release of Hybrid_App_Go, a Go port of the hybrid_app_ada project.
Both projects demonstrate identical architectural patterns:
- 5-layer hexagonal architecture
- Static dispatch for dependency injection
- Railway-oriented programming with Result monads
- Clean architecture boundary enforcement

---

## Release Notes Format

Each release will document changes in these categories:

- **Added** - New features
- **Changed** - Changes to existing functionality
- **Deprecated** - Soon-to-be-removed features
- **Removed** - Removed features
- **Fixed** - Bug fixes
- **Security** - Security vulnerability fixes

---

## License & Copyright

- **License**: BSD-3-Clause
- **Copyright**: (c) 2025 Michael Gardner, A Bit of Help, Inc.
- **SPDX-License-Identifier**: BSD-3-Clause

[1.0.0]: https://github.com/abitofhelp/hybrid_app_go/releases/tag/v1.0.0
