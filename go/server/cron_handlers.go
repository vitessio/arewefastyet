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

package server

import (
	"strings"

	"github.com/vitessio/arewefastyet/go/exec"
	"github.com/vitessio/arewefastyet/go/tools/git"
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
)

func (s *Server) branchCronHandler() {
	// update the local clone of vitess from remote
	s.vitessPathMu.Lock()
	defer s.vitessPathMu.Unlock()
	err := s.pullLocalVitess()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	mainBranchElements, err := s.mainBranchCronHandler()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	releaseBranchElements, err := s.releaseBranchesCronHandler()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	execElements := append(mainBranchElements, releaseBranchElements...)
	for _, elem := range execElements {
		s.addToQueue(elem)
	}
}

func (s *Server) mainBranchCronHandler() ([]*executionQueueElement, error) {
	var elements []*executionQueueElement
	configs := s.getConfigFiles()
	vitessPath := s.getVitessPath()

	// getting the latest commit hash from local fork of Vitess
	ref, err := git.GetCommitHash(vitessPath)
	if err != nil {
		return nil, err
	}

	// getting the latest release from local fork of Vitess
	lastRelease, err := git.GetLastReleaseAndCommitHash(vitessPath)
	if err != nil {
		return nil, err
	}
	currVersion := git.Version{Major: lastRelease.Version.Major + 1}

	// We compare main with the previous hash of main and with the latest release
	for workload, config := range configs {
		if config.skip {
			continue
		}
		if minVersion := config.v.GetInt(keyMinimumVitessVersion); minVersion > currVersion.Major {
			continue
		}
		if workload == "micro" {
			_, previousGitRef, err := exec.GetPreviousFromSourceMicrobenchmark(s.dbClient, exec.SourceCron, ref)
			if err != nil {
				slog.Warn(err.Error())
				continue
			}
			elements = append(elements, s.createBranchElementWithComparisonOnPreviousAndRelease(config, ref, workload, previousGitRef, "", exec.SourceCron, lastRelease, currVersion)...)
		} else {
			for _, version := range git.GetPlannerVersions() {
				_, previousGitRef, err := exec.GetPreviousFromSourceMacrobenchmark(s.dbClient, exec.SourceCron, workload, string(version), ref)
				if err != nil {
					slog.Warn(err.Error())
					continue
				}
				elements = append(elements, s.createBranchElementWithComparisonOnPreviousAndRelease(config, ref, workload, previousGitRef, string(version), exec.SourceCron, lastRelease, currVersion)...)
			}
		}
	}
	return elements, nil
}

func (s *Server) releaseBranchesCronHandler() ([]*executionQueueElement, error) {
	var elements []*executionQueueElement
	configs := s.getConfigFiles()
	vitesLocalPath := s.getVitessPath()

	releases, err := git.GetLatestVitessReleaseBranchCommitHash(vitesLocalPath)
	if err != nil {
		slog.Warn(err.Error())
		return nil, err
	}

	// We compare release-branches with the previous hash of that branch and with the latest patch release of that version
	for _, release := range releases {
		ref := release.CommitHash
		source := exec.SourceReleaseBranch + release.Name
		lastPatchRelease, err := git.GetLastPatchReleaseAndCommitHash(vitesLocalPath, release.Version)
		if err != nil && !strings.Contains(err.Error(), "could not find the latest patch release") {
			slog.Warn(err.Error())
			continue
		}
		currVersion := release.Version
		if lastPatchRelease != nil {
			currVersion = lastPatchRelease.Version
			currVersion.Patch += 1
		}

		for workload, config := range configs {
			if config.skip {
				continue
			}
			if minVersion := config.v.GetInt(keyMinimumVitessVersion); minVersion > currVersion.Major {
				continue
			}

			if workload == "micro" {
				_, previousGitRef, err := exec.GetPreviousFromSourceMicrobenchmark(s.dbClient, source, ref)
				if err != nil {
					slog.Warn(err.Error())
					continue
				}

				elements = append(elements, s.createBranchElementWithComparisonOnPreviousAndRelease(config, ref, workload, previousGitRef, "", source, lastPatchRelease, currVersion)...)
			} else {
				versions := git.GetPlannerVersions()

				for _, version := range versions {
					_, previousGitRef, err := exec.GetPreviousFromSourceMacrobenchmark(s.dbClient, source, workload, string(version), ref)
					if err != nil {
						slog.Warn(err.Error())
						continue
					}

					elements = append(elements, s.createBranchElementWithComparisonOnPreviousAndRelease(config, ref, workload, previousGitRef, string(version), source, lastPatchRelease, currVersion)...)
				}
			}
		}
	}
	return elements, nil
}

func (s *Server) createBranchElementWithComparisonOnPreviousAndRelease(config benchmarkConfig, ref, workload, previousGitRef, plannerVersion, source string, lastRelease *git.Release, version git.Version) []*executionQueueElement {
	var elements []*executionQueueElement

	// creating a benchmark for the latest commit on the branch with SourceCron as a source
	// this benchmark will be compared with the previous benchmark on branch with the given source, and with the latest release
	newExecutionElement := s.createSimpleExecutionQueueElement(config, source, ref, workload, plannerVersion, false, 0, version)
	elements = append(elements, newExecutionElement)

	if previousGitRef != "" {
		// creating an execution queue element for the latest benchmark with SourceCron as source
		// this will not be executed since the benchmark already exist, we still create the element in order to compare
		previousElement := s.createSimpleExecutionQueueElement(config, source, previousGitRef, workload, plannerVersion, false, 0, version)
		previousElement.compareWith = append(previousElement.compareWith, newExecutionElement.identifier)
		newExecutionElement.compareWith = append(newExecutionElement.compareWith, previousElement.identifier)
		elements = append(elements, previousElement)
	}

	if lastRelease != nil && config.v.GetInt(keyMinimumVitessVersion) <= lastRelease.Version.Major {
		// creating an execution queue element for the latest release (comparing branch with the latest release)
		// this will probably not be executed the benchmark should already exist, we still create it to compare main once its benchmark is over
		lastReleaseElement := s.createSimpleExecutionQueueElement(config, exec.SourceTag+lastRelease.Name, lastRelease.CommitHash, workload, plannerVersion, false, 0, lastRelease.Version)
		lastReleaseElement.compareWith = append(lastReleaseElement.compareWith, newExecutionElement.identifier)
		newExecutionElement.compareWith = append(newExecutionElement.compareWith, lastReleaseElement.identifier)
		elements = append(elements, lastReleaseElement)
	}
	return elements
}

