/*
 *
 * Copyright 2023 The Vitess Authors.
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
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vitessio/arewefastyet/go/exec"
	"github.com/vitessio/arewefastyet/go/exec/metrics"
	"github.com/vitessio/arewefastyet/go/tools/git"
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
	"github.com/vitessio/arewefastyet/go/tools/microbench"
)

type ErrorAPI struct {
	Error string `json:"error"`
}

type RecentExecutions struct {
	UUID          string     `json:"uuid"`
	Source        string     `json:"source"`
	GitRef        string     `json:"git_ref"`
	Status        string     `json:"status"`
	TypeOf        string     `json:"type_of"`
	PullNb        int        `json:"pull_nb"`
	GolangVersion string     `json:"golang_version"`
	StartedAt     *time.Time `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"`
}

func (s *Server) getRecentExecutions(c *gin.Context) {
	execs, err := exec.GetRecentExecutions(s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	recentExecs := make([]RecentExecutions, 0, len(execs))
	for _, e := range execs {
		recentExecs = append(recentExecs, RecentExecutions{
			UUID:          e.UUID.String(),
			Source:        e.Source,
			GitRef:        e.GitRef,
			Status:        e.Status,
			TypeOf:        e.TypeOf,
			PullNb:        e.PullNB,
			GolangVersion: e.GolangVersion,
			StartedAt:     e.StartedAt,
			FinishedAt:    e.FinishedAt,
		})
	}
	c.JSON(http.StatusOK, recentExecs)
}

func (s *Server) getExecutionsQueue(c *gin.Context) {
	recentExecs := make([]RecentExecutions, 0, len(queue))
	for _, e := range queue {
		if e.Executing {
			continue
		}
		recentExecs = append(recentExecs, RecentExecutions{
			Source: e.Identifier.Source,
			GitRef: e.Identifier.GitRef,
			TypeOf: e.Identifier.BenchmarkType,
			PullNb: e.Identifier.PullNb,
		})
	}
	sort.Slice(recentExecs, func(i, j int) bool {
		return recentExecs[i].GitRef > recentExecs[j].GitRef && recentExecs[i].Source > recentExecs[j].Source
	})
	c.JSON(http.StatusOK, recentExecs)
}

type VitessGitRef struct {
	UUID          string     `json:"uuid"`
	Source        string     `json:"source"`
	GitRef        string     `json:"git_ref"`
	Status        string     `json:"status"`
	TypeOf        string     `json:"type_of"`
	PullNb        int        `json:"pull_nb"`
	GolangVersion string     `json:"golang_version"`
	StartedAt     *time.Time `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"`
}

func (s *Server) getLatestVitessGitRef(c *gin.Context) {
	allReleases, err := git.GetLatestVitessReleaseCommitHash(s.getVitessPath())
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	lastrunCronSHA, err := exec.GetLatestCronJobForMacrobenchmarks(s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	mainRelease := []*git.Release{{
		Name:       "main",
		CommitHash: lastrunCronSHA,
	}}
	allReleases = append(mainRelease, allReleases...)
	// get all the latest release branches as well
	allReleaseBranches, err := git.GetLatestVitessReleaseBranchCommitHash(s.getVitessPath())
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	allReleases = append(allReleases, allReleaseBranches...)

	c.JSON(http.StatusOK, allReleases)
}

type CompareMacrobench struct {
	Type string                `json:"type"`
	Diff macrobench.Comparison `json:"diff"`
}

func (s *Server) compareMacrobenchmarks(c *gin.Context) {
	rightSHA := c.Query("rtag")
	leftSHA := c.Query("ltag")

	// Compare Macrobenchmarks for the two given SHAs.
	macrosMatrices, err := macrobench.CompareMacroBenchmarks(s.dbClient, rightSHA, leftSHA, macrobench.Gen4Planner, s.benchmarkTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}

	cmpMacros := make([]CompareMacrobench, 0, len(macrosMatrices))
	for typeof, cmp := range macrosMatrices {
		cmpMacro := CompareMacrobench{
			Type: typeof,
		}
		if len(cmp) > 0 {
			cmpMacro.Diff = cmp[0]
		}
		cmpMacros = append(cmpMacros, cmpMacro)
	}

	sort.Slice(cmpMacros, func(i, j int) bool {
		return cmpMacros[i].Type < cmpMacros[j].Type
	})

	c.JSON(http.StatusOK, cmpMacros)
}

func (s *Server) compareMicrobenchmarks(c *gin.Context) {
	leftSHA := c.Query("ltag")
	rightSHA := c.Query("rtag")

	// Get the results from the SHAs
	leftMbd, err := microbench.GetResultsForGitRef(leftSHA, s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	leftMbd = leftMbd.ReduceSimpleMedianByName()
	rightMbd, err := microbench.GetResultsForGitRef(rightSHA, s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	rightMbd = rightMbd.ReduceSimpleMedianByName()

	matrix := microbench.MergeDetails(rightMbd, leftMbd)
	c.JSON(http.StatusOK, matrix)
}

type searchResult struct {
	Macros map[string]macrobench.DetailsArray
	Micro  microbench.DetailsArray
}

func (s *Server) searchBenchmarck(c *gin.Context) {
	gitRef := c.Query("git_ref")

	macros, err := macrobench.GetDetailsArraysFromAllTypes(gitRef, macrobench.Gen4Planner, s.dbClient, s.benchmarkTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}

	micro, err := microbench.GetResultsForGitRef(gitRef, s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	micro = micro.ReduceSimpleMedianByName()

	var res searchResult
	res.Macros = macros
	res.Micro = micro

	c.JSON(http.StatusOK, res)
}

func (s *Server) queriesCompareMacrobenchmarks(c *gin.Context) {
	leftGitRef := c.Query("ltag")
	rightGitRef := c.Query("rtag")
	macroType := macrobench.Type(c.Query("type"))

	if leftGitRef == "" || rightGitRef == "" || macroType == "" {
		c.JSON(http.StatusBadRequest, &ErrorAPI{Error: "The gitref left and right and where the macrotype are incorrect are missing. Please kindly add them."})
		return
	}

	plansLeft, err := macrobench.GetVTGateSelectQueryPlansWithFilter(leftGitRef, macroType, macrobench.Gen4Planner, s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	plansRight, err := macrobench.GetVTGateSelectQueryPlansWithFilter(rightGitRef, macroType, macrobench.Gen4Planner, s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	comparison := macrobench.CompareVTGateQueryPlans(plansLeft, plansRight)
	c.JSON(http.StatusOK, comparison)
}

type cronSingleSummary struct {
	Name string
	Data []macrobench.CronSummary
}

func (s *Server) getCronSummary(c *gin.Context) {
	var cronSummary []cronSingleSummary
	for _, benchmarkType := range s.benchmarkTypes {
		data, err := macrobench.GetSummaryForLastDays(benchmarkType, "cron", macrobench.Gen4Planner, 31, s.dbClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
			slog.Error(err)
			return
		}
		cronSummary = append(cronSummary, cronSingleSummary{
			Name: benchmarkType,
			Data: data,
		})
	}
	c.JSON(http.StatusOK, cronSummary)
}

func (s *Server) getCron(c *gin.Context) {
	benchmarkType := c.Query("type")
	data, err := macrobench.GetResultsForLastDays(benchmarkType, "cron", macrobench.Gen4Planner, 31, s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	for i, d := range data {
		m, err := metrics.GetExecutionMetricsSQL(s.dbClient, d.ExecUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
			slog.Error(err)
			return
		}
		data[i].Metrics = m
	}
	c.JSON(http.StatusOK, data)
}
