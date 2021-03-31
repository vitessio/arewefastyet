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

package macrobench

import (
	qt "github.com/frankban/quicktest"
	"testing"
)

func TestMacroBenchmarkType_String(t *testing.T) {
	tests := []struct {
		name   string
		mbtype MacroBenchmarkType
		want   string
	}{
		{name: "String TPCC", mbtype: TPCC, want: string(TPCC)},
		{name: "String OLTP", mbtype: OLTP, want: string(OLTP)},
		{name: "Simple string", mbtype: "simple", want: "simple"},
		{name: "Empty string", mbtype: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			c.Assert(tt.mbtype.String(), qt.Equals, tt.want)
		})
	}
}

func TestMacroBenchmarkType_ToUpper(t *testing.T) {
	tests := []struct {
		name   string
		mbtype MacroBenchmarkType
		want   MacroBenchmarkType
	}{
		{name: "String TPCC", mbtype: TPCC, want: MacroBenchmarkType("TPCC")},
		{name: "String OLTP", mbtype: OLTP, want: MacroBenchmarkType("OLTP")},
		{name: "Simple string", mbtype: "simple", want: MacroBenchmarkType("SIMPLE")},
		{name: "Empty string", mbtype: "", want: MacroBenchmarkType("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mbtype.ToUpper(); got != tt.want {
				t.Errorf("ToUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}
