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
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Release struct {
	Name       string
	CommitHash string
	Number     []int
	RCnumber   int
}

var (
	// regex pattern accepts v[Num].[Num].[Num] and v[Num].[Num]
	regexPatternRelease = regexp.MustCompile(`^v(\d+)\.(\d+)(\.(\d+))?(-rc(\d+))?$`)

	// regex pattern accepts refs/remotes/origin/release-[Num].[Num].[Num]
	regexPatternReleaseBranch = regexp.MustCompile(`^refs/remotes/origin/release-(\d+)\.(\d+)$`)
)

func GetPlannerVersionsForRelease(release *Release) []macrobench.PlannerVersion {
	versions := []macrobench.PlannerVersion{macrobench.V3Planner}
	if release.Number[0] >= 10 {
		versions = append(versions, macrobench.Gen4FallbackPlanner)
	}
	return versions
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
			newRelease.Number = append(newRelease.Number, num)
			num, err = strconv.Atoi(isMatched[2])
			if err != nil {
				return nil, err
			}
			newRelease.Number = append(newRelease.Number, num)
			if isMatched[4] != "" {
				num, err := strconv.Atoi(isMatched[4])
				if err != nil {
					return nil, err
				}
				newRelease.Number = append(newRelease.Number, num)
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
	for _, release := range allReleases {
		if release.Number[0] > 6 {
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
			newRelease.Number = append(newRelease.Number, num)
			num, err = strconv.Atoi(isMatched[2])
			if err != nil {
				return nil, err
			}
			newRelease.Number = append(newRelease.Number, num)
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
	for _, release := range res {
		if release.Number[0] > 6 {
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
func GetLastPatchReleaseAndCommitHash(repoDir string, releaseNumber []int) (*Release, error) {
	major := releaseNumber[0]
	minor := releaseNumber[1]
	res, err := GetAllVitessReleaseCommitHash(repoDir)
	if err != nil {
		return nil, err
	}
	for _, release := range res {
		if release.Number[0] == major && release.Number[1] == minor {
			return release, nil
		}
	}
	return nil, fmt.Errorf("could not find the latest patch release for %d.%d", major, minor)
}

// compareReleaseNumbers compares the two release numbers provided as input
// the result is as follows -
// 0, if release1 == release2
// 1, if release1 > release2
// -1, if release1 < release2
func compareReleaseNumbers(release1, release2 *Release) int {
	index := 0
	for index < len(release1.Number) && index < len(release2.Number) {
		if release1.Number[index] > release2.Number[index] {
			return 1
		}
		if release1.Number[index] < release2.Number[index] {
			return -1
		}
		index++
	}
	if len(release1.Number) > len(release2.Number) {
		return -1
	}
	if len(release1.Number) < len(release2.Number) {
		return 1
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

// GetCommitHash gets the commit hash of the current branch
func GetCommitHash(repoDir string) (hash string, err error) {
	out, err := ExecCmd(repoDir, "git", "log", "-1", "--format=%H")
	// Trimspace is used here to remove any whitespace characters after the hash
	return strings.TrimSpace(string(out)), err
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
