// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package podman

import (
	"fmt"
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestCheckParameters(t *testing.T) {
	createServiceImages := func(
		image string, cocoImage string, hubImage string, salineImage string, dbImage string,
	) map[string]string {
		return map[string]string{
			podman.ServerService:                  image,
			podman.ServerAttestationService + "@": cocoImage,
			podman.HubXmlrpcService:               hubImage,
			podman.SalineService:                  salineImage,
			podman.DBService:                      dbImage,
		}
	}
	type testData struct {
		serviceImages       map[string]string
		hasRemoteImages     map[string]bool
		expectedImage       string
		expectedCocoImage   string
		expectedHubImage    string
		expectedSalineImage string
		expectedDBImage     string
		expectedError       string
	}

	data := []testData{
		{
			createServiceImages("registry.suse.com/suse/manager/5.0/x86_64/server:5.0.0", "", "", "", ""),
			map[string]bool{},
			"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server:latest-ptf-5678",
			"",
			"",
			"",
			"",
			"",
		},
		{
			createServiceImages(
				"registry.suse.com/suse/manager/5.0/x86_64/server:5.0.0",
				"registry.suse.com/suse/manager/5.0/x86_64/server-attestation:5.0.0",
				"registry.suse.com/suse/manager/5.0/x86_64/server-hub-xmlrpc-api:5.0.0",
				"registry.suse.com/suse/manager/5.0/x86_64/server-saline:5.0.0",
				"registry.suse.com/suse/manager/5.0/x86_64/server-postgresql:5.0.0",
			),
			map[string]bool{
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-attestation:latest-ptf-5678":    true,
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-hub-xmlrpc-api:latest-ptf-5678": true,
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-saline:latest-ptf-5678":         true,
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-postgresql:latest-ptf-5678":     true,
			},
			"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server:latest-ptf-5678",
			"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-attestation:latest-ptf-5678",
			"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-hub-xmlrpc-api:latest-ptf-5678",
			"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-saline:latest-ptf-5678",
			"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-postgresql:latest-ptf-5678",
			"",
		},
		{
			createServiceImages(
				"registry.suse.com/suse/manager/5.0/x86_64/server:5.0.0",
				"registry.suse.com/suse/manager/5.0/x86_64/server-attestation:5.0.0",
				"",
				"",
				"",
			),
			map[string]bool{
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server:latest-ptf-5678":             true,
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-attestation:latest-ptf-5678": false,
			},
			"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server:latest-ptf-5678",
			"",
			"",
			"",
			"",
			"",
		},
		{
			createServiceImages(
				"",
				"",
				"registry.suse.com/suse/manager/5.0/x86_64/server-hub-xmlrpc-api:5.0.0",
				"",
				"",
			),
			map[string]bool{
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-hub-xmlrpc-api:latest-ptf-5678": true,
			},
			"",
			"",
			"",
			"",
			"",
			"failed to find server image",
		},
	}

	for i, test := range data {
		getServiceImage = func(service string) string {
			return test.serviceImages[service]
		}
		hasRemoteImage = func(image string) bool {
			return test.hasRemoteImages[image]
		}

		flags := podmanPTFFlags{
			PTFId:      "5678",
			CustomerID: "1234",
			ServerFlags: utils.ServerFlags{
				Installation: utils.InstallationFlags{
					SCC: types.SCCCredentials{
						Registry: "registry.suse.com",
					},
				},
			},
		}

		testCase := fmt.Sprintf("case #%d - ", i+1)
		actualError := flags.checkParameters()
		errMessage := ""
		if actualError != nil {
			errMessage = actualError.Error()
		}
		testutils.AssertEquals(t, testCase+"error didn't match the expected behavior",
			test.expectedError, errMessage,
		)
		testutils.AssertEquals(t, testCase+"unexpected image", test.expectedImage, flags.Image.Name)
		testutils.AssertEquals(t, testCase+"unexpected coco image", test.expectedCocoImage, flags.Coco.Image.Name)
		testutils.AssertEquals(t, testCase+"unexpected hub image", test.expectedHubImage, flags.HubXmlrpc.Image.Name)
		testutils.AssertEquals(t, testCase+"unexpected saline image", test.expectedSalineImage, flags.Saline.Image.Name)
		testutils.AssertEquals(t, testCase+"unexpected db image", test.expectedDBImage, flags.Pgsql.Image.Name)
	}
}
