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

import "github.com/vitessio/arewefastyet/go/storage/mysql"

// CompareMacroBenchmarks takes in 3 arguments, the database, and 2 SHAs. It reads from the database, the macrobenchmark
// results for the 2 SHAs and compares them. The result is a map with the key being the macrobenchmark name.
func CompareMacroBenchmarks(dbClient *mysql.Client, reference string, compare string) (map[Type]interface{}, error) {
	// Get macro benchmarks from all the different types
	SHAs := []string{reference, compare}
	var err error
	macros := map[string]map[Type]DetailsArray{}
	for _, sha := range SHAs {
		macros[sha], err = GetDetailsArraysFromAllTypes(sha, dbClient)
		if err != nil {
			return nil, err
		}
		for mtype := range macros[sha] {
			macros[sha][mtype] = macros[sha][mtype].ReduceSimpleMedian()
		}
	}
	macrosMatrixes := map[Type]interface{}{}
	for _, mtype := range Types {
		macrosMatrixes[mtype] = CompareDetailsArrays(macros[reference][mtype], macros[compare][mtype])
	}
	return macrosMatrixes, nil
}
