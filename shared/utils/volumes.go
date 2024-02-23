// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import "github.com/uyuni-project/uyuni-tools/shared/types"

// PgsqlRequiredVolumeMounts represents volumes mount used by PostgreSQL.
var PgsqlRequiredVolumeMounts = []types.VolumeMount{
	{MountPath: "/var/lib/pgsql", Name: "var-pgsql"},
	{MountPath: "/etc/rhn", Name: "etc-rhn"},
}

// PgsqlRequiredVolumes represents volumes used by PostgreSQL.
var PgsqlRequiredVolumes = []types.Volume{
	{Name: "var-pgsql", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "var-pgsql"}},
	{Name: "etc-rhn", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-rhn"}},
}

// CACertsVolumeMounts represents volume mounts used to store certificates.
var CACertsVolumeMounts = []types.VolumeMount{
	{MountPath: "/etc/pki/tls", Name: "etc-tls"},
	{MountPath: "/etc/pki/spacewalk-tls", Name: "tls-key"},
	{MountPath: "/etc/pki/trust/anchors", Name: "ca-cert"},
}

// CACertsVolumes represents volumes used to store certificates.
var CACertsVolumes = []types.Volume{
	{Name: "etc-tls", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-tls"}},
	{Name: "tls-key",
		Secret: &types.Secret{
			SecretName: "uyuni-cert", Items: []types.SecretItem{
				{Key: "tls.crt", Path: "spacewalk.crt"},
				{Key: "tls.key", Path: "spacewalk.key"},
			},
		},
	},
	{Name: "ca-cert", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "ca-cert"}},
}

// MigrationVolumeMounts represents volume mounts during migration.
var MigrationVolumeMounts = []types.VolumeMount{
	{MountPath: "/var/lib/pgsql", Name: "var-pgsql"},
	{MountPath: "/etc/rhn", Name: "etc-rhn"},
	{MountPath: "/etc/apache2", Name: "etc-apache2"},
	{MountPath: "/etc/systemd/system/multi-user.target.wants", Name: "etc-systemd-multi"},
	{MountPath: "/etc/systemd/system/sockets.target.wants", Name: "etc-systemd-sockets"},
	{MountPath: "/etc/salt", Name: "etc-salt"},
	{MountPath: "/etc/tomcat", Name: "etc-tomcat"},
	{MountPath: "/etc/cobbler", Name: "etc-cobbler"},
	{MountPath: "/etc/sysconfig", Name: "etc-sysconfig"},
	{MountPath: "/etc/postfix", Name: "etc-postfix"},
	{MountPath: "/etc/pki/tls", Name: "etc-tls"},
}

