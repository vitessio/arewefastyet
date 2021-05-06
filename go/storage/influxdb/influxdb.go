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
)

const (
	ErrorInvalidConfiguration = "invalid configuration"
)

type Client struct {
	influx influxdb2.Client
	Config *Config
}

func (c *Client) Select(query string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	queryAPI := c.influx.QueryAPI("")
	queryResult, err := queryAPI.Query(context.Background(), query)
	if err == nil {
		for queryResult.Next() {
			result = append(result, queryResult.Record().Values())
		}
		if queryResult.Err() != nil {
			return result, fmt.Errorf("Query error: %s\n", queryResult.Err().Error())
		}
	} else {
		return result, fmt.Errorf("Query error: %s\n", err.Error())
	}
	return result, nil
}
