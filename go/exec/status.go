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
	Total    int
	Finished int
	Last30Days   int
}

func GetBenchmarkStats(client storage.SQLClient) (BenchmarkStats, error) {
	rows, err := client.Select(`SELECT
            (SELECT COUNT(uuid) FROM execution) AS count_status,
            (SELECT COUNT(*) FROM execution WHERE status = 'finished') AS count_finished,
            (SELECT COUNT(*) FROM execution WHERE started_at >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)) AS count_all
        FROM
            execution
		LIMIT 1;`)

	if err != nil {
		return BenchmarkStats{}, err
	}

	defer rows.Close()

	var res BenchmarkStats
	if rows.Next() {
		err := rows.Scan(&res.Total, &res.Finished, &res.Last30Days)
		if err != nil {
			return BenchmarkStats{}, err
		}
	}
	return res, nil
}
