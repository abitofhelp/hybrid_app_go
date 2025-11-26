# Software Test Guide

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

This Software Test Guide describes the testing approach, test organization, execution procedures, and guidelines for the Hybrid_App_Go project.

### 1.2 Scope

This document covers:
- Test strategy and organization (unit/integration)
- Go testing framework usage
- Running tests via Make and go test
- Writing new tests
- Coverage analysis
- Test maintenance procedures

---

## 2. Test Strategy

### 2.1 Testing Levels

Hybrid_App_Go uses two levels of testing:

**Unit Tests** (42 assertions in 2 test functions)
- Test individual components in isolation
- Focus on Domain layer (pure functions)
- Predictable, deterministic results
- Fast execution
- Location: `domain/*_test.go`

**Integration Tests** (23 tests)
- Test complete CLI application flow
- Run actual greeter binary
- Verify stdout, stderr, and exit codes
- Black-box testing approach
- Location: `test/integration/`

### 2.2 Testing Philosophy

- **Integration-First**: CLI apps are best tested via actual execution
- **Domain Unit Tests**: Pure functions tested in isolation
- **No Mocks**: Integration tests run the real binary
- **Railway-Oriented**: Test both success and error paths
- **Comprehensive**: Cover normal, edge, and error cases
- **Automated**: All tests runnable via `make test`
- **Fast**: All tests execute in < 3 seconds

---

## 3. Test Organization

### 3.1 Directory Structure

```
hybrid_app_go/
|-- domain/
|   |-- error/
|   |   |-- result_test.go      # Result monad unit tests (19 assertions)
|   |   |-- main_test.go        # Test runner
|   |-- valueobject/
|       |-- person_test.go      # Person value object tests (23 assertions)
|       |-- main_test.go        # Test runner
|
|-- test/
    |-- integration/
        |-- greet_flow_test.go  # CLI integration tests (23 tests)
        |-- go.mod              # Integration test module
```

### 3.2 Test Naming Convention

- **Pattern**: `*_test.go`
- **Test Functions**: `Test<Component>_<Scenario>` or `Test<Component>`
- **Examples**:
  - `result_test.go` -> Tests `domain/error.Result[T]`
  - `person_test.go` -> Tests `domain/valueobject.Person`
  - `greet_flow_test.go` -> Tests CLI greeter flow

---

## 4. Go Testing Framework

### 4.1 Standard Testing

Hybrid_App_Go uses Go's standard `testing` package with `testify` assertions:

**Benefits**:
- Standard Go tooling
- Rich assertion library (testify)
- Table-driven test support
- Parallel test execution
- Coverage analysis built-in

### 4.2 Framework Usage

**Basic Test Structure**:
```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestComponent_Scenario(t *testing.T) {
    // Arrange
    input := "test input"

    // Act
    result := FunctionUnderTest(input)

    // Assert
    assert.True(t, result.IsOk(), "should succeed with valid input")
    assert.Equal(t, expected, result.Value(), "should have correct value")
}
```

### 4.3 Custom Assertion Helpers

The project uses a custom test helper for colored output:

```go
// From domain tests
func check(t *testing.T, condition bool, format string, args ...interface{}) {
    t.Helper()
    if condition {
        fmt.Printf("\033[1;92m[PASS]\033[0m %s\n", fmt.Sprintf(format, args...))
    } else {
        fmt.Printf("\033[1;91m[FAIL]\033[0m %s\n", fmt.Sprintf(format, args...))
        t.Errorf(format, args...)
    }
}
```

---

## 5. Running Tests

### 5.1 Quick Start

```bash
# Run all tests
make test

# Run domain unit tests only
go test -v ./domain/...

# Run integration tests only
go test -v -tags=integration ./test/integration/...
```

### 5.2 Make Targets

**Test Execution**:
```bash
make test               # Run all tests
make test-unit          # Domain unit tests only
make test-integration   # CLI integration tests only
make test-coverage      # Run with coverage analysis
```

**Build and Test**:
```bash
make build && make test  # Build then test
make all                 # Full build + test cycle
```

### 5.3 Direct Execution

```bash
# Domain unit tests
go test -v ./domain/...

# Integration tests (requires build tag)
go test -v -tags=integration ./test/integration/...

# Run specific test
go test -v -run TestGreeter_ValidName ./test/integration/...

# With race detector
go test -race ./domain/...
```

### 5.4 Expected Output

**Domain Tests**:
```
=== RUN   TestDomainErrorResult
[PASS] Ok construction - IsOk returns true
[PASS] Ok construction - IsError returns false
[PASS] Ok value extraction - correct value
[PASS] Error construction - IsError returns true
...
--- PASS: TestDomainErrorResult (0.00s)
PASS
```

