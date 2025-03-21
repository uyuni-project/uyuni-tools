// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
)

// Deploy self-signed issuer or CA Certificate and key.
const generatedCAIssuerTemplate = `apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: uyuni-issuer
  namespace: {{ .IssuerTemplate.Namespace }}
  labels:
    app: ` + kubernetes.ServerApp + `
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: uyuni-ca
  namespace: {{ .IssuerTemplate.Namespace }}
  labels:
    app: ` + kubernetes.ServerApp + `
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
  commonName: {{ .IssuerTemplate.FQDN }}
  dnsNames:
    - {{ .IssuerTemplate.FQDN }}
  secretName: ` + kubernetes.CASecretName + `
  privateKey:
    algorithm: ECDSA
    size: 256
  issuerRef:
    name: uyuni-issuer
    kind: Issuer
    group: cert-manager.io
---
`

func NewGeneratedCAIssuerTemplate(
	namespace string,
	fqdn string,
	country string,
	state string,
	city string,
	org string,
	orgUnit string,
	email string,
) GeneratedCAIssuerTemplate {
	template := GeneratedCAIssuerTemplate{
		IssuerTemplate: IssuerTemplate{
			Namespace: namespace,
			FQDN:      fqdn,
		},
		Country: country,
		State:   state,
		City:    city,
		Org:     org,
		OrgUnit: orgUnit,
		Email:   email,
	}
	template.template = template
	return template
}

// GeneratedCAIssuerTemplate is a template to render cert-manager issuers for a generated self-signed CA.
type GeneratedCAIssuerTemplate struct {
	IssuerTemplate
	Country string
	State   string
	City    string
	Org     string
	OrgUnit string
	Email   string
}

// Render writers the issuer text in the wr parameter.
func (data GeneratedCAIssuerTemplate) Render(wr io.Writer) error {
	t := template.Must(template.New("issuer").Parse(generatedCAIssuerTemplate + uyuniCAIssuer))
	return t.Execute(wr, data)
}
