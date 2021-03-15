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
	"errors"
	"fmt"
	"github.com/vitessio/arewefastyet/go/mysql"
	"regexp"
	"strconv"
	"time"
)

type BenchType string

const (
	ErrorLineUnrecognized = "the line format was unrecognized"
	ErrorLineMalformed    = "the format of the line is malformed"

	RegularBenchmark = BenchType("regular")
	AllocsBenchmark  = BenchType("allocs")
	BytesBenchmark   = BenchType("bytes")
)

var benchmarkResultsRegArray = map[BenchType]*regexp.Regexp{
	AllocsBenchmark:  regexp.MustCompile(`(Benchmark.+\b)\s*([0-9]+)\s+([\d\.]+)\s+ns\/op\s+([\d\.]+)\s+(.+)\/op\s+([\d\.]+)\s+allocs\/op\n`),
	BytesBenchmark:   regexp.MustCompile(`(Benchmark.+\b)\s*([0-9]+)\s+([\d\.]+)\s+ns\/op\s+([\d\.]+)\s+(.+)\/s\n`),
	RegularBenchmark: regexp.MustCompile(`(Benchmark.+\b)\s*([0-9]+)\s+([\d\.]+)\s+ns\/op\n`),
}

type benchmarkResult struct {
	Op              int
	NanosecondPerOp float64
	MBs             float64
	BytesPerOp      float64
	AllocsPerOp     float64
}

type benchmarkRunLine struct {
	Time    time.Time
	Action  string
	Package string
	Output  string
	Elapsed string

	id        int64
	name      string
	benchType BenchType
	submatch  []string
	results   benchmarkResult
}

func (line *benchmarkRunLine) Parse() error {
	line.applyRegularExpr()
	if line.benchType != "" {
		return line.parseSubmatch()
	}
	return nil
}

func (line *benchmarkRunLine) applyRegularExpr() {
	var submatch []string

	for benchType, reg := range benchmarkResultsRegArray {
		submatch = reg.FindStringSubmatch(line.Output)
		if len(submatch) > 0 {
			line.submatch = submatch
			line.name = submatch[1]
			line.benchType = benchType
			break
		}
	}
}

func (line *benchmarkRunLine) parseSubmatch() error {
	switch line.benchType {
	case RegularBenchmark:
		return line.parseRegularBenchmark()
	case AllocsBenchmark:
		return line.parseAllocsBenchmark()
	case BytesBenchmark:
		return line.parseBytesBenchmark()
	default:
		break
	}
	return errors.New(ErrorLineUnrecognized)
}

// The length of submatch for a RegularBenchmark line is 4.
// Index 0 is the whole line, index 1 is the name of the benchmark,
// index 2 is the number of iterations, index 3 is the ns/op.
func (line *benchmarkRunLine) parseRegularBenchmark() error {
	if len(line.submatch) != 4 {
		return fmt.Errorf("%s: expected 4 arguments but got %d", ErrorLineMalformed, len(line.submatch))
	}

	// get number of ops
	n, err := strconv.Atoi(line.submatch[2])
	if err != nil {
		return fmt.Errorf("%s: %s", ErrorLineMalformed, err.Error())
	}
	line.results.Op = n

	// get ns/op
	nsop, err := strconv.ParseFloat(line.submatch[3], 64)
	if err != nil {
		return fmt.Errorf("%s: %s", ErrorLineMalformed, err.Error())
	}
	line.results.NanosecondPerOp = nsop
	return nil
}

// parse line of type AllocsBenchmark
// The length of submatch for an AllocsBenchmark line is 7.
// The index 0 is the whole line, index 1 is the name of the benchmark,
// index 2 is the number of iterations, index 3 is the ns/op.
//
// The 5th index measures the value of the 4th index.
// For example:
// 		submatch[4] = 1
// 		submatch[5] = "B"
// can be read as "1 Byte".
//
// The 6th index is the number of allocation per operation.
func (line *benchmarkRunLine) parseAllocsBenchmark() error {
	if len(line.submatch) != 7 {
		return fmt.Errorf("%s: expected 7 arguments but got %d", ErrorLineMalformed, len(line.submatch))
	}

	// get number of ops
	n, err := strconv.Atoi(line.submatch[2])
	if err != nil {
		return fmt.Errorf("%s: %s", ErrorLineMalformed, err.Error())
	}
	line.results.Op = n

	// get ns/op
	nsop, err := strconv.ParseFloat(line.submatch[3], 64)
	if err != nil {
		return fmt.Errorf("%s: %s", ErrorLineMalformed, err.Error())
	}
	line.results.NanosecondPerOp = nsop

	// get bytes/op
	bytesPerOp, err := strconv.ParseFloat(line.submatch[4], 64)
	if err != nil {
		return fmt.Errorf("%s: %s", ErrorLineMalformed, err.Error())
	}
	line.results.BytesPerOp = bytesPerOp

	// TODO: handle the 5th index in a better and more dynamic way. Error for now if different than "B" (byte(s)).
	if line.submatch[5] != "B" {
		return fmt.Errorf("%s: index 5 of the line expected %s but got %s", ErrorLineMalformed, "B", line.submatch[5])
	}

	// get allocs
	allocsPerOp, err := strconv.ParseFloat(line.submatch[6], 64)
	if err != nil {
		return fmt.Errorf("%s: %s", ErrorLineMalformed, err.Error())
	}
	line.results.AllocsPerOp = allocsPerOp
	return nil
}

// parse line of type BytesBenchmark
// The length of submatch for a BytesBenchmark line is 6.
// The index 0 is the whole line, index 1 is the name of the benchmark,
// index 2 is the number of iterations, index 3 is the ns/op.
//
// The 5th index measures the value of the 4th index.
// For example:
// 		submatch[4] = 42
// 		submatch[5] = "MB"
// can be read as "42 MB".
func (line *benchmarkRunLine) parseBytesBenchmark() error {
	if len(line.submatch) != 6 {
		return fmt.Errorf("%s: expected 6 arguments but got %d", ErrorLineMalformed, len(line.submatch))
	}

	// get number of ops
	n, err := strconv.Atoi(line.submatch[2])
	if err != nil {
		return fmt.Errorf("%s: %s", ErrorLineMalformed, err.Error())
	}
	line.results.Op = n

	// get ns/op
	nsop, err := strconv.ParseFloat(line.submatch[3], 64)
	if err != nil {
		return fmt.Errorf("%s: %s", ErrorLineMalformed, err.Error())
	}
	line.results.NanosecondPerOp = nsop

	// get MB/S
	mbPerSec, err := strconv.ParseFloat(line.submatch[4], 64)
	if err != nil {
		return fmt.Errorf("%s: %s", ErrorLineMalformed, err.Error())
	}
	line.results.MBs = mbPerSec

	// TODO: handle the 5th index in a better and more dynamic way. Error for now if different than "MB" (byte(s)).
	if line.submatch[5] != "MB" {
		return fmt.Errorf("%s: index 5 of the line expected %s but got %s", ErrorLineMalformed, "MB", line.submatch[5])
	}
	return nil
}

func (line *benchmarkRunLine) InsertToMySQL(microBenchID int64, client *mysql.Client) error {
	query := "INSERT INTO microbenchmark_details(microbenchmark_no, name, bench_type, n, ns_per_op, mb_per_sec, bytes_per_op, allocs_per_op) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
	id, err := client.Insert(query, microBenchID, line.name, line.benchType, line.results.Op, line.results.NanosecondPerOp, line.results.MBs, line.results.BytesPerOp, line.results.AllocsPerOp)
	if err != nil {
		return err
	}
	line.id = id
	return nil
}
