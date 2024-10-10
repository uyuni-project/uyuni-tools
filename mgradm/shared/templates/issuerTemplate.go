// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

// Deploy self-signed issuer or CA Certificate and key.
const issuerTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .Namespace }}
  labels:
    name: {{ .Namespace }}
---
{{if and .Certificate .Key -}}
apiVersion: v1
kind: Secret
type: kubernetes.io/tls
metadata:
  name: uyuni-ca
  namespace: {{ .Namespace }}
data:
  ca.crt: {{ .RootCa }}
  tls.crt: {{ .Certificate }}
  tls.key: {{ .Key }}
{{- else }}
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: uyuni-issuer
  namespace: {{ .Namespace }}
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: uyuni-ca
  namespace: {{ .Namespace }}
spec:
  isCA: true
{{- if or .Country .State .City .Org .OrgUnit }}
  subject:
	{{- if .Country }}
    countries: ["{{ .Country }}"]
	{{- end }}
	{{- if .State }}
    provinces: ["{{ .State }}"]
	{{- end }}
	{{- if .City }}
    localities: ["{{ .City }}"]
	{{- end }}
	{{- if .Org }}
    organizations: ["{{ .Org }}"]
	{{- end }}
	{{- if .OrgUnit }}
    organizationalUnits: ["{{ .OrgUnit }}"]
	{{- end }}
{{- end }}
{{- if .Email }}
  emailAddresses:
    - {{ .Email }}
{{- end }}
  commonName: {{ .Fqdn }}
  dnsNames:
    - {{ .Fqdn }}
  secretName: uyuni-ca
  privateKey:
    algorithm: ECDSA
    size: 256
  issuerRef:
    name: uyuni-issuer
    kind: Issuer
    group: cert-manager.io
{{- end }}
---
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: uyuni-ca-issuer
  namespace: {{ .Namespace }}
spec:
  ca:
    secretName:
      uyuni-ca
`

// IssuerTemplateData represents information used to create issuer file.
type IssuerTemplateData struct {
	Namespace   string
	Country     string
	State       string
	City        string
	Org         string
	OrgUnit     string
	Email       string
	Fqdn        string
	RootCa      string
	Certificate string
	Key         string
}

// Render creates issuer file.
func (data IssuerTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("issuer").Parse(issuerTemplate))
	return t.Execute(wr, data)
}
