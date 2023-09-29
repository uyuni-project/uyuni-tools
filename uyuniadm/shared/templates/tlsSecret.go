package templates

import (
	"io"
	"text/template"
)

// Deploy self-signed issuer or CA Certificate and key
const tlsSecretTemplate = `apiVersion: v1
kind: Secret
type: kubernetes.io/tls
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
data:
  ca.crt: {{ .RootCa }}
  tls.crt: {{ .Certificate }}
  tls.key: {{ .Key }}
`

type TlsSecretTemplateData struct {
	Name        string
	Namespace   string
	RootCa      string
	Certificate string
	Key         string
}

func (data TlsSecretTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("secret").Parse(tlsSecretTemplate))
	return t.Execute(wr, data)
}
