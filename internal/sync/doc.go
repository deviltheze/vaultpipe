// Package sync orchestrates the end-to-end pipeline of fetching secrets
// from HashiCorp Vault and persisting them as a local .env file.
//
// Typical usage:
//
//	syncer, err := sync.New(cfg)
//	if err != nil { ... }
//	result, err := syncer.Run()
package sync
