// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

const postUpgradeScriptTemplate = `#!/bin/bash
{{ if .CobblerHost }}
sed 's/cobbler\.host.*/cobbler\.host = {{ .CobblerHost }}/' -i /etc/rhn/rhn.conf;
grep spacewalk_authentication_endpoint /etc/cobbler/settings.yaml
if [ $? -eq 1 ]; then
	echo 'spacewalk_authentication_endpoint: "http://localhost"' >> /etc/cobbler/settings.yaml
else
	sed 's/spacewalk_authentication_endpoint.*/spacewalk_authentication_endpoint: http:\/\/localhost/' -i /etc/cobbler/settings.yaml;
fi
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
