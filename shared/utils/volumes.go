// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import "github.com/uyuni-project/uyuni-tools/shared/types"

// PgsqlRequiredVolumeMounts represents volumes mount used by PostgreSQL.
var PgsqlRequiredVolumeMounts = []types.VolumeMount{
	{MountPath: "/etc/pki/tls", Name: "etc-tls", Size: "1Mi"},
	{MountPath: "/var/lib/pgsql", Name: "var-pgsql", Size: "50Gi"},
	{MountPath: "/etc/rhn", Name: "etc-rhn", Size: "1Mi"},
	{MountPath: "/etc/pki/spacewalk-tls", Name: "tls-key"},
}

// PgsqlRequiredVolumes represents volumes used by PostgreSQL.
var PgsqlRequiredVolumes = []types.Volume{
	{Name: "etc-tls", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-tls"}},
	{Name: "var-pgsql", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "var-pgsql"}},
	{Name: "etc-rhn", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-rhn"}},
	{Name: "tls-key",
		Secret: &types.Secret{
			SecretName: "uyuni-cert", Items: []types.SecretItem{
				{Key: "tls.crt", Path: "spacewalk.crt"},
				{Key: "tls.key", Path: "spacewalk.key"},
			},
		},
	},
}

// etcServerVolumeMounts represents volumes mounted in /etc folder.
var etcServerVolumeMounts = []types.VolumeMount{
	{MountPath: "/etc/apache2", Name: "etc-apache2", Size: "1Mi"},
	{MountPath: "/etc/systemd/system/multi-user.target.wants", Name: "etc-systemd-multi", Size: "1Mi"},
	{MountPath: "/etc/systemd/system/sockets.target.wants", Name: "etc-systemd-sockets", Size: "1Mi"},
	{MountPath: "/etc/salt", Name: "etc-salt", Size: "1Mi"},
	{MountPath: "/etc/tomcat", Name: "etc-tomcat", Size: "1Mi"},
	{MountPath: "/etc/cobbler", Name: "etc-cobbler", Size: "1Mi"},
	{MountPath: "/etc/sysconfig", Name: "etc-sysconfig", Size: "20Mi"},
	{MountPath: "/etc/postfix", Name: "etc-postfix", Size: "1Mi"},
	{MountPath: "/etc/sssd", Name: "etc-sssd", Size: "1Mi"},
}

// EtcServerVolumes represents volumes used for configuration.
var EtcServerVolumes = []types.Volume{
	{Name: "etc-apache2", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-apache2"}},
	{Name: "etc-systemd-multi", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-systemd-multi"}},
	{Name: "etc-systemd-sockets", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-systemd-sockets"}},
	{Name: "etc-salt", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-salt"}},
	{Name: "etc-tomcat", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-tomcat"}},
	{Name: "etc-cobbler", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-cobbler"}},
	{Name: "etc-sysconfig", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-sysconfig"}},
	{Name: "etc-postfix", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-postfix"}},
	{Name: "etc-rhn", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-rhn"}},
	{Name: "etc-sssd", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "etc-sssd"}},
}

var etcAndPgsqlVolumeMounts = append(PgsqlRequiredVolumeMounts, etcServerVolumeMounts[:]...)
var etcAndPgsqlVolumes = append(PgsqlRequiredVolumes, EtcServerVolumes[:]...)

// ServerVolumeMounts should match the volumes mapping from the container definition in both
// the helm chart and the systemctl services definitions.
var ServerVolumeMounts = append([]types.VolumeMount{
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
	{MountPath: "/root", Name: "root", Size: "1Mi"},
	{MountPath: "/etc/pki/trust/anchors/", Name: "ca-cert"},
	{MountPath: "/run/salt/master", Name: "run-salt-master"},
}, etcAndPgsqlVolumeMounts[:]...)

// ServerVolumes match the volume with Persistent Volume Claim.
var ServerVolumes = append([]types.Volume{
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
	{Name: "ca-cert", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "ca-cert"}},
	{Name: "run-salt-master", PersistentVolumeClaim: &types.PersistentVolumeClaim{ClaimName: "run-salt-master"}},
}, etcAndPgsqlVolumes[:]...)

// HubXmlrpcVolumeMounts represents volumes used by Hub Xmlrpc container.
var HubXmlrpcVolumeMounts = []types.VolumeMount{
	{MountPath: "/etc/pki/trust/anchors", Name: "ca-cert"},
}

// SalineVolumeMounts represents volumes used by Saline container.
var SalineVolumeMounts = []types.VolumeMount{
	{Name: "etc-salt", MountPath: "/etc/salt"},
	{Name: "etc-tls", MountPath: "/etc/pki/tls"},
	{Name: "run-salt-master", MountPath: "/run/salt/master"},
}

// ProxyHttpdVolumes volumes used by HTTPD in proxy.
var ProxyHttpdVolumes = []types.VolumeMount{
	{Name: "uyuni-proxy-rhn-cache", MountPath: "/var/cache/rhn:z"},
	{Name: "uyuni-proxy-tftpboot", MountPath: "/srv/tftpboot:z"},
}

// ProxySquidVolumes volumes used by Squid in  proxy.
var ProxySquidVolumes = []types.VolumeMount{
	{Name: "uyuni-proxy-squid-cache", MountPath: "/var/cache/squid:z"},
}

// ProxyTftpdVolumes used by TFTP in proxy.
var ProxyTftpdVolumes = []types.VolumeMount{
	{Name: "uyuni-proxy-tftpboot", MountPath: "/srv/tftpboot:ro,z"},
}
