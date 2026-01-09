// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package podman

import (
	"fmt"
	"strings"
	"testing"

	mgrpxy_podman "github.com/uyuni-project/uyuni-tools/mgrpxy/shared/podman"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestUpdateParametersValidation(t *testing.T) {
	// Save original functions and mock them to bypass the image processing logic
	// This ensures the unit test focuses only on the validation at the start of updateParameters.
	originalGetServiceImage := getServiceImage
	originalHasRemoteImage := hasRemoteImage

	// Mock getServiceImage to return empty string to skip image processing logic
	getServiceImage = func(_ string) string { return "" }
	hasRemoteImage = func(_ string, _ string) bool { return false }

	defer func() {
		// Restore original functions after test completion
		getServiceImage = originalGetServiceImage
		hasRemoteImage = originalHasRemoteImage
	}()

	tests := []struct {
		name         string
		ptfID        string
		testID       string
		customerID   string
		expectError  bool
		errorMessage string
	}{
		{
			name:        "Success_PTF",
			ptfID:       "ptf123",
			testID:      "",
			customerID:  "sccuser",
			expectError: false,
		},
		{
			name:        "Success_Test",
			ptfID:       "",
			testID:      "test123",
			customerID:  "sccuser",
			expectError: false,
		},
		{
			name:         "Error_BothSet",
			ptfID:        "ptf123",
			testID:       "test123",
			customerID:   "sccuser",
			expectError:  true,
			errorMessage: "ptf and test flags cannot be set simultaneously",
		},
		{
			name:         "Error_BothEmpty",
			ptfID:        "",
			testID:       "",
			customerID:   "sccuser",
			expectError:  true,
			errorMessage: "ptf and test flags cannot be empty simultaneously",
		},
		{
			name:         "Error_CustomerEmpty_WithPTF",
			ptfID:        "ptf123",
			testID:       "",
			customerID:   "",
			expectError:  true,
			errorMessage: "user flag cannot be empty",
		},
		{
			name:         "Error_CustomerEmpty_WithTest",
			ptfID:        "",
			testID:       "test123",
			customerID:   "",
			expectError:  true,
			errorMessage: "user flag cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := podmanPTFFlags{
				UpgradeFlags: mgrpxy_podman.PodmanProxyFlags{
					ProxyImageFlags: utils.ProxyImageFlags{
						Registry: types.Registry{
							Host: "registry.suse.com",
						},
					},
				},
				PTFId:      tt.ptfID,
				TestID:     tt.testID,
				CustomerID: tt.customerID,
			}

			err := updateParameters(&flags, "")

			if tt.expectError {
				if err == nil {
					t.Errorf("updateParameters() expected an error, but got nil")
				} else if !strings.Contains(err.Error(), tt.errorMessage) {
					// Check if the error message contains the expected string part
					t.Errorf("updateParameters() got error: %q, want error containing: %q", err.Error(), tt.errorMessage)
				}
			} else {
				if err != nil {
					t.Errorf("updateParameters() got an unexpected error: %s", err)
				}
			}
		})
	}
}

var allServices = []string{
	mgrpxy_podman.ServiceHTTPd,
	mgrpxy_podman.ServiceSSH,
	mgrpxy_podman.ServiceTFTFd,
	mgrpxy_podman.ServiceSaltBroker,
	mgrpxy_podman.ServiceSquid,
}
var existingImages = map[string]string{
	mgrpxy_podman.ServiceHTTPd:      "registry.suse.com/suse/multi-linux-manager/5.1/x86_64/httpd:old",
	mgrpxy_podman.ServiceSSH:        "registry.suse.com/suse/multi-linux-manager/5.1/x86_64/ssh:old",
	mgrpxy_podman.ServiceTFTFd:      "registry.suse.com/suse/multi-linux-manager/5.1/x86_64/tftpd:old",
	mgrpxy_podman.ServiceSaltBroker: "registry.suse.com/suse/multi-linux-manager/5.1/x86_64/salt-broker:old",
	mgrpxy_podman.ServiceSquid:      "registry.suse.com/suse/multi-linux-manager/5.1/x86_64/squid:old",
}

