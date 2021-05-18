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
	"net/http"

	"github.com/vitessio/arewefastyet/go/exec"

	"github.com/gin-gonic/gin"
	"github.com/vitessio/arewefastyet/go/tools/git"
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
	"github.com/vitessio/arewefastyet/go/tools/microbench"
)

func handleRenderErrors(c *gin.Context, err error) {
	if err == nil {
		return
	}
	slog.Error(err.Error())
	c.HTML(http.StatusOK, "error.tmpl", gin.H{
		"title": "Vitess benchmark - Error",
		"url":   c.FullPath(),
	})
}

func (s *Server) informationHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "information.tmpl", gin.H{
		"title": "Vitess benchmark",
	})
}

func (s *Server) homeHandler(c *gin.Context) {
	oltpData, err := macrobench.GetResultsForLastDays(macrobench.OLTP, "webhook", 31, s.dbClient)
	if err != nil {
		slog.Warn(err.Error())
	}

	tpccData, err := macrobench.GetResultsForLastDays(macrobench.TPCC, "webhook", 31, s.dbClient)
	if err != nil {
		slog.Warn(err.Error())
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":     "Vitess benchmark",
		"data_oltp": oltpData,
		"data_tpcc": tpccData,
	})
}

func (s *Server) compareHandler(c *gin.Context) {
	reference := c.Query("r")
	compare := c.Query("c")

	compareSHA := map[string]interface{}{
		"SHA":   compare,
		"short": git.ShortenSHA(compare),
	}
	referenceSHA := map[string]interface{}{
		"SHA":   reference,
		"short": git.ShortenSHA(reference),
	}
	if reference == "" || compare == "" {
		c.HTML(http.StatusOK, "compare.tmpl", gin.H{
			"title":     "Vitess benchmark",
			"reference": referenceSHA,
			"compare":   compareSHA,
		})
		return
	}

	// Compare Macrobenchmarks for the two given SHAs.
	macrosMatrices, err := macrobench.CompareMacroBenchmarks(s.dbClient, s.executionMetricsDBClient, reference, compare)
	if err != nil {
		handleRenderErrors(c, err)
		return
	}

	// Compare Microbenchmarks for the two given SHAs.
	microsMatrix, err := microbench.Compare(s.dbClient, reference, compare)
	if err != nil {
		handleRenderErrors(c, err)
		return
	}

	c.HTML(http.StatusOK, "compare.tmpl", gin.H{
		"title":          "Vitess benchmark",
		"reference":      referenceSHA,
		"compare":        compareSHA,
		"microbenchmark": microsMatrix,
		"macrobenchmark": macrosMatrices,
	})
}

func (s *Server) searchHandler(c *gin.Context) {
	search := c.Query("s")
	if search == "" {
		c.HTML(http.StatusOK, "search.tmpl", gin.H{
			"title": "Vitess benchmark",
		})
		return
	}

	macros, err := macrobench.GetDetailsArraysFromAllTypes(search, s.dbClient, s.executionMetricsDBClient)
	if err != nil {
		handleRenderErrors(c, err)
		return
	}

	micro, err := microbench.GetResultsForGitRef(search, s.dbClient)
	if err != nil {
		handleRenderErrors(c, err)
		return
	}
	micro = micro.ReduceSimpleMedianByName()

	c.HTML(http.StatusOK, "search.tmpl", gin.H{
		"title":          "Vitess benchmark",
		"search":         search,
		"shortSHA":       git.ShortenSHA(search),
		"microbenchmark": micro,
		"macrobenchmark": macros,
	})
}

func (s *Server) requestBenchmarkHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "request_benchmark.tmpl", gin.H{
		"title": "Vitess benchmark",
	})
}

