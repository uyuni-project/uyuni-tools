package templates

import (
	"io"
	"text/template"
)

const podmanMigrationScriptTemplate = `#!/bin/bash
set -e
for folder in {{range .Volumes}}{{.}} {{end}};
do
  rsync -e "ssh -A " --rsync-path='sudo rsync' -avz {{.SourceFqdn}}:$folder/ $folder;
done;
rm -f /srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT;
ln -s /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT /srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT;

ssh {{ .SourceFqdn }} timedatectl show -p Timezone >/var/lib/uyuni-tools/data

{{ if .Kubernetes }}
echo 'server.no_ssl = 1' >> /etc/rhn/rhn.conf;
sed 's/address=[^:]*:/address=uyuni:/' -i /etc/rhn/taskomatic.conf;
sed 's/address=[^:]*:/address=uyuni:/' -i /etc/sysconfig/tomcat;
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
