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

package endtoend

import (
	qt "github.com/frankban/quicktest"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/exec/metrics"
	"github.com/vitessio/arewefastyet/go/storage/psdb"
	"log"
	"math"
	"os"
	"testing"
)

var (
	v = viper.New()

	skip = ""

	dbConfig = psdb.Config{}
	dbClient = new(psdb.Client)
)

func TestMain(m *testing.M) {
	configFile := os.Getenv("CONFIG_FILE_PATH_ENDTOEND")
	if configFile == "" {
		configFile = "../../../../config/config.yaml"
	}

	// Checking if the configuration file exists or not.
	// Skipping the E2E tests if none were found.
	//
	// CI does not include configuration file / secrets, thus
	// these tests will be skipped.
	_, err := os.OpenFile(configFile, os.O_RDWR, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			skip = "no configuration file found."
			os.Exit(m.Run())
		}
		log.Fatal(err)
	}

	v.SetConfigFile(configFile)
	if err = v.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	dbConfig.AddToViper(v)
	dbClient, err = dbConfig.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestInsertExecutionMetrics(t *testing.T) {
	c := qt.New(t)
	if skip != "" {
		c.Skip(skip)
	}

	_, err := dbClient.Insert("DELETE FROM metrics WHERE exec_uuid='test_TestInsertExecutionMetrics'")
	if err != nil {
		return
	}

	cpu := map[string]float64{
		"vtgate":   111987.93,
		"vttablet": 789.12,
	}
	mem := map[string]float64{
		"vtgate":   789.15,
		"vttablet": 456789.3,
	}

	execMetrics := metrics.ExecutionMetrics{
		TotalComponentsCPUTime: cpu["vtgate"] + cpu["vttablet"],
		ComponentsCPUTime:      cpu,

		TotalComponentsMemStatsAllocBytes: mem["vtgate"] + mem["vttablet"],
		ComponentsMemStatsAllocBytes:      mem,
	}

	err = metrics.InsertExecutionMetrics(dbClient, "test_TestInsertExecutionMetrics", execMetrics)
	c.Assert(err, qt.IsNil)

	selectQ := "select name,value from metrics where exec_uuid='test_TestInsertExecutionMetrics'"
	rows, err := dbClient.Select(selectQ)
	c.Assert(err, qt.IsNil)

	res := []float64{
		cpu["vtgate"] + cpu["vttablet"],
		mem["vtgate"] + mem["vttablet"],
		cpu["vtgate"],
		cpu["vttablet"],
		mem["vtgate"],
		mem["vttablet"],
	}
	for i := 0; rows.Next(); i++ {
		var name string
		var value float64
		err = rows.Scan(&name, &value)
		c.Assert(err, qt.IsNil)
		diff := math.Abs(value - res[i])
		c.Assert(diff < 1, qt.IsTrue)
	}
}
