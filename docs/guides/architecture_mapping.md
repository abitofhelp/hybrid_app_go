# Architecture Mapping: Go ↔ Ada

**Purpose**: Educational tool showing how the same hexagonal architecture is implemented in Go and Ada using each language's idioms.

## Core Principle

Both implementations follow **identical architecture** (Hybrid DDD/Clean/Hexagonal) but use **language-specific idioms** for implementation.

---

## File-by-File Mapping

### Domain Layer (Pure Business Logic - ZERO external module dependencies)

| Go | Ada | Notes |
|----|-----|-------|
| `domain/person.go` | `domain/domain-person.ads/adb` | ✅ Same: Person entity/value object |
| N/A (uses Go errors) | `domain/domain-result.ads/adb` | **Ada-specific**: Generic Result type (Go uses `mo.Result` external library) |

**Key Differences**:
- **Go**: Uses standard `(T, error)` return pattern in domain (pure Go, no external deps)
- **Ada**: Uses `Generic_Result` monad in domain (still zero external deps - defined in domain itself)
- **Why**: Go has built-in multiple returns; Ada uses discriminated records for Result pattern

**Naming**:
- ✅ Go: `Person` → Ada: `Person` (type name matches)
- ✅ Go: `NewPerson(name)` → Ada: `Create(name)` (Ada idiom: use `Create` for constructors)
- ✅ Go: `Name()` → Ada: `Get_Name()` (Ada idiom: prefix accessors with `Get_`)
- ✅ Go: `GreetingMessage()` → Ada: `Greeting_Message()` (matches exactly)

---

### Application Layer (Use Cases + Ports)

| Go | Ada | Notes |
|----|-----|-------|
| `application/model/unit.go` | `application/application-unit_type.ads` | ✅ Same: Unit type for void returns |
| `application/port/outbound/writer.go` | `application/application-console_port.ads` | ⚠️ **Different pattern** (see below) |
| `application/command/greet.go` | **MISSING** | ❌ Missing DTO |
| `application/port/inbound/greet.go` | **MISSING** | ❌ Missing input port interface |
| `application/usecase/greet.go` | `application/usecase/application-usecase-greet.ads/adb` | ✅ Same concept, different idiom |

**Key Differences - Ports Pattern**:

**Go (Runtime Polymorphism with Static Dispatch)**:
```go
// Output Port Interface
type WriterPort interface {
    Write(ctx context.Context, message string) mo.Result[Unit]
}

// Input Port Interface
type GreetPort interface {
    Execute(ctx context.Context, cmd GreetCommand) mo.Result[Unit]
}

// Use Case with Generic (Static Dispatch)
type GreetUseCase[W WriterPort] struct {
    writer W  // Concrete type known at compile time
}
```

**Ada (Compile-time Polymorphism with Generics)**:
```ada
-- Output Port via Generic Function Parameter
generic
   with function Write (Message : String) return Unit_Result.Result;
package Port is
   function Write_Message (Message : String) return Unit_Result.Result
     renames Write;
end Port;

-- Use Case with Generic (Static Dispatch)
generic
   with function Console_Write (Message : String)
     return Console_Port.Unit_Result.Result;
package Use_Case is
   function Execute (Name : String) return Unit_Result.Result;
end Use_Case;
```

**Educational Explanation**:
- **Both use static dispatch** via generics (compile-time binding)
- **Go**: Uses interface constraints on generic types (`W WriterPort`)
- **Ada**: Uses generic formal parameters (`with function`)
- **Both achieve**: Zero runtime overhead, type safety, dependency inversion
- **Difference**: Go has interface + generic, Ada has pure generic (more direct)

**Missing Components**:
1. ❌ **GreetCommand DTO** - Go has explicit DTO struct, Ada passes `String` directly
2. ❌ **GreetPort interface** - Ada uses generic function parameter instead

**Should we add these for educational clarity?** Yes, for 1:1 mapping.

---

### Infrastructure Layer (Adapters)

| Go | Ada | Notes |
|----|-----|-------|
| `infrastructure/adapter/consolewriter.go` | `infrastructure/adapter/infrastructure-adapter-console_writer.ads/adb` | ✅ Same: Console output adapter |

**Naming**:
- ✅ Go: `ConsoleWriterAdapter` → Ada: `Console_Writer` (matches, Ada uses `_` separator)
- ✅ Go: `Write(ctx, message)` → Ada: `Write(message)` (Ada omits `context` - not idiomatic in Ada)
- ✅ Go: Returns `mo.Result[Unit]` → Ada: Returns `Unit_Result.Result` (same pattern)

**Key Difference - Context**:
- **Go**: Passes `context.Context` for cancellation/timeout
- **Ada**: Does not use context (Ada has built-in task cancellation, not typically passed as parameter)
- **Educational Note**: Context pattern is Go-specific; Ada uses tasks with built-in cancellation

