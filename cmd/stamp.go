package cmd

import (
	"github.com/puppetlabs/cat-team-github-metrics/internal/stamp"
	"github.com/spf13/cobra"
)

var stampCmd = &cobra.Command{
	Use:   "stamp",
	Short: "Updates the watermark table with the current timestamp",
	Long:  "Updates the watermark table with the current timestamp",
	RunE: func(cmd *cobra.Command, args []string) error {
		return stamp.Run()
	},
}
