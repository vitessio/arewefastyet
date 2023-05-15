/*
 *
 * Copyright 2021 The Vitess Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * /
 */

package git

import (
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/vitessio/arewefastyet/go/tools/macrobench"
)

type (
	Release struct {
		Name       string
		CommitHash string
		Version    Version
		RCnumber   int
	}

	Version struct {
		Major, Minor, Patch int
	}
)

var (
	// regex pattern accepts v[Num].[Num].[Num] and v[Num].[Num]
	regexPatternRelease = regexp.MustCompile(`^v(\d+)\.(\d+)(\.(\d+))?(-rc(\d+))?$`)

	// regex pattern accepts refs/remotes/origin/release-[Num].[Num].[Num]
	regexPatternReleaseBranch = regexp.MustCompile(`^refs/remotes/origin/release-(\d+)\.(\d+)$`)
)

func GetPlannerVersions() []macrobench.PlannerVersion {
	return []macrobench.PlannerVersion{macrobench.Gen4Planner}
}

// GetAllVitessReleaseCommitHash gets all the vitess releases and the commit hashes given the directory of the clone of vitess
func GetAllVitessReleaseCommitHash(repoDir string) ([]*Release, error) {
	out, err := ExecCmd(repoDir, "git", "show-ref", "--tags", "-d")
	if err != nil {
		return nil, err
	}
	releases := strings.Split(string(out), "\n")
	var res []*Release

	// prevMatched keeps track whether the last tag matched the regular expression or not
	prevMatched := false

	for _, release := range releases {
		// if the length of the line is less than 55 then it cannot have a relese tag since
		// 40 is commit hash length + 1 space + 11 for refs/tags/v + atleast 3 for num.num
		if len(release) < 55 {
			continue
		}
		commitHash := release[0:40]
		tag := release[51:]

		// tags ending with `^{}` show dereference pointers for the given tag, and these commit hashes must be used instead of the original
		// so we check if the previous tag matched the regex and if the current tag has these 3 characters at the end, then we replace the
		// last hash with the current one
		// For example for the given input
		// c970e775be7ec79066aeddd307d050107e66c698 refs/tags/v9.0.1
		// 42c38e56e4ae29012a5d603d8bc8c22c35b78b52 refs/tags/v9.0.1^{}
		// output should have
		// tag = 9.0.1
		// commitHash = 42c38e56e4ae29012a5d603d8bc8c22c35b78b52
		if prevMatched && tag[len(tag)-3:] == "^{}" {
			res[len(res)-1].CommitHash = commitHash
		}

		isMatched := regexPatternRelease.FindStringSubmatch(tag)
		prevMatched = false
		if len(isMatched) > 6 {
			newRelease := &Release{
				Name:       tag[1:],
				CommitHash: commitHash,
			}
			num, err := strconv.Atoi(isMatched[1])
			if err != nil {
				return nil, err
			}
			newRelease.Version.Major = num
			num, err = strconv.Atoi(isMatched[2])
			if err != nil {
				return nil, err
			}
			newRelease.Version.Minor = num
			if isMatched[4] != "" {
				num, err := strconv.Atoi(isMatched[4])
				if err != nil {
					return nil, err
				}
				newRelease.Version.Patch = num
			}
			if isMatched[6] != "" {
				num, err := strconv.Atoi(isMatched[6])
				if err != nil {
					return nil, err
				}
				newRelease.RCnumber = num
			}
			res = append(res, newRelease)
			prevMatched = true
		}
	}
	// sort the releases in descending order
	sort.Slice(res, func(i, j int) bool {
		return compareReleaseNumbers(res[i], res[j]) == 1
	})
	return res, nil
}

// GetLatestVitessReleaseCommitHash gets the lastest major vitess releases and the commit hashes given the directory of the clone of vitess
func GetLatestVitessReleaseCommitHash(repoDir string) ([]*Release, error) {
	allReleases, err := GetAllVitessReleaseCommitHash(repoDir)
	if err != nil || len(allReleases) == 0 {
		return nil, err
	}
	var latestReleases []*Release

	// We take the 2 latest major release
	// TODO: @Florent: use the three latest major releases once v17 is out
	minimumRelease := allReleases[0].Version.Major - 1
	for _, release := range allReleases {
		if release.Version.Major >= minimumRelease {
			latestReleases = append(latestReleases, release)
		}
	}
	return latestReleases, nil
}

// GetAllVitessReleaseBranchCommitHash gets all the vitess release branches and the commit hashes given the directory of the clone of vitess
func GetAllVitessReleaseBranchCommitHash(repoDir string) ([]*Release, error) {
	out, err := ExecCmd(repoDir, "git", "branch", "-r", "--format", `"%(objectname) %(refname)"`)
	if err != nil {
		return nil, err
	}
	releases := strings.Split(string(out), "\n")
	var res []*Release

	for _, release := range releases {
		// value is possibly quoted
		if s, err := strconv.Unquote(release); err == nil {
			release = s
		}
		// if the length of the line is less than 60 then it cannot have a relese branch since
		// 40 is commit hash length + 20 for refs/origin/release + 3 for num.num
		if len(release) < 63 {
			continue
		}
		commitHash := release[0:40]
		tag := release[41:]

		isMatched := regexPatternReleaseBranch.FindStringSubmatch(tag)
		if len(isMatched) > 2 {
			newRelease := &Release{
				Name:       tag[20:] + "-branch",
				CommitHash: commitHash,
			}
			num, err := strconv.Atoi(isMatched[1])
			if err != nil {
				return nil, err
			}
			newRelease.Version.Major = num
			num, err = strconv.Atoi(isMatched[2])
			if err != nil {
				return nil, err
			}
			newRelease.Version.Minor = num
			res = append(res, newRelease)
		}
	}
	// sort the releases in descending order
	sort.Slice(res, func(i, j int) bool {
		return compareReleaseNumbers(res[i], res[j]) == 1
	})
	return res, nil
}

