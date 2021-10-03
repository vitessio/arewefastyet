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
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/vitessio/arewefastyet/go/tools/macrobench"

	"github.com/robfig/cron/v3"
	"github.com/vitessio/arewefastyet/go/exec"
	"github.com/vitessio/arewefastyet/go/tools/git"
)

type (
	// CompareInfo has the details required to compare two git commits
	CompareInfo struct {
		// config is the configuration file to use
		config string
		// execMain contains the details of the execution of the main commit
		execMain *execInfo
		// execComp contains the details of the execution of the secondary commit
		// when execComp is nil, then there is no execution to be done
		// when execComp.source == "", then we are sure that the results already exist and do not need to rerun
		execComp *execInfo
		// retry is the number of times we want to retry failed comparisons
		retry int
		// plannerVersion is the vtgate planner that we should be using for testing
		plannerVersion string
		// typeOf takes 3 values = [oltp, micro, tpcc]
		typeOf string
		// name is the name of the comparison, used in sending slack message
		name string
		// ignoreNonRegression is true when we want to send a slack message even when there is no regression
		ignoreNonRegression bool
		// pullNb is the number of the related pull request if any
		pullNb int
	}

	// execInfo contains execution information regarding each exec, which is not common between the 2 executions
	execInfo struct {
		ref    string
		pullNB int
		source string
	}

	executionStatus int
)

const (
	executionFailed executionStatus = iota
	executionSucceeded
	executionExists

	// maxConcurJob is the maximum number of concurrent jobs that we can execute
	maxConcurJob = 5
)

var (
	execQueue        chan *CompareInfo
	currentCountExec int
	mtx              *sync.RWMutex
)

func createIndividualCron(schedule string, jobs []func()) error {
	if schedule == "" {
		return nil
	}

	c := cron.New()
	for _, job := range jobs {
		_, err := c.AddFunc(schedule, job)
		if err != nil {
			return err
		}
	}
	c.Start()
	return nil
}

func (s *Server) createCrons() error {
	if s.cronSchedule == "" {
		return nil
	}
	execQueue = make(chan *CompareInfo)
	mtx = &sync.RWMutex{}

	err := createIndividualCron(s.cronSchedule, []func(){
		s.cronBranchHandler,
		s.cronTags,
	})
	if err != nil {
		return err
	}
	err = createIndividualCron(s.cronSchedulePullRequests, []func(){s.cronPullRequests})
	if err != nil {
		return err
	}

	go s.cronExecutionQueueWatcher()
	return nil
}

func (s *Server) cronBranchHandler() {
	err := s.pullLocalVitess()
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	// main branch
	compareInfos, err := s.compareMainBranch()
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	// release branches
	compareInfosReleases, err := s.compareReleaseBranches()
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	compareInfos = append(compareInfos, compareInfosReleases...)
	s.cronPrepare(compareInfos)
}

func (s *Server) compareMainBranch() ([]*CompareInfo, error) {
	configs := s.getConfigFiles()

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
	// We compare main with the previous hash of main and with the latest release
	for configType, configFile := range configs {
		if configType == "micro" {
			_, previousGitRef, err := exec.GetPreviousFromSourceMicrobenchmark(s.dbClient, exec.SourceCron, ref)
			if err != nil {
				slog.Warn(err.Error())
			} else if previousGitRef != "" {
				compareInfos = append(compareInfos, newCompareInfo("Comparing main with previous main - micro", configFile, ref, exec.SourceCron, 0, previousGitRef, "", s.cronNbRetry, configType, "", false))
			}
			compareInfos = append(compareInfos, newCompareInfo("Comparing main with latest release - "+lastRelease.Name+" - micro", configFile, ref, exec.SourceCron, 0, lastRelease.CommitHash, "cron_tags_"+lastRelease.Name, s.cronNbRetry, configType, "", false))
		} else {
			for _, version := range macrobench.PlannerVersions {
				_, previousGitRef, err := exec.GetPreviousFromSourceMacrobenchmark(s.dbClient, exec.SourceCron, configType, string(version), ref)
				if err != nil {
					slog.Warn(err.Error())
				} else if previousGitRef != "" {
					compareInfos = append(compareInfos, newCompareInfo("Comparing main with previous main - "+configType+" - "+string(version), configFile, ref, exec.SourceCron, 0, previousGitRef, "", s.cronNbRetry, configType, string(version), false))
				}
				compareInfos = append(compareInfos, newCompareInfo("Comparing main with latest release - "+lastRelease.Name+" - "+configType+" - "+string(version), configFile, ref, exec.SourceCron, 0, lastRelease.CommitHash, "cron_tags_"+lastRelease.Name, s.cronNbRetry, configType, string(version), false))
			}
		}
	}
	return compareInfos, nil
}

