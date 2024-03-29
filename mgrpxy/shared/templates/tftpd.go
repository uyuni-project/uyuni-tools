// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

const tftpdTemplate = `# uyuni-proxy-tftpd.service, generated by mgrpxy
# Use an uyuni-proxy-tftpd.service.d/local.conf file to override

[Unit]
Description=Uyuni proxy tftpd container service
Wants=network.target
After=network-online.target
BindsTo=uyuni-proxy-pod.service
After=uyuni-proxy-pod.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Environment=UYUNI_IMAGE={{ .Image }}
{{- if .HttpProxyFile }}
EnvironmentFile={{ .HttpProxyFile }}
{{- end }}
Restart=on-failure
ExecStartPre=/bin/rm -f %t/uyuni-proxy-tftpd.pid %t/uyuni-proxy-tftpd.ctr-id

ExecStart=/usr/bin/podman run \
	--conmon-pidfile %t/uyuni-proxy-tftpd.pid \
	--cidfile %t/uyuni-proxy-tftpd.ctr-id \
	--cgroups=no-conmon \
	--pod-id-file %t/uyuni-proxy-pod.pod-id -d \
	--replace -dt \
	-v /etc/uyuni/proxy:/etc/uyuni:ro \
	{{- range $name, $path := .Volumes }}
	 -v {{ $name }}:{{ $path }} \
	{{- end }}
	--name uyuni-proxy-tftpd \
	${UYUNI_IMAGE}

ExecStop=/usr/bin/podman stop --ignore --cidfile %t/uyuni-proxy-tftpd.ctr-id -t 10
ExecStopPost=/usr/bin/podman rm --ignore -f --cidfile %t/uyuni-proxy-tftpd.ctr-id
PIDFile=%t/uyuni-proxy-tftpd.pid
TimeoutStopSec=60
Type=forking

[Install]
WantedBy=multi-user.target default.target
`

// TFTPDTemplateData represents information used to create TFTPD systemd configuration file.
type TFTPDTemplateData struct {
	Volumes       map[string]string
	HttpProxyFile string
	Image         string
}

// Render will create the TFTPD systemd configuration file.
func (data TFTPDTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("service").Parse(tftpdTemplate))
	return t.Execute(wr, data)
}
