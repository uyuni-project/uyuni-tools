#!/usr/bin/bash

# SPDX-FileCopyrightText: 2023 SUSE LLC
#
# SPDX-License-Identifier: Apache-2.0

mkdir -p ./bin

tag=$(git describe --tags --abbrev=0 | cut -f 3 -d '-')
offset=$(git rev-list --count ${tag})

VERSION_NAME=github.com/uyuni-project/uyuni-tools/shared/utils.Version

go build -tags netgo -ldflags "-X ${VERSION_NAME}=${tag}-${offset}" -o ./bin ./...
go build -tags netgo -o ./bin ./...

for shell in "bash" "zsh" "fish"; do
    COMPLETION_FILE="./bin/completion.${shell}"

    # generate and source shell completion scripts for mgradm and mgrctl
    ./bin/mgradm completion ${shell} > "${COMPLETION_FILE}"
    ./bin/mgrctl completion ${shell} >> "${COMPLETION_FILE}"
done
