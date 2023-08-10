package templates

import (
	"io"
	"text/template"
)

const MgrSetupScriptTemplate = `#!/bin/sh
{{- range $name, $value := .Env }}
export {{ $name }}={{ $value }}
{{- end }}

/usr/lib/susemanager/bin/mgr-setup -s -n

# clean before leaving
rm $0`

type MgrSetupScriptTemplateData struct {
	Env map[string]string
}

func (data MgrSetupScriptTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(MgrSetupScriptTemplate))
	return t.Execute(wr, data)
}
