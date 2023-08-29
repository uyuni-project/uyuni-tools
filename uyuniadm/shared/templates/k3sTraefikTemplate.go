package templates

import (
	"io"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const k3sTraefikConfigTemplate = `apiVersion: helm.cattle.io/v1
kind: HelmChartConfig
metadata:
  name: traefik
  namespace: kube-system
spec:
  valuesContent: |-
    ports:
{{- range .TcpPorts }}
      {{ .Name }}:
        port: {{ .Port }}
        expose: true
        exposedPort: {{ .Exposed }}
        protocol: TCP
{{- end }}
{{- range .UdpPorts }}
      {{ .Name }}:
        port: {{ .Port }}
        expose: true
        exposedPort: {{ .Exposed }}
        protocol: UDP
{{- end }}
`

type K3sTraefikConfigTemplateData struct {
	TcpPorts []utils.PortMap
	UdpPorts []utils.PortMap
}

func (data K3sTraefikConfigTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("k3sTraefikConfig").Parse(k3sTraefikConfigTemplate))
	return t.Execute(wr, data)
}
