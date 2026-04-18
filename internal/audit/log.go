// Package audit provides structured audit logging for secret sync operations.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log event.
type Entry struct {
	Timestamp  time.Time `json:"timestamp"`
	Event      string    `json:"event"`
	SecretPath string    `json:"secret_path,omitempty"`
	OutputFile string    `json:"output_file,omitempty"`
	KeyCount   int       `json:"key_count,omitempty"`
	Error      string    `json:"error,omitempty"`
}

// Logger writes audit entries as newline-delimited JSON.
type Logger struct {
	w io.Writer
}

// NewLogger creates a Logger writing to the given path.
// Pass an empty path to write to stderr.
func NewLogger(path string) (*Logger, error) {
	if path == "" {
		return &Logger{w: os.Stderr}, nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &Logger{w: f}, nil
}

// Log writes an audit entry.
func (l *Logger) Log(e Entry) error {
	e.Timestamp = time.Now().UTC()
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}

// LogSync is a convenience helper for successful sync events.
func (l *Logger) LogSync(secretPath, outputFile string, keyCount int) error {
	return l.Log(Entry{
		Event:      "sync_success",
		SecretPath: secretPath,
		OutputFile: outputFile,
		KeyCount:   keyCount,
	})
}

// LogError is a convenience helper for error events.
func (l *Logger) LogError(secretPath string, err error) error {
	return l.Log(Entry{
		Event:      "sync_error",
		SecretPath: secretPath,
		Error:      err.Error(),
	})
}
