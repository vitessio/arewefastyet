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

package microbench

import (
	"testing"
)

func Test_benchmarkRunLine_applyRegularExpr_checkBenchType(t *testing.T) {
	tests := []struct {
		name   string
		benchTypeWanted BenchType
		stringToParse string
	}{
		{name: "Valid regular benchmark 1", benchTypeWanted: RegularBenchmark, stringToParse: "BenchmarkEmpty-16 1000000000 0.2439 ns/op\n"},
		{name: "Valid regular benchmark 2", benchTypeWanted: RegularBenchmark, stringToParse: "BenchmarkTestRegular-8   \t18983  \t1.2 ns/op\n"},
		{name: "Valid regular benchmark 3", benchTypeWanted: RegularBenchmark, stringToParse: "BenchmarkEmptySingleT-1 98413 16.9016 ns/op\n"},
		{name: "Valid regular benchmark 4", benchTypeWanted: RegularBenchmark, stringToParse: "BenchmarkEmpty-16 \t101609937 \t   156.3899 ns/op\n"},

		{name: "Valid allocs benchmark 1", benchTypeWanted: AllocsBenchmark, stringToParse: "BenchmarkAllocs-16 1000000000 0.2439 ns/op \t90 B/op \t1 allocs/op\n"},
		{name: "Valid allocs benchmark 2", benchTypeWanted: AllocsBenchmark, stringToParse: "BenchmarkTestAllocs-8   \t18983  \t1.2 ns/op \t1292 B/op \t0 allocs/op\n"},
		{name: "Valid allocs benchmark 3", benchTypeWanted: AllocsBenchmark, stringToParse: "BenchmarkAllocsSingleT-1 98413 16.9016 ns/op \t0 B/op \t198 allocs/op\n"},
		{name: "Valid allocs benchmark 4", benchTypeWanted: AllocsBenchmark, stringToParse: "BenchmarkAllocs-4 \t101609937 \t   156.9 ns/op \t0 B/op \t0 allocs/op\n"},
		{name: "Valid allocs benchmark (only spaces)", benchTypeWanted: AllocsBenchmark, stringToParse: "BenchmarkAllocs-4 101609937 156.9 ns/op 0 B/op 0 allocs/op\n"},
		{name: "Valid allocs benchmark (only tabs)", benchTypeWanted: AllocsBenchmark, stringToParse: "BenchmarkAllocs-32\t2039889\t12.1224 ns/op\t10 B/op\t1 allocs/op\n"},

		{name: "Valid bytes benchmark 1", benchTypeWanted: BytesBenchmark, stringToParse: "BenchmarkBytes-16 1000000000 0.2439 ns/op \t3837911248885.89 MB/s\n"},
		{name: "Valid bytes benchmark 2", benchTypeWanted: BytesBenchmark, stringToParse: "BenchmarkTestBytes-8   \t9821481  \t1.0071 ns/op \t3806098708.20 MB/s\n"},
		{name: "Valid bytes benchmark 3", benchTypeWanted: BytesBenchmark, stringToParse: "BenchmarkFunctionBytes-87 189233 161.9016 ns/op \t39110351.97 MB/s\n"},
		{name: "Valid bytes benchmark 4", benchTypeWanted: BytesBenchmark, stringToParse: "BenchmarkBytes-32 \t2039889 \t   152.913 ns/op \t40572898.10 MB/s\n"},
		{name: "Valid bytes benchmark (only spaces)", benchTypeWanted: BytesBenchmark, stringToParse: "BenchmarkBytes-32 2039889 152.913 ns/op 40572898.10 MB/s\n"},
		{name: "Valid bytes benchmark (only tabs)", benchTypeWanted: BytesBenchmark, stringToParse: "BenchmarkBytes-32\t2039889\t152.913 ns/op\t40572898.10 MB/s\n"},

		{name: "Not a benchmark 1", benchTypeWanted: "", stringToParse: "BenchmarkEmpty-16 1000000000 0.2439 \n"},
		{name: "Not a benchmark 2", benchTypeWanted: "", stringToParse: "BenchmarkTestRegular-8   \t18983  \t1.2 ns/op"},
		{name: "Not a benchmark 3", benchTypeWanted: "", stringToParse: "BenchmarkEmptySingleT-1 98413 16.9016 ns/op  data\n"},
		{name: "Not a benchmark 4", benchTypeWanted: "", stringToParse: "BenchmarkFunctionBytes-16\t2039889\t152.913 ns/op\t40572898.10 \n"},
		{name: "Not a benchmark 5", benchTypeWanted: "", stringToParse: "BenchmarkFunctionBytes-16\t2039889\t152.913 ns/op\t40572898.10 /s\n"},
		{name: "Not a benchmark 6", benchTypeWanted: "", stringToParse: ""},
		{name: "Not a benchmark 7", benchTypeWanted: "", stringToParse: "\n"},
		{name: "Not a benchmark 8", benchTypeWanted: "", stringToParse: "Not a benchmark type\n"},
		{name: "Not a benchmark 9", benchTypeWanted: "", stringToParse: "BenchEmpty-16 \t101609937 \t   156.3899 ns/op\n"}, // Should start with "Benchmark"
		{name: "Not a benchmark 10", benchTypeWanted: "", stringToParse: "BenchmarkAllocs-4 101609937 156.9 ns/op 0 B/op 0 alloc/op\n"}, // Should be "allocs" not "alloc"
		{name: "Not a benchmark 11", benchTypeWanted: "", stringToParse: "BenchmarkAllocs-4 101609937 156.9ns/op 0 B /op 0 allocs/op\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := &benchmarkRunLine{Output: tt.stringToParse}
			line.applyRegularExpr()

			gotBenchType := line.benchType
			if gotBenchType != tt.benchTypeWanted {
				t.Errorf("line.applyRegularExpr() = %v, want %v", gotBenchType, tt.benchTypeWanted)
			}
		})
	}
}

func BenchmarkApplyRegularExprRegularBenchmark(b *testing.B) {
	b.ReportAllocs()
	line := &benchmarkRunLine{Output: "BenchmarkEmpty-16 1000000000 0.2439 ns/op\n"}

	for i := 0; i < b.N; i++ {
		line.applyRegularExpr()
		if line.benchType != RegularBenchmark {
			b.Errorf("Benchmark failed, got %v, want %v", line.benchType, RegularBenchmark)
		}
	}
}

func BenchmarkApplyRegularExprAllocsBenchmark(b *testing.B) {
	b.ReportAllocs()
	line := &benchmarkRunLine{Output: "BenchmarkAllocs-16 1000000000 0.2439 ns/op \t90 B/op \t1 allocs/op\n"}

	for i := 0; i < b.N; i++ {
		line.applyRegularExpr()
		if line.benchType != AllocsBenchmark {
			b.Errorf("Benchmark failed, got %v, want %v", line.benchType, AllocsBenchmark)
		}
	}
}

func BenchmarkApplyRegularExprBytesBenchmark(b *testing.B) {
	b.ReportAllocs()
	line := &benchmarkRunLine{Output: "BenchmarkBytes-16 1000000000 0.2439 ns/op \t3837911248885.89 MB/s\n"}

	for i := 0; i < b.N; i++ {
		line.applyRegularExpr()
		if line.benchType != BytesBenchmark {
			b.Errorf("Benchmark failed, got %v, want %v", line.benchType, BytesBenchmark)
		}
	}
}

