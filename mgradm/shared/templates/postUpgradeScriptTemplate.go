// SPDX-FileCopyrightText: 2025 SUSE LLC
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

echo "DONE"`

// PostUpgradeTemplateData represents information used to create post upgrade.
type PostUpgradeTemplateData struct {
}

// Render will create script for finalizing PostgreSQL upgrade.
func (data PostUpgradeTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(postUpgradeScriptTemplate))
	return t.Execute(wr, data)
}
