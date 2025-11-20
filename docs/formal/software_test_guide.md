# Software Test Guide

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

This Software Test Guide describes the testing approach, test organization, execution procedures, and guidelines for the Hybrid_App_Ada project.

### 1.2 Scope

This document covers:
- Test strategy and organization (unit/integration/e2e)
- Custom test framework design
- Running tests via Make and direct execution
- Writing new tests
- Coverage analysis with GNATcoverage
- Test maintenance procedures

---

## 2. Test Strategy

### 2.1 Testing Levels

Hybrid_App_Ada uses three levels of testing:

**Unit Tests** (48 tests)
- Test individual components in isolation
- Focus on Domain and Application logic
- Pure functions, predictable results
- Fast execution
- Location: `test/unit/`

**Integration Tests** (26 tests)
- Test cross-layer interactions
- Real Infrastructure adapters
- Application use cases with dependencies
- Verify wiring works correctly
- Location: `test/integration/`

**End-to-End Tests** (8 tests)
- Test entire system via CLI
- Black-box testing
- User scenarios
- Exit code verification
- Location: `test/e2e/`

### 2.2 Testing Philosophy

- **Test-Driven**: Tests written alongside or before code
- **Railway-Oriented**: Test both success and error paths
- **Comprehensive**: Cover normal, edge, and error cases
- **Automated**: All tests runnable via `make test-all`
- **Fast**: All 82 tests execute in < 5 seconds
- **Self-Contained**: No external dependencies (custom framework)

---

## 3. Test Organization

### 3.1 Directory Structure

```
test/
├── common/
│   └── test_framework.ads/adb     # Custom test framework
├── unit/
│   ├── test_domain_error.adb      # Domain error tests
│   ├── test_person.adb             # Person value object tests
│   ├── test_greet_command.adb      # Application command tests
│   ├── test_greet_usecase.adb      # Use case tests
│   ├── unit_runner.adb             # Unit test runner
│   └── unit_tests.gpr              # Unit tests project
├── integration/
│   ├── test_greet_integration.adb  # Full use case integration
│   ├── test_console_writer.adb     # Infrastructure adapter tests
│   ├── integration_runner.adb      # Integration test runner
│   └── integration_tests.gpr       # Integration tests project
└── e2e/
    ├── test_cli_success.adb        # CLI happy path tests
    ├── test_cli_errors.adb         # CLI error scenarios
    ├── e2e_runner.adb              # E2E test runner
    └── e2e_tests.gpr               # E2E tests project
```

### 3.2 Test Naming Convention

- **Pattern**: `test_<component>.adb`
- **Examples**:
  - `test_person.adb` → Tests `Domain.Value_Object.Person`
  - `test_greet_usecase.adb` → Tests `Application.Usecase.Greet`
  - `test_console_writer.adb` → Tests `Infrastructure.Adapter.Console_Writer`

---

## 4. Custom Test Framework

### 4.1 Design Rationale

Hybrid_App_Ada uses a **custom lightweight test framework** instead of AUnit:

**Benefits**:
- No external test framework dependency
- Simple, easy to understand (< 100 lines)
- Full control over assertion messages
- Educational value (see how test frameworks work)
- Fast compilation and execution

### 4.2 Framework API

Located in `test/common/test_framework.{ads,adb}`:

```ada
-- Initialize test tracking
procedure Start_Test_Suite (Name : String);

-- Assert with automatic counting
procedure Assert (
   Condition : Boolean;
   Test_Name : String);

-- Report results
procedure Report_Results (Exit_On_Failure : Boolean := True);
```

### 4.3 Example Test

```ada
pragma Ada_2022;
with Test_Framework; use Test_Framework;
with Domain.Value_Object.Person;

procedure Test_Person is
   use Domain.Value_Object.Person;
begin
   Start_Test_Suite ("Person Value Object Tests");

   -- Test successful creation
   declare
      Result : constant Person_Result.Result := Create ("Alice");
   begin
      Assert (Person_Result.Is_Ok (Result),
              "Create with valid name returns Ok");
      Assert (Get_Name (Person_Result.Value (Result)) = "Alice",
              "Created person has correct name");
   end;

   -- Test validation error
   declare
      Result : constant Person_Result.Result := Create ("");
   begin
      Assert (Person_Result.Is_Error (Result),
              "Create with empty name returns Error");
   end;

   Report_Results;
end Test_Person;
```

---

## 5. Running Tests

### 5.1 Quick Start

```bash
# Run all 82 tests
make test-all

# Run specific test level
make test-unit          # 48 unit tests
make test-integration   # 26 integration tests
make test-e2e           # 8 e2e tests
```

### 5.2 Make Targets

**Test Execution**:
```bash
make test               # Alias for test-all
make test-all           # Run all test executables
make test-unit          # Unit tests only
make test-integration   # Integration tests only
make test-e2e           # E2E tests only
```

**Build and Test**:
```bash
make build build-tests  # Build main + tests
make test-all           # Run tests
```

