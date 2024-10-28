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
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vitessio/arewefastyet/go/exec"
	"github.com/vitessio/arewefastyet/go/tools/git"
	"github.com/vitessio/arewefastyet/go/tools/github"
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
	"github.com/vitessio/arewefastyet/go/tools/microbench"
	"github.com/vitessio/arewefastyet/go/tools/server"
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
	Workloads []string `json:"workloads"`
	Sources   []string `json:"sources"`
	Statuses  []string `json:"statuses"`
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
	c.JSON(http.StatusOK, s.workloads)
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

type VitessGitRefReleases struct {
	Tags     []*git.Release `json:"tags"`
	Branches []*git.Release `json:"branches"`
}

func (s *Server) getLatestVitessGitRef(c *gin.Context) {
	var response VitessGitRefReleases
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
	mainRelease := &git.Release{
		Name:       "main",
		CommitHash: lastrunDailySHA,
	}
	response.Branches = append(response.Branches, mainRelease)
	// get all the latest release branches as well
	allReleaseBranches, err := git.GetLatestVitessReleaseBranchCommitHash(s.getVitessPath())
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}
	response.Branches = append(response.Branches, allReleaseBranches...)
	response.Tags = allReleases

	c.JSON(http.StatusOK, response)
}

type CompareMacrobench struct {
	Workload string                               `json:"workload"`
	Result   macrobench.StatisticalCompareResults `json:"result"`
}

