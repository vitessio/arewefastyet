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

package git

import (
	"net/url"
)

// GetPullRequestsFromGitHub fetches every pull requests in the provided repo that have the
// given set of labels. It then returns each pull request's head SHA.
// The format for repo is: "{USERNAME}/{REPO_NAME}", i.e "vitessio/vitess".
func GetPullRequestsFromGitHub(labels []string, repo string) ([]PRInfo, error) {
	pulls, err := getPullRequestsForLabels(labelsToURL(labels), repo)
	if err != nil {
		return nil, err
	}

	prInfos := []PRInfo{}
	for _, pull := range pulls {
		prInfo, err := getPullRequestHeadAndBase(pull)
		if err != nil {
			return nil, err
		}
		prInfos = append(prInfos, prInfo)
	}
	return prInfos, nil
}

func labelsToURL(labels []string) string {
	result := ""
	for i, label := range labels {
		result += "label:" + url.PathEscape(label)
		if i+1 < len(labels) {
			result += "+"
		}
	}
	return result
}
