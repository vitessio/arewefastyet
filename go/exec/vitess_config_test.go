/*
 *
 * Copyright 2022 The Vitess Authors.
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
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitessio/arewefastyet/go/tools/git"
)

func Test_prepareVitessConfiguration(t *testing.T) {
	type args struct {
		rawVitessConfig rawSingleVitessVersionConfig
		vitessVersion   git.Version
		vitessConfig    *vitessConfig
	}
	tests := []struct {
		name             string
		args             args
		wantErr          bool
		wantVitessConfig vitessConfig
	}{
		{name: "No configuration", args: args{
			rawVitessConfig: nil,
			vitessVersion:   git.Version{Major: 14},
			vitessConfig:    &vitessConfig{},
		}},

		{name: "Single configuration that is the same", args: args{
			rawVitessConfig: map[string]interface{}{
				"15-0-0": map[string]interface{}{
					"vtgate":   "--toto=1",
					"vttablet": "--titi=1",
				},
			},
			vitessVersion: git.Version{Major: 15},
			vitessConfig:  &vitessConfig{},
		}, wantVitessConfig: vitessConfig{vtgate: "--toto=1", vttablet: "--titi=1"}},

		{name: "Single configuration that is different", args: args{
			rawVitessConfig: map[string]interface{}{
				"15-0-0": map[string]interface{}{
					"vtgate":   "--toto=1",
					"vttablet": "--titi=1",
				},
			},
			vitessVersion: git.Version{Major: 15, Patch: 1},
			vitessConfig:  &vitessConfig{},
		}, wantVitessConfig: vitessConfig{vtgate: "--toto=1", vttablet: "--titi=1"}},

		{name: "Multiple more recent configuration, match old one", args: args{
			rawVitessConfig: map[string]interface{}{
				"14-0-2": map[string]interface{}{
					"vtgate":   "--toto=2",
					"vttablet": "--titi=2",
				},
				"15-0-0": map[string]interface{}{
					"vtgate":   "--toto=1",
					"vttablet": "--titi=1",
				},
				"14-0-1": map[string]interface{}{
					"vtgate":   "--toto=3",
					"vttablet": "--titi=3",
				},
			},
			vitessVersion: git.Version{Major: 14, Minor: 0, Patch: 4},
			vitessConfig:  &vitessConfig{},
		}, wantVitessConfig: vitessConfig{vtgate: "--toto=2", vttablet: "--titi=2"}},

		{name: "Multiple configurations, match most recent one", args: args{
			rawVitessConfig: map[string]interface{}{
				"14-0-2": map[string]interface{}{
					"vtgate":   "--toto=2",
					"vttablet": "--titi=2",
				},
				"14-0-0": map[string]interface{}{
					"vtgate": "--toto=0",
				},
				"15-0-1": map[string]interface{}{
					"vtgate":   "--toto=1",
					"vttablet": "--titi=1",
				},
				"15-0-0": map[string]interface{}{
					"vtgate": "--toto=19",
				},
			},
			vitessVersion: git.Version{Major: 16, Minor: 0, Patch: 0},
			vitessConfig:  &vitessConfig{},
		}, wantVitessConfig: vitessConfig{vtgate: "--toto=1", vttablet: "--titi=1"}},

		{name: "No matching configuration", args: args{
			rawVitessConfig: map[string]interface{}{
				"14-0-2": map[string]interface{}{
					"vtgate":   "--toto=2",
					"vttablet": "--titi=2",
				},
				"14-0-0": map[string]interface{}{
					"vtgate": "--toto=0",
				},
			},
			vitessVersion: git.Version{Major: 13, Minor: 1, Patch: 2},
			vitessConfig:  &vitessConfig{},
		}, wantVitessConfig: vitessConfig{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := prepareVitessConfiguration(tt.args.rawVitessConfig, tt.args.vitessVersion, tt.args.vitessConfig)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.wantVitessConfig, *tt.args.vitessConfig)
		})
	}
}
