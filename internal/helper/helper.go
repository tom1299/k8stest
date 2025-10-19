// Package helper provides internal helper functions for k8stest.
// This package is internal and should not be imported by external code.
package helper

import (
	"errors"
	"fmt"
	"strings"
)

// ValidateNamespace checks if a namespace is valid.
// Returns an error if the namespace is empty or invalid.
func ValidateNamespace(namespace string) error {
	if namespace == "" {
		return errors.New("namespace cannot be empty")
	}
	if strings.Contains(namespace, " ") {
		return errors.New("namespace cannot contain spaces")
	}
	return nil
}

// FormatNamespace formats a namespace string with a prefix.
func FormatNamespace(namespace string) string {
	return fmt.Sprintf("k8s-namespace:%s", namespace)
}
