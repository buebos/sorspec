package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type ModeFileAction func(dir string) (*os.File, error)

const base string = "sorspec"
const modeDefault string = "file"

var modes map[string]ModeFileAction = map[string]ModeFileAction{
	"file": func(dir string) (*os.File, error) {
		return os.Create(filepath.Join(dir, base+".yaml"))
	},
	"dir": func(dir string) (*os.File, error) {
		os.MkdirAll(filepath.Join(dir, base, "requirement"), os.ModeDir)
		return os.Create(filepath.Join(dir, base, "app.yaml"))
	},
}

var initialize = &cobra.Command{
	Use:   "init <path>",
	Run:   run,
	Args:  args,
	Short: "Generates a project inside the given directory",
	Long:  "Generates a project inside the given directory.",
}

func init() {
	initialize.Flags().StringP("mode", "m", modeDefault, "defines if sorspec config will be in a dir with multiple files or in a single file")
}

func args(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}

	mode, err := cmd.Flags().GetString("mode")

	if err != nil {
		return err
	}

	if modes[mode] == nil {
		var keys []string

		for key := range modes {
			keys = append(keys, key)
		}

		return fmt.Errorf("invalid mode: %s, available modes are: %s", mode, strings.Join(keys, ", "))
	}

	return nil
}

func run(cmd *cobra.Command, args []string) {
	dir := args[0]
	mode, _ := cmd.Flags().GetString("mode")

	os.MkdirAll(dir, os.ModeDir)

	config, _ := modes[mode](dir)

	config.Write([]byte("app:\n\tname: " + filepath.Base(dir)))

	os.Create(filepath.Join(dir, ".gitignore"))
	os.Create(filepath.Join(dir, "README.md"))
}
