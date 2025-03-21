// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"strconv"
	"strings"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// GetSetupEnv computes the environment variables required by the setup script from the flags.
// As the requirements are slightly different for kubernetes there is a toggle parameter for it.
func GetSetupEnv(mirror string, flags *InstallationFlags, fqdn string, kubernetes bool) map[string]string {
	dbPort := "5432"
	if flags.DB.Port != 0 {
		dbPort = strconv.Itoa(flags.DB.Port)
	}

	reportdbPort := "5432"
	if flags.ReportDB.Port != 0 {
		reportdbPort = strconv.Itoa(flags.ReportDB.Port)
	}

	env := map[string]string{
		"UYUNI_FQDN":          fqdn,
		"MANAGER_ADMIN_EMAIL": flags.Email,
		"MANAGER_MAIL_FROM":   flags.EmailFrom,
		"MANAGER_ENABLE_TFTP": boolToString(flags.Tftp),
		"MANAGER_DB_NAME":     flags.DB.Name,
		"MANAGER_DB_HOST":     flags.DB.Host,
		"MANAGER_DB_PORT":     dbPort,
		"REPORT_DB_NAME":      flags.ReportDB.Name,
		"REPORT_DB_HOST":      flags.ReportDB.Host,
		"REPORT_DB_PORT":      reportdbPort,
		"EXTERNALDB_PROVIDER": flags.DB.Provider,
		"ISS_PARENT":          flags.IssParent,
	}

	if kubernetes {
		env["NO_SSL"] = "Y"
	} else {
		// Only add the credentials for podman as we have secret for Kubernetes.
		env["ADMIN_USER"] = flags.Admin.Login
		env["ADMIN_PASS"] = flags.Admin.Password
		env["SCC_USER"] = flags.SCC.User
		env["SCC_PASS"] = flags.SCC.Password
	}

	if mirror != "" {
		env["MIRROR_PATH"] = "/mirror"
	}

	return env
}

func boolToString(value bool) string {
	if value {
		return "Y"
	}
	return "N"
}

// GenerateSetupScript creates a temporary folder with the setup script to execute in the container.
// The script exports all the needed environment variables and calls uyuni's mgr-setup.
func GenerateSetupScript(flags *InstallationFlags, nossl bool) (string, error) {
	template := templates.MgrSetupScriptTemplateData{
		DebugJava:      flags.Debug.Java,
		OrgName:        flags.Organization,
		AdminLogin:     "$ADMIN_USER",
		AdminPassword:  "$ADMIN_PASS",
		AdminFirstName: flags.Admin.FirstName,
		AdminLastName:  flags.Admin.LastName,
		AdminEmail:     flags.Admin.Email,
		NoSSL:          nossl,
	}

	// Prepare the script
	scriptBuilder := new(strings.Builder)
	if err := template.Render(scriptBuilder); err != nil {
		return "", utils.Error(err, L("failed to render setup script"))
	}
	return scriptBuilder.String(), nil
}
