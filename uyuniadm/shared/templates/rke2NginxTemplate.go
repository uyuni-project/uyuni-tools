// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/utils"
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

type Rke2NginxConfigTemplateData struct {
	Namespace string
	TcpPorts  []utils.PortMap
	UdpPorts  []utils.PortMap
}

func (data Rke2NginxConfigTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("rke2NginxConfig").Parse(rke2NginxConfigTemplate))
	return t.Execute(wr, data)
}
