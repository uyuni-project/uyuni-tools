#!/usr/bin/bash

mkdir -p ./bin

tag=$(git describe --tags --abbrev=0 | cut -f 3 -d '-')
offset=$(git rev-list --count ${tag})

VERSION_NAME=github.com/uyuni-project/uyuni-tools/shared/utils.Version

go build -tags netgo -ldflags "-X ${VERSION_NAME}=${tag}-${offset}" -o ./bin ./...
