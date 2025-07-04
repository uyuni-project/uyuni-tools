#!/bin/bash
# SPDX-FileCopyrightText: 2024 SUSE LLC
#
# SPDX-License-Identifier: Apache-2.0
go build $* ./...
go test $* $(go list ./... | grep -v 'tests/e2e')
