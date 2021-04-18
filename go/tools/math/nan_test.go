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

type s struct {
	Vf float64
}

func TestCheckForNaN(t *testing.T) {
	type args struct {
		data  s
		setTo float64
	}
	tests := []struct {
		name  string
		args  args
		wants float64
	}{
		{name: "Float NaN", args: args{data: s{Vf: math.NaN()}, setTo: 100}, wants: 100},
		{name: "Float not NaN", args: args{data: s{Vf: 1.50}, setTo: 100}, wants: 1.50},
		{name: "Float not NaN value 0", args: args{data: s{Vf: 0}, setTo: 100}, wants: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			CheckForNaN(&tt.args.data, tt.args.setTo)
			c.Assert(tt.args.data.Vf, qt.Equals, tt.wants)
		})
	}
}

func BenchmarkCheckForNaN(b *testing.B) {
	bNaN := s{Vf: math.NaN()}
	bNotNaN := s{Vf: 9}

	for i := 0; i < b.N; i++ {
		CheckForNaN(&bNaN, 50)
		CheckForNaN(&bNotNaN, 50)
	}
}