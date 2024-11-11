package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	Layer "github.com/gael-herrera/sorspec/layer"
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
const base = "sorspec"

var config Config = Config{
	mode: ConfigMode{
		fallback: "file",
		options: map[string]ConfigModeAction{
			"file": func(dir string) (*os.File, error) {
				return os.Create(filepath.Join(dir, base+".yaml"))
			},
			"dir": func(dir string) (*os.File, error) {
				os.MkdirAll(filepath.Join(dir, base, "requirement"), os.ModePerm)
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

		os.MkdirAll(dir, os.ModePerm)

		configFile, err := config.mode.options[mode](dir)

		if err != nil {
			fmt.Print(err)
			return
		}

		configFileCursor, err := configFile.Write([]byte(fmt.Sprintf("app:\n%sname: %s\n", tab(1), filepath.Base(dir))))

		if err != nil {
			fmt.Print(fmt.Errorf("could not write on %s", filepath.Join(dir, configFile.Name())))
			fmt.Print(err)
			return
		}

		for _, layer := range Layer.Layers {
			core, _ := cmd.Flags().GetString(layer)

			if core == "" {
				continue
			}

			field := fmt.Sprintf("%s%s:\n%score: %s\n", tab(1), layer, tab(2), core)

			count, _ := configFile.WriteAt([]byte(field), int64(configFileCursor))
			configFileCursor += count
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

			reqConfig = fmt.Sprintf("\n%s:\n%s%s\n", req.Id, tab(1), strings.ReplaceAll(reqConfig, "\n", "\n"+tab(1)))

			count, _ := configFile.WriteAt([]byte(reqConfig), int64(configFileCursor))
			configFileCursor += count
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

		for i := 1; i < len(args); i++ {
			req := requirement.Options[args[i]]

			if req == nil {
				return fmt.Errorf("no requirement named: '%s'", args[i])
			}
		}

		for _, layer := range Layer.Layers {
			core, _ := cmd.Flags().GetString(layer)

			/* The flag has default value and was not set by the user */
			if !cmd.Flag(layer).Changed {
				continue
			}

			if Layer.GetCoreLatest(layer, core) == "" {
				return fmt.Errorf("no core named: '%s' for layer %s", core, layer)
			}
		}

		return nil
	},
}

func init() {
	initialize.Flags().StringP("mode", "m", config.mode.fallback, "defines if sorspec config will be in a dir with multiple files or in a single file")

	initialize.Flags().StringP("database", "d", "", "defines the system to use for the database")
	initialize.Flags().StringP("server", "s", "", "defines the core in which to build the server")
	initialize.Flags().StringP("browser", "b", "", "defines the framework in which to build the frontend")
	initialize.Flags().StringP("android", "a", "", "defines the core to use for building the ios app")
	initialize.Flags().StringP("ios", "i", "", "defines the core to use for building the ios app")
}

func tab(count int) string {
	result := ""

	for range count {
		result += "  "
	}

	return result
}
