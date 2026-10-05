// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"bytes"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	sharedPodman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestDisableServicesDisablesAllEnabledServicesWithoutStopping(t *testing.T) {
	driver := &testutils.FakeSystemdDriver{
		Installed: []string{
			sharedPodman.ServerService,
			sharedPodman.DBService,
			sharedPodman.TFTPService,
			sharedPodman.ServerAttestationService + "@",
			sharedPodman.HubXmlrpcService + "@",
			sharedPodman.SalineService + "@",
		},
		Enabled: []string{
			sharedPodman.ServerService,
			sharedPodman.DBService,
			sharedPodman.TFTPService,
			sharedPodman.ServerAttestationService + "@0",
			sharedPodman.ServerAttestationService + "@1",
			sharedPodman.HubXmlrpcService + "@0",
			sharedPodman.SalineService + "@0",
		},
		Running: []string{
			sharedPodman.ServerService,
			sharedPodman.DBService,
			sharedPodman.TFTPService,
			sharedPodman.ServerAttestationService + "@0",
			sharedPodman.ServerAttestationService + "@1",
			sharedPodman.HubXmlrpcService + "@0",
			sharedPodman.SalineService + "@0",
		},
	}
	originalRunning := append([]string{}, driver.Running...)

	err := DisableServices(sharedPodman.NewSystemdWithDriver(driver))
	testutils.AssertNoError(t, "disable failed: ", err)
	testutils.AssertEquals(t, "services remain enabled", []string{}, driver.Enabled)
	testutils.AssertEquals(t, "disable changed the running services", originalRunning, driver.Running)
}

func TestEnableServicesEnablesOnlyCoreServicesWithoutStarting(t *testing.T) {
	driver := &testutils.FakeSystemdDriver{
		Installed: []string{
			sharedPodman.ServerService,
			sharedPodman.DBService,
			sharedPodman.TFTPService,
			sharedPodman.ServerAttestationService + "@",
		},
	}
	systemd := sharedPodman.NewSystemdWithDriver(driver)

	testutils.AssertNoError(t, "enable failed: ", EnableServices(systemd))
	testutils.AssertContains(t, "server was not enabled", driver.Enabled, sharedPodman.ServerService)
	testutils.AssertContains(t, "database was not enabled", driver.Enabled, sharedPodman.DBService)
	testutils.AssertNotContains(t, "TFTP was unexpectedly enabled", driver.Enabled, sharedPodman.TFTPService)
	testutils.AssertNotContains(t, "replica was unexpectedly enabled", driver.Enabled,
		sharedPodman.ServerAttestationService+"@0")
	testutils.AssertEquals(t, "enable started services", 0, len(driver.Running))
	testutils.AssertNoError(t, "second enable failed: ", EnableServices(systemd))
}

func TestEnableDisableServicesRequireInstalledServer(t *testing.T) {
	systemd := sharedPodman.NewSystemdWithDriver(&testutils.FakeSystemdDriver{})
	testutils.AssertError(t, "no installed server detected", EnableServices(systemd))
	testutils.AssertError(t, "no installed server detected", DisableServices(systemd))
}

func TestDisableServicesReturnsErrors(t *testing.T) {
	disableErr := errors.New("cannot disable database")
	driver := &testutils.FakeSystemdDriver{
		Installed:            []string{sharedPodman.ServerService, sharedPodman.DBService},
		Enabled:              []string{sharedPodman.ServerService, sharedPodman.DBService},
		DisableServiceErrors: map[string]error{sharedPodman.DBService: disableErr},
	}

	err := DisableServices(sharedPodman.NewSystemdWithDriver(driver))
	testutils.AssertError(t, disableErr.Error(), err)
	testutils.AssertNotContains(t, "server was not disabled", driver.Enabled, sharedPodman.ServerService)
	testutils.AssertContains(t, "database should remain enabled after its error", driver.Enabled, sharedPodman.DBService)
}

func TestWarnIfServicesDisabled(t *testing.T) {
	driver := &testutils.FakeSystemdDriver{
		Installed: []string{sharedPodman.ServerService},
	}
	var output bytes.Buffer
	oldLogger := log.Logger
	log.Logger = zerolog.New(&output)
	t.Cleanup(func() {
		log.Logger = oldLogger
	})

	WarnIfServicesDisabled(sharedPodman.NewSystemdWithDriver(driver))
	testutils.AssertStringContains(t, "disabled warning missing: ", output.String(),
		"Server services are disabled and will not start automatically at boot")
}
