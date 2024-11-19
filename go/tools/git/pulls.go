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
	"encoding/json"
	"fmt"
	"net/http"
)

type PRInfo struct {
	Number int
	Base   string
	SHA    string
}

func GetPullRequestHeadAndBase(url string) (PRInfo, error) {
	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return PRInfo{}, err
	}
	response, err := client.Do(request)
	if err != nil {
		return PRInfo{}, err
	}
	defer response.Body.Close()

	res := struct {
		Number int `json:"number"`
		Head   struct {
			SHA string `json:"sha"`
		} `json:"head"`
		Base struct {
			SHA string `json:"sha"`
		} `json:"base"`
	}{}
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return PRInfo{}, err
	}
	prInfo := PRInfo{
		Number: res.Number,
		Base:   res.Base.SHA,
		SHA:    res.Head.SHA,
	}
	return prInfo, nil
}

func getPullRequestsForLabels(labels, repo string) ([]string, error) {
	query := fmt.Sprintf("https://api.github.com/search/issues?q=repo:%s+is:pr+is:open", repo)
	if labels != "" {
		query += "+" + labels
	}
	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, query, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	res := struct {
		Items []struct {
			PullRequest struct {
				URL string `json:"url"`
			} `json:"pull_request"`
		} `json:"items"`
	}{}
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	var pulls []string
	for _, r := range res.Items {
		pulls = append(pulls, r.PullRequest.URL)
	}
	return pulls, nil
}
