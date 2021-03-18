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

package equinix

import (
	qt "github.com/frankban/quicktest"
	"github.com/vitessio/arewefastyet/go/infra"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

func TestEquinix_ValidConfig(t *testing.T) {
	newTmpDir := func() string {
		path, _ := ioutil.TempDir("", "")
		return path
	}
	tests := []struct {
		name    string
		e       Equinix
		wantErr string
	}{
		{name: "Valid Equinix configuration", e: Equinix{Token: "token", ProjectID: "projectID", InfraCfg: &infra.Config{Path: newTmpDir()}}},
		{name: "Invalid Equinix configuration with missing token", e: Equinix{Token: "", ProjectID: "projectID", InfraCfg: &infra.Config{Path: newTmpDir()}}, wantErr: infra.ErrorInvalidConfiguration},
		{name: "Invalid Equinix configuration with missing project id", e: Equinix{Token: "token", ProjectID: "", InfraCfg: &infra.Config{Path: newTmpDir()}}, wantErr: infra.ErrorInvalidConfiguration},
		{name: "Invalid Equinix configuration with InfraCfg path", e: Equinix{Token: "token", ProjectID: "projectID", InfraCfg: &infra.Config{Path: ""}}, wantErr: infra.ErrorPathMissing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			if tt.e.InfraCfg.Path != "" {
				c.Cleanup(func() { os.RemoveAll(tt.e.InfraCfg.Path) })
				if _, err := os.Stat(tt.e.InfraCfg.Path); err != nil {
					c.Skip("Internal test error:", err.Error())
				}
			}

			gotValid := tt.e.ValidConfig()
			if tt.wantErr != "" {
				c.Assert(gotValid, qt.Not(qt.IsNil))
				c.Assert(gotValid, qt.ErrorMatches, regexp.MustCompile(tt.wantErr+`.*`).String())
			} else {
				c.Assert(gotValid, qt.IsNil)
			}
		})
	}
}
