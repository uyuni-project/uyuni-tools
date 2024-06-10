// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

const postgreSQLVersionUpgradeScriptTemplate = `#!/bin/bash
set -e
echo "PostgreSQL version upgrade"

OLD_VERSION={{ .OldVersion }}
NEW_VERSION={{ .NewVersion }}
FAST_UPGRADE=--link

echo "Testing presence of postgresql$NEW_VERSION..."
test -d /usr/lib/postgresql$NEW_VERSION/bin
echo "Testing presence of postgresql$OLD_VERSION..."
test -d /usr/lib/postgresql$OLD_VERSION/bin

echo "Create a backup at /var/lib/pgsql/data-pg$OLD_VERSION..."
mv /var/lib/pgsql/data /var/lib/pgsql/data-pg$OLD_VERSION
echo "Create new database directory..."
mkdir -p /var/lib/pgsql/data
chown -R postgres:postgres /var/lib/pgsql
echo "Enforce key permission"
chown -R postgres:postgres /etc/pki/tls/private/pg-spacewalk.key
chown -R postgres:postgres /etc/pki/tls/certs/spacewalk.crt

echo "Initialize new postgresql $NEW_VERSION database..."
. /etc/sysconfig/postgresql 2>/dev/null # Load locale for SUSE
PGHOME=$(getent passwd postgres | cut -d ":" -f6)
#. $PGHOME/.i18n 2>/dev/null # Load locale for Enterprise Linux
if [ -z $POSTGRES_LANG ]; then
    POSTGRES_LANG="en_US.UTF-8"
    [ ! -z $LC_CTYPE ] && POSTGRES_LANG=$LC_CTYPE
fi

echo "Running initdb using postgres user"
echo "Any suggested command from the console should be run using postgres user"
su -s /bin/bash - postgres -c "initdb -D /var/lib/pgsql/data --locale=$POSTGRES_LANG"
echo "Successfully initialized new postgresql $NEW_VERSION database."
su -s /bin/bash - postgres -c "pg_upgrade --old-bindir=/usr/lib/postgresql$OLD_VERSION/bin --new-bindir=/usr/lib/postgresql$NEW_VERSION/bin --old-datadir=/var/lib/pgsql/data-pg$OLD_VERSION --new-datadir=/var/lib/pgsql/data $FAST_UPGRADE"

cp /var/lib/pgsql/data-pg$OLD_VERSION/pg_hba.conf /var/lib/pgsql/data
cp /var/lib/pgsql/data-pg$OLD_VERSION/postgresql.conf /var/lib/pgsql/data/

echo "DONE"`

// PostgreSQLVersionUpgradeTemplateData represents information used to create PostgreSQL upgrade script.
type PostgreSQLVersionUpgradeTemplateData struct {
	OldVersion string
	NewVersion string
	Kubernetes bool
}

// Render will create PostgreSQL upgrade script.
func (data PostgreSQLVersionUpgradeTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(postgreSQLVersionUpgradeScriptTemplate))
	return t.Execute(wr, data)
}
