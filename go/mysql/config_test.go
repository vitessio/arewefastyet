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

import "testing"

func TestConfigDB_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		config    ConfigDB
		wantValid bool
	}{
		{name: "Valid ConfigDB", config: ConfigDB{Host: "host", User: "user", Password: "password", Database: "database"}, wantValid: true},
		{name: "Invalid ConfigDB, missing host", config: ConfigDB{User: "user", Password: "password", Database: "database"}, wantValid: false},
		{name: "Invalid ConfigDB, missing user", config: ConfigDB{Host: "host", Password: "password", Database: "database"}, wantValid: false},
		{name: "Invalid ConfigDB, missing password", config: ConfigDB{Host: "host", User: "user", Database: "database"}, wantValid: false},
		{name: "Invalid ConfigDB, missing database", config: ConfigDB{Host: "host", User: "user", Password: "password"}, wantValid: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotValid := tt.config.IsValid(); gotValid != tt.wantValid {
				t.Errorf("IsValid() = %v, want %v", gotValid, tt.wantValid)
			}
		})
	}
}
