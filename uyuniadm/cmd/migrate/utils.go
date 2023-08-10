package migrate

import (
	"log"
	"os"
	"path/filepath"

	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/templates"
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

	data := templates.MigrateScriptTemplateData{
		Volumes:    utils.VOLUMES,
		SourceFqdn: sourceFqdn,
		Kubernetes: kubernetes,
	}

	scriptPath := filepath.Join(scriptDir, "migrate.sh")
	if err = utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		log.Fatalf("Failed to generate migration script: %s\n", err)
	}

	return scriptDir
}
