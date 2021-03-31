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

package macrobench

import (
	"github.com/vitessio/arewefastyet/go/mysql"
	"time"
)

type (
	// QPS represents the QPS table. This table contains the raw
	// results of a macro benchmark.
	QPS struct {
		ID     int
		RefID  int
		Total  float64
		Reads  float64
		Writes float64
		Other  float64
	}

	// MacroBenchmarkResult represents both OLTP and TPCC tables.
	// The two tables share the same schema and can thus be grouped
	// under an unique go struct.
	MacroBenchmarkResult struct {
		QPS
		TPS        float64
		Latency    float64
		Errors     float64
		Reconnects float64
		Time       int
		Threads    float64
	}

	BenchmarkID struct {
		ID        int
		Source    string
		CreatedAt time.Time
	}

	MacroBenchmarkDetails struct {
		BenchmarkID

		// refers to commit
		GitRef string
		Result MacroBenchmarkResult
	}

	MacroBenchmarkDetailsArray []MacroBenchmarkDetails
)

// GetResultsForGitRef
func GetResultsForGitRef(ref string, client *mysql.Client) (err error) {
	return err
}
