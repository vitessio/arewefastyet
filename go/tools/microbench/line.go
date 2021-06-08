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
	"github.com/vitessio/arewefastyet/go/storage"
	"regexp"
	"strconv"
	"time"
)

type microType string

const (
	ErrorLineUnrecognized = "the line format was unrecognized"
	ErrorLineMalformed    = "the format of the line is malformed"

	// GeneralBenchmark are results that fits with the generalBenchTypeRegExpr regular expression.
	GeneralBenchmark = microType("general")
)

// benchmarkResultsRegexp's submatch array contains 7 indexes defined below:
//
// submatch[0]:  whole string
// submatch[1]:  benchmark name
// submatch[2]:  number of ops (number)
// submatch[3]:  ns/op (number)
// submatch[4]:  MB/s (number)
// submatch[5]:  B/op (number)
// submatch[6]:  allocs/op (number)
//
// https://regex101.com/r/JdCpno/3
var generalBenchTypeRegExpr = regexp.MustCompile(`(Benchmark.+\b)\s*([0-9]+)\s+([\d\.]+)\s+ns\/op(?:\s*)?(?:([\d\.]+) MB/s+\s*)?(?:([\d\.]+) B/op+\s*)?(?:([0-9]+) allocs/op\s*)?\n`)

var benchTypeRegExprs = map[microType]*regexp.Regexp{
	GeneralBenchmark: generalBenchTypeRegExpr,
}

type lineResult struct {
	Op              int
	NanosecondPerOp float64
	MBs             float64
	BytesPerOp      float64
	AllocsPerOp     float64
}

type lineRun struct {
	Time    time.Time
	Action  string
	Package string
	Output  string
	Elapsed string

	id        int64
	name      string
	benchType microType
	submatch  []string
	results   lineResult
}

func (line *lineRun) Parse() error {
	line.applyRegularExpr()
	if line.benchType != "" {
		return line.parseSubmatch()
	}
	return nil
}

func (line *lineRun) applyRegularExpr() {
	var submatch []string

	for benchType, r := range benchTypeRegExprs {
		submatch = r.FindStringSubmatch(line.Output)
		if len(submatch) > 0 {
			line.submatch = submatch
			line.name = submatch[1]
			line.benchType = benchType
			break
		}
	}
}

func (line *lineRun) parseSubmatch() error {
	switch line.benchType {
	case GeneralBenchmark:
		return line.parseGeneralBenchmark()
	default:
		break
	}
	return errors.New(ErrorLineUnrecognized)
}

func getFloatSubMatch(parent, child string) (float64, error) {
	if len(parent) == 0 {
		return 0, nil
	}
	val, err := strconv.ParseFloat(child, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (line *lineRun) parseGeneralBenchmark() error {
	if len(line.submatch) != 7 {
		return fmt.Errorf("%s: expected 7 arguments but got %d", ErrorLineMalformed, len(line.submatch))
	}

	var err error

	// get number of ops
	line.results.Op, err = strconv.Atoi(line.submatch[2])
	if err != nil {
		return fmt.Errorf("%s: %s", ErrorLineMalformed, err.Error())
	}

	// get ns/op
	line.results.NanosecondPerOp, err = strconv.ParseFloat(line.submatch[3], 64)
	if err != nil {
		return fmt.Errorf("%s: %s", ErrorLineMalformed, err.Error())
	}

	line.results.MBs, err = getFloatSubMatch(line.submatch[4], line.submatch[4])
	if err != nil {
		return fmt.Errorf("%s: %s", ErrorLineMalformed, err.Error())
	}

	line.results.BytesPerOp, err = getFloatSubMatch(line.submatch[5], line.submatch[5])
	if err != nil {
		return fmt.Errorf("%s: %s", ErrorLineMalformed, err.Error())
	}

	line.results.AllocsPerOp, err = getFloatSubMatch(line.submatch[6], line.submatch[6])
	if err != nil {
		return fmt.Errorf("%s: %s", ErrorLineMalformed, err.Error())
	}
	return nil
}

func (line *lineRun) InsertToMySQL(microBenchID int64, client storage.SQLClient) error {
	query := "INSERT INTO microbenchmark_details(microbenchmark_no, name, bench_type, n, ns_per_op, mb_per_sec, bytes_per_op, allocs_per_op) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
	res, err := client.Insert(query, microBenchID, line.name, line.benchType, line.results.Op, line.results.NanosecondPerOp, line.results.MBs, line.results.BytesPerOp, line.results.AllocsPerOp)
	if err != nil {
		return err
	}
	line.id = res
	return nil
}
