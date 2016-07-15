// Package validation contains validation tools.
package validation

import (
	"k8s.io/kubernetes/pkg/util/validation"
)

// ValidateUID returns nil if valid otherwise a string list of issues.
func ValidateUID(uid string) (issues []string) {
	if len(uid) != 36 {
		issues = append(issues, "UID must be 36 characters long")
	}
	// TODO: Add more validations.
	return issues
}

// ValidateQualitfiedName validates a Qualified Name used in Kubernetes.
func ValidateQualitfiedName(name string) []string {
	return validation.IsQualifiedName(name)
}
