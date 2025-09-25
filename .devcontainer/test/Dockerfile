FROM registry.opensuse.org/uyuni/init:latest
ARG TARGETARCH

# Install necessary packages, clean up afterwards
RUN set -euo pipefail; zypper -n in --no-recommends hostname curl git vi tar make gcc podman; zypper -n clean; rm -rf /var/log/*

# Ensure Podman uses netavark
RUN mkdir -p /etc/containers && \
    echo '[network]' > /etc/containers/containers.conf && \
    echo 'network_backend = "netavark"' >> /etc/containers/containers.conf

# Install Go
RUN curl -LO https://go.dev/dl/go1.21.13.linux-${TARGETARCH}.tar.gz
RUN tar -C /usr/local -xzf go1.21.13.linux-${TARGETARCH}.tar.gz
ENV PATH="/usr/local/go/bin:$PATH"

# Install Ginkgo
RUN go install github.com/onsi/ginkgo/v2/ginkgo@v2.20.1

# Location of the source code
WORKDIR /workspace

CMD ["/usr/lib/systemd/systemd"]
