// Package sync provides the core synchronisation logic for vaultpipe.
//
// syncer_hook.go implements pre- and post-sync lifecycle hooks.
// A hook is an optional shell command defined in the config that fires
// before or after the Vault secret sync completes.  Every invocation —
// whether or not a command is configured — is recorded as a HookRecord
// in a JSON sidecar file adjacent to the output .env file, enabling
// full auditability of when syncs occurred and what scripts ran.
//
// Hook failures are non-fatal: they are logged as warnings and recorded
// with an error field, but the sync pipeline continues normally.
package sync
