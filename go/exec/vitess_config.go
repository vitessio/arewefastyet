/*
 *
 * Copyright 2022 The Vitess Authors.
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
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/vitessio/arewefastyet/go/infra/ansible"
	"github.com/vitessio/arewefastyet/go/tools/git"
)

type (
	vitessConfig struct {
		vtgate, vttablet string
	}

	rawVitessVersionConfig struct {
		version git.Version
		value   rawSingleVitessVersionConfig
	}

	rawSingleVitessVersionConfig = map[string]interface{}
)

func prepareVitessConfiguration(rawVitessConfig rawSingleVitessVersionConfig, vitessVersion git.Version, vitessConfig *vitessConfig) error {
	if len(rawVitessConfig) == 0 {
		return nil
	}

	var data []rawVitessVersionConfig

	for rawVersion, value := range rawVitessConfig {
		parsedVersion := strings.Split(rawVersion, "-")
		if len(parsedVersion) == 0 {
			continue
		}
		var err error
		v := git.Version{}
		if len(parsedVersion) > 0 {
			v.Major, err = strconv.Atoi(parsedVersion[0])
			if err != nil {
				return err
			}
		}
		if len(parsedVersion) > 1 {
			v.Minor, err = strconv.Atoi(parsedVersion[1])
			if err != nil {
				return err
			}
		}
		if len(parsedVersion) > 2 {
			v.Patch, err = strconv.Atoi(parsedVersion[2])
			if err != nil {
				return err
			}
		}
		data = append(data, rawVitessVersionConfig{
			version: v,
			value:   value.(map[string]interface{}),
		})
	}

	sort.Slice(data, func(i, j int) bool {
		return git.CompareVersionNumbers(data[i].version, data[j].version) == 1
	})

	lastElem := rawVitessVersionConfig{}
	for i := 0; i < len(data); i++ {
		cmp := git.CompareVersionNumbers(vitessVersion, data[i].version)
		if cmp == 1 || cmp == 0 {
			lastElem = data[i]
			break
		}
	}

	err := getVitessConfigFromMap(lastElem.value, "vtgate", &vitessConfig.vtgate)
	if err != nil {
		return err
	}

	err = getVitessConfigFromMap(lastElem.value, "vttablet", &vitessConfig.vttablet)
	if err != nil {
		return err
	}

	return nil
}

func getVitessConfigFromMap(value map[string]interface{}, key string, set *string) error {
	if elem, ok := value[key]; ok {
		elem, str := elem.(string)
		if !str {
			return fmt.Errorf("could not parse the %s value for this vitess configuration", key)
		}
		*set = elem
	}
	return nil
}

func (vcfg vitessConfig) addToAnsible(ansibleCfg *ansible.Config) {
	if len(vcfg.vtgate) != 0 {
		ansibleCfg.AddExtraVar(ansible.KeyExtraFlagsVTGate, vcfg.vtgate)
	}
	if len(vcfg.vttablet) != 0 {
		ansibleCfg.AddExtraVar(ansible.KeyExtraFlagsVTTablet, vcfg.vttablet)
	}
}
