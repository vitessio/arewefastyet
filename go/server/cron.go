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
	"sync"
	"time"

	"github.com/vitessio/arewefastyet/go/tools/macrobench"

	"github.com/robfig/cron/v3"
	"github.com/vitessio/arewefastyet/go/exec"
	"github.com/vitessio/arewefastyet/go/tools/git"
)

type CompareInfo struct {
	Config              string
	execMain            *execInfo
	execComp            *execInfo
	Retry               int
	PlannerVersion      string
	ExecType            string
	Name                string
	IgnoreNonRegression bool
}

type execInfo struct {
	ref    string
	source string
}

const (
	maxConcurJob = 5
)

var (
	execQueue        chan *CompareInfo
	currentCountExec int
	mtx              *sync.RWMutex
)

func (s *Server) createNewCron() error {
	if s.cronSchedule == "" {
		return nil
	}
	execQueue = make(chan *CompareInfo)
	mtx = &sync.RWMutex{}

	c := cron.New()

	jobs := []func(){
		s.cronBranchHandler,
		s.cronPRLabels,
		s.cronTags,
	}
	for _, job := range jobs {
		_, err := c.AddFunc(s.cronSchedule, job)
		if err != nil {
			return err
		}
	}
	c.Start()
	go s.cron()
	return nil
}

func (s *Server) cronBranchHandler() {
	configs := map[string]string{
		"micro": s.microbenchConfigPath,
		"oltp":  s.macrobenchConfigPathOLTP,
		"tpcc":  s.macrobenchConfigPathTPCC,
	}

	err := s.pullLocalVitess()
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	// master branch
	compareInfos, err := s.compareMasterBranch(configs)
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	// release branches
	compareInfosReleases, err := s.compareReleaseBranches(configs)
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	compareInfos = append(compareInfos, compareInfosReleases...)
	s.cronPrepare(compareInfos)
}

func (s *Server) compareMasterBranch(configs map[string]string) ([]*CompareInfo, error) {
	var compareInfos []*CompareInfo
	ref, err := git.GetCommitHash(s.getVitessPath())
	if err != nil {
		slog.Warn(err.Error())
		return nil, err
	}
	lastRelease, err := git.GetLastReleaseAndCommitHash(s.getVitessPath())
	if err != nil {
		slog.Warn(err.Error())
		return nil, nil
	}
	for configType, configFile := range configs {
		if configType == "micro" {
			_, previousGitRef, err := exec.GetPreviousFromSourceMicrobenchmark(s.dbClient, "cron", ref)
			if err != nil {
				slog.Warn(err.Error())
				return nil, err
			}
			compareInfos = append(compareInfos, newCompareInfo(configFile, ref, "cron", previousGitRef, "", s.cronNbRetry, configType, "", false))
			compareInfos = append(compareInfos, newCompareInfo(configFile, ref, "cron", lastRelease.CommitHash, "cron_tags_"+lastRelease.Name, s.cronNbRetry, configType, "", false))
		} else {
			for _, version := range macrobench.PlannerVersions {
				_, previousGitRef, err := exec.GetPreviousFromSourceMacrobenchmark(s.dbClient, "cron", configType, string(version), ref)
				if err != nil {
					slog.Warn(err.Error())
					return nil, err
				}
				compareInfos = append(compareInfos, newCompareInfo(configFile, ref, "cron", previousGitRef, "", s.cronNbRetry, configType, string(version), false))
				compareInfos = append(compareInfos, newCompareInfo(configFile, ref, "cron", lastRelease.CommitHash, "cron_tags_"+lastRelease.Name, s.cronNbRetry, configType, string(version), false))
			}
		}
	}
	return compareInfos, nil
}

