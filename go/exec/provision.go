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

package exec

import (
	"encoding/json"
	"github.com/vitessio/arewefastyet/go/infra"
)

func provision(infra infra.Infra) (IPs []string, err error) {
	if err = infra.Prepare(); err != nil {
		return nil, err
	}

	out, err := infra.Create("device_public_ip")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(out["device_public_ip"]), &IPs)
	if err != nil {
		return nil, err
	}
	return IPs, nil
}
