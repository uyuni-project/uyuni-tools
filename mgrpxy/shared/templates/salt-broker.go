// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

const saltBrokerTemplate = `# uyuni-proxy-salt-broker.service, generated by mgrpxy
# Use an uyuni-proxy-salt-broker.service.d/local.conf file to override

[Unit]
Description=Uyuni proxy Salt broker container service
Wants=network.target
After=network-online.target
BindsTo=uyuni-proxy-pod.service
After=uyuni-proxy-pod.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
{{- if .HttpProxyFile }}
EnvironmentFile={{ .HttpProxyFile }}
{{- end }}
Restart=on-failure
ExecStartPre=/bin/rm -f %t/uyuni-proxy-salt-broker.pid %t/uyuni-proxy-salt-broker.ctr-id

ExecStart=/bin/sh -c '/usr/bin/podman run \
	--conmon-pidfile %t/uyuni-proxy-salt-broker.pid \
	--cidfile %t/uyuni-proxy-salt-broker.ctr-id \
	--cgroups=no-conmon \
	--pod-id-file %t/uyuni-proxy-pod.pod-id -d \
	--replace -dt \
	-v /etc/uyuni/proxy:/etc/uyuni:ro \
	--name uyuni-proxy-salt-broker \
	${UYUNI_IMAGE}'

ExecStop=/usr/bin/podman stop --ignore --cidfile %t/uyuni-proxy-salt-broker.ctr-id -t 10
ExecStopPost=/usr/bin/podman rm --ignore -f --cidfile %t/uyuni-proxy-salt-broker.ctr-id
PIDFile=%t/uyuni-proxy-salt-broker.pid
TimeoutStopSec=60
Type=forking

[Install]
WantedBy=multi-user.target default.target
`

// SaltBrokerTemplateData represents Salt Broker information to create systemd file.
type SaltBrokerTemplateData struct {
	HttpProxyFile string
}

// Render will create the systemd configuration file.
func (data SaltBrokerTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("service").Parse(saltBrokerTemplate))
	return t.Execute(wr, data)
}
