# Software Requirements Specification (SRS)

**Project:** Hybrid_App_Ada - Ada 2022 Application Starter
**Version:** 1.0.0
**Date:** November 18, 2025
**SPDX-License-Identifier:** BSD-3-Clause
**License File:** See the LICENSE file in the project root.
**Copyright:** © 2025 Michael Gardner, A Bit of Help, Inc.
**Status:** Released

---

## 1. Introduction

### 1.1 Purpose

This Software Requirements Specification (SRS) describes the functional and non-functional requirements for Hybrid_App_Ada, a professional Ada 2022 application starter template demonstrating hexagonal architecture with functional programming principles.

### 1.2 Scope

Hybrid_App_Ada provides:
- Professional 5-layer hexagonal architecture implementation
- Static dependency injection via Ada generics
- Railway-oriented programming with Result monads
- Clean architecture boundary enforcement
- Comprehensive test suite with custom framework
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

### 1.4 References

- Ada 2022 Language Reference Manual
- Clean Architecture (Robert C. Martin)
- Domain-Driven Design (Eric Evans)
- Railway-Oriented Programming (Scott Wlaschin)
- Hexagonal Architecture (Alistair Cockburn)

---

## 2. Overall Description

### 2.1 Product Perspective

Hybrid_App_Ada is a standalone application starter template implementing professional architectural patterns:

**Architecture Layers**:
1. **Domain**: Pure business logic (zero dependencies)
2. **Application**: Use cases and port definitions
3. **Infrastructure**: Driven adapters (implementations)
4. **Presentation**: Driving adapters (user interfaces)
5. **Bootstrap**: Composition root (dependency wiring)

### 2.2 Product Features

1. **Hexagonal Architecture**: 5-layer clean architecture
2. **Static Dependency Injection**: Compile-time wiring via generics
3. **Railway-Oriented Programming**: Result monad error handling
4. **Architecture Enforcement**: Automated boundary validation
5. **Test Infrastructure**: Custom test framework, 82 tests
6. **Build Automation**: Comprehensive Makefile
7. **Documentation**: Complete SRS, SDS, Test Guide

### 2.3 User Classes

- **Application Developers**: Learn hexagonal architecture patterns
- **Team Leads**: Adopt architectural standards
- **Educators**: Teach clean architecture principles
- **Ada Developers**: Start new projects with best practices

### 2.4 Operating Environment

- **Platforms**: Linux, macOS, BSD, Windows
- **Compiler**: GNAT FSF 13+ or GNAT Pro (Ada 2022 support)
- **Build System**: Alire 2.0+, GNU Make
- **Ada Version**: Ada 2022

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
- FR-01.5: Result monads must handle all errors

**Test Coverage**: 48 unit tests

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

**Test Coverage**: 26 integration tests

### 3.3 Infrastructure Layer (FR-03)

**Priority**: High
**Description**: Concrete adapter implementations

**Requirements**:
- FR-03.1: Implement outbound port interfaces
- FR-03.2: Adapt external systems to Domain types
- FR-03.3: Handle infrastructure exceptions at boundaries
- FR-03.4: Convert exceptions to Result errors
- FR-03.5: Provide console writer adapter

**Test Coverage**: Covered by integration tests

### 3.4 Presentation Layer (FR-04)

**Priority**: High
**Description**: User interface adapters (CLI)

**Requirements**:
- FR-04.1: Cannot access Domain layer directly
- FR-04.2: Must use Application.Error for error handling
- FR-04.3: Must use Application.Model for output
- FR-04.4: Command line argument parsing
- FR-04.5: User-friendly error messages
- FR-04.6: Exit code mapping (0=success, 1=error)

**Test Coverage**: 8 E2E tests

### 3.5 Bootstrap Layer (FR-05)

**Priority**: High
**Description**: Composition root with dependency wiring

**Requirements**:
- FR-05.1: Wire all generic instantiations
- FR-05.2: Connect ports to adapters
- FR-05.3: Minimal main procedure (delegate to Bootstrap)
- FR-05.4: Single-project structure (no aggregates)
- FR-05.5: Static wiring (compile-time resolution)

**Test Coverage**: Covered by E2E tests

### 3.6 Error Handling (FR-06)

**Priority**: Critical
**Description**: Railway-oriented programming with Result monad

**Requirements**:
- FR-06.1: No exceptions across layer boundaries
- FR-06.2: Result monad for all fallible operations
- FR-06.3: Error types with kind and message
- FR-06.4: Validation_Error for business rule violations
- FR-06.5: Infrastructure_Error for system failures
- FR-06.6: Is_Ok/Is_Error predicates
- FR-06.7: Value/Error_Info accessors

**Test Coverage**: All tests verify error handling

### 3.7 Dependency Injection (FR-07)

**Priority**: Critical
**Description**: Static DI via Ada generics

