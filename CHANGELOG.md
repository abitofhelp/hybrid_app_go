# Changelog

**Project:** Hybrid_App_Ada - Ada 2022 Application Starter
**SPDX-License-Identifier:** BSD-3-Clause
**License File:** See the LICENSE file in the project root.
**Copyright:** © 2025 Michael Gardner, A Bit of Help, Inc.

All notable changes to Hybrid_App_Ada will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Fixed
- **Code Quality**: Resolved all compiler warnings and style violations (357 automated fixes)
  - Fixed array aggregate syntax to Ada 2022 standard (`[]` for aggregates)
  - Fixed array constraint syntax (kept `()` for type constraints)
  - Removed unused with/use clauses across test suite
  - Fixed long lines exceeding 80 characters (proper line breaks at logical points)
  - Updated short-circuit operators (`and` → `and then`, `or` → `or else`)
  - Standardized comment separator lines to 79 characters
- **Code Safety**: Eliminated unsafe code in Infrastructure.Adapter.Console_Writer
  - Removed module-level mutable state and `Unchecked_Access` usage
  - Migrated to safe parameterized Functional.Try pattern
  - Enhanced Functional library with `Try_To_Result_With_Param` and `Try_To_Option_With_Param`
  - Added backwards-compatible child packages for existing tests
  - All 83 Functional library tests passing, all 105 hybrid_app_ada tests passing
- **Architecture Validation**: Fixed path detection in `arch_guard.py` after script reorganization
- **Build System**: Corrected Makefile FORMAT_DIRS syntax (missing backslash continuation)

### Changed
- **Script Organization**: Reorganized scripts into purpose-driven subdirectories
  - Created `scripts/makefile/` for Makefile-related scripts
  - Created `scripts/release/` for release automation
  - Updated all Makefile paths accordingly
- **GNATcoverage Runtime**: Automated runtime build from sources
  - Created `build_gnatcov_runtime.py` for reproducible builds
  - Removed vendored runtime with hardcoded paths
  - Added `make build-coverage-runtime` target
  - Updated coverage workflow documentation

### Removed
- Deleted tzif-specific scripts not applicable to hybrid_app_ada:
  - `check_tzif_parser.py`
  - `generate_version_ada.py`
  - `compare_tzif_versions.py`
  - `test_parser.py`
  - `generate_tzdata_version.py`

### Added
- 5-layer hexagonal architecture (Domain, Application, Infrastructure, Presentation, Bootstrap)
- Static dispatch dependency injection via generics (zero runtime overhead)
- Railway-oriented programming with Result monads (functional error handling)
- Application.Error re-export pattern (Presentation isolation from Domain)
- Single-project structure (simplified from aggregate project)
- Comprehensive UML diagrams (5 PlantUML diagrams with SVG generation)
- Architecture validation script (arch_guard.py enforces layer boundaries)
- **Comprehensive test suite with 105 tests achieving 90%+ code coverage**:
  - 81 unit tests (Domain + Application layers, 100% coverage)
  - 16 integration tests (cross-layer flows with real components)
  - 8 E2E tests (black-box CLI binary testing)
- **Professional test framework** (test/common/test_framework):
  - Color-coded test output (bright green success, bright red failure)
  - `Print_Category_Summary` function for consistent, professional test reporting
  - Exit code integration for CI/CD pipelines
  - Reusable test helpers eliminating code duplication
- **Test organization by architectural layer**:
  - test/unit/ - Unit tests with unit_tests.gpr
  - test/integration/ - Integration tests with integration_tests.gpr
  - test/e2e/ - E2E tests with e2e_tests.gpr
- Professional documentation (SDS, SRS, Test Guide, Quick Start)
- Makefile with comprehensive build/test/coverage targets and color-coded test summaries
- Code coverage support with vendored GNATcoverage RTS
- PlantUML diagram generation tooling
- Example greeter application demonstrating all patterns

### Changed
- Renamed test directory from `tests/` to `test/` (Ada standard convention)
- Enhanced Makefile test-all target with professional bordered success/failure indicators

### Architecture Patterns
- **Static Dependency Injection**: Generic packages with function parameters (compile-time DI)
- **Result Monad**: Railway-oriented error handling (no exceptions across boundaries)
- **Presentation Isolation**: Only Domain layer is shareable across apps
- **Minimal Entry Point**: Main delegates to Bootstrap.CLI.Run (3 lines)
- **Ports & Adapters**: Application defines ports, Infrastructure implements adapters

### Technical Details
- Ada 2022 with aspects (not obsolescent pragmas)
- Bounded strings (no heap allocation in domain)
- Pre/Post contracts on all public operations
- Pure domain logic with zero dependencies
- Functional composition via generics

### Documentation
- README with architecture overview and static dispatch explanation
- 5 UML diagrams:
  - Layer dependencies with architectural rules
  - Application.Error re-export pattern
  - Package structure by layer
  - Error handling flow (railway-oriented)
  - Static vs dynamic dispatch comparison
- Software Design Specification
- Software Requirements Specification
- Software Test Guide
- Quick Start Guide

---

## [1.0.0] - TBD

_First stable release - Professional Ada 2022 application starter template._

Coming soon: First stable release demonstrating hybrid DDD/Clean/Hexagonal architecture with functional programming principles in Ada 2022.

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
- **Copyright**: © 2025 Michael Gardner, A Bit of Help, Inc.
- **SPDX-License-Identifier**: BSD-3-Clause
