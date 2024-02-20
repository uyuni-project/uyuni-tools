// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"io"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/types"
)

const rke2NginxConfigTemplate = `apiVersion: helm.cattle.io/v1
kind: HelmChartConfig
metadata:
  name: rke2-ingress-nginx
  namespace: kube-system
spec:
  valuesContent: |-
    controller:
      config:
        hsts: "false"
    tcp:
{{- range .TcpPorts }}
      {{ .Exposed }}: "{{ $.Namespace }}/uyuni-tcp:{{ .Port }}"
{{- end }}
    udp:
{{- range .UdpPorts }}
      {{ .Exposed }}: "{{ $.Namespace }}/uyuni-udp:{{ .Port }}"
{{- end }}
`

// Rke2NginxConfigTemplateData represents information used to create Rke2 Ngix helm chart.
type Rke2NginxConfigTemplateData struct {
	Namespace string
	TcpPorts  []types.PortMap
	UdpPorts  []types.PortMap
}

// Render will create the helm chart configuation for Rke2 Nginx.
func (data Rke2NginxConfigTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("rke2NginxConfig").Parse(rke2NginxConfigTemplate))
	return t.Execute(wr, data)
}
