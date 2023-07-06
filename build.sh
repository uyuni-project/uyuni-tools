#!/usr/bin/bash

mkdir -p ./bin
podman run --rm -v $PWD:/src \
    --workdir /src/ \
    registry.suse.com/bci/golang:latest \
    go build -o ./bin ./...
