#!/bin/bash
set -e

trap 'killall distribkv dlv; rm -f *.db' SIGINT

cd $(dirname $0)

killall distribkv dlv || true
sleep 0.1

# Build the binary explicitly
go build -gcflags="all=-N -l" -o distribkv

# Run one instance with delve
dlv exec ./distribkv --headless --listen=:2345 --api-version=2 --log -- \
    -db-location=moscow.db -http-addr=127.0.0.1:8080 -config-file=sharding.toml -shard=Moscow &

# Run other instances normally
./distribkv -db-location=minsk.db -http-addr=127.0.0.1:8081 -config-file=sharding.toml -shard=Minsk &
./distribkv -db-location=kiev.db -http-addr=127.0.0.1:8082 -config-file=sharding.toml -shard=Kiev &
./distribkv -db-location=tashkent.db -http-addr=127.0.0.1:8083 -config-file=sharding.toml -shard=Tashkent &

wait
