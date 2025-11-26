<!-- SPDX-License-Identifier: BSD-3-Clause -->

# Presentation Layer

Implements inbound adapters that handle user interaction and call application use cases.

## Responsibilities

- Handle user input (CLI arguments, HTTP requests, etc.)
- Validate input format before passing to application layer
- Call application use cases via inbound ports
- Format and display results to users
- Convert application results to appropriate output format

## Key Packages

- `adapter/cli/command/` - CLI command handlers

## Architectural Rules

- **Can import**: application layer only (via ports and error re-exports)
- **CANNOT import**: domain layer directly
- Access domain types only through application layer re-exports
- Use `application/error` for error types, not `domain/error`

## Why No Domain Access?

The presentation layer must remain decoupled from domain internals:

```go
// CORRECT - use application error re-exports
import apperr "github.com/.../application/error"

// WRONG - direct domain access
import domerr "github.com/.../domain/error"  // NOT ALLOWED
```

This ensures domain changes don't ripple to presentation layer.

## Example Command Handler

```go
type GreetCommand[UC inbound.GreetPort] struct {
    useCase UC
}

func (c *GreetCommand[UC]) Run(args []string) int {
    cmd := command.NewGreetCommand(args[1])
    result := c.useCase.Execute(ctx, cmd)

    if result.IsError() {
        fmt.Fprintln(os.Stderr, result.ErrorInfo().Message)
        return 1
    }
    return 0
}
```
