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
	"time"

	"github.com/google/uuid"
	"github.com/vitessio/arewefastyet/go/exec"
)

func (s *Server) executeSingle(config benchmarkConfig, identifier executionIdentifier, nextIsSame, lastIsSame bool) (err error) {
	var e *exec.Exec
	defer func() {
		if e != nil {
			if err != nil {
				err = fmt.Errorf("%v", err)
			}
			if errSuccess := e.Success(); errSuccess != nil {
				err = errSuccess
				return
			}
			slog.Info("Finished execution: UUID: [", e.UUID.String(), "], Git Ref: [", identifier.GitRef, "], Workload: [", identifier.Workload, "]")
		}
	}()

	e, err = exec.NewExecWithConfig(config.file, identifier.UUID)

	if err != nil {
		nErr := fmt.Errorf("new exec error: %v", err)
		slog.Error(nErr.Error())
		return nErr
	}
	e.Source = identifier.Source
	e.GitRef = identifier.GitRef
	e.VtgatePlannerVersion = identifier.PlannerVersion
	e.PullNB = identifier.PullNb
	e.VitessVersion = identifier.Version
	e.NextBenchmarkIsTheSame = nextIsSame
	e.ProfileInformation = identifier.Profile
	e.RepoDir = s.getVitessPath()

	// Check if the previous benchmark is the same and if it is
	// safe to execute this new benchmark without a preparatory cleanup phase.
	e.PreviousBenchmarkIsTheSame = lastIsSame
	if lastIsSame {
		lastBenchmarkWasClean, err := exec.IsLastExecutionFinished(s.dbClient)
		if err != nil {
			return err
		}
		if !lastBenchmarkWasClean {
			e.PreviousBenchmarkIsTheSame = false
		}
	}

	slog.Info("Starting execution: UUID: [", e.UUID.String(), "], Git Ref: [", identifier.GitRef, "], Workload: [", identifier.Workload, "]")
	err = e.Prepare()
	if err != nil {
		nErr := fmt.Errorf("prepare error: %v", err)
		slog.Error(nErr.Error())
		return nErr
	}

	err = e.SetOutputToDefaultPath()
	if err != nil {
		nErr := fmt.Errorf("prepare output error: %v", err)
		slog.Error(nErr.Error())
		return nErr
	}

	timeout := 1 * time.Hour
	if identifier.Workload == "micro" {
		timeout = 4 * time.Hour
	}
	err = e.ExecuteWithTimeout(timeout)
	if err != nil {
		nErr := fmt.Errorf("execute with timeout error: %v", err)
		slog.Error(nErr.Error())
		return nErr
	}
	return nil
}

func (s *Server) executeElement(element *executionQueueElement, nextIsSame bool, lastIsSame bool) {
	if element.retry < 0 {
		if _, found := queue[element.identifier]; found {
			// removing the element from the queue since we are done with it
			mtx.Lock()
			delete(queue, element.identifier)
			mtx.Unlock()
		}
		decrementNumberOfOnGoingExecution()
		return
	}

	// execute with the given configuration file and exec identifier
	err := s.executeSingle(element.config, element.identifier, nextIsSame, lastIsSame)
	if err != nil {
		slog.Error(err.Error())

		// execution failed, we retry
		element.retry -= 1
		element.identifier.UUID = uuid.NewString()

		// Here we set lastIsSame as false since the previous benchmark has failed
		// That allows us to avoid executing one more database request to check if
		// the previous benchmark was successful or not.
		s.executeElement(element, nextIsSame, false)
		return
	}

	go func() {
		// removing the element from the queue since we are done with it
		mtx.Lock()
		delete(queue, element.identifier)
		mtx.Unlock()

		// we will wait for the benchmarks we need to compare it against and notify users if needed
		s.compareElement(element)
	}()

	decrementNumberOfOnGoingExecution()
}