// MigrationVolume represents volumes used during migration.
var MigrationVolumes = []types.Volume{
	{Name: "var-pgsql", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "var-pgsql"}},
	{Name: "etc-rhn", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-rhn"}},
	{Name: "etc-apache2", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-apache2"}},
	{Name: "etc-systemd-multi", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-systemd-multi"}},
	{Name: "etc-systemd-sockets", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-systemd-sockets"}},
	{Name: "etc-salt", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-salt"}},
	{Name: "etc-tomcat", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-tomcat"}},
	{Name: "etc-cobbler", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-cobbler"}},
	{Name: "etc-sysconfig", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-sysconfig"}},
	{Name: "etc-postfix", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-postfix"}},
	{Name: "etc-tls", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-tls"}},
}

// EtcVolumeMounts represents volumes mounted in /etc folder.
var EtcVolumeMounts = []types.VolumeMount{
	{MountPath: "/etc/apache2", Name: "etc-apache2"},
	{MountPath: "/etc/systemd/system/multi-user.target.wants", Name: "etc-systemd-multi"},
	{MountPath: "/etc/systemd/system/sockets.target.wants", Name: "etc-systemd-sockets"},
	{MountPath: "/etc/salt", Name: "etc-salt"},
	{MountPath: "/etc/rhn", Name: "etc-rhn"},
	{MountPath: "/etc/tomcat", Name: "etc-tomcat"},
	{MountPath: "/etc/cobbler", Name: "etc-cobbler"},
	{MountPath: "/etc/sysconfig", Name: "etc-sysconfig"},
	{MountPath: "/etc/postfix", Name: "etc-postfix"},
}

// EtcServerVolumeMounts represents volumes used for configuration.
var EtcVolumes = []types.Volume{
	{Name: "etc-apache2", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-apache2"}},
	{Name: "etc-systemd-multi", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-systemd-multi"}},
	{Name: "etc-systemd-sockets", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-systemd-sockets"}},
	{Name: "etc-salt", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-salt"}},
	{Name: "etc-tomcat", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-tomcat"}},
	{Name: "etc-cobbler", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-cobbler"}},
	{Name: "etc-sysconfig", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-sysconfig"}},
	{Name: "etc-postfix", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-postfix"}},
	{Name: "etc-rhn", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-rhn"}},
}

// ServerVolumeMounts should match the volumes mapping from the container definition in both
// the helm chart and the systemctl services definitions.
var ServerVolumeMounts = append(append(append([]types.VolumeMount{
	{MountPath: "/var/lib/cobbler", Name: "var-cobbler"},
	{MountPath: "/var/lib/salt", Name: "var-salt"},
	{MountPath: "/var/cache", Name: "var-cache"},
	{MountPath: "/var/spacewalk", Name: "var-spacewalk"},
	{MountPath: "/var/log", Name: "var-log"},
	{MountPath: "/srv/salt", Name: "srv-salt"},
	{MountPath: "/srv/www/", Name: "srv-www"},
	{MountPath: "/srv/tftpboot", Name: "srv-tftpboot"},
	{MountPath: "/srv/formula_metadata", Name: "srv-formulametadata"},
	{MountPath: "/srv/pillar", Name: "srv-pillar"},
	{MountPath: "/srv/susemanager", Name: "srv-susemanager"},
	{MountPath: "/srv/spacewalk", Name: "srv-spacewalk"},
	{MountPath: "/root", Name: "root"},
}, EtcVolumeMounts[:]...), PgsqlRequiredVolumeMounts[:]...), CACertsVolumeMounts[:]...) //FIXME golang 1.22 will introduce slices.Concat.

// ServerVolumes match the volume with Persistent Volume Claim.
var ServerVolumes = append(append(append([]types.Volume{
	{Name: "var-cobbler", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "var-cobbler"}},
	{Name: "var-salt", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "var-salt"}},
	{Name: "var-cache", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "var-cache"}},
	{Name: "var-spacewalk", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "var-spacewalk"}},
	{Name: "var-log", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "var-log"}},
	{Name: "srv-salt", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "srv-salt"}},
	{Name: "srv-www", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "srv-www"}},
	{Name: "srv-tftpboot", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "srv-tftpboot"}},
	{Name: "srv-formulametadata", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "srv-formulametadata"}},
	{Name: "srv-pillar", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "srv-pillar"}},
	{Name: "srv-susemanager", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "srv-susemanager"}},
	{Name: "srv-spacewalk", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "srv-spacewalk"}},
	{Name: "root", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "root"}},
}, EtcVolumes[:]...), PgsqlRequiredVolumes[:]...), CACertsVolumes[:]...) //FIXME golang 1.22 will introduce slices.Concat.

// PROXY_HTTPD_VOLUMES volumes used by HTTPD in proxy.
var PROXY_HTTPD_VOLUMES = map[string]string{
	"uyuni-proxy-rhn-cache": "/var/cache/rhn",
	"uyuni-proxy-tftpboot":  "/srv/tftpboot",
}

// PROXY_HTTPD_VOLUMES volumes used by Squid in  proxy.
var PROXY_SQUID_VOLUMES = map[string]string{
	"uyuni-proxy-squid-cache": "/var/cache/squid",
}

// PROXY_TFTPD_VOLUMES volumes used by TFTP in proxy.
var PROXY_TFTPD_VOLUMES = map[string]string{
	"uyuni-proxy-tftpboot": "/srv/tftpboot:ro",
}
