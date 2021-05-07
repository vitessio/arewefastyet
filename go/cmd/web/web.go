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

package web

import (
	"github.com/spf13/cobra"
	"github.com/vitessio/arewefastyet/go/server"
)

func WebCmd() *cobra.Command {
	var srv server.Server

	cmd := &cobra.Command{
		Use:   "web",
		Short: "Starts the HTTP web server",
		Long: `web command starts the HTTP web server, the credentials for which are provided as command line arguments or the configuration file. 
It uses a MySQL and InfluxDB instances to read the metrics that will be displayed. It has an interface for macrobenchmark and microbenchmark runs.`,
		Example: `arewefastyet web --db-database benchmark --db-host localhost --db-password <db-password> --db-user <db-username>  
--influx-database benchmark-influx --influx-hostname localhost --influx-password <influx-password>
--influx-port <influx-port> --influx-username <influx-username> --web-api-key <web-api-key>
--web-mode production --web-port <web-port> --web-static-path ./server/static --web-template-path ./server/template
--web-webhook-config ./config.yaml`,
		Aliases: []string{"w"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return srv.Run()
		},
	}

	srv.AddToCommand(cmd)

	return cmd
}
