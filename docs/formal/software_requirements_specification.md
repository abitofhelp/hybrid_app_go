# Software Requirements Specification (SRS)

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

This Software Requirements Specification (SRS) describes the functional and non-functional requirements for Hybrid_App_Go, a professional Go 1.22+ application starter template demonstrating hexagonal architecture with functional programming principles.

### 1.2 Scope

Hybrid_App_Go provides:
- Professional 5-layer hexagonal architecture implementation
- Static dependency injection via Go generics
- Railway-oriented programming with Result[T] monads
- Clean architecture boundary enforcement
- Integration test suite for CLI verification
- Domain unit tests with comprehensive assertions
- Production-ready code quality standards
- Educational documentation and examples

### 1.3 Definitions and Acronyms

- **DDD**: Domain-Driven Design
- **SRS**: Software Requirements Specification
- **SDS**: Software Design Specification
- **CLI**: Command Line Interface
- **DI**: Dependency Injection
- **Result Monad**: Functional programming error handling pattern
- **Hexagonal Architecture**: Also known as Ports and Adapters or Clean Architecture
- **Static Dispatch**: Compile-time method resolution (no vtable)

### 1.4 References

- Go 1.22 Language Specification
- Clean Architecture (Robert C. Martin)
- Domain-Driven Design (Eric Evans)
- Railway-Oriented Programming (Scott Wlaschin)
- Hexagonal Architecture (Alistair Cockburn)
- Companion project: hybrid_app_ada (Ada 2022 implementation)

---

## 2. Overall Description

### 2.1 Product Perspective

Hybrid_App_Go is a standalone application starter template implementing professional architectural patterns:

**Architecture Layers**:
1. **Domain**: Pure business logic (zero dependencies)
2. **Application**: Use cases and port definitions
3. **Infrastructure**: Driven adapters (implementations)
4. **Presentation**: Driving adapters (user interfaces)
5. **Bootstrap**: Composition root (dependency wiring)

### 2.2 Product Features

1. **Hexagonal Architecture**: 5-layer clean architecture
2. **Static Dependency Injection**: Compile-time wiring via generics
3. **Railway-Oriented Programming**: Result[T] monad error handling
4. **Architecture Enforcement**: Automated boundary validation (arch_guard.py)
5. **Test Infrastructure**: Domain unit tests + CLI integration tests
6. **Build Automation**: Comprehensive Makefile
7. **Documentation**: Complete SRS, SDS, Test Guide

### 2.3 User Classes

- **Application Developers**: Learn hexagonal architecture patterns
- **Team Leads**: Adopt architectural standards
- **Educators**: Teach clean architecture principles
- **Go Developers**: Start new projects with best practices

### 2.4 Operating Environment

- **Platforms**: Linux, macOS, BSD, Windows
- **Go Version**: Go 1.22+ (generics with type inference)
- **Build System**: go build, make
- **Testing**: go test, testify

---

## 3. Functional Requirements

### 3.1 Domain Layer (FR-01)

**Priority**: Critical
**Description**: Pure business logic with zero external dependencies

**Requirements**:
- FR-01.1: Value objects must be immutable
- FR-01.2: Validation must occur in value object creation
- FR-01.3: Domain must have zero infrastructure dependencies
- FR-01.4: Business rules must be pure functions
- FR-01.5: Result[T] monads must handle all errors

**Test Coverage**: 42 unit test assertions (2 test functions)

### 3.2 Application Layer (FR-02)

**Priority**: Critical
**Description**: Use case orchestration and port definitions

**Requirements**:
- FR-02.1: Define inbound ports (use case interfaces)
- FR-02.2: Define outbound ports (infrastructure interfaces)
- FR-02.3: Implement use cases using Domain logic
- FR-02.4: Commands must be immutable DTOs
- FR-02.5: Models must be immutable output DTOs
- FR-02.6: Re-export Domain.Error for Presentation access
- FR-02.7: Use generics for static dispatch

**Test Coverage**: Covered by integration tests

### 3.3 Infrastructure Layer (FR-03)

**Priority**: High
**Description**: Concrete adapter implementations

