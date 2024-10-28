package command

import (
	"github.com/spf13/cobra"
)

var generate = &cobra.Command{
	Use:   "generate <path>",
	Short: "Generates a project based on the requirements found on <path>",
	Long: "Generates a project based on the requirements found on <path>. Path must be a\n" +
		"directory which contains a config.yaml file. By default the command will search\n" +
		"on: '.sorspec' and 'sorspec'.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
