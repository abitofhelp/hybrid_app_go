# Bootstrap Layer

The composition root that wires all dependencies together using static dispatch.

## Responsibilities

- Instantiate concrete implementations (adapters)
- Wire dependencies using Go generics (static dispatch)
- Create and configure the application entry point
- Handle application lifecycle (startup, shutdown)

## Key Packages

- `cli/` - CLI application bootstrap and runner

## Architectural Rules

- **Can import**: all layers (domain, application, infrastructure, presentation)
- This is the only layer that knows about all concrete types
- Dependencies are wired at compile-time via generics

## Static Dispatch Pattern

```go
// 1. Create concrete adapter
writer := &adapter.ConsoleWriter{}

// 2. Instantiate use case with concrete type
useCase := usecase.NewGreetUseCase[*adapter.ConsoleWriter](writer)

// 3. Instantiate command handler with concrete use case type
cmd := command.NewGreetCommand[*usecase.GreetUseCase[*adapter.ConsoleWriter]](useCase)

// 4. Run
return cmd.Run(os.Args)
```

## Benefits

- **Zero runtime overhead** - no interface dispatch, no reflection
- **Compile-time safety** - type mismatches caught at build time
- **Clear dependency graph** - all wiring visible in one place
