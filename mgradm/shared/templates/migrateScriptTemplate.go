// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/types"
)

const migrationScriptTemplate = `#!/bin/bash
set -e
SSH_CONFIG=""
if test -e /tmp/ssh_config; then
  SSH_CONFIG="-F /tmp/ssh_config"
fi
SSH="ssh -o User={{ .User }} -A $SSH_CONFIG "

{{ if .Prepare }}
echo "Preparing migration..."
$SSH {{ .SourceFqdn }} "sudo systemctl start postgresql.service"
{{ else }}
echo "Stopping spacewalk service..."
$SSH {{ .SourceFqdn }} "sudo spacewalk-service stop ; sudo systemctl start postgresql.service"
{{ end }}

$SSH {{ .SourceFqdn }} \
 "echo \"COPY (SELECT MIN(CONCAT(org_id, '-', label)) AS target, base_path FROM rhnKickstartableTree GROUP BY base_path) TO STDOUT WITH CSV;\" \
 |sudo spacewalk-sql --select-mode - " > distros

{{ if not .Prepare }}
echo "Stopping posgresql service..."
$SSH {{ .SourceFqdn }} "sudo systemctl stop postgresql.service"
{{ end }}

while IFS="," read -r target path ; do
    echo "-/ $path"
done < distros > exclude_list

# exclude all config files which already exist and are not marked noreplace
rpm -qa --qf '[%{fileflags},%{filenames}\n]' |grep ",/etc/" | while IFS="," read -r flags path ; do
    # config(noreplace) is 1<<4 (from lib/rpmlib.h)
    if [ $(( $flags & 16 )) -eq 0 -a -f "$path" ] ; then
        echo "-/ $path" >> exclude_list
    fi
done

# exclude tomcat default configuration. All settings should be store in /etc/tomcat/conf.d/
echo "-/ /etc/sysconfig/tomcat" >> exclude_list
echo "-/ /etc/tomcat/tomcat.conf" >> exclude_list

# exclude schema migration files
echo "-/ /etc/sysconfig/rhn/reportdb-schema-upgrade" >> exclude_list
echo "-/ /etc/sysconfig/rhn/schema-upgrade" >> exclude_list

# exclude lastlog - it is not needed and can be too large
echo "-/ /var/log/lastlog" >> exclude_list

for folder in {{ range .Volumes }}{{ .MountPath }} {{ end }};
do
  if $SSH {{ .SourceFqdn }} test -e $folder; then
    echo "Copying $folder..."
    rsync -e "$SSH" --rsync-path='sudo rsync' -avz --trust-sender -f "merge exclude_list" {{ .SourceFqdn }}:$folder/ $folder;
  else
    echo "Skipping missing $folder..."
  fi
done;

sed -i -e 's|appBase="webapps"|appBase="/usr/share/susemanager/www/tomcat/webapps"|' /etc/tomcat/server.xml
sed -i -e 's|DocumentRoot\s*"/srv/www/htdocs"|DocumentRoot "/usr/share/susemanager/www/htdocs"|' /etc/apache2/vhosts.d/vhost-ssl.conf

echo "Migrating auto-installable distributions..."

while IFS="," read -r target path ; do
  if $SSH -n {{ .SourceFqdn }} test -e $path ; then
    echo "Copying distribution $target from $path"
    mkdir -p "/srv/www/distributions/$target"
    rsync -e "$SSH" --rsync-path='sudo rsync' -avz "{{ .SourceFqdn }}:$path/" "/srv/www/distributions/$target"
  else
    echo "Skipping missing distribution $path..."
  fi
done < distros

if $SSH {{ .SourceFqdn }} test -e /etc/tomcat/conf.d; then
  echo "Copying tomcat configuration.."
  mkdir -p /etc/tomcat/conf.d
  rsync -e "$SSH" --rsync-path='sudo rsync' -avz {{ .SourceFqdn }}:/etc/tomcat/conf.d /etc/tomcat/;
else
  echo "Skipping tomcat configuration.."
fi

rm -f /srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT;
ln -s /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT /srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT;

echo "Extracting time zone..."
$SSH {{ .SourceFqdn }} timedatectl show -p Timezone >/var/lib/uyuni-tools/data

echo "Extracting postgresql versions..."
echo "image_pg_version=$(rpm -qa --qf '%{VERSION}\n' 'name=postgresql[0-8][0-9]-server'  | cut -d. -f1 | sort -n | tail -1)" >> /var/lib/uyuni-tools/data
echo "current_pg_version=$(cat /var/lib/pgsql/data/PG_VERSION)" >> /var/lib/uyuni-tools/data

grep '^db_user' /etc/rhn/rhn.conf | sed 's/[ \t]//g' >>/var/lib/uyuni-tools/data
grep '^db_password' /etc/rhn/rhn.conf | sed 's/[ \t]//g' >>/var/lib/uyuni-tools/data
grep '^db_name' /etc/rhn/rhn.conf | sed 's/[ \t]//g' >>/var/lib/uyuni-tools/data
grep '^db_port' /etc/rhn/rhn.conf | sed 's/[ \t]//g' >>/var/lib/uyuni-tools/data

$SSH {{ .SourceFqdn }} sh -c "systemctl list-unit-files | grep hub-xmlrpc-api | grep -q active && echo has_hubxmlrpc=true || echo has_hubxmlrpc=false" >>/var/lib/uyuni-tools/data

echo "Altering configuration for domain resolution..."
sed 's/report_db_host = {{ .SourceFqdn }}/report_db_host = localhost/' -i /etc/rhn/rhn.conf;
sed 's/server\.jabber_server/java\.hostname/' -i /etc/rhn/rhn.conf;
sed 's/client_use_localhost: false/client_use_localhost: true/' -i /etc/cobbler/settings.yaml;

echo "Altering configuration for container environment..."
sed 's/address=[^:]*:/address=*:/' -i /etc/rhn/taskomatic.conf;

echo "Altering tomcat configuration..."
sed 's/--add-modules java.annotation,com.sun.xml.bind://' -i /etc/tomcat/conf.d/*
sed 's/-XX:-UseConcMarkSweepGC//' -i /etc/tomcat/conf.d/*
test -f /etc/tomcat/conf.d/remote_debug.conf && sed 's/address=[^:]*:/address=*:/' -i /etc/tomcat/conf.d/remote_debug.conf

{{ if .Kubernetes }}
echo 'server.no_ssl = 1' >> /etc/rhn/rhn.conf;
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
  rsync -e "$SSH" --rsync-path='sudo rsync' -avz {{ .SourceFqdn }}:/etc/pki/tls/private/spacewalk.key /var/lib/uyuni-tools/
  rsync -e "$SSH" --rsync-path='sudo rsync' -avz {{ .SourceFqdn }}:/etc/pki/tls/certs/spacewalk.crt /var/lib/uyuni-tools/

fi

echo "Removing useless ssl-build folder..."
rm -rf /root/ssl-build

# The content of this folder will be a RO mount from a configmap
rm /etc/pki/trust/anchors/*
{{ end }}

echo "DONE"`

// MigrateScriptTemplateData represents migration information used to create migration script.
type MigrateScriptTemplateData struct {
	Volumes    []types.VolumeMount
	SourceFqdn string
	User       string
	Kubernetes bool
	Prepare    bool
}

// Render will create migration script.
func (data MigrateScriptTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(migrationScriptTemplate))
	return t.Execute(wr, data)
}
