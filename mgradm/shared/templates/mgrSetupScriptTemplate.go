// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

//nolint:lll
const mgrSetupScriptTemplate = `#!/bin/sh
{{- range $name, $value := .Env }}
export {{ $name }}='{{ $value }}'
{{- end }}

{{- if .DebugJava }}
echo 'JAVA_OPTS=" $JAVA_OPTS -Xdebug -Xrunjdwp:transport=dt_socket,address=*:8003,server=y,suspend=n" ' >> /etc/tomcat/conf.d/remote_debug.conf
echo 'JAVA_OPTS=" $JAVA_OPTS -Xdebug -Xrunjdwp:transport=dt_socket,address=*:8001,server=y,suspend=n" ' >> /etc/rhn/taskomatic.conf
echo 'JAVA_OPTS=" $JAVA_OPTS -Xdebug -Xrunjdwp:transport=dt_socket,address=*:8002,server=y,suspend=n" ' >> /usr/share/rhn/config-defaults/rhn_search_daemon.conf
{{- end }}

/usr/lib/susemanager/bin/mgr-setup -s -n
RESULT=$?

# The CA needs to be added to the database for Kickstart use.
/usr/bin/rhn-ssl-dbstore --ca-cert=/etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT

if test -n "{{ .AdminPassword }}"; then
	{{ if .NoSSL }}
	CURL_SCHEME="http"
	{{ else }}
	CURL_SCHEME="-k https"
	{{ end }}

	HTTP_CODE=$(curl -o /dev/null -s -w %{http_code} $CURL_SCHEME://localhost/rhn/newlogin/CreateFirstUser.do)
	if test "$HTTP_CODE" == "200"; then
		echo "Creating administration user"
		curl -s -o /tmp/curl_out \
			-d "orgName={{ .OrgName }}" \
			-d "adminLogin={{ .AdminLogin }}" \
			-d "adminPassword={{ .AdminPassword }}" \
			-d "firstName={{ .AdminFirstName }}" \
			-d "lastName={{ .AdminLastName }}" \
			-d "email={{ .AdminEmail }}" \
			$CURL_SCHEME://localhost/rhn/manager/api/org/createFirst
		if ! grep -q '^{"success":true' /tmp/curl_out ; then
			echo "Failed to create the administration user"
			cat /tmp/curl_out
		fi
		rm -f /tmp/curl_out
	elif test "$HTTP_CODE" == "403"; then
		echo "Administration user already exists, reusing"
	fi
fi

# clean before leaving
rm $0
exit $RESULT
`

// MgrSetupScriptTemplateData represents information used to create setup script.
type MgrSetupScriptTemplateData struct {
	Env            map[string]string
	NoSSL          bool
	DebugJava      bool
	AdminPassword  string
	AdminLogin     string
	AdminFirstName string
	AdminLastName  string
	AdminEmail     string
	OrgName        string
}

// Render will create setup script.
func (data MgrSetupScriptTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("script").Parse(mgrSetupScriptTemplate))
	return t.Execute(wr, data)
}
