package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	Layer "github.com/gael-herrera/sorspec/layer"
	Requirement "github.com/gael-herrera/sorspec/requirement"
	"github.com/spf13/cobra"
)

/*
Base scope name for project configuration. Could result in a file or a directory
depending on the config mode.
*/
const base = "sorspec"

var initialize = &cobra.Command{
	Use:   "init <path> <...requirements>",
	Short: "Initializes a project configuration inside the given directory",
	Long:  "Initializes a project configuration inside the given directory.",
	Run: func(cmd *cobra.Command, args []string) {
		/*
			This will get the absolute path for the relative path passed in. It
			is used to find out the name of the parent directory which will give
			a name to the app.
		*/
		dir, err := filepath.Abs(args[0])

		if err != nil {
			fmt.Print(err)
			return
		}

		/* Creates the dirs all the way to the target dir */
		err = os.MkdirAll(dir, os.ModePerm)

		if err != nil {
			fmt.Print(err)
			return
		}

		configFile, err := os.Create(filepath.Join(dir, base+".yaml"))

		if err != nil {
			fmt.Print(err)
			return
		}

		/*
			Will be the last dir's name in the path chain. If the directory passed
			in is '.' (the current dir), it will assign the name to the current
			parent directory.
		*/
		appname := filepath.Base(dir)
		tab := func(count int) string {
			result := ""

			for range count {
				result += "  "
			}

			return result
		}

		/*
			This will result in something like:

			app:
				name: <app-name>
		*/
		configFileCursor, err := configFile.Write([]byte(fmt.Sprintf("app:\n%sname: %s\n", tab(1), appname)))

		if err != nil {
			fmt.Print(fmt.Errorf("could not write on %s", dir))
			fmt.Print(err)
			return
		}

		/*
			Layers are the different components of the system and. Check out layered
			architecture concept for more context. The program handles layers
			not estrictly in this manner (because it includes platform concepts like
			ios and android). But the main thing is that each layer can have a core
			engine which can be a different variety of tools: languages, frameworks,
			dbms, sdks, etc.
		*/
		for _, layer := range Layer.Layers {
			core, _ := cmd.Flags().GetString(layer)

			/* The user did not assign the flag's value */
			if !cmd.Flag(layer).Changed {
				continue
			}

			/*
				This will result in something like:

				database:
					core: postgres
				server:
					core: go
			*/
			field := fmt.Sprintf("%s%s:\n%score: %s\n", tab(1), layer, tab(2), core)

			/* This should keep the cursor at the end of the file */
			count, _ := configFile.WriteAt([]byte(field), int64(configFileCursor))
			configFileCursor += count
		}

		/*
			Requirements are concepts or cross-cutting concerns of the app;
			authentication, authorization, etc. They are passed in after
			the path for the command (that's why i starts at 1).
		*/
		for i := 1; i < len(args); i++ {
			/*
				The requirements are registered in a hashmap at that package
				with the same keys that should be passed in to this command.
			*/
			requirement := Requirement.Requirements[args[i]]
			/* yaml string containing the configuration for the requirement */
			reqConfig, err := requirement.Config.GetDefault()

			if err != nil {
				fmt.Println(fmt.Errorf("could not get default config for %s", requirement.Id))
				fmt.Println(err)
				return
			}

			/*
				The replace all call is to offset one tab all the requirement
				config since it will be inside the tag of that requirement's
				name. This will result in something like:

				<name-of-req>:
					...<all-config-of-req>
			*/
			reqConfig = fmt.Sprintf("\n%s:\n%s%s\n", requirement.Id, tab(1), strings.ReplaceAll(reqConfig, "\n", "\n"+tab(1)))

			count, _ := configFile.WriteAt([]byte(reqConfig), int64(configFileCursor))
			configFileCursor += count
		}

		os.Create(filepath.Join(dir, ".gitignore"))
		os.Create(filepath.Join(dir, "README.md"))
	},
	/*
		Handles validation before the 'Run' function.
	*/
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}

		for i := 1; i < len(args); i++ {
			requirement := Requirement.Requirements[args[i]]

			if requirement == nil {
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
	initialize.Flags().StringP("database", "d", "", "defines the system to use for the database")
	initialize.Flags().StringP("server", "s", "", "defines the core in which to build the server")
	initialize.Flags().StringP("browser", "b", "", "defines the framework in which to build the frontend")
	initialize.Flags().StringP("android", "a", "", "defines the core to use for building the android app")
	initialize.Flags().StringP("ios", "i", "", "defines the core to use for building the ios app")
}