**Integration Tests**:
```
=== RUN   TestGreeter_ValidName_Success
--- PASS: TestGreeter_ValidName_Success (0.01s)
=== RUN   TestGreeter_EmptyName_ValidationError
--- PASS: TestGreeter_EmptyName_ValidationError (0.01s)
...
PASS
ok      github.com/abitofhelp/hybrid_app_go/test/integration    0.523s
```

---

## 6. Writing New Tests

### 6.1 Unit Test Template

```go
package valueobject_test

import (
    "testing"
    "github.com/abitofhelp/hybrid_app_go/domain/valueobject"
)

func TestYourComponent(t *testing.T) {
    t.Run("success case", func(t *testing.T) {
        result := valueobject.CreateThing("valid")

        check(t, result.IsOk(), "Valid input should return Ok")
        check(t, result.Value().GetField() == "expected", "Field should match")
    })

    t.Run("error case", func(t *testing.T) {
        result := valueobject.CreateThing("")

        check(t, result.IsError(), "Empty input should return Error")
        check(t, result.ErrorInfo().Kind == domerr.ValidationError, "Should be ValidationError")
    })
}
```

### 6.2 Integration Test Template

```go
//go:build integration

package integration

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestYourScenario_Description(t *testing.T) {
    // Run the greeter binary
    stdout, stderr, exitCode := runGreeter("input")

    // Verify results
    assert.Equal(t, 0, exitCode, "exit code should be 0")
    assert.Equal(t, "Expected Output\n", stdout, "stdout should match")
    assert.Empty(t, stderr, "stderr should be empty")
}
```

### 6.3 Table-Driven Tests

```go
func TestGreeter_ValidNames_TableDriven(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"simple name", "Alice", "Hello, Alice!\n"},
        {"name with space", "John Doe", "Hello, John Doe!\n"},
        {"unicode name", "Jose", "Hello, Jose!\n"},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            stdout, stderr, exitCode := runGreeter(tc.input)

            require.Equal(t, 0, exitCode, "exit code should be 0")
            assert.Equal(t, tc.expected, stdout)
            assert.Empty(t, stderr)
        })
    }
}
```

### 6.4 Adding Tests

**For Domain Tests**:
1. Add test file in appropriate domain package
2. Follow `*_test.go` naming
3. Run with `go test -v ./domain/...`

**For Integration Tests**:
1. Add test function in `test/integration/greet_flow_test.go`
2. Use `//go:build integration` tag
3. Run with `go test -v -tags=integration ./test/integration/...`

---

## 7. Test Coverage

### 7.1 Coverage Goals

- **Target**: > 80% line coverage for Domain layer
- **Critical Code**: 100% coverage for error handling paths
- **Domain Layer**: High coverage (pure functions, easy to test)

### 7.2 Running Coverage Analysis

```bash
# Generate coverage report
go test -cover ./domain/...

# Detailed coverage with HTML report
go test -coverprofile=coverage.out ./domain/...
go tool cover -html=coverage.out -o coverage.html

# Via Make
make test-coverage
```

### 7.3 Coverage Output

```
ok      github.com/abitofhelp/hybrid_app_go/domain/error        0.002s  coverage: 95.0% of statements
ok      github.com/abitofhelp/hybrid_app_go/domain/valueobject  0.002s  coverage: 100.0% of statements
```

---

## 8. Test Maintenance

### 8.1 When to Update Tests

- **New Features**: Add tests before implementing
- **Bug Fixes**: Add regression test first
- **Refactoring**: Ensure tests still pass
- **Requirements Change**: Update affected tests

### 8.2 Test Quality Guidelines

- **Clear Names**: Test names explain what's being verified
- **One Assertion Per Concept**: Group related assertions together
- **Arrange-Act-Assert**: Structure tests clearly
- **No Business Logic**: Tests should be simple
- **Fast Execution**: Avoid slow operations

### 8.3 Debugging Failed Tests

```bash
# Run single test for debugging
go test -v -run TestGreeter_ValidName ./test/integration/...

# Run with verbose output
go test -v ./domain/...

# Run with race detector
go test -race ./domain/...

# Use delve for debugging
dlv test ./domain/valueobject -- -test.run TestDomainValueObjectPerson
```

---

## 9. Integration Test Details

### 9.1 How Integration Tests Work

The integration tests in `test/integration/greet_flow_test.go`:

1. **TestMain**: Builds the greeter binary before tests run
2. **runGreeter**: Executes the binary with arguments
3. **Captures**: stdout, stderr, and exit code
4. **Verifies**: Expected output and exit behavior

