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
	"sync"

	"github.com/kjnsn/tim/lib"
	"github.com/kjnsn/tim/lib/message"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [plugin]",
	Short: "Upgrades plugins",
	Long: `Upgrades plugins.

To upgrade all plugins run "upgrade".

To check if any updates are available without modifying any versions,
pass the "--check" flag.
	
Either a single plugin can be specified, or all plugins
will be affected.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := ""
		if len(args) > 0 {
			pluginName = strings.ToLower(strings.TrimSpace(args[0]))
		}
		upgradeCommand(strings.ToLower(strings.TrimSpace(pluginName)))
	},
}

var (
	uCheckFlag bool
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().BoolVar(&uCheckFlag, "check", false, "Check if any upgrades are available without upgrading anything.")
}

func upgradeCommand(pluginName string) {
	lockFile, err := lib.GetLockfile(cfgFile)
	if err != nil {
		message.Error(err.Error())
	}
	defer lockFile.Close()

	if pluginName != "" {
		plugin := lockFile.GetPlugin(pluginName)
		if plugin == nil {
			message.Error("Plugin %s not found", pluginName)
		}
		upgradePlugin(plugin)
		if !uCheckFlag {
			lockFile.PluginSpecs[plugin.Name] = plugin.Version.GitRef()
		}
	} else {
		var wg sync.WaitGroup
		lockSync := new(sync.Mutex)

		for _, plugin := range lockFile.Plugins() {
			wg.Add(1)
			go func(plugin lib.Plugin) {
				defer wg.Done()

				upgradePlugin(&plugin)
				if !uCheckFlag {
					lockSync.Lock()
					lockFile.PluginSpecs[plugin.Name] = plugin.Version.GitRef()
					lockSync.Unlock()
				}
			}(plugin)
		}

		wg.Wait()
	}

	if !uCheckFlag {
		if err := lockFile.Save(); err != nil {
			message.Error(err.Error())
		}
	}
}

func upgradePlugin(plugin *lib.Plugin) error {
	newVersion, err := newVersion(plugin)
	if err != nil {
		return err
	}

	if newVersion == nil {
		message.Info("Plugin %s up-to-date", plugin.Name)
		return nil
	}

	oldVersion := plugin.Version.String()
	message.Info("Plugin %s has upgrade available: %s -> %s", plugin.Name, oldVersion, newVersion)

	if uCheckFlag {
		return nil
	}

	pluginDir, err := plugin.Dir()
	if err != nil {
		return err
	}
	err = newVersion.Upgrade(pluginDir)
	if err != nil {
		return err
	}

	plugin.Version = newVersion
	message.Info("Plugin %s upgraded from %s to %s", plugin.Name, oldVersion, newVersion)

	return nil
}

// Returns the version to upgrade to. Will be non-empty
// if an upgrade should occur.
func newVersion(plugin *lib.Plugin) (lib.Version, error) {
	pluginDir, err := plugin.Dir()
	if err != nil {
		return nil, err
	}

	message.Debug("Checking plugin %s for a new version", plugin.Name)

	if err := plugin.Version.Check(pluginDir); err != nil {
		return nil, err
	}

	if ok, newVersion := plugin.Version.HasUpgrade(); ok {
		return newVersion, nil
	}
	return nil, nil
}
