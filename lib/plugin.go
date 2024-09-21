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
	Name string `json:"name"`

	// Semantic version of the plugin as currently installed.
	// Only one of version and branch can be specified at the same time.
	Version string `json:"version"`

	// Git branch of the plugin as currently installed.
	// Only one of version and branch can be specified at the same time.
	Branch string `json:"branch"`
}

// Checks that the given plugin is installed. Returns a nil error if successful.
func (p *Plugin) CheckInstalled() error {
	pluginsDir, err := GetPluginsDir()
	if err != nil {
		return err
	}

	pluginDir := path.Join(pluginsDir, p.Name)
	fsInfo, err := os.Stat(pluginDir)
	if os.IsNotExist(err) || !fsInfo.IsDir() {
		return ErrPluginNotInstalled
	}
	if err != nil {
		return err
	}

	return nil
}

// Installs the given plugin with git.
func (p *Plugin) Install() error {
	pluginsDir, err := GetPluginsDir()
	if err != nil {
		return err
	}

	pluginDir := path.Join(pluginsDir, p.Name)
	if err := os.MkdirAll(pluginDir, 0750); err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", "https://github.com/"+p.Name+".git", pluginDir)
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (p *Plugin) AvailableVersions() ([]string, error) {
	if err := p.CheckInstalled(); err != nil {
		return []string{}, err
	}

	pluginsDir, err := GetPluginsDir()
	if err != nil {
		return []string{}, err
	}

	var out strings.Builder
	cmd := exec.Command("git", "tag", "--list", "v*")
	cmd.Dir = pluginsDir
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return []string{}, err
	}

	return strings.Split(out.String(), "\n"), nil
}

type Lockfile struct {
	file *os.File

	Plugins []Plugin `json:"plugins"`
}

// Closes all resources associated with this lock file.
func (lf *Lockfile) Close() {
	lf.file.Close()
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
	lockPath := path.Join(configDir, "/timlock.json")

	actualLockFile, err := os.OpenFile(lockPath, os.O_CREATE, 0600)
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
