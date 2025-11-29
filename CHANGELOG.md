# Changelog

**Version:** 1.0.2
**Date:** November 28, 2025
**SPDX-License-Identifier:** BSD-3-Clause
**License File:** See the LICENSE file in the project root.
**Copyright:** (c) 2025 Michael Gardner, A Bit of Help, Inc.
**Status:** Released


All notable changes to Hybrid_App_Go will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Changed

### Added

### Removed

### Fixed

---

## [1.0.2] - 2025-11-28

### Changed
- Standardize README.md to new documentation format with proper section order
- Sync release scripts with hybrid_app_ada template

### Added
- AI Assistance & Authorship section in README.md (legally required disclosure)
- Embed version generation in Ada adapter (no external script needed)
- AI Assistance section validation (Step 0d) checks new placement rule
  (after Contributing, before License)

### Fixed
- Release script AI validation now checks correct section placement

---

## [1.0.1] - 2025-11-26

### Added

- **Documentation Validation**: New Step 0c in release scripts validates documentation consistency
  - Auto-detects project type (library vs application)
  - Checks for incorrect terminology (e.g., library terms in app projects)
  - Validates file references exist
  - Validates directory trees in documentation
  - Prompts user to review discrepancies before continuing

---

## [1.0.0] - 2025-11-26

_First stable release - Professional Go 1.23+ application starter template demonstrating hybrid DDD/Clean/Hexagonal architecture with functional programming principles._

### Added

#### Architecture
- **5-Layer Hexagonal Architecture**: Domain, Application, Infrastructure, Presentation, Bootstrap
- **Static Dispatch via Generics**: Zero-overhead dependency injection using Go generics
- **Port Abstraction**: Inbound (GreetPort) and Outbound (WriterPort) port interfaces
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

#### Release Tooling
- **Release Script Validation**: Makefile target and documentation link validation
- **brand_project.py**: Graceful handling of non-existent output directories
- **SPDX Headers**: All layer READMEs include BSD-3-Clause license identifier

### Changed

#### Port Naming Conventions
- Renamed `application/port/inward` to `application/port/inbound` (driving adapters)
- Renamed `application/port/outward` to `application/port/outbound` (driven adapters)
- Updated all imports, documentation, and PlantUML diagrams

#### Presentation Layer Structure
- Restructured `presentation/cli/command` to `presentation/adapter/cli/command`
- Consistent with adapter pattern used in infrastructure layer

#### Testing Improvements
- Added summary banner to integration tests (matching e2e test format)
- Integration tests now display pass/fail count in colored banner

#### Tooling
- Fixed arch_guard.py go.mod parsing to handle multi-line `require (...)` blocks
- Bootstrap layer now correctly shows all dependencies in architecture validation

#### Documentation
- Updated `docs/index.md` and `docs/quick_start.md` for Go project (was Ada)
- Updated README.md to reflect static dispatch via generics pattern
- All formal documentation and diagrams updated with current paths

### Fixed

- **.gitattributes**: Replaced placeholder prose with actual git configuration
- **architecture_mapping.md**: Corrected file paths to match current project structure
- **ports_mapping.md**: Corrected file paths to match current project structure

### Removed

- **APP_VS_LIB_ARCHITECTURE.md**: Not needed for application starter template

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

- **Go Version**: 1.23+ (generics with type inference)
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
