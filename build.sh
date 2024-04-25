#!/usr/bin/bash

# SPDX-FileCopyrightText: 2024 SUSE LLC
#
# SPDX-License-Identifier: Apache-2.0
set -e
mkdir -p ./bin

tag=$(git describe --tags --abbrev=0)
version=$(git describe --tags --abbrev=0 | cut -f 3 -d '-')
offset=$(git rev-list --count ${tag}..)
commit_id=$(git rev-parse --short HEAD)

VERSION_NAME=github.com/uyuni-project/uyuni-tools/shared/utils.Version

CGO_ENABLED=0 go build -ldflags "-X \"${VERSION_NAME}=${version}-${offset} (${commit_id})\"" -o ./bin ./...

for shell in "bash" "zsh" "fish"; do
    COMPLETION_FILE="./bin/completion.${shell}"

    # generate and source shell completion scripts for mgradm and mgrctl
    ./bin/mgradm completion ${shell} > "${COMPLETION_FILE}"
    ./bin/mgrctl completion ${shell} >> "${COMPLETION_FILE}"
    ./bin/mgrpxy completion ${shell} >> "${COMPLETION_FILE}"
done

golangci-lint run
./check_localizable
echo "DONE"
