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
		{name: "Not a benchmark 9", benchTypeWanted: "", stringToParse: "BenchEmpty-16 \t101609937 \t   156.3899 ns/op\n"},              // Should start with "Benchmark"
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

func Test_benchmarkRunLine_parseRegularBenchmark(t *testing.T) {
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

		{name: "regular benchmark 1", submatch: []string{"BenchmarkEmpty-16 1000000000 0.2439 ns/op", "BenchmarkEmpty-16", "1000000000", "0.2439"}, wantResults: benchmarkResult{Op: 1000000000, NanosecondPerOp: 0.2439}},
		{name: "regular benchmark 2", submatch: []string{"BenchmarkEmpty-16 19836 178.2396 ns/op", "BenchmarkEmpty-16", "19836", "178.2396"}, wantResults: benchmarkResult{Op: 19836, NanosecondPerOp: 178.2396}},

		{name: "invalid number of ops", submatch: []string{"BenchmarkEmpty-16 wrong 178.2396 ns/op", "BenchmarkEmpty-16", "wrong", "178.2396"}, wantErr: true},
		{name: "invalid number of ns/op", submatch: []string{"BenchmarkEmpty-16 1000000000 wrong ns/op", "BenchmarkEmpty-16", "1000000000", "wrong"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			line := &benchmarkRunLine{submatch: tt.submatch}
			if err = line.parseRegularBenchmark(); (err != nil) != tt.wantErr {
				t.Errorf("parseRegularBenchmark() error = %v, wantErr %v", err, tt.wantErr)
			} else if (err != nil) == tt.wantErr {
				return
			}
			gotResults := line.results
			if gotResults.Op != tt.wantResults.Op {
				t.Errorf("parseRegularBenchmark() results.Op = %v, want %v", gotResults.Op, tt.wantResults.Op)
			}
			if gotResults.NanosecondPerOp != tt.wantResults.NanosecondPerOp {
				t.Errorf("parseRegularBenchmark() results.NanosecondPerOp = %v, want %v", gotResults.NanosecondPerOp, tt.wantResults.NanosecondPerOp)
			}
		})
	}
}

func BenchmarkParseRegularBenchmark(b *testing.B) {
	var err error
	line := benchmarkRunLine{submatch: []string{"BenchmarkEmpty-16 1000000000 0.2439 ns/op", "BenchmarkEmpty-16", "1000000000", "0.2439"}}

	for i := 0; i < b.N; i++ {
		err = line.parseRegularBenchmark()
		if err != nil {
			b.Errorf("Got an error: %v", err)
		}
	}
}

func Test_benchmarkRunLine_parseAllocsBenchmark(t *testing.T) {
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

		{name: "allocs benchmark 1", submatch: []string{"BenchmarkAllocs-16 1000000000 0.2439 ns/op \t90 B/op \t1 allocs/op\n", "BenchmarkAllocs-16", "1000000000", "0.2439", "90", "B", "1"}, wantResults: benchmarkResult{Op: 1000000000, NanosecondPerOp: 0.2439, BytesPerOp: 90, AllocsPerOp: 1}},
		{name: "allocs benchmark 2", submatch: []string{"BenchmarkAllocs-16 19836 178.2396 ns/op \t40489 B/op \t190 allocs/op\n", "BenchmarkAllocs-16", "19836", "178.2396", "40489", "B", "190"}, wantResults: benchmarkResult{Op: 19836, NanosecondPerOp: 178.2396, BytesPerOp: 40489, AllocsPerOp: 190}},

		{name: "invalid number of ops", submatch: []string{"BenchmarkAllocs-16 wrong 0.2439 ns/op \t90 B/op \t1 allocs/op\n", "BenchmarkAllocs-16", "wrong", "0.2439", "90", "B", "1"}, wantErr: true},
		{name: "invalid ns/op", submatch: []string{"BenchmarkAllocs-16 19836 178.wrong ns/op \t40489 B/op \t190 allocs/op\n", "BenchmarkAllocs-16", "19836", "178.wrong", "40489", "B", "190"}, wantErr: true},
		{name: "invalid bytes/op", submatch: []string{"BenchmarkAllocs-16 1000000000 0.2439 ns/op \twrong B/op \t1 allocs/op\n", "BenchmarkAllocs-16", "1000000000", "0.2439", "wrong", "B", "1"}, wantErr: true},
		{name: "unsupported 5th index", submatch: []string{"BenchmarkAllocs-16 19836 178.2396 ns/op \t40489 O/op \t190 allocs/op\n", "BenchmarkAllocs-16", "19836", "178.2396", "40489", "O", "190"}, wantErr: true},
		{name: "invalid allocs/op", submatch: []string{"BenchmarkAllocs-16 19836 178.2396 ns/op \t40489 B/op \twrong allocs/op\n", "BenchmarkAllocs-16", "19836", "178.2396", "40489", "B", "wrong"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			line := &benchmarkRunLine{submatch: tt.submatch}
			if err = line.parseAllocsBenchmark(); (err != nil) != tt.wantErr {
				t.Errorf("parseAllocsBenchmark() error = %v, wantErr %v", err, tt.wantErr)
			} else if (err != nil) == tt.wantErr {
				return
			}
			gotResults := line.results
			if gotResults.Op != tt.wantResults.Op {
				t.Errorf("parseAllocsBenchmark() results.Op = %v, want %v", gotResults.Op, tt.wantResults.Op)
			}
			if gotResults.NanosecondPerOp != tt.wantResults.NanosecondPerOp {
				t.Errorf("parseAllocsBenchmark() results.NanosecondPerOp = %v, want %v", gotResults.NanosecondPerOp, tt.wantResults.NanosecondPerOp)
			}
			if gotResults.BytesPerOp != tt.wantResults.BytesPerOp {
				t.Errorf("parseAllocsBenchmark() results.BytesPerOp = %v, want %v", gotResults.BytesPerOp, tt.wantResults.BytesPerOp)
			}
			if gotResults.AllocsPerOp != tt.wantResults.AllocsPerOp {
				t.Errorf("parseAllocsBenchmark() results.AllocsPerOp = %v, want %v", gotResults.AllocsPerOp, tt.wantResults.AllocsPerOp)
			}
		})
	}
}

