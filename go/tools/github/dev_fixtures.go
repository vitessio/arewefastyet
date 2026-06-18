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
 */

package github

import "time"

// devPRFixture is canned pull-request metadata used in local development, where
// no real GitHub App is configured (appID == 0) and the GitHub API cannot be
// reached. It lets the /pr and /pr/:nb pages render with plausible data so the
// site can be exercised offline.
type devPRFixture struct {
	Author  string
	Title   string
	DaysAgo int
}

// devPRFixtures maps the seeded pull-request numbers (see the matching
// docker-seed SQL in the dev database) to their fake metadata. Keep this in
// sync with any cron_pr executions seeded into the local database.
var devPRFixtures = map[int]devPRFixture{
	17234: {Author: "harshit-gangal", Title: "vtgate: improve query routing for sharded keyspaces", DaysAgo: 2},
	17198: {Author: "systay", Title: "planner: optimize OLTP transaction handling", DaysAgo: 4},
	17156: {Author: "deepthi", Title: "pools: fix connection pool leak under high concurrency", DaysAgo: 6},
	17102: {Author: "frouioui", Title: "tpcc: reduce allocations on the hot transaction path", DaysAgo: 9},
	17045: {Author: "GuptaManan100", Title: "planner: improve plan cache hit rate", DaysAgo: 13},
	16987: {Author: "mattlord", Title: "schema: speed up schema tracking reloads", DaysAgo: 18},
	16921: {Author: "rohit-nayak-ps", Title: "vstream: batch events for lower replication latency", DaysAgo: 24},
	16884: {Author: "shlomi-noach", Title: "onlineddl: cut over with fewer table locks", DaysAgo: 29},
	16830: {Author: "ajm188", Title: "vtctld: parallelize GetSchema across shards", DaysAgo: 33},
	16777: {Author: "vmg", Title: "evalengine: vectorize numeric comparisons", DaysAgo: 38},
	16742: {Author: "dbussink", Title: "mysql: faster packet decoding in the wire protocol", DaysAgo: 42},
	16705: {Author: "harshit-gangal", Title: "vtgate: cache prepared statement plans per session", DaysAgo: 47},
	16668: {Author: "systay", Title: "planner: push down aggregations to vttablet", DaysAgo: 53},
	16610: {Author: "frouioui", Title: "throttler: lower latency of metrics collection", DaysAgo: 60},
}

// devPRInfo returns deterministic fixture metadata for a pull request when the
// GitHub App is not configured. Known seeded PRs get curated titles/authors;
// any other number falls back to a generic, deterministic entry so the pages
// never error out in local development.
func devPRInfo(prNumber int) PRInfo {
	f, ok := devPRFixtures[prNumber]
	if !ok {
		f = devPRFixture{
			Author:  "vitess-bot",
			Title:   "Benchmarked pull request",
			DaysAgo: 1 + prNumber%30,
		}
	}
	created := time.Now().AddDate(0, 0, -f.DaysAgo)
	return PRInfo{
		ID:        prNumber,
		Author:    f.Author,
		Title:     f.Title,
		CreatedAt: &created,
		BaseRef:   "main",
		IsMerged:  false,
	}
}