func (s *Server) compareReleaseBranches() ([]*CompareInfo, error) {
	configs := s.getConfigFiles()

	var compareInfos []*CompareInfo
	releases, err := git.GetLatestVitessReleaseBranchCommitHash(s.getVitessPath())
	if err != nil {
		slog.Warn(err.Error())
		return nil, err
	}

	// We compare release-branches with the previous hash of that branch and with the latest patch release of that version
	for _, release := range releases {
		ref := release.CommitHash
		source := exec.SourceReleaseBranch + release.Name
		lastPathRelease, err := git.GetLastPatchReleaseAndCommitHash(s.getVitessPath(), release.Number)
		if err != nil {
			slog.Warn(err.Error())
		}
		for configType, configFile := range configs {
			if configType == "micro" {
				_, previousGitRef, err := exec.GetPreviousFromSourceMicrobenchmark(s.dbClient, source, ref)
				if err != nil {
					slog.Warn(err.Error())
					return nil, err
				} else if previousGitRef != "" {
					compareInfos = append(compareInfos, newCompareInfo("Comparing "+release.Name+" with previous commit - micro", configFile, ref, source, 0, previousGitRef, "", s.cronNbRetry, configType, "", false))
				}
				if lastPathRelease != nil {
					compareInfos = append(compareInfos, newCompareInfo("Comparing "+release.Name+" with last path release "+lastPathRelease.Name+" - micro", configFile, ref, source, 0, lastPathRelease.CommitHash, "cron_tags_"+lastPathRelease.Name, s.cronNbRetry, configType, "", false))
				}
			} else {
				versions := []macrobench.PlannerVersion{macrobench.V3Planner}
				if release.Number[0] >= 10 {
					versions = append(versions, macrobench.Gen4FallbackPlanner)
				}
				for _, version := range versions {
					_, previousGitRef, err := exec.GetPreviousFromSourceMacrobenchmark(s.dbClient, source, configType, string(version), ref)
					if err != nil {
						slog.Warn(err.Error())
					} else if previousGitRef != "" {
						compareInfos = append(compareInfos, newCompareInfo("Comparing "+release.Name+" with previous commit - "+configType+" - "+string(version), configFile, ref, source, 0, previousGitRef, "", s.cronNbRetry, configType, string(version), false))
					}
					if lastPathRelease != nil {
						compareInfos = append(compareInfos, newCompareInfo("Comparing "+release.Name+" with last path release "+lastPathRelease.Name+" - "+configType+" - "+string(version), configFile, ref, source, 0, lastPathRelease.CommitHash, "cron_tags_"+lastPathRelease.Name, s.cronNbRetry, configType, string(version), false))
					}
				}
			}
		}
	}

	return compareInfos, nil
}

func (s Server) cronPullRequests() {
	configs := s.getConfigFiles()
	prLabelsInfo := []struct {
		label   string
		useGen4 bool
	}{
		{label: s.prLabelTrigger, useGen4: true},
		{label: s.prLabelTriggerV3, useGen4: false},
	}

	// a slice of compareInfo's
	var compareInfos []*CompareInfo

	for _, labelInfo := range prLabelsInfo {
		prInfos, err := git.GetPullRequestsFromGitHub([]string{labelInfo.label}, "vitessio/vitess")
		if err != nil {
			slog.Error(err)
			return
		}

		// We compare PRs with the base of the PR
		for _, prInfo := range prInfos {
			for configType, configFile := range configs {
				ref := prInfo.SHA
				previousGitRef := prInfo.Base
				pullNb := prInfo.Number
				if configType == "micro" {
					compareInfos = append(compareInfos, newCompareInfo(
						"Comparing pull request number - "+strconv.Itoa(pullNb)+" - micro",
						configFile,
						ref,
						exec.SourcePullRequest,
						pullNb,
						previousGitRef,
						exec.SourcePullRequestBase,
						s.cronNbRetry,
						configType,
						"",
						true,
					))
				} else {
					versions := []macrobench.PlannerVersion{macrobench.V3Planner}
					if labelInfo.useGen4 {
						versions = append(versions, macrobench.Gen4FallbackPlanner)
					}
					for _, version := range versions {
						compareInfos = append(compareInfos, newCompareInfo(
							"Comparing pull request number - "+strconv.Itoa(pullNb)+" - "+configType+" - "+string(version),
							configFile,
							ref,
							exec.SourcePullRequest,
							pullNb,
							previousGitRef,
							exec.SourcePullRequestBase,
							s.cronNbRetry,
							configType,
							string(version),
							true,
						))
					}
				}
			}
		}
	}
	s.cronPrepare(compareInfos)
}