func BenchmarkParseAllocsBenchmark(b *testing.B) {
	var err error
	line := benchmarkRunLine{submatch: []string{"BenchmarkAllocs-16 1000000000 0.2439 ns/op \t90 B/op \t1 allocs/op\n", "BenchmarkAllocs-16", "1000000000", "0.2439", "90", "B", "1"}}

	for i := 0; i < b.N; i++ {
		err = line.parseAllocsBenchmark()
		if err != nil {
			b.Errorf("Got an error: %v", err)
		}
	}
}

func Test_benchmarkRunLine_parseBytesBenchmark(t *testing.T) {
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

		{name: "bytes benchmark 1", submatch: []string{"BenchmarkBytes-16 1000000000 0.2439 ns/op \t3837911248885.89 MB/s\n", "BenchmarkBytes-16", "1000000000", "0.2439", "3837911248885.89", "MB"}, wantResults: benchmarkResult{Op: 1000000000, NanosecondPerOp: 0.2439, MBs: 3837911248885.89}},
		{name: "bytes benchmark 2", submatch: []string{"BenchmarkBytes-16 19836 178.2396 ns/op \t985291124.89213 MB/s\n", "BenchmarkBytes-16", "19836", "178.2396", "985291124.89213", "MB"}, wantResults: benchmarkResult{Op: 19836, NanosecondPerOp: 178.2396, MBs: 985291124.89213}},

		{name: "invalid number of ops", submatch: []string{"BenchmarkBytes-16 wrong 0.2439 ns/op \t3837911248885.89 MB/s\n", "BenchmarkBytes-16", "wrong", "0.2439", "3837911248885.89", "MB"}, wantErr: true},
		{name: "invalid ns/op", submatch: []string{"BenchmarkBytes-16 19836 wrong ns/op \t985291124.89213 MB/s\n", "BenchmarkBytes-16", "19836", "wrong", "985291124.89213", "MB"}, wantErr: true},
		{name: "invalid MB/s", submatch: []string{"BenchmarkBytes-16 19836 178.2396 ns/op \twrong.89213 MB/s\n", "BenchmarkBytes-16", "19836", "178.2396", "wrong.89213", "MB"}, wantErr: true},
		{name: "unsupported 5th index", submatch: []string{"BenchmarkBytes-16 19836 178.2396 ns/op \t985291124.89213 PB/s\n", "BenchmarkBytes-16", "19836", "178.2396", "985291124.89213", "PB"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			line := &benchmarkRunLine{submatch: tt.submatch}
			if err = line.parseBytesBenchmark(); (err != nil) != tt.wantErr {
				t.Errorf("parseBytesBenchmark() error = %v, wantErr %v", err, tt.wantErr)
			} else if (err != nil) == tt.wantErr {
				return
			}
			gotResults := line.results
			if gotResults.Op != tt.wantResults.Op {
				t.Errorf("parseBytesBenchmark() results.Op = %v, want %v", gotResults.Op, tt.wantResults.Op)
			}
			if gotResults.NanosecondPerOp != tt.wantResults.NanosecondPerOp {
				t.Errorf("parseBytesBenchmark() results.NanosecondPerOp = %v, want %v", gotResults.NanosecondPerOp, tt.wantResults.NanosecondPerOp)
			}
			if gotResults.MBs != tt.wantResults.MBs {
				t.Errorf("parseBytesBenchmark() results.MBs = %v, want %v", gotResults.MBs, tt.wantResults.MBs)
			}
		})
	}
}

func BenchmarkParseBytesBenchmark(b *testing.B) {
	var err error
	line := benchmarkRunLine{submatch: []string{"BenchmarkBytes-16 1000000000 0.2439 ns/op \t3837911248885.89 MB/s\n", "BenchmarkBytes-16", "1000000000", "0.2439", "3837911248885.89", "MB"}}

	for i := 0; i < b.N; i++ {
		err = line.parseBytesBenchmark()
		if err != nil {
			b.Errorf("Got an error: %v", err)
		}
	}
}
