// Copyright 2021 The Vitess Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

variable "auth_token" {
  description = "Equinix Metal auth token"
}

variable "project_id" {
  description = "Project ID"
}

variable "hostname" {
  description = "Hostname given to the new node"
  default = "benchmark-node-terraform"
}

variable "operating_system" {
  description = "Operating system on which to start the node"
  default = "centos_8"
}

variable "instance_type" {
  description = "Equinix Metal instance type that will be used"
  default = "t1.small.x86"
}

variable "facilities" {
  description = "Equinix Metal facility used to run the server"
  default = ["ams1"]
}