func (s *Server) compareElement(element *executionQueueElement) {
	// map that contains all the comparison we saw and analyzed
	seen := map[executionIdentifier]bool{}
	done := 0
	for done != len(element.compareWith) {
		time.Sleep(1 * time.Second)
		for _, comparer := range element.compareWith {
			// checking if we have already seen this comparison, if we did, we can skip it.
			if _, ok := seen[comparer]; ok {
				continue
			}
			comparerUUID, err := exec.GetFinishedExecution(s.dbClient, comparer.GitRef, comparer.Source, comparer.Workload, comparer.PlannerVersion, comparer.PullNb)
			if err != nil {
				slog.Error(err)
				return
			}
			if comparerUUID != "" {
				seen[comparer] = true
				done++
			}
		}
	}
}

func (s *Server) getNumberOfBenchmarksInDB(identifier executionIdentifier) (int, error) {
	var nb int
	var err error
	if identifier.Workload == "micro" {
		var exists bool
		exists, err = exec.Exists(s.dbClient, identifier.GitRef, identifier.Source, identifier.Workload, exec.StatusFinished)
		if exists {
			nb = 1
		}
	} else {
		nb, err = exec.CountMacroBenchmark(s.dbClient, identifier.GitRef, identifier.Source, identifier.Workload, exec.StatusFinished, identifier.PlannerVersion)
	}
	if err != nil {
		slog.Error(err)
		return 0, err
	}
	return nb, nil
}

// cronExecutionQueueWatcher runs an infinite loop that watches the execution queue
// it will send an item to the Executor based on different priority rules ordered that way:
//  0. No executions that are in progress will get executed
//  1. Admin executions always get executed first no matter what
//  2. Execution of the same type (workload/commit) will be executed sequentially
//  3. If none of this priority match, a random element is picked
func (s *Server) cronExecutionQueueWatcher() {
	var lastExecutedId executionIdentifier
	queueWatch := func() {
		mtx.Lock()
		defer mtx.Unlock()
		if currentCountExec >= maxConcurJob || len(queue) == 0 {
			return
		}

		// Look for what's in the queue, specifically here is what we look for:
		// 	- An execution that matches the previous execution
		// 	- The first admin execution
		var lastBenchmarkIsTheSame bool
		var nextExecuteElement *executionQueueElement
		var firstAdminExecuteElement *executionQueueElement
		for _, element := range queue {
			if element.Executing {
				continue
			}

			// Look if there is a matching execution in the queue
			if nextExecuteElement == nil && element.identifier.equalWithoutUUID(lastExecutedId) {
				nextExecuteElement = element
				lastBenchmarkIsTheSame = true
			}

			// Look if there is any admin execution in the queue
			if firstAdminExecuteElement == nil && element.identifier.Source == "admin" {
				firstAdminExecuteElement = element
			}

			// If both the first admin and next execution are not nil, we found what we want, we can stop
			if nextExecuteElement != nil && firstAdminExecuteElement != nil {
				break
			}
		}

		// If we found an admin execution and the next execution we found is either nil or not an admin
		// we force execute the first admin execution we found.
		if firstAdminExecuteElement != nil && (nextExecuteElement == nil || nextExecuteElement.identifier.Source != "admin") {
			nextExecuteElement = firstAdminExecuteElement
			lastBenchmarkIsTheSame = false
		}

		// If we did not find any matching element just go to the first one which is not executing
		if firstAdminExecuteElement == nil && nextExecuteElement == nil {
			for _, element := range queue {
				if !element.Executing {
					nextExecuteElement = element
					break
				}
			}
		}

		if nextExecuteElement != nil {
			// Find out if there is another element in queue that match the one we want to execute
			var nextBenchmarkIsTheSame bool
			for _, element := range queue {
				if element.Executing {
					continue
				}
				if nextExecuteElement.identifier.UUID != element.identifier.UUID && element.identifier.equalWithoutUUID(nextExecuteElement.identifier) {
					nextBenchmarkIsTheSame = true
					break
				}
			}

			// Execute the element if found
			currentCountExec++
			lastExecutedId = nextExecuteElement.identifier

			// setting this element to `Executing = true`, so we do not execute it twice in the future
			nextExecuteElement.Executing = true
			go s.executeElement(nextExecuteElement, nextBenchmarkIsTheSame, lastBenchmarkIsTheSame)
			return
		}
	}

	for {
		queueWatch()
	}
}

func decrementNumberOfOnGoingExecution() {
	mtx.Lock()
	currentCountExec--
	mtx.Unlock()
}
