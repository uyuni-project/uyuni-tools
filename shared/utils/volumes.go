// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import "github.com/uyuni-project/uyuni-tools/shared/types"

// EtcRhnVolumeMount defines the /etc/rhn volume mount.
var EtcRhnVolumeMount = types.VolumeMount{MountPath: "/etc/rhn", Name: "etc-rhn", Size: "1Mi"}

// VarPgsqlDataVolumeMount defines the /var/lib/pgsql/data volume mount.
var VarPgsqlDataVolumeMount = types.VolumeMount{MountPath: "/var/lib/pgsql/data", Name: "var-pgsql", Size: "50Gi"}

// RootVolumeMount defines the /root volume mount.
var RootVolumeMount = types.VolumeMount{MountPath: "/root", Name: "root", Size: "1Mi"}

// PgsqlRequiredVolumeMounts represents volumes mount used by PostgreSQL.
var PgsqlRequiredVolumeMounts = []types.VolumeMount{
	VarPgsqlDataVolumeMount,
}

// CaCertVolumeMount represents volume for CA certificates.
var CaCertVolumeMount = types.VolumeMount{MountPath: "/etc/pki/trust/anchors/", Name: "ca-cert"}

// EtcTLSTmpVolumeMount represents temporary volume for SSL certificates.
var EtcTLSTmpVolumeMount = types.VolumeMount{MountPath: "/etc/pki/tls/", Name: "etc-tls", Size: "1Mi"}

// ServerVolumeMounts should match the volumes mapping from the container definition in both
// the helm chart and the systemctl services definitions.
var ServerVolumeMounts = []types.VolumeMount{
	{MountPath: "/var/lib/cobbler", Name: "var-cobbler", Size: "10Mi"},
	{MountPath: "/var/lib/rhn/search", Name: "var-search", Size: "10Gi"},
	{MountPath: "/var/lib/salt", Name: "var-salt", Size: "10Mi"},
	{MountPath: "/var/cache", Name: "var-cache", Size: "10Gi"},
	{MountPath: "/var/spacewalk", Name: "var-spacewalk", Size: "100Gi"},
	{MountPath: "/var/log", Name: "var-log", Size: "2Gi"},
	{MountPath: "/srv/salt", Name: "srv-salt", Size: "10Mi"},
	{MountPath: "/srv/www/", Name: "srv-www", Size: "100Gi"},
	{MountPath: "/srv/tftpboot", Name: "srv-tftpboot", Size: "300Mi"},
	{MountPath: "/srv/formula_metadata", Name: "srv-formulametadata", Size: "10Mi"},
	{MountPath: "/srv/pillar", Name: "srv-pillar", Size: "10Mi"},
	{MountPath: "/srv/susemanager", Name: "srv-susemanager", Size: "1Mi"},
	{MountPath: "/srv/spacewalk", Name: "srv-spacewalk", Size: "10Mi"},
	RootVolumeMount,
	CaCertVolumeMount,
	{MountPath: "/run/salt/master", Name: "run-salt-master"},
	{MountPath: "/etc/apache2", Name: "etc-apache2", Size: "1Mi"},
	{MountPath: "/etc/systemd/system/multi-user.target.wants", Name: "etc-systemd-multi", Size: "1Mi"},
	{MountPath: "/etc/systemd/system/sockets.target.wants", Name: "etc-systemd-sockets", Size: "1Mi"},
	{MountPath: "/etc/salt", Name: "etc-salt", Size: "1Mi"},
	{MountPath: "/etc/tomcat", Name: "etc-tomcat", Size: "1Mi"},
	{MountPath: "/etc/cobbler", Name: "etc-cobbler", Size: "1Mi"},
	{MountPath: "/etc/sysconfig", Name: "etc-sysconfig", Size: "20Mi"},
	{MountPath: "/etc/postfix", Name: "etc-postfix", Size: "1Mi"},
	{MountPath: "/etc/sssd", Name: "etc-sssd", Size: "1Mi"},
	EtcRhnVolumeMount,
}

// ServerMigrationVolumeMounts match server + postgres volume mounts, used for migration.
var ServerMigrationVolumeMounts = append(ServerVolumeMounts, VarPgsqlDataVolumeMount, EtcTLSTmpVolumeMount)

// DatabaseMigrationVolumeMounts match database + etc/rhn volume mounts, used for database migration.
var DatabaseMigrationVolumeMounts = []types.VolumeMount{EtcRhnVolumeMount, VarPgsqlDataVolumeMount}

// SalineVolumeMounts represents volumes used by Saline container.
var SalineVolumeMounts = []types.VolumeMount{
	{Name: "etc-salt", MountPath: "/etc/salt"},
	{Name: "run-salt-master", MountPath: "/run/salt/master"},
}

// ProxyHttpdVolumes volumes used by HTTPD in proxy.
var ProxyHttpdVolumes = []types.VolumeMount{
	{Name: "uyuni-proxy-rhn-cache", MountPath: "/var/cache/rhn"},
	{Name: "uyuni-proxy-tftpboot", MountPath: "/srv/tftpboot"},
}

// ProxySquidVolumes volumes used by Squid in  proxy.
var ProxySquidVolumes = []types.VolumeMount{
	{Name: "uyuni-proxy-squid-cache", MountPath: "/var/cache/squid"},
}

// ProxyTftpdVolumes used by TFTP in proxy.
var ProxyTftpdVolumes = []types.VolumeMount{
	{Name: "uyuni-proxy-tftpboot", MountPath: "/srv/tftpboot:ro"},
}
