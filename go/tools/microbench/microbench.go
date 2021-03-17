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
	"errors"
	"fmt"
	"github.com/vitessio/arewefastyet/go/mysql"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	errorInvalidProfileType = "invalid profile type"

	profileCPU = "cpu"
	profileMem = "mem"
)

type benchmark struct {
	id               int64
	filePath         string
	name             string
	pkgPath, pkgName string
	sql              *mysql.Client
}

func (b *benchmark) registerToMySQL(client *mysql.Client) error {
	query := "INSERT INTO microbenchmark(test_no, pkg_name, name) VALUES(?, ?, ?)"
	id, err := client.Insert(query, 0, b.pkgName, b.name)
	if err != nil {
		return err
	}
	b.id = id
	return nil
}

func (b *benchmark) execute(rootDir string, w *os.File) error {
	command := exec.Command("go", "test", "-bench=^"+b.name+"$", "-run==", "-json", "-count=10", b.pkgPath)
	command.Dir = rootDir
	out, err := command.Output()

	if err != nil {
		return err
	}

	if b.sql != nil {
		if err := b.registerToMySQL(b.sql); err != nil {
			return err
		}
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		var benchLine benchmarkRunLine
		err := json.Unmarshal([]byte(line), &benchLine)
		if err != nil || benchLine.Output == "" {
			continue
		}

		err = benchLine.Parse()
		if err != nil {
			return err
		}

		if benchLine.benchType != "" {
			fmt.Printf("%s - %s %f ns/op\n", b.pkgName, benchLine.name, benchLine.results.NanosecondPerOp)
			fmt.Fprintf(w, "%s - %s %f ns/op\n", b.pkgName, benchLine.name, benchLine.results.NanosecondPerOp)
			if b.sql != nil {
				err = benchLine.InsertToMySQL(b.id, b.sql)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (b benchmark) executeProfile(rootDir, profileType string, w *os.File) error {
	if profileType != profileCPU && profileType != profileMem {
		return errors.New(errorInvalidProfileType)
	}
	profileName := fmt.Sprintf("%sprof_%s.%s.out", profileType, b.pkgName, b.name)
	command := exec.Command("go", "test", "-bench=^"+b.name+"$", "-run==", "-count=1", b.pkgPath, fmt.Sprintf("-%sprofile=%s", profileType, profileName))
	command.Dir = rootDir

	_, err := command.Output()
	if err != nil {
		return err
	}
	fmt.Printf("CPU profile generated %s\n", profileName)
	fmt.Fprintf(w, "CPU profile generated %s\n", profileName)
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
		Dir:   cfg.RootDir,
	}, cfg.Package)
	if err != nil {
		panic(err)
	}

	w, err := os.Create(cfg.Output)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	benchmarks := findBenchmarks(loaded)

	for _, benchmark := range benchmarks {
		benchmark.sql = sqlClient

		fmt.Println(benchmark.pkgPath)

		err := benchmark.execute(cfg.RootDir, w)
		if err != nil {
			fmt.Println(err.Error())
		}

		profiles := []string{profileMem, profileCPU}
		for _, profile := range profiles {
			err = benchmark.executeProfile(cfg.RootDir, profile, w)
			if err != nil && err.Error() != errorInvalidProfileType {
				fmt.Println(err.Error())
			}
		}
		fmt.Println()
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
					pkgPath:  f.Pkg().Path(),
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
