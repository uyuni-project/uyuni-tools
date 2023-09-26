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

type PortMap struct {
	Name     string
	Exposed  int
	Port     int
	Protocol string
}

func NewPortMap(name string, exposed int, port int) PortMap {
	return PortMap{
		Name:    name,
		Exposed: exposed,
		Port:    port,
	}
}

// The port names should be less than 15 characters long and lowercased for traefik to eat them
var TCP_PORTS = []PortMap{
	NewPortMap("postgres", 5432, 5432),
	NewPortMap("salt-publish", 4505, 4505),
	NewPortMap("salt-request", 4506, 4506),
	NewPortMap("cobbler", 25151, 25151),
	NewPortMap("psql-mtrx", 9187, 9187),
	NewPortMap("tasko-jmx-mtrx", 5556, 5556),
	NewPortMap("tomcat-jmx-mtrx", 5557, 5557),
}

var DEBUG_PORTS = []PortMap{
	// We can't expose on port 8000 since traefik already uses it
	NewPortMap("tomcat-debug", 8003, 8003),
	NewPortMap("tasko-debug", 8001, 8001),
	NewPortMap("search-debug", 8002, 8002),
}

var UDP_PORTS = []PortMap{
	{
		Name:     "tftp",
		Exposed:  69,
		Port:     69,
		Protocol: "udp",
	},
}
