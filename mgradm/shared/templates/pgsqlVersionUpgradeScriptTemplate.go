// SPDX-FileCopyrightText: 2026 SUSE LLC
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

OLD_VERSION={{ .OldVersion }}
NEW_VERSION={{ .NewVersion }}

echo "Testing presence of postgresql$NEW_VERSION..."
test -d /usr/lib/postgresql$NEW_VERSION/bin
echo "Testing presence of postgresql$OLD_VERSION..."
test -d /usr/lib/postgresql$OLD_VERSION/bin

BACKUP_DIR={{ .BackupDir }}/backup

restore_database() {
    echo "Migration failed. Restoring original database..."
    if [ -d "$BACKUP_DIR" ] && [ "$(ls -A "$BACKUP_DIR")" ]; then
        if [ -f "$BACKUP_DIR/postgresql.conf.bak" ]; then
            echo "Restoring postgresql.conf..."
            mv "$BACKUP_DIR/postgresql.conf.bak" "$BACKUP_DIR/postgresql.conf"
        fi
        if [ -f "$BACKUP_DIR/pg_hba.conf.bak" ]; then
            echo "Restoring pg_hba.conf..."
            mv "$BACKUP_DIR/pg_hba.conf.bak" "$BACKUP_DIR/pg_hba.conf"
        fi
        echo "Cleaning up /var/lib/pgsql/data..."
        rm -rf /var/lib/pgsql/data/*
        echo "Restoring from $BACKUP_DIR..."
        shopt -s dotglob
        mv "$BACKUP_DIR"/* /var/lib/pgsql/data/
        shopt -u dotglob
        echo "Database restored."
    fi
}

trap restore_database ERR


echo "Create a database backup at $BACKUP_DIR"
test -d "$BACKUP_DIR" && mv "$BACKUP_DIR" "${BACKUP_DIR}$(date '+%Y%m%d_%H%M%S')"
mkdir -p "$BACKUP_DIR"
chown postgres:postgres "$BACKUP_DIR"
chmod 700 "$BACKUP_DIR"
shopt -s dotglob
mv /var/lib/pgsql/data/* "$BACKUP_DIR"
shopt -u dotglob

echo "Create new database directory..."
chown -R postgres:postgres /var/lib/pgsql

if [ -e /etc/pki/tls/private/pg-spacewalk.key ]; then
	echo "Enforce key permission"
	chown postgres:postgres /etc/pki/tls/private/pg-spacewalk.key
	chown postgres:postgres /etc/pki/tls/certs/spacewalk.crt
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
su -s /bin/bash - postgres -c "pg_checksums --disable --pgdata /var/lib/pgsql/data"

echo "Successfully initialized new postgresql $NEW_VERSION database."

echo "Temporarily disable SSL in the old posgresql configuration"
cp "${BACKUP_DIR}/postgresql.conf" "${BACKUP_DIR}/postgresql.conf.bak"
sed 's/^ssl/#ssl/' -i "${BACKUP_DIR}/postgresql.conf"

echo "Temporarily allow local postgres user connection "
cp "${BACKUP_DIR}/pg_hba.conf" "${BACKUP_DIR}/pg_hba.conf.bak"
echo "local all postgres trust" >> "${BACKUP_DIR}/pg_hba.conf"
cp "/var/lib/pgsql/data/pg_hba.conf" "/var/lib/pgsql/data/pg_hba.conf.bak"
echo "local all postgres trust" >> "/var/lib/pgsql/data/pg_hba.conf"

su -s /bin/bash - postgres -c "pg_upgrade --old-bindir=/usr/lib/postgresql$OLD_VERSION/bin --new-bindir=/usr/lib/postgresql$NEW_VERSION/bin --old-datadir=\"$BACKUP_DIR\" --new-datadir=/var/lib/pgsql/data"

echo "Enable SSL again"
mv "${BACKUP_DIR}/postgresql.conf.bak" "${BACKUP_DIR}/postgresql.conf"

echo "Restore pg_hba.conf and postgresql.conf"
mv "${BACKUP_DIR}/pg_hba.conf.bak" "${BACKUP_DIR}/pg_hba.conf"
cp "${BACKUP_DIR}/pg_hba.conf" "/var/lib/pgsql/data/pg_hba.conf"
cp "${BACKUP_DIR}/postgresql.conf" "/var/lib/pgsql/data/postgresql.conf"

echo "Reenabling checksums"
su -s /bin/bash - postgres -c "pg_checksums --enable --pgdata /var/lib/pgsql/data"

echo "DONE"`

// PostgreSQLVersionUpgradeTemplateData represents information used to create PostgreSQL upgrade script.
type PostgreSQLVersionUpgradeTemplateData struct {
	OldVersion string
	NewVersion string
	BackupDir  string
}

// Render will create PostgreSQL upgrade script.
func (data PostgreSQLVersionUpgradeTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(postgreSQLVersionUpgradeScriptTemplate))
	return t.Execute(wr, data)
}
