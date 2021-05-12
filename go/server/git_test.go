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

package server

import (
	"io/ioutil"
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestSetupLocalVitess(t *testing.T) {
	// Create a temporary folder and try setup vitess
	tmpDir, err := ioutil.TempDir("", "setup_vitess_*")
	defer os.RemoveAll(tmpDir)
	s := Server{
		localVitessPath: tmpDir,
	}
	err = s.setupLocalVitess()
	qt.Assert(t, err, qt.IsNil)
	// read the directory and indeed verify that the .git folder exists
	files, err := ioutil.ReadDir(s.getVitessPath())
	qt.Assert(t, err, qt.IsNil)
	foundGit := false
	for _, file := range files {
		if file.Name() == ".git" {
			foundGit = true
		}
	}
	qt.Assert(t, foundGit, qt.IsTrue)

	// Create a temporary directory and create a vitess folder manually
	tmpDir, err = ioutil.TempDir("", "setup_vitess_*")
	defer os.RemoveAll(tmpDir)
	s = Server{
		localVitessPath: tmpDir,
	}
	err = os.Mkdir(s.getVitessPath(), 0777)
	qt.Assert(t, err, qt.IsNil)
	err = s.setupLocalVitess()
	qt.Assert(t, err, qt.IsNil)
	// assert that if the vitess folder already exists, then it is not cloned again
	files, err = ioutil.ReadDir(s.getVitessPath())
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, len(files), qt.Equals, 0)
}

func TestGetVitessPath(t *testing.T) {
	testcases := []struct {
		locDirectory      string
		expectedVitessDir string
	}{{
		locDirectory:      "/",
		expectedVitessDir: "/vitess",
	}, {
		locDirectory:      "/r",
		expectedVitessDir: "/r/vitess",
	}, {
		locDirectory:      "/r/",
		expectedVitessDir: "/r/vitess",
	}}
	for _, testcase := range testcases {
		t.Run(testcase.locDirectory, func(t *testing.T) {
			s := Server{
				localVitessPath: testcase.locDirectory,
			}
			out := s.getVitessPath()
			qt.Assert(t, out, qt.Equals, testcase.expectedVitessDir)
		})
	}
}

func TestFetchLocalVitess(t *testing.T) {
	// Create a temporary folder and try setup vitess
	tmpDir, err := ioutil.TempDir("", "setup_vitess_*")
	defer os.RemoveAll(tmpDir)
	s := Server{
		localVitessPath: tmpDir,
	}
	err = s.setupLocalVitess()
	qt.Assert(t, err, qt.IsNil)
	err = s.fetchLocalVitess()
	qt.Assert(t, err, qt.IsNil)
}
