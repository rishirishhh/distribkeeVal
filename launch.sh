#!/bin/bash
set -e

trap 'killall -SIGINT go' SIGINT

cd "$(dirname "$0")"

killall -q go || true
sleep 0.1

# Run multiple instances using `go run`
go run main.go -db-location="$PWD/moscow.db" -http-addr=127.0.0.1:8080 -config-file="$PWD/sharding.toml" -shard=Moscow &
go run main.go -db-location="$PWD/minsk.db" -http-addr=127.0.0.1:8081 -config-file="$PWD/sharding.toml" -shard=Minsk &
go run main.go -db-location="$PWD/kiev.db" -http-addr=127.0.0.1:8082 -config-file="$PWD/sharding.toml" -shard=Kiev &

wait
