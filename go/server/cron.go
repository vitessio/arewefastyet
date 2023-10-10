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

	"github.com/robfig/cron/v3"
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

func createIndividualDaily(schedule string, job func()) error {
	if schedule == "" {
		return nil
	}

	c := cron.New()
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
		err := createIndividualDaily(c.schedule, c.f)
		if err != nil {
			return err
		}
	}
	go s.dailyExecutionQueueWatcher()
	return nil
}

func (s *Server) getConfigFiles() map[string]benchmarkConfig {
	return s.benchmarkConfig
}

func (s *Server) removePRFromQueue(element *executionQueueElement) {
	mtx.Lock()
	defer mtx.Unlock()

	for id, e := range queue {
		if e.Executing == false && id.PullNb == element.identifier.PullNb && id.BenchmarkType == element.identifier.BenchmarkType && id.Source == element.identifier.Source && id.GitRef != element.identifier.GitRef {
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

	if len(s.sourceFilter) > 0 && !slices.Contains(s.sourceFilter, element.identifier.Source) {
		return
	}
	if len(s.excludeSourceFilter) > 0 && slices.Contains(s.excludeSourceFilter, element.identifier.Source) {
		return
	}

	_, found := queue[element.identifier]

	if found {
		return
	}
	exists, err := s.checkIfExecutionExists(element.identifier)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if !exists {
		queue[element.identifier] = element
		slog.Infof("%+v is added to the queue", element.identifier)

		// we sleep here to avoid adding too many similar elements to the queue at the same time.
		time.Sleep(2 * time.Second)
	}
}
