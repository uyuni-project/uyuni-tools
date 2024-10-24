// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"io"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/types"
)

const k3sTraefikConfigTemplate = `apiVersion: helm.cattle.io/v1
kind: HelmChartConfig
metadata:
  name: traefik
  namespace: kube-system
spec:
  valuesContent: |-
    ports:
{{- range .TCPPorts }}
      {{ .Name }}:
        port: {{ .Port }}
        expose: true
        exposedPort: {{ .Exposed }}
        protocol: TCP
{{- end }}
{{- range .UDPPorts }}
      {{ .Name }}:
        port: {{ .Port }}
        expose: true
        exposedPort: {{ .Exposed }}
        protocol: UDP
{{- end }}
`

// K3sTraefikConfigTemplateData represents information used to create K3s Traefik helm chart.
type K3sTraefikConfigTemplateData struct {
	TCPPorts []types.PortMap
	UDPPorts []types.PortMap
}

// Render will create the helm chart configuation for K3sTraefik.
func (data K3sTraefikConfigTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("k3sTraefikConfig").Parse(k3sTraefikConfigTemplate))
	return t.Execute(wr, data)
}
