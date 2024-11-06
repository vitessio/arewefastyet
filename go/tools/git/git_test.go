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
	"fmt"
	"os"
	"path"
	"testing"

	qt "github.com/frankban/quicktest"
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

func TestGetAllVitessReleaseCommitHashOrdering(t *testing.T) {
	tmpDir, vitessPath, err := createTemporaryVitessClone()
	defer os.RemoveAll(tmpDir)
	qt.Assert(t, err, qt.IsNil)
	s, err := getAllVitessReleases(vitessPath)
	qt.Assert(t, err, qt.IsNil)

	// ordered releases from v10.0.1 to v9.0.0-rc1
	// these should come one after the other but there can be intermediate values in between owing to new releases
	tc := []*Release{
		{Name: "10.0.1", CommitHash: "f7304cd1893accfefee0525910098a8e0e68deec", Version: Version{10, 0, 1}},
		{Name: "10.0.0", CommitHash: "48dccf56282dc79903c0ab0b1d0177617f927403", Version: Version{10, 0, 0}},
		{Name: "10.0.0-rc1", CommitHash: "29a494f7b45faf26eaaa3e6727b452a2ef254101", Version: Version{10, 0, 0}, RCnumber: 1},
		{Name: "9.0.1", CommitHash: "42c38e56e4ae29012a5d603d8bc8c22c35b78b52", Version: Version{9, 0, 1}},
		{Name: "9.0.0", CommitHash: "daa60859822ff85ce18e2d10c61a27b7797ec6b8", Version: Version{9, 0, 0}},
		{Name: "9.0.0-rc1", CommitHash: "0472d4728ff4b5a0b91834331ff16ab9b0057da8", Version: Version{9, 0, 0}, RCnumber: 1},
	}

	idx := 0
	for _, release := range tc {
		found := false
		// keep iterating over all the releases gotten and try to find the next one
		for ; idx < len(s); idx++ {
			if s[idx].Name == release.Name {
				found = true
				break
			}
		}
		// assert that the current release is found and also assert that its contents match the expectations
		qt.Assert(t, found, qt.IsTrue)
		qt.Assert(t, s[idx], qt.DeepEquals, release)
	}
}

func TestGetAllVitessReleaseCommitHash(t *testing.T) {
	tmpDir, vitessPath, err := createTemporaryVitessClone()
	defer os.RemoveAll(tmpDir)
	qt.Assert(t, err, qt.IsNil)
	s, err := getAllVitessReleases(vitessPath)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "5.0.1",
		CommitHash: "5165f851ecce1e58d12461ce17e401c2b7788139",
		Version:    Version{5, 0, 1},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "7.0.3",
		CommitHash: "5f293938aa637e073231e24fe97448f3b6f2579a",
		Version:    Version{7, 0, 3},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "9.0.0",
		CommitHash: "daa60859822ff85ce18e2d10c61a27b7797ec6b8",
		Version:    Version{9, 0, 0},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "9.0.1",
		CommitHash: "42c38e56e4ae29012a5d603d8bc8c22c35b78b52",
		Version:    Version{9, 0, 1},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "7.0.2",
		CommitHash: "aea21dcbfab3d01fedf2ad4b42f9c7727bc47128",
		Version:    Version{7, 0, 2},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "5.0.0",
		CommitHash: "1b384b8a7c96b1c0ca4fdec62af7295004df9eab",
		Version:    Version{5, 0, 0},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "4.0.0",
		CommitHash: "cc07de2a374699e645fd1273c48b0948bdd38fca",
		Version:    Version{4, 0, 0},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "7.0.0",
		CommitHash: "a3a52322d4d24bac4f020ec6fd95418f88276662",
		Version:    Version{7, 0, 0},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "0.7.0",
		CommitHash: "a3a52322d4d24bac4f020ec6fd95418f88276662",
		Version:    Version{0, 7, 0},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "2.1.0",
		CommitHash: "6c06e70a5d7828ad9f79488a704359661bdb996b",
		Version:    Version{2, 1, 0},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "10.0.0",
		CommitHash: "48dccf56282dc79903c0ab0b1d0177617f927403",
		Version:    Version{10, 0, 0},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "2.1.1",
		CommitHash: "89cc312b4da3004d3b5382cf2b37d70e901b1c36",
		Version:    Version{2, 1, 1},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "10.0.1",
		CommitHash: "f7304cd1893accfefee0525910098a8e0e68deec",
		Version:    Version{10, 0, 1},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "0.8.0",
		CommitHash: "7e09d0c20ca1e535b9b3f2d96ff2b1ab907d96e8",
		Version:    Version{0, 8, 0},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "8.0.0",
		CommitHash: "7e09d0c20ca1e535b9b3f2d96ff2b1ab907d96e8",
		Version:    Version{8, 0, 0},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "4.0.1",
		CommitHash: "0629f0da20ab3a78459951137a8482ed804da8b9",
		Version:    Version{4, 0, 1},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "2.0.0",
		CommitHash: "d429f4015ace1f1366acb28e996172dc6693515c",
		Version:    Version{2, 0, 0},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "7.0.1",
		CommitHash: "19c92a5eabefe4556ae23154e1fee12f977ed1ec",
		Version:    Version{7, 0, 1},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "0.9.0",
		CommitHash: "daa60859822ff85ce18e2d10c61a27b7797ec6b8",
		Version:    Version{0, 9, 0},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "3.0",
		CommitHash: "4f192d1003d128e6d399f0a3b37747d9b970d70c",
		Version:    Version{3, 0, 0},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "2.2",
		CommitHash: "66e84fadcc1a7e956e7ffcebcaaba0b04132ca1f",
		Version:    Version{2, 2, 0},
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "9.0.0-rc1",
		CommitHash: "0472d4728ff4b5a0b91834331ff16ab9b0057da8",
		Version:    Version{9, 0, 0},
		RCnumber:   1,
	})
}

