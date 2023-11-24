set -euxo pipefail

go mod vendor && tar czvf vendor.tar.gz vendor >/dev/null && rm -rf vendor

echo "vendor.tar.gz"
