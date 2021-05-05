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

package metrics

import (
	"fmt"
	"github.com/vitessio/arewefastyet/go/storage/influxdb"
	"log"
)

func GetCPU(client influxdb.Client, start, end, execUUID, component string) error {
	resultClient, err := client.Select(fmt.Sprintf(`from(bucket:"%s")
			|> range(start: %s, stop: %s)
			|> filter(fn:(r) => r._measurement == "process_cpu_seconds_total" and r.exec_uuid == "%s" and r.component == "%s")`,
		client.Config.Database, start, end, execUUID, component))
	if err != nil {
		return err
	}
	result := resultClient.(influxdb.SelectResult)
	log.Println(len(result.Values))
	return nil
}
