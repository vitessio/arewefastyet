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

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/vitessio/arewefastyet/go/exec"
	"github.com/vitessio/arewefastyet/go/tools/git"
	"golang.org/x/exp/slices"
)

type (
	executionQueueElement struct {
		config                  benchmarkConfig
		retry                   int
		identifier              executionIdentifier
		compareWith             []executionIdentifier
		notifyAlways, Executing bool
	}

	executionIdentifier struct {
		GitRef, Source, BenchmarkType, PlannerVersion string
		PullNb                                        int
		PullBaseRef                                   string
		Version                                       git.Version
		UUID                                          string
	}

	executionQueue map[executionIdentifier]*executionQueueElement
)

const (
	// maxConcurJob is the maximum number of concurrent jobs that we can execute
	maxConcurJob = 1
)

var (
	currentCountExec int
	mtx              sync.RWMutex
	queue            executionQueue
)

func (ei executionIdentifier) equalWithoutUUID(id executionIdentifier) bool {
	ei.UUID = ""
	id.UUID = ""
	return ei == id
}

func createIndividualCRON(schedule string, job func()) error {
	if schedule == "" {
		return nil
	}

	c := cron.New(cron.WithLogger(cron.DefaultLogger))
	_, err := c.AddFunc(schedule, job)
	if err != nil {
		return err
	}
	c.Start()
	return nil
}

func (s *Server) createCrons() error {
	queue = make(executionQueue)

	crons := []struct {
		schedule string
		f        func()
		name     string
	}{
		{name: "branch", schedule: s.cronSchedule, f: s.branchCronHandler},
		{name: "pull_requests", schedule: s.cronSchedulePullRequests, f: s.pullRequestsCronHandler},
		{name: "tags", schedule: s.cronScheduleTags, f: s.tagsCronHandler},
	}
	for _, c := range crons {
		if c.schedule == "none" {
			continue
		}
		slog.Info("Starting the CRON ", c.name, " with schedule: ", c.schedule)
		err := createIndividualCRON(c.schedule, c.f)
		if err != nil {
			return err
		}

		// Trigger CRONs upon creation of the server
		go c.f()
	}
	go s.cronExecutionQueueWatcher()
	return nil
}

func (s *Server) getConfigFiles() map[string]benchmarkConfig {
	return s.benchmarkConfig
}

func (s *Server) removePRFromQueue(element *executionQueueElement) {
	mtx.Lock()
	defer mtx.Unlock()

	for id, e := range queue {
		if !e.Executing && id.PullNb == element.identifier.PullNb && id.BenchmarkType == element.identifier.BenchmarkType && id.Source == element.identifier.Source && id.GitRef != element.identifier.GitRef {
			slog.Infof("%+v is removed from the queue", id)
			delete(queue, id)
		}
	}
}

func (s *Server) addToQueue(element *executionQueueElement) {
	mtx.Lock()
	defer func() {
		mtx.Unlock()
	}()

	// Check if the benchmark we are trying to add is part of exclusion rules
	if len(s.sourceFilter) > 0 && !slices.Contains(s.sourceFilter, element.identifier.Source) {
		return
	}
	if len(s.excludeSourceFilter) > 0 && slices.Contains(s.excludeSourceFilter, element.identifier.Source) {
		return
	}

	// Duplication mechanism to multiply the execution queue element depending
	// on how many times we want to execute the same benchmark and on how many
	// times it already exists in the database.
	var execElements []*executionQueueElement
	if element.identifier.BenchmarkType == "micro" {
		execElements = append(execElements, element)
	} else {
		nb, err := s.getNumberOfBenchmarksInDB(element.identifier)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		countInQueue := 0
		for identifier, _ := range queue {
			if identifier.equalWithoutUUID(element.identifier) {
				countInQueue++
			}
		}

		multiplyFactor := exec.MaximumBenchmarkWithSameConfig - nb - countInQueue
		if multiplyFactor <= 0 {
			slog.Infof("not adding %+v to the queue, already full", element.identifier)
			return
		}
		for i := 0; i < multiplyFactor; i++ {
			newElement := *element
			newElement.identifier.UUID = uuid.NewString()
			execElements = append(execElements, &newElement)
		}
	}

	// Add all the elements to the queue
	for _, execElement := range execElements {
		// Check if the exact same benchmark is already in the queue
		_, found := queue[execElement.identifier]
		if found {
			slog.Infof("not adding %+v, already in the queue", execElement.identifier)
			return
		}

		queue[execElement.identifier] = execElement
		slog.Infof("%+v is added to the queue", execElement.identifier)

		// We sleep here to avoid adding too many similar elements to the queue at the same time.
		time.Sleep(100 * time.Millisecond)
	}
}
