package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{
		Use:          "scripts",
		Short:        "Build panascais/node images and update their configuration",
		SilenceUsage: true,
	}
	command.AddCommand(buildCommand(), planCommand(), publishCommand(), updateCommand())
	command.CompletionOptions.DisableDefaultCmd = true

	if command.Execute() != nil {
		os.Exit(1)
	}
}
