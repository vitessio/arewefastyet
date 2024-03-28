/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	_ "math/rand"
)

func main() {
	// Example data for two independent samples
	// Example data for two independent samples with a 20% difference
	sample2 := []float64{75, 79, 82, 88, 72, 85, 92, 88, 78, 95, 81, 87, 90, 83, 91}

	// Generate sample1 with values 20% higher than sample2
	sample1 := make([]float64, len(sample2))
	for i, val := range sample2 {
		sample1[i] = val * 1.2
	}

}