func (s *Server) cronTags() {
	configs := s.getConfigFiles()

	releases, err := git.GetLatestVitessReleaseCommitHash(s.getVitessPath())
	if err != nil {
		slog.Error(err)
		return
	}

	// a slice of compareInfo's
	var compareInfos []*CompareInfo

	// We add single executions for the tags, we do not compare them against anything
	for _, release := range releases {
		for configType, configFile := range configs {
			if configType == "micro" {
				compareInfos = append(compareInfos, newSingleExecution("Tag run", configFile, release.CommitHash, exec.SourceTag+release.Name, s.cronNbRetry, configType, ""))
			} else {
				versions := []macrobench.PlannerVersion{macrobench.V3Planner}
				if release.Number[0] >= 10 {
					versions = append(versions, macrobench.Gen4FallbackPlanner)
				}
				for _, version := range versions {
					compareInfos = append(compareInfos, newSingleExecution("Tag run", configFile, release.CommitHash, exec.SourceTag+release.Name, s.cronNbRetry, configType, string(version)))
				}
			}
		}
	}
	s.cronPrepare(compareInfos)
}

func (s *Server) getConfigFiles() map[string]string {
	configs := map[string]string{
		"micro": s.microbenchConfigPath,
		"oltp":  s.macrobenchConfigPathOLTP,
		"tpcc":  s.macrobenchConfigPathTPCC,
	}
	return configs
}

func (s *Server) cronPrepare(compareInfos []*CompareInfo) {
	for _, info := range compareInfos {
		execQueue <- info
		slog.Infof("New Comparison - Name: %s (config: %s, refMain: %s, sourceMain: %s, planner: %s) added to the queue (length: %d)", info.name, info.config, info.execMain.ref, info.execMain.source, info.plannerVersion, len(execQueue))
	}
}

// checkIfExists is used to check if the results for a given configuration exists or not
// wantOnlyOld defines wether we only want a day old results or not.
// wantOnlyOld = true -> check existence for a day older results
// wantOnlyOld = false -> check existence for any time
func (s *Server) checkIfExists(ref, typeOf, plannerVersion, source string, wantOnlyOld bool) (bool, error) {
	if typeOf == "micro" {
		exist, err := exec.Exists(s.dbClient, ref, source, typeOf, exec.StatusFinished, wantOnlyOld)
		if err != nil {
			slog.Error(err)
			return false, err
		}
		return exist, nil
	}
	exist, err := exec.ExistsMacrobenchmark(s.dbClient, ref, source, typeOf, exec.StatusFinished, plannerVersion, wantOnlyOld)
	if err != nil {
		slog.Error(err)
		return false, err
	}
	return exist, nil
}

// checkIfInProgress returns whether there is an execution with a given configuration already in progress for todays cron job
func (s *Server) checkIfInProgress(ref, typeOf, plannerVersion, source string) (bool, error) {
	if typeOf == "micro" {
		exist, err := exec.ExistsStartedToday(s.dbClient, ref, source, typeOf)
		if err != nil {
			slog.Error(err)
			return false, err
		}
		return exist, nil
	}
	exist, err := exec.ExistsMacrobenchmarkStartedToday(s.dbClient, ref, source, typeOf, plannerVersion)
	if err != nil {
		slog.Error(err)
		return false, err
	}
	return exist, nil
}

