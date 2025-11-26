// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: cli
// Description: CLI bootstrap and dependency wiring

// Package cli provides the composition root for the CLI application.
// This is where all dependencies are wired together via GENERIC INSTANTIATION.
//
// Architecture Notes:
//   - Part of the BOOTSTRAP layer (composition root)
//   - Depends on ALL layers to wire dependencies together
//   - This is the ONLY place where all layers meet
//   - Performs STATIC DEPENDENCY INJECTION via generics
//   - No business logic here (only wiring)
//   - Enables STATIC DISPATCH (compile-time method resolution)
//
// Static Dispatch Pattern:
//   - Infrastructure: *adapter.ConsoleWriter implements WriterPort
//   - Use Case: usecase.GreetUseCase[*adapter.ConsoleWriter]
//   - Command: command.GreetCommand[*usecase.GreetUseCase[*adapter.ConsoleWriter]]
//   - All method calls are resolved at compile time (no vtable)
//
// Mapping to Ada:
//   - Ada: Bootstrap.CLI instantiates generic packages in declarative region
//   - Go: Bootstrap.Run instantiates generic types with concrete parameters
//   - Both achieve: static dispatch, compile-time resolution, zero overhead
//
// Dependency Wiring Flow:
//  1. Infrastructure → Application ports (ConsoleWriter implements WriterPort)
//  2. Application → Domain (GreetUseCase[*ConsoleWriter] coordinates domain)
//  3. Presentation → Application (GreetCommand[*GreetUseCase[*ConsoleWriter]])
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
	"github.com/abitofhelp/hybrid_app_go/presentation/adapter/cli/command"
)

// Run is the composition root that wires all dependencies and executes the application.
//
// This function demonstrates STATIC DEPENDENCY INJECTION via generics:
//
//	Step 1: Create Infrastructure adapter
//	  - adapter.NewConsoleWriter() returns *adapter.ConsoleWriter
//	  - ConsoleWriter implements WriterPort interface
//
//	Step 2: Instantiate Use Case with concrete type
//	  - usecase.NewGreetUseCase[*adapter.ConsoleWriter](writer)
//	  - Compiler knows concrete type → static dispatch
//
//	Step 3: Instantiate Command with concrete use case type
//	  - command.NewGreetCommand[*usecase.GreetUseCase[*adapter.ConsoleWriter]](uc)
//	  - Full type chain is known at compile time
//
//	Step 4: Run the application
//	  - Call GreetCommand.Run with command-line arguments
//	  - All method calls are statically dispatched
//
// Mapping to Ada Bootstrap:
//
//	Ada (bootstrap-cli.adb):
//	  package Writer_Port_Instance is new Generic_Writer(Write => Console_Writer.Write);
//	  package Greet_UC_Instance is new Application.Usecase.Greet(Writer => Writer_Port_Instance.Write_Message);
//	  package Greet_Cmd_Instance is new Presentation.CLI.Command.Greet(Execute => Greet_UC_Instance.Execute);
//
//	Go (this file):
//	  consoleWriter := adapter.NewConsoleWriter()
//	  greetUseCase := usecase.NewGreetUseCase[*adapter.ConsoleWriter](consoleWriter)
//	  greetCommand := command.NewGreetCommand[*usecase.GreetUseCase[*adapter.ConsoleWriter]](greetUseCase)
//
// Flow of data through the architecture:
//
//  1. User runs: ./greeter Alice
//  2. Main calls Bootstrap.Run with os.Args
//  3. Bootstrap instantiates generics with concrete types (this function)
//  4. GreetCommand parses args and extracts "Alice"
//  5. GreetCommand creates GreetCommand DTO
//  6. GreetCommand calls GreetUseCase.Execute(GreetCommand) [STATIC DISPATCH]
//  7. GreetUseCase extracts name from DTO
//  8. GreetUseCase calls Domain.Person.CreatePerson("Alice")
//  9. Domain validates the name
//  10. GreetUseCase gets greeting message from Person
//  11. GreetUseCase calls ConsoleWriter.Write("Hello, Alice!") [STATIC DISPATCH]
//  12. ConsoleWriter.Write() writes to stdout
//  13. Result flows back through layers:
//     Writer → UseCase → Command → Bootstrap → Main
//  14. Main returns exit code to shell
//
// Architectural Benefits:
//   - STATIC DISPATCH: All method calls resolved at compile time
//   - Zero runtime overhead (no interface vtable lookups)
//   - All layers remain independent (loose coupling)
//   - Dependencies point inward (Dependency Rule)
//   - Type safety verified at compile time
//   - Easy to swap implementations (change generic parameters)
//   - Testable (inject mock implementations as type parameters)
//
// Contract:
//   - Pre: args is os.Args (program name + arguments)
//   - Post: Returns 0 if application succeeded
//   - Post: Returns non-zero if application failed
func Run(args []string) int {
	// ========================================================================
	// Step 1: Create Infrastructure adapter
	// ========================================================================

	// DEPENDENCY INVERSION in action:
	// - Application.Port.Outward.WriterPort defines the interface (port)
	// - Infrastructure.Adapter.ConsoleWriter implements the interface
	// - We instantiate the concrete type here in the composition root
	consoleWriter := adapter.NewConsoleWriter()

	// ========================================================================
	// Step 2: Instantiate Use Case with concrete writer type
	// ========================================================================

	// STATIC DISPATCH via generics:
	// - GreetUseCase[*adapter.ConsoleWriter] knows the concrete writer type
	// - All calls to writer.Write() are statically dispatched
	// - Equivalent to Ada: package Greet_UC is new Greet(Writer => Console_Writer.Write)
	greetUseCase := usecase.NewGreetUseCase[*adapter.ConsoleWriter](consoleWriter)

	// ========================================================================
	// Step 3: Instantiate Command with concrete use case type
	// ========================================================================

	// STATIC DISPATCH continues through the chain:
	// - GreetCommand knows the exact use case type
	// - All calls to useCase.Execute() are statically dispatched
	// - The entire call chain is resolved at compile time
	greetCommand := command.NewGreetCommand[*usecase.GreetUseCase[*adapter.ConsoleWriter]](greetUseCase)

	// ========================================================================
	// Step 4: Run the application and return exit code
	// ========================================================================

	// Call the Greet Command to start the application.
	// The command will:
	//   1. Parse command-line arguments
	//   2. Create GreetCommand DTO
	//   3. Call the use case (STATIC DISPATCH to Execute)
	//   4. Use case calls writer (STATIC DISPATCH to Write)
	//   5. Return an exit code
	return greetCommand.Run(args)
}
