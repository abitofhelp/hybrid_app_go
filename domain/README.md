# Domain Layer

The innermost layer containing pure business logic with **zero external dependencies**.

## Responsibilities

- Define core business entities and value objects
- Implement business rules and validation
- Define error types and Result monad for functional error handling
- Remain completely isolated from infrastructure concerns

## Key Packages

- `error/` - Error types and Result[T] monad implementation
- `valueobject/` - Immutable value objects (Person, Option[T])

## Architectural Rules

- **No imports** from application, infrastructure, presentation, or bootstrap layers
- **No external dependencies** - only Go standard library
- All types should be immutable where possible
- Use Result[T] for operations that can fail

## Example

```go
// Creating a Person value object with validation
person, err := valueobject.NewPerson("Alice")
if err != nil {
    // Handle validation error
}
```
