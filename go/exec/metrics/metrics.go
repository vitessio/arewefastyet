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
)

func GetCPU(client influxdb.Client, execUUID string) {
	client.Select(fmt.Sprintf(`from(bucket:"%s") |> range(start:-48h) |> filter(fn:(r) => r._measurement == "process_cpu_seconds_total" and r.exec_uuid == "%s")`, client.Config.Database, execUUID))
}
