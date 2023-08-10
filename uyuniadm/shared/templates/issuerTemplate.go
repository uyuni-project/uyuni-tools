package templates

import (
	"io"
	"text/template"
)

// Deploy self-signed issuer
const issuerTemplate = `apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: uyuni-issuer
  namespace: default
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: uyuni-ca
  namespace: default
spec:
  isCA: true
  subject:
    countries: ["{{ .Country }}"]
    provinces: ["{{ .State }}"]
    localities: ["{{ .City }}"]
    organizations: ["{{ .Org }}"]
    organizationalUnits: ["{{ .OrgUnit }}"]
  emailAddresses:
    - {{ .Email }}
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
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: uyuni-ca-issuer
  namespace: default
spec:
  ca:
    secretName:
      uyuni-ca
`

type IssuerTemplateData struct {
	Country string
	State   string
	City    string
	Org     string
	OrgUnit string
	Email   string
	Fqdn    string
}

func (data IssuerTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("issuer").Parse(issuerTemplate))
	return t.Execute(wr, data)
}
