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

package report

import (
	"github.com/vitessio/arewefastyet/go/storage/psdb"
	"log"

	"github.com/spf13/cobra"
	"github.com/vitessio/arewefastyet/go/storage/influxdb"
	"github.com/vitessio/arewefastyet/go/tools/report"
)

func GenerateReport() *cobra.Command {
	var reportFile string
	var toSHA string
	var fromSHA string
	var dbConfig psdb.Config
	var metricsCfg influxdb.Config
	cmd := &cobra.Command{
		Use:     "report",
		Short:   "Generate comparison between two sha commits of Vitess",
		Long:    "Command to generate a pdf-report that compares the microbenchmark and macrobenchmark runs between the two sha commits of Vitess provided",
		Example: "arewefastyet gen report --compare-from sha1 --compare-to sha2 --report-file report.pdf",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("Generating file.")
			client, err := dbConfig.NewClient()
			if err != nil {
				return err
			}

			metricsClient, err := metricsCfg.NewClient()
			if err != nil {
				return err
			}
			err = report.GenerateCompareReport(client, metricsClient, fromSHA, toSHA, reportFile)
			if err != nil {
				return err
			}
			log.Println("Report generated successfully.")
			return nil
		},
	}
	cmd.Flags().StringVar(&reportFile, "report-file", "./report.pdf", "File created that stores the report.")
	cmd.Flags().StringVar(&toSHA, "compare-to", "", "SHA for Vitess that we want to compare to")
	_ = cmd.MarkFlagRequired("compare-to")
	cmd.Flags().StringVar(&fromSHA, "compare-from", "", "SHA for Vitess that we want to compare from")
	_ = cmd.MarkFlagRequired("compare-from")
	dbConfig.AddToCommand(cmd)
	metricsCfg.AddToCommand(cmd)
	return cmd
}
