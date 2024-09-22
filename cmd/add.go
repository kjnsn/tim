/*
Copyright © 2024 Kaley Main <kaleymain@google.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/kjnsn/tim/lib"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [plugin]",
	Short: "Adds a plugin",
	Long: `Installs a plugin to the plugin directory.
	
Plugins are github URLs of the format <username>/<repo>.

So "add user123/my-cool-plugin" installs github.com/user123/my-cool-plugin.

The repository will be scanned for releases and tags,
and the latest installed by default.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addCommand(args[0])
	},
}

var (
	// The desired branch of the plugin.
	branchFlag  string
	versionFlag string
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&branchFlag, "branch", "", "Branch to checkout.")
	addCmd.Flags().StringVar(&versionFlag, "version", "", "Version to use. Only semver 2.0 compliant strings are supported.")
	addCmd.MarkFlagsMutuallyExclusive("branch", "version")
}

func addCommand(pluginName string) {
	lockFile, err := lib.GetLockfile()
	if err != nil {
		log.Fatal(err)
	}
	defer lockFile.Close()

	plugin := lib.Plugin{
		Name: pluginName,
	}
	if branchFlag != "" {
		plugin.Branch = branchFlag
	}
	if versionFlag != "" {
		plugin.Version = versionFlag
	}

	err = plugin.CheckInstalled()
	if err == nil {
		fmt.Printf("Plugin %s is already installed\n", pluginName)
		os.Exit(1)
	}
	if err != nil && err != lib.ErrPluginNotInstalled {
		log.Fatal(err)
	}

	if err := plugin.Install(); err != nil {
		log.Fatal(err)
	}

	lockFile.Plugins = append(lockFile.Plugins, plugin)
	if err := lockFile.Save(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Plugin %s successfully installed\n", pluginName)
}
