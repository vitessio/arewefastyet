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

package infra

import (
	"context"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/infra/ansible"
)

const (
	ErrorInvalidConfiguration = "invalid configuration"
	ErrorProvision            = "provision failed"
)

type Infra interface {
	AddToCommand(cmd *cobra.Command)
	AddToViper(v *viper.Viper)
	CleanUp() error
	Create(wantOutputs ...string) (output map[string]string, err error)
	ValidConfig() error
	Prepare() error
	TerraformVarArray() (vars []*tfexec.VarOption)
	Run(ansible *ansible.Config) error
	SetConfig(config *Config)
	SetTags(tags map[string]string)
	SetExecUUID(uuid uuid.UUID)
}

func getTerraformExecPath(installPath string) (string, error) {
	execPath, err := tfinstall.Find(context.Background(), tfinstall.LatestVersion(installPath, false))
	if err != nil {
		return "", err
	}
	return execPath, nil
}
