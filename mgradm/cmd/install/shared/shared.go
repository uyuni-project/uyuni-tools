// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/org"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const setupName = "setup.sh"

// RunSetup execute the setup.
func RunSetup(cnx *shared.Connection, flags *InstallFlags, fqdn string, env map[string]string) error {
	// Containers should be running now, check storage if it is using volume from already configured server
	preconfigured := false
	if isServerConfigured(cnx) {
		log.Warn().Msg(
			L("Server appears to be already configured. Installation will continue, but installation options may be ignored."),
		)
		preconfigured = true
	}

	tmpFolder, err := generateSetupScript(flags, fqdn, env)
	defer os.RemoveAll(tmpFolder)
	if err != nil {
		return err
	}

	if err := cnx.Copy(filepath.Join(tmpFolder, setupName), "server:/tmp/setup.sh", "root", "root"); err != nil {
		return utils.Errorf(err, L("cannot copy /tmp/setup.sh"))
	}

	err = adm_utils.ExecCommand(zerolog.InfoLevel, cnx, "/tmp/setup.sh")
	if err != nil && !preconfigured {
		return utils.Errorf(err, L("error running the setup script"))
	}
	if err := cnx.CopyCaCertificate(fqdn); err != nil {
		return utils.Errorf(err, L("failed to add SSL CA certificate to host trusted certificates"))
	}

	// Call the org.createFirst api if flags are passed
	// This should not happen since the password is queried and enforced
	if flags.Admin.Password != "" {
		apiCnx := api.ConnectionDetails{
			Server:   fqdn,
			Insecure: false,
			User:     flags.Admin.Login,
			Password: flags.Admin.Password,
		}

		// Check if there is already admin user with given password and organization with same name
		client, err := api.Init(&apiCnx)
		if err != nil {
			log.Error().Err(err).Msgf(L("unable to prepare API client"))
		}
		if err = client.Login(); err == nil {
			if _, err := org.GetOrganizationDetails(&apiCnx, flags.Organization); err == nil {
				log.Info().Msgf(L("Server organization already exists, reusing"))
			} else {
				log.Debug().Err(err).Msg("Error returned by server")
				log.Warn().Msgf(L("Administration user already exists, but organization %s could not be found"), flags.Organization)
			}
		} else {
			var connError *url.Error
			if errors.As(err, &connError) {
				// We were not able to connect to the server at all
				return err
			}
			// We do not have any user existing, create one. CreateFirst skip user login
			_, err := org.CreateFirst(&apiCnx, flags.Organization, &flags.Admin)
			if err != nil {
				if preconfigured {
					log.Warn().Msgf(L("Administration user already exists, but provided credentials are not valid"))
				} else {
					return err
				}
			}
		}
	}

	log.Info().Msgf(L("Server set up, login on https://%[1]s with %[2]s user"), fqdn, flags.Admin.Login)
	return nil
}

// generateSetupScript creates a temporary folder with the setup script to execute in the container.
// The script exports all the needed environment variables and calls uyuni's mgr-setup.
// Podman or kubernetes-specific variables can be passed using extraEnv parameter.
func generateSetupScript(flags *InstallFlags, fqdn string, extraEnv map[string]string) (string, error) {
	localHostValues := []string{
		"localhost",
		"127.0.0.1",
		"::1",
		fqdn,
	}

	localDB := utils.Contains(localHostValues, flags.DB.Host)

	dbHost := flags.DB.Host
	reportdbHost := flags.ReportDB.Host

	if localDB {
		dbHost = "localhost"
		if reportdbHost == "" {
			reportdbHost = "localhost"
		}
	}
	env := map[string]string{
		"UYUNI_FQDN":            fqdn,
		"MANAGER_USER":          flags.DB.User,
		"MANAGER_PASS":          flags.DB.Password,
		"MANAGER_ADMIN_EMAIL":   flags.Email,
		"MANAGER_MAIL_FROM":     flags.EmailFrom,
		"MANAGER_ENABLE_TFTP":   boolToString(flags.Tftp),
		"LOCAL_DB":              boolToString(localDB),
		"MANAGER_DB_NAME":       flags.DB.Name,
		"MANAGER_DB_HOST":       dbHost,
		"MANAGER_DB_PORT":       strconv.Itoa(flags.DB.Port),
		"MANAGER_DB_PROTOCOL":   flags.DB.Protocol,
		"REPORT_DB_NAME":        flags.ReportDB.Name,
		"REPORT_DB_HOST":        reportdbHost,
		"REPORT_DB_PORT":        strconv.Itoa(flags.ReportDB.Port),
		"REPORT_DB_USER":        flags.ReportDB.User,
		"REPORT_DB_PASS":        flags.ReportDB.Password,
		"EXTERNALDB_ADMIN_USER": flags.DB.Admin.User,
		"EXTERNALDB_ADMIN_PASS": flags.DB.Admin.Password,
		"EXTERNALDB_PROVIDER":   flags.DB.Provider,
		"ISS_PARENT":            flags.IssParent,
		"ACTIVATE_SLP":          "N", // Deprecated, will be removed soon
		"SCC_USER":              flags.Scc.User,
		"SCC_PASS":              flags.Scc.Password,
	}
	if flags.Mirror != "" {
		env["MIRROR_PATH"] = "/mirror"
	}

	// Add the extra environment variables
	for key, value := range extraEnv {
		env[key] = value
	}

	scriptDir, err := utils.TempDir()
	if err != nil {
		return "", err
	}

	dataTemplate := templates.MgrSetupScriptTemplateData{
		Env:       env,
		DebugJava: flags.Debug.Java,
	}

	scriptPath := filepath.Join(scriptDir, setupName)
	if err = utils.WriteTemplateToFile(dataTemplate, scriptPath, 0555, true); err != nil {
		return "", utils.Errorf(err, L("Failed to generate setup script"))
	}

	return scriptDir, nil
}

func boolToString(value bool) string {
	if value {
		return "Y"
	}
	return "N"
}

func isServerConfigured(cnx *shared.Connection) bool {
	return cnx.TestExistenceInPod("/root/.MANAGER_SETUP_COMPLETE")
}
