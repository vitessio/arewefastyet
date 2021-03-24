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

package construct

import (
	qt "github.com/frankban/quicktest"
	"github.com/vitessio/arewefastyet/go/infra"
	"github.com/vitessio/arewefastyet/go/infra/equinix"
	"reflect"
	"testing"
)

func TestNewInfra(t *testing.T) {
	tests := []struct {
		name    string
		infraName    string
		want    infra.Infra
		wantErr bool
	}{
		{name: "Valid Equinix infra", infraName: equinix.Name, want: infra.Infra(&equinix.Equinix{})},
		{name: "Invalid infra name", infraName: "iFakeInfra", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			got, err := NewInfra(tt.infraName)

			if tt.wantErr == true {
				c.Assert(err, qt.Not(qt.IsNil))
				return
			}
			c.Assert(reflect.DeepEqual(got, tt.want), qt.IsTrue)
		})
	}
}
