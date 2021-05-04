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
	"errors"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go"
)

const (
	ErrorInvalidConfiguration = "invalid configuration"
)

type Client struct {
	influx *influxdb2.Client
	config *Config
}

// New creates a new InfluxDB Client using the given Config.
// If no port is defined, the port will be set to its default ("8086").
func New(config *Config) (client *Client, err error) {
	if !config.IsValid() {
		return nil, errors.New(ErrorInvalidConfiguration)
	}
	if config.Port == "" {
		config.Port = "8086"
	}
	client = &Client{
		config: config,
	}
	influxclient := influxdb2.NewClient(config.Host+":"+config.Port, fmt.Sprintf("%s:%s", config.User, config.Password))
	client.influx = &influxclient
	return client, nil
}
