// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
	//	"github.com/uyuni-project/uyuni-tools/shared/types"
)

const pgsqlConfigTemplate = `#!/bin/sh -x
POSTGRESQL=/var/lib/pgsql/data/postgresql.conf
SSL_CERT=/etc/pki/tls/certs/spacewalk.crt
SSL_KEY=/etc/pki/tls/private/pg-spacewalk.key

postgres_reconfig() {
    if grep -E "^$1[[:space:]]*=" $POSTGRESQL >/dev/null; then
        sed -i "s|^$1[[:space:]]*=.*|$1 = $2|" $POSTGRESQL
    else
        echo "$1 = $2" >> $POSTGRESQL
    fi
}


postgres_reconfig effective_cache_size 1152MB
postgres_reconfig maintenance_work_mem 96MB
postgres_reconfig max_connections 600
postgres_reconfig shared_buffers 384MB
postgres_reconfig wal_buffers 4MB
postgres_reconfig work_mem 2560kB
postgres_reconfig jit off

if [ -f $SSL_KEY ] ; then
    chown postgres $SSL_KEY
    chmod 400 $SSL_KEY
    postgres_reconfig "ssl" "on"
    postgres_reconfig "ssl_cert_file" "'$SSL_CERT'"
    postgres_reconfig "ssl_key_file" "'$SSL_KEY'"
fi

echo "postgresql.conf updated"
`

// PgsqlConfigTemplateData POD information to create systemd file.
type PgsqlConfigTemplateData struct {
}

// Render will create the systemd configuration file.
func (data PgsqlConfigTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(pgsqlConfigTemplate))
	return t.Execute(wr, data)
}
