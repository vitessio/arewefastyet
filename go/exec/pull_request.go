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

package exec

import "github.com/vitessio/arewefastyet/go/storage"

func GetPullRequestList(client storage.SQLClient) ([]int, error) {
	rows, err := client.Read("select pull_nb from execution where pull_nb > 0 group by pull_nb order by pull_nb desc")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res []int
	for rows.Next() {
		var pullNumber int
		err = rows.Scan(&pullNumber)
		if err != nil {
			return nil, err
		}
		res = append(res, pullNumber)
	}
	return res, nil
}

type pullRequestInfo struct {
	Main string
	PR   string
}

func GetPullRequestInfo(client storage.SQLClient, pullNumber int) (pullRequestInfo, error) {
	rows, err := client.Read("select cron_pr.git_ref as pr, cron_pr_base.git_ref as main from (select git_ref from execution where pull_nb = ? and status = 'finished' and source = 'cron_pr' order by started_at desc limit 1) cron_pr , (select git_ref from execution where pull_nb = ? and status = 'finished' and source = 'cron_pr_base' order by started_at desc limit 1) cron_pr_base ", pullNumber, pullNumber)
	if err != nil {
		return pullRequestInfo{}, err
	}

	defer rows.Close()

	var res pullRequestInfo
	if rows.Next() {
		err = rows.Scan(&res.PR, &res.Main)
		if err != nil {
			return pullRequestInfo{}, err
		}
	}
	return res, nil
}
