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

package integration

import (
	"fmt"
	"os"
	"testing"
)

func BenchmarkEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

func BenchmarkMulti(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(b.Name())))

	for i := 0; i < b.N; i++ {

	}
}

func BenchmarkAllocs(b *testing.B) {
	b.ReportAllocs()
	null, err := os.Open(os.DevNull)
	if err != nil {
		b.Error(err.Error())
	}
	for i := 0; i < b.N; i++ {
		fmt.Fprintln(null, i)
	}
}

func BenchmarkBytes(b *testing.B) {
	null, err := os.Open(os.DevNull)
	if err != nil {
		b.Error(err.Error())
	}
	b.SetBytes(int64(len(b.Name())))
	for i := 0; i < b.N; i++ {
		fmt.Fprintln(null, b.Name())
	}
}