**Requirements**:
- FR-03.1: Implement outbound port interfaces
- FR-03.2: Adapt external systems to Domain types
- FR-03.3: Handle panics at boundaries via defer/recover
- FR-03.4: Convert panics to Result errors
- FR-03.5: Provide console writer adapter
- FR-03.6: Support context.Context for cancellation

**Test Coverage**: Covered by CLI integration tests

### 3.4 Presentation Layer (FR-04)

**Priority**: High
**Description**: User interface adapters (CLI)

**Requirements**:
- FR-04.1: Cannot access Domain layer directly
- FR-04.2: Must use application/error for error handling
- FR-04.3: Must use application/model for output
- FR-04.4: Command line argument parsing
- FR-04.5: User-friendly error messages
- FR-04.6: Exit code mapping (0=success, 1=error)
- FR-04.7: Use generics for static dispatch

**Test Coverage**: 23 CLI integration tests

### 3.5 Bootstrap Layer (FR-05)

**Priority**: High
**Description**: Composition root with dependency wiring

**Requirements**:
- FR-05.1: Wire all generic instantiations
- FR-05.2: Connect ports to adapters
- FR-05.3: Minimal main() (delegate to Bootstrap)
- FR-05.4: Single go.work workspace structure
- FR-05.5: Static wiring (compile-time resolution)

**Test Coverage**: Covered by CLI integration tests

### 3.6 Error Handling (FR-06)

**Priority**: Critical
**Description**: Railway-oriented programming with Result[T] monad

**Requirements**:
- FR-06.1: No panics across layer boundaries
- FR-06.2: Result[T] monad for all fallible operations
- FR-06.3: Error types with kind and message
- FR-06.4: ValidationError for business rule violations
- FR-06.5: InfrastructureError for system failures
- FR-06.6: IsOk()/IsError() predicates
- FR-06.7: Value()/ErrorInfo() accessors
- FR-06.8: UnwrapOr() for default values

**Test Coverage**: All tests verify error handling

### 3.7 Dependency Injection (FR-07)

**Priority**: Critical
**Description**: Static DI via Go generics

**Requirements**:
- FR-07.1: Generic structs for use cases and commands
- FR-07.2: Interface constraints for port abstraction
- FR-07.3: Compile-time instantiation in Bootstrap
- FR-07.4: Zero runtime overhead (no interface dispatch)
- FR-07.5: Type-safe wiring
- FR-07.6: Concrete type parameters at instantiation

**Test Coverage**: Verified by compilation success

### 3.8 Architecture Validation (FR-08)

**Priority**: High
**Description**: Automated boundary enforcement

**Requirements**:
- FR-08.1: Validate Presentation cannot access Domain
- FR-08.2: Validate Infrastructure can access Domain
- FR-08.3: Validate Application accesses Domain only
- FR-08.4: Validate Domain has zero dependencies
- FR-08.5: Python script for validation (arch_guard.py)
- FR-08.6: Make target integration (make check-arch)

**Test Coverage**: Python unit tests for arch_guard.py

### 3.9 Build System (FR-09)

**Priority**: High
**Description**: Comprehensive build automation

**Requirements**:
- FR-09.1: Development build target (go build)
- FR-09.2: Test execution targets (go test)
- FR-09.3: Integration test target (go test -tags=integration)
- FR-09.4: Architecture validation target
- FR-09.5: Clean targets
- FR-09.6: Statistics target

**Test Coverage**: Manual verification of all targets

### 3.10 Test Framework (FR-10)

**Priority**: High
**Description**: Go testing infrastructure

