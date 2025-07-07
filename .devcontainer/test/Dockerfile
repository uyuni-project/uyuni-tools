FROM registry.opensuse.org/uyuni/init:latest

# Install necessary packages, clean up afterwards
RUN set -euo pipefail; zypper -n in --no-recommends curl git tar gzip make gcc podman; zypper -n clean; rm -rf /var/log/*

# Install Go
RUN curl -LO https://go.dev/dl/go1.21.13.linux-${TARGETARCH}.tar.gz
RUN tar -C /usr/local -xzf go1.21.13.linux-${TARGETARCH}.tar.gz
ENV PATH="/usr/local/go/bin:$PATH"

# Install Ginkgo
RUN go install github.com/onsi/ginkgo/v2/ginkgo@v2.20.1

# Location of the source code
WORKDIR /workspace

CMD ["/usr/lib/systemd/systemd"]