// GetLatestVitessReleaseBranchCommitHash gets the latest vitess release branches and the commit hashes given the directory of the clone of vitess
func GetLatestVitessReleaseBranchCommitHash(repoDir string) ([]*Release, error) {
	res, err := GetAllVitessReleaseBranchCommitHash(repoDir)
	if err != nil {
		return nil, err
	}
	var latestReleaseBranches []*Release
	// We take the 2 latest major release
	// TODO: @Florent: use the three latest major releases once v17 is out
	minimumRelease := res[0].Version.Major - 1
	for _, release := range res {
		if release.Version.Major >= minimumRelease {
			latestReleaseBranches = append(latestReleaseBranches, release)
		}
	}
	return latestReleaseBranches, nil
}

// GetLastReleaseAndCommitHash gets the last release number along with the commit hash given the directory of the clone of vitess
func GetLastReleaseAndCommitHash(repoDir string) (*Release, error) {
	res, err := GetAllVitessReleaseCommitHash(repoDir)
	if err != nil {
		return nil, err
	}
	return res[0], nil
}

// GetLastPatchReleaseAndCommitHash gets the last release number given the major and minor release number along with the commit hash given the directory of the clone of vitess
func GetLastPatchReleaseAndCommitHash(repoDir string, version Version) (*Release, error) {
	res, err := GetAllVitessReleaseCommitHash(repoDir)
	if err != nil {
		return nil, err
	}
	for _, release := range res {
		if release.Version.Major == version.Major && release.Version.Minor == version.Minor {
			return release, nil
		}
	}
	return nil, fmt.Errorf("could not find the latest patch release for %d.%d", version.Major, version.Minor)
}

// compareReleaseNumbers compares the two release numbers provided as input
// the result is as follows -
// 0, if release1 == release2
// 1, if release1 > release2
// -1, if release1 < release2
func compareReleaseNumbers(release1, release2 *Release) int {
	r := CompareVersionNumbers(release1.Version, release2.Version)
	if r != 0 {
		return r
	}
	if release1.RCnumber == 0 && release2.RCnumber == 0 {
		return 0
	}
	if release1.RCnumber == 0 {
		return 1
	}
	if release2.RCnumber == 0 {
		return -1
	}
	if release1.RCnumber > release2.RCnumber {
		return 1
	}
	if release1.RCnumber < release2.RCnumber {
		return -1
	}
	return 0
}

// CompareVersionNumbers compares the two version numbers provided as input
// the result is as follows -
// 0, if version1 == version2
// 1, if version1 > version2
// -1, if version1 < version2
func CompareVersionNumbers(version1, version2 Version) int {
	if version1.Major > version2.Major {
		return 1
	}
	if version1.Major < version2.Major {
		return -1
	}
	if version1.Minor > version2.Minor {
		return 1
	}
	if version1.Minor < version2.Minor {
		return -1
	}
	if version1.Patch > version2.Patch {
		return 1
	}
	if version1.Patch < version2.Patch {
		return -1
	}
	return 0
}

// GetCommitHash gets the commit hash of the current branch
func GetCommitHash(repoDir string) (hash string, err error) {
	out, err := ExecCmd(repoDir, "git", "log", "-1", "--format=%H")
	// Trimspace is used here to remove any whitespace characters after the hash
	return strings.TrimSpace(string(out)), err
}

// GetBranchesForCommit gets the branches where the commit belong
func GetBranchesForCommit(repoDir, sha string) (branches []string, err error) {
	out, err := ExecCmd(repoDir, "git", "branch", "-a", "--contains", sha)
	branches = strings.Split(string(out), "\n")
	return
}

// ShortenSHA will return the first 7 characters of a SHA.
// If the given SHA is too short, it will be returned untouched.
func ShortenSHA(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}
	return sha
}

// ExecCmd is used to execute a git command in the given directory
func ExecCmd(dir string, name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		execErr, ok := err.(*exec.ExitError)
		if ok {
			return nil, fmt.Errorf("%s:\nstderr: %s\nstdout: %s", err.Error(), execErr.Stderr, out)
		}
		if strings.Contains(err.Error(), " executable file not found in") {
			return nil, fmt.Errorf("the command `git` seems to be missing. Please install it first")
		}
		return nil, err
	}
	return out, nil
}

func GetVersionForCommitSHA(repoDir, sha string) (Version, error) {
	branches, err := GetBranchesForCommit(repoDir, sha)
	if err != nil {
		return Version{}, err
	}
	matchRelease := regexp.MustCompile(`release-([0-9]+).0`)
	for _, branch := range branches {
		if strings.Contains(branch, "origin/main") {
			lastRelease, err := GetLastReleaseAndCommitHash(repoDir)
			if err != nil {
				return Version{}, err
			}
			version := Version{
				Major: lastRelease.Version.Major + 1,
			}
			return version, nil
		}
		matches := matchRelease.FindStringSubmatch(branch)
		if len(matches) == 2 {
			majorV, err := strconv.Atoi(matches[1])
			if err != nil {
				return Version{}, err
			}
			lastPatch, err := GetLastPatchReleaseAndCommitHash(repoDir, Version{Major: majorV})
			if err != nil {
				return Version{}, err
			}
			version := lastPatch.Version
			version.Patch += 1
			return version, nil
		}
	}
	return Version{}, fmt.Errorf("release not found for %s", sha)
}
