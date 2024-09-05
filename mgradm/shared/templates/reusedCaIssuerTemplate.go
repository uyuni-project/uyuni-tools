// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

const uyuniCaIssuer = `apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: uyuni-ca-issuer
  namespace: {{ .Namespace }}
spec:
  ca:
    secretName: uyuni-ca
`

const reusedCaIssuerTemplate = `apiVersion: v1
kind: Secret
type: kubernetes.io/tls
metadata:
  name: uyuni-ca
  namespace: {{ .Namespace }}
data:
  ca.crt: {{ .Certificate }}
  tls.crt: {{ .Certificate }}
  tls.key: {{ .Key }}
---
`

// ReusedCaIssuerTemplateData is a template to render cert-manager issuer from an existing root CA.
type ReusedCaIssuerTemplateData struct {
	Namespace   string
	Certificate string
	Key         string
}

// Render creates issuer file.
func (data ReusedCaIssuerTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("issuer").Parse(reusedCaIssuerTemplate + uyuniCaIssuer))
	return t.Execute(wr, data)
}
