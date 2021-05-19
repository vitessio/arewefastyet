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
	"io/ioutil"
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
	s, err := GetAllVitessReleaseCommitHash(vitessPath)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, s, qt.DeepEquals, []*Release{
		{Name: "10.0.1", CommitHash: "f7304cd1893accfefee0525910098a8e0e68deec"},
		{Name: "10.0.0", CommitHash: "48dccf56282dc79903c0ab0b1d0177617f927403"},
		{Name: "10.0.0-rc1", CommitHash: "29a494f7b45faf26eaaa3e6727b452a2ef254101"},
		{Name: "9.0.1", CommitHash: "42c38e56e4ae29012a5d603d8bc8c22c35b78b52"},
		{Name: "9.0.0", CommitHash: "daa60859822ff85ce18e2d10c61a27b7797ec6b8"},
		{Name: "9.0.0-rc1", CommitHash: "0472d4728ff4b5a0b91834331ff16ab9b0057da8"},
		{Name: "8.0.0", CommitHash: "7e09d0c20ca1e535b9b3f2d96ff2b1ab907d96e8"},
		{Name: "8.0.0-rc1", CommitHash: "098845592f16d235cb440bdafe070b4dcc0eb07b"},
		{Name: "7.0.3", CommitHash: "5f293938aa637e073231e24fe97448f3b6f2579a"},
		{Name: "7.0.2", CommitHash: "aea21dcbfab3d01fedf2ad4b42f9c7727bc47128"},
		{Name: "7.0.1", CommitHash: "19c92a5eabefe4556ae23154e1fee12f977ed1ec"},
		{Name: "7.0.0", CommitHash: "a3a52322d4d24bac4f020ec6fd95418f88276662"},
	})
}

func TestGetAllVitessReleaseCommitHash(t *testing.T) {
	tmpDir, vitessPath, err := createTemporaryVitessClone()
	defer os.RemoveAll(tmpDir)
	qt.Assert(t, err, qt.IsNil)
	s, err := GetAllVitessReleaseCommitHash(vitessPath)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "7.0.3",
		CommitHash: "5f293938aa637e073231e24fe97448f3b6f2579a",
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "9.0.0",
		CommitHash: "daa60859822ff85ce18e2d10c61a27b7797ec6b8",
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "9.0.1",
		CommitHash: "42c38e56e4ae29012a5d603d8bc8c22c35b78b52",
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "7.0.2",
		CommitHash: "aea21dcbfab3d01fedf2ad4b42f9c7727bc47128",
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "7.0.0",
		CommitHash: "a3a52322d4d24bac4f020ec6fd95418f88276662",
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "10.0.0",
		CommitHash: "48dccf56282dc79903c0ab0b1d0177617f927403",
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "10.0.1",
		CommitHash: "f7304cd1893accfefee0525910098a8e0e68deec",
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "8.0.0",
		CommitHash: "7e09d0c20ca1e535b9b3f2d96ff2b1ab907d96e8",
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "7.0.1",
		CommitHash: "19c92a5eabefe4556ae23154e1fee12f977ed1ec",
	})
	qt.Assert(t, s, qt.Any(qt.DeepEquals), &Release{
		Name:       "9.0.0-rc1",
		CommitHash: "0472d4728ff4b5a0b91834331ff16ab9b0057da8",
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
	}, {
		versionString:   "10.0.0-rc1",
		expectedVersion: []int{10, 0, 0, 1},
	}}

	for _, s := range testcase {
		t.Run(s.versionString, func(t *testing.T) {
			out, err := getVersionNumbersFromString(s.versionString)
			qt.Assert(t, err, qt.IsNil)
			qt.Assert(t, s.expectedVersion, qt.DeepEquals, out)
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
	}, {
		versionString1:     "9.0.0-rc1",
		versionString2:     "9.0.0",
		expectedComparison: -1,
	}}

	for _, s := range testcase {
		t.Run(s.versionString1+"-"+s.versionString2, func(t *testing.T) {
			out, err := compareReleaseNumbers(s.versionString1, s.versionString2)
			qt.Assert(t, err, qt.IsNil)
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
	tmpDir, err := ioutil.TempDir("", "setup_vitess_*")
	if err != nil {
		return "", "", err
	}
	_, err = ExecCmd(tmpDir, "git", "clone", "https://github.com/vitessio/vitess.git")
	vitessPath := path.Join(tmpDir, "vitess")
	return tmpDir, vitessPath, err
}
