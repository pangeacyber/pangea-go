package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var ExamplesCmd = &cobra.Command{
	Use: "examples",
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func Execute() {
	if err := ExamplesCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
