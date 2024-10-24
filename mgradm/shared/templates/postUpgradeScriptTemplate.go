// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

const postUpgradeScriptTemplate = `#!/bin/bash
sed 's/cobbler\.host.*/cobbler\.host = localhost/' -i /etc/rhn/rhn.conf;
grep uyuni_authentication_endpoint /etc/cobbler/settings.yaml
if [ $? -eq 1 ]; then
	echo 'uyuni_authentication_endpoint: "http://localhost"' >> /etc/cobbler/settings.yaml
else
	sed 's/uyuni_authentication_endpoint.*/uyuni_authentication_endpoint: http:\/\/localhost/' \
        -i /etc/cobbler/settings.yaml;
fi

grep pam_auth_service /etc/rhn/rhn.conf
if [ $? -eq 1 ]; then
	echo 'pam_auth_service = susemanager' >> /etc/rhn/rhn.conf
else
	sed 's/pam_auth_service.*/pam_auth_service = susemanager/' -i /etc/rhn/rhn.conf;
fi

# (bsc#1231206) fix error happened during migration to 5.0.0. 
if [ -f /var/lib/pgsql/data-pg14/pg_hba.conf ]; then
    echo "Migrating pgsql 14 pg_hba.conf to pgsql 16"
    cp /var/lib/pgsql/data-pg14/pg_hba.conf /var/lib/pgsql/data
    mv /var/lib/pgsql/data-pg14/pg_hba.conf /var/lib/pgsql/data-pg14/pg_hba.conf.migrated
    db_user=$(sed -n '/^db_user/{s/^.*=[ \t]\+\(.*\)$/\1/ ; p}' /etc/rhn/rhn.conf)
    db_name=$(sed -n '/^db_name/{s/^.*=[ \t]\+\(.*\)$/\1/ ; p}' /etc/rhn/rhn.conf)
    ip=$(ip -o -4 addr show up scope global | head -1 | awk '{print $4}' || true)
    echo "host $db_name $db_user $ip scram-sha-256" >> /var/lib/pgsql/data/pg_hba.conf
fi
if [ -f /var/lib/pgsql/data-pg14/postgresql.conf ]; then
    echo "Migrating pgsql 14 postgresql.conf to pgsql 16"
    cp /var/lib/pgsql/data-pg14/postgresql.conf /var/lib/pgsql/data/
    mv /var/lib/pgsql/data-pg14/postgresql.conf /var/lib/pgsql/data-pg14/postgresql.conf.migrated
fi
# end (bsc#1231206)

echo "DONE"`

// PostUpgradeTemplateData represents information used to create post upgrade.
type PostUpgradeTemplateData struct {
}

// Render will create script for finalizing PostgreSQL upgrade.
func (data PostUpgradeTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(postUpgradeScriptTemplate))
	return t.Execute(wr, data)
}
