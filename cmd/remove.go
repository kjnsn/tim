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

var removeCmd = &cobra.Command{
	Use:   "remove [plugin]",
	Short: "Removes a plugin",
	Long:  `Uninstalls a plugin`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		removeCommand(strings.ToLower(strings.TrimSpace(args[0])))
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func removeCommand(pluginName string) {
	lockFile, err := lib.GetLockfile()
	if err != nil {
		message.Error(err.Error())
	}
	defer lockFile.Close()

	plugin := lockFile.GetPlugin(pluginName)
	if plugin == nil {
		message.Error("plugin %s not found", pluginName)
	}

	if err := plugin.Uninstall(); err != nil {
		message.Error(err.Error())
	}

	delete(lockFile.PluginSpecs, pluginName)

	if err := lockFile.Save(); err != nil {
		message.Error(err.Error())
	}

	message.Info("Successfully uninstalled plugin %s", pluginName)
}
