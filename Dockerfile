# Copyright 2023 The Vitess Authors.
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

FROM golang:1.22.5-bookworm AS build-stage

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build arewefastyet
RUN CGO_ENABLED=0 GOOS=linux go build -o /arewefastyetcli ./go/main.go

FROM debian:bookworm AS run-stage

# Install Git, Golang, and Python
RUN apt-get update && apt-get install -y \
    git \
    python3 \
    python3-pip \
    python3-venv \
    wget \
    gnutls-bin

# Set up Python virtual environment
RUN python3 -m venv /venv
ENV PATH="/venv/bin:$PATH"

# Upgrade pip and install requirements
RUN pip3 install --upgrade pip
COPY requirements.txt .
RUN pip3 install -r requirements.txt

# Install ansible add-ons
RUN ansible-galaxy install cloudalchemy.node_exporter && ansible-galaxy install cloudalchemy.prometheus

# Copy the source code to the working directory
COPY --from=build-stage /arewefastyetcli /arewefastyetcli

EXPOSE 8080

# Needed for Ansible to execute sub-processes
ENV OBJC_DISABLE_INITIALIZE_FORK_SAFETY=YES

# Make sure all directories are created
RUN mkdir -p /config /exec

# Configuration files MUST be attached to the container using a volume.
# The configuration files are not mounted on the Docker image for obvious
# security reasons.
CMD ["/arewefastyetcli", "api", "--config", "/config/config.yaml", "--secrets", "/config/secrets.yaml"]