func (s *Server) compareReleaseBranches(configs map[string]string) ([]*CompareInfo, error) {
	var compareInfos []*CompareInfo
	releases, err := git.GetLatestVitessReleaseBranchCommitHash(s.getVitessPath())
	if err != nil {
		slog.Warn(err.Error())
		return nil, err
	}
	for _, release := range releases {
		ref := release.CommitHash
		source := "cron_" + release.Name
		lastPathRelease, err := git.GetLastPatchReleaseAndCommitHash(s.getVitessPath(), release.Number)
		if err != nil {
			slog.Warn(err.Error())
			return nil, nil
		}
		for configType, configFile := range configs {
			if configType == "micro" {
				_, previousGitRef, err := exec.GetPreviousFromSourceMicrobenchmark(s.dbClient, source, ref)
				if err != nil {
					slog.Warn(err.Error())
					return nil, err
				}
				compareInfos = append(compareInfos, newCompareInfo(configFile, ref, source, previousGitRef, "", s.cronNbRetry, configType, "", false))
				compareInfos = append(compareInfos, newCompareInfo(configFile, ref, source, lastPathRelease.CommitHash, "cron_tags_"+lastPathRelease.Name, s.cronNbRetry, configType, "", false))
			} else {
				for _, version := range macrobench.PlannerVersions {
					_, previousGitRef, err := exec.GetPreviousFromSourceMacrobenchmark(s.dbClient, source, configType, string(version), ref)
					if err != nil {
						slog.Warn(err.Error())
						return nil, err
					}
					compareInfos = append(compareInfos, newCompareInfo(configFile, ref, source, previousGitRef, "", s.cronNbRetry, configType, string(version), false))
					compareInfos = append(compareInfos, newCompareInfo(configFile, ref, source, lastPathRelease.CommitHash, "cron_tags_"+lastPathRelease.Name, s.cronNbRetry, configType, string(version), false))
				}
			}
		}
	}

	return compareInfos, nil
}

func (s Server) cronPRLabels() {
	configs := map[string]string{
		"micro": s.microbenchConfigPath,
		"oltp":  s.macrobenchConfigPathOLTP,
		"tpcc":  s.macrobenchConfigPathTPCC,
	}

	prInfos, err := git.GetPullRequestHeadForLabels([]string{s.prLabelTrigger}, "vitessio/vitess")
	if err != nil {
		slog.Error(err)
		return
	}

	// a slice of compareInfo's
	var compareInfos []*CompareInfo
	source := "cron_pr"
	for _, prInfo := range prInfos {
		for configType, configFile := range configs {
			ref := prInfo.SHA
			previousGitRef := prInfo.Base
			if configType == "micro" {
				compareInfos = append(compareInfos, newCompareInfo(configFile, ref, source, previousGitRef, source, s.cronNbRetry, configType, "", true))
			} else {
				for _, version := range macrobench.PlannerVersions {
					compareInfos = append(compareInfos, newCompareInfo(configFile, ref, source, previousGitRef, source, s.cronNbRetry, configType, string(version), true))
				}
			}
		}
	}
	s.cronPrepare(compareInfos)
}

func (s *Server) cronTags() {
	configs := map[string]string{
		"micro": s.microbenchConfigPath,
		"oltp":  s.macrobenchConfigPathOLTP,
		"tpcc":  s.macrobenchConfigPathTPCC,
	}

	releases, err := git.GetLatestVitessReleaseCommitHash(s.getVitessPath())
	if err != nil {
		slog.Error(err)
		return
	}

	// a slice of compareInfo's
	var compareInfos []*CompareInfo
	for _, release := range releases {
		for configType, configFile := range configs {
			if configType == "micro" {
				compareInfos = append(compareInfos, newSingleExecution(configFile, release.CommitHash, "cron_tags_"+release.Name, s.cronNbRetry, configType, ""))
			} else {
				for _, version := range macrobench.PlannerVersions {
					compareInfos = append(compareInfos, newSingleExecution(configFile, release.CommitHash, "cron_tags_"+release.Name, s.cronNbRetry, configType, string(version)))
				}
			}
		}
	}
	s.cronPrepare(compareInfos)
}

