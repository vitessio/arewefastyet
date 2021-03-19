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

func TestRun(t *testing.T) {
	type args struct {
		port         string
		templatePath string
		staticPath   string
		apiKey       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     string
	}{
		{name: "Missing port", args: args{templatePath: "./", staticPath: "./", apiKey: "api_key"}, wantErr: true, err: ErrorIncorrectConfiguration},
		{name: "Missing template path", args: args{port: "8888", staticPath: "./", apiKey: "424242-848484-ABC"}, wantErr: true, err: ErrorIncorrectConfiguration},
		{name: "Missing static path", args: args{port: "9999", templatePath: "./", apiKey: "my key"}, wantErr: true, err: ErrorIncorrectConfiguration},
		{name: "Missing api key", args: args{port: "8080", templatePath: "./", staticPath: "./static"}, wantErr: true, err: ErrorIncorrectConfiguration},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			gotErr := Run(tt.args.port, tt.args.templatePath, tt.args.staticPath, tt.args.apiKey)
			if tt.wantErr == true {
				c.Assert(gotErr, qt.Not(qt.IsNil))
				c.Assert(gotErr, qt.ErrorMatches, tt.err)
			}
		})
	}
}

func TestServer_Run(t *testing.T) {
	tests := []struct {
		name    string
		s       Server
		wantErr bool
		err     string
	}{
		{name: "Server not ready", s: Server{}, wantErr: true, err: ErrorIncorrectConfiguration},
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
		s    Server
		want bool
	}{
		{name: "Server fully ready", s: Server{port: "8080", templatePath: "./", staticPath: "./", apiKey: "api_key"}, want: true},
		{name: "Missing port", s: Server{templatePath: "./", staticPath: "./", apiKey: "api_key"}},
		{name: "Missing template path", s: Server{port: "8888", staticPath: "./", apiKey: "424242-848484-ABC"}},
		{name: "Missing static path", s: Server{port: "9999", templatePath: "./", apiKey: "my key"}},
		{name: "Missing api key", s: Server{port: "8080", templatePath: "./", staticPath: "./static"}},
		{name: "Missing multiple elements (1)", s: Server{port: "8080", staticPath: "", apiKey: ""}},
		{name: "Missing multiple elements (2)", s: Server{templatePath: "", staticPath: "./", apiKey: "428484242-4284842-IUN81B5465"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			gotReady := tt.s.isReady()
			c.Assert(gotReady, qt.Equals, tt.want)
		})
	}
}
