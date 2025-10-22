This is a Go-based library that provides minimal helpers for composing Kubernetes test resources with a fluent API. It is designed to simplify the creation and management of Kubernetes resources (Deployments, ConfigMaps, Secrets) in test environments. Please follow these guidelines when contributing:

## Code Standards

### Required Before Each Commit
- Run `go fmt` before committing any changes to ensure proper code formatting
- Ensure all tests pass with `go test ./...`

### Development Flow
- Test: `go test ./...`
- Format: `go fmt ./...`
- Vet: `go vet ./...`

## Repository Structure
- Root package (`k8stest.go`): Core fluent API for building and managing Kubernetes resources
- `internal/`: Internal helper functions for building Kubernetes clients
- `k8stest_test.go`: Test cases demonstrating usage patterns

## Key Guidelines
1. Follow Go best practices and idiomatic patterns
2. Maintain the fluent API design - all methods should be chainable
3. Keep the library minimal and focused on testing scenarios
4. Write unit tests for new functionality using table-driven tests when possible
5. All resource operations should work gracefully with non-existent resources (e.g., delete operations should not fail if resource doesn't exist)
6. Document public APIs with clear examples
7. Default namespace for k8s tests is "default" for all tests
