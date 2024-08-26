/*
 *
 * Copyright 2024 The Vitess Authora.
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

package admin

import (
	"os"
	"path"

	"github.com/vitessio/arewefastyet/go/tools/git"
)

// setupLocalVitess is used to setup the local clone of vitess
func (a *Admin) setupLocalVitess() error {
	files, err := os.ReadDir(a.localVitessPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.Name() == "vitess" && file.IsDir() {
			return nil
		}
	}

	_, err = git.ExecCmd(a.localVitessPath, "git", "clone", "https://github.com/vitessio/vitess.git")

	return err
}

// getVitessPath is used to find the path of the directory where the local vitess clone should exits
func (a *Admin) getVitessPath() string {
	return path.Join(a.localVitessPath, "vitess")
}

// pullLocalVitess is used to execute
func (a *Admin) pullLocalVitess() error {
	_, err := git.ExecCmd(a.getVitessPath(), "git", "fetch", "origin", "--tags")
	if err != nil {
		return err
	}
	_, err = git.ExecCmd(a.getVitessPath(), "git", "reset", "--hard", "origin/main")
	return err
}
