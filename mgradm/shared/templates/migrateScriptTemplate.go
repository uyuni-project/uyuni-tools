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
const migrationScriptTemplate = `
set -e
SSH_CONFIG=""
if test -e /tmp/ssh_config; then
  SSH_CONFIG="-F /tmp/ssh_config"
fi
SSH="ssh -o User={{ .User }} -A $SSH_CONFIG "

if $SSH {{ .SourceFqdn }} "[[ ! -f /etc/susemanager-release && ! -f /etc/uyuni-release ]]"; then
  echo "Cannot find neither /etc/susemanager-release nor /etc/uyuni-release. Is the source a no-containerized server?"
  exit 1
fi

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
  if $SSH -n {{ .SourceFqdn }} test -e "$path" ; then
    echo "-/ $path"

    # protect the targets that can be already synced in --prepare phase
    echo "P/ /srv/www/distributions/$target"
  fi
done < distros > exclude_list

# exclude all config files which already exist and are not marked noreplace
rpm -qa --qf '[%{fileflags},%{filenames}\n]' |grep ",/etc/" | while IFS="," read -r flags path ; do
    # config(noreplace) is 1<<4 (from lib/rpmlib.h)
    if [ $(( $flags & 16 )) -eq 0 -a -f "$path" ] ; then
        echo "-/ $path" >> exclude_list
    fi
done

# No need to migrate zypper's cache
echo "-/ /var/cache/zypp/**" >> exclude_list

# Migrating the reposync cache files doesn't bring value and contains dangling symlinks (bsc#1231769)
echo "-/ /var/cache/rhn/reposync/**" >> exclude_list

# exclude mgr-sync configuration file, in this way it would be re-generated (bsc#1228685)
echo "-/ /root/.mgr-sync" >> exclude_list

# exclude tomcat default configuration. All settings should be store in /etc/tomcat/conf.d/
echo "-/ /etc/sysconfig/tomcat" >> exclude_list
echo "-/ /etc/tomcat/tomcat.conf" >> exclude_list

# exclude krb5 configuration. All settings should be store in /etc/rhn/krb5.conf.d
echo "-/ /etc/krb5.conf.d" >> exclude_list

# exclude /etc/rhn/krb5.conf.d configuration.
# This folder should not exists in 4.3, but it's created by the Dockerfile in a persisted folder
# This should prevents "rsync --delete" to delete it
echo "-/ /etc/rhn/krb5.conf.d" >> exclude_list

# exclude schema migration files
echo "-/ /etc/sysconfig/rhn/reportdb-schema-upgrade" >> exclude_list
echo "-/ /etc/sysconfig/rhn/schema-upgrade" >> exclude_list

# exclude lastlog - it is not needed and can be too large
echo "-/ /var/log/lastlog" >> exclude_list

# exclude systemd units as they will be recreated later
echo "-/ /etc/systemd/**" >> exclude_list

# Exclude py2*-compat-salt.conf as they can't work in the container
echo "-/ /etc/salt/master.d/py2*-compat-salt.conf" >> exclude_list

# uyuni issue #10055. Some old system might have this file
echo "-/ /etc/apache2/vhosts.d/cobbler.conf" >> exclude_list

$SSH {{ .SourceFqdn }} 'cat /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT >/etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT-nolink'

for folder in {{ range .Volumes }}{{ .MountPath }} {{ end }};
do
  RSYNC_ARGS=-l
  # Those folders used to have symlinks in the cloud images.
  if test "$folder" = "/var/cache" -o "$folder" = "/var/spacewalk" -o \
          "$folder" = "/var/lib/pgsql"; then
    RSYNC_ARGS=-L
  fi
  if $SSH {{ .SourceFqdn }} test -e $folder; then
    echo "Copying $folder..."
    rsync --delete -e "$SSH" --rsync-path='sudo rsync' -avz $RSYNC_ARGS --trust-sender -f 'merge exclude_list' {{ .SourceFqdn }}:$folder/ $folder;
  else
    echo "Skipping missing $folder..."
  fi
done;

spacewalk-service enable
if $SSH {{ .SourceFqdn }} systemctl is-enabled tftp.socket; then
  echo "Enabling tftp socket..."
  systemctl enable tftp.socket
fi

sed -i -e 's|appBase="webapps"|appBase="/usr/share/susemanager/www/tomcat/webapps"|' /etc/tomcat/server.xml
sed -i -e 's|DocumentRoot\s*"/srv/www/htdocs"|DocumentRoot "/usr/share/susemanager/www/htdocs"|' /etc/apache2/vhosts.d/vhost-ssl.conf

echo "Migrating auto-installable distributions..."

while IFS="," read -r target path ; do
  if $SSH -n {{ .SourceFqdn }} test -e $path ; then
    echo "Copying distribution $target from $path"
    mkdir -p "/srv/www/distributions/$target"
    rsync --delete -e "$SSH" --rsync-path='sudo rsync' -avz "{{ .SourceFqdn }}:$path/" "/srv/www/distributions/$target"
	# Adjust cobbler distro configuration
	for config in $(grep "$path/" -r /var/lib/cobbler/collections/distros/ -l); do
		sed "s|$path/|/srv/www/distributions/$target/|g" -i $config
	done
  else
    echo "Skipping missing distribution $path..."
  fi
done < distros

echo "Migrating auto-installation snippets..."
$SSH {{ .SourceFqdn }} "find /var/lib/cobbler/snippets/spacewalk/* -type d" > snippets_dirs
while read -r snippets_dir ; do
  if $SSH -n {{ .SourceFqdn }} test -e $snippets_dir; then
    echo "Copying autoinstallation snippets from $snippets_dir..."
    mkdir -p "$snippets_dir"
    rsync --delete -e "$SSH" --rsync-path='sudo rsync' -avz "{{ .SourceFqdn }}:$snippets_dir" "$snippets_dir";
  else
    echo "Skipping autoinstallation snippets from $snippets_dir.."
  fi
done < snippets_dirs

if $SSH {{ .SourceFqdn }} test -e /etc/tomcat/conf.d; then
  echo "Copying tomcat configuration.."
  mkdir -p /etc/tomcat/conf.d
  rsync --delete -e "$SSH" --rsync-path='sudo rsync' -avz {{ .SourceFqdn }}:/etc/tomcat/conf.d /etc/tomcat/;
else
  echo "Skipping tomcat configuration.."
fi

if $SSH {{ .SourceFqdn }} test -e /etc/krb5.conf.d; then
  echo "Copying krb5 configuration.."
  mkdir -p /etc/rhn/krb5.conf.d
  rsync --delete -e "$SSH" --rsync-path='sudo rsync' -avz {{ .SourceFqdn }}:/etc/krb5.conf.d /etc/rhn/;
else
  echo "Skipping krb5 configuration.."
fi

echo "Migrating monitoring configuration..."
if test -f /etc/sysconfig/prometheus-postgres_exporter; then
  echo "The server monitoring configuration is too old, reapply after the migration"
  rm /etc/sysconfig/prometheus-postgres_exporter
fi
declare -A monitoring_conf
monitoring_conf+=([/usr/lib/systemd/system/tomcat.service.d/jmx.conf]="/etc/sysconfig/tomcat/systemd/")
monitoring_conf+=([/usr/lib/systemd/system/taskomatic.service.d/jmx.conf]="/etc/sysconfig/taskomatic/systemd/")
monitoring_conf+=([/etc/postgres_exporter/postgres_exporter_queries.yaml]="/etc/sysconfig/prometheus-postgres_exporter/")
monitoring_conf+=([/etc/systemd/system/prometheus-postgres_exporter.service.d/60-server.conf]="/etc/sysconfig/prometheus-postgres_exporter/systemd/")
monitoring_conf+=([/etc/postgres_exporter/pg_passwd]="/etc/sysconfig/prometheus-postgres_exporter/")
for config_file in ${!monitoring_conf[@]}
do
  if $SSH {{ .SourceFqdn }} test -e ${config_file}; then
    mkdir -p "${monitoring_conf[${config_file}]}"
    rsync --delete -e "$SSH" --rsync-path='sudo rsync' -avz "{{ .SourceFqdn }}:${config_file}" "${monitoring_conf[${config_file}]}";
  fi
done
rm -f /etc/systemd/system/multi-user.target.wants/prometheus-postgres_exporter.service
if $SSH {{ .SourceFqdn }} systemctl is-enabled prometheus-postgres_exporter.service; then
  if test ! -e /etc/systemd/system/multi-user.target.wants/prometheus-postgres_exporter.service; then
    echo "Enabling prometheus-postgres_exporter service..."
    ln -s /usr/lib/systemd/system/prometheus-postgres_exporter.service /etc/systemd/system/multi-user.target.wants/prometheus-postgres_exporter.service
  fi
fi

echo "Migrating custom SSL CA certificates..."
rsync -e "$SSH" --rsync-path='sudo rsync' -avz --trust-sender \
    --exclude='/etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT' \
    --ignore-missing-args \
    {{ .SourceFqdn }}:/usr/share/pki/trust/anchors/ \
    {{ .SourceFqdn }}:/etc/pki/trust/anchors/ \
    /etc/pki/trust/anchors/

rm -f /srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT;
rm /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT
mv /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT-nolink /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT
ln -s /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT /srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT;

echo "Extracting time zone..."
$SSH {{ .SourceFqdn }} timedatectl show -p Timezone >/var/lib/uyuni-tools/data

echo "Extracting postgresql versions..."
echo "current_pg_version=$(cat /var/lib/pgsql/data/PG_VERSION)" >> /var/lib/uyuni-tools/data
echo "current_libc_version=2.31" >> /var/lib/uyuni-tools/data

grep '^db_user' /etc/rhn/rhn.conf | sed 's/[ \t]//g' >>/var/lib/uyuni-tools/data
grep '^db_password' /etc/rhn/rhn.conf | sed 's/[ \t]//g' >>/var/lib/uyuni-tools/data
grep '^db_name' /etc/rhn/rhn.conf | sed 's/[ \t]//g' >>/var/lib/uyuni-tools/data
grep '^db_port' /etc/rhn/rhn.conf | sed 's/[ \t]//g' >>/var/lib/uyuni-tools/data
grep '^report_db_user' /etc/rhn/rhn.conf | sed 's/[ \t]//g' >>/var/lib/uyuni-tools/data
grep '^report_db_password' /etc/rhn/rhn.conf | sed 's/[ \t]//g' >>/var/lib/uyuni-tools/data

$SSH {{ .SourceFqdn }} sh -c "systemctl list-unit-files | grep hub-xmlrpc-api | grep -q active && echo has_hubxmlrpc=true || echo has_hubxmlrpc=false" >>/var/lib/uyuni-tools/data
(test $($SSH {{ .SourceFqdn }} grep jdwp -r /etc/tomcat/conf.d/ /etc/rhn/taskomatic.conf | wc -l) -gt 0 && echo debug=true || echo debug=false) >>/var/lib/uyuni-tools/data

echo "Altering configuration for domain resolution..."
sed 's/^db_host = .*/db_host = {{ .DBHost }}/' -i /etc/rhn/rhn.conf;
sed 's/^report_db_host = .*/report_db_host = {{ .ReportDBHost }}/' -i /etc/rhn/rhn.conf;
if ! grep -q '^java.hostname *=' /etc/rhn/rhn.conf; then
    sed 's/server\.jabber_server/java\.hostname/' -i /etc/rhn/rhn.conf;
fi
echo 'client_use_localhost: true' >> /etc/cobbler/settings.d/zz-uyuni.settings;

echo "Altering configuration for container environment..."
sed 's/address=[^:]*:/address=*:/' -i /etc/rhn/taskomatic.conf;

echo "Altering tomcat configuration..."
sed 's/--add-modules java.annotation,com.sun.xml.bind://' -i /etc/tomcat/conf.d/*
sed 's/-XX:-UseConcMarkSweepGC//' -i /etc/tomcat/conf.d/*
test -f /etc/tomcat/conf.d/remote_debug.conf && sed 's/address=[^:]*:/address=*:/' -i /etc/tomcat/conf.d/remote_debug.conf

# Alter rhn.conf to ensure mirror is set to /mirror if set at all
sed 's/server.susemanager.fromdir =.*/server.susemanager.fromdir = \/mirror/' -i /etc/rhn/rhn.conf

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
  rsync --delete -e "$SSH" --rsync-path='sudo rsync' -avz {{ .SourceFqdn }}:/etc/pki/tls/private/spacewalk.key /var/lib/uyuni-tools/
  rsync --delete -e "$SSH" --rsync-path='sudo rsync' -avz {{ .SourceFqdn }}:/etc/pki/tls/certs/spacewalk.crt /var/lib/uyuni-tools/
fi

echo "Removing useless ssl-build folder..."
rm -rf /root/ssl-build

# The content of this folder will be a RO mount from a configmap
rm /etc/pki/trust/anchors/*
{{ end }}

echo "DONE"`

// MigrateScriptTemplateData represents migration information used to create migration script.
type MigrateScriptTemplateData struct {
	Volumes      []types.VolumeMount
	SourceFqdn   string
	User         string
	Kubernetes   bool
	Prepare      bool
	DBHost       string
	ReportDBHost string
}

// Render will create migration script.
func (data MigrateScriptTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(migrationScriptTemplate))
	return t.Execute(wr, data)
}
