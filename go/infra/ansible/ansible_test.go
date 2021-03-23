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

package ansible

import (
	qt "github.com/frankban/quicktest"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestConfig_MoveRootFolder(t *testing.T) {
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
		{name: "Valid directory", cfg: Config{RootDir: newTmpDir()}, dir: newTmpDir()},
		{name: "Invalid directory", cfg: Config{RootDir: newTmpDir()}, dir: "invalid", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			c.Cleanup(func() {
				os.RemoveAll(tt.cfg.RootDir)
				os.RemoveAll(tt.dir)
			})

			var testFile string
			if tt.cfg.RootDir != "" {
				testFile = path.Join(tt.cfg.RootDir, "test.txt")
				os.Create(testFile)
			}

			err := tt.cfg.CopyRootFolder(tt.dir)
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
		})
	}
}
