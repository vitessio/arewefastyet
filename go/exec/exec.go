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
	"io"
	"os"
)

const (
	// stderrFile = "exec-stderr.log"
	// stdoutFile = "exec-stdout.log"
)

type Exec struct {
	UUID          uuid.UUID
	InfraConfig   infra.Config
	AnsibleConfig ansible.Config
	Infra         infra.Infra

	rootDir string
	dirPath string

	stdout io.Writer
	stderr io.Writer
}

// SetStdout sets the standard output of Exec.
func (e *Exec) SetStdout(stdout io.Writer) {
	e.stdout = stdout
}

// SetStderr sets the standard error output of Exec.
func (e *Exec) SetStderr(stderr io.Writer) {
	e.stderr = stderr
}

// SetOutputToDefaultPath sets both outputs to their default files (stdoutFile and
// stderrFile). If they can't be found in exec.dirPath, they will be created there.
func (e Exec) SetOutputToDefaultPath() error {
	return nil
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

// CleanUp cleans and removes all things required only during the execution flow
// and not after it is done.
func (e Exec) CleanUp() error {
	err := e.Infra.CleanUp()
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
		UUID:  uuid.New(),
		Infra: inf,

		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	ex.Infra.SetConfig(&ex.InfraConfig)

	return &ex, nil
}

// NewExecWithConfig will create a new Exec using the NewExec method, and will
// use viper.Viper to apply the configuration located at pathConfig.
func NewExecWithConfig(pathConfig string) (*Exec, error) {
	e, err := NewExec()
	if err != nil {
		return nil, err
	}
	v := viper.New()

	v.SetConfigFile(pathConfig)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	err = e.AddToViper(v)
	if err != nil {
		return nil, err
	}

	return e, nil
}
