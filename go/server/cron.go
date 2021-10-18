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
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type (
	executionQueueElement struct {
		name, config                       string
		retry                              int
		identifier                         executionIdentifier
		compareWith                        []executionIdentifier
		notifyAlways, done, primary bool
	}

	executionIdentifier struct {
		uuid, gitRef, source, benchmarkType, plannerVersion string
		pullNb                                              int
	}

	executionQueue map[executionIdentifier]*executionQueueElement

	// CompareInfo has the details required to compare two git commits
	CompareInfo struct {
		// config is the configuration file to use
		config string

		// execInfo contains the details of the execution of the main commit
		execInfo *execInfo

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
	execQueue = make(chan *CompareInfo)
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
	str := fmt.Sprintf("Identifying: %+v\t", element.identifier)

	mtx.Lock()
	defer func() {
		log.Println(str)
		mtx.Unlock()
	}()

	_, found := queue[element.identifier]

	if found {
		str = fmt.Sprintf("%s WAS FOUND", str)
		return
	}
	exists, err := s.checkIfExecutionExists(element.identifier)
	if err != nil {
		str = fmt.Sprintf("%s GOT ERROR, stop", str)
		slog.Error(err.Error())
		return
	}
	if !exists {
		// we sleep here to avoid adding too many similar elements to the queue at the same time.
		time.Sleep(1 * time.Second)

		queue[element.identifier] = element
		str = fmt.Sprintf("%s IS ADDED TO THE QUEUE", str)
	} else {
		str = fmt.Sprintf("%s ALREADY EXISTS IN DATABASE", str)
	}
}

func (s *Server) addCompare(compareInfos []*CompareInfo) {
	for _, info := range compareInfos {
		execQueue <- info
		slog.Infof("New Comparison - Name: %s (config: %s, refMain: %s, sourceMain: %s, planner: %s) added to the queue (length: %d)", info.name, info.config, info.execInfo.ref, info.execInfo.source, info.plannerVersion, len(execQueue))
	}
}

func newCompareInfo(name, configFile, ref, source string, pullNB int, compareRef, compareSource string, retry int, configType, plannerVersion string, ignoreNonRegression bool) *CompareInfo {
	return &CompareInfo{
		config: configFile,
		execInfo: &execInfo{
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
		execInfo: &execInfo{
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
