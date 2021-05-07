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
	"net/http"
	"sort"

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
	macrosMatrices, err := macrobench.CompareMacroBenchmarks(s.dbClient, reference, compare)
	if err != nil {
		handleRenderErrors(c, err)
		return
	}

	// Compare Microbenchmarks for the two given SHAs.
	microsMatrix, err := microbench.CompareMicroBenchmarks(s.dbClient, reference, compare)
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

	macros, err := macrobench.GetDetailsArraysFromAllTypes(search, s.dbClient)
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
	// get current results
	currentSHA := c.Query("csh")
	if currentSHA == "" {
		// https://github.com/vitessio/vitess/commit/92584e9bf60f354a6b980717c86a336465ab0354
		currentSHA = "92584e9bf60f354a6b980717c86a336465ab0354" // todo: dynamic value
	}
	currentMbd, err := microbench.GetResultsForGitRef(currentSHA, s.dbClient)
	if err != nil {
		handleRenderErrors(c, err)
		return
	}
	currentMbd = currentMbd.ReduceSimpleMedianByName()

	// get last release results
	lastReleaseSHA := c.Query("lrsh")
	if lastReleaseSHA == "" {
		// https://github.com/vitessio/vitess/commit/daa60859822ff85ce18e2d10c61a27b7797ec6b8
		lastReleaseSHA = "daa60859822ff85ce18e2d10c61a27b7797ec6b8" // todo: dynamic value
	}
	lastReleaseMbd, err := microbench.GetResultsForGitRef(lastReleaseSHA, s.dbClient)
	if err != nil {
		handleRenderErrors(c, err)
		return
	}
	lastReleaseMbd = lastReleaseMbd.ReduceSimpleMedianByName()

	matrix := microbench.MergeMicroBenchmarkDetails(currentMbd, lastReleaseMbd)
	sort.SliceStable(matrix, func(i, j int) bool {
		return !(matrix[i].Current.NSPerOp < matrix[j].Current.NSPerOp)
	})

	c.HTML(http.StatusOK, "microbench.tmpl", gin.H{
		"title":          "Vitess benchmark - microbenchmark",
		"currentSHA":     currentSHA,
		"lastReleaseSHA": lastReleaseSHA,
		"resultMatrix":   matrix,
	})
}

func (s *Server) microbenchmarkSingleResultsHandler(c *gin.Context) {
	name := c.Param("name")

	results, err := microbench.GetLatestResultsFor(name, s.dbClient)
	if err != nil {
		handleRenderErrors(c, err)
		return
	}
	results = results.ReduceSimpleMedianByGitRef()

	c.HTML(http.StatusOK, "microbench_single.tmpl", gin.H{
		"title":   "Vitess benchmark - microbenchmark - " + name,
		"name":    name,
		"results": results,
	})
}
