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
	// MacroBenchmarkType determines the type of a macro benchmark.
	// For instance a macro benchmark can be OLTP or TPCC.
	MacroBenchmarkType string
)

const (
	OLTP = MacroBenchmarkType("oltp")
	TPCC = MacroBenchmarkType("tpcc")

	IncorrectMacroBenchmarkType = "incorrect macrobenchmark type"
)

// Set implements Cobra flag.Value interface.
func (mbtype *MacroBenchmarkType) Set(s string) error {
	*mbtype = MacroBenchmarkType(s)
	return nil
}

// ToUpper returns a new MacroBenchmarkType in upper case.
func (mbtype MacroBenchmarkType) ToUpper() MacroBenchmarkType {
	return MacroBenchmarkType(strings.ToUpper(string(mbtype)))
}

// String returns the given MacroBenchmarkType as a string.
func (mbtype MacroBenchmarkType) String() string {
	return string(mbtype)
}