func (s *Server) pullRequestsCronHandler() {
	vitesLocalPath := s.getVitessPath()
	configs := s.getConfigFiles()
	prLabelsInfo := []struct {
		label   string
		useGen4 bool
	}{
		{label: s.prLabelTrigger, useGen4: true},
		{label: s.prLabelTriggerV3, useGen4: false},
	}

	var elements []*executionQueueElement

	for _, labelInfo := range prLabelsInfo {
		prInfos, err := git.GetPullRequestsFromGitHub([]string{labelInfo.label}, "vitessio/vitess")
		if err != nil {
			slog.Warn(err)
			continue
		}

		for _, prInfo := range prInfos {
			ref := prInfo.SHA
			previousGitRef := prInfo.Base
			pullNb := prInfo.Number
			if ref == "" || pullNb == 0 {
				continue
			}
			currVersion, err := git.GetVersionForCommitSHA(vitesLocalPath, previousGitRef)
			if err != nil {
				slog.Warn(err)
				continue
			}
			for workload, config := range configs {
				if config.skip {
					continue
				}
				if minVersion := config.v.GetInt(keyMinimumVitessVersion); minVersion > currVersion.Major {
					continue
				}

				if workload == "micro" {
					elements = append(elements, s.createPullRequestElementWithBaseComparison(config, ref, workload, previousGitRef, "", pullNb, currVersion)...)
				} else {
					versions := []macrobench.PlannerVersion{macrobench.V3Planner}
					if labelInfo.useGen4 {
						versions = []macrobench.PlannerVersion{macrobench.Gen4Planner}
					}
					for _, version := range versions {
						elements = append(elements, s.createPullRequestElementWithBaseComparison(config, ref, workload, previousGitRef, version, pullNb, currVersion)...)
					}
				}
			}
		}
	}
	for _, element := range elements {
		s.removePRFromQueue(element)
		s.addToQueue(element)
	}
}

func (s *Server) createPullRequestElementWithBaseComparison(config benchmarkConfig, ref, workload, previousGitRef string, plannerVersion macrobench.PlannerVersion, pullNb int, gitVersion git.Version) []*executionQueueElement {
	var elements []*executionQueueElement

	newExecutionElement := s.createSimpleExecutionQueueElement(config, exec.SourcePullRequest, ref, workload, string(plannerVersion), true, pullNb, gitVersion)
	newExecutionElement.identifier.PullBaseRef = previousGitRef
	elements = append(elements, newExecutionElement)

	if previousGitRef != "" {
		previousElement := s.createSimpleExecutionQueueElement(config, exec.SourcePullRequestBase, previousGitRef, workload, string(plannerVersion), false, pullNb, gitVersion)
		previousElement.compareWith = append(previousElement.compareWith, newExecutionElement.identifier)
		newExecutionElement.compareWith = append(newExecutionElement.compareWith, previousElement.identifier)
		elements = append(elements, previousElement)
	}
	return elements
}

func (s *Server) tagsCronHandler() {
	// update the local clone of vitess from remote
	s.vitessPathMu.Lock()
	defer s.vitessPathMu.Unlock()
	err := s.pullLocalVitess()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	configs := s.getConfigFiles()

	releases, err := git.GetLatestVitessReleaseCommitHash(s.getVitessPath())
	if err != nil {
		slog.Error(err)
		return
	}

	var elements []*executionQueueElement

	// We add single executions for the tags, we do not compare them against anything
	for _, release := range releases {
		source := exec.SourceTag + release.Name
		for workload, config := range configs {
			if config.skip {
				continue
			}
			if minVersion := config.v.GetInt(keyMinimumVitessVersion); minVersion > release.Version.Major {
				continue
			}
			if workload == "micro" {
				elements = append(elements, s.createSimpleExecutionQueueElement(config, source, release.CommitHash, workload, "", true, 0, release.Version))
			} else {
				versions := git.GetPlannerVersions()
				for _, version := range versions {
					elements = append(elements, s.createSimpleExecutionQueueElement(config, source, release.CommitHash, workload, string(version), true, 0, release.Version))
				}
			}
		}
	}
	for _, element := range elements {
		s.addToQueue(element)
	}
}

func (s *Server) createSimpleExecutionQueueElement(config benchmarkConfig, source, ref, workload, plannerVersion string, notify bool, pullNb int, version git.Version) *executionQueueElement {
	return &executionQueueElement{
		config:       config,
		retry:        s.cronNbRetry,
		notifyAlways: notify,
		identifier: executionIdentifier{
			GitRef:         ref,
			Source:         source,
			Workload:       workload,
			PlannerVersion: plannerVersion,
			PullNb:         pullNb,
			Version:        version,
		},
	}
}
