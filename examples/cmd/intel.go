package cmd

import (
	"github.com/spf13/cobra"
)

var intelCmd = &cobra.Command{
	Use:   "intel",
	Short: "Intel service examples",
}

func init() {
	ExamplesCmd.AddCommand(intelCmd)
}
