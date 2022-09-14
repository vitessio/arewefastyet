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
)

func TestRun(t *testing.T) {
	type args struct {
		port            string
		templatePath    string
		staticPath      string
		localVitessPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     string
	}{
		{name: "Missing port", args: args{templatePath: "./", staticPath: "./", localVitessPath: "~/"}, wantErr: true, err: ErrorIncorrectConfiguration},
		{name: "Missing template path", args: args{port: "8888", staticPath: "./", localVitessPath: "~/"}, wantErr: true, err: ErrorIncorrectConfiguration},
		{name: "Missing static path", args: args{port: "9999", templatePath: "./", localVitessPath: "~/"}, wantErr: true, err: ErrorIncorrectConfiguration},
		{name: "Missing local vitess path", args: args{port: "8080", templatePath: "./", staticPath: "./static"}, wantErr: true, err: ErrorIncorrectConfiguration},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			gotErr := Run(tt.args.port, tt.args.templatePath, tt.args.staticPath, tt.args.localVitessPath)
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
		{name: "Server not ready", s: &Server{}, wantErr: true, err: ErrorIncorrectConfiguration},
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
		{name: "Server fully ready", s: &Server{port: "8080", templatePath: "./", staticPath: "./", localVitessPath: "~/"}, want: true},
		{name: "Missing port", s: &Server{templatePath: "./", staticPath: "./", localVitessPath: "~/"}},
		{name: "Missing template path", s: &Server{port: "8888", staticPath: "./", localVitessPath: "~/"}},
		{name: "Missing static path", s: &Server{port: "9999", templatePath: "./", localVitessPath: "~/"}},
		{name: "Missing api key", s: &Server{port: "8080", templatePath: "./", staticPath: "./static", localVitessPath: "~/"}, want: true},
		{name: "Missing local vitess path", s: &Server{port: "8080", templatePath: "./", staticPath: "./static"}},
		{name: "Missing multiple elements (1)", s: &Server{port: "8080", staticPath: "", localVitessPath: "~/"}},
		{name: "Missing multiple elements (2)", s: &Server{templatePath: "", staticPath: "./", localVitessPath: "~/"}},
		{name: "Missing execution configuration paths", s: &Server{port: "8080", templatePath: "./", staticPath: "./", localVitessPath: "~/"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			gotReady := tt.s.isReady()
			c.Assert(gotReady, qt.Equals, tt.want)
		})
	}
}
