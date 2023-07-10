package migrate

import (
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func getSshAuthSocket() string {
	path := os.Getenv("SSH_AUTH_SOCK")
	if len(path) == 0 {
		log.Fatal("SSH_AUTH_SOCK is not defined, start an ssh agent and try again")
	}
	return path
}

func generateMigrationScript(sourceFqdn string, kubernetes bool) string {
	scriptDir, err := os.MkdirTemp("", "uyuniadm-*")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %s\n", err)
	}

	const scriptTemplate = `#!/bin/bash
set -e
for folder in {{range .Volumes}}{{.}} {{end}};
do
  rsync -e "ssh -A " --rsync-path='sudo rsync' -avz {{.SourceFqdn}}:$folder/ $folder;
done;
rm -f /srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT;
ln -s /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT /srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT;

{{ if .Kubernetes }}
echo 'server.no_ssl = 1' >> /etc/rhn/rhn.conf;
sed 's/address=[^:]*:/address=uyuni:/' -i /etc/rhn/taskomatic.conf;
sed 's/address=[^:]*:/address=uyuni:/' -i /etc/sysconfig/tomcat;
{{ end }}
echo "DONE"`

	model := struct {
		Volumes    map[string]string
		SourceFqdn string
		Kubernetes bool
	}{
		Volumes:    utils.VOLUMES,
		SourceFqdn: sourceFqdn,
		Kubernetes: kubernetes,
	}

	file, err := os.OpenFile(filepath.Join(scriptDir, "migrate.sh"), os.O_WRONLY|os.O_CREATE, 0555)
	if err != nil {
		log.Fatalf("Fail to open migration script: %s\n", err)
	}
	defer file.Close()

	t := template.Must(template.New("script").Parse(scriptTemplate))
	if err = t.Execute(file, model); err != nil {
		log.Fatalf("Failed to generate migration script: %s\n", err)
	}

	return scriptDir
}
