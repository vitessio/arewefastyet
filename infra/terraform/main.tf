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

provider "packet" {
  auth_token = var.auth_token
}

resource "packet_device" "node" {
  hostname         = var.hostname
  plan             = var.instance_type
  facilities       = var.facilities
  operating_system = var.operating_system
  billing_cycle    = "hourly"
  project_id       = var.project_id
}