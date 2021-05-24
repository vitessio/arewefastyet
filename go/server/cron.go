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
	releases, err := git.GetLatestVitessReleaseBranchCommitHash(s.getVitessPath())
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	for _, release := range releases {
		ref = release.CommitHash
		for configType, configFile := range configs {
			if configType == "micro" {
				compareInfos = append(compareInfos, newExecInfo(configFile, ref, s.cronNbRetry, "cron_"+release.Name, "", configType))
			} else {
				for _, version := range macrobench.PlannerVersions {
					compareInfos = append(compareInfos, newExecInfo(configFile, ref, s.cronNbRetry, "cron_"+release.Name, string(version), configType))
				}
			}
		}
	}
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
			compareInfos = append(compareInfos, newCompareInfo(configFile, ref, "cron", previousGitRef, "", s.cronNbRetry, configType, ""))
			compareInfos = append(compareInfos, newCompareInfo(configFile, ref, "cron", lastRelease.CommitHash, "cron_tags_"+lastRelease.Name, s.cronNbRetry, configType, ""))
		} else {
			for _, version := range macrobench.PlannerVersions {
				_, previousGitRef, err := exec.GetPreviousFromSourceMacrobenchmark(s.dbClient, "cron", configType, string(version), ref)
				if err != nil {
					slog.Warn(err.Error())
					return nil, err
				}
				compareInfos = append(compareInfos, newCompareInfo(configFile, ref, "cron", previousGitRef, "", s.cronNbRetry, configType, string(version)))
				compareInfos = append(compareInfos, newCompareInfo(configFile, ref, "cron", lastRelease.CommitHash, "cron_tags_"+lastRelease.Name, s.cronNbRetry, configType, string(version)))
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

	// a slice of execInfos
	var execInfos []execInfo
	source := "cron_pr"
	for _, prInfo := range prInfos {
		for configType, configFile := range configs {
			if configType == "micro" {
				execInfos = append(execInfos, newExecInfo(configFile, prInfo.SHA, s.cronNbRetry, source, "", configType))
			} else {
				for _, version := range macrobench.PlannerVersions {
					execInfos = append(execInfos, newExecInfo(configFile, prInfo.SHA, s.cronNbRetry, source, string(version), configType))
				}
			}
		}
	}
	s.cronPrepare(execInfos)
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

	// a slice of execInfos
	var execInfos []execInfo
	for _, release := range releases {
		for configType, configFile := range configs {
			if configType == "micro" {
				execInfos = append(execInfos, newExecInfo(configFile, release.CommitHash, s.cronNbRetry, "cron_tags_"+release.Name, "", configType))
			} else {
				for _, version := range macrobench.PlannerVersions {
					execInfos = append(execInfos, newExecInfo(configFile, release.CommitHash, s.cronNbRetry, "cron_tags_"+release.Name, string(version), configType))
				}
			}
		}
	}
	s.cronPrepare(execInfos)
}

func (s *Server) cronPrepare(execInfos []execInfo) {
	for _, info := range execInfos {
		exists, err := s.checkIfExists(info)
		if err != nil || exists {
			continue
		}
		execQueue <- info
		slog.Infof("New Execution (config: %s, ref: %s, source: %s, planner: %s) added to the queue (length: %d)", info.config, info.ref, info.source, info.plannerVersion, len(execQueue))
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

func newCompareInfo(configFile, ref, source, compareRef, compareSource string, retry int, configType, plannerVersion string) *CompareInfo {
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
		IgnoreNonRegression: false,
	}
}
