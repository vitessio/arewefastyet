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

	_, err := c.AddFunc(s.cronSchedule, s.cronMasterHandler)
	if err != nil {
		return err
	}
	c.Start()
	return nil
}

func (s *Server) cronMasterHandler() {
	configs := []string{
		s.microbenchConfigPath,
		s.macrobenchConfigPathOLTP,
		s.macrobenchConfigPathTPCC,
	}

	for _, config := range configs {
		e, err := exec.NewExecWithConfig(config)
		if err != nil {
			slog.Warn(err.Error())
			return
		}
		slog.Info("Created new execution: ", e.UUID.String())

		ref, err := git.GetLatestVitessCommitHash()
		if err != nil {
			slog.Warn(err.Error())
			return
		}

		e.Source = "cron"
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

			err = e.CleanUp()
			if err != nil {
				slog.Error("Clean step", err.Error())
				return
			}
			slog.Info("Finished execution: ", e.UUID.String())
		}()
	}
}
