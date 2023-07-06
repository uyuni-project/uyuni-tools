package migrate

import (
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func getSshAuthSocket() string {
	path := os.Getenv("SSH_AUTH_SOCK")
	if len(path) == 0 {
		log.Fatal("SSH_AUTH_SOCK is not defined, start an ssh agent and try again")
	}
	return path
}

func generateMigrationScript(sourceFqdn string) string {
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
echo "DONE"`

	model := struct {
		Volumes    map[string]string
		SourceFqdn string
	}{
		Volumes:    VOLUMES,
		SourceFqdn: sourceFqdn,
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

// This map should match the volumes mapping from the container definition in both
// the helm chart and the systemctl services definitions
var VOLUMES = map[string]string{
	"var-lib-cobbler":     "/var/lib/cobbler",
	"var-pgsql":           "/var/lib/pgsql",
	"var-cache":           "/var/cache",
	"var-spacewalk":       "/var/spacewalk",
	"var-log":             "/var/log",
	"srv-salt":            "/srv/salt",
	"srv-www-pub":         "/srv/www/htdocs/pub",
	"srv-www-cobbler":     "/srv/www/cobbler",
	"srv-www-osimages":    "/srv/www/os-images",
	"srv-tftpboot":        "/srv/tftpboot",
	"srv-formulametadata": "/srv/formula_metadata",
	"srv-pillar":          "/srv/pillar",
	"srv-susemanager":     "/srv/susemanager",
	"srv-spacewalk":       "/srv/spacewalk",
	"root":                "/root",
	"etc-apache2":         "/etc/apache2",
	"etc-rhn":             "/etc/rhn",
	"etc-systemd":         "/etc/systemd/system/multi-user.target.wants",
	"etc-salt":            "/etc/salt",
	"etc-tomcat":          "/etc/tomcat",
	"etc-cobbler":         "/etc/cobbler",
	"etc-sysconfig":       "/etc/sysconfig",
	"etc-tls":             "/etc/pki/tls",
	"ca-cert":             "/etc/pki/trust/anchors/",
}
