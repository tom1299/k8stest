# k8stest

A minimal Go library demonstrating root/internal package structure for Kubernetes testing utilities.

## Features

- Root package with exported types and functions
- Internal package with helper functions (not accessible to external consumers)
- Comprehensive test coverage
- Example usage documentation

## Installation

```bash
go get github.com/tom1299/k8stest
```

## Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/tom1299/k8stest"
)

func main() {
    // Create a new config with defaults
    config := k8stest.NewConfig()
    
    // Validate the config
    if err := config.ValidateConfig(); err != nil {
        log.Fatal(err)
    }
    
    // Get formatted namespace
    fmt.Println(config.GetFormattedNamespace())
}
```

## Package Structure

```
k8stest/
├── k8stest.go           # Root package with exported API
├── k8stest_test.go      # Tests for root package
├── example_test.go      # Example usage
└── internal/
    └── helper/
        ├── helper.go       # Internal helper functions
        └── helper_test.go  # Tests for internal package
```

The `internal/` directory is special in Go - packages inside it can only be imported by code in the parent directory and its subdirectories, making it perfect for implementation details that shouldn't be part of the public API.

## Testing

Run tests with:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Building

This is a library package, so there's no main binary to build. To verify everything compiles:

```bash
go build ./...
```

## License

This is a minimal example repository for demonstration purposes.