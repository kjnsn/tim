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
	"github.com/kjnsn/tim/lib"
	"github.com/kjnsn/tim/lib/message"
	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Loads all plugins into tmux",
	Long:  `Loads all plugins, running their scripts.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		loadCommand()
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
}

func loadCommand() {
	lockFile, err := lib.GetLockfile()
	if err != nil {
		message.Error(err.Error())
	}
	defer lockFile.Close()

	for _, plugin := range lockFile.Plugins() {
		if err := plugin.Load(); err != nil {
			message.Error(err.Error())
		}
		message.Info("loaded plugin %s", plugin.Name)
	}
}