**Requirements**:
- FR-07.1: Generic packages for use cases
- FR-07.2: Generic function parameters for ports
- FR-07.3: Compile-time instantiation in Bootstrap
- FR-07.4: Zero runtime overhead
- FR-07.5: Type-safe wiring
- FR-07.6: No heap allocation for DI

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
- FR-09.1: Development build target
- FR-09.2: Optimized build target (O2)
- FR-09.3: Release build target
- FR-09.4: Test execution targets (unit/integration/e2e)
- FR-09.5: Coverage analysis target
- FR-09.6: Architecture validation target
- FR-09.7: Clean targets (normal/deep/coverage)
- FR-09.8: Statistics target

**Test Coverage**: Manual verification of all targets

### 3.10 Test Framework (FR-10)

**Priority**: High
**Description**: Custom lightweight testing infrastructure

**Requirements**:
- FR-10.1: No AUnit dependency
- FR-10.2: Simple assertion API
- FR-10.3: Test count and pass tracking
- FR-10.4: Result reporting
- FR-10.5: Exit code support (0=pass, 1=fail)
- FR-10.6: Test organization (unit/integration/e2e)
- FR-10.7: Test runners for each level

**Test Coverage**: Self-verifying (tests use the framework)

### 3.11 Documentation (FR-11)

**Priority**: High
**Description**: Complete project documentation

**Requirements**:
- FR-11.1: Software Requirements Specification (this document)
- FR-11.2: Software Design Specification
- FR-11.3: Software Test Guide
- FR-11.4: Quick Start Guide
- FR-11.5: Development Roadmap
- FR-11.6: UML diagrams (PlantUML sources + SVG)
- FR-11.7: Inline code documentation
- FR-11.8: README with examples

**Test Coverage**: Documentation review process

### 3.12 Code Quality (FR-12)

**Priority**: High
**Description**: Professional code standards

**Requirements**:
- FR-12.1: Zero compiler warnings
- FR-12.2: Zero style violations
- FR-12.3: Ada 2022 aspect syntax (no obsolescent pragmas)
- FR-12.4: Consistent naming conventions
- FR-12.5: File headers with copyright and SPDX
- FR-12.6: Comprehensive docstrings

**Test Coverage**: Build verification

---

## 4. Non-Functional Requirements

### 4.1 Performance (NFR-01)

**Priority**: Medium

- NFR-01.1: Static dispatch overhead: 0 (compile-time resolution)
- NFR-01.2: No heap allocation in hot paths
- NFR-01.3: Result monad overhead: minimal (stack-based)
- NFR-01.4: Build time: < 30 seconds (clean build)
- NFR-01.5: Test execution: < 5 seconds (all tests)

**Verification**: Benchmarks, profiling

### 4.2 Reliability (NFR-02)

**Priority**: High

- NFR-02.1: All 82 tests must pass (100% pass rate)
- NFR-02.2: No memory leaks
- NFR-02.3: Deterministic error handling (no exceptions)
- NFR-02.4: Type-safe boundaries (compile-time verification)

**Verification**: Test suite, static analysis

### 4.3 Portability (NFR-03)

**Priority**: High

- NFR-03.1: Support POSIX platforms (Linux, macOS, BSD)
- NFR-03.2: Support Windows
- NFR-03.3: Standard Ada 2022 (no compiler-specific features)
- NFR-03.4: Alire-compatible project structure
- NFR-03.5: No platform-specific code in Domain/Application

**Verification**: Multi-platform CI testing

### 4.4 Maintainability (NFR-04)

**Priority**: Critical

- NFR-04.1: Clear layer separation (enforced by arch_guard.py)
- NFR-04.2: Self-documenting code with docstrings
- NFR-04.3: Comprehensive test coverage (82 tests)
- NFR-04.4: Standard file naming conventions
- NFR-04.5: Consistent code style
- NFR-04.6: Version control friendly (single project)

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
- NFR-06.3: Custom test framework (no external dependencies)
- NFR-06.4: Test organization by type (unit/integration/e2e)
- NFR-06.5: Coverage analysis support (GNATcoverage)

**Verification**: Test suite execution

---

## 5. System Constraints

### 5.1 Technical Constraints

- **SC-01**: Must compile with GNAT FSF 13+ or GNAT Pro
- **SC-02**: Must use Ada 2022 language features
- **SC-03**: Must be Alire-compatible
- **SC-04**: Single-project structure (no aggregate projects)
- **SC-05**: Must use functional crate v1.0.0 for Result monad
- **SC-06**: No AUnit dependency (custom test framework)

### 5.2 Design Constraints

- **SC-07**: Must enforce hexagonal architecture boundaries
- **SC-08**: Presentation cannot access Domain directly
- **SC-09**: Domain must have zero external dependencies
- **SC-10**: Must use static dispatch (generics, not interfaces)
- **SC-11**: No exceptions across layer boundaries
- **SC-12**: All errors via Result monad

