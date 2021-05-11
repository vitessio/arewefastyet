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
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/stretchr/testify/require"
)

func TestShortenSHA(t *testing.T) {
	tests := []struct {
		name string
		sha  string
		want string
	}{
		{name: "Regular SHA", sha: "5a504473aec6176b2523bf935ffe4217f61e9928", want: "5a50447"},
		{name: "Short SHA", sha: "5a504473", want: "5a50447"},
		{name: "Tiny SHA", sha: "5a50", want: "5a50"},
		{name: "Empty", sha: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			gotSHA := ShortenSHA(tt.sha)
			c.Assert(gotSHA, qt.Equals, tt.want)
		})
	}
}

func TestGetAllVitessReleaseCommitHash(t *testing.T) {
	s, err := GetAllVitessReleaseCommitHash()
	require.NoError(t, err)
	require.Contains(t, s, &Release{
		Name:       "5.0.1",
		CommitHash: "5165f851ecce1e58d12461ce17e401c2b7788139",
	})
	require.Contains(t, s, &Release{
		Name:       "7.0.3",
		CommitHash: "5f293938aa637e073231e24fe97448f3b6f2579a",
	})
	require.Contains(t, s, &Release{
		Name:       "9.0.1",
		CommitHash: "c970e775be7ec79066aeddd307d050107e66c698",
	})
	require.Contains(t, s, &Release{
		Name:       "7.0.2",
		CommitHash: "aea21dcbfab3d01fedf2ad4b42f9c7727bc47128",
	})
	require.Contains(t, s, &Release{
		Name:       "5.0.0",
		CommitHash: "1b384b8a7c96b1c0ca4fdec62af7295004df9eab",
	})
	require.Contains(t, s, &Release{
		Name:       "4.0.0",
		CommitHash: "cc07de2a374699e645fd1273c48b0948bdd38fca",
	})
	require.Contains(t, s, &Release{
		Name:       "7.0.0",
		CommitHash: "a3a52322d4d24bac4f020ec6fd95418f88276662",
	})
	require.Contains(t, s, &Release{
		Name:       "0.7.0",
		CommitHash: "a3a52322d4d24bac4f020ec6fd95418f88276662",
	})
	require.Contains(t, s, &Release{
		Name:       "2.1.0",
		CommitHash: "5f18b1ed2140b1a7ffcf8d50df69f97cafe38f60",
	})
	require.Contains(t, s, &Release{
		Name:       "10.0.0",
		CommitHash: "48dccf56282dc79903c0ab0b1d0177617f927403",
	})
	require.Contains(t, s, &Release{
		Name:       "2.1.1",
		CommitHash: "405183279f617f941f42d7cbcb54259a3c1a6315",
	})
	require.Contains(t, s, &Release{
		Name:       "10.0.1",
		CommitHash: "05e745812189a3ebbf11b7ade0329510de47dcc3",
	})
	require.Contains(t, s, &Release{
		Name:       "0.8.0",
		CommitHash: "7e09d0c20ca1e535b9b3f2d96ff2b1ab907d96e8",
	})
	require.Contains(t, s, &Release{
		Name:       "8.0.0",
		CommitHash: "7e09d0c20ca1e535b9b3f2d96ff2b1ab907d96e8",
	})
	require.Contains(t, s, &Release{
		Name:       "4.0.1",
		CommitHash: "0629f0da20ab3a78459951137a8482ed804da8b9",
	})
	require.Contains(t, s, &Release{
		Name:       "2.0.0",
		CommitHash: "51ce1ea9e6c70d050be3111d209330885df9c7e3",
	})
	require.Contains(t, s, &Release{
		Name:       "7.0.1",
		CommitHash: "19c92a5eabefe4556ae23154e1fee12f977ed1ec",
	})
	require.Contains(t, s, &Release{
		Name:       "0.9.0",
		CommitHash: "22d6fc0962e366f87b7039efd0f78e2a8c13091f",
	})
}

func TestGetVersionNumbersFromString(t *testing.T) {
	testcase := []struct {
		versionString   string
		expectedVersion []int
	}{{
		versionString:   "7.0.1",
		expectedVersion: []int{7, 0, 1},
	}, {
		versionString:   "10.0.1",
		expectedVersion: []int{10, 0, 1},
	}, {
		versionString:   "10.0",
		expectedVersion: []int{10, 0},
	}, {
		versionString:   "101.132.14.6",
		expectedVersion: []int{101, 132, 14, 6},
	}}

	for _, s := range testcase {
		t.Run(s.versionString, func(t *testing.T) {
			out, err := getVersionNumbersFromString(s.versionString)
			require.NoError(t, err)
			require.Equal(t, s.expectedVersion, out)
		})
	}
}

func TestCompareReleaseNumbers(t *testing.T) {
	testcase := []struct {
		versionString1     string
		versionString2     string
		expectedComparison int
	}{{
		versionString1:     "10.0.1",
		versionString2:     "9.0.1",
		expectedComparison: 1,
	}, {
		versionString1:     "10.0.1",
		versionString2:     "10.0.1",
		expectedComparison: 0,
	}, {
		versionString1:     "8.10",
		versionString2:     "9.0.1",
		expectedComparison: -1,
	}, {
		versionString1:     "9.0.4",
		versionString2:     "9.0.1",
		expectedComparison: 1,
	}, {
		versionString1:     "9.1.1",
		versionString2:     "9.4.1",
		expectedComparison: -1,
	}}

	for _, s := range testcase {
		t.Run(s.versionString1+"-"+s.versionString2, func(t *testing.T) {
			out, err := compareReleaseNumbers(s.versionString1, s.versionString2)
			require.NoError(t, err)
			require.Equal(t, s.expectedComparison, out)
		})
	}
}
