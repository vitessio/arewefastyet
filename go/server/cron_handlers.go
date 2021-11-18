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

	// We compare main with the previous hash of main and with the latest release
	for configType, configFile := range configs {
		if configType == "micro" {
			_, previousGitRef, err := exec.GetPreviousFromSourceMicrobenchmark(s.dbClient, exec.SourceCron, ref)
			if err != nil {
				slog.Warn(err.Error())
				continue
			}
			elements = append(elements, s.createBranchElementWithComparisonOnPreviousAndRelease(configFile, ref, configType, previousGitRef, "", exec.SourceCron, lastRelease)...)
		} else {
			for _, version := range macrobench.PlannerVersions {
				_, previousGitRef, err := exec.GetPreviousFromSourceMacrobenchmark(s.dbClient, exec.SourceCron, configType, string(version), ref)
				if err != nil {
					slog.Warn(err.Error())
					continue
				}
				elements = append(elements, s.createBranchElementWithComparisonOnPreviousAndRelease(configFile, ref, configType, previousGitRef, string(version), exec.SourceCron, lastRelease)...)
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
		lastPatchRelease, err := git.GetLastPatchReleaseAndCommitHash(vitesLocalPath, release.Number)
		if err != nil {
			slog.Warn(err.Error())
			continue
		}

		for configType, configFile := range configs {
			if configType == "micro" {
				_, previousGitRef, err := exec.GetPreviousFromSourceMicrobenchmark(s.dbClient, source, ref)
				if err != nil {
					slog.Warn(err.Error())
					continue
				}

				elements = append(elements, s.createBranchElementWithComparisonOnPreviousAndRelease(configFile, ref, configType, previousGitRef, "", source, lastPatchRelease)...)
			} else {
				versions := git.GetPlannerVersionsForRelease(release)

				for _, version := range versions {
					_, previousGitRef, err := exec.GetPreviousFromSourceMacrobenchmark(s.dbClient, source, configType, string(version), ref)
					if err != nil {
						slog.Warn(err.Error())
						continue
					}

					elements = append(elements, s.createBranchElementWithComparisonOnPreviousAndRelease(configFile, ref, configType, previousGitRef, string(version), source, lastPatchRelease)...)
				}
			}
		}
	}
	return elements, nil
}

func (s *Server) createBranchElementWithComparisonOnPreviousAndRelease(configFile, ref, configType, previousGitRef, plannerVersion, source string, lastRelease *git.Release) []*executionQueueElement {
	var elements []*executionQueueElement

	// creating a benchmark for the latest commit on the branch with SourceCron as a source
	// this benchmark will be compared with the previous benchmark on branch with the given source, and with the latest release
	newExecutionElement := s.createSimpleExecutionQueueElement(source, configFile, ref, configType, plannerVersion, false, 0)
	elements = append(elements, newExecutionElement)

	if previousGitRef != "" {
		// creating an execution queue element for the latest benchmark with SourceCron as source
		// this will not be executed since the benchmark already exist, we still create the element in order to compare
		previousElement := s.createSimpleExecutionQueueElement(source, configFile, previousGitRef, configType, plannerVersion, false, 0)
		previousElement.compareWith = append(previousElement.compareWith, newExecutionElement.identifier)
		newExecutionElement.compareWith = append(newExecutionElement.compareWith, previousElement.identifier)
		elements = append(elements, previousElement)
	}

	if lastRelease != nil {
		// creating an execution queue element for the latest release (comparing branch with the latest release)
		// this will probably not be executed the benchmark should already exist, we still create it to compare main once its benchmark is over
		lastReleaseElement := s.createSimpleExecutionQueueElement(exec.SourceTag+lastRelease.Name, configFile, lastRelease.CommitHash, configType, plannerVersion, false, 0)
		lastReleaseElement.compareWith = append(lastReleaseElement.compareWith, newExecutionElement.identifier)
		newExecutionElement.compareWith = append(newExecutionElement.compareWith, lastReleaseElement.identifier)
		elements = append(elements, lastReleaseElement)
	}
	return elements
}

func (s *Server) createSimpleExecutionQueueElement(source, configFile, ref, configType, plannerVersion string, notify bool, pullNb int) *executionQueueElement {
	return &executionQueueElement{
		config:       configFile,
		retry:        s.cronNbRetry,
		notifyAlways: notify,
		identifier: executionIdentifier{
			GitRef:         ref,
			Source:         source,
			BenchmarkType:  configType,
			PlannerVersion: plannerVersion,
			PullNb:         pullNb,
		},
	}
}

func (s Server) pullRequestsCronHandler() {
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
			for configType, configFile := range configs {
				ref := prInfo.SHA
				previousGitRef := prInfo.Base
				pullNb := prInfo.Number
				if configType == "micro" {
					elements = append(elements, s.createPullRequestElementWithBaseComparison(configFile, ref, configType, previousGitRef, "", pullNb)...)
				} else {
					versions := []macrobench.PlannerVersion{macrobench.V3Planner}
					if labelInfo.useGen4 {
						versions = append(versions, macrobench.Gen4FallbackPlanner)
					}
					for _, version := range versions {
						elements = append(elements, s.createPullRequestElementWithBaseComparison(configFile, ref, configType, previousGitRef, version, pullNb)...)
					}
				}
			}
		}
	}
	for _, element := range elements {
		s.addToQueue(element)
	}
}

func (s Server) createPullRequestElementWithBaseComparison(configFile, ref, configType, previousGitRef string, version macrobench.PlannerVersion, pullNb int) []*executionQueueElement {
	var elements []*executionQueueElement

	newExecutionElement := s.createSimpleExecutionQueueElement(exec.SourcePullRequest, configFile, ref, configType, string(version), true, pullNb)
	elements = append(elements, newExecutionElement)

	if previousGitRef != "" {
		previousElement := s.createSimpleExecutionQueueElement(exec.SourcePullRequestBase, configFile, previousGitRef, configType, string(version), false, pullNb)
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
		source := exec.SourceTag+release.Name
		for configType, configFile := range configs {
			if configType == "micro" {
				elements = append(elements, s.createSimpleExecutionQueueElement(source, configFile, release.CommitHash, configType, "", true, 0))
			} else {
				versions := git.GetPlannerVersionsForRelease(release)
				for _, version := range versions {
					elements = append(elements, s.createSimpleExecutionQueueElement(source, configFile, release.CommitHash, configType, string(version), true, 0))
				}
			}
		}
	}
	for _, element := range elements {
		s.addToQueue(element)
	}
}
