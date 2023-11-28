// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

const podmanMigrationScriptTemplate = `#!/bin/bash
set -e
SSH_CONFIG=""
if test -e /tmp/ssh_config; then
  SSH_CONFIG="-F /tmp/ssh_config"
fi
SSH="ssh -A $SSH_CONFIG "
SCP="scp -A $SSH_CONFIG "

echo "Stopping spacewalk service..."
$SSH {{ .SourceFqdn }} "spacewalk-service stop ; systemctl start postgresql.service"

$SSH {{ .SourceFqdn }} \
 "echo \"COPY (SELECT MIN(CONCAT(org_id, '-', label)) AS target, base_path FROM rhnKickstartableTree GROUP BY base_path) TO STDOUT WITH CSV;\" \
 |spacewalk-sql --select-mode - " > distros

echo "Stopping posgresql service..."
$SSH {{ .SourceFqdn }} "systemctl stop postgresql.service"

while IFS="," read -r target path ; do
    echo "-/ $path"
done < distros > exclude_list

for folder in {{ range .Volumes }}{{ . }} {{ end }};
do
  if $SSH {{ .SourceFqdn }} test -e $folder; then
    echo "Copying $folder..."
    rsync -e "$SSH" --rsync-path='sudo rsync' -avz -f "merge exclude_list" {{ .SourceFqdn }}:$folder/ $folder;
  else
    echo "Skipping missing $folder..."
  fi
done;

echo "Migrating auto-installable distributions..."
while IFS="," read -r target path ; do
  if $SSH -A {{ .SourceFqdn }} test -e $path; then
    echo "Copying distribution $target from $path"
    mkdir -p "/srv/www/distributions/$target"
    rsync -e "$SSH" --rsync-path='sudo rsync' -avz "{{ .SourceFqdn }}:$path/" "/srv/www/distributions/$target"
  else
    echo "Skipping missing distribution $path..."
  fi
done < distros

rm -f /srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT;
ln -s /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT /srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT;

echo "Extracting time zone..."
$SSH {{ .SourceFqdn }} timedatectl show -p Timezone >/var/lib/uyuni-tools/data

echo "Altering configuration for domain resolution..."
sed 's/report_db_host = {{ .SourceFqdn }}/report_db_host = localhost/' -i /etc/rhn/rhn.conf;
sed 's/server\.jabber_server/java\.hostname/' -i /etc/rhn/rhn.conf;
sed 's/client_use_localhost: false/client_use_localhost: true/' -i /etc/cobbler/settings.yaml;

{{ if .Kubernetes }}
echo "Altering configuration for kubernetes..."
echo 'server.no_ssl = 1' >> /etc/rhn/rhn.conf;
sed 's/address=[^:]*:/address=*:/' -i /etc/rhn/taskomatic.conf;

if test ! -f /etc/tomcat/conf.d/remote_debug.conf -a -f /etc/sysconfig/tomcat; then
  mv /etc/sysconfig/tomcat /etc/tomcat/conf.d/remote_debug.conf
fi

sed 's/address=[^:]*:/address=*:/' -i /etc/tomcat/conf.d/remote_debug.conf

echo "Extracting SSL certificate and authority"
extractedSSL=
if test -d /root/ssl-build; then
  # We may have an old unused ssl-build folder, check if the CA matches the deployed one
  buildCaFingerprint=
  if test -e /root/ssl-build/RHN-ORG-TRUSTED-SSL-CERT; then
    buildCaFingerprint=$(openssl x509 -in /root/ssl-build/RHN-ORG-TRUSTED-SSL-CERT -noout -fingerprint)
  fi
  caFingerprint=$(openssl x509 -in /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT -noout -fingerprint)

  if test "$buildCaFingerprint" == "$caFingerprint"; then
    echo "Extracting SSL Root CA key..."
    # Extract the SSL CA certificate and key.
    # The server certificate will be auto-generated by cert-manager using it, so no need to copy it.
    cp /root/ssl-build/RHN-ORG-PRIVATE-SSL-KEY /var/lib/uyuni-tools/

    extractedSSL="1"
  fi
fi

# This Root CA file is common to both cases
cp /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT /var/lib/uyuni-tools/RHN-ORG-TRUSTED-SSL-CERT

if test "extractedSSL" != "1"; then
  # For third party certificates, the CA chain is in the certificate file.
  $SCP {{ .SourceFqdn }}:/etc/pki/tls/private/spacewalk.key /var/lib/uyuni-tools/
  $SCP {{ .SourceFqdn }}:/etc/pki/tls/certs/spacewalk.crt /var/lib/uyuni-tools/
fi

echo "Removing useless ssl-build folder..."
rm -rf /root/ssl-build

# The content of this folder will be a RO mount from a configmap
rm /etc/pki/trust/anchors/*
{{ end }}
echo "DONE"`

type MigrateScriptTemplateData struct {
	Volumes    map[string]string
	SourceFqdn string
	Kubernetes bool
}

func (data MigrateScriptTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(podmanMigrationScriptTemplate))
	return t.Execute(wr, data)
}
