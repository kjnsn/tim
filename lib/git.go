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
	"os"
	"os/exec"
	"strings"
)

// Returns the default branch of the given repo at basedir (what does the upstream default to).
func DefaultBranch(basedir string) (string, error) {
	return GetRef(basedir, "--abbrev-ref", "origin/HEAD")
}

// Runs `get rev-parse` with the given flag and pathspec.
func GetRef(basedir, flag, pathspec string) (string, error) {
	return RunGitCommand(basedir, "rev-parse", flag, pathspec)
}

// Checks out and updates `branch` from the remote.
func UpdateBranch(baseDir, branch string) error {
	_, err := RunGitCommand(baseDir, "checkout", "-f", branch)
	if err != nil {
		return err
	}

	_, err = RunGitCommand(baseDir, "pull", "--ff-only", "-q")
	return err
}

// Runs the given git command.
func RunGitCommand(basedir string, args ...string) (string, error) {
	var out strings.Builder
	cmd := exec.Command("git", args...)
	cmd.Dir = basedir
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
