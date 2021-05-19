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

package exec

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestGetComparisonLink(t *testing.T) {
	testcases := []struct {
		leftSHA  string
		rightSHA string
		link     string
	}{
		{
			leftSHA:  "71126fe3286a2f0c25f0ab1be1f19ae4664e5571",
			rightSHA: "daa60859822ff85ce18e2d10c61a27b7797ec6b8",
			link:     "https://benchmark.vitess.io/compare?r=71126fe3286a2f0c25f0ab1be1f19ae4664e5571&c=daa60859822ff85ce18e2d10c61a27b7797ec6b8",
		}, {
			leftSHA:  "aea21dcbfab3d01fedf2ad4b42f9c7727bc47128",
			rightSHA: "cc07de2a374699e645fd1273c48b0948bdd38fca",
			link:     "https://benchmark.vitess.io/compare?r=aea21dcbfab3d01fedf2ad4b42f9c7727bc47128&c=cc07de2a374699e645fd1273c48b0948bdd38fca",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.link, func(t *testing.T) {
			out := getComparisonLink(testcase.leftSHA, testcase.rightSHA)
			qt.Assert(t, out, qt.Equals, testcase.link)
		})
	}
}
