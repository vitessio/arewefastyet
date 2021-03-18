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

package infra

import (
	qt "github.com/frankban/quicktest"
	"github.com/hashicorp/terraform-exec/tfexec"
	"testing"
)

func TestPopulateTfOption(t *testing.T) {
	type args struct {
		vars []*tfexec.VarOption
		opts interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "PlanOption with a single var", args: args{vars: []*tfexec.VarOption{tfexec.Var("auth_token=token")}, opts: &[]tfexec.PlanOption{}}, wantErr: false},
		{name: "ApplyOption with a single var", args: args{vars: []*tfexec.VarOption{tfexec.Var("amb=false")}, opts: &[]tfexec.ApplyOption{}}, wantErr: false},
		{name: "Invalid opts type", args: args{vars: []*tfexec.VarOption{tfexec.Var("import=true")}, opts: &[]tfexec.ImportOption{}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			err := PopulateTfOption(tt.args.vars, tt.args.opts)
			if tt.wantErr == false {
				c.Assert(err, qt.IsNil)
			} else {
				c.Assert(err, qt.Not(qt.IsNil))
			}
		})
	}
}
