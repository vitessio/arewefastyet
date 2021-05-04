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

// Config defines the required configuration used to authenticate
// to an InfluxDB database.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// IsValid return true if Config is ready to be used, and false otherwise.
func (cfg Config) IsValid() bool {
	return cfg.Host != ""
}
