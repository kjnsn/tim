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
package lib

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"
)

var ErrPluginNotInstalled = errors.New("Plugin not installed")

// Gets the tim directory, creating it if it does not already exist.
// The tim directory is inside xdg-config-home, usually "~/.config".
// Directories ~/.config/tim and ~/.config/tim/plugins are created.
func GetTimDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	// Explicitly override the path containing `Application Support` on macos.
	// ~/.config/... is (subjectively) more conventional.
	if strings.Contains(configDir, "Application Support") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = path.Join(homeDir, ".config/")
	}

	timDir := path.Join(configDir, "/tim")

	pluginDir := path.Join(timDir, "/plugins")
	err = os.MkdirAll(pluginDir, 0750)
	if err != nil {
		return "", err
	}

	return timDir, nil
}

// Returns the plugins install directory. This is `GetTimDir() + "/plugins"`
func GetPluginsDir() (string, error) {
	timDir, err := GetTimDir()
	if err != nil {
		return "", err
	}

	return path.Join(timDir, "/plugins"), nil
}

type Plugin struct {
	// Name of the plugin in the form <username>/<repo>
	Name string

	// Semantic version of the plugin as currently installed.
	Version Version
}

// Loads the plugin by running all of it's scripts.
func (p *Plugin) Load() error {
	pluginDir, err := p.Dir()
	if err != nil {
		return err
	}

	dir := os.DirFS(pluginDir)
	entries, err := fs.ReadDir(dir, ".")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".tmux") && entry.Type().IsRegular() {
			cmd := exec.Command(path.Join(pluginDir, entry.Name()))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Returns the absolute path to this plugin's directory.
func (p *Plugin) Dir() (string, error) {
	pluginsDir, err := GetPluginsDir()
	if err != nil {
		return "", err
	}

	return path.Join(pluginsDir, p.Name), nil
}

// Checks that the given plugin is installed. Returns a nil error if successful.
func (p *Plugin) CheckInstalled() error {
	pluginDir, err := p.Dir()
	if err != nil {
		return err
	}
	fsInfo, err := os.Stat(pluginDir)
	if os.IsNotExist(err) || !fsInfo.IsDir() {
		return ErrPluginNotInstalled
	}
	if err != nil {
		return err
	}

	return nil
}

// Installs the given plugin with git, overwriting any existing configuration.
// Uses the given version spec to install at the provided version.
func (p *Plugin) Install(versionSpec string) error {
	pluginDir, err := p.Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(pluginDir, 0750); err != nil {
		return err
	}

	if _, err := RunGitCommand(pluginDir, "clone", "https://github.com/"+p.Name+".git", pluginDir); err != nil {
		return err
	}

	bestVersion, err := FindBestVersion(pluginDir)
	if err != nil {
		return err
	}
	p.Version = bestVersion

	if p.Version != nil {
		return p.CheckoutVersion(p.Version)
	}

	return nil
}

// Checks out the given git ref, could be a branch or tag.
func (p *Plugin) CheckoutVersion(version Version) error {
	pluginDir, err := p.Dir()
	if err != nil {
		return err
	}

	_, err = RunGitCommand(pluginDir, "checkout", "-q", version.GitRef())
	return err
}

// Removes all files related to this plugin from the filesystem.
func (p *Plugin) Uninstall() error {
	pluginDir, err := p.Dir()
	if err != nil {
		return err
	}

	return os.RemoveAll(pluginDir)
}
