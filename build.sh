#!/usr/bin/sh

# SPDX-FileCopyrightText: 2024 SUSE LLC
#
# SPDX-License-Identifier: Apache-2.0
set -e
mkdir -p ./bin

build_tags=$1
if [ -n "${build_tags}" ]; then
    build_tags="-tags ${build_tags}"
fi

tag=$(git describe --tags --abbrev=0)
version=$(git describe --tags --abbrev=0 | cut -f 3 -d '-')
offset=$(git rev-list --count ${tag}..)

VERSION_NAME=github.com/uyuni-project/uyuni-tools/shared/utils.Version

CGO_ENABLED=0 go build ${build_tags} -ldflags "-X ${VERSION_NAME}=${tag}-${offset} ${GO_LDFLAGS}" -o ./bin ./...

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
