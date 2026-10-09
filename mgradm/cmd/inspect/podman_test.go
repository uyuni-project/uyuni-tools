// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestPodmanInspectServerNotInitialized(t *testing.T) {
	originalInspectHost := inspectHost
	originalPodmanLogin := podmanLogin
	originalPrepareImages := prepareImages
	originalInspectImages := inspectImages
	defer func() {
		inspectHost = originalInspectHost
		podmanLogin = originalPodmanLogin
		prepareImages = originalPrepareImages
		inspectImages = originalInspectImages
	}()

	loginCalled := false
	prepareCalled := false
	inspectCalled := false

	inspectHost = func() (*podman.HostInspectData, error) {
		return &podman.HostInspectData{HasUyuniServer: false}, nil
	}
	podmanLogin = func(
		_ *podman.HostInspectData,
		_ types.Registry,
		_ types.SCCCredentials,
	) (string, func(), error) {
		loginCalled = true
		return "", func() {}, nil
	}
	prepareImages = func(_ string, _ types.ImageFlags, _ types.PgsqlFlags) (string, string, error) {
		prepareCalled = true
		return "", "", nil
	}
	inspectImages = func(_ string, _ string) (*utils.InspectData, error) {
		inspectCalled = true
		return nil, nil
	}

	err := podmanInspect(&types.GlobalFlags{}, &inspectFlags{}, nil, nil)

	testutils.AssertError(t, "", err)
	testutils.AssertEquals(t, "wrong error", "server is not initialized.", err.Error())
	testutils.AssertTrue(t, "podman login should not be called", !loginCalled)
	testutils.AssertTrue(t, "image preparation should not be called", !prepareCalled)
	testutils.AssertTrue(t, "inspect should not be called", !inspectCalled)
}
