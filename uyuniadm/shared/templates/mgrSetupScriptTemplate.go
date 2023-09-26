package templates

import (
	"io"
	"text/template"
)

const MgrSetupScriptTemplate = `#!/bin/sh
{{- range $name, $value := .Env }}
export {{ $name }}={{ $value }}
{{- end }}

{{- if .DebugJava }}
echo 'JAVA_OPTS=" $JAVA_OPTS -Xdebug -Xrunjdwp:transport=dt_socket,address=*:8002,server=y,suspend=n" ' >> /etc/tomcat/conf.d/remote_debug.conf
echo 'JAVA_OPTS=" $JAVA_OPTS -Xdebug -Xrunjdwp:transport=dt_socket,address=*:8001,server=y,suspend=n" ' >> /etc/rhn/taskomatic.conf
{{- end }}

/usr/lib/susemanager/bin/mgr-setup -s -n

# clean before leaving
rm $0`

type MgrSetupScriptTemplateData struct {
	Env       map[string]string
	DebugJava bool
}

func (data MgrSetupScriptTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(MgrSetupScriptTemplate))
	return t.Execute(wr, data)
}
