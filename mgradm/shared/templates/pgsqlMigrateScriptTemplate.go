// SPDX-FileCopyrightText: 2025 SUSE LLC
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
	rsync -a --exclude=pg_hba.conf /var/lib/pgsql/data/data/ /var/lib/pgsql/data/ 2>/dev/null
    rm -rf /var/lib/pgsql/data/data

    echo "Adding database access for other containers..."
    db_user=$(sed -n '/^db_user/{s/^.*=[ \t]\+\(.*\)$/\1/ ; p}' /etc/rhn/rhn.conf)
    db_name=$(sed -n '/^db_name/{s/^.*=[ \t]\+\(.*\)$/\1/ ; p}' /etc/rhn/rhn.conf)
    echo "host $db_name $db_user all scram-sha-256" >> /var/lib/pgsql/data/pg_hba.conf
    report_db_user=$(sed -n '/^report_db_user/{s/^.*=[ \t]\+\(.*\)$/\1/ ; p}' /etc/rhn/rhn.conf)
    report_db_name=$(sed -n '/^report_db_name/{s/^.*=[ \t]\+\(.*\)$/\1/ ; p}' /etc/rhn/rhn.conf)
    echo "host $report_db_name $report_db_user all scram-sha-256" >> /var/lib/pgsql/data/pg_hba.conf

    echo "host postgres postgres 127.0.0.1/32 scram-sha-256" >> /var/lib/pgsql/data/pg_hba.conf
    echo "host postgres postgres ::1/128 scram-sha-256" >> /var/lib/pgsql/data/pg_hba.conf
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
