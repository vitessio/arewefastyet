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

package microbench

import (
	qt "github.com/frankban/quicktest"
	"testing"
)

func TestMicroBenchmark(t *testing.T) {
	tests := []struct {
		name        string
		cfg         Config
		wantErr     bool
		errContains string
	}{
		{name: "Invalid package path", cfg: Config{Package: "invalid", Output: "test.txt"}, wantErr: true, errContains: errorInvalidPackageParsing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			gotErr := Run(tt.cfg)
			if tt.wantErr {
				c.Assert(gotErr, qt.Not(qt.IsNil))

				// regexp to catch the whole error
				c.Assert(gotErr, qt.ErrorMatches, tt.errContains+`(.*\n)*`)
				return
			}
			c.Assert(gotErr, qt.IsNil)
		})
	}
}
