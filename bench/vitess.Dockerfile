ARG VITESS_REF=main

# Stage 1: Build Vitess
FROM golang:1.26 AS builder

ARG VITESS_REF

RUN apt-get update && apt-get install -y --no-install-recommends \
    git make curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

RUN git clone --filter=blob:none https://github.com/vitessio/vitess.git && \
    cd vitess && \
    git checkout ${VITESS_REF} && \
    MAKEFLAGS="-j2" GOFLAGS="-p=2" make build

# Download mysqld_exporter (version 0.12.1, matching production)
RUN curl -sL "https://github.com/prometheus/mysqld_exporter/releases/download/v0.12.1/mysqld_exporter-0.12.1.linux-amd64.tar.gz" \
    | tar -xz -C /tmp && \
    cp /tmp/mysqld_exporter-0.12.1.linux-amd64/mysqld_exporter /usr/local/bin/mysqld_exporter

# Stage 2: Runtime image
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y --no-install-recommends \
    mysql-server-8.0 \
    curl \
    jq \
    etcd-client \
    gosu \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/vitess/bin/vtgate       /usr/local/bin/vtgate
COPY --from=builder /build/vitess/bin/vttablet     /usr/local/bin/vttablet
COPY --from=builder /build/vitess/bin/vtctld       /usr/local/bin/vtctld
COPY --from=builder /build/vitess/bin/vtctl        /usr/local/bin/vtctl
COPY --from=builder /build/vitess/bin/vtctlclient   /usr/local/bin/vtctlclient
COPY --from=builder /build/vitess/bin/vtctldclient /usr/local/bin/vtctldclient
COPY --from=builder /build/vitess/bin/mysqlctld    /usr/local/bin/mysqlctld
COPY --from=builder /build/vitess/bin/mysqlctl     /usr/local/bin/mysqlctl
COPY --from=builder /build/vitess/config           /vt/config/
COPY --from=builder /usr/local/bin/mysqld_exporter /usr/local/bin/mysqld_exporter

ENV VTDATAROOT=/vt

RUN groupadd -r vitess && useradd -r -g vitess -s /bin/bash vitess && \
    mkdir -p /vt && chown -R vitess:vitess /vt
