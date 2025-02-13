#!/bin/sh
#
# SPDX-FileCopyrightText: 2025 SUSE LLC
#
# SPDX-License-Identifier: Apache-2.0

current_year=$(date +%Y)
for changed_file in $@; do
    sed -i -E "s/\/\/ SPDX-FileCopyrightText: [0-9]{4}/\/\/ SPDX-FileCopyrightText: $current_year/" $changed_file
done

if test $(git status --porcelain | wc -l) -ne 0 ; then
    echo "✨ Copyright year adjusted, please commit"
fi
