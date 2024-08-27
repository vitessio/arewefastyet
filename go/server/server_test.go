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
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/vitessio/arewefastyet/go/tools/server"
)

func TestRun(t *testing.T) {
	type args struct {
		port            string
		localVitessPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     string
	}{
		{name: "Missing port", args: args{localVitessPath: "~/"}, wantErr: true, err: server.ErrorIncorrectConfiguration},
		{name: "Missing local vitess path", args: args{port: "8080"}, wantErr: true, err: server.ErrorIncorrectConfiguration},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			gotErr := Run(tt.args.port, tt.args.localVitessPath)
			if tt.wantErr == true {
				c.Assert(gotErr, qt.Not(qt.IsNil))
				c.Assert(gotErr, qt.ErrorMatches, tt.err)
			} else {
				c.Assert(gotErr, qt.IsNil)
			}
		})
	}
}

func TestServer_Run(t *testing.T) {
	tests := []struct {
		name    string
		s       *Server
		wantErr bool
		err     string
	}{
		{name: "Server not ready", s: &Server{}, wantErr: true, err: server.ErrorIncorrectConfiguration},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			gotErr := tt.s.Run()
			if tt.wantErr == true {
				c.Assert(gotErr, qt.Not(qt.IsNil))
				c.Assert(gotErr, qt.ErrorMatches, tt.err)
			}
		})
	}
}

func TestServer_isReady(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
		want bool
	}{
		{name: "Server fully ready", s: &Server{port: "8080", localVitessPath: "~/"}, want: true},
		{name: "Missing port", s: &Server{localVitessPath: "~/"}},
		{name: "Missing local vitess path", s: &Server{port: "8080"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			gotReady := tt.s.isReady()
			c.Assert(gotReady, qt.Equals, tt.want)
		})
	}
}
