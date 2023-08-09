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

package api

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vitessio/arewefastyet/go/server"
)

func ApiCmd() *cobra.Command {
	var srv server.Server

	cmd := &cobra.Command{
		Use:   "api",
		Short: "Starts the api server of arewefastyet and the CRON service",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("RunE")
			err := srv.Init()
			if err != nil {
				log.Println(err)
				return err
			}
			log.Println("Srv Run")
			return srv.Run()
		},
	}

	srv.AddToCommand(cmd)

	return cmd
}
