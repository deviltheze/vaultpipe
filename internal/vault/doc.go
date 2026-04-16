// Package vault provides a thin wrapper around the HashiCorp Vault API client
// used by vaultpipe to authenticate and retrieve secrets.
//
// It supports both KV v1 and KV v2 secret engines and returns secrets as a
// flat map[string]string suitable for writing to .env files.
package vault
