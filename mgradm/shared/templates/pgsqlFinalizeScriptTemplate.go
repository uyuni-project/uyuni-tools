// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

//nolint:lll
const postgresFinalizeScriptTemplate = `#!/bin/bash
set -e

{{ if .RunReindex }}
echo "Reindexing database. This may take a while, please do not cancel it!"
database=$(sed -n "s/^\s*db_name\s*=\s*\([^ ]*\)\s*$/\1/p" /etc/rhn/rhn.conf)
spacewalk-sql --select-mode - <<<"REINDEX DATABASE \"${database}\";"
{{ end }}

{{ if .RunSchemaUpdate }}
echo "Schema update..."
/usr/sbin/spacewalk-startup-helper check-database
{{ end }}

{{ if .Migration }}
echo "Updating auto-installable distributions..."
spacewalk-sql --select-mode - <<EOT
SELECT MIN(CONCAT(org_id, '-', label)) AS target, base_path INTO TEMP TABLE dist_map FROM rhnKickstartableTree GROUP BY base_path;
UPDATE rhnKickstartableTree SET base_path = CONCAT('/srv/www/distributions/', target)
    from dist_map WHERE dist_map.base_path = rhnKickstartableTree.base_path;
DROP TABLE dist_map;
EOT

echo "Schedule a system list update task..."
spacewalk-sql --select-mode - <<EOT
insert into rhnTaskQueue (id, org_id, task_name, task_data)
SELECT nextval('rhn_task_queue_id_seq'), 1, 'update_system_overview', s.id
from rhnserver s
where not exists (select 1 from rhntaskorun r join rhntaskotemplate t on r.template_id = t.id
join rhntaskobunch b on t.bunch_id = b.id where b.name='update-system-overview-bunch' limit 1);
EOT
{{ end }}
`

// FinalizePostgresTemplateData represents information used to create PostgreSQL migration script.
type FinalizePostgresTemplateData struct {
	RunReindex      bool
	RunSchemaUpdate bool
	Migration       bool
	Kubernetes      bool
}

// Render will create script for finalizing PostgreSQL upgrade.
func (data FinalizePostgresTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(postgresFinalizeScriptTemplate))
	return t.Execute(wr, data)
}
