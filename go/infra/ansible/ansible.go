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
	"io"
	"os"
	"path"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ErrorPathUnknown = "path does not exist"

	flagAnsibleRoot   = "ansible-root-directory"
	flagInventoryFile = "ansible-inventory-file"
	flagPlaybookFile  = "ansible-playbook-file"
)

type Config struct {
	RootDir       string
	InventoryFile string
	PlaybookFile  string

	stdout io.Writer
	stderr io.Writer

	// ExtraVars is a key value map that will be passed to Ansible
	// as extra variable using --extra-vars.
	// The corresponding keys are defined as constants in the `vars.go` file of this package.
	// This map gets filled by the Executor before the execution of Ansible.
	// It contains all the required information to run the Ansible roles and tasks properly.
	ExtraVars map[string]interface{}
}

func NewConfig() Config {
	return Config{
		ExtraVars: map[string]interface{}{},
	}
}

func (c *Config) AddExtraVar(key string, value interface{}) {
	c.ExtraVars[key] = value
}

func (c *Config) AddToViper(v *viper.Viper) {
	_ = v.UnmarshalKey(flagAnsibleRoot, &c.RootDir)
	_ = v.UnmarshalKey(flagInventoryFile, &c.InventoryFile)
	_ = v.UnmarshalKey(flagPlaybookFile, &c.PlaybookFile)
}

func (c *Config) AddToPersistentCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.RootDir, flagAnsibleRoot, "", "Root directory of Ansible")
	cmd.Flags().StringVar(&c.InventoryFile, flagInventoryFile, "", "Inventory file used by Ansible")
	cmd.Flags().StringVar(&c.PlaybookFile, flagPlaybookFile, "", "Playbook file used by Ansible")

	_ = viper.BindPFlag(flagAnsibleRoot, cmd.Flags().Lookup(flagAnsibleRoot))
	_ = viper.BindPFlag(flagInventoryFile, cmd.Flags().Lookup(flagInventoryFile))
	_ = viper.BindPFlag(flagPlaybookFile, cmd.Flags().Lookup(flagPlaybookFile))
}

func applyRootToFiles(root string, file *string) {
	if !path.IsAbs(*file) {
		*file = path.Join(root, *file)
	}
}

func Run(c *Config) error {
	applyRootToFiles(c.RootDir, &c.PlaybookFile)
	applyRootToFiles(c.RootDir, &c.InventoryFile)

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		User:          "root",
		SSHCommonArgs: "-o StrictHostKeyChecking=no",
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: c.InventoryFile,
		ExtraVars: c.ExtraVars,
	}

	ansiblePlaybookPrivilegeEscalationOptions := &options.AnsiblePrivilegeEscalationOptions{
		Become: true,
	}

	plb := &playbook.AnsiblePlaybookCmd{
		Playbooks:                  []string{c.PlaybookFile},
		ConnectionOptions:          ansiblePlaybookConnectionOptions,
		PrivilegeEscalationOptions: ansiblePlaybookPrivilegeEscalationOptions,
		Options:                    ansiblePlaybookOptions,
		Exec: execute.NewDefaultExecute(
			execute.WithShowDuration(),
			execute.WithWrite(c.stdout),
			execute.WithWriteError(c.stderr),
		),
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

func (c *Config) SetStdout(stdout *os.File) {
	c.stdout = stdout
}

func (c *Config) SetStderr(stderr *os.File) {
	c.stderr = stderr
}

func (c *Config) SetOutputs(stdout, stderr *os.File) {
	c.stdout = stdout
	c.stderr = stderr
}