func (s *Server) cronPrepare(compareInfos []*CompareInfo) {
	for _, info := range compareInfos {
		execQueue <- info
		slog.Infof("New Comparison Execution - Name: %s (config: %s, refMain: %s, sourceMain: %s, planner: %s) added to the queue (length: %d)", info.Name, info.Config, info.execMain.ref, info.execMain.source, info.PlannerVersion, len(execQueue))
	}
}

func (s *Server) checkIfExists(execIn execInfo) (bool, error) {
	if execIn.execType == "micro" {
		exist, err := exec.Exists(s.dbClient, execIn.ref, "cron", execIn.execType, exec.StatusFinished)
		if err != nil {
			slog.Error(err)
			return false, err
		}
		return exist, nil
	}
	exist, err := exec.ExistsMacrobenchmark(s.dbClient, execIn.ref, "cron", execIn.execType, exec.StatusFinished, execIn.plannerVersion)
	if err != nil {
		slog.Error(err)
		return false, err
	}
	return exist, nil
}

func (s *Server) cron() {
	for {
		time.Sleep(time.Second * 1)
		mtx.RLock()
		if currentCountExec >= maxConcurJob {
			mtx.RUnlock()
			continue
		}
		mtx.RUnlock()
		e := <-execQueue
		mtx.Lock()
		currentCountExec++
		mtx.Unlock()
		go s.cronExecution(e)
	}
}

func (s *Server) cronExecution(eI execInfo) {
	var err error
	var e *exec.Exec
	defer func() {
		if e != nil {
			err = e.CleanUp()
			if err != nil {
				slog.Errorf("CleanUp step: %v", err)
			}
		}
		if err != nil {
			// Retry after execution failure if the counter is above zero.
			if eI.retry > 0 {
				eI.retry--
				slog.Info("Retrying execution for source: ", eI.source, " git ref: ", eI.ref, " with configuration: ", eI.config, ". Number of retry left: ", eI.retry)
				s.cronExecution(eI)
				return
			}
		}
		e.Success()
		mtx.Lock()
		currentCountExec--
		mtx.Unlock()
	}()

	e, err = exec.NewExecWithConfig(eI.config)
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	slog.Info("Created new execution: ", e.UUID.String())
	e.Source = eI.source
	e.GitRef = eI.ref
	e.VtgatePlannerVersion = eI.plannerVersion

	slog.Info("Started execution: ", e.UUID.String(), ", with git ref: ", eI.ref)
	err = e.Prepare()
	if err != nil {
		slog.Errorf("Prepare step: %v", err)
		return
	}

	err = e.SetOutputToDefaultPath()
	if err != nil {
		slog.Errorf("Prepare outputs step: %v", err)
		return
	}

	err = e.Execute()
	if err != nil {
		slog.Errorf("Execution step: %v", err)
		return
	}

	err = e.SendNotificationForRegression()
	if err != nil {
		slog.Errorf("Send notification: %v", err)
		return
	}
	slog.Info("Finished execution: ", e.UUID.String())
}

func newCompareInfo(configFile, ref, source, compareRef, compareSource string, retry int, configType, plannerVersion string, ignoreNonRegression bool) *CompareInfo {
	return &CompareInfo{
		Config: configFile,
		execMain: &execInfo{
			ref:    ref,
			source: source,
		},
		execComp: &execInfo{
			ref:    compareRef,
			source: compareSource,
		},
		Retry:               retry,
		PlannerVersion:      plannerVersion,
		ExecType:            configType,
		IgnoreNonRegression: ignoreNonRegression,
	}
}

func newSingleExecution(configFile, ref, source string, retry int, configType, plannerVersion string) *CompareInfo {
	return &CompareInfo{
		Config: configFile,
		execMain: &execInfo{
			ref:    ref,
			source: source,
		},
		execComp:            nil,
		Retry:               retry,
		PlannerVersion:      plannerVersion,
		ExecType:            configType,
		IgnoreNonRegression: false,
	}
}
