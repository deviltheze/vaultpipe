package dotenv

import "strings"

// SensitivePatterns lists substrings that indicate a key holds sensitive data.
var SensitivePatterns = []string{
	"PASSWORD", "PASSWD", "SECRET", "TOKEN", "API_KEY", "APIKEY",
	"PRIVATE_KEY", "PRIVATE", "CREDENTIAL", "AUTH",
}

// RedactedValue is the placeholder used in place of sensitive values.
const RedactedValue = "********"

// Redact returns a copy of secrets where values whose keys match any
// SensitivePattern are replaced with RedactedValue. Useful for logging
// and display purposes.
func Redact(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if isSensitive(k) {
			out[k] = RedactedValue
		} else {
			out[k] = v
		}
	}
	return out
}

// isSensitive reports whether key contains any sensitive pattern.
func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range SensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}
