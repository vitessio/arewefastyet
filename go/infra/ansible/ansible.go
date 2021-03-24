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
	"context"
	"errors"
	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

const (
	ErrorPathUnknown = "path does not exist"

	flagAnsibleRoot    = "ansible-root-directory"
	flagInventoryFiles = "ansible-inventory-files"
	flagPlaybookFiles  = "ansible-playbook-files"
)

type Config struct {
	RootDir        string
	InventoryFiles []string
	PlaybookFiles  []string
}

func (c *Config) AddToPersistentCommand(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&c.RootDir, flagAnsibleRoot, "", "Root directory of Ansible")
	cmd.PersistentFlags().StringSliceVar(&c.InventoryFiles, flagInventoryFiles, []string{}, "List of inventory files used by Ansible")
	cmd.PersistentFlags().StringSliceVar(&c.PlaybookFiles, flagPlaybookFiles, []string{}, "List of playbook files used by Ansible")

	viper.BindPFlag(flagAnsibleRoot, cmd.Flags().Lookup(flagAnsibleRoot))
	viper.BindPFlag(flagInventoryFiles, cmd.Flags().Lookup(flagInventoryFiles))
	viper.BindPFlag(flagPlaybookFiles, cmd.Flags().Lookup(flagPlaybookFiles))
}

func applyRootToFiles(root string, files *[]string) {
	for i, file := range *files {
		if path.IsAbs(file) == false {
			(*files)[i] = path.Join(root, file)
		}
	}
}

func inventoryFilesToString(invFiles []string) string {
	var res string
	for i, inv := range invFiles {
		if i > 0 {
			res = res + ", "
		}
		res = res + inv
	}
	return res
}

func Run(c *Config) error {
	applyRootToFiles(c.RootDir, &c.PlaybookFiles)
	applyRootToFiles(c.RootDir, &c.InventoryFiles)

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		User:          "root",
		SSHCommonArgs: "-o StrictHostKeyChecking=no",
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: inventoryFilesToString(c.InventoryFiles),
	}

	ansiblePlaybookPrivilegeEscalationOptions := &options.AnsiblePrivilegeEscalationOptions{
		Become: true,
	}

	plb := &playbook.AnsiblePlaybookCmd{
		Playbooks:                  c.PlaybookFiles,
		ConnectionOptions:          ansiblePlaybookConnectionOptions,
		PrivilegeEscalationOptions: ansiblePlaybookPrivilegeEscalationOptions,
		Options:                    ansiblePlaybookOptions,
		Exec:                       execute.NewDefaultExecute(execute.WithShowDuration()),
	}

	err := plb.Run(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) CopyRootDirectory(directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return errors.New(ErrorPathUnknown)
	}

	err := copy.Copy(c.RootDir, directory)
	if err != nil {
		return err
	}
	c.RootDir = directory
	return nil
}
