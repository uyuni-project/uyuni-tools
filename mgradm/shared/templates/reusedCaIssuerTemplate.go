// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
)

const uyuniCAIssuer = `apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: ` + kubernetes.CAIssuerName + `
  namespace: {{ .IssuerTemplate.Namespace }}
  labels:
    app: ` + kubernetes.ServerApp + `
spec:
  ca:
    secretName: ` + kubernetes.CASecretName + `
---
`

const reusedCAIssuerTemplate = `apiVersion: v1
kind: Secret
type: kubernetes.io/tls
metadata:
  name: ` + kubernetes.CASecretName + `
  namespace: {{ .IssuerTemplate.Namespace }}
  labels:
    app: ` + kubernetes.ServerApp + `
data:
  ca.crt: {{ .Certificate }}
  tls.crt: {{ .Certificate }}
  tls.key: {{ .Key }}
---
`

// NewReusedCAIssuerTemplate creates a new ReusedCAIssuerTemplate object.
func NewReusedCAIssuerTemplate(
	namespace string,
	fqdn string,
	certificate string,
	key string,
) ReusedCAIssuerTemplate {
	template := ReusedCAIssuerTemplate{
		IssuerTemplate: IssuerTemplate{
			Namespace: namespace,
			FQDN:      fqdn,
		},
		Certificate: certificate,
		Key:         key,
	}
	template.template = template
	return template
}

// ReusedCAIssuerTemplate is a template to render cert-manager issuer from an existing root CA.
type ReusedCAIssuerTemplate struct {
	IssuerTemplate
	Certificate string
	Key         string
	Template    string
}

// Render writers the issuer text in the wr parameter.
func (data ReusedCAIssuerTemplate) Render(wr io.Writer) error {
	t := template.Must(template.New("issuer").Parse(reusedCAIssuerTemplate + uyuniCAIssuer))
	return t.Execute(wr, data)
}