```go
func TestMain(m *testing.M) {
    // Build the greeter binary
    projectRoot := findProjectRoot()
    greeterPath = filepath.Join(projectRoot, "greeter_test_binary")

    cmd := exec.Command("go", "build", "-o", greeterPath, "./cmd/greeter")
    cmd.Dir = projectRoot
    if output, err := cmd.CombinedOutput(); err != nil {
        panic("Failed to build greeter: " + err.Error())
    }

    // Run tests
    code := m.Run()

    // Cleanup
    os.Remove(greeterPath)
    os.Exit(code)
}
```

### 9.2 runGreeter Helper

```go
func runGreeter(args ...string) (stdout, stderr string, exitCode int) {
    cmd := exec.Command(greeterPath, args...)

    var stdoutBuf, stderrBuf bytes.Buffer
    cmd.Stdout = &stdoutBuf
    cmd.Stderr = &stderrBuf

    err := cmd.Run()

    stdout = stdoutBuf.String()
    stderr = stderrBuf.String()

    if err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            exitCode = exitErr.ExitCode()
        }
    }

    return
}
```

---

## 10. Continuous Integration

### 10.1 CI Testing Strategy

All tests run on every commit:

```bash
# CI pipeline equivalent
make clean
make build
make test
make check-arch
```

### 10.2 Success Criteria

All must pass:
- Zero build errors
- Zero go vet warnings
- All domain unit tests pass (42 assertions)
- All integration tests pass (23 tests)
- Architecture validation passes
- Exit code 0 from all commands

---

## 11. Test Statistics

### 11.1 Current Test Metrics (v1.0.0)

**Test Count**:
- Total: 25 test functions
  - Domain Unit: 2 functions (42 assertions)
  - Integration: 23 tests
- Pass Rate: 100%

**Coverage**:
- Domain/error: ~95%
- Domain/valueobject: ~100%

**Execution Time**:
- Domain tests: < 0.1 seconds
- Integration tests: < 2 seconds
- Total: < 3 seconds

---

## 12. Common Testing Patterns

### 12.1 Testing Result[T] Monads

```go
// Test success path
t.Run("success case", func(t *testing.T) {
    result := FunctionUnderTest(validInput)

    check(t, result.IsOk(), "Should succeed with valid input")
    check(t, result.Value() == expected, "Should have correct value")
})

// Test error path
t.Run("error case", func(t *testing.T) {
    result := FunctionUnderTest(invalidInput)

    check(t, result.IsError(), "Should fail with invalid input")
    check(t, result.ErrorInfo().Kind == domerr.ValidationError, "Should be ValidationError")
    check(t, strings.Contains(result.ErrorInfo().Message, "expected"), "Should have descriptive message")
})
```

### 12.2 Testing Value Objects

```go
// Test validation
t.Run("validation", func(t *testing.T) {
    valid := valueobject.CreatePerson("Alice")
    invalid := valueobject.CreatePerson("")

    check(t, valid.IsOk(), "Valid input accepted")
    check(t, invalid.IsError(), "Invalid input rejected")
})

// Test immutability (compile-time check)
// person.name = "New" // Won't compile - unexported field
```

### 12.3 Testing CLI Output

```go
func TestGreeter_ValidName_Success(t *testing.T) {
    stdout, stderr, exitCode := runGreeter("Alice")

    assert.Equal(t, 0, exitCode, "exit code should be 0")
    assert.Equal(t, "Hello, Alice!\n", stdout, "stdout should contain greeting")
    assert.Empty(t, stderr, "stderr should be empty")
}

func TestGreeter_EmptyName_ValidationError(t *testing.T) {
    stdout, stderr, exitCode := runGreeter("")

    assert.Equal(t, 1, exitCode, "exit code should be 1")
    assert.Empty(t, stdout, "stdout should be empty")
    assert.Contains(t, stderr, "Error:", "stderr should contain error")
}
```

---

## 13. Troubleshooting

### 13.1 Common Issues

**Q: Integration tests fail to build**

A: Ensure you're using the build tag:
```bash
go test -v -tags=integration ./test/integration/...
```

**Q: Tests fail with "cannot find package"**

A: Update go.work to include test modules:
```bash
go work sync
```

**Q: Integration tests hang**

A: Check if binary path is correct. The tests build from project root.

**Q: Coverage doesn't include all packages**

A: Specify packages explicitly:
```bash
go test -coverprofile=coverage.out ./domain/...
```

---

## 14. Future Enhancements

### 14.1 Planned Improvements

- **Benchmark Tests**: Performance regression detection
- **Fuzz Testing**: Go 1.18+ fuzzing for edge cases
- **Parallel Integration Tests**: Speed up CI
- **Property-Based Testing**: Test invariants with random inputs

---

**Document Control**:
- Version: 1.0.0
- Last Updated: November 25, 2025
- Status: Released
- Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
- License: BSD-3-Clause
- SPDX-License-Identifier: BSD-3-Clause
