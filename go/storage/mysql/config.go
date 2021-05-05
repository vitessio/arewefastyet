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

package mysql

import (
	"database/sql"
	"fmt"
)

type ConfigDB struct {
	Host     string
	User     string
	Password string
	Database string
}

func (cfg ConfigDB) NewClient() (*Client, error) {
	var err error
	client := &Client{}
	client.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Database))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (cfg ConfigDB) IsValid() bool {
	return !(cfg.Database == "" || cfg.User == "" || cfg.Host == "")
}

