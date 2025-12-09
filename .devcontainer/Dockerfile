# SPDX-FileCopyrightText: 2024 SUSE LLC
#
# SPDX-License-Identifier: Apache-2.0

FROM registry.opensuse.org/opensuse/tumbleweed:latest

RUN set -euo pipefail; \
    curl -o uyuni-tools.gpg https://build.opensuse.org/projects/systemsmanagement:Uyuni:Utils/signing_keys/download?kind=gpg && \
    rpm --import uyuni-tools.gpg && \
    zypper ar https://download.opensuse.org/repositories/systemsmanagement:/Uyuni:/Utils/openSUSE_Tumbleweed/ uyuni-utils && \
    zypper -n in \
    awk \
    bash-completion \
    curl \
    gettext-tools \
    git-core \
    go1.25 \
    go1.25-doc \
    gpg2 \
    grep \
    less \
    neovim \
    openssh \
    python313 \
    python313-pre-commit \
    sed \
    uyuni-releng-tools \
    util-linux

RUN curl -sSfL https://github.com/golangci/golangci-lint/releases/download/v2.7.1/golangci-lint-2.7.1-linux-amd64.tar.gz | tar -xz &&  \
    mv golangci-lint-2.7.1-linux-amd64/golangci-lint /usr/local/bin/ && \
    rm -rf golangci-lint-2.7.1-linux-amd64

RUN set -euo pipefail; zypper -n clean -a; \
    rm -rf {/target,}/var/log/{alternatives.log,lastlog,tallylog,zypper.log,zypp/history,YaST2}
