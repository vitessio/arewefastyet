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
	// Type determines the type of a macro benchmark.
	Type string
)

// Type implements Cobra flag.Value interface.
func (mbtype *Type) Type() string {
	return "Type"
}

// Set implements Cobra flag.Value interface.
func (mbtype *Type) Set(s string) error {
	*mbtype = Type(s)
	return nil
}

// ToUpper returns a new Type in upper case.
func (mbtype Type) ToUpper() Type {
	return Type(strings.ToUpper(string(mbtype)))
}

// String returns the given Type as a string.
// It also implements the flag.Value interface.
func (mbtype Type) String() string {
	return string(mbtype)
}
