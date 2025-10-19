# k8stest

Minimal Go helpers for composing Kubernetes test resources with a fluent API.

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
# k8stest

Minimal Go helpers for composing Kubernetes test resources with a fluent API.

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

## Git tips: Undo a commit but keep the changes

- If the commit is local (not pushed) and you want to keep the changes staged for commit:
  - git reset --soft HEAD~1
  - Effect: Moves HEAD back one commit, leaves all your changes staged in the index.

- If the commit is local (not pushed) and you want to keep the changes but unstaged (in your working tree):
  - git reset HEAD~1
  - Same as --mixed (default). Your files become regular modifications you can edit and re-stage.

- If you need to undo multiple recent commits locally, adjust the count:
  - git reset --soft HEAD~3   # keep changes from last 3 commits staged
  - git reset HEAD~3          # keep them unstaged

- If the commit has already been pushed and you must keep history intact (recommended for shared branches):
  - git revert <commit>
  - Effect: creates a new commit that undoes the specified commit. Safer for collaborators.

- If you pushed but truly need to rewrite history (last resort):
  - Save collaborators first, then:
    - git reset --soft <commit^>   # move HEAD before the commit, keep changes staged
    - git push --force-with-lease  # rewrite remote history safely
  - Warning: this rewrites history and can disrupt others. Prefer revert on shared branches.

Notes
- --soft keeps changes staged; --mixed (default) keeps them unstaged; --hard discards changes (dangerous; do not use if you need the changes).
- You can replace HEAD~1 with an explicit commit hash or a different ancestor spec.

## License

This repository is for demonstration and testing utilities.
