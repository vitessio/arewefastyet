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

type execInfo struct {
	config         string
	ref            string
	retry          int
	source         string
	plannerVersion string
	execType       string
}

const (
	maxConcurJob = 5
)

var (
	execQueue        chan execInfo
	currentCountExec int
	mtx              *sync.Mutex
)

func (s *Server) createNewCron() error {
	if s.cronSchedule == "" {
		return nil
	}
	execQueue = make(chan execInfo)
	mtx = &sync.Mutex{}

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
	ref, err := git.GetCommitHash(s.getVitessPath())
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	var execInfos []execInfo
	for configType, configFile := range configs {
		if configType == "micro" {
			execInfos = append(execInfos, execInfo{
				config: configFile,
				ref:    ref,
				retry:  s.cronNbRetry,
				source: "cron",
			})
		} else {
			for _, version := range macrobench.PlannerVersions {
				execInfos = append(execInfos, execInfo{
					config:         configFile,
					ref:            ref,
					retry:          s.cronNbRetry,
					source:         "cron",
					plannerVersion: version,
				})
			}
		}
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
				execInfos = append(execInfos, execInfo{
					config: configFile,
					ref:    ref,
					retry:  s.cronNbRetry,
					source: "cron_" + release.Name,
				})
			} else {
				for _, version := range macrobench.PlannerVersions {
					execInfos = append(execInfos, execInfo{
						config:         configFile,
						ref:            ref,
						retry:          s.cronNbRetry,
						source:         "cron_" + release.Name,
						plannerVersion: version,
					})
				}
			}
		}
		s.cronPrepare(execInfos)
	}
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
				execInfos = append(execInfos, execInfo{
					config: configFile,
					ref:    prInfo.SHA,
					retry:  s.cronNbRetry,
					source: source,
				})
			} else {
				for _, version := range macrobench.PlannerVersions {
					execInfos = append(execInfos, execInfo{
						config:         configFile,
						ref:            prInfo.SHA,
						retry:          s.cronNbRetry,
						source:         source,
						plannerVersion: version,
					})
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
				execInfos = append(execInfos, execInfo{
					config: configFile,
					ref:    release.CommitHash,
					retry:  s.cronNbRetry,
					source: "cron_tags_" + release.Name,
				})
			} else {
				for _, version := range macrobench.PlannerVersions {
					execInfos = append(execInfos, execInfo{
						config:         configFile,
						ref:            release.CommitHash,
						retry:          s.cronNbRetry,
						source:         "cron_tags_" + release.Name,
						plannerVersion: version,
					})
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
		if currentCountExec >= maxConcurJob {
			continue
		}
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
