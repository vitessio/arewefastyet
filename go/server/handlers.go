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
	"github.com/gin-gonic/gin"
	"github.com/vitessio/arewefastyet/go/tools/microbench"
	"log"
	"net/http"
)

func (s *Server) informationHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "information.tmpl", gin.H{
		"title": "Vitess benchmark",
	})
}

func (s *Server) homeHanlder(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Vitess benchmark",
	})
}

func (s *Server) searchCompareHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "search_compare.tmpl", gin.H{
		"title": "Vitess benchmark",
	})
}

func (s *Server) requestBenchmarkHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "request_benchmark.tmpl", gin.H{
		"title": "Vitess benchmark",
	})
}

func (s *Server) microbenchmarkResultsHandler(c *gin.Context) {
	// get current results
	currentSHA := "36110d81be87fdd11e0626eeff529ce67df27661"
	currentMbd, err := microbench.GetResultsForGitRef(currentSHA, s.dbClient)
	if err != nil {
		log.Println(err)
		return
	}
	currentMbd = currentMbd.ReduceSimpleMedian()

	// get last release results
	lastReleaseSHA := "daa60859822ff85ce18e2d10c61a27b7797ec6b8"
	lastReleaseMbd, err := microbench.GetResultsForGitRef(lastReleaseSHA, s.dbClient)
	if err != nil {
		log.Println(err)
		return
	}
	lastReleaseMbd = lastReleaseMbd.ReduceSimpleMedian()

	matrix := microbench.MergeMicroBenchmarkDetails(currentMbd, lastReleaseMbd)

	c.HTML(http.StatusOK, "microbench.tmpl", gin.H{
		"title":          "Vitess benchmark - microbenchmark",
		"currentSHA":     currentSHA,
		"lastReleaseSHA": lastReleaseSHA,
		"resultMatrix":   matrix,
	})
}
