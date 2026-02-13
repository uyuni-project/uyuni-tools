// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

//nolint:lll
const splitContainerSettingsScriptTemplate = `
set -e -x

{{ if .ReportDBHost }}
sed 's/^report_db_host = .*/report_db_host = {{ .ReportDBHost }}/' -i /etc/rhn/rhn.conf;
{{ end }}

{{ if .DBHost }}
sed 's/^db_host = .*/db_host = {{ .DBHost }}/' -i /etc/rhn/rhn.conf;
{{ end }}

echo "DONE"`

// SplitContainerSettingsScriptTemplateData represents the information that needs
// to be changed when db container is separated by the server.
type SplitContainerSettingsScriptTemplateData struct {
	DBHost       string
	ReportDBHost string
}

// Render will create migration script.
func (data SplitContainerSettingsScriptTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(splitContainerSettingsScriptTemplate))
	return t.Execute(wr, data)
}
