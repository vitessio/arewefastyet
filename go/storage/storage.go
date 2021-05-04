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

package storage

import (
	"errors"
	"github.com/spf13/cobra"
)

const (
	ErrorInvalidConfiguration = "invalid configuration"
)

type Insertion interface {
	Empty() bool
}

type Selection interface {
	Empty() bool
}

type Client interface {
	Close() error
	Insert(query string, args ...interface{}) (Insertion, error)
	Select(query string, args ...interface{}) (Selection, error)
}

type Configuration interface {
	AddToCommand(cmd *cobra.Command)
	NewClient() (Client, error)
	IsValid() bool
}

func Create(config Configuration) (client Client, err error) {
	if config == nil {
		return
	}

	if config.IsValid() {
		client, err = config.NewClient()
		if err != nil {
			return
		}
	} else {
		return nil, errors.New(ErrorInvalidConfiguration)
	}
	return
}
