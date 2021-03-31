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

package infra

import (
	qt "github.com/frankban/quicktest"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestConfig_Close(t *testing.T) {
	newTmpDir := func() string {
		path, _ := ioutil.TempDir("", "")
		return path
	}
	tests := []struct {
		name string
		cfg  Config
	}{
		{name: "Close with valid pathInstallTF", cfg: Config{pathInstallTF: newTmpDir()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			c := qt.New(t)
			c.Cleanup(func() { os.RemoveAll(tt.cfg.pathInstallTF) })

			if _, err = os.Stat(tt.cfg.pathInstallTF); err != nil {
				c.Skip("Internal test error:", err.Error())
			}

			err = tt.cfg.Close()
			c.Assert(err, qt.IsNil)
			_, err = os.Stat(tt.cfg.pathInstallTF)
			c.Assert(err, qt.Not(qt.IsNil))
			c.Assert(os.IsNotExist(err), qt.IsTrue)
		})
	}
}

func TestConfig_Valid(t *testing.T) {
	newTmpDir := func() string {
		path, _ := ioutil.TempDir("", "")
		return path
	}
	tests := []struct {
		name    string
		cfg     Config
		wantErr string
	}{
		{name: "Valid configuration", cfg: Config{Path: newTmpDir()}},
		{name: "Invalid configuration with unknown path", cfg: Config{Path: "/unknown"}, wantErr: ErrorPathUnknown},
		{name: "Invalid configuration with missing path", cfg: Config{Path: ""}, wantErr: ErrorPathMissing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			c.Cleanup(func() { os.RemoveAll(tt.cfg.Path) })

			gotValid := tt.cfg.Valid()

			c.Assert((gotValid == nil) == (tt.wantErr == ""), qt.IsTrue)
			if tt.wantErr != "" && gotValid != nil {
				c.Assert(gotValid, qt.ErrorMatches, tt.wantErr)
			}
		})
	}
}

func TestConfig_CopyTerraformDirectory(t *testing.T) {
	newTmpDir := func() string {
		path, _ := ioutil.TempDir("", "")
		return path
	}
	tests := []struct {
		name    string
		cfg     Config
		dir     string
		wantErr bool
	}{
		{name: "Valid directory", cfg: Config{Path: newTmpDir()}, dir: newTmpDir()},
		{name: "Invalid directory", cfg: Config{Path: newTmpDir()}, dir: "invalid", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			c.Cleanup(func() {
				os.RemoveAll(tt.cfg.Path)
				os.RemoveAll(tt.dir)
			})

			var testFile string
			if tt.cfg.Path != "" {
				testFile = path.Join(tt.cfg.Path, "test.txt")
				_, err := os.Create(testFile)
				if err != nil {
					c.Skip("Internal test error", err.Error())
					return
				}
			}

			err := tt.cfg.CopyTerraformDirectory(tt.dir)
			if tt.wantErr == false {
				c.Assert(err, qt.IsNil)
			} else {
				c.Assert(err, qt.Not(qt.IsNil))
				return
			}

			mustFile := path.Join(tt.dir, "test.txt")
			stat, err := os.Stat(mustFile)
			c.Assert(err, qt.IsNil)
			c.Assert(stat.Name(), qt.Equals, "test.txt")
			c.Assert(tt.cfg.Path, qt.Equals, tt.dir)
		})
	}
}
