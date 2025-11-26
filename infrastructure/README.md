<!-- SPDX-License-Identifier: BSD-3-Clause -->

# Infrastructure Layer

Implements outbound ports with concrete adapters for external services.

## Responsibilities

- Implement outbound port interfaces defined in application layer
- Handle I/O operations (console, file, network, database)
- Convert infrastructure errors to domain error types
- Provide panic recovery at system boundaries

## Key Packages

- `adapter/` - Concrete implementations of outbound ports

## Architectural Rules

- **Can import**: domain, application layers
- **Cannot import**: presentation, bootstrap
- Must implement interfaces from `application/port/outbound/`
- Convert all errors to domain.error.ErrorType
- Recover from panics and convert to Result errors

## Example Adapter

```go
// ConsoleWriter implements outbound.WriterPort
type ConsoleWriter struct{}

func (w *ConsoleWriter) Write(ctx context.Context, msg string) Result[model.Unit] {
    // Panic recovery wrapper
    defer func() {
        if r := recover(); r != nil {
            // Convert panic to Result error
        }
    }()

    fmt.Println(msg)
    return Ok(model.Unit{})
}
```
