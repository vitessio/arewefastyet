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

package math

import "sort"

// MedianInt computes the median of the given int array.
func MedianInt(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	sort.Ints(values)
	middle := len(values) / 2
	if len(values)%2 == 1 {
		return float64(values[middle])
	}
	return float64(values[middle-1] + values[middle]) / 2
}

// MedianFloat computes the median of the given float64 array.
func MedianFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sort.Float64s(values)
	middle := len(values) / 2
	if len(values)%2 == 1 {
		return values[middle]
	}
	return (values[middle-1] + values[middle]) / 2
}
