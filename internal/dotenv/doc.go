// Package dotenv provides utilities for reading, writing, transforming,
// validating and managing .env files used by vaultpipe.
//
// Subfeatures include:
//   - Reading and writing key=value pairs
//   - Merging, diffing, and comparing secret maps
//   - Backup, rotation, rollback, and pruning of .env files
//   - Encryption and decryption of secret values at rest
//   - Filtering, masking, sanitising, and redacting secrets
//   - Checksum and drift detection
//   - TTL, expiry, lock, pin, alias, namespace, scope, and chain management
//   - Template generation and lineage tracking
//   - Snapshot and annotation support
//   - Policy enforcement for secret compliance
package dotenv