**Coverage**:
```bash
make test-coverage      # Run with GNATcoverage analysis
```

### 5.3 Direct Execution

```bash
# Build tests first
make build-tests

# Run test runners directly
./test/bin/unit_runner
./test/bin/integration_runner
./test/bin/e2e_runner

# Run individual test executable
./test/bin/test_person
```

### 5.4 Expected Output

**Success**:
```
========================================
Test Suite: Person Value Object Tests
========================================

  [PASS] Create with valid name returns Ok
  [PASS] Created person has correct name
  [PASS] Create with empty name returns Error
  [PASS] Error has correct kind
  [PASS] Error has descriptive message

========================================
  Results: 5 / 5 passed
  Status: ALL TESTS PASSED
========================================
```

**Failure**:
```
========================================
Test Suite: Person Value Object Tests
========================================

  [PASS] Create with valid name returns Ok
  [FAIL] Created person has correct name
  [PASS] Create with empty name returns Error

========================================
  Results: 2 / 3 passed
  Status: TESTS FAILED
========================================
```

---

## 6. Writing New Tests

### 6.1 Unit Test Template

```ada
pragma Ada_2022;
with Test_Framework; use Test_Framework;
with Your.Component;

procedure Test_Your_Component is
   use Your.Component;
begin
   Start_Test_Suite ("Your Component Tests");

   -- Test case 1: Success scenario
   declare
      Result : constant Result_Type := Some_Function (Valid_Input);
   begin
      Assert (Is_Ok (Result), "Valid input produces Ok result");
   end;

   -- Test case 2: Error scenario
   declare
      Result : constant Result_Type := Some_Function (Invalid_Input);
   begin
      Assert (Is_Error (Result), "Invalid input produces Error result");
   end;

   Report_Results;
end Test_Your_Component;
```

### 6.2 Integration Test Template

```ada
pragma Ada_2022;
with Test_Framework; use Test_Framework;
with Application.Usecase.Your_UseCase;
with Infrastructure.Adapter.Your_Adapter;

procedure Test_Your_UseCase_Integration is
   use Test_Framework;
begin
   Start_Test_Suite ("Your Use Case Integration Tests");

   -- Test with real infrastructure
   declare
      Cmd    : constant Your_Command := Make_Command ("test_data");
      Result : constant Unit_Result.Result := Execute_UseCase (Cmd);
   begin
      Assert (Unit_Result.Is_Ok (Result),
              "Use case succeeds with real adapter");
   end;

   Report_Results;
end Test_Your_UseCase_Integration;
```

### 6.3 E2E Test Template

```ada
pragma Ada_2022;
with Test_Framework; use Test_Framework;
with Ada.Command_Line;
with Bootstrap.CLI;

procedure Test_CLI_Scenario is
begin
   Start_Test_Suite ("CLI End-to-End Tests");

   -- Simulate CLI execution
   declare
      Exit_Code : Integer;
   begin
      -- Set up command line args
      -- (In real tests, this would be more sophisticated)

      Exit_Code := Bootstrap.CLI.Run;

      Assert (Exit_Code = 0, "CLI returns success exit code");
   end;

   Report_Results;
end Test_CLI_Scenario;
```

### 6.4 Adding Tests to Runners

**Unit Test Runner** (`test/unit/unit_runner.adb`):
```ada
with Test_Person;
with Test_Your_New_Component;  -- Add this

procedure Unit_Runner is
begin
   Test_Person;
   Test_Your_New_Component;      -- Add this
end Unit_Runner;
```

Update `test/unit/unit_tests.gpr`:
```gprbuild
project Unit_Tests is
   for Main use (
      "unit_runner.adb",
      "test_person.adb",
      "test_your_new_component.adb"  -- Add this
   );
end Unit_Tests;
```

---

## 7. Test Coverage

### 7.1 Coverage Goals

- **Target**: > 90% line coverage
- **Critical Code**: 100% coverage for error handling paths
- **Domain Layer**: Near 100% coverage (pure functions, easy to test)

### 7.2 Running Coverage Analysis

```bash
# Generate coverage report
make test-coverage

# Coverage data generated in coverage/ directory
```

### 7.3 Coverage Process

The `make test-coverage` target:
1. Builds instrumented binaries
2. Runs all test suites
3. Collects coverage traces
4. Generates HTML report
5. Outputs summary to console

### 7.4 Interpreting Results

Coverage report shows:
- Lines executed vs total
- Branches taken
- Functions covered
- Per-file statistics

---

## 8. Test Maintenance

### 8.1 When to Update Tests

- **New Features**: Add tests before implementing
- **Bug Fixes**: Add regression test first
- **Refactoring**: Ensure tests still pass
- **Requirements Change**: Update affected tests

### 8.2 Test Quality Guidelines

- **One Assertion Per Test Case**: Clear failure messages
- **Descriptive Names**: Test names explain what's being verified
- **Arrange-Act-Assert**: Structure tests clearly
- **No Business Logic**: Tests should be simple
- **Fast Execution**: Avoid slow operations

