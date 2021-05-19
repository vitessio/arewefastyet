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
)

func (s *Server) createNewCron() error {
	if s.cronSchedule == "" {
		return nil
	}
	c := cron.New()

	jobs := []func(){
		s.cronBranchHandler,
		s.cronPRLabels,
	}
	for _, job := range jobs {
		_, err := c.AddFunc(s.cronSchedule, job)
		if err != nil {
			return err
		}
	}
	c.Start()
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
	s.cronExecution(configFiles, ref, "cron", s.cronNbRetry)

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
		s.cronExecution(configFiles, ref, "cron_release", s.cronNbRetry)
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
		s.cronExecution(configArray, gitRef, "cron_pr", s.cronNbRetry)
	}
}

func (s *Server) cronExecution(configs []string, ref, source string, retry int) {
	for _, config := range configs {
		e, err := exec.NewExecWithConfig(config)
		if err != nil {
			slog.Warn(err.Error())
			return
		}
		slog.Info("Created new execution: ", e.UUID.String())
		e.Source = source
		e.GitRef = ref

		config := config
		go func() {
			defer func() {
				err = e.CleanUp()
				if err != nil {
					slog.Errorf("CleanUp step: %v", err)
				}
				if err != nil {
					// Retry after execution failure if the counter is above zero.
					if retry > 0 {
						slog.Info("Retrying execution for source: ", source, " git ref: ", ref, " with configuration: ", config, ". Number of retry left: ", retry-1)
						s.cronExecution([]string{config}, ref, source, retry-1)
					}
				}
			}()

			slog.Info("Started execution: ", e.UUID.String(), ", with git ref: ", ref)
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
		}()
	}
}
