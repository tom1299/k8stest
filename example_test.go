package k8stest_test

import (
	"fmt"
	"log"

	"github.com/tom1299/k8stest"
)

// Example demonstrates basic usage of the k8stest library.
func Example() {
	// Create a new config with defaults
	config := k8stest.NewConfig()
	fmt.Printf("Default namespace: %s\n", config.Namespace)
	fmt.Printf("Default timeout: %d\n", config.Timeout)

	// Validate the config
	if err := config.ValidateConfig(); err != nil {
		log.Fatal(err)
	}

	// Get formatted namespace
	fmt.Printf("Formatted: %s\n", config.GetFormattedNamespace())

	// Output:
	// Default namespace: default
	// Default timeout: 30
	// Formatted: k8s-namespace:default
}

// ExampleConfig_ValidateConfig demonstrates config validation.
func ExampleConfig_ValidateConfig() {
	config := &k8stest.Config{
		Namespace: "my-namespace",
		Timeout:   60,
	}

	if err := config.ValidateConfig(); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		return
	}

	fmt.Println("Config is valid")
	// Output: Config is valid
}

// ExampleNewConfig demonstrates creating a new config.
func ExampleNewConfig() {
	config := k8stest.NewConfig()
	fmt.Printf("Namespace: %s, Timeout: %d\n", config.Namespace, config.Timeout)
	// Output: Namespace: default, Timeout: 30
}