func TestCompareReleaseNumbers(t *testing.T) {
	testcase := []struct {
		release1           *Release
		release2           *Release
		expectedComparison int
	}{{
		release1: &Release{
			Version: Version{10, 0, 1},
		},
		release2: &Release{
			Version: Version{9, 0, 1},
		},
		expectedComparison: 1,
	}, {
		release1: &Release{
			Version: Version{10, 0, 1},
		},
		release2: &Release{
			Version: Version{10, 0, 1},
		},
		expectedComparison: 0,
	}, {
		release1: &Release{
			Version: Version{8, 10, 0},
		},
		release2: &Release{
			Version: Version{9, 0, 1},
		},
		expectedComparison: -1,
	}, {
		release1: &Release{
			Version: Version{9, 0, 4},
		},
		release2: &Release{
			Version: Version{9, 0, 1},
		},
		expectedComparison: 1,
	}, {
		release1: &Release{
			Version: Version{9, 1, 1},
		},
		release2: &Release{
			Version: Version{9, 4, 1},
		},
		expectedComparison: -1,
	}, {
		release1: &Release{
			Version:  Version{9, 0, 0},
			RCnumber: 1,
		},
		release2: &Release{
			Version: Version{9, 0, 0},
		},
		expectedComparison: -1,
	}}

	for _, s := range testcase {
		t.Run(fmt.Sprintf("%v", s.release1.Version)+"-"+fmt.Sprintf("%v", s.release2.Version), func(t *testing.T) {
			out := compareReleaseNumbers(s.release1, s.release2)
			qt.Assert(t, s.expectedComparison, qt.Equals, out)
		})
	}
}

func TestGetCommitHash(t *testing.T) {
	tmpDir, vitessPath, err := createTemporaryVitessClone()
	defer os.RemoveAll(tmpDir)
	qt.Assert(t, err, qt.IsNil)
	out, err := GetCommitHash(vitessPath)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, len(out), qt.Equals, 40)
}

// createTemporaryVitessClone creates a temporary vitess clone
func createTemporaryVitessClone() (string, string, error) {
	// Create a temporary folder and clone vitess repo
	tmpDir, err := os.MkdirTemp("", "setup_vitess_*")
	if err != nil {
		return "", "", err
	}
	_, err = ExecCmd(tmpDir, "git", "clone", "https://github.com/vitessio/vitess.git")
	vitessPath := path.Join(tmpDir, "vitess")
	return tmpDir, vitessPath, err
}

func TestGetAllVitessReleaseBranchCommitHash(t *testing.T) {
	tmpDir, vitessPath, err := createTemporaryVitessClone()
	defer os.RemoveAll(tmpDir)
	qt.Assert(t, err, qt.IsNil)
	out, err := GetAllVitessReleaseBranchCommitHash(vitessPath)
	qt.Assert(t, err, qt.IsNil)
	for _, release := range out {
		qt.Assert(t, len(release.CommitHash), qt.Equals, 40)
		qt.Assert(t, release.Name, qt.Contains, "release-")
	}
}

func TestGetLatestVitessReleaseBranchCommitHash(t *testing.T) {
	tmpDir, vitessPath, err := createTemporaryVitessClone()
	defer os.RemoveAll(tmpDir)
	qt.Assert(t, err, qt.IsNil)
	out, err := GetLatestVitessReleaseBranchCommitHash(vitessPath)
	qt.Assert(t, err, qt.IsNil)
	for _, release := range out {
		qt.Assert(t, len(release.CommitHash), qt.Equals, 40)
		qt.Assert(t, release.Name, qt.Contains, "release-")
		qt.Assert(t, release.Version.Major >= 7, qt.IsTrue)
	}
}

func TestGetLatestVitessReleaseCommitHash(t *testing.T) {
	tmpDir, vitessPath, err := createTemporaryVitessClone()
	defer os.RemoveAll(tmpDir)
	qt.Assert(t, err, qt.IsNil)
	out, err := GetSupportedVitessReleases(vitessPath)
	qt.Assert(t, err, qt.IsNil)
	for _, release := range out {
		qt.Assert(t, len(release.CommitHash), qt.Equals, 40)
		qt.Assert(t, release.Version.Major >= 7, qt.IsTrue)
	}
}
