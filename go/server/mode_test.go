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
	qt "github.com/frankban/quicktest"
	"testing"
)

func TestMode_correct(t *testing.T) {
	tests := []struct {
		name string
		m    Mode
		want bool
	}{
		{name: "Production mode", m: ProductionMode, want: true},
		{name: "Development mode", m: DevelopmentMode, want: true},
		{name: "Default mode", m: DefaultMode, want: true},
		{name: "Invalid mode", m: Mode("invalid"), want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			got := tt.m.correct()
			c.Assert(got, qt.Equals, tt.want)
		})
	}
}

func TestMode_useDefault(t *testing.T) {
	tests := []struct {
		name string
		m    Mode
		want Mode
	}{
		{name: "Production mode", m: ProductionMode, want: DefaultMode},
		{name: "Development mode", m: DevelopmentMode, want: DefaultMode},
		{name: "Empty mode", m: "", want: DefaultMode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			tt.m.useDefault()
			c.Assert(tt.m, qt.Equals, tt.want)
		})
	}
}
