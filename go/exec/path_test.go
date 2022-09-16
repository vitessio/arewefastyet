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

package exec

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/google/uuid"
)

func Test_createDirFromUUID(t *testing.T) {
	newTmpDir := func() string {
		path, _ := ioutil.TempDir("", "")
		return path
	}
	type args struct {
		uuid uuid.UUID
		root string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Valid absolute root and UUID", args: args{uuid: uuid.New(), root: newTmpDir()}},
		{name: "Valid non existing absolute root and UUID", args: args{uuid: uuid.New(), root: path.Join(newTmpDir(), "subdir")}},
		{name: "Valid relative root and UUID", args: args{uuid: uuid.New(), root: newTmpDir()}},
		{name: "Valid with unknown relative root", args: args{uuid: uuid.New(), root: "unknown"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			c.Cleanup(func() {
				os.RemoveAll(tt.args.root)
			})

			gotDirPath, err := createDirFromUUID(tt.args.uuid, tt.args.root)

			if tt.wantErr == true {
				c.Assert(err, qt.Not(qt.IsNil))
				return
			}
			wantPath := path.Join(tt.args.root, execDir, tt.args.uuid.String())
			wantPath, _ = filepath.Abs(wantPath)
			c.Assert(gotDirPath, qt.Equals, wantPath)
			stat, err := os.Stat(wantPath)
			if err != nil {
				c.Fatal(err)
			}
			c.Assert(stat.IsDir(), qt.IsTrue)
		})
	}
}
