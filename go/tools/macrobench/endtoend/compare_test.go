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

	qt "github.com/frankban/quicktest"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/storage/influxdb"
	"github.com/vitessio/arewefastyet/go/storage/psdb"
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
)

var (
	v = viper.New()

	skip = ""

	dbConfig = psdb.Config{}
	dbClient = new(psdb.Client)

	execMetricsConfig = influxdb.Config{}
	execMetricsClient = new(influxdb.Client)
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

	execMetricsConfig.AddToViper(v)
	execMetricsClient, err = execMetricsConfig.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func TestCompareMacroBenchmark(t *testing.T) {
	c := qt.New(t)
	if skip != "" {
		c.Skip(skip)
	}
	types := []string{"OLTP", "TPCC"}
	macrosMatrices, err := macrobench.CompareMacroBenchmarks(dbClient, "dff8d632908583cae5940b25a962eaa2e6550508", "f7304cd1893accfefee0525910098a8e0e68deec", macrobench.V3Planner, types)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(macrosMatrices, qt.HasLen, len(types))
}

func BenchmarkCompareMacroBenchmark(b *testing.B) {
	c := qt.New(b)

	run := func(b *testing.B, planner macrobench.PlannerVersion, reference, compare string) {
		if skip != "" {
			c.Skip(skip)
		}
		b.ReportAllocs()
		types := []string{"OLTP", "TPCC"}
		for i := 0; i < b.N; i++ {
			macrosMatrices, err := macrobench.CompareMacroBenchmarks(dbClient, reference, compare, planner, types)
			if err != nil {
				c.Fatal(err)
			}
			c.Assert(macrosMatrices, qt.HasLen, len(types))
		}
	}

	b.Run("Gen4 planner", func(b *testing.B) {
		run(b, macrobench.Gen4FallbackPlanner, "48dccf56282dc79903c0ab0b1d0177617f927403", "f7304cd1893accfefee0525910098a8e0e68deec")
	})
	b.Run("V3 planner", func(b *testing.B) {
		run(b, macrobench.V3Planner, "48dccf56282dc79903c0ab0b1d0177617f927403", "f7304cd1893accfefee0525910098a8e0e68deec")
	})
}
