// Package k8stest provides utilities for Kubernetes testing.
package k8stest

import "github.com/tom1299/k8stest/internal/helper"

// Config holds configuration for the k8s test utilities.
type Config struct {
	Namespace string
	Timeout   int
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		Namespace: "default",
		Timeout:   30,
	}
}

// ValidateConfig validates a Config and returns any errors.
// It uses internal helper functions to perform validation.
func (c *Config) ValidateConfig() error {
	return helper.ValidateNamespace(c.Namespace)
}

// GetFormattedNamespace returns a formatted namespace string.
// This demonstrates the use of an internal package function.
func (c *Config) GetFormattedNamespace() string {
	return helper.FormatNamespace(c.Namespace)
}
