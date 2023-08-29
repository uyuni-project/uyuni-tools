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
	"etc-postfix":         "/etc/postfix",
	"ca-cert":             "/etc/pki/trust/anchors/",
}

type PortMap struct {
	Name     string
	Exposed  int
	Port     int
	Protocol string
}

func newPortMap(name string, exposed int, port int) PortMap {
	return PortMap{
		Name:    name,
		Exposed: exposed,
		Port:    port,
	}
}

var TCP_PORTS = []PortMap{
	newPortMap("postgres", 5432, 5432),
	newPortMap("salt-publish", 4505, 4505),
	newPortMap("salt-request", 4506, 4506),
	newPortMap("cobbler", 25151, 25151),
	newPortMap("tomcat-debug", 8000, 8080),
	newPortMap("tasko-debug", 8001, 8081),
	newPortMap("psql-metrics", 9187, 9187),
	newPortMap("node-metrics", 9101, 9101),
}

var UDP_PORTS = []PortMap{
	newPortMap("tftp", 69, 69),
}
