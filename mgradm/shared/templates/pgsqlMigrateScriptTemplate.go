// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

//nolint:lll
const pgsqlMigrationScriptTemplate = `#!/bin/bash
set -e -x

if [ -d /var/lib/pgsql/data/data ] ; then
    shopt -s dotglob
    mv /var/lib/pgsql/data/data/* /var/lib/pgsql/data
    rmdir /var/lib/pgsql/data/data

    echo "Adding database access for other containers..."
    db_user=$(sed -n '/^db_user/{s/^.*=[ \t]\+\(.*\)$/\1/ ; p}' /etc/rhn/rhn.conf)
    db_name=$(sed -n '/^db_name/{s/^.*=[ \t]\+\(.*\)$/\1/ ; p}' /etc/rhn/rhn.conf)
    ip=$(ip -o -4 addr show up scope global | head -1 | awk '{print $4}' || true)
    echo "host $db_name $db_user all scram-sha-256" >> /var/lib/pgsql/data/pg_hba.conf

    ls -la /var/lib/pgsql/data
    
fi

{{ if .ReportDBHost }}
sed 's/^report_db_host = .*/report_db_host = {{ .ReportDBHost }}/' -i /etc/rhn/rhn.conf;
{{ end }}

{{ if .DBHost }}
sed 's/^db_host = .*/db_host = {{ .DBHost }}/' -i /etc/rhn/rhn.conf;
{{ end }}


echo "DONE"`

// MigrateScriptTemplateData represents migration information used to create migration script.
type PgsqlMigrateScriptTemplateData struct {
	DBHost       string
	ReportDBHost string
}

// Render will create migration script.
func (data PgsqlMigrateScriptTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(pgsqlMigrationScriptTemplate))
	return t.Execute(wr, data)
}