func (s *Server) compareMacroBenchmarks(c *gin.Context) {
	oldSHA := c.Query("old")
	newSHA := c.Query("new")

	results, err := macrobench.Compare(s.dbClient, oldSHA, newSHA, s.workloads, macrobench.Gen4Planner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}

	resultsSlice := make([]CompareMacrobench, 0, len(results))
	for workload, res := range results {
		resultsSlice = append(resultsSlice, CompareMacrobench{
			Workload: workload,
			Result:   res,
		})
	}

	sort.Slice(resultsSlice, func(i, j int) bool {
		return resultsSlice[i].Workload < resultsSlice[j].Workload
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

	results, err := macrobench.Search(s.dbClient, sha, s.workloads, macrobench.Gen4Planner)
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
	workload := macrobench.Workload(c.Query("workload"))

	if leftGitRef == "" || rightGitRef == "" || workload == "" {
		c.JSON(http.StatusBadRequest, &ErrorAPI{Error: "The gitref left and right and/or the workload are missing"})
		return
	}

	comparison := s.queriesCompare(c, leftGitRef, rightGitRef, workload, workload)
	if comparison == nil {
		return
	}

	c.JSON(http.StatusOK, comparison)
}

func (s *Server) fkQueriesCompareMacrobenchmarks(c *gin.Context) {
	gitRef := c.Query("gitRef")
	oldWorkload := macrobench.Workload(c.Query("oldWorkload"))
	newWorkload := macrobench.Workload(c.Query("newWorkload"))

	if gitRef == "" || oldWorkload == "" || newWorkload == "" {
		c.JSON(http.StatusBadRequest, &ErrorAPI{Error: "The gitref and the two workloads are incorrect or missing. Please kindly add them."})
		return
	}

	comparison := s.queriesCompare(c, gitRef, gitRef, oldWorkload, newWorkload)
	if comparison == nil {
		return
	}

	c.JSON(http.StatusOK, comparison)
}

func (s *Server) queriesCompare(c *gin.Context, oldGitRef, newGitRef string, oldWorkload, newWorkload macrobench.Workload) []macrobench.VTGateQueryPlanComparer {
	oldPlans, err := macrobench.GetVTGateSelectQueryPlansWithFilter(oldGitRef, oldWorkload, macrobench.Gen4Planner, s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return nil
	}
	newPlans, err := macrobench.GetVTGateSelectQueryPlansWithFilter(newGitRef, newWorkload, macrobench.Gen4Planner, s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return nil
	}
	return macrobench.CompareVTGateQueryPlans(oldPlans, newPlans)
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
		workloads = s.workloads
	} else {
		for _, workload := range workloads {
			workload = strings.ToUpper(workload)
			if !slices.Contains(s.workloads, workload) {
				c.JSON(http.StatusBadRequest, &ErrorAPI{Error: "Wrong workload specified"})
				return
			}
		}
	}
	results, err := macrobench.SearchForLast30DaysQPSOnly(s.dbClient, workloads, macrobench.Gen4Planner)
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
	workload := c.Query("workload")
	data, err := macrobench.SearchForLast30Days(s.dbClient, workload, macrobench.Gen4Planner)
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
	workload := c.Query("workload")
	sha := c.Query("sha")
	pswd := c.Query("key")
	v := c.Query("version")

	errStrFmt := "missing argument: %s"
	if workload == "" {
		errStr := fmt.Sprintf(errStrFmt, "workload")
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
	cfg, ok := configs[strings.ToLower(workload)]
	if !ok {
		errMsg := "unknown benchmark workload: " + strings.ToUpper(workload)
		c.JSON(http.StatusBadRequest, &ErrorAPI{Error: errMsg})
		slog.Error(errMsg)
		return
	}

	// create execution element
	elem := s.createSimpleExecutionQueueElement(cfg, "custom_run", sha, workload, string(macrobench.Gen4Planner), false, 0, currVersion)

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
	newWorkload := c.Query("newWorkload")
	oldWorkload := c.Query("oldWorkload")

	results, err := macrobench.CompareFKs(s.dbClient, oldWorkload, newWorkload, sha, macrobench.Gen4Planner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}

	c.JSON(http.StatusOK, results)
}

func (s *Server) getHistory(c *gin.Context) {
	results, err := exec.GetHistory(s.dbClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorAPI{Error: err.Error()})
		slog.Error(err)
		return
	}

	c.JSON(http.StatusOK, results)
}

type ExecutionRequest struct {
	Auth               string   `json:"auth"`
	Source             string   `json:"source"`
	SHA                string   `json:"sha"`
	Workloads          []string `json:"workloads"`
	NumberOfExecutions string   `json:"number_of_executions"`
}

func (s *Server) addExecutions(c *gin.Context) {
	var req ExecutionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := s.handleAuthentication(c, req.Auth); err != nil {
		c.JSON(http.StatusUnauthorized, &ErrorAPI{Error: err.Error()})
		return
	}

	if req.Source == "" || req.SHA == "" || len(req.Workloads) == 0 || req.NumberOfExecutions == "" {
		c.JSON(http.StatusBadRequest, &ErrorAPI{Error: "missing argument"})
		return
	}

	if len(req.Workloads) == 1 && req.Workloads[0] == "all" {
		req.Workloads = s.workloads
	}
	execs, err := strconv.Atoi(req.NumberOfExecutions)
	if err != nil {
		c.JSON(http.StatusBadRequest, &ErrorAPI{Error: "numberOfExecutions must be an integer"})
		return
	}
	newElements := make([]*executionQueueElement, 0, execs*len(req.Workloads))

	for _, workload := range req.Workloads {
		for i := 0; i < execs; i++ {
			elem := s.createSimpleExecutionQueueElement(s.benchmarkConfig[strings.ToLower(workload)], req.Source, req.SHA, workload, string(macrobench.Gen4Planner), false, 0, git.Version{})
			elem.identifier.UUID = uuid.NewString()
			newElements = append(newElements, elem)
		}
	}

	s.appendToQueue(newElements)

	c.JSON(http.StatusCreated, "")
}

func (s *Server) clearExecutionQueue(c *gin.Context) {
	var req struct {
		Auth string `json:"auth"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := s.handleAuthentication(c, req.Auth); err != nil {
		c.JSON(http.StatusUnauthorized, &ErrorAPI{Error: err.Error()})
		return
	}

	s.clearQueue()

	c.JSON(http.StatusAccepted, "")
}

func (s *Server) handleAuthentication(c *gin.Context, auth string) error {
	decryptedToken, err := server.Decrypt(auth, s.ghTokenSalt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &ErrorAPI{Error: "Unauthorized"})
		return errors.New("unauthenticated")
	}

	isUserAuthenticated, err := IsUserAuthenticated(decryptedToken)
	if err != nil || !isUserAuthenticated {
		return errors.New("unauthenticated")
	}
	return nil
}

func IsUserAuthenticated(accessToken string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		slog.Error("Error creating request to Github: %v", err)
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error making request to Github: %v", err)
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
