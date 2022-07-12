package cmd

import (
	"github.com/puppetlabs/cat-team-github-metrics/internal/export"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports GitHub metrics to BigQuery",
	Long:  "Exports GitHub metrics to BigQuery",
	RunE: func(cmd *cobra.Command, args []string) error {
		return export.Run()
	},
}
