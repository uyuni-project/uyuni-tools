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

// ServerKubernetesFlagsTestArgs are the expected values for AssertServerKubernetesFlags.
var ServerKubernetesFlagsTestArgs = []string{
	"--kubernetes-uyuni-namespace", "uyunins",
	"--kubernetes-certmanager-namespace", "certmanagerns",
	"--kubernetes-certmanager-chart", "oci://srv/certmanager",
	"--kubernetes-certmanager-version", "4.5.6",
	"--kubernetes-certmanager-values", "certmanager/values.yaml",
}

// AssertServerKubernetesFlags checks that all Kubernetes flags are parsed correctly.
func AssertServerKubernetesFlags(t *testing.T, flags *utils.KubernetesFlags) {
	testutils.AssertEquals(t, "Error parsing --helm-uyuni-namespace", "uyunins", flags.Uyuni.Namespace)
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

// VolumesFlagsTestExpected is the expected values for AssertVolumesFlags.
var VolumesFlagsTestExpected = []string{
	"--volumes-class", "MyStorageClass",
	"--volumes-mirror", "mirror-pv",
	"--volumes-database-size", "123Gi",
	"--volumes-database-class", "dbclass",
	"--volumes-packages-size", "456Gi",
	"--volumes-packages-class", "pkgclass",
	"--volumes-www-size", "123Mi",
	"--volumes-www-class", "wwwclass",
	"--volumes-cache-size", "789Gi",
	"--volumes-cache-class", "cacheclass",
}

// AssertVolumesFlags checks that all the volumes flags are parsed correctly.
func AssertVolumesFlags(t *testing.T, flags *utils.VolumesFlags) {
	testutils.AssertEquals(t, "Error parsing --volumes-class", "MyStorageClass", flags.Class)
	testutils.AssertEquals(t, "Error parsing --volumes-mirror", "mirror-pv", flags.Mirror)
	testutils.AssertEquals(t, "Error parsing --volumes-database-size", "123Gi", flags.Database.Size)
	testutils.AssertEquals(t, "Error parsing --volumes-database-class", "dbclass", flags.Database.Class)
	testutils.AssertEquals(t, "Error parsing --volumes-packages-size", "456Gi", flags.Packages.Size)
	testutils.AssertEquals(t, "Error parsing --volumes-packages-class", "pkgclass", flags.Packages.Class)
	testutils.AssertEquals(t, "Error parsing --volumes-www-size", "123Mi", flags.Www.Size)
	testutils.AssertEquals(t, "Error parsing --volumes-www-class", "wwwclass", flags.Www.Class)
	testutils.AssertEquals(t, "Error parsing --volumes-cache-size", "789Gi", flags.Cache.Size)
	testutils.AssertEquals(t, "Error parsing --volumes-cache-class", "cacheclass", flags.Cache.Class)
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

// DBUpdateImageFlagTestArgs is the expected values for AssertDBUpgradeImageFlag.
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

// SalineoFlagsTestArgs is the expected values for AssertSalineFlag.
var SalineFlagsTestArgs = []string{
	"--saline-image", "salineimg",
	"--saline-tag", "salinetag",
	"--saline-replicas", "1",
	"--saline-port", "8226",
}

// AssertSalineFlag asserts that all saline flags are parsed correctly.
func AssertSalineFlag(t *testing.T, flags *utils.SalineFlags) {
	testutils.AssertEquals(t, "Error parsing --saline-image", "salineimg", flags.Image.Name)
	testutils.AssertEquals(t, "Error parsing --saline-tag", "salinetag", flags.Image.Tag)
	testutils.AssertEquals(t, "Error parsing --saline-replicas", 1, flags.Replicas)
	testutils.AssertEquals(t, "Error parsing --saline-port", 8226, flags.Port)
	testutils.AssertTrue(t, "Saline should be changed", flags.IsChanged)
}
