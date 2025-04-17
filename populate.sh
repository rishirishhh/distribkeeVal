#!/bin/bash

for shard in localhost:8080 localhost:8081; do
    echo $shard
    for i in {1..100}; do
        curl "http://$shard/set?key=key-$RANDOM&value=value-$RANDOM"
    done
done