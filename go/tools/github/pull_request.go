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

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v53/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
)

type pullRequestHandler struct {
	githubapp.ClientCreator
}

func (p pullRequestHandler) Handles() []string {
	return []string{"pull_request"}
}

const (
	prefixCommentBenchmarkMePR = `Hello! :wave:

This Pull Request is now handled by arewefastyet. The current HEAD and future commits will be benchmarked.

`
	suffixCommentBenchmarkMePR = `You can find the performance comparison on the [arewefastyet website](https://benchmark.vitess.io:8000/pr/%d).`

	fullBenchmarkMePRComment = prefixCommentBenchmarkMePR + suffixCommentBenchmarkMePR
)

func (p pullRequestHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.PullRequestEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse issue comment event payload")
	}

	switch event.GetAction() {
	case "labeled":
		lbl := event.GetLabel()
		if lbl == nil {
			break
		}
		if lbl.GetName() != "Benchmark me" {
			break
		}

		// create client
		installationID := githubapp.GetInstallationIDFromEvent(&event)
		client, err := p.NewInstallationClient(installationID)
		if err != nil {
			return err
		}

		repo := event.GetRepo()
		ctx, logger := githubapp.PreparePRContext(ctx, installationID, repo, event.GetNumber())

		logger.Info().Msgf("Benchmark me label applied on #%d", event.GetNumber())

		// use client to get comments
		var allComments []*github.PullRequestComment
		perPage := 100
		for page := 1; true; page++ {
			comments, _, err := client.PullRequests.ListComments(ctx, repo.GetOwner().GetLogin(), repo.GetName(), event.GetNumber(), &github.PullRequestListCommentsOptions{
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: perPage,
				},
			})
			if err != nil {
				logger.Error().Err(err).Msgf("failed to get comments on Pull Request #%d", event.GetNumber())
				return err
			}
			allComments = append(allComments, comments...)
			if len(comments) < perPage {
				break
			}
		}

		// look through comments
		for _, comment := range allComments {
			body := comment.GetBody()
			if strings.Contains(body, prefixCommentBenchmarkMePR) {
				logger.Info().Msgf("arewefastyet comment already added to Pull Request #%d", event.GetNumber())
				return nil
			}
		}

		// add comment to PR
		body := fmt.Sprintf(fullBenchmarkMePRComment, event.GetNumber())
		_, _, err = client.Issues.CreateComment(ctx, repo.GetOwner().GetLogin(), repo.GetName(), event.GetNumber(), &github.IssueComment{
			Body: &body,
		})
		if err != nil {
			logger.Error().Err(err).Msgf("failed to add comment to Pull Request #%d", event.GetNumber())
			return err
		}
	}
	return nil
}
