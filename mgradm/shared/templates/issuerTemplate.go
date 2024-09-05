// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

// Deploy self-signed issuer or CA Certificate and key.
const generatedCaIssuerTemplate = `apiVersion: cert-manager.io/v1
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
---
`

// GeneratedCaIssuerTemplateData is a template to render cert-manager issuers for a generated self-signed CA.
type GeneratedCaIssuerTemplateData struct {
	Namespace string
	Country   string
	State     string
	City      string
	Org       string
	OrgUnit   string
	Email     string
	Fqdn      string
}

// Render creates issuer file.
func (data GeneratedCaIssuerTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("issuer").Parse(generatedCaIssuerTemplate + uyuniCaIssuer))
	return t.Execute(wr, data)
}
