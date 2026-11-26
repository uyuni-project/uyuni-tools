// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package ptf

import (
	"testing"

	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestCheckParameters(t *testing.T) {
	createServiceImages := func(
		image string, cocoImage string, hubImage string, salineImage string, dbImage string, tftpImage string,
	) map[string]string {
		return map[string]string{
			podman.ServerService:                  image,
			podman.ServerAttestationService + "@": cocoImage,
			podman.HubXmlrpcService + "@":         hubImage,
			podman.SalineService + "@":            salineImage,
			podman.DBService:                      dbImage,
			podman.TFTPService:                    tftpImage,
		}
	}
	type testData struct {
		name                string
		serviceImages       map[string]string
		hasRemoteImages     map[string]bool
		expectedImage       string
		expectedCocoImage   string
		expectedHubImage    string
		expectedSalineImage string
		expectedDBImage     string
		expectedTFTPImage   string
		expectedError       string
	}

	data := []testData{
		{
			name: "server PTF only",
			serviceImages: createServiceImages(
				"registry.suse.com/suse/manager/5.0/x86_64/server:5.0.0", "", "", "", "", "",
			),
			hasRemoteImages:     map[string]bool{},
			expectedImage:       "registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server:latest-ptf-5678",
			expectedCocoImage:   "",
			expectedHubImage:    "",
			expectedSalineImage: "",
			expectedDBImage:     "",
			expectedTFTPImage:   "",
			expectedError:       "",
		},
		{
			name: "PTF for all services",
			serviceImages: createServiceImages(
				"registry.suse.com/suse/manager/5.0/x86_64/server:5.0.0",
				"registry.suse.com/suse/manager/5.0/x86_64/server-attestation:5.0.0",
				"registry.suse.com/suse/manager/5.0/x86_64/server-hub-xmlrpc-api:5.0.0",
				"registry.suse.com/suse/manager/5.0/x86_64/server-saline:5.0.0",
				"registry.suse.com/suse/manager/5.0/x86_64/server-postgresql:5.0.0",
				"registry.suse.com/suse/manager/5.0/x86_64/proxy-tftpd:5.0.0",
			),
			hasRemoteImages: map[string]bool{
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-attestation:latest-ptf-5678":    true,
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-hub-xmlrpc-api:latest-ptf-5678": true,
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-saline:latest-ptf-5678":         true,
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-postgresql:latest-ptf-5678":     true,
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/proxy-tftpd:latest-ptf-5678":           true,
			},
			expectedImage:       "registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server:latest-ptf-5678",
			expectedCocoImage:   "registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-attestation:latest-ptf-5678",
			expectedHubImage:    "registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-hub-xmlrpc-api:latest-ptf-5678",
			expectedSalineImage: "registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-saline:latest-ptf-5678",
			expectedDBImage:     "registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-postgresql:latest-ptf-5678",
			expectedTFTPImage:   "registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/proxy-tftpd:latest-ptf-5678",
			expectedError:       "",
		},
		{
			name: "server only PTF and other services running",
			serviceImages: createServiceImages(
				"registry.suse.com/suse/manager/5.0/x86_64/server:5.0.0",
				"registry.suse.com/suse/manager/5.0/x86_64/server-attestation:5.0.0",
				"",
				"",
				"",
				"",
			),
			hasRemoteImages: map[string]bool{
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server:latest-ptf-5678":             true,
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-attestation:latest-ptf-5678": false,
			},
			expectedImage:       "registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server:latest-ptf-5678",
			expectedCocoImage:   "",
			expectedHubImage:    "",
			expectedSalineImage: "",
			expectedDBImage:     "",
			expectedTFTPImage:   "",
			expectedError:       "",
		},
		{
			name: "missing server service",
			serviceImages: createServiceImages(
				"",
				"",
				"registry.suse.com/suse/manager/5.0/x86_64/server-hub-xmlrpc-api:5.0.0",
				"",
				"",
				"",
			),
			hasRemoteImages: map[string]bool{
				"registry.suse.com/a/1234/5678/suse/manager/5.0/x86_64/server-hub-xmlrpc-api:latest-ptf-5678": true,
			},
			expectedImage:       "",
			expectedCocoImage:   "",
			expectedHubImage:    "",
			expectedSalineImage: "",
			expectedDBImage:     "",
			expectedTFTPImage:   "",
			expectedError:       "failed to find server image",
		},
	}

	for _, test := range data {
		t.Run(test.name, func(t *testing.T) {
			getServiceImage = func(service string) string {
				return test.serviceImages[service]
			}
			hasRemoteImage = func(image string, _ string) bool {
				return test.hasRemoteImages[image]
			}

			flags := podmanPTFFlags{
				PTFId:      "5678",
				CustomerID: "1234",
				ServerFlags: adm_utils.ServerFlags{
					Image: types.ImageFlags{
						Registry: types.Registry{
							Host: "registry.suse.com/suse/manager/5.0/x86_64/",
						},
					},
				},
			}
			actualError := flags.checkParameters("")
			errMessage := ""
			if actualError != nil {
				errMessage = actualError.Error()
			}
			testutils.AssertEquals(t, "error didn't match the expected behavior",
				test.expectedError, errMessage,
			)
			testutils.AssertEquals(t, "unexpected image", test.expectedImage, flags.Image.Name)
			testutils.AssertEquals(t, "unexpected coco image", test.expectedCocoImage, flags.Coco.Image.Name)
			testutils.AssertEquals(t, "unexpected hub image", test.expectedHubImage, flags.HubXmlrpc.Image.Name)
			testutils.AssertEquals(t, "unexpected saline image", test.expectedSalineImage, flags.Saline.Image.Name)
			testutils.AssertEquals(t, "unexpected db image", test.expectedDBImage, flags.Pgsql.Image.Name)
		})
	}
}
