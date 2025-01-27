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
if test -e /root/.MANAGER_SETUP_COMPLETE; then
	echo "Server appears to be already configured. Installation options may be ignored."
	exit 0
fi

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
    echo "starting tomcat..."
	(su -s /usr/bin/sh -g tomcat -G www -G susemanager tomcat /usr/lib/tomcat/server start)&

	echo "starting apache2..."
	/usr/sbin/start_apache2 -k start

	echo "starting taskomatic..."
	set -a
	. /usr/share/rhn/config-defaults/rhn_taskomatic_daemon.conf
	set +a
	/usr/sbin/taskomatic &

	echo "Creating first user..."
	{{ if .NoSSL }}
	CURL_SCHEME="http"
	{{ else }}
	CURL_SCHEME="-L -k https"
	{{ end }}

	curl -o /tmp/curl-retry -s --retry 7 $CURL_SCHEME://localhost/rhn/newlogin/CreateFirstUser.do

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
	else
		RESULT=1
	fi
fi

exit $RESULT
`

// MgrSetupScriptTemplateData represents information used to create setup script.
type MgrSetupScriptTemplateData struct {
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
