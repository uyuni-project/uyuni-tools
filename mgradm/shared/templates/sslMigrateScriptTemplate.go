// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/types"
)

//nolint:lll
const sslMigrationScriptTemplate = `
set -e
SSH_CONFIG=""
if test -e /tmp/ssh_config; then
  SSH_CONFIG="-F /tmp/ssh_config"
fi
SSH="ssh -o User={{ .User }} -A $SSH_CONFIG "
$SSH {{ .SourceFqdn }} 'cat /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT >/etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT-nolink'

for folder in {{ range .Volumes }}{{ .MountPath }} {{ end }};
do
  if $SSH {{ .SourceFqdn }} test -e $folder; then
    echo "Copying $folder..."
    rsync --delete -e "$SSH" --rsync-path='sudo rsync' -avzl --trust-sender {{ .SourceFqdn }}:$folder/ $folder;
  else
    echo "Skipping missing $folder..."
  fi
done;

rm -f /srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT;
rm /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT
mv /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT-nolink /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT
ln -s /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT /srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT;

if test -d /root/ssl-build; then
  # We may have an old unused ssl-build folder, check if the CA matches the deployed one
  buildCaFingerprint=
  if test -e /root/ssl-build/RHN-ORG-TRUSTED-SSL-CERT; then
    buildCaFingerprint=$(openssl x509 -in /root/ssl-build/RHN-ORG-TRUSTED-SSL-CERT -noout -fingerprint)
  fi
  caFingerprint=$(openssl x509 -in /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT -noout -fingerprint)

  if test "$buildCaFingerprint" != "$caFingerprint"; then
    echo "Removing unused ssl-build folder"
    rm -r /root/ssl-build/
  fi
fi

echo "Extracting data..."
$SSH {{ .SourceFqdn }} timedatectl show -p Timezone >/var/lib/uyuni-tools/data
$SSH {{ .SourceFqdn }} 'echo "current_pg_version=$(cat /var/lib/pgsql/data/PG_VERSION)"' >> /var/lib/uyuni-tools/data
echo "current_libc_version=2.31" >> /var/lib/uyuni-tools/data

$SSH {{ .SourceFqdn }} 'grep "^db_user" /etc/rhn/rhn.conf | sed "s/[ \t]//g"' >>/var/lib/uyuni-tools/data
$SSH {{ .SourceFqdn }} 'grep "^db_password" /etc/rhn/rhn.conf | sed "s/[ \t]//g"' >>/var/lib/uyuni-tools/data
$SSH {{ .SourceFqdn }} 'grep "^db_name" /etc/rhn/rhn.conf | sed "s/[ \t]//g"' >>/var/lib/uyuni-tools/data
$SSH {{ .SourceFqdn }} 'grep "^db_port" /etc/rhn/rhn.conf | sed "s/[ \t]//g"' >>/var/lib/uyuni-tools/data
$SSH {{ .SourceFqdn }} 'grep "^report_db_user" /etc/rhn/rhn.conf | sed "s/[ \t]//g"' >>/var/lib/uyuni-tools/data
$SSH {{ .SourceFqdn }} 'grep "^report_db_password" /etc/rhn/rhn.conf | sed "s/[ \t]//g"' >>/var/lib/uyuni-tools/data

$SSH {{ .SourceFqdn }} "systemctl list-unit-files | grep hub-xmlrpc-api | grep -q active && echo has_hubxmlrpc=true || echo has_hubxmlrpc=false" >>/var/lib/uyuni-tools/data
(test $($SSH {{ .SourceFqdn }} grep jdwp -r /etc/tomcat/conf.d/ /etc/rhn/taskomatic.conf | wc -l) -gt 0 && echo debug=true || echo debug=false) >>/var/lib/uyuni-tools/data
`

// SSLMigrateScriptTemplateData represents migration information used to create the podman SSL migration script.
type SSLMigrateScriptTemplateData struct {
	Volumes    []types.VolumeMount
	SourceFqdn string
	User       string
}

// Render will create the ssl migration script.
func (data SSLMigrateScriptTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(sslMigrationScriptTemplate))
	return t.Execute(wr, data)
}
