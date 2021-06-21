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

func (s *Server) sendNotificationForRegression(compInfo *CompareInfo) (err error) {
	// regression header, appender to header in the event of a regression
	regressionHeader := `*Observed a regression.*
`

	// header of the message, before the regression explanation
	header := compInfo.name + "\n\n"
	if compInfo.execMain.pullNB > 0 {
		header += fmt.Sprintf(`Benchmarked PR #<https://github.com/vitessio/vitess/pull/%d>.`, compInfo.execMain.pullNB)
	} else {
		header += `Comparing: recent commit <https://github.com/vitessio/vitess/commit/` + compInfo.execMain.ref + `|` + git.ShortenSHA(compInfo.execMain.ref) + `> with old commit <https://github.com/vitessio/vitess/commit/` + compInfo.execComp.ref + `|` + git.ShortenSHA(compInfo.execComp.ref) + `>.`
	}
	header += `Comparison can be seen at : ` + getComparisonLink(compInfo.execMain.ref, compInfo.execComp.ref) + `

`

	if compInfo.typeOf == "micro" {
		microBenchmarks, err := microbench.Compare(s.dbClient, compInfo.execMain.ref, compInfo.execComp.ref)
		if err != nil {
			return err
		}
		regression := microBenchmarks.Regression()
		err = s.sendMessageIfRegression(compInfo.ignoreNonRegression, regression, header, regressionHeader)
		if err != nil {
			return err
		}
	} else if compInfo.typeOf == "oltp" || compInfo.typeOf == "tpcc" {
		macrosMatrices, err := macrobench.CompareMacroBenchmarks(s.dbClient, compInfo.execMain.ref, compInfo.execComp.ref, macrobench.PlannerVersion(compInfo.plannerVersion))
		if err != nil {
			return err
		}

		macroResults := macrosMatrices[macrobench.Type(compInfo.typeOf)].(macrobench.ComparisonArray)
		if len(macroResults) == 0 {
			return fmt.Errorf("no macrobenchmark result")
		}

		regression := macroResults[0].Regression()
		err = s.sendMessageIfRegression(compInfo.ignoreNonRegression, regression, header, regressionHeader)
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
