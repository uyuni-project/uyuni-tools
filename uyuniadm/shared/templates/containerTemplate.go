package templates

import (
	"io"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const containerTemplate = `
[Unit]
Description=Uyuni server image container service
Wants=network.target
After=network-online.target

[Service]
AddCapability=NET_RAW
ContainerName={{ .NamePrefix }}-server
Environment=TZ={{ .Timezone }}
Image={{ .Image }}
Tmpfs=/run
Network={{ .Network }}
{{- range .Ports }}
PublishPort={{ .Exposed }}:{{ .Port }}{{if .Protocol}}/{{ .Protocol }}{{end}} 
{{- end }}
{{- range $name, $path := .Volumes }}
Volume={{ $name }}:{{ $path }}
{{- end }}
PodmanArgs=--cgroups no-conmon --hostname=uyuni-server {{ .Args }}

[Service]
TimeoutStopSec=180
TimeoutStartSec=900
Type=forking

[Install]
WantedBy=multi-user.target default.target
`

type PodmanContainerTemplateData struct {
	Volumes    map[string]string
	NamePrefix string
	Args       string
	Ports      []utils.PortMap
	Timezone   string
	Image      string
	Network    string
}

func (data PodmanContainerTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("container").Parse(containerTemplate))
	return t.Execute(wr, data)
}
