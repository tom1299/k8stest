# k8stest

Minimal Go helpers for composing Kubernetes test resources with a fluent API.

> [!CAUTION]
> This code is mainly AI-generated and not intended for production use.


## Features

- Fluent builders for common resources: Deployments, ConfigMaps, and Secrets
- Chainable methods to attach ConfigMaps/Secrets to Deployments
- Simple Create helper that uses Kubernetes clients to create resources in the "default" namespace
- Designed for use in tests

## Installation

```bash
go get github.com/tom1299/k8stest
```

## Quick start (in tests)

This library is intended to be used from tests. The example below shows how tests in this repository use the fluent API to define resources and create them in the current kube context.

```go
package k8stest

import (
    "context"
    "testing"
)

func TestExample(t *testing.T) {
    ctx := context.Background()

    // TestClients provides access to Kubernetes APIs used by Create.
    // In this repository we construct it via a local helper.
    clients := setupTestClients(t)

    // Define resources and create them in the cluster
    _, err := (&Resources{}).
        WithDeployment("web").
            WithConfigMap("app-config").
            WithSecret("tls").
            And().
        WithConfigMap("shared-config").
        Create(clients, ctx)
    if err != nil {
        t.Fatalf("create failed: %v", err)
    }
}
```

Notes
- The Create helper currently targets the "default" namespace and requires valid kubeconfig for your current context.
- The TestClients type is exported, but its constructor in this repo (setupTestClients) is a test helper. External consumers can create their own clients and adapt as needed.

## Package structure

```
k8stest/
├── k8stest.go          # Fluent builders and Create helper
├── k8stest_test.go     # Tests using the fluent API
├── show_output_test.go # Example test that logs output with t.Log
└── README.md
```

## Testing

Run tests with:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

See verbose output from tests (including log lines):

```bash
# t.Log/t.Logf messages are shown when using -v
go test -v ./...
```

Note: Example functions (in files like `example_test.go`) have their standard output compared against the `// Output:` comments. When the output matches, Go suppresses printing it, even with `-v`. If you want to see runtime output during `go test -v`, use `t.Log`/`t.Logf` inside normal `Test*` functions (see `show_output_test.go`).

## Building

This is a library package, so there's no main binary to build. To verify everything compiles:

```bash
go build ./...
```

## License

This repository is for demonstration and testing utilities.