**Requirements**:
- FR-10.1: Go standard testing package
- FR-10.2: testify for assertions
- FR-10.3: Domain unit tests (pure function testing)
- FR-10.4: CLI integration tests (binary execution)
- FR-10.5: Table-driven tests for comprehensive coverage
- FR-10.6: Build tag separation (//go:build integration)

**Test Coverage**: Self-verifying (tests use the framework)

### 3.11 Documentation (FR-11)

**Priority**: High
**Description**: Complete project documentation

**Requirements**:
- FR-11.1: Software Requirements Specification (this document)
- FR-11.2: Software Design Specification
- FR-11.3: Software Test Guide
- FR-11.4: Quick Start Guide
- FR-11.5: UML diagrams (PlantUML sources + SVG)
- FR-11.6: Inline code documentation (godoc)
- FR-11.7: README with examples

**Test Coverage**: Documentation review process

### 3.12 Code Quality (FR-12)

**Priority**: High
**Description**: Professional code standards

**Requirements**:
- FR-12.1: Zero go vet warnings
- FR-12.2: Zero staticcheck violations
- FR-12.3: Consistent naming conventions (Go style)
- FR-12.4: File headers with copyright and SPDX
- FR-12.5: Comprehensive godoc documentation
- FR-12.6: Context support for cancellation

**Test Coverage**: Build verification

---

## 4. Non-Functional Requirements

### 4.1 Performance (NFR-01)

**Priority**: Medium

- NFR-01.1: Static dispatch overhead: 0 (compile-time resolution)
- NFR-01.2: Result monad overhead: minimal (value types)
- NFR-01.3: Build time: < 5 seconds (clean build)
- NFR-01.4: Test execution: < 3 seconds (all tests)

**Verification**: Benchmarks, profiling

### 4.2 Reliability (NFR-02)

**Priority**: High

- NFR-02.1: All tests must pass (100% pass rate)
- NFR-02.2: No memory leaks
- NFR-02.3: Deterministic error handling (no panics)
- NFR-02.4: Type-safe boundaries (compile-time verification)
- NFR-02.5: Context cancellation support

**Verification**: Test suite, race detector

### 4.3 Portability (NFR-03)

**Priority**: High

- NFR-03.1: Support POSIX platforms (Linux, macOS, BSD)
- NFR-03.2: Support Windows
- NFR-03.3: Standard Go 1.22 (no CGO dependencies)
- NFR-03.4: go.work compatible project structure
- NFR-03.5: No platform-specific code in Domain/Application

**Verification**: Multi-platform CI testing

### 4.4 Maintainability (NFR-04)

**Priority**: Critical

- NFR-04.1: Clear layer separation (enforced by arch_guard.py)
- NFR-04.2: Self-documenting code with godoc
- NFR-04.3: Comprehensive test coverage
- NFR-04.4: Standard file naming conventions
- NFR-04.5: Consistent code style (gofmt)
- NFR-04.6: Version control friendly

**Verification**: Architecture validation, code review

### 4.5 Usability (NFR-05)

**Priority**: High

- NFR-05.1: Quick Start Guide for beginners
- NFR-05.2: Working examples in under 5 minutes
- NFR-05.3: Clear error messages
- NFR-05.4: Comprehensive documentation
- NFR-05.5: Educational UML diagrams
- NFR-05.6: Make target help system

**Verification**: User documentation review

### 4.6 Testability (NFR-06)

**Priority**: Critical

- NFR-06.1: Pure functions in Domain (easy to test)
- NFR-06.2: Port abstraction for test doubles
- NFR-06.3: Standard Go testing framework
- NFR-06.4: Test organization by type (unit/integration)
- NFR-06.5: Coverage analysis support (go test -cover)

**Verification**: Test suite execution

---

## 5. System Constraints

### 5.1 Technical Constraints

- **SC-01**: Must compile with Go 1.22+
- **SC-02**: Must use Go generics for static dispatch
- **SC-03**: Must be go.work compatible
- **SC-04**: Uses testify for test assertions
- **SC-05**: No external runtime dependencies in production

### 5.2 Design Constraints

- **SC-06**: Must enforce hexagonal architecture boundaries
- **SC-07**: Presentation cannot access Domain directly
- **SC-08**: Domain must have zero external dependencies
- **SC-09**: Must use static dispatch (generics, not interfaces for DI)
- **SC-10**: No panics across layer boundaries
- **SC-11**: All errors via Result[T] monad
- **SC-12**: Context support for cancellation

### 5.3 Regulatory Constraints

- **SC-13**: BSD-3-Clause license
- **SC-14**: SPDX identifiers in all source files
- **SC-15**: Copyright attribution to Michael Gardner, A Bit of Help, Inc.

---

## 6. Verification and Validation

### 6.1 Test Coverage Matrix

| Requirement | Test Type | Test Count | Status |
|-------------|-----------|------------|--------|
| FR-01 (Domain) | Unit | 42 assertions | Pass |
| FR-02 (Application) | Integration | Covered | Pass |
| FR-03 (Infrastructure) | Integration | Covered | Pass |
| FR-04 (Presentation) | Integration | 23 tests | Pass |
| FR-05 (Bootstrap) | Integration | Covered | Pass |
| FR-06 (Error Handling) | All | Verified | Pass |
| FR-07 (Static DI) | Compile-time | N/A | Verified |
| FR-08 (Arch Validation) | Python Unit | Tests | Pass |
| FR-09 (Build System) | Manual | All targets | Verified |
| FR-10 (Test Framework) | Self-test | N/A | Pass |
| FR-11 (Documentation) | Review | Complete | Verified |
| FR-12 (Code Quality) | Build | 0 warnings | Verified |

### 6.2 Verification Methods

- **Code Review**: All code reviewed before release
- **Static Analysis**: go vet, staticcheck
- **Dynamic Testing**: Domain unit tests + CLI integration tests
- **Architecture Validation**: arch_guard.py enforcement
- **Coverage Analysis**: go test -cover
- **Documentation Review**: Complete formal specifications

---

## 7. Appendices

### 7.1 Project Statistics

**Source Code**:
- Go source files (.go): ~25
- Total lines of code: ~2,500

**Tests**:
- Domain unit tests: 2 functions, 42 assertions
- CLI integration tests: 23 tests
- Pass rate: 100%

**Documentation**:
- Formal specs: 3 (SRS, SDS, Test Guide)
- Guides: 2 (Quick Start, Architecture Mapping)
- UML diagrams: 5
- README: Complete with examples

**Build System**:
- Makefile targets: 20+
- Dependencies: testify (test only)

### 7.2 Layer Responsibilities Summary

| Layer | Responsibilities | Dependencies | Tests |
|-------|------------------|--------------|-------|
| Domain | Business logic, validation | NONE | Unit (42) |
| Application | Use cases, ports | Domain | Integration |
| Infrastructure | Adapters (driven) | App + Domain | Integration |
| Presentation | UI (driving) | Application ONLY | Integration (23) |
| Bootstrap | DI wiring | ALL | Integration |

### 7.3 Static Dispatch Pattern

```go
// Bootstrap instantiates with concrete types
consoleWriter := adapter.NewConsoleWriter()
greetUseCase := usecase.NewGreetUseCase[*adapter.ConsoleWriter](consoleWriter)
greetCommand := command.NewGreetCommand[*usecase.GreetUseCase[*adapter.ConsoleWriter]](greetUseCase)

// All method calls are statically dispatched (no vtable)
exitCode := greetCommand.Run(args)
```

### 7.4 Dependency Graph

```
Bootstrap
    |
Presentation -> Application -> Domain
    |              |
Infrastructure ----+
```

**Critical Rule**: Presentation cannot access Domain directly (enforced by arch_guard.py)

---

## 8. Traceability Matrix

| FR ID | Design Element | Test Coverage | Status |
|-------|---------------|---------------|--------|
| FR-01 | domain/valueobject/person.go | 42 unit tests | Pass |
| FR-02 | application/usecase/greet.go | Integration | Pass |
| FR-03 | infrastructure/adapter/consolewriter.go | Integration | Pass |
| FR-04 | presentation/cli/command/greet.go | 23 integration | Pass |
| FR-05 | bootstrap/cli/cli.go | Integration | Pass |
| FR-06 | domain/error/result.go | All tests | Pass |
| FR-07 | Generic instantiation | Compile-time | Verified |
| FR-08 | arch_guard.py | Python tests | Pass |
| FR-09 | Makefile | Manual | Verified |
| FR-10 | go test + testify | Self-test | Pass |
| FR-11 | docs/ directory | Review | Verified |
| FR-12 | Build verification | 0 warnings | Verified |

---

**Document Control**:
- Version: 1.0.0
- Last Updated: November 25, 2025
- Status: Released
- Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
- License: BSD-3-Clause
- SPDX-License-Identifier: BSD-3-Clause
