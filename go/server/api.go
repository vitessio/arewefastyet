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
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vitessio/arewefastyet/go/exec"
	"github.com/vitessio/arewefastyet/go/tools/git"
	"github.com/vitessio/arewefastyet/go/tools/github"
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
	"github.com/vitessio/arewefastyet/go/tools/microbench"
	"golang.org/x/exp/slices"
)

type ErrorAPI struct {
	Error string `json:"error"`
}

type ExecutionQueue struct {
	Source   string `json:"source"`
	GitRef   string `json:"git_ref"`
	Workload string `json:"workload"`
	PullNb   int    `json:"pull_nb"`
}

type RecentExecutions struct {
	UUID          string     `json:"uuid"`
	Source        string     `json:"source"`
	GitRef        string     `json:"git_ref"`
	Status        string     `json:"status"`
	Workload      string     `json:"workload"`
	PullNb        int        `json:"pull_nb"`
	GolangVersion string     `json:"golang_version"`
	StartedAt     *time.Time `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"`
}

type ExecutionMetadatas struct {
	Workloads  []string           `json:"workloads"`
	Sources    []string           `json:"sources"`
	Statuses   []string           `json:"statuses"`
}

type RecentExecutionsResponse struct {
	Executions []RecentExecutions `json:"executions"`
	ExecutionMetadatas
}


type ExecutionQueueResponse struct {
	Executions []ExecutionQueue `json:"executions"`
	ExecutionMetadatas
}

func (s *Server) getWorkloadList(c *gin.Context) {
	c.JSON(http.StatusOK, s.benchmarkTypes)
}

