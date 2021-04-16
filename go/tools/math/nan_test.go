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

package math

import (
	qt "github.com/frankban/quicktest"
	"math"
	"testing"
)

func TestCheckForNaN(t *testing.T) {
	type s struct {
		Vf float64
	}

	type args struct {
		data  s
		setTo float64
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Float NaN", args: args{data: s{Vf: math.NaN()}, setTo: 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			CheckForNaN(&tt.args.data, tt.args.setTo)
			c.Assert(tt.args.data.Vf, qt.Equals, tt.args.setTo)
		})
	}
}
