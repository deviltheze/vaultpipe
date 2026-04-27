package dotenv

import (
	"fmt"
	"time"
)

// ChainRecord represents a single entry in a secret chain log,
// recording the source path, destination path, and when the chain was applied.
type ChainRecord struct {
	Source    string    `json:"source"`
	Dest      string    `json:"dest"`
	Keys      []string  `json:"keys"`
	AppliedAt time.Time `json:"applied_at"`
}

// ChainFile holds a list of chain records for an output file.
type ChainFile struct {
	Records []ChainRecord `json:"records"`
}

// ChainPath returns the conventional path for the chain file
// associated with the given output .env file.
func ChainPath(envPath string) string {
	return envPath + ".chain.json"
}

// WriteChainRecord appends a ChainRecord to the chain file at chainPath.
// The file is created if it does not exist.
func WriteChainRecord(chainPath string, rec ChainRecord) error {
	if rec.AppliedAt.IsZero() {
		rec.AppliedAt = time.Now().UTC()
	}

	cf, err := ReadChainFile(chainPath)
	if err != nil {
		return fmt.Errorf("chain: read existing: %w", err)
	}

	cf.Records = append(cf.Records, rec)
	return writeJSON(chainPath, cf)
}

// ReadChainFile reads the chain file at path. If the file does not exist,
// an empty ChainFile is returned without error.
func ReadChainFile(path string) (ChainFile, error) {
	var cf ChainFile
	if err := readJSONOptional(path, &cf); err != nil {
		return ChainFile{}, fmt.Errorf("chain: read: %w", err)
	}
	return cf, nil
}
