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

package exec

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vitessio/arewefastyet/go/exec"
)

func ExecCmd() *cobra.Command {
	ex, err := exec.NewExec()
	if err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:     "exec",
		Aliases: []string{"e"},
		Short:   "Execute a task",
		Long: `Execute a task based on the given terraform and ansible configuration.
It handles the creation, configuration, and cleanup of the infrastructure.`,
		Example: `arewefastyet exec --exec-git-ref 4a70d3d226113282554b393a97f893d133486b94  --db-database benchmark --db-host localhost --db-password <db-password> --db-user <db-username>
--exec-source config_micro_remote --ansible-inventory-files microbench_inventory.yml --ansible-playbook-files microbench.yml --ansible-root-directory ./ansible/
--equinix-instance-type m2.xlarge.x86 --equinix-token tok --equinix-project-id id
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// prepare
			if err := ex.Prepare(); err != nil {
				return err
			}

			// execute
			if err := ex.Execute(); err != nil {
				return err
			}

			// cleanup
			if err := ex.CleanUp(); err != nil {
				return err
			}
			return nil
		},
	}

	ex.AddToCommand(cmd)
	return cmd
}
