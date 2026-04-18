// Package dotenv provides utilities for reading, writing, and merging
// .env files used by vaultpipe to persist secrets locally.
//
// Read parses an existing .env file into a key/value map.
// Writer writes a sorted, safely-quoted .env file from a map.
// Merge combines two maps according to a configurable MergeStrategy,
// returning the merged result and a summary of changes.
package dotenv
