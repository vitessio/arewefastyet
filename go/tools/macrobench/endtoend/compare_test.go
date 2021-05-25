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
	"github.com/vitessio/arewefastyet/go/storage/influxdb"
	"github.com/vitessio/arewefastyet/go/storage/mysql"
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
	"log"
	"os"
	"testing"
)

var (
	v = viper.New()

	dbConfig = mysql.ConfigDB{}
	dbClient = new(mysql.Client)

	execMetricsConfig = influxdb.Config{}
	execMetricsClient = new(influxdb.Client)
)

func TestMain(m *testing.M) {
	configFile := os.Getenv("CONFIG_FILE_PATH_ENDTOEND")
	if configFile == "" {
		configFile = "../../../../config/config.yaml"
	}
	var err error

	v.SetConfigFile(configFile)
	if err = v.ReadInConfig(); err != nil {
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
	macrosMatrices, err := macrobench.CompareMacroBenchmarks(dbClient, execMetricsClient, "dff8d632908583cae5940b25a962eaa2e6550508", "f7304cd1893accfefee0525910098a8e0e68deec", macrobench.V3Planner)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(macrosMatrices, qt.HasLen, len(macrobench.Types))
}

func BenchmarkCompareMacroBenchmark(b *testing.B) {
	c := qt.New(b)

	run := func(b *testing.B, planner macrobench.PlannerVersion, reference, compare string) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			macrosMatrices, err := macrobench.CompareMacroBenchmarks(dbClient, execMetricsClient, reference, compare, planner)
			if err != nil {
				c.Fatal(err)
			}
			c.Assert(macrosMatrices, qt.HasLen, len(macrobench.Types))
		}
	}

	b.Run("Gen4 planner", func(b *testing.B) {
		run(b, macrobench.Gen4FallbackPlanner, "48dccf56282dc79903c0ab0b1d0177617f927403", "f7304cd1893accfefee0525910098a8e0e68deec")
	})
	b.Run("V3 planner", func(b *testing.B) {
		run(b, macrobench.V3Planner, "48dccf56282dc79903c0ab0b1d0177617f927403", "f7304cd1893accfefee0525910098a8e0e68deec")
	})
}
