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

	"github.com/kjnsn/tim/lib/message"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tim",
	Short: "Tmux plugIn Manager",
	Long: `A humble plugin manager for tmux, batteries included.

Tim manages plugins for tmux and optionaly ensures that the tmux
configuration is setup with opinionated defaults.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		message.DebugEnabled = enableVerbose
	},
}

var cfgFile string
var enableVerbose bool
var TimVersion string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ver string) {
	TimVersion = ver
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.config/tim/tim.json)")
	rootCmd.PersistentFlags().BoolVarP(&enableVerbose, "verbose", "v", false, "print verbose information")
}
