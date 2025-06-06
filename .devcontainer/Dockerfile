# SPDX-FileCopyrightText: 2024 SUSE LLC
#
# SPDX-License-Identifier: Apache-2.0

FROM opensuse/tumbleweed:latest

RUN set -euo pipefail; \
    curl -o uyuni-tools.gpg https://build.opensuse.org/projects/systemsmanagement:Uyuni:Utils/signing_keys/download?kind=gpg && \
    rpm --import uyuni-tools.gpg && \
    zypper ar https://download.opensuse.org/repositories/systemsmanagement:/Uyuni:/Utils/openSUSE_Tumbleweed/ uyuni-utils && \
    zypper -n in go1.21 go1.21-doc golangci-lint git-core curl neovim python313-pre-commit python313 gpg2 \
        sed grep gettext-tools util-linux less awk uyuni-releng-tools

RUN set -euo pipefail; zypper -n clean -a; \
    rm -rf {/target,}/var/log/{alternatives.log,lastlog,tallylog,zypper.log,zypp/history,YaST2}
