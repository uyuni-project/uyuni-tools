package install

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const SETUP_NAME = "setup.sh"

// generateSetupScript creates a temporary folder with the setup script to execute in the container.
// The script exports all the needed environment variables and calls uyuni's mgr-setup.
// Podman or kubernetes-specific variables can be passed using extraEnv parameter.
func generateSetupScript(viper *viper.Viper, fqdn string, extraEnv map[string]string) string {
	localHostValues := []string{
		"localhost",
		"127.0.0.1",
		"::1",
		fqdn,
	}

	localDb := utils.Contains(localHostValues, viper.GetString("db.host"))

	dbHost := viper.GetString("db.host")
	reportdbHost := viper.GetString("reportdb.host")

	if localDb {
		// For now the setup script expects the localhost value for local DB
		// but the FQDN is required for the report db even if it's local
		dbHost = "localhost"
		if viper.GetString("reportdb.host") == "" {
			reportdbHost = fqdn
		}
	}
	env := map[string]string{
		"UYUNI_FQDN":            fqdn,
		"MANAGER_USER":          viper.GetString("db.user"),
		"MANAGER_PASS":          viper.GetString("db.password"),
		"MANAGER_ADMIN_EMAIL":   viper.GetString("email"),
		"MANAGER_MAIL_FROM":     viper.GetString("emailFrom"),
		"MANAGER_ENABLE_TFTP":   boolToString(viper.GetBool("enableTftp")),
		"LOCAL_DB":              boolToString(localDb),
		"MANAGER_DB_NAME":       viper.GetString("db.name"),
		"MANAGER_DB_HOST":       dbHost,
		"MANAGER_DB_PORT":       strconv.Itoa(viper.GetInt("db.port")),
		"MANAGER_DB_PROTOCOL":   viper.GetString("db.protocol"),
		"REPORT_DB_NAME":        viper.GetString("reportdb.name"),
		"REPORT_DB_HOST":        reportdbHost,
		"REPORT_DB_PORT":        strconv.Itoa(viper.GetInt("reportdb.port")),
		"REPORT_DB_USER":        viper.GetString("reportdb.user"),
		"REPORT_DB_PASS":        viper.GetString("reportdb.password"),
		"EXTERNALDB_ADMIN_USER": viper.GetString("db.admin.user"),
		"EXTERNALDB_ADMIN_PASS": viper.GetString("db.admin.password"),
		"EXTERNALDB_PROVIDER":   viper.GetString("db.provider"),
		"ISS_PARENT":            viper.GetString("issParent"),
		"MIRROR_PATH":           viper.GetString("mirrorPath"),
		"ACTIVATE_SLP":          "N", // Deprecated, will be removed soon
		"SCC_USER":              viper.GetString("scc.user"),
		"SCC_PASS":              viper.GetString("scc.password"),
	}

	// Add the extra environment variables
	for key, value := range extraEnv {
		env[key] = value
	}

	scriptDir, err := os.MkdirTemp("", "uyuniadm-*")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %s\n", err)
	}

	const scriptTemplate = `#!/bin/sh
{{- range $name, $value := .Env }}
export {{ $name }}={{ $value }}
{{- end }}

/usr/lib/susemanager/bin/mgr-setup -s -n

# clean before leaving
rm $0`

	model := struct {
		Env map[string]string
	}{
		Env: env,
	}

	file, err := os.OpenFile(filepath.Join(scriptDir, SETUP_NAME), os.O_WRONLY|os.O_CREATE, 0555)
	if err != nil {
		log.Fatalf("Fail to open setup script: %s\n", err)
	}
	defer file.Close()

	t := template.Must(template.New("script").Parse(scriptTemplate))
	if err = t.Execute(file, model); err != nil {
		log.Fatalf("Failed to generate setup script: %s\n", err)
	}

	return scriptDir
}
