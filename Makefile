# Copyright 2021 The Vitess Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.ONESHELL:
SHELL = /bin/bash

PY_VERSION = 3.7
VIRTUALENV_PATH = benchmark

BIN_NAME = arewefastyetcli

.PHONY: install virtual_env molecule_converge_all

all: build

virtual_env:
	test -d $(VIRTUALENV_PATH) || virtualenv -p python$(PY_VERSION) $(VIRTUALENV_PATH)

install: requirements.txt virtual_env $(VIRTUALENV_PATH)
	source $(VIRTUALENV_PATH)/bin/activate && \
	pip install -r ./requirements.txt && \
	ansible-galaxy install cloudalchemy.prometheus && \
	ansible-galaxy install cloudalchemy.node_exporter

build:
	go build -o $(BIN_NAME) ./go/main.go

clean:
	rm -f $(BIN_NAME)

gendoc: build exec_gendoc clean

exec_gendoc:
	./$(BIN_NAME) gen doc

test:
	go test -v ./go/...

molecule_converge_all: virtual_env $(VIRTUALENV_PATH)
	source $(VIRTUALENV_PATH)/bin/activate && \
	cd $(ANSIBLE_PATH)/roles/vitess_build && \
	molecule create --scenario-name all && \
	OBJC_DISABLE_INITIALIZE_FORK_SAFETY=YES molecule converge --scenario-name all

molecule_dryrun_all: virtual_env $(VIRTUALENV_PATH)
	source $(VIRTUALENV_PATH)/bin/activate && \
	cd $(ANSIBLE_PATH)/roles/vitess_build && \
	OBJC_DISABLE_INITIALIZE_FORK_SAFETY=YES molecule converge --scenario-name all -- --skip-tags="oltp,tpcc"

tools:
	echo Installing git hooks
	chmod +x bootstrap.sh
	./bootstrap.sh
