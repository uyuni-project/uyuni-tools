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
	rsync -a --exclude=/var/lib/pgsql/data/data/pg_hba.conf /var/lib/pgsql/data/data/ /var/lib/pgsql/data/ 2>/dev/null
    rm -rf /var/lib/pgsql/data/data
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
