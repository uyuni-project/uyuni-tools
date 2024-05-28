// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/types"
)

const hubXmlrpcServiceTemplate = `# uyuni-uyuni-hub-xmlrpc.service, generated by mgradm
# Use an uyuni-hub-xmlrpc.service.d/local.conf file to override

[Unit]
Description=Uyuni Hub XMLRPC API container service
Wants=network.target
After=network-online.target

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Environment=HUB_API_URL=http://{{ .NamePrefix }}-server.mgr.internal:80/rpc/api
Environment=HUB_CONNECT_USING_SSL=true
Restart=on-failure
ExecStartPre=/bin/rm -f %t/uyuni-hub-xmlrpc-%i.pid %t/%n.ctr-id
ExecStartPre=/usr/bin/podman rm --ignore --force -t 10 {{ .NamePrefix }}-hub-xmlrpc-%i
ExecStart=/usr/bin/podman run \
	--conmon-pidfile %t/uyuni-hub-xmlrpc-%i.pid \
	--cidfile=%t/%n-%i.ctr-id \
	--cgroups=no-conmon \
	--sdnotify=conmon \
	-d \
	--replace \
	{{- range .Ports }}
        -p {{ .Exposed }}:{{ .Port }}{{if .Protocol}}/{{ .Protocol }}{{end}} \
        {{- end }}
        {{- range .Volumes }}
        -v {{ .Name }}:{{ .MountPath }} \
        {{- end }}
	-e HUB_API_URL \
	-e HUB_CONNECT_TIMEOUT \
	-e HUB_REQUEST_TIMEOUT \
	-e HUB_CONNECT_USING_SSL \
	--name {{ .NamePrefix }}-hub-xmlrpc-%i \
	--hostname {{ .NamePrefix }}-hub-xmlrpc-%i.mgr.internal \
	--network {{ .Network }} \
	${UYUNI_IMAGE}

ExecStop=/usr/bin/podman stop --ignore -t 10 --cidfile=%t/%n-%i.ctr-id
ExecStopPost=/usr/bin/podman rm -f --ignore -t 10 --cidfile=%t/%n-%i.ctr-id
PIDFile=%t/uyuni-hub-xmlrpc-%i.pid
TimeoutStopSec=60
TimeoutStartSec=60
Type=forking

[Install]
WantedBy=multi-user.target default.target
`

// PodmanServiceTemplateData POD information to create systemd file.
type HubXmlrpcServiceTemplateData struct {
	Volumes    []types.VolumeMount
	Ports      []types.PortMap
	NamePrefix string
	Image      string
	Network    string
}

// Render will create the systemd configuration file.
func (data HubXmlrpcServiceTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("service").Parse(hubXmlrpcServiceTemplate))
	return t.Execute(wr, data)
}
