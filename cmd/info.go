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
	"slices"

	"github.com/kjnsn/tim/lib"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info [plugin]",
	Short: "Displays information about installed plugins",
	Long: `Displays information about the given installed plugin,
or without an argument shows information about all plugins.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := ""
		if len(args) > 0 {
			pluginName = args[0]
		}
		infoCommand(pluginName)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func infoCommand(pluginName string) {
	lockFile, err := lib.GetLockfile()
	if err != nil {
		log.Fatal(err)
	}
	defer lockFile.Close()

	if pluginName != "" {
		i := slices.IndexFunc(lockFile.Plugins, func(plugin lib.Plugin) bool {
			return plugin.Name == pluginName
		})
		if i == -1 {
			log.Fatalf("Plugin %s not installed\n", pluginName)
		}
		printPluginInfo(lockFile.Plugins[i])
		return
	}

	for _, plugin := range lockFile.Plugins {
		printPluginInfo(plugin)
	}
}

func printPluginInfo(plugin lib.Plugin) {
	if err := plugin.CheckInstalled(); err != nil {
		log.Fatal(err)
	}

	str := ""

	str += fmt.Sprintf("Name: %s\n", plugin.Name)
	if plugin.Branch != "" {
		str += fmt.Sprintf("Using branch: %s\n", plugin.Branch)
	}
	if plugin.Version != "" {
		str += fmt.Sprintf("Version: %s\n", plugin.Version)
	}

	fmt.Println(str)
}
