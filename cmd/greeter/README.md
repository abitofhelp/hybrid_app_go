<!-- SPDX-License-Identifier: BSD-3-Clause -->

# Greeter CLI Application

Main entry point for the greeter application.

## Usage

```bash
# Build
go build -o greeter ./cmd/greeter

# Run
./greeter Alice
# Output: Hello, Alice!

# Error case
./greeter ""
# Output: Error: name cannot be empty
```

## Structure

```
cmd/greeter/
└── main.go    # Entry point - delegates to bootstrap.CLI.Run()
```

## Implementation

The main function is minimal - it simply delegates to the bootstrap layer:

```go
func main() {
    exitCode := cli.Run(os.Args)
    os.Exit(exitCode)
}
```

All dependency wiring and application logic lives in the bootstrap layer,
keeping `main.go` as a thin entry point.
