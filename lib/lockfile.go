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
	"io"
	"os"
	"path"
)

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
func GetLockfile(cfgOverride string) (*Lockfile, error) {
	lockPath, err := lockfilePath(cfgOverride)
	if err != nil {
		return nil, err
	}

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

// Returns the path to the lockfile.
// Preferences, in order:
// - pathOverride
// - ~/.config/tim/tim.json
func lockfilePath(pathOverride string) (string, error) {
	// Ensure the correct directories are created.
	timDir, err := GetTimDir()
	if err != nil {
		return "", err
	}

	if pathOverride != "" {
		return pathOverride, nil
	}

	return path.Join(timDir, "tim.json"), nil
}
