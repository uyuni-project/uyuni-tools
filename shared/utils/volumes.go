// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

// This map should match the volumes mapping from the container definition in both
// the helm chart and the systemctl services definitions
var VOLUMES = map[string]string{
	"var-lib-cobbler":       "/var/lib/cobbler",
	"var-pgsql":             "/var/lib/pgsql",
	"var-cache":             "/var/cache",
	"var-spacewalk":         "/var/spacewalk",
	"var-log":               "/var/log",
	"srv-salt":              "/srv/salt",
	"srv-www-pub":           "/srv/www/htdocs/pub",
	"srv-www-cobbler":       "/srv/www/cobbler",
	"srv-www-osimages":      "/srv/www/os-images",
	"srv-www-distributions": "/srv/www/distributions",
	"srv-tftpboot":          "/srv/tftpboot",
	"srv-formulametadata":   "/srv/formula_metadata",
	"srv-pillar":            "/srv/pillar",
	"srv-susemanager":       "/srv/susemanager",
	"srv-spacewalk":         "/srv/spacewalk",
	"root":                  "/root",
	"etc-apache2":           "/etc/apache2",
	"etc-rhn":               "/etc/rhn",
	"etc-systemd-multi":     "/etc/systemd/system/multi-user.target.wants",
	"etc-systemd-sockets":   "/etc/systemd/system/sockets.target.wants",
	"etc-salt":              "/etc/salt",
	"etc-tomcat":            "/etc/tomcat",
	"etc-cobbler":           "/etc/cobbler",
	"etc-sysconfig":         "/etc/sysconfig",
	"etc-tls":               "/etc/pki/tls",
	"etc-postfix":           "/etc/postfix",
	"ca-cert":               "/etc/pki/trust/anchors",
}

var PROXY_HTTPD_VOLUMES = map[string]string{
	"uyuni-proxy-rhn-cache": "/var/cache/rhn",
	"uyuni-proxy-tftpboot":  "/srv/tftpboot",
}

var PROXY_SQUID_VOLUMES = map[string]string{
	"uyuni-proxy-squid-cache": "/var/cache/squid",
}

var PROXY_TFTPD_VOLUMES = map[string]string{
	"uyuni-proxy-tftpboot": "/srv/tftpboot:ro",
}
