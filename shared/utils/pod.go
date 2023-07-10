package utils

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
