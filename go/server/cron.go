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
	"strconv"
	"sync"
	"time"

	"github.com/vitessio/arewefastyet/go/tools/macrobench"

	"github.com/robfig/cron/v3"
	"github.com/vitessio/arewefastyet/go/exec"
	"github.com/vitessio/arewefastyet/go/tools/git"
)

type CompareInfo struct {
	config              string
	execMain            *execInfo
	execComp            *execInfo
	retry               int
	plannerVersion      string
	typeOf              string
	name                string
	ignoreNonRegression bool
}

type execInfo struct {
	ref    string
	pullNB int
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
			compareInfos = append(compareInfos, newCompareInfo("Comparing master with previous master - micro", configFile, ref, "cron", 0, previousGitRef, "", s.cronNbRetry, configType, "", false))
			compareInfos = append(compareInfos, newCompareInfo("Comparing master with latest release - "+lastRelease.Name+" - micro", configFile, ref, "cron", 0, lastRelease.CommitHash, "cron_tags_"+lastRelease.Name, s.cronNbRetry, configType, "", false))
		} else {
			for _, version := range macrobench.PlannerVersions {
				_, previousGitRef, err := exec.GetPreviousFromSourceMacrobenchmark(s.dbClient, "cron", configType, string(version), ref)
				if err != nil {
					slog.Warn(err.Error())
					return nil, err
				}
				compareInfos = append(compareInfos, newCompareInfo("Comparing master with previous master - "+configType+" - "+string(version), configFile, ref, "cron", 0, previousGitRef, "", s.cronNbRetry, configType, string(version), false))
				compareInfos = append(compareInfos, newCompareInfo("Comparing master with latest release - "+lastRelease.Name+" - "+configType+" - "+string(version), configFile, ref, "cron", 0, lastRelease.CommitHash, "cron_tags_"+lastRelease.Name, s.cronNbRetry, configType, string(version), false))
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
				compareInfos = append(compareInfos, newCompareInfo("Comparing "+release.Name+" with previous commit - micro", configFile, ref, source, 0, previousGitRef, "", s.cronNbRetry, configType, "", false))
				compareInfos = append(compareInfos, newCompareInfo("Comparing "+release.Name+" with last path release "+lastPathRelease.Name+" - micro", configFile, ref, source, 0, lastPathRelease.CommitHash, "cron_tags_"+lastPathRelease.Name, s.cronNbRetry, configType, "", false))
			} else {
				for _, version := range macrobench.PlannerVersions {
					_, previousGitRef, err := exec.GetPreviousFromSourceMacrobenchmark(s.dbClient, source, configType, string(version), ref)
					if err != nil {
						slog.Warn(err.Error())
						return nil, err
					}
					compareInfos = append(compareInfos, newCompareInfo("Comparing "+release.Name+" with previous commit - "+configType+" - "+string(version), configFile, ref, source, 0, previousGitRef, "", s.cronNbRetry, configType, string(version), false))
					compareInfos = append(compareInfos, newCompareInfo("Comparing "+release.Name+" with last path release "+lastPathRelease.Name+" - "+configType+" - "+string(version), configFile, ref, source, 0, lastPathRelease.CommitHash, "cron_tags_"+lastPathRelease.Name, s.cronNbRetry, configType, string(version), false))
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
			// TODO : Get non zero PR number and use it in execution
			pullNb := 0
			if configType == "micro" {
				compareInfos = append(compareInfos, newCompareInfo("Comparing pull request number - "+strconv.Itoa(pullNb)+" - micro", configFile, ref, source, pullNb, previousGitRef, source, s.cronNbRetry, configType, "", true))
			} else {
				for _, version := range macrobench.PlannerVersions {
					compareInfos = append(compareInfos, newCompareInfo("Comparing pull request number - "+strconv.Itoa(pullNb)+" - "+configType+" - "+string(version), configFile, ref, source, pullNb, previousGitRef, source, s.cronNbRetry, configType, string(version), true))
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
				compareInfos = append(compareInfos, newSingleExecution("Tag run", configFile, release.CommitHash, "cron_tags_"+release.Name, s.cronNbRetry, configType, ""))
			} else {
				for _, version := range macrobench.PlannerVersions {
					compareInfos = append(compareInfos, newSingleExecution("Tag run", configFile, release.CommitHash, "cron_tags_"+release.Name, s.cronNbRetry, configType, string(version)))
				}
			}
		}
	}
	s.cronPrepare(compareInfos)
}

func (s *Server) cronPrepare(compareInfos []*CompareInfo) {
	for _, info := range compareInfos {
		execQueue <- info
		slog.Infof("New Comparison Execution - Name: %s (config: %s, refMain: %s, sourceMain: %s, planner: %s) added to the queue (length: %d)", info.name, info.config, info.execMain.ref, info.execMain.source, info.plannerVersion, len(execQueue))
	}
}

func (s *Server) checkIfExists(ref, typeOf, plannerVersion string) (bool, error) {
	if typeOf == "micro" {
		exist, err := exec.Exists(s.dbClient, ref, "cron", typeOf, exec.StatusFinished)
		if err != nil {
			slog.Error(err)
			return false, err
		}
		return exist, nil
	}
	exist, err := exec.ExistsMacrobenchmark(s.dbClient, ref, "cron", typeOf, exec.StatusFinished, plannerVersion)
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

func (s *Server) cronExecution(compInfo *CompareInfo) {
	var err error
	defer func() {
		if err != nil {
			// Retry after any failure if the counter is above zero.
			if compInfo.retry > 0 {
				compInfo.retry--
				slog.Infof("Retrying Cron Execution - Name: %s (config: %s, refMain: %s, sourceMain: %s, planner: %s) retries left: %d", compInfo.name, compInfo.config, compInfo.execMain.ref, compInfo.execMain.source, compInfo.plannerVersion, len(execQueue), compInfo.retry)
				s.cronExecution(compInfo)
				return
			}
		}
	}()

	err = s.executeSingle(compInfo.config, compInfo.execMain.source, compInfo.execMain.ref, compInfo.typeOf, compInfo.plannerVersion)
	if err != nil {
		slog.Errorf("Error while single execution: %v", err)
		return
	}

	if compInfo.execComp != nil || compInfo.execComp.source == "" {
		err = s.executeSingle(compInfo.config, compInfo.execComp.source, compInfo.execComp.ref, compInfo.typeOf, compInfo.plannerVersion)
		if err != nil {
			slog.Errorf("Error while single execution: %v", err)
			return
		}
	}

	err = s.sendNotificationForRegression(compInfo)
	if err != nil {
		slog.Errorf("Send notification: %v", err)
		return
	}
}

func (s *Server) executeSingle(config, source, ref, typeOf, plannerVersion string) (err error) {
	var e *exec.Exec
	var exists bool
	defer func() {
		if e != nil {
			err2 := e.CleanUp()
			if err2 != nil {
				slog.Errorf("CleanUp step: %v", err2)
				err = err2
				return
			}
			e.Success()
			slog.Info("Finished execution: ", e.UUID.String())
			mtx.Lock()
			currentCountExec--
			mtx.Unlock()
		}
	}()

	exists, err = s.checkIfExists(ref, typeOf, plannerVersion)
	if exists || err != nil {
		return err
	}

	e, err = exec.NewExecWithConfig(config)
	if err != nil {
		slog.Warn(err.Error())
		return err
	}
	slog.Info("Created new execution: ", e.UUID.String())
	e.Source = source
	e.GitRef = ref
	e.VtgatePlannerVersion = plannerVersion

	slog.Info("Started execution: ", e.UUID.String(), ", with git ref: ", ref)
	err = e.Prepare()
	if err != nil {
		slog.Errorf("Prepare step: %v", err)
		return err
	}

	err = e.SetOutputToDefaultPath()
	if err != nil {
		slog.Errorf("Prepare outputs step: %v", err)
		return err
	}

	err = e.Execute()
	if err != nil {
		slog.Errorf("Execution step: %v", err)
		return err
	}
	return nil
}

func newCompareInfo(name, configFile, ref, source string, pullNB int, compareRef, compareSource string, retry int, configType, plannerVersion string, ignoreNonRegression bool) *CompareInfo {
	return &CompareInfo{
		config: configFile,
		execMain: &execInfo{
			ref:    ref,
			pullNB: pullNB,
			source: source,
		},
		execComp: &execInfo{
			ref:    compareRef,
			pullNB: pullNB,
			source: compareSource,
		},
		retry:               retry,
		plannerVersion:      plannerVersion,
		typeOf:              configType,
		ignoreNonRegression: ignoreNonRegression,
		name:                name,
	}
}

func newSingleExecution(name, configFile, ref, source string, retry int, configType, plannerVersion string) *CompareInfo {
	return &CompareInfo{
		config: configFile,
		execMain: &execInfo{
			ref:    ref,
			source: source,
		},
		execComp:            nil,
		retry:               retry,
		plannerVersion:      plannerVersion,
		typeOf:              configType,
		ignoreNonRegression: false,
	}
}
