// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/types"
)

const podTemplate = `# uyuni-proxy-pod.service, generated by mgrpxy

[Unit]
Description=Podman uyuni-proxy-pod.service
Wants=network.target
After=network-online.target
Requires=uyuni-proxy-httpd.service
Requires=uyuni-proxy-salt-broker.service
Requires=uyuni-proxy-squid.service
Requires=uyuni-proxy-ssh.service
Requires=uyuni-proxy-tftpd.service
Before=uyuni-proxy-httpd.service
Before=uyuni-proxy-salt-broker.service
Before=uyuni-proxy-squid.service
Before=uyuni-proxy-ssh.service
Before=uyuni-proxy-tftpd.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
{{- if .HTTPProxyFile }}
EnvironmentFile={{ .HTTPProxyFile }}
{{- end }}
Restart=on-failure
ExecStartPre=/bin/rm -f %t/uyuni-proxy-pod.pid %t/uyuni-proxy-pod.pod-id

ExecStartPre=/bin/sh -c '/usr/bin/podman pod create --infra-conmon-pidfile %t/uyuni-proxy-pod.pid \
		--pod-id-file %t/uyuni-proxy-pod.pod-id --name uyuni-proxy-pod \
		--network {{ .Network }} \
        {{- range .Ports }}
        -p {{ .Exposed }}:{{ .Port }}{{ if .Protocol }}/{{ .Protocol }}{{ end }} \
        {{- if $.IPV6Enabled }}
	-p [::]:{{ .Exposed }}:{{ .Port }}{{if .Protocol}}/{{ .Protocol }}{{end}} \
        {{- end }}
        {{- end }}
		--replace ${PODMAN_EXTRA_ARGS}'

ExecStart=/usr/bin/podman pod start --pod-id-file %t/uyuni-proxy-pod.pod-id
ExecStop=/usr/bin/podman pod stop --ignore --pod-id-file %t/uyuni-proxy-pod.pod-id -t 10
ExecStopPost=/usr/bin/podman pod rm --ignore -f --pod-id-file %t/uyuni-proxy-pod.pod-id

PIDFile=%t/uyuni-proxy-pod.pid
TimeoutStopSec=60
Type=forking

[Install]
WantedBy=multi-user.target default.target
`

// PodTemplateData POD information to create systemd file.
type PodTemplateData struct {
	Ports         []types.PortMap
	HTTPProxyFile string
	Network       string
	IPV6Enabled   bool
}

// Render will create the systemd configuration file.
func (data PodTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("service").Parse(podTemplate))
	return t.Execute(wr, data)
}
