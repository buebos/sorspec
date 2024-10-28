package command

import (
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "sorspec",
	Short: "Generate web apps based on requirement specification files",
	Long:  "Generate web apps based on requirement specification files",
}

func Execute() {
	err := root.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	root.AddCommand(generate)
	root.AddCommand(initialize)
}
