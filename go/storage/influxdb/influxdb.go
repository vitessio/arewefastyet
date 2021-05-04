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
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/vitessio/arewefastyet/go/storage"
)

const (
	ErrorInvalidConfiguration = "invalid configuration"
)

type Client struct {
	influx influxdb2.Client
	config *Config
}

func (c *Client) Close() error {
	panic("implement me")
}

func (c *Client) Insert(query string, args ...interface{}) (storage.Insertion, error) {
	panic("implement me")
}

func (c *Client) Select(query string, args ...interface{}) (storage.Selection, error) {
	queryAPI := c.influx.QueryAPI("")
	// Supply string in a form database/retention-policy as a bucket. Skip retention policy for the default one, use just a database name (without the slash character)
	result, err := queryAPI.Query(context.Background(), `from(bucket:"test")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "stat")`)
	if err == nil {
		for result.Next() {
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			fmt.Printf("row: %s\n", result.Record().String())
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	} else {
		fmt.Printf("Query error: %s\n", err.Error())
	}
	return nil, nil
}

