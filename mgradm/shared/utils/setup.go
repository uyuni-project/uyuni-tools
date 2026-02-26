// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"strconv"
)

// GetSetupEnv computes the environment variables required by the setup script from the flags.
func GetSetupEnv(mirror string, flags *InstallationFlags, fqdn string) map[string]string {
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
		"MANAGER_DB_NAME":     flags.DB.Name,
		"MANAGER_DB_HOST":     flags.DB.Host,
		"MANAGER_DB_PORT":     dbPort,
		"REPORT_DB_NAME":      flags.ReportDB.Name,
		"REPORT_DB_HOST":      flags.ReportDB.Host,
		"REPORT_DB_PORT":      reportdbPort,
		"EXTERNALDB_PROVIDER": flags.DB.Provider,
		"ISS_PARENT":          flags.IssParent,
		"DEBUG_JAVA":          strconv.FormatBool(flags.Debug.Java),
		"ORG_NAME":            flags.Organization,
		"ADMIN_USER":          flags.Admin.Login,
		"ADMIN_PASS":          flags.Admin.Password,
		"ADMIN_FIRST_NAME":    flags.Admin.FirstName,
		"ADMIN_LAST_NAME":     flags.Admin.LastName,
		"SCC_USER":            flags.SCC.User,
		"SCC_PASS":            flags.SCC.Password,
		"NO_SSL":              "N",
	}

	if mirror != "" {
		env["MIRROR_PATH"] = "/mirror"
	}

	return env
}
