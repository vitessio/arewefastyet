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

package equinix

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/infra"
)

const (
	flagToken     = "equinix-token"
	flagProjectID = "equinix-project-id"
)

type Equinix struct {
	Token     string
	ProjectID string
	InfraCfg  *infra.Config
}

func (e *Equinix) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&e.Token, flagToken, "", "Auth Token for Equinix Metal")
	cmd.Flags().StringVar(&e.ProjectID, flagProjectID, "", "Project ID to use for Equinix Metal")

	viper.BindPFlag(flagToken, cmd.Flags().Lookup(flagToken))
	viper.BindPFlag(flagProjectID, cmd.Flags().Lookup(flagProjectID))
}

func (e Equinix) Create() error {
	if err := e.ValidConfig(); err != nil {
		return err
	}
	// create
	return nil
}

func (e Equinix) ValidConfig() error {
	if e.Token == "" {
		return fmt.Errorf("%s: missing token", infra.ErrorInvalidConfiguration)
	} else if e.ProjectID == "" {
		return fmt.Errorf("%s: missing project id", infra.ErrorInvalidConfiguration)
	} else if err := e.InfraCfg.Valid(); err != nil {
		return err
	}
	return nil
}

func (e *Equinix) Prepare() error {
	return nil
}

func (e *Equinix) Run() error {
	return nil
}
