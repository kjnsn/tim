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
	"slices"
	"strings"

	"golang.org/x/mod/semver"
)

var ErrNoVersions = errors.New("no versions available")

// Version represents a version of a plugin, which
// can either be a git based branch & commit, or a semver
// version.
//
// HasUpgrade returns true and a new Version when there is
// an upgrade to a new version available.
//
// Check checks if there is an upgrade available.
//
// Upgrade checks out and switches to this version.
type Version interface {
	HasUpgrade() (bool, Version)

	Upgrade(pluginDir string) error

	Check(pluginDir string) error

	String() string

	GitRef() string
}

func VersionFromSpec(spec string) Version {
	if semver.IsValid(spec) {
		return &SemanticVersion{
			currentVersion: spec,
		}
	}

	if semver.IsValid("v" + spec) {
		return &SemanticVersion{
			currentVersion: "v" + spec,
		}
	}

	return &GitVersion{
		branch: spec,
	}
}

// Finds the best version of the plugin at the given pluginDir,
// preferencing semver over git.
func FindBestVersion(pluginDir string) (Version, error) {
	_, err := RunGitCommand(pluginDir, "fetch", "-t")
	if err != nil {
		return nil, err
	}

	versions, err := RunGitCommand(pluginDir, "tag", "--list", "v*")
	if err != nil {
		return nil, err
	}

	highestSemver := maxVersion(strings.Split(versions, "\n"))
	if highestSemver != "" {
		return &SemanticVersion{
			currentVersion: highestSemver,
		}, nil
	}

	branch, err := DefaultBranch(pluginDir)
	if err != nil {
		return nil, err
	}
	currentHash, err := GetRef(pluginDir, "--short", branch)
	if err != nil {
		return nil, err
	}
	return &GitVersion{
		currentHash: currentHash,
		branch:      branch,
	}, nil
}

type SemanticVersion struct {
	currentVersion string
	latestVersion  string
}

func (sv *SemanticVersion) HasUpgrade() (bool, Version) {
	if sv.latestVersion != "" && semver.Compare(sv.latestVersion, sv.currentVersion) == 1 {
		return true, &SemanticVersion{
			currentVersion: sv.latestVersion,
		}
	}
	return false, nil
}

// Checks to see if there is an upgrade, returning ErrNoVersions if no
// semantic versions are available.
func (sv *SemanticVersion) Check(pluginDir string) error {
	_, err := RunGitCommand(pluginDir, "fetch", "-t")
	if err != nil {
		return err
	}

	versions, err := RunGitCommand(pluginDir, "tag", "--list", "v*")
	if err != nil {
		return err
	}

	latest := maxVersion(strings.Split(versions, "\n"))
	if latest == "" {
		return ErrNoVersions
	}
	sv.latestVersion = latest
	if semver.Compare(latest, sv.currentVersion) == 1 {
		sv.latestVersion = latest
	}
	return nil
}

func (sv *SemanticVersion) Upgrade(pluginDir string) error {
	_, err := RunGitCommand(pluginDir, "checkout", "-f", sv.GitRef())
	return err
}

// Finds the maximum semver in the given slice of versions.
// Returns an empty string if no valid versions are present in the slice.
func maxVersion(versions []string) string {
	if len(versions) == 0 {
		return ""
	}

	return slices.MaxFunc(versions, func(a, b string) int {
		return semver.Compare(a, b)
	})
}

func (sv *SemanticVersion) String() string {
	return sv.currentVersion
}

func (sv *SemanticVersion) GitRef() string {
	return sv.currentVersion
}

type GitVersion struct {
	currentHash string
	branch      string

	// The latest hash available, empty if currently checked out
	latestHash string
}

func (gv *GitVersion) HasUpgrade() (bool, Version) {
	if gv.currentHash != "" && gv.latestHash != "" && gv.currentHash != gv.latestHash {
		return true, &GitVersion{
			currentHash: gv.latestHash,
			branch:      gv.branch,
		}
	}
	return false, nil
}

func (gv *GitVersion) Check(pluginDir string) error {
	_, err := RunGitCommand(pluginDir, "fetch", "-t")
	if err != nil {
		return err
	}

	gv.latestHash, err = GetRef(pluginDir, "--verify", "@{u}")
	return err
}

func (sv *GitVersion) Upgrade(pluginDir string) error {
	if _, err := RunGitCommand(pluginDir, "checkout", "-f", sv.GitRef()); err != nil {
		return err
	}

	_, err := RunGitCommand(pluginDir, "branch", "-f", sv.branch, "origin/"+sv.branch)
	return err
}

func (gv *GitVersion) String() string {
	return fmt.Sprintf("%s@%10s", gv.branch, gv.currentHash)
}

func (gv *GitVersion) GitRef() string {
	return gv.branch
}
