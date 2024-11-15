// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flagstests

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// ServerHelmFlagsTestArgs is the expected values for AssertHelmInstallFlags.
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

// AssertServerHelmFlags checks that all Helm flags are parsed correctly.
func AssertServerHelmFlags(t *testing.T, flags *utils.HelmFlags) {
	testutils.AssertEquals(t, "Error parsing --helm-uyuni-namespace", "uyunins", flags.Uyuni.Namespace)
	testutils.AssertEquals(t, "Error parsing --helm-uyuni-chart", "oci://srv/uyuni", flags.Uyuni.Chart)
	testutils.AssertEquals(t, "Error parsing --helm-uyuni-version", "1.2.3", flags.Uyuni.Version)
	testutils.AssertEquals(t, "Error parsing --helm-uyuni-values", "uyuni/values.yaml", flags.Uyuni.Values)
	testutils.AssertEquals(t, "Error parsing --helm-certmanager-namespace",
		"certmanagerns", flags.CertManager.Namespace,
	)
	testutils.AssertEquals(t, "Error parsing --helm-certmanager-chart",
		"oci://srv/certmanager", flags.CertManager.Chart,
	)
	testutils.AssertEquals(t, "Error parsing --helm-certmanager-version", "4.5.6", flags.CertManager.Version)
	testutils.AssertEquals(t, "Error parsing --helm-certmanager-values",
		"certmanager/values.yaml", flags.CertManager.Values,
	)
}

// ImageFlagsTestArgs is the expected values for AssertImageFlag.
var ImageFlagsTestArgs = []string{
	"--image", "path/to/image",
	"--registry", "myregistry",
	"--tag", "v1.2.3",
	"--pullPolicy", "never",
}

// AssertImageFlag checks that all image flags are parsed correctly.
func AssertImageFlag(t *testing.T, flags *types.ImageFlags) {
	testutils.AssertEquals(t, "Error parsing --image", "path/to/image", flags.Name)
	testutils.AssertEquals(t, "Error parsing --registry", "myregistry", flags.Registry)
	testutils.AssertEquals(t, "Error parsing --tag", "v1.2.3", flags.Tag)
	testutils.AssertEquals(t, "Error parsing --pullPolicy", "never", flags.PullPolicy)
}

// DBUpdateImageFlagTestArgs is the expected values for AssertDbUpgradeImageFlag.
var DBUpdateImageFlagTestArgs = []string{
	"--dbupgrade-image", "dbupgradeimg",
	"--dbupgrade-tag", "dbupgradetag",
}

// AssertDBUpgradeImageFlag asserts that all DB upgrade image flags are parsed correctly.
func AssertDBUpgradeImageFlag(t *testing.T, flags *types.ImageFlags) {
	testutils.AssertEquals(t, "Error parsing --dbupgrade-image", "dbupgradeimg", flags.Name)
	testutils.AssertEquals(t, "Error parsing --dbupgrade-tag", "dbupgradetag", flags.Tag)
}

// MirrorFlagTestArgs is the expected values for AssertMirrorFlag.
var MirrorFlagTestArgs = []string{
	"--mirror", "/path/to/mirror",
}

// AssertMirrorFlag asserts that all mirror flags are parsed correctly.
func AssertMirrorFlag(t *testing.T, value string) {
	testutils.AssertEquals(t, "Error parsing --mirror", "/path/to/mirror", value)
}

// CocoFlagsTestArgs is the expected values for AssertCocoFlag.
var CocoFlagsTestArgs = []string{
	"--coco-image", "cocoimg",
	"--coco-tag", "cocotag",
	"--coco-replicas", "2",
}

// AssertCocoFlag asserts that all confidential computing flags are parsed correctly.
func AssertCocoFlag(t *testing.T, flags *utils.CocoFlags) {
	testutils.AssertEquals(t, "Error parsing --coco-image", "cocoimg", flags.Image.Name)
	testutils.AssertEquals(t, "Error parsing --coco-tag", "cocotag", flags.Image.Tag)
	testutils.AssertEquals(t, "Error parsing --coco-replicas", 2, flags.Replicas)
	testutils.AssertTrue(t, "Coco should be changed", flags.IsChanged)
}

// HubXmlrpcFlagsTestArgs is the expected values for AssertHubXmlrpcFlag.
var HubXmlrpcFlagsTestArgs = []string{
	"--hubxmlrpc-image", "hubimg",
	"--hubxmlrpc-tag", "hubtag",
	"--hubxmlrpc-replicas", "1",
}

// AssertHubXmlrpcFlag asserts that all hub XML-RPC API flags are parsed correctly.
func AssertHubXmlrpcFlag(t *testing.T, flags *utils.HubXmlrpcFlags) {
	testutils.AssertEquals(t, "Error parsing --hubxmlrpc-image", "hubimg", flags.Image.Name)
	testutils.AssertEquals(t, "Error parsing --hubxmlrpc-tag", "hubtag", flags.Image.Tag)
	testutils.AssertEquals(t, "Error parsing --hubxmlrpc-replicas", 1, flags.Replicas)
	testutils.AssertTrue(t, "Hub should be changed", flags.IsChanged)
}
