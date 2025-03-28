// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

//nolint:lll
const postgreSQLVersionUpgradeScriptTemplate = `#!/bin/bash
set -e
echo "PostgreSQL version upgrade"

ls -la /var/lib/pgsql/data/
OLD_VERSION={{ .OldVersion }}
NEW_VERSION={{ .NewVersion }}

echo "Testing presence of postgresql$NEW_VERSION..."
test -d /usr/lib/postgresql$NEW_VERSION/bin
echo "Testing presence of postgresql$OLD_VERSION..."
test -d /usr/lib/postgresql$OLD_VERSION/bin

echo "Create a backup at /var/lib/pgsql/data-backup..."
mkdir -p /var/lib/pgsql/data-backup
chown postgres:postgres /var/lib/pgsql/data-backup
chmod 700 /var/lib/pgsql/data-backup
shopt -s dotglob
mv /var/lib/pgsql/data/* /var/lib/pgsql/data-backup

echo "Create new database directory..."
chown -R postgres:postgres /var/lib/pgsql

if [ -e /etc/pki/tls/private/pg-spacewalk.key ]; then
	echo "Enforce key permission"
	chown -R postgres:postgres /etc/pki/tls/private/pg-spacewalk.key
	chown -R postgres:postgres /etc/pki/tls/certs/spacewalk.crt
fi

echo "Initialize new postgresql $NEW_VERSION database..."
. /etc/sysconfig/postgresql 2>/dev/null # Load locale for SUSE
PGHOME=$(getent passwd postgres | cut -d ":" -f6)
if [ -z $POSTGRES_LANG ]; then
    POSTGRES_LANG="en_US.UTF-8"
    [ ! -z $LC_CTYPE ] && POSTGRES_LANG=$LC_CTYPE
fi

echo "Running initdb using postgres user"
echo "Any suggested command from the console should be run using postgres user"
su -s /bin/bash - postgres -c "initdb -D /var/lib/pgsql/data --locale=$POSTGRES_LANG"
echo "Successfully initialized new postgresql $NEW_VERSION database."

echo "Temporarily disable SSL in the old posgresql configuration"
cp /var/lib/pgsql/data-backup/postgresql.conf /var/lib/pgsql/data-backup/postgresql.conf.bak
sed 's/^ssl/#ssl/' -i /var/lib/pgsql/data-backup/postgresql.conf

su -s /bin/bash - postgres -c "pg_upgrade --old-bindir=/usr/lib/postgresql$OLD_VERSION/bin --new-bindir=/usr/lib/postgresql$NEW_VERSION/bin --old-datadir=/var/lib/pgsql/data-backup --new-datadir=/var/lib/pgsql/data"

echo "Enable SSL again"
cp /var/lib/pgsql/data-backup/postgresql.conf.bak /var/lib/pgsql/data-backup/postgresql.conf

echo "DONE"`

// PostgreSQLVersionUpgradeTemplateData represents information used to create PostgreSQL upgrade script.
type PostgreSQLVersionUpgradeTemplateData struct {
	OldVersion string
	NewVersion string
}

// Render will create PostgreSQL upgrade script.
func (data PostgreSQLVersionUpgradeTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(postgreSQLVersionUpgradeScriptTemplate))
	return t.Execute(wr, data)
}