### 8.3 Debugging Failed Tests

```bash
# Run single test for debugging
./test/bin/test_person

# Add debug output (temporarily)
Ada.Text_IO.Put_Line ("Debug: Variable = " & Variable'Image);

# Use gdb for step-through debugging
gdb ./test/bin/test_person
```

---

## 9. Continuous Integration

### 9.1 CI Testing Strategy

All tests run on every commit:

```bash
# CI pipeline equivalent
make clean
make build
make build-tests
make test-all
make check-arch
```

### 9.2 Success Criteria

All must pass:
- ✅ Zero build warnings
- ✅ All 82 tests pass (100% pass rate)
- ✅ Architecture validation passes
- ✅ Exit code 0 from test runners

---

## 10. Test Statistics

### 10.1 Current Test Metrics (v1.0.0)

**Test Count**:
- Total: 82 tests
  - Unit: 48 (58%)
  - Integration: 26 (32%)
  - E2E: 8 (10%)
- Pass Rate: 100%

**Coverage**:
- Domain Layer: High (pure functions, fully tested)
- Application Layer: High (use cases covered)
- Infrastructure Layer: Adequate (integration tests)
- Presentation Layer: Good (E2E tests)
- Bootstrap Layer: Verified (E2E tests)

**Execution Time**:
- Unit tests: < 1 second
- Integration tests: < 2 seconds
- E2E tests: < 2 seconds
- Total: < 5 seconds (all 82 tests)

---

## 11. Common Testing Patterns

### 11.1 Testing Result Monads

```ada
-- Test success path
declare
   Result : constant Result_Type := Function_Under_Test (Valid_Input);
begin
   Assert (Result_Type.Is_Ok (Result), "Success case returns Ok");
   Assert (Result_Type.Value (Result) = Expected_Value,
           "Result has correct value");
end;

-- Test error path
declare
   Result : constant Result_Type := Function_Under_Test (Invalid_Input);
begin
   Assert (Result_Type.Is_Error (Result), "Error case returns Error");

   declare
      Err : constant Error_Type := Result_Type.Error_Info (Result);
   begin
      Assert (Err.Kind = Validation_Error,
              "Error has correct kind");
   end;
end;
```

### 11.2 Testing Value Objects

```ada
-- Test validation
declare
   Valid   : constant Result := Create ("Valid Name");
   Invalid : constant Result := Create ("");
begin
   Assert (Is_Ok (Valid), "Valid input accepted");
   Assert (Is_Error (Invalid), "Invalid input rejected");
end;

-- Test immutability (compile-time check)
-- This won't compile (good!):
-- Set_Name (Person, "New Name");  -- No setter exists
```

### 11.3 Testing Use Cases

```ada
-- Test with mock/stub adapter
package Mock_Writer is
   Was_Called : Boolean := False;
   Last_Message : Unbounded_String;

   function Write (Msg : String) return Unit_Result.Result is
   begin
      Was_Called := True;
      Last_Message := To_Unbounded_String (Msg);
      return Unit_Result.Ok (Unit_Value);
   end Write;
end Mock_Writer;

-- Instantiate use case with mock
package UseCase_Under_Test is new Application.Usecase.Greet
  (Writer => Mock_Writer.Write);

-- Test
declare
   Cmd    : constant Greet_Command := Make_Command ("Alice");
   Result : constant Unit_Result.Result := UseCase_Under_Test.Execute (Cmd);
begin
   Assert (Unit_Result.Is_Ok (Result), "Use case succeeds");
   Assert (Mock_Writer.Was_Called, "Writer adapter was called");
   Assert (Index (Mock_Writer.Last_Message, "Alice") > 0,
           "Message contains name");
end;
```

---

## 12. Troubleshooting

### 12.1 Common Issues

**Q: Tests fail to compile**

A: Ensure you've built the main project first:
```bash
make build
make build-tests
```

**Q: Test runner not found**

A: Build tests explicitly:
```bash
make build-tests
ls test/bin/
```

**Q: All tests fail with similar errors**

A: Check if test framework was updated. Rebuild all:
```bash
make clean
make build build-tests
```

**Q: Coverage target fails**

A: Ensure GNATcoverage runtime is built:
```bash
make build-coverage-runtime
make test-coverage
```

---

## 13. Future Enhancements

### 13.1 Planned Improvements (Roadmap)

- **Property-Based Testing**: Test invariants with random inputs
- **Mutation Testing**: Verify test suite catches bugs
- **Performance Benchmarks**: Track performance regressions
- **Parallel Test Execution**: Speed up test runs
- **Test Data Builders**: Simplify test setup

---

**Document Control**:
- Version: 1.0.0
- Last Updated: November 18, 2025
- Status: Released
- Copyright © 2025 Michael Gardner, A Bit of Help, Inc.
- License: BSD-3-Clause
- SPDX-License-Identifier: BSD-3-Clause
