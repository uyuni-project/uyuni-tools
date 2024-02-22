// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

const postUpgradeScriptTemplate = `#!/bin/bash
set -e
{{ if .CobblerHost }}
sed 's/cobbler\.host.*/cobbler\.host = {{ .CobblerHost }}/' -i /etc/rhn/rhn.conf;
sed 's/redhat_management_server.*/redhat_management_server = {{ .CobblerHost }}/' -i /etc/cobbler/settings.yaml;
{{ end }}
`

// PostUpgradeTemplateData represents information used to create post upgrade.
type PostUpgradeTemplateData struct {
	CobblerHost string
}

// Render will create script for finalizing PostgreSQL upgrade.
func (data PostUpgradeTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(postUpgradeScriptTemplate))
	return t.Execute(wr, data)
}