func (s *Server) getRecentExecutions(c *gin.Context) {
	execs, err := exec.GetRecentExecutions(s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	response := RecentExecutionsResponse{
		Executions: make([]RecentExecutions, 0, len(execs)),
	}
	for _, e := range execs {
		response.Executions = append(response.Executions, RecentExecutions{
			UUID:          e.RawUUID,
			Source:        e.Source,
			GitRef:        e.GitRef,
			Status:        e.Status,
			Workload:      e.Workload,
			PullNb:        e.PullNB,
			GolangVersion: e.GolangVersion,
			StartedAt:     e.StartedAt,
			FinishedAt:    e.FinishedAt,
		})
		if !slices.Contains(response.Workloads, e.Workload) {
			response.Workloads = append(response.Workloads, e.Workload)
		}
		if !slices.Contains(response.Statuses, e.Status) {
			response.Statuses = append(response.Statuses, e.Status)
		}
		if !slices.Contains(response.Sources, e.Source) {
			response.Sources = append(response.Sources, e.Source)
		}
	}
	c.JSON(http.StatusOK, response)
}

func (s *Server) getExecutionsQueue(c *gin.Context) {
	response := ExecutionQueueResponse{
		Executions: make([]ExecutionQueue, 0, len(queue)),
	}
	for _, e := range queue {
		if e.Executing {
			continue
		}
		response.Executions = append(response.Executions, ExecutionQueue{
			Source:   e.identifier.Source,
			GitRef:   e.identifier.GitRef,
			Workload: e.identifier.Workload,
			PullNb:   e.identifier.PullNb,
		})
		if !slices.Contains(response.Workloads, e.identifier.Workload) {
			response.Workloads = append(response.Workloads, e.identifier.Workload)
		}
		if !slices.Contains(response.Sources, e.identifier.Source) {
			response.Sources = append(response.Sources, e.identifier.Source)
		}
	}
	sort.Slice(response.Executions, func(i, j int) bool {
		return response.Executions[i].GitRef > response.Executions[j].GitRef && response.Executions[i].Source > response.Executions[j].Source
	})
	c.JSON(http.StatusOK, response)
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
	lastrunDailySHA, err := exec.GetLatestDailyJobForMacrobenchmarks(s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	mainRelease := []*git.Release{{
		Name:       "main",
		CommitHash: lastrunDailySHA,
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
	Type   string                               `json:"type"`
	Result macrobench.StatisticalCompareResults `json:"result"`
}

func (s *Server) compareMacroBenchmarks(c *gin.Context) {
	oldSHA := c.Query("old")
	newSHA := c.Query("new")

	results, err := macrobench.Compare(s.dbClient, oldSHA, newSHA, s.benchmarkTypes, macrobench.Gen4Planner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}

	resultsSlice := make([]CompareMacrobench, 0, len(results))
	for typeof, res := range results {
		resultsSlice = append(resultsSlice, CompareMacrobench{
			Type:   typeof,
			Result: res,
		})
	}

	sort.Slice(resultsSlice, func(i, j int) bool {
		return resultsSlice[i].Type < resultsSlice[j].Type
	})

	c.JSON(http.StatusOK, resultsSlice)
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
	Macros map[string]macrobench.StatisticalSingleResult
}

func (s *Server) searchBenchmark(c *gin.Context) {
	sha := c.Query("sha")

	results, err := macrobench.Search(s.dbClient, sha, s.benchmarkTypes, macrobench.Gen4Planner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}

	var res searchResult
	res.Macros = results

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

func (s *Server) getPullRequest(c *gin.Context) {
	prNumbers, err := exec.GetPullRequestList(s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	var prs []github.PRInfo
	for _, prNumber := range prNumbers {
		newPRInfo, err := s.ghApp.GetPullRequestInfo(prNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
			slog.Error(err)
			return
		}
		prs = append(prs, newPRInfo)
	}
	c.JSON(http.StatusOK, prs)
}

func (s *Server) getPullRequestInfo(c *gin.Context) {
	pullNbStr := c.Param("nb")
	pullNb, err := strconv.Atoi(pullNbStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	gitPRInfo, err := exec.GetPullRequestInfo(s.dbClient, pullNb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}

	prInfo, err := s.ghApp.GetPullRequestInfo(pullNb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	prInfo.Base = gitPRInfo.Main
	prInfo.Head = gitPRInfo.PR
	c.JSON(http.StatusOK, prInfo)
}

type dailySummaryResp struct {
	Name string                                    `json:"name"`
	Data []macrobench.ShortStatisticalSingleResult `json:"data"`
}

func (s *Server) getDailySummary(c *gin.Context) {
	// Query array allows to get multiple values for the same key
	// For example: /api/daily/summary?workloads=TPCC&workloads=OLTP
	workloads := c.QueryArray("workloads")
	if len(workloads) == 0 {
		workloads = s.benchmarkTypes
	} else {
		for _, workload := range workloads {
			workload = strings.ToUpper(workload)
			if !slices.Contains(s.benchmarkTypes, workload) {
				c.JSON(http.StatusBadRequest, &ErrorAPI{Error: "Wrong workload specified"})
				return
			}
		}
	}
	results, err := macrobench.SearchForLastDaysQPSOnly(s.dbClient, workloads, macrobench.Gen4Planner, 31)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	var resp []dailySummaryResp
	for name, result := range results {
		resp = append(resp, dailySummaryResp{
			Name: name,
			Data: result,
		})
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Name < resp[j].Name
	})
	c.JSON(http.StatusOK, resp)
}

func (s *Server) getDaily(c *gin.Context) {
	benchmarkType := c.Query("workload")
	data, err := macrobench.SearchForLastDays(s.dbClient, benchmarkType, macrobench.Gen4Planner, 31)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	c.JSON(http.StatusOK, data)
}

func (s *Server) getStatusStats(c *gin.Context) {
	stats, err := exec.GetBenchmarkStats(s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (s *Server) requestRun(c *gin.Context) {
	benchmarkType := c.Query("type")
	sha := c.Query("sha")
	pswd := c.Query("key")
	v := c.Query("version")

	errStrFmt := "missing argument: %s"
	if benchmarkType == "" {
		errStr := fmt.Sprintf(errStrFmt, "type")
		c.JSON(http.StatusBadRequest, &ErrorAPI{Error: errStr})
		slog.Error(errStr)
		return
	}

	if sha == "" {
		errStr := fmt.Sprintf(errStrFmt, "sha")
		c.JSON(http.StatusBadRequest, &ErrorAPI{Error: errStr})
		slog.Error(errStr)
		return
	}

	if v == "" {
		errStr := fmt.Sprintf(errStrFmt, "version")
		c.JSON(http.StatusBadRequest, &ErrorAPI{Error: errStr})
		slog.Error(errStr)
		return
	}

	// check request run key is correct
	if pswd != s.requestRunKey {
		errStr := "unauthorized, wrong key"
		c.JSON(http.StatusUnauthorized, &ErrorAPI{Error: errStr})
		slog.Error(errStr)
		return
	}

	// get version from URL
	version, err := strconv.Atoi(v)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	currVersion := git.Version{Major: version}

	configs := s.getConfigFiles()
	cfg, ok := configs[strings.ToLower(benchmarkType)]
	if !ok {
		errMsg := "unknown benchmark type: " + strings.ToUpper(benchmarkType)
		c.JSON(http.StatusBadRequest, &ErrorAPI{Error: errMsg})
		slog.Error(errMsg)
		return
	}

	// create execution element
	elem := s.createSimpleExecutionQueueElement(cfg, "custom_run", sha, benchmarkType, string(macrobench.Gen4Planner), false, 0, currVersion)

	// to new element to the queue
	s.addToQueue(elem)

	c.JSON(http.StatusCreated, "created")
}

func (s *Server) deleteRun(c *gin.Context) {
	uuid := c.Query("uuid")
	sha := c.Query("sha")
	pswd := c.Query("key")

	errStrFmt := "missing argument: %s"
	if uuid == "" {
		errStr := fmt.Sprintf(errStrFmt, "uuid")
		c.JSON(http.StatusBadRequest, &ErrorAPI{Error: errStr})
		slog.Error(errStr)
		return
	}

	if sha == "" {
		errStr := fmt.Sprintf(errStrFmt, "sha")
		c.JSON(http.StatusBadRequest, &ErrorAPI{Error: errStr})
		slog.Error(errStr)
		return
	}

	// check request run key is correct
	if pswd != s.requestRunKey {
		errStr := "unauthorized, wrong key"
		c.JSON(http.StatusUnauthorized, &ErrorAPI{Error: errStr})
		slog.Error(errStr)
		return
	}

	err := exec.DeleteExecution(s.dbClient, sha, uuid, "custom_run")
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	c.JSON(http.StatusOK, "deleted")
}

func (s *Server) compareBenchmarkFKs(c *gin.Context) {
	sha := c.Query("sha")

	var mtypes []string
	for _, benchmarkType := range s.benchmarkTypes {
		if strings.Contains(benchmarkType, "TPCC") {
			mtypes = append(mtypes, benchmarkType)
		}
	}

	results, err := macrobench.Search(s.dbClient, sha, mtypes, macrobench.Gen4Planner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	c.JSON(http.StatusOK, results)
}
