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
	"github.com/vitessio/arewefastyet/go/exec"
	"time"
)

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

func (s *Server) cronExecution(compInfo *CompareInfo) {
	var err error
	var execStatusMain executionStatus
	var execStatusComp executionStatus
	defer func() {
		if err != nil {
			// Retry after any failure if the counter is above zero.
			if compInfo.retry > 0 {
				compInfo.retry--
				slog.Infof("Retrying Comparison - Name: %s (config: %s, refMain: %s, sourceMain: %s, planner: %s) retries left: %d", compInfo.name, compInfo.config, compInfo.execInfo.ref, compInfo.execInfo.source, compInfo.plannerVersion, len(execQueue), compInfo.retry)
				s.cronExecution(compInfo)
				return
			}
		}
		mtx.Lock()
		currentCountExec--
		mtx.Unlock()
	}()

	// check and execute the main execution
	execStatusMain, err = s.checkAndExecuteSingle(compInfo.config, compInfo.execInfo.source, compInfo.execInfo.ref, compInfo.typeOf, compInfo.plannerVersion, compInfo.pullNb)
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

func (s *Server) checkIfExecutionExists(identifier executionIdentifier) (bool, error) {
	checkStatus := []string{exec.StatusStarted, exec.StatusFinished, exec.StatusCreated}
	for _, status := range checkStatus {
		var exists bool
		var err error
		if identifier.benchmarkType == "micro" {
			exists, err = exec.Exists(s.dbClient, identifier.gitRef, identifier.source, identifier.benchmarkType, status, true)
		} else {
			exists, err = exec.ExistsMacrobenchmark(s.dbClient, identifier.gitRef, identifier.source, identifier.benchmarkType, status, identifier.plannerVersion, true)
		}
		if err != nil {
			slog.Error(err)
			return false, err
		}
		if exists {
			return true, nil
		}
	}
	return false, nil
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
