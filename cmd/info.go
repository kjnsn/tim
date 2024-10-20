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
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/kjnsn/tim/lib"
	"github.com/kjnsn/tim/lib/message"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info [plugin]",
	Short: "Displays information about installed plugins and tim itself",
	Long: `Displays information about the given installed plugin,
or without an argument shows information about all plugins.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := ""
		if len(args) > 0 {
			pluginName = strings.ToLower(strings.TrimSpace(args[0]))
		}
		infoCommand(pluginName)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func infoCommand(pluginName string) {
	lockFile, err := lib.GetLockfile(cfgFile)
	if err != nil {
		message.Error(err.Error())
	}
	defer lockFile.Close()

	if pluginName != "" {
		i := slices.IndexFunc(lockFile.Plugins(), func(plugin lib.Plugin) bool {
			return plugin.Name == pluginName
		})
		if i == -1 {
			message.Warning("Plugin %s not installed", pluginName)
		} else {
			printPluginInfo(lockFile.Plugins()[i])
		}

		return
	}

	// Print some generic information about tim.
	message.Info("Tim Version: %s", TimVersion)
	message.Info("Lockfile: %s", lockFile.Path())

	for _, plugin := range lockFile.Plugins() {
		printPluginInfo(plugin)
	}
}

func printPluginInfo(plugin lib.Plugin) {
	pluginDir, err := plugin.Dir()
	if err != nil {
		message.Error(err.Error())
	}

	str := ""

	name := message.Hyperlink("https://github.com/"+plugin.Name, plugin.Name)
	str += fmt.Sprintf("\nName: %s\n", name)
	if plugin.Version != nil {
		switch version := plugin.Version.(type) {
		case *lib.SemanticVersion:
			ver := message.Hyperlink("https://github.com/"+plugin.Name+"/releases/tag/"+version.GitRef(), version.String())
			str += fmt.Sprintf("Version: %s\n", ver)
		case *lib.GitVersion:
			ver := message.Hyperlink("https://github.com/"+plugin.Name+"/tree/"+version.GitRef(), version.String())
			str += fmt.Sprintf("Version: %s\n", ver)
		}
	}
	fmt.Println(str)

	err = plugin.CheckInstalled()
	if err != nil {
		if errors.Is(err, lib.ErrPluginNotInstalled) {
			message.Warning("Plugin %s is present in the config file but not installed.\n"+
				"  Run \"tim add\" to install it.", plugin.Name)
		} else {
			message.Error(err.Error())
		}
	} else {
		str += fmt.Sprintf("Installed to: %s", pluginDir)
	}
}
