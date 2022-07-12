// Package cmd holds all of the entry points for the cobra app.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var errSilent = errors.New("ErrSilent")

var rootCmd = &cobra.Command{
	Use:           "collector",
	Short:         "Collects GitHub metrics",
	Long:          "Collects GitHub metrics and exports them to BigQuery",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run:           nil,
}

func init() {
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(stampCmd)
}

// Execute is called from main and is responsible for processing
// requests to the application and handling exit codes appropriately
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		if err != errSilent {
			fmt.Fprintln(os.Stderr, fmt.Errorf("❌ %s", err))
		}
		return 1
	}
	return 0
}
