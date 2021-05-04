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

package influxdb

import (
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/vitessio/arewefastyet/go/storage"
)

const (
	ErrorInvalidConfiguration = "invalid configuration"
)

type Client struct {
	influx *influxdb2.Client
	config *Config
}

func (c *Client) Close() error {
	panic("implement me")
}

func (c *Client) Insert(query string, args ...interface{}) (storage.Insertion, error) {
	panic("implement me")
}

func (c *Client) Select(query string, args ...interface{}) (storage.Selection, error) {
	panic("implement me")
}

