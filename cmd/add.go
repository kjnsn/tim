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
	"os"
	"strings"

	"github.com/kjnsn/tim/lib"
	"github.com/kjnsn/tim/lib/message"
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
		addCommand(strings.ToLower(strings.TrimSpace(args[0])))
	},
}

var (
	versionSpecFlag string
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&versionSpecFlag, "version", "",
		"Version to use. Only semver 2.0 compliant strings and branch names are supported.")
}

func addCommand(pluginName string) {
	lockFile, err := lib.GetLockfile(cfgFile)
	if err != nil {
		message.Error(err.Error())
	}
	defer lockFile.Close()

	plugin := lib.Plugin{
		Name: pluginName,
	}

	err = plugin.CheckInstalled()
	if err == nil {
		message.Error("Plugin %s is already installed\n", pluginName)
		os.Exit(1)
	}
	if err != nil && err != lib.ErrPluginNotInstalled {
		message.Error(err.Error())
	}

	if err := plugin.Install(versionSpecFlag); err != nil {
		message.Error(err.Error())
	}

	lockFile.PluginSpecs[plugin.Name] = plugin.Version.GitRef()
	if err := lockFile.Save(); err != nil {
		message.Error(err.Error())
	}

	message.Info("Plugin %s successfully installed\n", pluginName)
}
