// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

// Procedure from the https://www.postgresql.org/docs/16/continuous-archiving.html#BACKUP-PITR-RECOVERY
const postgresRestoreScriptTemplate = `
/usr/bin/bash -e

if [[ ! -f {{ .Basebackup }} ]]; then
    echo "Missing base backup file!"
	exit 1
fi

echo "Removing old cluster files..."
find "{{ .Datadir }}" -mindepth 1 -delete

echo "Restoring basebackup..."
tar -xf "{{ .Basebackup }}" -C "{{ .Datadir }}"
chown -R postgres:postgres "{{ .Datadir }}"

echo "Signal postgresql to start in recovery mode"
touch "{{ .Datadir }}/recovery.signal"
`

// FinalizePostgresTemplateData represents information used to create PostgreSQL migration script.
type RestorePostgresTemplateData struct {
	Datadir    string
	Basebackup string
}

// Render will create script for finalizing PostgreSQL upgrade.
func (data RestorePostgresTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(postgresRestoreScriptTemplate))
	return t.Execute(wr, data)
}