---

### Presentation Layer (UI/CLI)

| Go | Ada | Notes |
|----|-----|-------|
| `presentation/adapter/cli/command/greet.go` | `presentation/presentation-cli_controller.ads/adb` | ✅ Same: CLI command handler |

**Naming**:
- ⚠️ Go: `GreetCommand` → Ada: `CLI_Controller` (different name, same role)
- ✅ Go: `Execute(ctx, name)` → Ada: `Run()` (different signature, same purpose)

**Key Differences**:
- **Go**: `GreetCommand` struct holds `GreetPort` interface field
- **Ada**: `CLI_Controller` is a generic package that receives execute function as parameter
- **Go**: `Execute` method returns `int` (exit code)
- **Ada**: `Run` function returns `Integer` (exit code)

**Should align**: Rename Ada's `CLI_Controller` to `Greet_Command` for consistency.

---

### Bootstrap/Composition Root

| Go | Ada | Notes |
|----|-----|-------|
| `bootstrap/cli/cli.go` | **Embedded in main** | ⚠️ Different structure |
| `cmd/greeter/main.go` | `src/greeter.adb` | ✅ Same: Main entry point |

**Key Differences**:
- **Go**: Has separate `bootstrap/cli/cli.go` package that wires dependencies
- **Ada**: Wiring done directly in `greeter.adb` (no separate bootstrap package)
- **Both**: Perform same function - instantiate generics and wire dependencies

**Educational Note**:
- Go separates bootstrap into its own package for reusability
- Ada typically does composition in main (simpler, common pattern)
- Could extract Ada bootstrap to separate package for exact matching

---

## Architectural Patterns Comparison

### Static Dispatch via Generics (SAME in both)

**Go (1.18+ Generics)**:
```go
type GreetUseCase[W WriterPort] struct {
    writer W  // Compile-time binding
}

// Bootstrap instantiates with concrete type
useCase := NewGreetUseCase[*ConsoleWriterAdapter](consoleWriter)
```

**Ada (Generics since Ada 83)**:
```ada
generic
   with function Console_Write(...) return Result;
package Use_Case is
   -- Compile-time binding
end Use_Case;

-- Bootstrap instantiates with concrete function
package Greet_UC is new Use_Case(Console_Write => Writer.Write);
```

**Educational Point**: Both achieve compile-time polymorphism with zero runtime overhead.

---

## Missing Ada Components (For Educational Parity)

### 1. ❌ GreetCommand DTO
**Go has**: `application/command/greet.go`
```go
type GreetCommand struct {
    Name string
}
```

**Ada should add**: `application/application-greet_command.ads`
```ada
type Greet_Command is record
   Name : String(1..Max_Name_Length);
end record;
```

**Why**: Shows proper DTO pattern for crossing architectural boundaries.

### 2. ❌ GreetPort Input Port Interface
**Go has**: `application/port/inbound/greet.go`
```go
type GreetPort interface {
    Execute(ctx context.Context, cmd GreetCommand) mo.Result[Unit]
}
```

**Ada equivalent**: Already handled via generic parameter, but could document as conceptual "Input Port"

**Educational Note**: Ada's generic `with function` parameter **IS** the input port contract. It's just expressed differently.

---

## Summary: What Maps Directly vs What Differs

### ✅ Direct Mappings (Same Concept, Same Name):
1. Person entity/value object
2. Unit type
3. Console Writer Adapter
4. Use Case (Greet)
5. Main entry point
6. Result monad pattern

### ⚠️ Same Concept, Different Idiom:
1. **Ports**: Go uses interfaces, Ada uses generic formal parameters (both achieve static dispatch)
2. **Context**: Go passes explicitly, Ada uses built-in task features
3. **Error Handling**: Go domain uses `(T, error)`, Ada domain uses `Result[T]`
4. **Bootstrap**: Go has separate package, Ada embedded in main

### ❌ Missing in Ada (Should Add for Parity):
1. **GreetCommand DTO** - for showing explicit boundary crossing
2. **Separate Bootstrap package** - for showing composition root pattern clearly

### 📚 Educational Value:
- Students learn how **same architecture** can be implemented with **different language idioms**
- Shows **static dispatch** can be achieved via interfaces+generics (Go) or pure generics (Ada)
- Demonstrates **architectural patterns transcend language syntax**

---

## Recommendations for Educational Parity

1. ✅ **Add GreetCommand DTO** to Ada version
2. ✅ **Rename CLI_Controller → Greet_Command** to match Go
3. ✅ **Extract Bootstrap to separate package** in Ada (optional but clearer)
4. ✅ **Add architecture mapping comments** in both codebases
5. ✅ **Create side-by-side examples** showing equivalent patterns

Would you like me to implement these changes?
