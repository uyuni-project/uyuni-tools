// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/types"
)

const serviceTemplate = `# uyuni-server.service, generated by mgradm
# Use an uyuni-server.service.d/local.conf file to override

[Unit]
Description=Uyuni server image container service
Wants=network.target
After=network-online.target
RequiresMountsFor=%t/containers

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Environment=TZ={{ .Timezone }}
Restart=on-failure
ExecStartPre=/bin/rm -f %t/uyuni-server.pid %t/%n.ctr-id
ExecStartPre=/usr/bin/podman rm --ignore --force -t 10 {{ .NamePrefix }}-server
ExecStart=/usr/bin/podman run \
	--conmon-pidfile %t/uyuni-server.pid \
	--cidfile=%t/%n.ctr-id \
	--cgroups=no-conmon \
	--sdnotify=conmon \
	-d \
	--name {{ .NamePrefix }}-server \
	--hostname {{ .NamePrefix }}-server \
	{{ .Args }} \
	{{- range .Ports }}
	-p {{ .Exposed }}:{{ .Port }}{{if .Protocol}}/{{ .Protocol }}{{end}} \
	{{- end }}
	{{- range $name, $path := .Volumes }}
	-v {{ $name }}:{{ $path }} \
	{{- end }}
	-e TZ=${TZ} \
	--network {{ .Network }} \
	${UYUNI_IMAGE}
ExecStop=/usr/bin/podman stop \
	--ignore -t 10 \
	--cidfile=%t/%n.ctr-id
ExecStopPost=/usr/bin/podman rm \
	-f \
	--ignore -t 10 \
	--cidfile=%t/%n.ctr-id

PIDFile=%t/uyuni-server.pid
TimeoutStopSec=180
TimeoutStartSec=900
Type=forking

[Install]
WantedBy=multi-user.target default.target
`

type PodmanServiceTemplateData struct {
	Volumes    map[string]string
	NamePrefix string
	Args       string
	Ports      []types.PortMap
	Timezone   string
	Image      string
	Network    string
}

func (data PodmanServiceTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("service").Parse(serviceTemplate))
	return t.Execute(wr, data)
}
