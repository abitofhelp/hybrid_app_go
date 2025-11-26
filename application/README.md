# Application Layer

Orchestrates use cases and defines port interfaces for dependency inversion.

## Responsibilities

- Define inbound ports (use case interfaces)
- Define outbound ports (dependency interfaces)
- Implement use cases that coordinate domain logic
- Define command/DTO types for crossing layer boundaries
- Re-export domain errors for presentation layer access

## Key Packages

- `port/inbound/` - Use case interfaces (what we offer)
- `port/outbound/` - Dependency interfaces (what we need)
- `usecase/` - Use case implementations
- `command/` - Command/DTO types
- `model/` - Application-specific models (Unit type)
- `error/` - Re-exports domain errors

## Architectural Rules

- **Can import**: domain layer only
- **Cannot import**: infrastructure, presentation, bootstrap
- Use cases accept commands, return Result types
- All dependencies injected via generic type parameters

## Port Pattern

```go
// Inbound port - what clients call
type GreetPort interface {
    Execute(ctx context.Context, cmd command.GreetCommand) Result[model.Unit]
}

// Outbound port - what we need
type WriterPort interface {
    Write(ctx context.Context, message string) Result[model.Unit]
}
```
