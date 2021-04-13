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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagRootExec = "exec-root-dir"
)

func (e *Exec) AddToViper(v *viper.Viper) (err error) {
	err = v.UnmarshalKey(flagRootExec, &e.rootDir)
	if err != nil {
		return err
	}
	e.AnsibleConfig.AddToViper(v)
	e.InfraConfig.AddToViper(v)
	e.Infra.AddToViper(v)
	return nil
}

func (e *Exec) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&e.rootDir, flagRootExec, "", "Path to the root directory of exec")
	_ = viper.BindPFlag(flagRootExec, cmd.Flags().Lookup(flagRootExec))

	e.AnsibleConfig.AddToPersistentCommand(cmd)
	e.InfraConfig.AddToPersistentCommand(cmd)
	e.Infra.AddToCommand(cmd)
	e.configDB.AddToCommand(cmd)
}