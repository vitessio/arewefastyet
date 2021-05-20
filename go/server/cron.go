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
	"github.com/robfig/cron/v3"
	"github.com/vitessio/arewefastyet/go/exec"
	"github.com/vitessio/arewefastyet/go/tools/git"
	"sync"
	"time"
)

type execInfo struct {
	config string
	ref    string
	retry  int
	source string
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
	var configFiles []string
	for configType, configFile := range configs {
		exist, err := exec.Exists(s.dbClient, ref, "cron", configType, exec.StatusFinished)
		if err != nil {
			slog.Error(err)
			continue
		}
		if !exist {
			configFiles = append(configFiles, configFile)
		}
	}
	s.cronPrepare(configFiles, ref, "cron")

	// release branches
	releases, err := git.GetLatestVitessReleaseBranchCommitHash(s.getVitessPath())
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	for _, release := range releases {
		ref = release.CommitHash
		configFiles = nil
		for configType, configFile := range configs {
			exist, err := exec.Exists(s.dbClient, ref, "cron_release", configType, exec.StatusFinished)
			if err != nil {
				slog.Error(err)
				continue
			}
			if !exist {
				configFiles = append(configFiles, configFile)
			}
		}
		s.cronPrepare(configFiles, ref, "cron_release")
	}
}

func (s Server) cronPRLabels() {
	configs := map[string]string{
		"micro": s.microbenchConfigPath,
		"oltp":  s.macrobenchConfigPathOLTP,
		"tpcc":  s.macrobenchConfigPathTPCC,
	}

	SHAs, err := git.GetPullRequestHeadForLabels([]string{s.prLabelTrigger}, "vitessio/vitess")
	if err != nil {
		slog.Error(err)
		return
	}

	// map with the git_ref as key and a slice of configuration file as value
	toExec := map[string][]string{}
	for _, sha := range SHAs {
		for configType, configFile := range configs {
			exist, err := exec.Exists(s.dbClient, sha, "cron_pr", configType, exec.StatusFinished)
			if err != nil {
				slog.Error(err)
				continue
			}
			if !exist {
				toExec[sha] = append(toExec[sha], configFile)
			}
		}
	}
	for gitRef, configArray := range toExec {
		s.cronPrepare(configArray, gitRef, "cron_pr")
	}
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

	// map with the git_ref as key and a slice of configuration file as value
	toExec := map[string][]string{}
	source := "cron_tags"
	for _, release := range releases {
		for configType, configFile := range configs {
			exist, err := exec.Exists(s.dbClient, release.CommitHash, source, configType, exec.StatusFinished)
			if err != nil {
				slog.Error(err)
				continue
			}
			if !exist {
				toExec[release.CommitHash] = append(toExec[release.CommitHash], configFile)
			}
		}
	}
	for gitRef, configArray := range toExec {
		s.cronPrepare(configArray, gitRef, source)
	}
}

func (s *Server) cronPrepare(configs []string, ref, source string) {
	for _, config := range configs {
		info := execInfo{
			config: config,
			ref:    ref,
			retry:  s.cronNbRetry,
			source: source,
		}
		execQueue <- info
		slog.Infof("New Execution (config: %s, ref: %s, source: %s) added to the queue (length: %d)", info.config, info.ref, info.source, len(execQueue))
	}
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
