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

package metrics

import (
	"github.com/spf13/cobra"
	"github.com/vitessio/arewefastyet/go/exec"
	"github.com/vitessio/arewefastyet/go/exec/metrics"
	"github.com/vitessio/arewefastyet/go/storage/influxdb"
	"github.com/vitessio/arewefastyet/go/storage/psdb"
	"log"
)

func GenExecMetricsCmd() *cobra.Command {
	dbConfig := &psdb.Config{}
	metricsDBConfig := &influxdb.Config{}

	cmd := &cobra.Command{
		Use:     "exec_metrics",
		Short:   "For each execution, fetches the metrics from influxDB and store them to SQL if not already present.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientSQL, err := dbConfig.NewClient()
			if err != nil {
				return err
			}
			clientMetrics, err := metricsDBConfig.NewClient()
			if err != nil {
				return err
			}

			rowsExecUUIDs, err := clientSQL.Select("select uuid from execution where status = ?", exec.StatusFinished)
			if err != nil {
				return err
			}
			defer rowsExecUUIDs.Close()

			for rowsExecUUIDs.Next() {
				var uuid string
				err = rowsExecUUIDs.Scan(&uuid)
				if err != nil {
					return err
				}

				rowsExecMetrics, err := clientSQL.Select("select id from metrics where exec_uuid = ? limit 1", uuid)
				if err != nil {
					return err
				}
				exist := rowsExecMetrics.Next()
				rowsExecMetrics.Close()
				if exist {
					continue
				}

				executionMetrics, err := metrics.GetExecutionMetrics(*clientMetrics, uuid)
				if err != nil {
					return err
				}

				log.Println("found metrics for:", uuid)

				err = metrics.InsertExecutionMetrics(clientSQL, uuid, executionMetrics)
				if err != nil {
					return err
				}
				log.Println("inserted new metric for:", uuid)
			}
			return
		},
	}

	dbConfig.AddToCommand(cmd)
	metricsDBConfig.AddToCommand(cmd)
	return cmd
}