// Helper function to create an ImageFlags map for easy iteration.
func imageFlagsMap(flags *podmanPTFFlags) map[string]*types.ImageFlags {
	return map[string]*types.ImageFlags{
		mgrpxy_podman.ServiceHTTPd:      &flags.UpgradeFlags.Httpd,
		mgrpxy_podman.ServiceSSH:        &flags.UpgradeFlags.SSH,
		mgrpxy_podman.ServiceTFTFd:      &flags.UpgradeFlags.Tftpd,
		mgrpxy_podman.ServiceSaltBroker: &flags.UpgradeFlags.SaltBroker,
		mgrpxy_podman.ServiceSquid:      &flags.UpgradeFlags.Squid,
	}
}

func TestUpdateParametersImageLogicOverrideOneImage(t *testing.T) {
	// Save original functions for cleanup
	originalGetServiceImage := getServiceImage
	originalHasRemoteImage := hasRemoteImage

	// Restore original functions after test completion
	defer func() {
		getServiceImage = originalGetServiceImage
		hasRemoteImage = originalHasRemoteImage
	}()

	getServiceImage = func(service string) string {
		return existingImages[service]
	}

	// Iterate through each service to test the scenario where *only* that service is remote
	for _, service := range allServices {
		t.Run("Override_Only_"+service, func(t *testing.T) {
			// Mock hasRemoteImage to return true only for the image of the current service
			hasRemoteImage = func(image string, _ string) bool {
				return strings.Contains(image, strings.ReplaceAll(service, "uyuni-proxy-", ""))
			}

			// Initialize the flags structure
			flags := podmanPTFFlags{
				UpgradeFlags: mgrpxy_podman.PodmanProxyFlags{
					ProxyImageFlags: utils.ProxyImageFlags{
						Registry: types.Registry{
							Host: "registry.suse.com",
						},
					},
				},
				PTFId:      "ptf999",
				TestID:     "",
				CustomerID: "sccuser",
			}

			err := updateParameters(&flags, "")

			if err != nil {
				t.Fatalf("updateParameters() failed with unexpected error: %s", err)
			}

			for checkService, imageFlag := range imageFlagsMap(&flags) {
				if checkService == service {
					// The image that SHOULD be overridden: must be set to the PTF image.
					actual := imageFlag.Name
					testutils.AssertTrue(t, fmt.Sprintf("Expected the PTF image for service %s, got %s", service, actual),
						strings.Contains(actual, "ptf999"))
				} else {
					actual := imageFlag.Name
					testutils.AssertEquals(
						t, fmt.Sprintf("Image shouldn't have changed as the remote is not available: %s", checkService),
						"", actual,
					)
				}
			}
		})
	}
}

func TestUpdateParametersImageLogicNoRemoteImage(t *testing.T) {
	t.Run("No_Image_Is_Remote", func(t *testing.T) {
		// Save original functions for cleanup
		originalGetServiceImage := getServiceImage
		originalHasRemoteImage := hasRemoteImage

		// Restore original functions after test completion
		defer func() {
			getServiceImage = originalGetServiceImage
			hasRemoteImage = originalHasRemoteImage
		}()

		// Mock getServiceImage to return the existing image name for every service
		getServiceImage = func(service string) string {
			return existingImages[service]
		}

		// Mock hasRemoteImage to always return false
		hasRemoteImage = func(_ string, _ string) bool {
			return false
		}

		// Initialize the flags structure
		flags := podmanPTFFlags{
			UpgradeFlags: mgrpxy_podman.PodmanProxyFlags{
				ProxyImageFlags: utils.ProxyImageFlags{
					Registry: types.Registry{
						Host: "registry.suse.com",
					},
				},
			},
			PTFId:      "ptf999",
			TestID:     "",
			CustomerID: "sccuser",
		}

		err := updateParameters(&flags, "")

		if err != nil {
			t.Fatalf("updateParameters() failed with unexpected error: %s", err)
		}

		for checkService, imageFlag := range imageFlagsMap(&flags) {
			actual := imageFlag.Name

			testutils.AssertEquals(
				t, fmt.Sprintf("Image shouldn't have changed as the remote is not available: %s", checkService),
				"", actual,
			)
		}
	})
}
