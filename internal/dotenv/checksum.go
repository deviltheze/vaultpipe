package dotenv

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Checksum computes a deterministic SHA-256 hex digest over the given secrets map.
// Keys are sorted before hashing so the result is order-independent.
func Checksum(secrets map[string]string) string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, secrets[k])
	}
	return hex.EncodeToString(h.Sum(nil))
}

// WriteChecksumFile writes the checksum of secrets to path.
func WriteChecksumFile(path string, secrets map[string]string) error {
	sum := Checksum(secrets)
	return os.WriteFile(path, []byte(sum+"\n"), 0600)
}

// VerifyChecksumFile reads the checksum at path and compares it against secrets.
// Returns true when the file matches, false when it does not or does not exist.
func VerifyChecksumFile(path string, secrets map[string]string) (bool, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	stored := strings.TrimSpace(string(data))
	return stored == Checksum(secrets), nil
}
