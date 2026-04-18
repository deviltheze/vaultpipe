// Package dotenv provides utilities for reading, writing, and merging
// .env files used by vaultpipe to persist secrets locally.
//
// # Overview
//
// Read parses an existing .env file into a key/value map, handling
// quoted values, escaped characters, and inline comments.
//
// Write serializes a key/value map to a sorted, safely-quoted .env file,
// ensuring values containing spaces or special characters are wrapped
// in double quotes.
//
// Merge combines two maps (base and overlay) according to a configurable
// [MergeStrategy], returning the merged result and a [ChangeSet] that
// summarises which keys were added, updated, or left unchanged.
package dotenv
