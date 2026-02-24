# Build planetscale/sysbench from source
FROM ubuntu:22.04 AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    git make automake libtool pkg-config \
    libmysqlclient-dev \
    libssl-dev \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

RUN git clone https://github.com/planetscale/sysbench.git && \
    cd sysbench && \
    ./autogen.sh && \
    ./configure --with-mysql && \
    make -j$(nproc) && \
    make install

# Runtime image
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y --no-install-recommends \
    libmysqlclient-dev \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/local/bin/sysbench /usr/local/bin/sysbench
COPY --from=builder /usr/local/share/sysbench /usr/local/share/sysbench

# Clone sysbench-tpcc lua scripts for TPCC workloads
RUN apt-get update && apt-get install -y --no-install-recommends \
    git ca-certificates \
    && git clone https://github.com/planetscale/sysbench-tpcc.git /src/sysbench-tpcc \
    && rm -rf /var/lib/apt/lists/*

COPY sysbench/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY sysbench/olap_sort.lua /usr/local/share/sysbench/

ENTRYPOINT ["/entrypoint.sh"]
