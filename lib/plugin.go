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
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"
)

var ErrPluginNotInstalled = errors.New("Plugin not installed")

// Gets the plugins directory, creating it if it does not already exist.
// The plugin directory is inside xdg-config-home, usually "~/.config".
// Directories ~/.config/tim and ~/.config/tim/plugins are created.
func GetPluginsDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	pluginDir := path.Join(configDir, "/tim/plugins")
	err = os.MkdirAll(pluginDir, 0750)
	if err != nil {
		return "", err
	}

	return pluginDir, nil
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

type Lockfile struct {
	file *os.File

	PluginSpecs map[string]string `json:"plugins"`
}

func (lf *Lockfile) Plugins() []Plugin {
	plugins := make([]Plugin, 0)
	for name, versionSpec := range lf.PluginSpecs {
		plugins = append(plugins, Plugin{
			Name:    name,
			Version: VersionFromSpec(versionSpec),
		})
	}
	return plugins
}

// Attempts to find a plugin with the given name. Returns nil if the given plugin cannot be found.
func (lf *Lockfile) GetPlugin(name string) *Plugin {
	for _, plugin := range lf.Plugins() {
		if plugin.Name == name {
			return &plugin
		}
	}

	return nil
}

// Closes all resources associated with this lock file.
func (lf *Lockfile) Close() {
	lf.file.Close()
}

// Writes the lock file to disk.
func (lf *Lockfile) Save() error {
	if err := lf.file.Truncate(0); err != nil {
		return err
	}
	if _, err := lf.file.Seek(0, 0); err != nil {
		return err
	}

	defer lf.file.Sync()

	encoder := json.NewEncoder(lf.file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(lf)
}

// Loads the lockfile, creating one if required.
func GetLockfile() (*Lockfile, error) {
	// Ensure the correct directories are created.
	_, err := GetPluginsDir()
	if err != nil {
		return nil, err
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	lockPath := path.Join(configDir, "/tim/timlock.json")

	actualLockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	lockFileContents, err := io.ReadAll(actualLockFile)
	if err != nil {
		return nil, err
	}

	lockFile := &Lockfile{
		file: actualLockFile,
	}

	// Only try and parse the contents if the file is non-empty.
	if len(lockFileContents) > 0 {
		if err := json.Unmarshal(lockFileContents, lockFile); err != nil {
			return nil, err
		}
	}

	return lockFile, nil
}
