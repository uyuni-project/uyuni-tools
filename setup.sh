# SPDX-FileCopyrightText: 2023 SUSE LLC
#
# SPDX-License-Identifier: Apache-2.0

set -euxo pipefail

go mod vendor && tar czvf vendor.tar.gz vendor >/dev/null && rm -rf vendor

echo "vendor.tar.gz"
