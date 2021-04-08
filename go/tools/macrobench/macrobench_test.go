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
	"strings"
	"testing"
)

func TestBuildSysbenchArgString(t *testing.T) {
	type cfg = map[string]string
	tts := []struct {
		step string
		m    cfg
		want []string
	}{
		{step: "prepare", m: cfg{"all_luajit-cmd": "off"}, want: []string{"--luajit-cmd=off"}},
		{step: "prepare", m: cfg{"prepare_luajit-cmd": "on"}, want: []string{"--luajit-cmd=on"}},
		{step: "prepare", m: cfg{"luajit-cmd": "on"}, want: []string{}},
		{step: "warmup", m: cfg{"prepare_luajit-cmd": "on", "warmup_luajit-cmd": "off"}, want: []string{"--luajit-cmd=off"}},
		{step: "prepare", m: cfg{"all_luajit-cmd": "off", "prepare_luajit-cmd": "on"}, want: []string{"--luajit-cmd=on"}},
		{m: cfg{}, want: []string{}},
	}
	for _, tt := range tts {
		t.Run(strings.Join(tt.want, " "), func(t *testing.T) {
			argString := buildSysbenchArgString(tt.m, tt.step)
			c := qt.New(t)
			if len(tt.want) == 0 {
				c.Assert(argString, qt.HasLen, 0)
			} else {

				c.Assert(argString, qt.DeepEquals, tt.want)
			}
		})
	}
}
