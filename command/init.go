package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gael-herrera/sorspec/requirement"
	"github.com/spf13/cobra"
)

type ConfigModeAction func(dir string) (*os.File, error)

type ConfigMode struct {
	fallback string
	options  map[string]ConfigModeAction
}

type Config struct {
	mode ConfigMode
}

/*
Base scope name for project configuration. Could result in a file or a directory
depending on the config mode.
*/
const base string = "sorspec"

var config Config = Config{
	mode: ConfigMode{
		fallback: "file",
		options: map[string]ConfigModeAction{
			"file": func(dir string) (*os.File, error) {
				return os.Create(filepath.Join(dir, base+".yaml"))
			},
			"dir": func(dir string) (*os.File, error) {
				os.MkdirAll(filepath.Join(dir, base, "requirement"), os.ModeDir)
				return os.Create(filepath.Join(dir, base, "app.yaml"))
			},
		},
	},
}

var initialize = &cobra.Command{
	Use:   "init <path> <...requirements>",
	Short: "Initializes a project configuration inside the given directory",
	Long:  "Initializes a project configuration inside the given directory.",
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		mode, _ := cmd.Flags().GetString("mode")

		os.MkdirAll(dir, os.ModeDir)

		configFile, _ := config.mode.options[mode](dir)

		configFileEnd, err := configFile.Write([]byte("app:\n\tname: " + filepath.Base(dir)))

		if err != nil {
			fmt.Print(fmt.Errorf("could not write on %s", filepath.Join(dir, configFile.Name())))
			return
		}

		for i := 1; i < len(args); i++ {
			req := requirement.Options[args[i]]
			/* yaml string containing the configuration for the requirement */
			reqConfig, err := req.Config.GetDefault()

			if err != nil {
				fmt.Println(fmt.Errorf("could not get default config for %s", req.Id))
				fmt.Println(err)
				return
			}

			reqConfig = fmt.Sprintf("\n\n" + req.Id + ":\n\t" + strings.ReplaceAll(reqConfig, "\n", "\n\t"))

			configFileEnd, _ = configFile.WriteAt([]byte(reqConfig), int64(configFileEnd))
		}

		os.Create(filepath.Join(dir, ".gitignore"))
		os.Create(filepath.Join(dir, "README.md"))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}

		mode, err := cmd.Flags().GetString("mode")

		if err != nil {
			return err
		}
		if config.mode.options[mode] == nil {
			var keys []string

			for key := range config.mode.options {
				keys = append(keys, key)
			}

			return fmt.Errorf("invalid mode: %s, available modes are: %s", mode, strings.Join(keys, ", "))
		}

		if err != nil {
			return err
		}
		for i := 1; i < len(args); i++ {
			req := requirement.Options[args[i]]

			if req == nil {
				return fmt.Errorf("no requirement named: '%s'", args[i])
			}
		}

		return nil
	},
}

func init() {
	initialize.Flags().StringP("mode", "m", config.mode.fallback, "defines if sorspec config will be in a dir with multiple files or in a single file")
}
