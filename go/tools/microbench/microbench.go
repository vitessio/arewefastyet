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

package microbench

import (
	"encoding/json"
	"fmt"
	"github.com/vitessio/arewefastyet/go/mysql"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type BenchType string

const (
	RegularBenchmark = BenchType("regular")
	AllocsBenchmark  = BenchType("allocs")
	BytesBenchmark   = BenchType("bytes")
)

var benchmarkResultsRegArray = map[BenchType]*regexp.Regexp{
	AllocsBenchmark:  regexp.MustCompile(`(Benchmark.+\b) \s*([0-9]+)\s+([\d\.]+) ns\/op\s+([\d\.]+) (.+)\/op\s+([\d\.]+) allocs\/op`),
	BytesBenchmark:   regexp.MustCompile(`(Benchmark.+\b) \s*([0-9]+)\s+([\d\.]+) ns\/op\s+([\d\.]+) (.+)\/s`),
	RegularBenchmark: regexp.MustCompile(`(Benchmark.+\b) \s*([0-9]+)\s+([\d\.]+) ns/op`),
}

type benchmark struct {
	id                      int64
	filePath, pkgName, name string
}

func (b *benchmark) RegisterToMySQL(client *mysql.Client) error {
	query := "INSERT INTO microbenchmark(test_no, pkg_name, name) VALUES(?, ?, ?)"
	id, err := client.Insert(query, 0, b.pkgName, b.name)
	if err != nil {
		return err
	}
	b.id = id
	return nil
}

// MicroBenchmark runs "go test bench" on the given package (pkg) and outputs
// the results to outputPath.
// Profiling files will be written to the current working directory.
func MicroBenchmark(cfg MicroBenchConfig) {
	var sqlClient *mysql.Client
	var err error

	if cfg.DatabaseConfig.IsValid() {
		sqlClient, err = mysql.New(*cfg.DatabaseConfig)
		if err != nil {
			log.Fatal(err)
		}
		defer sqlClient.Close()
	}

	loaded, err := packages.Load(&packages.Config{
		Mode:  packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedImports | packages.NeedModule,
		Tests: true,
	}, cfg.Package)
	if err != nil {
		panic(err)
	}

	w, err := os.Create(cfg.Output)
	if err != nil {
		panic(err)
	}

	benchmarks := findBenchmarks(loaded)

	for _, benchmark := range benchmarks {
		idx := strings.LastIndex(benchmark.filePath, "/")
		path := benchmark.filePath[:idx+1]
		command := exec.Command("go", "test", "-bench=^"+benchmark.name+"$", "-run==", "-json", "-count=10", path)
		b, err := command.Output()

		if err == nil {
			if sqlClient != nil {
				if err := benchmark.RegisterToMySQL(sqlClient); err != nil {
					log.Fatal(err)
				}
			}

			lines := strings.Split(string(b), "\n")
			for _, line := range lines {
				var benchLine benchmarkRunLine
				err := json.Unmarshal([]byte(line), &benchLine)
				if err != nil || benchLine.Output == "" {
					continue
				}

				err = benchLine.Parse()
				if err != nil {
					log.Fatal(err)
				}

				if benchLine.benchType != "" {
					fmt.Printf("%s - %s: %s ns/op\n", benchmark.pkgName, benchLine.submatch[1], benchLine.submatch[3])
					fmt.Fprintf(w, "%s - %s: %s ns/op\n", benchmark.pkgName, benchLine.submatch[1], benchLine.submatch[3])
					err = benchLine.InsertToMySQL(benchmark.id, sqlClient)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		} else {
			fmt.Println(err.Error())
		}

		cpuprofFile := fmt.Sprintf("cpuprof_%s.%s.out", benchmark.pkgName, benchmark.name)
		command = exec.Command("go", "test", "-bench=^"+benchmark.name+"$", "-run==", "-count=1", path, fmt.Sprintf("-cpuprofile=%s", cpuprofFile))
		b, err = command.Output()
		if err == nil {
			fmt.Printf("CPU profile generated %s\n", cpuprofFile)
			fmt.Fprintf(w, "CPU profile generated %s\n", cpuprofFile)

			// TODO: insert to MySQL
		}
	}
}

func findBenchmarks(loaded []*packages.Package) []benchmark {
	var benchmarks []benchmark
	for _, pkg := range loaded {
		scope := pkg.Types.Scope()
		for _, typName := range scope.Names() {
			f, ok := scope.Lookup(typName).(*types.Func)
			if ok && isBenchmark(f) {
				fs := pkg.Fset.File(f.Pos())
				benchmarks = append(benchmarks, benchmark{
					pkgName:  f.Pkg().Name(),
					name:     f.Name(),
					filePath: fs.Name(),
				})
			}
		}
	}
	return benchmarks
}

func isBenchmark(f *types.Func) bool {
	return strings.HasPrefix(f.Name(), "Bench") && f.Type().String() == "func(b *testing.B)"
}
