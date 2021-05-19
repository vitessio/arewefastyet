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
		s.cronMasterHandler,
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

func (s *Server) cronMasterHandler() {
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

	s.cronExecution(configFiles, ref, "cron")
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
		s.cronExecution(configArray, gitRef, "cron_pr")
	}
}

func (s *Server) cronExecution(configs []string, ref, source string) {
	for _, config := range configs {
		e, err := exec.NewExecWithConfig(config)
		if err != nil {
			slog.Warn(err.Error())
			return
		}
		slog.Info("Created new execution: ", e.UUID.String())
		e.Source = source
		e.GitRef = ref

		go func() {
			slog.Info("Started execution: ", e.UUID.String(), ", with git ref: ", ref)
			err = e.Prepare()
			if err != nil {
				slog.Error("Prepare step", err.Error())
				return
			}

			err = e.SetOutputToDefaultPath()
			if err != nil {
				slog.Error("Prepare outputs step", err.Error())
				return
			}

			err = e.Execute()
			if err != nil {
				slog.Error("Execution step", err.Error())
				return
			}

			err = e.SendNotificationForRegression()
			if err != nil {
				slog.Error("Send notification", err.Error())
				return
			}

			err = e.CleanUp()
			if err != nil {
				slog.Error("Clean step", err.Error())
				return
			}
			slog.Info("Finished execution: ", e.UUID.String())
		}()
	}
}
