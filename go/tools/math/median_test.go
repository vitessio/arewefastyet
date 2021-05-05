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
	"testing"
)

func TestMedianFloat(t *testing.T) {
	tests := []struct {
		name string
		arr []float64
		want float64
	}{
		{name: "No element array", arr: []float64{0}, want: 0},
		{name: "Single element array", arr: []float64{5.00}, want: 5.00},
		{name: "Two elements array", arr: []float64{5.00, 10.00}, want: 7.50},

		{name: "Odd number of elements array (1)", arr: []float64{1.00, 3.00, 5.00, 10.00, 11.00}, want: 5.00},
		{name: "Odd number of elements array (2)", arr: []float64{1, 20.5, 45.3, 78.1, 90.5}, want: 45.3},
		{name: "Odd number of elements in unordered array", arr: []float64{78.1, 1, 90.5, 20.5, 45.3}, want: 45.3},

		{name: "Even number of elements array", arr: []float64{4.5, 6.79, 55.3, 78, 86, 110.99}, want: 66.65},
		{name: "Even number of elements in unordered array", arr: []float64{86, 4.5, 78, 6.79, 55.3, 110.99}, want: 66.65},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			got := MedianFloat(tt.arr)
			c.Assert(got, qt.Equals, tt.want)
		})
	}
}

func TestMedianInt(t *testing.T) {
	tests := []struct {
		name string
		arr []int
		want float64
	}{
		{name: "No element array", arr: []int{}, want: 0},
		{name: "Single element array", arr: []int{5}, want: 5},
		{name: "Two elements array", arr: []int{5.00, 10.00}, want: 7.5},

		{name: "Odd number of elements array (1)", arr: []int{1, 3, 5, 10, 11}, want: 5.00},
		{name: "Odd number of elements array (2)", arr: []int{1, 20, 45, 78, 90}, want: 45},
		{name: "Odd number of elements in unordered array", arr: []int{78, 1, 90, 20, 45}, want: 45},

		{name: "Even number of elements array", arr: []int{4, 6, 55, 78, 86, 110}, want: 66.5},
		{name: "Even number of elements in unordered array", arr: []int{86, 4, 78, 6, 55, 110}, want: 66.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			got := MedianInt(tt.arr)
			c.Assert(got, qt.Equals, tt.want)
		})
	}
}
