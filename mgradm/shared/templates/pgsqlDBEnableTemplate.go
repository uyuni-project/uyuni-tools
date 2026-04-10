// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

const postgresEnableScriptTemplate = `
/usr/bin/bash -e

if [[ ! -d {{ .BackupDir }} ]]; then
    echo "Missing backup directory!"
	exit 1
fi

echo "Removing old backup files..."
find "{{ .BackupDir }}" -mindepth 1 -delete

echo "Adjusting ownership..."
chown postgres:postgres {{ .BackupDir }}
chmod u=rwx,og= {{ .BackupDir }}

# Inspired by the original smdba
# https://github.com/SUSE/smdba/blob/master/src/smdba/postgresqlgate.py#L853C110-L853C120
echo "Performing initial base backup..."
su postgres -c '/usr/bin/pg_basebackup -U postgres -D {{ .BackupDir }} -Ft -z -c fast -X fetch -v'
`

// FinalizePostgresTemplateData represents information used to create PostgreSQL migration script.
type EnablePostgresTemplateData struct {
	BackupDir string
}

// Render will create script for finalizing PostgreSQL upgrade.
func (data EnablePostgresTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(postgresEnableScriptTemplate))
	return t.Execute(wr, data)
}
