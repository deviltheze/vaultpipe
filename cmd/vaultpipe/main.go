package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpipe/internal/config"
	"github.com/vaultpipe/internal/dotenv"
	"github.com/vaultpipe/internal/sync"
	"github.com/vaultpipe/internal/vault"
)

var cfgFile string

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "vaultpipe",
	Short: "Sync secrets from HashiCorp Vault into local .env files",
	RunE:  runSync,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "vaultpipe.yaml", "config file path")
}

func runSync(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	writer, err := dotenv.NewWriter(cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("creating writer: %w", err)
	}

	syncer := sync.New(client, writer)
	if err := syncer.Run(cmd.Context(), cfg); err != nil {
		return fmt.Errorf("syncing secrets: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "secrets written to %s\n", cfg.OutputFile)
	return nil
}
