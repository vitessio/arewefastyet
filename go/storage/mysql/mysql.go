/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	ErrorClientConnectionNotInitialized = "the client connection to the database is not initialized"
)

type Client struct {
	db *sql.DB
}

// New creates a new Client based on the given ConfigDB.
func New(config ConfigDB) (client *Client, err error) {
	client = &Client{}
	client.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", config.User, config.Password, config.Host, config.Database))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) Close() error {
	if c.db == nil {
		return errors.New(ErrorClientConnectionNotInitialized)
	}
	return c.db.Close()
}

func (c *Client) Write(query string, args ...interface{}) (int64, error) {
	if c.db == nil {
		return 0, errors.New(ErrorClientConnectionNotInitialized)
	}
	stms, err := c.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stms.Close()

	res, err := stms.Exec(args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (c *Client) Read(query string, args ...interface{}) (*sql.Rows, error) {
	if c.db == nil {
		return nil, errors.New(ErrorClientConnectionNotInitialized)
	}
	rows, err := c.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
