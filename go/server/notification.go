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

	"github.com/vitessio/arewefastyet/go/slack"

	"github.com/vitessio/arewefastyet/go/tools/git"
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
	"github.com/vitessio/arewefastyet/go/tools/microbench"
)

func (s *Server) sendNotificationForRegression(leftSource, rightSource, leftRef, rightRef, plannerVersion, benchmarkType string, pullNb int, notifyAlways bool) (err error) {
	// regression header, appender to header in the event of a regression
	regressionHeader := `*Observed a regression.*
`

	// header of the message, before the regression explanation
	header := fmt.Sprintf("Comparing %s with %s, with the %s benchmark", leftSource, rightSource, benchmarkType)
	if benchmarkType != "micro" {
		header += fmt.Sprintf(" using the %s query planner", plannerVersion)
	}
	header += "\n\n"
	if pullNb > 0 {
		header += fmt.Sprintf(`Benchmarked PR #<https://github.com/vitessio/vitess/pull/%d>. `, pullNb)
	} else {
		header += `Comparing: recent commit <https://github.com/vitessio/vitess/commit/` + leftRef + `|` + git.ShortenSHA(leftRef) + `> with old commit <https://github.com/vitessio/vitess/commit/` + rightRef + `|` + git.ShortenSHA(rightRef) + `>. `
	}
	header += `Comparison can be seen at : ` + getComparisonLink(leftRef, rightRef) + `

`

	if benchmarkType == "micro" {
		microBenchmarks, err := microbench.Compare(s.dbClient, leftRef, rightRef)
		if err != nil {
			return err
		}
		regression := microBenchmarks.Regression()
		err = s.sendMessageIfRegression(notifyAlways, regression, header, regressionHeader)
		if err != nil {
			return err
		}
	} else {
		macrosMatrices, err := macrobench.CompareMacroBenchmarks(s.dbClient, leftRef, rightRef, macrobench.PlannerVersion(plannerVersion), s.benchmarkTypes)
		if err != nil {
			return err
		}

		macroResults := macrosMatrices[benchmarkType]
		if len(macroResults) == 0 {
			return fmt.Errorf("no macrobenchmark result")
		}

		regression := macroResults[0].Regression()
		err = s.sendMessageIfRegression(notifyAlways, regression, header, regressionHeader)
		if err != nil {
			return err
		}
	}
	return nil
}

func getComparisonLink(leftSHA, rightSHA string) string {
	return "https://benchmark.vitess.io/compare?r=" + leftSHA + "&c=" + rightSHA
}

func (s *Server) sendMessageIfRegression(ignoreNonRegression bool, regression, header, regressionHeader string) error {
	if regression != "" || ignoreNonRegression {
		hd := header
		if regression != "" {
			hd = regressionHeader + header
		}
		err := s.sendSlackMessage(regression, hd)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) sendSlackMessage(regression, header string) error {
	content := header + regression
	msg := slack.TextMessage{Content: content}
	err := msg.Send(s.slackConfig)
	if err != nil {
		return err
	}
	return nil
}
