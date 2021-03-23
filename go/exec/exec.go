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
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/infra"
	"github.com/vitessio/arewefastyet/go/infra/ansible"
	"github.com/vitessio/arewefastyet/go/infra/construct"
	"github.com/vitessio/arewefastyet/go/infra/equinix"
)

type Exec struct {
	UUID          uuid.UUID
	InfraConfig   infra.Config
	AnsibleConfig ansible.Config
	Infra         infra.Infra

	rootDir string
	dirPath string
}

func (e *Exec) Prepare() error {
	err := e.prepareDirectories()
	if err != nil {
		return err
	}

	IPs, err := provision(e.Infra)
	if err != nil {
		return err
	}

	err = ansible.AddIPsToFiles(IPs, e.AnsibleConfig)
	if err != nil {
		return err
	}
	err = ansible.AddLocalConfigPathToFiles(viper.ConfigFileUsed(), e.AnsibleConfig)
	if err != nil {
		return err
	}
	return nil
}

func (e *Exec) Execute() error {
	err := e.Infra.Run(&e.AnsibleConfig)
	if err != nil {
		return err
	}
	return nil
}

func NewExec() (*Exec, error) {
	// todo: dynamic choice for infra provider
	inf, err := construct.NewInfra(equinix.Name)
	if err != nil {
		return nil, err
	}

	ex := Exec{
		UUID: uuid.New(),
		Infra: inf,
	}

	ex.Infra.SetConfig(&ex.InfraConfig)

	return &ex, nil
}

