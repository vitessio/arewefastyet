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
	"github.com/google/uuid"
	"os"
	"path"
	"path/filepath"
)

const (
	execDir      = "./exec/"
	ansibleDir = "./ansible"
)

func createDirFromUUID(uuid uuid.UUID, root string) (dirPath string, err error) {
	dirPath = path.Join(root, execDir, uuid.String())
	dirPath, err = filepath.Abs(dirPath)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return "", err
	}
	return dirPath, nil
}

func createSubDir(rootDir, subDir string, call func(dir string) error) error {
	subDir = path.Join(rootDir, subDir)
	subdirAbs, err := filepath.Abs(subDir)
	if err != nil {
		return err
	}

	err = os.Mkdir(subdirAbs, 0755)
	if err != nil {
		return err
	}

	err = call(subdirAbs)
	if err != nil {
		return err
	}
	return nil
}

func (e *Exec) copyAllDirs() error {
	err := createSubDir(e.dirPath, ansibleDir, func(dir string) error {
		return e.AnsibleConfig.CopyRootDirectory(dir)
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *Exec) prepareDirectories() error {
	dirPath, err := createDirFromUUID(e.UUID, e.rootDir)
	if err != nil {
		return err
	}
	e.dirPath = dirPath

	err = e.copyAllDirs()
	if err != nil {
		return err
	}
	return nil
}