func (s *Server) cronExecutionQueueWatcher() {
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
	var execStatusMain executionStatus
	var execStatusComp executionStatus
	defer func() {
		if err != nil {
			// Retry after any failure if the counter is above zero.
			if compInfo.retry > 0 {
				compInfo.retry--
				slog.Infof("Retrying Comparison - Name: %s (config: %s, refMain: %s, sourceMain: %s, planner: %s) retries left: %d", compInfo.name, compInfo.config, compInfo.execMain.ref, compInfo.execMain.source, compInfo.plannerVersion, len(execQueue), compInfo.retry)
				s.cronExecution(compInfo)
				return
			}
		}
		mtx.Lock()
		currentCountExec--
		mtx.Unlock()
	}()

	// check and execute the main execution
	execStatusMain, err = s.checkAndExecuteSingle(compInfo.config, compInfo.execMain.source, compInfo.execMain.ref, compInfo.typeOf, compInfo.plannerVersion, compInfo.pullNb)
	if err != nil {
		slog.Errorf("Error while single execution: %v", err)
		return
	}

	// only execute the secondary execution if it is not nil and the source is set
	// when source is nil, we already know that an execution already exists from previous cron jobs
	if compInfo.execComp != nil {
		if compInfo.execComp.source == "" {
			execStatusComp = executionExists
		} else {
			execStatusComp, err = s.checkAndExecuteSingle(compInfo.config, compInfo.execComp.source, compInfo.execComp.ref, compInfo.typeOf, compInfo.plannerVersion, compInfo.pullNb)
			if err != nil {
				slog.Errorf("Error while single execution: %v", err)
				return
			}
		}
	}

	// only try to send a regression message when there is a secondary execution to compare against
	if compInfo.execComp != nil {
		if execStatusMain == executionExists && execStatusComp == executionExists {
			return
		}
		err = s.sendNotificationForRegression(compInfo)
		if err != nil {
			slog.Errorf("Send notification: %v", err)
			return
		}
	}
}

// checkAndExecuteSingle checks whether there already exists a run or if one is in progress.
// It runs the execution if there are no results.
// It returns the executionStatus and an error. Returned values are
// executionExists -> there already exist results from yesterdays cron jobs ; err will be nil
// executionSucceeded -> results have been found today during cron jobs ; err will be nil
// executionFailed -> an error occurred while reading results or execution ; requires not nil err
func (s *Server) checkAndExecuteSingle(config, source, ref, typeOf, plannerVersion string, pullNb int) (executionStatus, error) {
	// First check if an old execution exists or not
	exists, err := s.checkIfExists(ref, typeOf, plannerVersion, source, true)
	if err != nil {
		return executionFailed, err
	}
	if exists {
		return executionExists, nil
	}

	// Now check if an execution is in progress or not
	isStarted, err := s.checkIfInProgress(ref, typeOf, plannerVersion, source)
	if err != nil {
		return executionFailed, err
	}
	if isStarted {
		// Wait for the execution to finish, with a timeout
		timeOut := time.After(2 * time.Hour)

		for {
			select {
			case <-timeOut:
				// return error due to timeout
				return executionFailed, fmt.Errorf("timed out waiting for existing execution to finish for source: %s, ref: %s, typeOf: %s, plannerVersion: %s", source, ref, typeOf, plannerVersion)
			case <-time.After(1 * time.Minute):
				// check again after every minute if execution finished or not, if it did then exit the for loop
				isStarted, err = s.checkIfInProgress(ref, typeOf, plannerVersion, source)
				if err != nil {
					return executionFailed, err
				}
			}
			if !isStarted {
				break
			}
		}
	}

	// check if there are results already. These would only be from this days cron jobs since we already checked for older results earlier
	exists, err = s.checkIfExists(ref, typeOf, plannerVersion, source, false)
	if err != nil {
		return executionFailed, err
	}
	if exists {
		return executionSucceeded, nil
	}

	// try executing given the configuration.
	err = s.executeSingle(config, source, ref, plannerVersion, pullNb)
	if err != nil {
		return executionFailed, err
	}
	return executionSucceeded, nil
}

func (s *Server) executeSingle(config, source, ref, plannerVersion string, pullNb int) (err error) {
	var e *exec.Exec
	defer func() {
		if e != nil {
			errCleanUp := e.CleanUp()
			if errCleanUp != nil {
				slog.Errorf("CleanUp step: %v", errCleanUp)
				if err != nil {
					err = fmt.Errorf("%v: %v", errCleanUp, err)
				} else {
					err = errCleanUp
				}
			}
			e.Success()
			slog.Info("Finished execution: ", e.UUID.String())
		}
	}()

	e, err = exec.NewExecWithConfig(config)
	if err != nil {
		slog.Warn(err.Error())
		return err
	}
	slog.Info("Created new execution: ", e.UUID.String())
	e.Source = source
	e.GitRef = ref
	e.VtgatePlannerVersion = plannerVersion
	e.PullNB = pullNb

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

	err = e.ExecuteWithTimeout(time.Hour * 2)
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
		pullNb:              pullNB,
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
		name:                name,
	}
}
