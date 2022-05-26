package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	dataDir          string
	numberOfBlocks   uint64
	accountPrefix    string
	connectionString string
	appName          = "cosmos-snapshot-parser"
)

// NewRootCmd returns the root command for parser.
func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   appName,
		Short: "cosmos-snapshot-parser indexes the state of a cosmos-sdk based blockchain using a relational database",
	}

	// --blocks flag
	rootCmd.PersistentFlags().Uint64VarP(&numberOfBlocks, "blocks", "b", 10, "set the amount of blocks to keep (default=10)")

	// --account-prefix flag
	rootCmd.PersistentFlags().StringVarP(&accountPrefix, "account-prefix", "a", "cosmos", "set the account account-prefix (default=cosmos)")

	// --connection-string flag
	rootCmd.PersistentFlags().StringVarP(&connectionString, "connection-string", "c", "", "connection string of that contains all the info needed to connect to the psql db: (example: postgresql://plural:plural@localhost:5432/chain?sslmode=disable)")

	// --db-dir flag
	rootCmd.PersistentFlags().StringVarP(&dataDir, "db-dir", "d", "", "the data direction that contains 'application.db' and 'blocks.db")

	rootCmd.MarkPersistentFlagRequired("blocks")
	rootCmd.MarkPersistentFlagRequired("account-prefix")
	rootCmd.MarkPersistentFlagRequired("connection-string")
	rootCmd.MarkPersistentFlagRequired("db-dir")

	rootCmd.AddCommand(
		parseCmd(),
	)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd()
	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
