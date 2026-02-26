// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/uyuni-project/uyuni-tools/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// NewContainerInspector creates a new templates.InspectTemplateData for the big container inspection.
func NewContainerInspector() templates.InspectTemplateData {
	return templates.InspectTemplateData{
		Values: []types.InspectData{
			types.NewInspectData(
				"uyuni_release",
				"cat /etc/*release | grep 'Uyuni release' | cut -d ' ' -f3 || true"),
			types.NewInspectData(
				"suse_manager_release",
				`[ -f /etc/susemanager-release ] && sed 's/.*(\([0-9.]\+\).*/\1/g' /etc/susemanager-release || true`),
			types.NewInspectData(
				"fqdn",
				`sed -n '/^java\.hostname/{s/^java\.hostname[[:space:]]*=[[:space:]]*\(.*\)/\1/;p}' /etc/rhn/rhn.conf || true`),
			types.NewInspectData("pg_version",
				"(psql -V | awk '{print $3}' | cut -d. -f1) || true"),
			types.NewInspectData("libc_version", "ldd --version | head -n1 | sed 's/^ldd (GNU libc) //'"),
			types.NewInspectData("db_user",
				`sed -n '/^db_user/{s/^db_user[[:space:]]*=[[:space:]]*\(.*\)/\1/;p}' /etc/rhn/rhn.conf || true`),
			types.NewInspectData("db_password",
				`sed -n '/^db_password/{s/^db_password[[:space:]]*=[[:space:]]*\(.*\)/\1/;p}' /etc/rhn/rhn.conf || true`),
			types.NewInspectData("db_name",
				`sed -n '/^db_name/{s/^db_name[[:space:]]*=[[:space:]]*\(.*\)/\1/;p}' /etc/rhn/rhn.conf || true`),
			types.NewInspectData("db_port",
				`sed -n '/^db_port/{s/^db_port[[:space:]]*=[[:space:]]*\(.*\)/\1/;p}' /etc/rhn/rhn.conf || true`),
			types.NewInspectData("db_host",
				`sed -n '/^db_host/{s/^db_host[[:space:]]*=[[:space:]]*\(.*\)/\1/;p}' /etc/rhn/rhn.conf || true`),
			types.NewInspectData("report_db_host",
				`sed -n '/^report_db_host/{s/^report_db_host[[:space:]]*=[[:space:]]*\(.*\)/\1/;p}' /etc/rhn/rhn.conf || true`),
		},
	}
}

// ContainerInspectData are the data extracted by a server inspector.
type ContainerInspectData struct {
	DBUser             string `mapstructure:"db_user"`
	DBPassword         string `mapstructure:"db_password"`
	DBName             string `mapstructure:"db_name"`
	DBPort             int    `mapstructure:"db_port"`
	DBHost             string `mapstructure:"db_host"`
	ReportDBUser       string `mapstructure:"report_db_user"`
	ReportDBPassword   string `mapstructure:"report_db_password"`
	ReportDBHost       string `mapstructure:"report_db_host"`
	UyuniRelease       string `mapstructure:"uyuni_release"`
	SuseManagerRelease string `mapstructure:"suse_manager_release"`
	PgVersion          string `mapstructure:"pg_version"`
	LibcVersion        string `mapstructure:"libc_version"`
	Fqdn               string `mapstructure:"fqdn"`
}
