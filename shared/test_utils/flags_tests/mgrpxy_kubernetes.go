// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flags_tests

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
)

// ProxyHelmFlagsTestArgs is the slice of parameters to use with AssertHelmFlags.
var ProxyHelmFlagsTestArgs = []string{
	"--helm-proxy-namespace", "uyunins",
	"--helm-proxy-chart", "oci://srv/proxy-helm",
	"--helm-proxy-version", "v1.2.3",
	"--helm-proxy-values", "path/value.yaml",
}

// AssertProxyHelmFlags checks that the proxy helm flags are parsed correctly.
func AssertProxyHelmFlags(t *testing.T, cmd *cobra.Command, flags *kubernetes.HelmFlags) {
	test_utils.AssertEquals(t, "Error parsing --helm-proxy-namespace", "uyunins", flags.Proxy.Namespace)
	test_utils.AssertEquals(t, "Error parsing --helm-proxy-chart", "oci://srv/proxy-helm", flags.Proxy.Chart)
	test_utils.AssertEquals(t, "Error parsing --helm-proxy-version", "v1.2.3", flags.Proxy.Version)
	test_utils.AssertEquals(t, "Error parsing --helm-proxy-values", "path/value.yaml", flags.Proxy.Values)
}

// ImageProxyFlagsTestArgs is the slice of parameters to use with AssertImageFlags.
var ImageProxyFlagsTestArgs = []string{
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
}

// AssertProxyImageFlags checks that all image flags are parsed correctly.
func AssertProxyImageFlags(t *testing.T, cmd *cobra.Command, flags *utils.ProxyImageFlags) {
	test_utils.AssertEquals(t, "Error parsing --tag", "v1.2.3", flags.Tag)
	test_utils.AssertEquals(t, "Error parsing --pullPolicy", "never", flags.PullPolicy)
	test_utils.AssertEquals(t, "Error parsing --httpd-image", "path/to/httpd", flags.Httpd.Name)
	test_utils.AssertEquals(t, "Error parsing --httpd-tag", "httpd-tag", flags.Httpd.Tag)
	test_utils.AssertEquals(t, "Error parsing --saltbroker-image", "path/to/saltbroker", flags.SaltBroker.Name)
	test_utils.AssertEquals(t, "Error parsing --saltbroker-tag", "saltbroker-tag", flags.SaltBroker.Tag)
	test_utils.AssertEquals(t, "Error parsing --squid-image", "path/to/squid", flags.Squid.Name)
	test_utils.AssertEquals(t, "Error parsing --squid-tag", "squid-tag", flags.Squid.Tag)
	test_utils.AssertEquals(t, "Error parsing --ssh-image", "path/to/ssh", flags.Ssh.Name)
	test_utils.AssertEquals(t, "Error parsing --ssh-tag", "ssh-tag", flags.Ssh.Tag)
	test_utils.AssertEquals(t, "Error parsing --tftpd-image", "path/to/tftpd", flags.Tftpd.Name)
	test_utils.AssertEquals(t, "Error parsing --tftpd-tag", "tftpd-tag", flags.Tftpd.Tag)
	test_utils.AssertEquals(t, "Error parsing --tuning-httpd", "path/to/httpd.conf", flags.Tuning.Httpd)
	test_utils.AssertEquals(t, "Error parsing --tuning-squid", "path/to/squid.conf", flags.Tuning.Squid)
}
