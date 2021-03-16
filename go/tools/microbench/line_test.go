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
		name            string
		benchTypeWanted BenchType
		stringToParse   string
	}{
		// regular benchmarks
		{name: "Valid benchmark with only space (1)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkEmpty-16 1000000000 0.2439 ns/op\n"},
		{name: "Valid benchmark with only space (2)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkEmptySingleT-1 98413 16.9016 ns/op\n"},
		{name: "Valid benchmark with normal go test output (1)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkTestRegular-8   \t18983  \t1.2 ns/op\n"},
		{name: "Valid benchmark with normal go test output (2)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkEmpty-16 \t101609937 \t   156.3899 ns/op\n"},

		// allocs benchmarks
		{name: "Valid allocs benchmark with normal go test output (1)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkTestAllocs-8   \t18983  \t1.2 ns/op \t1292 B/op \t0 allocs/op\n"},
		{name: "Valid allocs benchmark with normal go test output (2)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkAllocsSingleT-1 98413 16.9016 ns/op \t0 B/op \t198 allocs/op\n"},
		{name: "Valid allocs benchmark with normal go test output (3)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkAllocs-4 \t101609937 \t   156.9 ns/op \t0 B/op \t0 allocs/op\n"},
		{name: "Valid allocs benchmark with only spaces", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkAllocs-4 101609937 156.9 ns/op 0 B/op 0 allocs/op\n"},
		{name: "Valid allocs benchmark with only tabs", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkAllocs-32\t2039889\t12.1224 ns/op\t10 B/op\t1 allocs/op\n"},
		{name: "Valid allocs benchmark with spaces and tabs combined", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkAllocs-16   \t1000000000 0.2439 ns/op  \t90 B/op    \t1 allocs/op\n"},

		// SetBytes benchmarks
		{name: "Valid bytes benchmark with normal go test output (1)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkTestBytes-8   \t9821481  \t1.0071 ns/op \t3806098708.20 MB/s\n"},
		{name: "Valid bytes benchmark with normal go test output (2)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkBytes-32 \t2039889 \t   152.913 ns/op \t40572898.10 MB/s\n"},
		{name: "Valid bytes benchmark with spaces and tabs combined (1)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkBytes-16 1000000000  0.2439 ns/op \t3837911248885.89 MB/s\n"},
		{name: "Valid bytes benchmark with spaces and tabs combined (2)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkFunctionBytes-87   189233   161.9016 ns/op    \t39110351.97 MB/s\n"},
		{name: "Valid bytes benchmark with only spaces", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkBytes-32 2039889 152.913 ns/op 40572898.10 MB/s\n"},
		{name: "Valid bytes benchmark with only tabs", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkBytes-32\t2039889\t152.913 ns/op\t40572898.10 MB/s\n"},

		// both allocs and SetBytes benchmarks
		{name: "Valid mixed benchmark with normal go test output (1)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkBytes-16 1000000000 0.2439 ns/op \t3837911248885.89 MB/s \t90 B/op \t1 allocs/op\n"},
		{name: "Valid mixed benchmark with normal go test output (2)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkBytes-32 \t2039889 \t   152.913 ns/op \t40572898.10 MB/s \t0 B/op \t0 allocs/op\n"},
		{name: "Valid mixed benchmark with spaces and tabs combined (1)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkBytes-16 1000000000 0.2439 ns/op \t3837911248885.89 MB/s\t1292 B/op \t0 allocs/op\n"},
		{name: "Valid mixed benchmark with spaces and tabs combined (2)", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkFunctionBytes-87 189233 161.9016 ns/op   \t39110351.97 MB/s\t1292 B/op   \t0 allocs/op\n"},
		{name: "Valid mixed benchmark with only spaces", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkBytes-32 2039889 152.913 ns/op 40572898.10 MB/s 0 B/op 0 allocs/op\n"},
		{name: "Valid mixed benchmark with only tabs", benchTypeWanted: GeneralBenchmark, stringToParse: "BenchmarkBytes-32\t2039889\t152.913 ns/op\t40572898.10 MB/s\t0 B/op\t0 allocs/op\n"},

		// wrong benchmarks
		{name: "Wrong benchmark with extra space (1)", benchTypeWanted: "", stringToParse: "BenchmarkEmpty-16 1000000000 0.2439 \n"},
		{name: "Wrong benchmark with no new line char (2)", benchTypeWanted: "", stringToParse: "BenchmarkTestRegular-8   \t18983  \t1.2 ns/op"},
		{name: "Wrong benchmark with noise data (3)", benchTypeWanted: "", stringToParse: "BenchmarkEmptySingleT-1 98413 16.9016 ns/op  data\n"},
		{name: "Wrong benchmark with missing attribute name (4)", benchTypeWanted: "", stringToParse: "BenchmarkFunctionBytes-16\t2039889\t152.913 ns/op\t40572898.10 \n"},
		{name: "Wrong benchmark with partially missing attribute name (5)", benchTypeWanted: "", stringToParse: "BenchmarkFunctionBytes-16\t2039889\t152.913 ns/op\t40572898.10 /s\n"},
		{name: "Wrong benchmark with empty string (6)", benchTypeWanted: "", stringToParse: ""},
		{name: "Wrong benchmark with empty new line (7)", benchTypeWanted: "", stringToParse: "\n"},
		{name: "Wrong benchmark with unrelated content (8)", benchTypeWanted: "", stringToParse: "Not a benchmark type\n"},
		{name: "Wrong benchmark with wrong benchmark name prefix (9)", benchTypeWanted: "", stringToParse: "BenchEmpty-16 \t101609937 \t   156.3899 ns/op\n"},              // Should start with "Benchmark"
		{name: "Wrong benchmark with misspelled attribute name (10)", benchTypeWanted: "", stringToParse: "BenchmarkAllocs-4 101609937 156.9 ns/op 0 B/op 0 alloc/op\n"}, // Should be "allocs" not "alloc"
		{name: "Wrong benchmark with missing chars (11)", benchTypeWanted: "", stringToParse: "BenchmarkAllocs-4 101609937 156.9ns/op 0 B /op 0 allocs/op\n"},
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
		if line.benchType != GeneralBenchmark {
			b.Errorf("Benchmark failed, got %v, want %v", line.benchType, GeneralBenchmark)
		}
	}
}

func BenchmarkApplyRegularExprAllocsBenchmark(b *testing.B) {
	b.ReportAllocs()
	line := &benchmarkRunLine{Output: "BenchmarkAllocs-16 1000000000 0.2439 ns/op \t90 B/op \t1 allocs/op\n"}

	for i := 0; i < b.N; i++ {
		line.applyRegularExpr()
		if line.benchType != GeneralBenchmark {
			b.Errorf("Benchmark failed, got %v, want %v", line.benchType, GeneralBenchmark)
		}
	}
}

func BenchmarkApplyRegularExprBytesBenchmark(b *testing.B) {
	b.ReportAllocs()
	line := &benchmarkRunLine{Output: "BenchmarkBytes-16 1000000000 0.2439 ns/op \t3837911248885.89 MB/s\n"}

	for i := 0; i < b.N; i++ {
		line.applyRegularExpr()
		if line.benchType != GeneralBenchmark {
			b.Errorf("Benchmark failed, got %v, want %v", line.benchType, GeneralBenchmark)
		}
	}
}

func BenchmarkApplyRegularExprMixedBenchmark(b *testing.B) {
	b.ReportAllocs()
	line := &benchmarkRunLine{Output: "BenchmarkBytes-16 1000000000 0.2439 ns/op \t3837911248885.89 MB/s \t90 B/op \t1 allocs/op\n"}

	for i := 0; i < b.N; i++ {
		line.applyRegularExpr()
		if line.benchType != GeneralBenchmark {
			b.Errorf("Benchmark failed, got %v, want %v", line.benchType, GeneralBenchmark)
		}
	}
}

func Test_benchmarkRunLine_parseGeneralBenchmarkInvalidSubmatchLen(t *testing.T) {
	tests := []struct {
		name        string
		submatch    []string
		wantErr     bool
		wantResults benchmarkResult
	}{
		{name: "submatch length 0", submatch: []string{}, wantErr: true},
		{name: "submatch length 1", submatch: []string{""}, wantErr: true},
		{name: "submatch length 2", submatch: []string{"", ""}, wantErr: true},
		{name: "submatch length 3", submatch: []string{"", "", ""}, wantErr: true},
		{name: "submatch length 4", submatch: []string{"", "", "", ""}, wantErr: true},
		{name: "submatch length 5", submatch: []string{"", "", "", "", ""}, wantErr: true},
		{name: "submatch length 6", submatch: []string{"", "", "", "", "", ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			line := &benchmarkRunLine{submatch: tt.submatch}
			if err = line.parseGeneralBenchmark(); (err != nil) != tt.wantErr {
				t.Errorf("parseGeneralBenchmark() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_benchmarkRunLine_parseGeneralBenchmark(t *testing.T) {
	tests := []struct {
		name        string
		submatch    []string
		wantErr     bool
		wantResults benchmarkResult
	}{
		{name: "regular benchmark (1)", submatch: []string{"BenchmarkEmpty-16 1000000000 0.2439 ns/op", "BenchmarkEmpty-16", "1000000000", "0.2439", "", "", ""}, wantResults: benchmarkResult{Op: 1000000000, NanosecondPerOp: 0.2439}},
		{name: "regular benchmark (2)", submatch: []string{"BenchmarkEmpty-16 19836 178.2396 ns/op", "BenchmarkEmpty-16", "19836", "178.2396", "", "", ""}, wantResults: benchmarkResult{Op: 19836, NanosecondPerOp: 178.2396}},

		{name: "allocs benchmark (1)", submatch: []string{"BenchmarkAllocs-16 1000000000 0.2439 ns/op \t90 B/op \t1 allocs/op\n", "BenchmarkAllocs-16", "1000000000", "0.2439", "", "90", "1"}, wantResults: benchmarkResult{Op: 1000000000, NanosecondPerOp: 0.2439, BytesPerOp: 90, AllocsPerOp: 1}},
		{name: "allocs benchmark (2)", submatch: []string{"BenchmarkAllocs-16 19836 178.2396 ns/op \t40489 B/op \t190 allocs/op\n", "BenchmarkAllocs-16", "19836", "178.2396", "", "40489", "190"}, wantResults: benchmarkResult{Op: 19836, NanosecondPerOp: 178.2396, BytesPerOp: 40489, AllocsPerOp: 190}},

		{name: "bytes benchmark (1)", submatch: []string{"BenchmarkBytes-16 1000000000 0.2439 ns/op \t3837911248885.89 MB/s\n", "BenchmarkBytes-16", "1000000000", "0.2439", "3837911248885.89", "", ""}, wantResults: benchmarkResult{Op: 1000000000, NanosecondPerOp: 0.2439, MBs: 3837911248885.89}},
		{name: "bytes benchmark (2)", submatch: []string{"BenchmarkBytes-16 19836 178.2396 ns/op \t985291124.89213 MB/s\n", "BenchmarkBytes-16", "19836", "178.2396", "985291124.89213", "", ""}, wantResults: benchmarkResult{Op: 19836, NanosecondPerOp: 178.2396, MBs: 985291124.89213}},

		{name: "mixed benchmark (1)", submatch: []string{"BenchmarkBytes-16 1000000000 0.2439 ns/op \t3837911248885.89 MB/s\n", "BenchmarkBytes-16", "1000000000", "0.2439", "3837911248885.89", "95", "2"}, wantResults: benchmarkResult{Op: 1000000000, NanosecondPerOp: 0.2439, MBs: 3837911248885.89, BytesPerOp: 95, AllocsPerOp: 2}},
		{name: "mixed benchmark (2)", submatch: []string{"BenchmarkBytes-16 19836 178.2396 ns/op \t985291124.89213 MB/s\n", "BenchmarkBytes-16", "19836", "178.2396", "985291124.89213", "173", "9"}, wantResults: benchmarkResult{Op: 19836, NanosecondPerOp: 178.2396, MBs: 985291124.89213, BytesPerOp: 173, AllocsPerOp: 9}},

		{name: "invalid number of ops", submatch: []string{"BenchmarkAllocs-16 wrong 0.2439 ns/op \t90 B/op \t1 allocs/op\n", "BenchmarkAllocs-16", "wrong", "0.2439", "", "90", "1"}, wantErr: true},
		{name: "invalid ns/op", submatch: []string{"BenchmarkAllocs-16 19836 178.wrong ns/op \t40489 B/op \t190 allocs/op\n", "BenchmarkAllocs-16", "19836", "178.wrong", "", "40489", "190"}, wantErr: true},
		{name: "invalid MB/s", submatch: []string{"BenchmarkBytes-16 19836 178.2396 ns/op \twrong.89213 MB/s\n", "BenchmarkBytes-16", "19836", "178.2396", "wrong.89213", "MB"}, wantErr: true},
		{name: "invalid bytes/op", submatch: []string{"BenchmarkAllocs-16 1000000000 0.2439 ns/op \twrong B/op \t1 allocs/op\n", "BenchmarkAllocs-16", "1000000000", "0.2439", "", "wrong", "1"}, wantErr: true},
		{name: "invalid allocs/op", submatch: []string{"BenchmarkAllocs-16 19836 178.2396 ns/op \t40489 B/op \twrong allocs/op\n", "BenchmarkAllocs-16", "19836", "178.2396", "", "40489", "wrong"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			line := &benchmarkRunLine{submatch: tt.submatch}
			if err = line.parseGeneralBenchmark(); (err != nil) != tt.wantErr {
				t.Errorf("parseGeneralBenchmark() error = %v, wantErr %v", err, tt.wantErr)
			} else if (err != nil) == tt.wantErr {
				return
			}
			gotResults := line.results
			if gotResults.Op != tt.wantResults.Op {
				t.Errorf("parseGeneralBenchmark() results.Op = %v, want %v", gotResults.Op, tt.wantResults.Op)
			}
			if gotResults.NanosecondPerOp != tt.wantResults.NanosecondPerOp {
				t.Errorf("parseGeneralBenchmark() results.NanosecondPerOp = %v, want %v", gotResults.NanosecondPerOp, tt.wantResults.NanosecondPerOp)
			}
			if gotResults.MBs != tt.wantResults.MBs {
				t.Errorf("parseGeneralBenchmark() results.MBs = %v, want %v", gotResults.MBs, tt.wantResults.MBs)
			}
			if gotResults.BytesPerOp != tt.wantResults.BytesPerOp {
				t.Errorf("parseGeneralBenchmark() results.BytesPerOp = %v, want %v", gotResults.BytesPerOp, tt.wantResults.BytesPerOp)
			}
			if gotResults.AllocsPerOp != tt.wantResults.AllocsPerOp {
				t.Errorf("parseGeneralBenchmark() results.AllocsPerOp = %v, want %v", gotResults.AllocsPerOp, tt.wantResults.AllocsPerOp)
			}
		})
	}
}

func BenchmarkParseSimpleGeneralBenchmark(b *testing.B) {
	var err error
	line := benchmarkRunLine{submatch: []string{"BenchmarkEmpty-16 1000000000 0.2439 ns/op", "BenchmarkEmpty-16", "1000000000", "0.2439", "", "", ""}}

	for i := 0; i < b.N; i++ {
		err = line.parseGeneralBenchmark()
		if err != nil {
			b.Errorf("Got an error: %v", err)
		}
	}
}

func BenchmarkParseAllocsGeneralBenchmark(b *testing.B) {
	var err error
	line := benchmarkRunLine{submatch: []string{"BenchmarkAllocs-16 1000000000 0.2439 ns/op \t90 B/op \t1 allocs/op\n", "BenchmarkAllocs-16", "1000000000", "0.2439", "", "90", "1"}}

	for i := 0; i < b.N; i++ {
		err = line.parseGeneralBenchmark()
		if err != nil {
			b.Errorf("Got an error: %v", err)
		}
	}
}

func BenchmarkParseBytesGeneralBenchmark(b *testing.B) {
	var err error
	line := benchmarkRunLine{submatch: []string{"BenchmarkBytes-16 1000000000 0.2439 ns/op \t3837911248885.89 MB/s\n", "BenchmarkBytes-16", "1000000000", "0.2439", "3837911248885.89", "", ""}}

	for i := 0; i < b.N; i++ {
		err = line.parseGeneralBenchmark()
		if err != nil {
			b.Errorf("Got an error: %v", err)
		}
	}
}

func BenchmarkParseMixedGeneralBenchmark(b *testing.B) {
	var err error
	line := benchmarkRunLine{submatch: []string{"BenchmarkBytes-16 1000000000 0.2439 ns/op \t3837911248885.89 MB/s\n", "BenchmarkBytes-16", "1000000000", "0.2439", "3837911248885.89", "645", "14"}}

	for i := 0; i < b.N; i++ {
		err = line.parseGeneralBenchmark()
		if err != nil {
			b.Errorf("Got an error: %v", err)
		}
	}
}