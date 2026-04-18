// Package audit provides structured, append-only audit logging for
// vaultpipe sync operations. Each event is written as a newline-delimited
// JSON record containing a UTC timestamp, event type, and relevant metadata.
//
// Audit logs can be directed to a file (recommended for production) or
// to stderr for development use.
package audit