func (s *Server) microbenchmarkResultsHandler(c *gin.Context) {
	var err error

	// get all the releases and the last cron job for master
	allReleases, err := git.GetAllVitessReleaseCommitHash(s.getVitessPath())
	if err != nil {
		handleRenderErrors(c, err)
		return
	}
	lastrunCronSHA, err := exec.GetLatestCronJobForMicrobenchmarks(s.dbClient)
	if err != nil {
		handleRenderErrors(c, err)
		return
	}
	allReleases = append(allReleases, &git.Release{
		Name:       "master",
		CommitHash: lastrunCronSHA,
	})

	// initialize left tag and the corresponding sha
	leftTag := c.Query("ltag")
	leftSHA := ""
	if leftTag == "" {
		// get the latest cron job if leftTag is not specified
		leftTag = "master"
		leftSHA = lastrunCronSHA
	} else {
		leftSHA, err = findSHA(allReleases, leftTag)
		if err != nil {
			handleRenderErrors(c, err)
			return
		}
	}

	// initialize right tag and the corresponding sha
	rightTag := c.Query("rtag")
	rightSHA := ""
	if rightTag == "" {
		// get the last release sha if rightSHA is not specified
		rel, err := git.GetLastReleaseAndCommitHash(s.getVitessPath())
		if err != nil {
			handleRenderErrors(c, err)
			return
		}
		rightSHA = rel.CommitHash
		rightTag = rel.Name
	} else {
		rightSHA, err = findSHA(allReleases, rightTag)
		if err != nil {
			handleRenderErrors(c, err)
			return
		}
	}

	// Get the results from the SHAs
	leftMbd, err := microbench.GetResultsForGitRef(leftSHA, s.dbClient)
	if err != nil {
		handleRenderErrors(c, err)
		return
	}
	leftMbd = leftMbd.ReduceSimpleMedianByName()
	rightMbd, err := microbench.GetResultsForGitRef(rightSHA, s.dbClient)
	if err != nil {
		handleRenderErrors(c, err)
		return
	}
	rightMbd = rightMbd.ReduceSimpleMedianByName()

	matrix := microbench.MergeDetails(leftMbd, rightMbd)
	c.HTML(http.StatusOK, "microbench.tmpl", gin.H{
		"title":        "Vitess benchmark - microbenchmark",
		"leftSHA":      leftSHA,
		"rightSHA":     rightSHA,
		"leftTag":      leftTag,
		"rightTag":     rightTag,
		"allReleases":  allReleases,
		"resultMatrix": matrix,
	})
}

func findSHA(releases []*git.Release, tag string) (string, error) {
	for _, release := range releases {
		if release.Name == tag {
			return release.CommitHash, nil
		}
	}
	return "", fmt.Errorf("unknown tag provided %s", tag)
}

func (s *Server) microbenchmarkSingleResultsHandler(c *gin.Context) {
	name := c.Param("name")
	subBenchmarkName := c.Query("subBenchmarkName")

	results, err := microbench.GetLatestResultsFor(name, subBenchmarkName, 10, s.dbClient)
	if err != nil {
		handleRenderErrors(c, err)
		return
	}
	results = results.ReduceSimpleMedianByGitRef()

	c.HTML(http.StatusOK, "microbench_single.tmpl", gin.H{
		"title":            "Vitess benchmark - microbenchmark - " + name,
		"name":             name,
		"subBenchmarkName": subBenchmarkName,
		"results":          results,
	})
}

func (s *Server) macrobenchmarkResultsHandler(c *gin.Context) {
	var err error

	// get all the releases and the last cron job for master
	allReleases, err := git.GetAllVitessReleaseCommitHash(s.getVitessPath())
	if err != nil {
		handleRenderErrors(c, err)
		return
	}
	lastrunCronSHA, err := exec.GetLatestCronJobForMacrobenchmarks(s.dbClient)
	if err != nil {
		handleRenderErrors(c, err)
		return
	}
	allReleases = append(allReleases, &git.Release{
		Name:       "master",
		CommitHash: lastrunCronSHA,
	})

	// initialize left tag and the corresponding sha
	leftTag := c.Query("ltag")
	leftSHA := ""
	if leftTag == "" {
		// get the latest cron job if leftTag is not specified
		leftTag = "master"
		leftSHA = lastrunCronSHA
	} else {
		leftSHA, err = findSHA(allReleases, leftTag)
		if err != nil {
			handleRenderErrors(c, err)
			return
		}
	}

	// initialize right tag and the corresponding sha
	rightTag := c.Query("rtag")
	rightSHA := ""
	if rightTag == "" {
		// get the last release sha if rightSHA is not specified
		rel, err := git.GetLastReleaseAndCommitHash(s.getVitessPath())
		if err != nil {
			handleRenderErrors(c, err)
			return
		}
		rightSHA = rel.CommitHash
		rightTag = rel.Name
	} else {
		rightSHA, err = findSHA(allReleases, rightTag)
		if err != nil {
			handleRenderErrors(c, err)
			return
		}
	}

	// Compare Macrobenchmarks for the two given SHAs.
	macrosMatrices, err := macrobench.CompareMacroBenchmarks(s.dbClient, s.executionMetricsDBClient, leftSHA, rightSHA)
	if err != nil {
		handleRenderErrors(c, err)
		return
	}
	c.HTML(http.StatusOK, "macrobench.tmpl", gin.H{
		"title":          "Vitess benchmark - macrobenchmark",
		"leftSHA":        leftSHA,
		"rightSHA":       rightSHA,
		"leftTag":        leftTag,
		"rightTag":       rightTag,
		"allReleases":    allReleases,
		"macrobenchmark": macrosMatrices,
	})
}
