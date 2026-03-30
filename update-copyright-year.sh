#!/bin/bash
#
# SPDX-FileCopyrightText: 2026 SUSE LLC
#
# SPDX-License-Identifier: Apache-2.0

need_fix="n"
changed="n"
current_year=$(date +%Y)
for changed_file in $@; do
    case "$changed_file" in
        uyuni-tools.changes.*)
            continue
            ;;
    esac


    if ! grep -q "^[/# ]*SPDX-FileCopyrightText: $current_year" $changed_file; then
        lines=`grep "^[/# ]*SPDX-FileCopyrightText: " $changed_file | wc -l`
        if [ "z$lines" = "z1" ]; then
            echo "🛠️ Fixed copyright year on $changed_file"
            sed -i -E "s/^([/# ]*)SPDX-FileCopyrightText: [0-9]{4}/\1SPDX-FileCopyrightText: $current_year/" $changed_file
            changed="y"
        else
            echo "🛑 Cannot update the copyright year in $changed_file (multiple lines)"
            need_fix="y"
        fi
    fi
done

if [ "$need_fix" = "y" ]; then
    exit 1
fi

if [ "$changed" = "y" ]; then
    echo "✨ Copyright year adjusted, please commit"
    exit 1
fi
