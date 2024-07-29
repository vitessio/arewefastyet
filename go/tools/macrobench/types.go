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
	"strings"
)

type (
	// Workload determines the workload of a macro-benchmark.
	Workload string
)

// Type implements Cobra flag.Value interface.
func (mbtype *Workload) Type() string {
	return "Workload"
}

// Set implements Cobra flag.Value interface.
func (mbtype *Workload) Set(s string) error {
	*mbtype = Workload(s)
	return nil
}

// ToUpper returns a new Workload in upper case.
func (mbtype Workload) ToUpper() Workload {
	return Workload(strings.ToUpper(string(mbtype)))
}

// String returns the given Workload as a string.
// It also implements the flag.Value interface.
func (mbtype Workload) String() string {
	return string(mbtype)
}