### 5.3 Regulatory Constraints

- **SC-13**: BSD-3-Clause license
- **SC-14**: SPDX identifiers in all source files
- **SC-15**: Copyright attribution to Michael Gardner, A Bit of Help, Inc.

---

## 6. Verification and Validation

### 6.1 Test Coverage Matrix

| Requirement | Test Type | Test Count | Status |
|-------------|-----------|------------|--------|
| FR-01 (Domain) | Unit | 48 | ✅ Pass |
| FR-02 (Application) | Integration | 26 | ✅ Pass |
| FR-03 (Infrastructure) | Integration | Covered | ✅ Pass |
| FR-04 (Presentation) | E2E | 8 | ✅ Pass |
| FR-05 (Bootstrap) | E2E | Covered | ✅ Pass |
| FR-06 (Error Handling) | All | 82 | ✅ Pass |
| FR-07 (DI) | Compile-time | N/A | ✅ Verified |
| FR-08 (Arch Validation) | Python Unit | 5 | ✅ Pass |
| FR-09 (Build System) | Manual | All targets | ✅ Verified |
| FR-10 (Test Framework) | Self-test | N/A | ✅ Pass |
| FR-11 (Documentation) | Review | Complete | ✅ Verified |
| FR-12 (Code Quality) | Build | 0 warnings | ✅ Verified |

### 6.2 Verification Methods

- **Code Review**: All code reviewed before release
- **Static Analysis**: Zero compiler warnings enforced
- **Dynamic Testing**: 82 automated tests (100% pass rate)
- **Architecture Validation**: arch_guard.py enforcement
- **Coverage Analysis**: GNATcoverage support
- **Documentation Review**: Complete formal specifications

---

## 7. Appendices

### 7.1 Project Statistics

**Source Code**:
- Ada specification files (.ads): 26
- Ada implementation files (.adb): 11
- Total lines of code: ~3,500

**Tests**:
- Total test files: 11
- Total tests: 82
  - Unit tests: 48
  - Integration tests: 26
  - E2E tests: 8
- Pass rate: 100%

**Documentation**:
- Formal specs: 3 (SRS, SDS, Test Guide)
- Guides: 2 (Quick Start, Roadmap)
- UML diagrams: 5
- README: Complete with examples

**Build System**:
- Makefile targets: 30+
- Dependencies: 3 (functional, aunit, gnatcov)

### 7.2 Layer Responsibilities Summary

| Layer | Responsibilities | Dependencies | Tests |
|-------|------------------|--------------|-------|
| Domain | Business logic, validation | NONE | Unit (48) |
| Application | Use cases, ports | Domain | Integration (26) |
| Infrastructure | Adapters (driven) | App + Domain | Integration |
| Presentation | UI (driving) | Application ONLY | E2E (8) |
| Bootstrap | DI wiring | ALL | E2E |

### 7.3 Error Handling Flow

```
Domain Validation Error
    ↓
Result[Value, Error] returned
    ↓
Application checks Is_Error
    ↓
If error: propagate up via Result
    ↓
Presentation pattern matches Error_Kind
    ↓
User-friendly message displayed
    ↓
Exit code 1
```

### 7.4 Dependency Graph

```
Bootstrap
    ↓
Presentation → Application → Domain
    ↓              ↓
Infrastructure ────┘
```

**Critical Rule**: Presentation cannot access Domain directly (enforced by arch_guard.py)

---

## 8. Traceability Matrix

| FR ID | Design Element | Test Coverage | Status |
|-------|---------------|---------------|--------|
| FR-01 | Domain.Value_Object.Person | 48 unit tests | ✅ |
| FR-02 | Application.Usecase.Greet | 26 integration | ✅ |
| FR-03 | Infrastructure.Adapter.Console_Writer | Integration | ✅ |
| FR-04 | Presentation.CLI.Command.Greet | 8 E2E tests | ✅ |
| FR-05 | Bootstrap.CLI | E2E tests | ✅ |
| FR-06 | Domain.Error.Result | All tests | ✅ |
| FR-07 | Generic instantiation | Compile-time | ✅ |
| FR-08 | arch_guard.py | Python tests | ✅ |
| FR-09 | Makefile | Manual | ✅ |
| FR-10 | test/common/test_framework | Self-test | ✅ |
| FR-11 | docs/ directory | Review | ✅ |
| FR-12 | Build verification | 0 warnings | ✅ |

---

**Document Control**:
- Version: 1.0.0
- Last Updated: November 18, 2025
- Status: Released
- Copyright © 2025 Michael Gardner, A Bit of Help, Inc.
- License: BSD-3-Clause
- SPDX-License-Identifier: BSD-3-Clause
