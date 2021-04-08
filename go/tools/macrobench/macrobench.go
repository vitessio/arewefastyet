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

package macrobench

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const (
	prefixMacrobenchSysbenchConfig = "macrobench_"
)

func buildSysbenchArgString(m map[string]string, step string) []string {
	output := map[string]string{}
	for k, v := range m {
		idx := strings.Index(k, "_")
		if idx < 0 {
			continue
		}
		head := k[:idx]
		tail := k[idx+1:]
		if head == step || head == "all" {
			if _, exists := output[tail]; exists {
				if head == "all" {
					continue
				}
			}
			output[tail] = v
		}
	}

	var results []string
	for k, v := range output {
		results = append(results, fmt.Sprintf("--%s=%s", k, v))
	}
	return results
}

func MacroBench(mabcfg MacroBenchConfig) error {
	var results []MacroBenchmarkResult
	var resStr []byte

	mabcfg.parseIntoMap(prefixMacrobenchSysbenchConfig)

	newSteps := skipSteps(steps, mabcfg.SkipSteps)
	for _, step := range newSteps {
		log.Println("Step", step.name)
		args := buildSysbenchArgString(mabcfg.M, step.name)
		args = append(args, mabcfg.WorkloadPath, step.sysbenchName)
		log.Println(strings.Join(args, " "))
		command := exec.Command(mabcfg.SysbenchExec, args...)
		out, err := command.Output()
		if err != nil {
			return err
		}
		if step.name == stepRun {
			resStr = out
		}
	}

	err := json.Unmarshal(resStr, &results)
	if err != nil {
		return fmt.Errorf("unmarshal results: %+v\n", err)
	}

	log.Printf("%+v\n", results[0])
	return nil
}
