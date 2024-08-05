# SPDX-FileCopyrightText: 2024 SUSE LLC
#
# SPDX-License-Identifier: Apache-2.0

FROM mcr.microsoft.com/devcontainers/go:1.20
WORKDIR /build
ENV GOMODCACHE=/opt/golang/
RUN wget https://raw.githubusercontent.com/uyuni-project/uyuni-tools/main/go.mod
RUN wget https://raw.githubusercontent.com/uyuni-project/uyuni-tools/main/go.sum
RUN go mod download
RUN rm go.mod
RUN rm go.sum
