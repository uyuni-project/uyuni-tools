package shared

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/templates"
	adm_utils "github.com/uyuni-project/uyuni-tools/uyuniadm/shared/utils"
)

const SETUP_NAME = "setup.sh"

func RunSetup(cnx *utils.Connection, flags *InstallFlags, fqdn string, env map[string]string) {
	tmpFolder := generateSetupScript(flags, fqdn, env)
	defer os.RemoveAll(tmpFolder)

	utils.Copy(cnx, filepath.Join(tmpFolder, SETUP_NAME), "server:/tmp/setup.sh", "root", "root")

	err := adm_utils.ExecCommand(zerolog.InfoLevel, cnx, "/tmp/setup.sh")
	if err != nil {
		log.Fatal().Err(err).Msg("error running the setup script")
	}

	log.Info().Msg("Server set up")
}

// generateSetupScript creates a temporary folder with the setup script to execute in the container.
// The script exports all the needed environment variables and calls uyuni's mgr-setup.
// Podman or kubernetes-specific variables can be passed using extraEnv parameter.
func generateSetupScript(flags *InstallFlags, fqdn string, extraEnv map[string]string) string {
	localHostValues := []string{
		"localhost",
		"127.0.0.1",
		"::1",
		fqdn,
	}

	localDb := utils.Contains(localHostValues, flags.Db.Host)

	dbHost := flags.Db.Host
	reportdbHost := flags.ReportDb.Host

	if localDb {
		dbHost = "localhost"
		if reportdbHost == "" {
			reportdbHost = "localhost"
		}
	}
	env := map[string]string{
		"UYUNI_FQDN":            fqdn,
		"MANAGER_USER":          flags.Db.User,
		"MANAGER_PASS":          flags.Db.Password,
		"MANAGER_ADMIN_EMAIL":   flags.Email,
		"MANAGER_MAIL_FROM":     flags.EmailFrom,
		"MANAGER_ENABLE_TFTP":   boolToString(flags.Tftp),
		"LOCAL_DB":              boolToString(localDb),
		"MANAGER_DB_NAME":       flags.Db.Name,
		"MANAGER_DB_HOST":       dbHost,
		"MANAGER_DB_PORT":       strconv.Itoa(flags.Db.Port),
		"MANAGER_DB_PROTOCOL":   flags.Db.Protocol,
		"REPORT_DB_NAME":        flags.ReportDb.Name,
		"REPORT_DB_HOST":        reportdbHost,
		"REPORT_DB_PORT":        strconv.Itoa(flags.ReportDb.Port),
		"REPORT_DB_USER":        flags.ReportDb.User,
		"REPORT_DB_PASS":        flags.ReportDb.Password,
		"EXTERNALDB_ADMIN_USER": flags.Db.Admin.User,
		"EXTERNALDB_ADMIN_PASS": flags.Db.Admin.Password,
		"EXTERNALDB_PROVIDER":   flags.Db.Provider,
		"ISS_PARENT":            flags.IssParent,
		"ACTIVATE_SLP":          "N", // Deprecated, will be removed soon
		"SCC_USER":              flags.Scc.User,
		"SCC_PASS":              flags.Scc.Password,
	}
	if flags.MirrorPath != "" {
		env["MIRROR_PATH"] = "/mirror"
	}

	// Add the extra environment variables
	for key, value := range extraEnv {
		env[key] = value
	}

	scriptDir, err := os.MkdirTemp("", "uyuniadm-*")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create temporary directory")
	}

	dataTemplate := templates.MgrSetupScriptTemplateData{
		Env:       env,
		DebugJava: flags.Debug.Java,
	}

	scriptPath := filepath.Join(scriptDir, SETUP_NAME)
	if err = utils.WriteTemplateToFile(dataTemplate, scriptPath, 0555, true); err != nil {
		log.Fatal().Err(err).Msg("Failed to generate setup script")
	}

	return scriptDir
}

func boolToString(value bool) string {
	if value {
		return "Y"
	}
	return "N"
}
