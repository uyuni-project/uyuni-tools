// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flags_tests

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// Expected values for AssertHelmInstallFlags.
var ServerHelmFlagsTestArgs = []string{
	"--helm-uyuni-namespace", "uyunins",
	"--helm-uyuni-chart", "oci://srv/uyuni",
	"--helm-uyuni-version", "1.2.3",
	"--helm-uyuni-values", "uyuni/values.yaml",
	"--helm-certmanager-namespace", "certmanagerns",
	"--helm-certmanager-chart", "oci://srv/certmanager",
	"--helm-certmanager-version", "4.5.6",
	"--helm-certmanager-values", "certmanager/values.yaml",
}

// Assert that all Helm flags are parsed correctly.
func AssertServerHelmFlags(t *testing.T, cmd *cobra.Command, flags *utils.HelmFlags) {
	test_utils.AssertEquals(t, "Error parsing --helm-uyuni-namespace", "uyunins", flags.Uyuni.Namespace)
	test_utils.AssertEquals(t, "Error parsing --helm-uyuni-chart", "oci://srv/uyuni", flags.Uyuni.Chart)
	test_utils.AssertEquals(t, "Error parsing --helm-uyuni-version", "1.2.3", flags.Uyuni.Version)
	test_utils.AssertEquals(t, "Error parsing --helm-uyuni-values", "uyuni/values.yaml", flags.Uyuni.Values)
	test_utils.AssertEquals(t, "Error parsing --helm-certmanager-namespace",
		"certmanagerns", flags.CertManager.Namespace,
	)
	test_utils.AssertEquals(t, "Error parsing --helm-certmanager-chart",
		"oci://srv/certmanager", flags.CertManager.Chart,
	)
	test_utils.AssertEquals(t, "Error parsing --helm-certmanager-version", "4.5.6", flags.CertManager.Version)
	test_utils.AssertEquals(t, "Error parsing --helm-certmanager-values",
		"certmanager/values.yaml", flags.CertManager.Values,
	)
}

// Expected values for AssertImageFlag.
var ImageFlagsTestArgs = []string{
	"--image", "path/to/image",
	"--registry", "myregistry",
	"--tag", "v1.2.3",
	"--pullPolicy", "never",
}

// Assert that all image flags are parsed correctly.
func AssertImageFlag(t *testing.T, cmd *cobra.Command, flags *types.ImageFlags) {
	test_utils.AssertEquals(t, "Error parsing --image", "path/to/image", flags.Name)
	test_utils.AssertEquals(t, "Error parsing --registry", "myregistry", flags.Registry)
	test_utils.AssertEquals(t, "Error parsing --tag", "v1.2.3", flags.Tag)
	test_utils.AssertEquals(t, "Error parsing --pullPolicy", "never", flags.PullPolicy)
}

// Expected values for AssertDbUpgradeImageFlag.
var DbUpdateImageFlagTestArgs = []string{
	"--dbupgrade-image", "dbupgradeimg",
	"--dbupgrade-tag", "dbupgradetag",
}

// Assert that all DB upgrade image flags are parsed correctly.
func AssertDbUpgradeImageFlag(t *testing.T, cmd *cobra.Command, flags *types.ImageFlags) {
	test_utils.AssertEquals(t, "Error parsing --dbupgrade-image", "dbupgradeimg", flags.Name)
	test_utils.AssertEquals(t, "Error parsing --dbupgrade-tag", "dbupgradetag", flags.Tag)
}

// Expected values for AssertMirrorFlag.
var MirrorFlagTestArgs = []string{
	"--mirror", "/path/to/mirror",
}

// Assert that all mirror flags are parsed correctly.
func AssertMirrorFlag(t *testing.T, cmd *cobra.Command, value string) {
	test_utils.AssertEquals(t, "Error parsing --mirror", "/path/to/mirror", value)
}

// Expected values for AssertCocoFlag.
var CocoFlagsTestArgs = []string{
	"--coco-image", "cocoimg",
	"--coco-tag", "cocotag",
	"--coco-replicas", "2",
}

// Assert that all confidential computing flags are parsed correctly.
func AssertCocoFlag(t *testing.T, cmd *cobra.Command, flags *utils.CocoFlags) {
	test_utils.AssertEquals(t, "Error parsing --coco-image", "cocoimg", flags.Image.Name)
	test_utils.AssertEquals(t, "Error parsing --coco-tag", "cocotag", flags.Image.Tag)
	test_utils.AssertEquals(t, "Error parsing --coco-replicas", 2, flags.Replicas)
	test_utils.AssertTrue(t, "Coco should be changed", flags.IsChanged)
}

// Expected values for AssertHubXmlrpcFlag.
var HubXmlrpcFlagsTestArgs = []string{
	"--hubxmlrpc-image", "hubimg",
	"--hubxmlrpc-tag", "hubtag",
	"--hubxmlrpc-replicas", "1",
}

// Assert that all hub XML-RPC API flags are parsed correctly.
func AssertHubXmlrpcFlag(t *testing.T, cmd *cobra.Command, flags *utils.HubXmlrpcFlags) {
	test_utils.AssertEquals(t, "Error parsing --hubxmlrpc-image", "hubimg", flags.Image.Name)
	test_utils.AssertEquals(t, "Error parsing --hubxmlrpc-tag", "hubtag", flags.Image.Tag)
	test_utils.AssertEquals(t, "Error parsing --hubxmlrpc-replicas", 1, flags.Replicas)
	test_utils.AssertTrue(t, "Hub should be changed", flags.IsChanged)
}
