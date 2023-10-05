#!/usr/bin/bash

mkdir -p ./bin
go build -tags netgo -o ./bin ./...
