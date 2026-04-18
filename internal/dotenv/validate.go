package dotenv

import (
	"fmt"
	"regexp"
	"strings"
)

// validKeyRe matches POSIX-style env var names.
var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidationError holds all key-level issues found during validation.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(e.Issues, "; "))
}

// Validate checks that every key in secrets is a valid environment variable
// name and that no value exceeds maxValueLen bytes. Pass maxValueLen <= 0 to
// skip the length check.
func Validate(secrets map[string]string, maxValueLen int) error {
	var issues []string

	for k, v := range secrets {
		if !validKeyRe.MatchString(k) {
			issues = append(issues, fmt.Sprintf("invalid key %q", k))
		}
		if maxValueLen > 0 && len(v) > maxValueLen {
			issues = append(issues, fmt.Sprintf("value for key %q exceeds %d bytes", k, maxValueLen))
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}
