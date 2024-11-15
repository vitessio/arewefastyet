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

import (
	"errors"
	"regexp"

	"github.com/vitessio/arewefastyet/go/storage"
	"github.com/vitessio/arewefastyet/go/tools/github"
)

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
	Base string
	Head string
}

func getSourceFromBranchName(branch string) (string, error) {
	if branch == "main" {
		return SourceCron, nil
	}
	matchRelease := regexp.MustCompile(`(release-[0-9]+.0)`)
	matches := matchRelease.FindStringSubmatch(branch)
	if len(matches) == 2 {
		return SourceReleaseBranch + matches[1] + "-branch", nil
	}
	return "", errors.New("no match found")
}

func GetPullRequestInfo(client storage.SQLClient, pullNumber int, info github.PRInfo) (pullRequestInfo, error) {
	rows, err := client.Read("select git_ref from execution where pull_nb = ? and status = 'finished' and source = 'cron_pr' order by started_at desc limit 1", pullNumber)
	if err != nil {
		return pullRequestInfo{}, err
	}

	defer rows.Close()

	var res pullRequestInfo
	if rows.Next() {
		err = rows.Scan(&res.Head)
		if err != nil {
			return pullRequestInfo{}, err
		}
	}

	source, err := getSourceFromBranchName(info.BaseRef)
	if err != nil {
		return res, nil
	}

	if info.IsMerged {
		res.Base, err = getGitRefOfFinishedMatchingSourceGivenTimestamp(client, source, info.MergedTime)
	} else {
		res.Base, err = getGitRefOfLatestFinishedMatchingSource(client, source)
	}
	if err != nil {
		return pullRequestInfo{}, err
	}
	return res, nil
}
