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
	"io"
	"log"
	"os"

	"github.com/kjnsn/tim/lib"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		initCommand(cmd, args)
	},
}

var cfgFlag string

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&cfgFlag, "config", "", "Path to the tmux config file.")
}

func initCommand(cmd *cobra.Command, args []string) {
	var tmuxConfigPath string
	var err error

	// Find the path to the config file.
	if cfgFlag != "" {
		tmuxConfigPath = cfgFlag
	} else {
		// Try and find the file automatically.
		tmuxConfigPath, err = lib.GetTmuxConfigPath()
		if err != nil {
			log.Fatal(err)
		}
	}

	tmuxConfig, err := os.Open(tmuxConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	contents, err := io.ReadAll(tmuxConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(contents))

}
