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

package exec

import (
	"github.com/vitessio/arewefastyet/go/storage"
)

const (
	StatusCreated  = "created"
	StatusStarted  = "started"
	StatusFailed   = "failed"
	StatusFinished = "finished"
)

type BenchmarkStats struct {
	Total       int
	Last30Days  int
	Commits     int
	AvgDuration float64
	Last7Days   []int
}

func GetBenchmarkStats(client storage.SQLClient) (BenchmarkStats, error) {
	rows, err := client.Read(`SELECT
			(SELECT COUNT(uuid) FROM execution) AS count_status,
			(SELECT COUNT(DISTINCT git_ref) FROM execution) AS count_commits,
			(SELECT COUNT(*) FROM execution WHERE started_at >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)) AS count_all,
			(SELECT IFNULL(AVG(TIMESTAMPDIFF(MINUTE, started_at, finished_at)), 0) AS avg_duration_minutes FROM execution WHERE profile_binary IS NULL AND started_at IS NOT NULL AND finished_at IS NOT NULL AND status NOT IN ('failed', 'started') ORDER BY avg_duration_minutes ASC) AS avg_duration_minutes
		FROM 
			execution
		LIMIT 1;`)

	if err != nil {
		return BenchmarkStats{}, err
	}

	defer rows.Close()

	var res BenchmarkStats
	if rows.Next() {
		err := rows.Scan(&res.Total, &res.Commits, &res.Last30Days, &res.AvgDuration)
		if err != nil {
			return BenchmarkStats{}, err
		}
	}

	rows.Close()

	rows, err = client.Read(`SELECT 
			COUNT(*)
		FROM
			execution
		WHERE
			started_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY) 
		GROUP BY
			DATE_FORMAT(started_at, '%Y%m%d')
		ORDER BY
			DATE_FORMAT(started_at, '%Y%m%d') ASC
		LIMIT 7;`)

	if err != nil {
		return BenchmarkStats{}, err
	}

	for rows.Next() {
		var count int
		err := rows.Scan(&count)
		if err != nil {
			return BenchmarkStats{}, err
		}
		res.Last7Days = append(res.Last7Days, count)
	}

	return res, nil
}
