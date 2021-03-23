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

package ansible

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagInventoryFiles = "ansible-inventory-files"
	flagPlaybookFiles = "ansible-playbook-files"
)

type AnsibleConfig struct {
	InventoryFiles []string
	PlaybookFiles  []string
}

func (a *AnsibleConfig) AddToPersistentCommand(cmd *cobra.Command) {
	cmd.PersistentFlags().StringSliceVar(&a.InventoryFiles, flagInventoryFiles, []string{}, "List of inventory files used by Ansible")
	cmd.PersistentFlags().StringSliceVar(&a.PlaybookFiles, flagPlaybookFiles, []string{}, "List of playbook files used by Ansible")

	viper.BindPFlag(flagInventoryFiles, cmd.Flags().Lookup(flagInventoryFiles))
	viper.BindPFlag(flagPlaybookFiles, cmd.Flags().Lookup(flagPlaybookFiles))
}