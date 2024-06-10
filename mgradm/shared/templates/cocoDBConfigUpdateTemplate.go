// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

const cocoDBConfigScriptTemplate = `#!/bin/bash
echo "host {{ .DBName }} {{ .User }} {{ .Ip }} scram-sha-256" >> /var/lib/pgsql/data/pg_hba.conf
`

// CocoDBConfigTemplateData.
type CocoDBConfigTemplateData struct {
	Ip     string
	User   string
	DBName string
}

// Render will create script for updating db access permissions.
func (data CocoDBConfigTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(cocoDBConfigScriptTemplate))
	return t.Execute(wr, data)
}
