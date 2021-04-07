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

type (
	// Mode defines the type of mode on which to run the server.
	// Either testing, development, or production.
	Mode string
)

const (
	// ErrorIncorrectMode indicates that the given mode is not correct.
	ErrorIncorrectMode = "incorrect mode"

	// ProductionMode runs the server in production mode.
	ProductionMode  = Mode("production")

	// DevelopmentMode runs the server in development mode.
	DevelopmentMode = Mode("development")

	// DefaultMode to use if none is specified.
	DefaultMode = DevelopmentMode
)

func (m *Mode) useDefault() {
	*m = DefaultMode
}

func (m Mode) correct() bool {
	modes := []Mode{ProductionMode, DevelopmentMode}
	for _, mode := range modes {
		if mode == m {
			return true
		}
	}
	return false
}
