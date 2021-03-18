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
	"errors"
	"io/ioutil"
	"os"
)

const (
	ErrorPathUnknown = "path does not exist"
	ErrorPathMissing = "path is missing"
)

type Config struct {
	Path       string
	PathExecTF string

	pathInstallTF string
}

func (c Config) Valid() error {
	if c.Path == "" {
		return errors.New(ErrorPathMissing)
	} else if _, err := os.Stat(c.Path); os.IsNotExist(err) {
		return errors.New(ErrorPathUnknown)
	}
	return nil
}

func (c *Config) Prepare() error {
	pathInstallTF, err := ioutil.TempDir("", "tfinstall")
	if err != nil {
		return err
	}
	c.pathInstallTF = pathInstallTF
	pathExecTF, err := getTerraformExecPath(c.pathInstallTF)
	if err != nil {
		return err
	}
	c.PathExecTF = pathExecTF
	return nil
}

func (c *Config) Close() error {
	return os.RemoveAll(c.pathInstallTF)
}
