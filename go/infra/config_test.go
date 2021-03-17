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
	"testing"
)

func TestConfig_Close(t *testing.T) {
	newTmpDir := func() string {
		path, _ := ioutil.TempDir("","")
		return path
	}
	tests := []struct {
		name    string
		cfg  Config
	}{
		{name: "Close with valid pathInstallTF", cfg: Config{pathInstallTF: newTmpDir()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			c := qt.New(t)
			c.Cleanup(func () {os.RemoveAll(tt.cfg.pathInstallTF)})

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
