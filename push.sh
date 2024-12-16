#!/usr/bin/bash

# SPDX-FileCopyrightText: 2024 SUSE LLC
#
# SPDX-License-Identifier: Apache-2.0

# This script is called by push-packages-to-obs

OSCAPI=$1
GIT_DIR=$2
PKG_NAME=$3

SRPM_PKG_DIR=$(dirname "$0")

pushd ${GIT_DIR}
REMOTE_BRANCH=$(git for-each-ref --format='%(upstream:lstrip=-1)' $(git rev-parse --symbolic-full-name HEAD))
COMMIT_ID=$(git rev-parse --short HEAD)
popd

if [ "${OSCAPI}" == "https://api.suse.de" ]; then
  VERSION="HEAD"
  case ${REMOTE_BRANCH} in Manager-*)
    VERSION="${REMOTE_BRANCH#Manager-}"
  esac

# Define the default tag to use
  sed 's/^tag=%{!?_default_tag:latest}/tag=5.0.2/' -i ${SRPM_PKG_DIR}/uyuni-tools.spec

  sed "s/namespace='%{_default_namespace}'/namespace='%{_default_namespace}\/%{_arch}'/" -i ${SRPM_PKG_DIR}/uyuni-tools.spec

else

  pushd ${GIT_DIR}
  VERSION=$(git tag --points-at HEAD Uyuni-*)
  popd

  if test -z "${VERSION}"; then
      case ${REMOTE_BRANCH} in Uyuni-*)
        VERSION="${REMOTE_BRANCH#Uyuni-}"
      esac

      if test -z "${VERSION}"; then
        VERSION="Master"
      fi
  fi
fi

# Add the version_details value for use in the version tag
sed "/^%global productname.*$/a%global version_details ${VERSION} $COMMIT_ID" -i ${SRPM_PKG_DIR}/uyuni-tools.spec
