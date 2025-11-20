// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: cli
// Description: CLI bootstrap and dependency wiring

// Package cli provides the composition root for the CLI application.
// This is where all dependencies are wired together (dependency injection).
//
// Architecture Notes:
//   - Part of the BOOTSTRAP layer (composition root)
//   - Depends on ALL layers to wire dependencies together
//   - This is the ONLY place where all layers meet
//   - Performs static dependency injection
//   - No business logic here (only wiring)
//
// Dependency Wiring Flow:
//  1. Infrastructure → Application ports (Writer adapter)
//  2. Application → Domain (Use case with domain logic)
//  3. Presentation → Application (CLI command with use case)
//  4. Main → Bootstrap (Entry point calls Run)
//
// Usage:
//
//	import "github.com/abitofhelp/hybrid_app_go/bootstrap/cli"
//
//	func main() {
//	    exitCode := cli.Run(os.Args)
//	    os.Exit(exitCode)
//	}
package cli

import (
	"github.com/abitofhelp/hybrid_app_go/application/usecase"
	"github.com/abitofhelp/hybrid_app_go/infrastructure/adapter"
	"github.com/abitofhelp/hybrid_app_go/presentation/cli/command"
)

// Run is the composition root that wires all dependencies and executes the application.
//
// This function demonstrates the complete dependency injection flow:
//
//	Step 1: Wire Infrastructure → Application ports
//	  - Infrastructure.ConsoleWriter → Application.WriterFunc (output port)
//
//	Step 2: Wire Application use case with injected dependencies
//	  - Application.GreetUseCase with WriterFunc (from step 1)
//
//	Step 3: Wire Presentation command with use case
//	  - Presentation.GreetCommand with GreetUseCase.Execute (from step 2)
//
//	Step 4: Run the application
//	  - Call GreetCommand.Run with command-line arguments
//	  - Return exit code to caller
//
// Flow of data through the architecture:
//
//	1. User runs: ./greeter Alice
//	2. Main calls Bootstrap.Run with os.Args
//	3. Bootstrap wires all dependencies (this function)
//	4. GreetCommand parses args and extracts "Alice"
//	5. GreetCommand creates GreetCommand DTO
//	6. GreetCommand calls GreetUseCase.Execute(GreetCommand)
//	7. GreetUseCase extracts name from DTO
//	8. GreetUseCase calls Domain.Person.CreatePerson("Alice")
//	9. Domain validates the name
//	10. GreetUseCase gets greeting message from Person
//	11. GreetUseCase calls WriterFunc("Hello, Alice!")
//	12. WriterFunc routes to ConsoleWriter (via injection)
//	13. ConsoleWriter calls fmt.Println
//	14. Result flows back through layers:
//	    Writer → UseCase → Command → Bootstrap → Main
//	15. Main returns exit code to shell
//
// Architectural Benefits:
//   - All layers remain independent (loose coupling)
//   - Dependencies point inward (Dependency Rule)
//   - Easy to swap implementations (e.g., different writer)
//   - Testable (inject mock implementations)
//   - Clear separation of concerns
//
// Contract:
//   - Pre: args is os.Args (program name + arguments)
//   - Post: Returns 0 if application succeeded
//   - Post: Returns non-zero if application failed
func Run(args []string) int {
	// ========================================================================
	// Step 1: Wire Infrastructure → Application ports
	// ========================================================================

	// DEPENDENCY INVERSION in action:
	// - Application.Port.Outward.WriterFunc defines the interface (port)
	// - Infrastructure.Adapter.ConsoleWriter provides implementation
	// - We wire them together here in the composition root
	consoleWriter := adapter.NewConsoleWriter()

	// ========================================================================
	// Step 2: Wire Application use case with injected dependencies
	// ========================================================================

	// The use case receives the Writer function through constructor injection.
	// This is FUNCTION INJECTION - a lightweight Go pattern for dependency injection.
	greetUseCase := usecase.NewGreetUseCase(consoleWriter)

	// ========================================================================
	// Step 3: Wire Presentation command with use case
	// ========================================================================

	// Wire the Presentation layer to the Application layer.
	// The command receives the Execute function from the use case.
	// Again, function injection - zero runtime overhead.
	greetCommand := command.NewGreetCommand(greetUseCase.Execute)

	// ========================================================================
	// Step 4: Run the application and return exit code
	// ========================================================================

	// Call the Greet Command to start the application.
	// The command will:
	//   1. Parse command-line arguments
	//   2. Create GreetCommand DTO
	//   3. Call the use case (which calls domain and console port)
	//   4. Return an exit code
	return greetCommand.Run(args)
}
