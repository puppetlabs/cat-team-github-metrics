package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var stampCmd = &cobra.Command{
	Use:   "stamp",
	Short: "Updates the watermark table with the current timestamp",
	Long:  "Updates the watermark table with the current timestamp",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(time.Now().Format(time.RFC3339))
		return nil
	},
}
