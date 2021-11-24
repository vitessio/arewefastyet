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
)

type (
	executionQueueElement struct {
		config                  string
		retry                   int
		identifier              executionIdentifier
		compareWith             []executionIdentifier
		notifyAlways, executing bool
	}

	executionIdentifier struct {
		GitRef, Source, BenchmarkType, PlannerVersion string
		PullNb                                        int
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
	queue = make(executionQueue)

	err := createIndividualCron(s.cronSchedule, []func(){
		s.branchCronHandler,
		s.tagsCronHandler,
	})
	if err != nil {
		return err
	}
	err = createIndividualCron(s.cronSchedulePullRequests, []func(){s.pullRequestsCronHandler})
	if err != nil {
		return err
	}

	go s.cronExecutionQueueWatcher()
	return nil
}

func (s *Server) getConfigFiles() map[string]string {
	configs := map[string]string{
		"micro": s.microbenchConfigPath,
		"oltp":  s.macrobenchConfigPathOLTP,
		"tpcc":  s.macrobenchConfigPathTPCC,
	}
	return configs
}

func (s *Server) addToQueue(element *executionQueueElement) {
	if element.identifier.BenchmarkType == "micro" {
		return
	}

	mtx.Lock()
	defer func() {
		mtx.Unlock()
	}()

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
