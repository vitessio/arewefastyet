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
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/exec/metrics"
	"github.com/vitessio/arewefastyet/go/storage/psdb"
)

var (
	v = viper.New()

	skip = ""

	dbConfig = psdb.Config{}
	dbClient = new(psdb.Client)
)

func TestMain(m *testing.M) {
	configFile, secretsFile := os.Getenv("CONFIG_FILE_PATH_ENDTOEND"), os.Getenv("SECRETS_FILE_PATH_ENDTOEND")
	if configFile == "" {
		configFile = "../../../../config/dev/config.yaml"
		secretsFile = "../../../../config/dev/secrets.yaml"
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

	_, err = os.OpenFile(secretsFile, os.O_RDWR, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			skip = "no secrets file found."
			os.Exit(m.Run())
		}
		log.Fatal(err)
	}

	v.SetConfigFile(configFile)
	if err = v.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	v.SetConfigFile(secretsFile)
	err = v.MergeInConfig()
	if err != nil {
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

	uuid := "test_TestInsertExecutionMetrics"
	_, err := dbClient.Insert("DELETE FROM metrics WHERE exec_uuid=?", uuid)
	defer func() {
		_, err := dbClient.Insert("DELETE FROM metrics WHERE exec_uuid=?", uuid)
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		return
	}

	cpu := map[string]float64{
		"vtgate":   1987.93,
		"vttablet": 789.12,
	}
	mem := map[string]float64{
		"vtgate":   789,
		"vttablet": 456789,
	}

	execMetrics := metrics.ExecutionMetrics{
		TotalComponentsCPUTime: cpu["vtgate"] + cpu["vttablet"],
		ComponentsCPUTime:      cpu,

		TotalComponentsMemStatsAllocBytes: mem["vtgate"] + mem["vttablet"],
		ComponentsMemStatsAllocBytes:      mem,
	}

	err = metrics.InsertExecutionMetrics(dbClient, uuid, execMetrics)
	c.Assert(err, qt.IsNil)

	result, err := metrics.GetExecutionMetricsSQL(dbClient, uuid)
	c.Assert(err, qt.IsNil)
	c.Assert(result, qt.DeepEquals, execMetrics)
}
