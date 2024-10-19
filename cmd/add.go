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
and the latest installed by default.

If no plugin names are given, then plugins are installed according to the
configuration file ~/.config/tim/tim.json`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			addPlugin(strings.ToLower(strings.TrimSpace(args[0])))
		} else {
			syncPlugins()
		}
	},
}

var (
	versionSpec string
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&versionSpec, "version", "",
		"Version to use. Only semver 2.0 compliant strings and branch names are supported.")
}

func syncPlugins() {
	lockFile, err := lib.GetLockfile(cfgFile)
	if err != nil {
		message.Error(err.Error())
	}
	defer lockFile.Close()

	for pluginName, spec := range lockFile.PluginSpecs {
		plugin := lib.Plugin{
			Name: pluginName,
		}

		if err := plugin.Install(spec); err != nil {
			message.Error(err.Error())
		}
		message.Info("Plugin %s successfully installed at version %s", pluginName, plugin.Version)
	}
}

func addPlugin(pluginName string) {
	lockFile, err := lib.GetLockfile(cfgFile)
	if err != nil {
		message.Error(err.Error())
	}
	defer lockFile.Close()

	// If a version has not been explicitly specified,
	// try and find the plugin in the lockfile,
	// and if it exists use that version spec.
	if versionSpec == "" {
		versionSpec = lockFile.PluginSpecs[pluginName]
	}

	plugin := lib.Plugin{
		Name: pluginName,
	}

	if err := plugin.Install(versionSpec); err != nil {
		message.Error(err.Error())
	}

	lockFile.PluginSpecs[plugin.Name] = plugin.Version.GitRef()
	if err := lockFile.Save(); err != nil {
		message.Error(err.Error())
	}

	message.Info("Plugin %s successfully installed at version %s", pluginName, plugin.Version)
}
