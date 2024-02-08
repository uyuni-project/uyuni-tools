#!/usr/bin/bash

# SPDX-FileCopyrightText: 2024 SUSE LLC
#
# SPDX-License-Identifier: Apache-2.0

# This script is called by push-packages-to-obs

OSCAPI=$1
GIT_DIR=$2
PKG_NAME=$3

SRPM_PKG_DIR=$(dirname "$0")

if [ "${OSCAPI}" == "https://api.suse.de" ]; then
  sed 's/^tag=%{!?_default_tag:latest}/tag=5.0.0-beta1/' -i ${SRPM_PKG_DIR}/uyuni-tools.spec
fi
