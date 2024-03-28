#!/bin/sh
#
# SPDX-FileCopyrightText: 2024 SUSE LLC
#
# SPDX-License-Identifier: Apache-2.0

PREFIX=$1
locales_dir=$(dirname $0)/
if test "x${PREFIX}" == "x"; then
    PREFIX=${locales_dir}
fi

for domain in mgrctl mgradm mgrpxy; do
    for po_file in `ls ${locales_dir}/${domain}/*.po`; do
        lang=$(basename ${po_file} | sed 's/\.po$//')
        locale_dir=${PREFIX}${lang}/LC_MESSAGES
        install -vd -m 0755 ${locale_dir}

        msgcat -o ${locale_dir}/${domain}.po ${po_file} ${locales_dir}/shared/${lang}.po
        msgfmt -c -o ${locale_dir}/${domain}.mo ${locale_dir}/${domain}.po
        if test $? -ne 0;
        then
            echo "Broken ${po_file}"
            exit 1
        fi
        rm ${locale_dir}/${domain}.po
    done
done
