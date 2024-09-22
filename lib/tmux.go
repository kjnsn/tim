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
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

var ErrNoTmuxConfig = errors.New("no tmux.conf file found")

// Ensures that tmux is installed, and returns the version as a string.
func GetTmuxVersion() (string, error) {
	cmd := exec.Command("tmux", "-V")
	var out strings.Builder
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// The command outputs `tmux <version>`. Ensure
	// both parts are correct.
	splits := strings.Split(out.String(), " ")
	if len(splits) != 2 {
		return "", fmt.Errorf("bad output of command 'tmux -V': %s", out.String())
	}
	if splits[0] != "tmux" {
		return "", fmt.Errorf("bad output of first part of command 'tmux -V': %s", out.String())
	}

	return splits[1], nil
}

// Finds the path of the current tmux configuration file.
func GetTmuxConfigPath() (string, error) {
	var configPath string

	// Check ~/.tmux.conf first.
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath = path.Join(home, ".tmux.conf")
	_, err = os.Stat(configPath)
	if err == nil {
		// Config found! Return the path.
		return configPath, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	// Check the alternative ~/.config/tmux/tmux.conf.
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	_, err = os.Stat(path.Join(configDir, "tmux/tmux.conf"))
	if err == nil {
		// Config found! Return the path.
		return configPath, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	// No more available candidates.
	return "", ErrNoTmuxConfig
}
