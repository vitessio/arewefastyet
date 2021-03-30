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

package errors

import (
	"errors"
	qt "github.com/frankban/quicktest"
	"testing"
)

func TestConcat(t *testing.T) {
	tests := []struct {
		name    string
		errs    []error
		wantErr bool
		wantStr string
	}{
		{name: "Short list of errors", errs: []error{errors.New("err1"), errors.New("err2")}, wantErr: true, wantStr: "err1\nerr2\n"},
		{name: "Empty list of errors", errs: []error{}, wantErr: false},
		{name: "List of errors with empty strings", errs: []error{errors.New(""), errors.New(""), errors.New("")}, wantErr: true, wantStr: "\n\n\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			gotErr := Concat(tt.errs)
			if !tt.wantErr {
				c.Assert(gotErr, qt.IsNil)
				return
			}
			c.Assert(gotErr, qt.Not(qt.IsNil))
			c.Assert(gotErr.Error(), qt.Equals, tt.wantStr)
		})
	}
}
