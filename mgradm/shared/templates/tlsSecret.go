// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
)

// Deploy self-signed issuer or CA Certificate and key.
const tlsSecretTemplate = `apiVersion: v1
kind: Secret
type: kubernetes.io/tls
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
  labels:
    app: ` + kubernetes.ServerApp + `
data:
  ca.crt: {{ .RootCa }}
  tls.crt: {{ .Certificate }}
  tls.key: {{ .Key }}
`

// TLSSecretTemplateData contains information to create secret configuration file.
type TLSSecretTemplateData struct {
	Name        string
	Namespace   string
	RootCa      string
	Certificate string
	Key         string
}

// Render creates secret configuration file.
func (data TLSSecretTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("secret").Parse(tlsSecretTemplate))
	return t.Execute(wr, data)
}
