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
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	qt "github.com/frankban/quicktest"
)

// TestGetVTGateQueryPlans_PascalCaseFields verifies that stats emitted by
// VTGate's /debug/query_plans endpoint are parsed correctly.
//
// Vitess serialises plan statistics with PascalCase JSON keys (ExecCount,
// ExecTime, RowsReturned, …) but VTGateQueryPlanValue currently declares
// snake_case struct tags (exec_count, exec_time, rows_returned, …).  Because
// the keys do not match, Go's JSON decoder silently leaves every stats field
// at zero, so nothing meaningful ever reaches the database.
//
// This test will FAIL until the JSON tags in VTGateQueryPlanValue are updated
// to match the names that Vitess actually emits.
func TestGetVTGateQueryPlans_PascalCaseFields(t *testing.T) {
	// JSON exactly as returned by Vitess's /debug/query_plans endpoint.
	// Field names are PascalCase with no underscores, and ExecTime is an
	// integer number of nanoseconds.
	const vitessResponse = `{
		"select :v1 from dual": {
			"QueryType":    "SELECT",
			"Original":     "select 1 from dual",
			"Instructions": {"OperatorType": "Route", "Variant": "Unsharded"},
			"ExecCount":    42,
			"ExecTime":     1500000,
			"ShardQueries": 10,
			"RowsReturned": 100,
			"RowsAffected": 0,
			"Errors":       2
		}
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, vitessResponse)
	}))
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("parsing test server URL: %v", err)
	}

	plans, err := getVTGateQueryPlans(u.Port())
	c := qt.New(t)
	c.Assert(err, qt.IsNil)
	c.Assert(plans, qt.HasLen, 1)

	plan, ok := plans["select :v1 from dual"]
	c.Assert(ok, qt.IsTrue)

	// All of these assertions currently fail because the snake_case JSON tags
	// in VTGateQueryPlanValue do not match Vitess's PascalCase output, so the
	// decoder leaves every field at its zero value.
	c.Assert(plan.ExecCount, qt.Equals, 42)
	c.Assert(plan.ExecTime, qt.Equals, 1500000)
	c.Assert(plan.RowsReturned, qt.Equals, 100)
	c.Assert(plan.Errors, qt.Equals, 2)
	c.Assert(plan.Instructions, qt.Equals, "{\n\t\"OperatorType\": \"Route\",\n\t\"Variant\": \"Unsharded\"\n}")
}
