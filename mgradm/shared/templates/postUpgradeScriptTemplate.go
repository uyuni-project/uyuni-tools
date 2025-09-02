// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

const postUpgradeScriptTemplate = `
sed 's/cobbler\.host.*/cobbler\.host = localhost/' -i /etc/rhn/rhn.conf;
if grep -q uyuni_authentication_endpoint /etc/cobbler/settings.d/zz-uyuni.settings; then
	echo 'uyuni_authentication_endpoint: "http://localhost"' >> /etc/cobbler/settings.d/zz-uyuni.settings
else
	sed 's/uyuni_authentication_endpoint.*/uyuni_authentication_endpoint: http:\/\/localhost/' \
        -i /etc/cobbler/settings.d/zz-uyuni.settings;
fi

if grep -q pam_auth_service /etc/rhn/rhn.conf; then
	echo 'pam_auth_service = susemanager' >> /etc/rhn/rhn.conf
else
	sed 's/pam_auth_service.*/pam_auth_service = susemanager/' -i /etc/rhn/rhn.conf;
fi

if test -e /etc/sysconfig/prometheus-postgres_exporter/systemd/60-server.conf; then
        sed 's/\/etc\/postgres_exporter\//\/etc\/sysconfig\/prometheus-postgres_exporter\//' \
        -i /etc/sysconfig/prometheus-postgres_exporter/systemd/60-server.conf;
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
