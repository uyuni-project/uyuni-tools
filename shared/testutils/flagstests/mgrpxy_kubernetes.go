// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flagstests

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

// ImageProxyFlagsTestArgs is the slice of parameters to use with AssertImageFlags.
var ImageProxyFlagsTestArgs = []string{
	"--registry", "myoldregistry.com",
	"--tag", "v1.2.3",
	"--pullPolicy", "never",
	"--httpd-image", "path/to/httpd",
	"--httpd-tag", "httpd-tag",
	"--saltbroker-image", "path/to/saltbroker",
	"--saltbroker-tag", "saltbroker-tag",
	"--squid-image", "path/to/squid",
	"--squid-tag", "squid-tag",
	"--ssh-image", "path/to/ssh",
	"--ssh-tag", "ssh-tag",
	"--tftpd-image", "path/to/tftpd",
	"--tftpd-tag", "tftpd-tag",
	"--tuning-httpd", "path/to/httpd.conf",
	"--tuning-squid", "path/to/squid.conf",
	"--tuning-ssh", "path/to/ssh.conf",
	"--registry-host", "myregistry.com",
	"--registry-user", "user",
	"--registry-password", "password",
}

// AssertProxyImageFlags checks that all image flags are parsed correctly.
func AssertProxyImageFlags(t *testing.T, flags *utils.ProxyImageFlags) {
	testutils.AssertEquals(t, "Error parsing --registry", "myoldregistry.com", flags.Registry.Host)
	testutils.AssertEquals(t, "Error parsing --registry-host", "myoldregistry.com", flags.Registry.Host)
	testutils.AssertEquals(t, "Error parsing --registry-user", "user", flags.Registry.User)
	testutils.AssertEquals(t, "Error parsing --registry-password", "password", flags.Registry.Password)
	testutils.AssertEquals(t, "Error parsing --tag", "v1.2.3", flags.Tag)
	testutils.AssertEquals(t, "Error parsing --pullPolicy", "never", flags.PullPolicy)
	testutils.AssertEquals(t, "Error parsing --httpd-image", "path/to/httpd", flags.Httpd.Name)
	testutils.AssertEquals(t, "Error parsing --httpd-tag", "httpd-tag", flags.Httpd.Tag)
	testutils.AssertEquals(t, "Error parsing --saltbroker-image", "path/to/saltbroker", flags.SaltBroker.Name)
	testutils.AssertEquals(t, "Error parsing --saltbroker-tag", "saltbroker-tag", flags.SaltBroker.Tag)
	testutils.AssertEquals(t, "Error parsing --squid-image", "path/to/squid", flags.Squid.Name)
	testutils.AssertEquals(t, "Error parsing --squid-tag", "squid-tag", flags.Squid.Tag)
	testutils.AssertEquals(t, "Error parsing --ssh-image", "path/to/ssh", flags.SSH.Name)
	testutils.AssertEquals(t, "Error parsing --ssh-tag", "ssh-tag", flags.SSH.Tag)
	testutils.AssertEquals(t, "Error parsing --tftpd-image", "path/to/tftpd", flags.Tftpd.Name)
	testutils.AssertEquals(t, "Error parsing --tftpd-tag", "tftpd-tag", flags.Tftpd.Tag)
	testutils.AssertEquals(t, "Error parsing --tuning-httpd", "path/to/httpd.conf", flags.Tuning.Httpd)
	testutils.AssertEquals(t, "Error parsing --tuning-squid", "path/to/squid.conf", flags.Tuning.Squid)
	testutils.AssertEquals(t, "Error parsing --tuning-ssh", "path/to/ssh.conf", flags.Tuning.SSH)
}
