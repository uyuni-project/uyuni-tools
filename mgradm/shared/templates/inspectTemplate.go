// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"

	"github.com/uyuni-project/uyuni-tools/shared/types"
)

const inspectTemplate = `#!/bin/bash
# inspect.sh, generated by mgradm
{{- range .Param }}
echo "{{ .Variable }}=$({{ .CLI }})" >> {{ $.OutputFile }}
{{- end }}
exit 0
`

type InspectTemplateData struct {
	Param      []types.InspectData
	OutputFile string
}

func (data InspectTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("inspect").Parse(inspectTemplate))
	return t.Execute(wr, data)
}
